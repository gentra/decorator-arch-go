package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/gentra/decorator-arch-go/internal/events"
	"github.com/gentra/decorator-arch-go/internal/notification"
	"github.com/gentra/decorator-arch-go/internal/token"
	"github.com/gentra/decorator-arch-go/internal/user"
)

// Dependencies defines external services that the usecase layer depends on
type Dependencies struct {
	NotificationService notification.Service
	TokenService        token.Service
	EventPublisher      events.Service
}

// service implements the user.Service interface with business logic
type service struct {
	next user.Service
	deps Dependencies
}

// NewService creates a new usecase service with business logic
func NewService(next user.Service, deps Dependencies) user.Service {
	return &service{
		next: next,
		deps: deps,
	}
}

// Register creates a new user with business logic and orchestration
func (s *service) Register(ctx context.Context, data user.RegisterData) (*user.User, error) {
	// Call next service to create the user
	result, err := s.next.Register(ctx, data)
	if err != nil {
		return nil, err
	}

	// Business logic: Send welcome email (non-blocking)
	go func() {
		// Use a background context to avoid cancellation affecting this operation
		backgroundCtx := context.Background()
		if err := s.deps.NotificationService.SendWelcomeEmail(
			backgroundCtx,
			result.Email,
			result.GetFullName(),
		); err != nil {
			// Log error but don't fail the registration
			log.Printf("Failed to send welcome email to %s: %v", result.Email, err)
		}
	}()

	// Publish user registered event using events domain service
	event := events.Event{
		Type:          events.EventTypeUserRegistered,
		AggregateID:   result.ID.String(),
		AggregateType: "user",
		Data: map[string]interface{}{
			"user_id":       result.ID.String(),
			"email":         result.Email,
			"first_name":    result.FirstName,
			"last_name":     result.LastName,
			"registered_at": result.CreatedAt,
		},
	}

	if err := s.deps.EventPublisher.Publish(ctx, event); err != nil {
		// Log event publishing failure but don't fail the operation
		log.Printf("Failed to publish UserRegistered event: %v", err)
	}

	return result, nil
}

// Login authenticates a user with business logic and token generation
func (s *service) Login(ctx context.Context, email, password string) (*user.AuthResult, error) {
	// Call next service to authenticate
	result, err := s.next.Login(ctx, email, password)
	if err != nil {
		return nil, err
	}

	// Business logic: Generate tokens
	token, expiresAt, err := s.deps.TokenService.GenerateAuthToken(
		ctx,
		result.User.ID.String(),
		result.User.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth token: %w", err)
	}

	refreshToken, err := s.deps.TokenService.GenerateRefreshToken(
		ctx,
		result.User.ID.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Update auth result with tokens
	result.Token = token
	result.RefreshToken = refreshToken
	result.ExpiresAt = expiresAt

	// Publish login event using events domain service
	loginEvent := events.Event{
		Type:          events.EventTypeUserLoggedIn,
		AggregateID:   result.User.ID.String(),
		AggregateType: "user",
		Data: map[string]interface{}{
			"user_id":  result.User.ID.String(),
			"email":    result.User.Email,
			"login_at": time.Now(),
		},
	}

	if err := s.deps.EventPublisher.Publish(ctx, loginEvent); err != nil {
		log.Printf("Failed to publish UserLoggedIn event: %v", err)
	}

	return result, nil
}

// GetByID retrieves a user by ID (no additional business logic needed)
func (s *service) GetByID(ctx context.Context, id string) (*user.User, error) {
	return s.next.GetByID(ctx, id)
}

// UpdateProfile updates user profile with business logic
func (s *service) UpdateProfile(ctx context.Context, id string, data user.UpdateProfileData) (*user.User, error) {
	// Get current user data for comparison
	currentUser, err := s.next.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Call next service to update profile
	result, err := s.next.UpdateProfile(ctx, id, data)
	if err != nil {
		return nil, err
	}

	// Business logic: Determine what changed and send notifications
	changes := s.detectProfileChanges(currentUser, result, data)

	if len(changes) > 0 {
		// Send notification about profile changes (non-blocking)
		go func() {
			backgroundCtx := context.Background()
			if err := s.deps.NotificationService.SendProfileUpdateNotification(
				backgroundCtx,
				result.ID.String(),
				changes,
			); err != nil {
				log.Printf("Failed to send profile update notification: %v", err)
			}
		}()

		// Publish profile updated event using events domain service
		updateEvent := events.Event{
			Type:          events.EventTypeUserUpdated,
			AggregateID:   result.ID.String(),
			AggregateType: "user",
			Data: map[string]interface{}{
				"user_id":    result.ID.String(),
				"updated_at": result.UpdatedAt,
				"changes":    changes,
			},
		}

		if err := s.deps.EventPublisher.Publish(ctx, updateEvent); err != nil {
			log.Printf("Failed to publish ProfileUpdated event: %v", err)
		}
	}

	return result, nil
}

// GetPreferences retrieves user preferences with business logic for defaults
func (s *service) GetPreferences(ctx context.Context, userID string) (*user.UserPreferences, error) {
	// Try to get preferences from next service
	prefs, err := s.next.GetPreferences(ctx, userID)
	if err != nil {
		// If preferences don't exist, create default ones
		if err == user.ErrPreferencesNotFound {
			return s.createDefaultPreferencesForUser(ctx, userID)
		}
		return nil, err
	}

	// Business logic: Ensure preferences have all required fields
	return s.ensureCompletePreferences(prefs), nil
}

// UpdatePreferences updates user preferences with business logic
func (s *service) UpdatePreferences(ctx context.Context, userID string, prefs user.UserPreferences) error {
	// Get current preferences for comparison
	currentPrefs, _ := s.next.GetPreferences(ctx, userID)

	// Call next service to update preferences
	err := s.next.UpdatePreferences(ctx, userID, prefs)
	if err != nil {
		return err
	}

	// Business logic: Determine what changed
	if currentPrefs != nil {
		changes := s.detectPreferencesChanges(currentPrefs, &prefs)

		if len(changes) > 0 {
			// Publish preferences updated event using events domain service
			prefsEvent := events.Event{
				Type:          events.EventTypeUserPrefsUpdated,
				AggregateID:   userID,
				AggregateType: "user",
				Data: map[string]interface{}{
					"user_id":    userID,
					"updated_at": time.Now(),
					"preferences": map[string]interface{}{
						"theme":               prefs.Theme,
						"language":            prefs.Language,
						"timezone":            prefs.Timezone,
						"email_notifications": prefs.EmailNotifications,
						"push_notifications":  prefs.PushNotifications,
						"sms_notifications":   prefs.SMSNotifications,
						"notification_types":  prefs.NotificationTypes,
					},
				},
			}

			if err := s.deps.EventPublisher.Publish(ctx, prefsEvent); err != nil {
				log.Printf("Failed to publish PreferencesUpdated event: %v", err)
			}
		}
	}

	return nil
}

// Helper methods for business logic

func (s *service) detectProfileChanges(current, updated *user.User, data user.UpdateProfileData) map[string]interface{} {
	changes := make(map[string]interface{})

	if data.Email != nil && current.Email != updated.Email {
		changes["email"] = map[string]string{
			"old": current.Email,
			"new": updated.Email,
		}
	}

	if data.FirstName != nil && current.FirstName != updated.FirstName {
		changes["first_name"] = map[string]string{
			"old": current.FirstName,
			"new": updated.FirstName,
		}
	}

	if data.LastName != nil && current.LastName != updated.LastName {
		changes["last_name"] = map[string]string{
			"old": current.LastName,
			"new": updated.LastName,
		}
	}

	return changes
}

func (s *service) detectPreferencesChanges(current, updated *user.UserPreferences) map[string]interface{} {
	changes := make(map[string]interface{})

	if current.EmailNotifications != updated.EmailNotifications {
		changes["email_notifications"] = updated.EmailNotifications
	}

	if current.PushNotifications != updated.PushNotifications {
		changes["push_notifications"] = updated.PushNotifications
	}

	if current.SMSNotifications != updated.SMSNotifications {
		changes["sms_notifications"] = updated.SMSNotifications
	}

	if current.Theme != updated.Theme {
		changes["theme"] = updated.Theme
	}

	if current.Language != updated.Language {
		changes["language"] = updated.Language
	}

	if current.Timezone != updated.Timezone {
		changes["timezone"] = updated.Timezone
	}

	// Compare notification types
	if !equalNotificationTypeMaps(current.NotificationTypes, updated.NotificationTypes) {
		changes["notification_types"] = updated.NotificationTypes
	}

	return changes
}

func (s *service) createDefaultPreferencesForUser(ctx context.Context, userID string) (*user.UserPreferences, error) {
	// Parse user ID
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	// Create default preferences
	defaultPrefs := user.DefaultUserPreferences(parsedUserID)

	// Save the default preferences
	err = s.next.UpdatePreferences(ctx, userID, *defaultPrefs)
	if err != nil {
		return nil, err
	}

	return defaultPrefs, nil
}

func (s *service) ensureCompletePreferences(prefs *user.UserPreferences) *user.UserPreferences {
	// Ensure notification types map is not nil and has all required types
	if prefs.NotificationTypes == nil {
		prefs.NotificationTypes = make(map[string]bool)
	}

	// Add missing notification types with default values
	defaultTypes := map[string]bool{
		"task_assigned":   true,
		"task_due_soon":   true,
		"project_updated": true,
		"project_invite":  true,
		"system_updates":  false,
		"marketing":       false,
	}

	for notificationType, defaultValue := range defaultTypes {
		if _, exists := prefs.NotificationTypes[notificationType]; !exists {
			prefs.NotificationTypes[notificationType] = defaultValue
		}
	}

	// Ensure required fields have default values
	if prefs.Theme == "" {
		prefs.Theme = "light"
	}
	if prefs.Language == "" {
		prefs.Language = "en"
	}
	if prefs.Timezone == "" {
		prefs.Timezone = "UTC"
	}

	return prefs
}

// Utility functions

func equalNotificationTypeMaps(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}

	for key, valueA := range a {
		if valueB, exists := b[key]; !exists || valueA != valueB {
			return false
		}
	}

	return true
}
