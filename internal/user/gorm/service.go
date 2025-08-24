package gorm

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/gentra/decorator-arch-go/internal/user"
)

// service implements the user.Service interface using GORM
type service struct {
	db *gorm.DB
}

// NewService creates a new GORM-based user service
func NewService(db *gorm.DB) user.Service {
	return &service{
		db: db,
	}
}

// Register creates a new user in the database
func (s *service) Register(ctx context.Context, data user.RegisterData) (*user.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user model
	userModel := UserModel{
		Email:        data.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    data.FirstName,
		LastName:     data.LastName,
	}

	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create user
	if err := tx.Create(&userModel).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, user.ErrEmailAlreadyExists
		}
		return nil, err
	}

	// Create default preferences for the user
	defaultPrefs := user.DefaultUserPreferences(userModel.ID)
	notificationTypesJSON, err := json.Marshal(defaultPrefs.NotificationTypes)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	prefsModel := UserPreferencesModel{
		UserID:             userModel.ID,
		EmailNotifications: defaultPrefs.EmailNotifications,
		PushNotifications:  defaultPrefs.PushNotifications,
		SMSNotifications:   defaultPrefs.SMSNotifications,
		Theme:              defaultPrefs.Theme,
		Language:           defaultPrefs.Language,
		Timezone:           defaultPrefs.Timezone,
		NotificationTypes:  notificationTypesJSON,
	}

	if err := tx.Create(&prefsModel).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Convert to domain model
	return s.toDomainUser(&userModel), nil
}

// Login authenticates a user and returns auth result
func (s *service) Login(ctx context.Context, email, password string) (*user.AuthResult, error) {
	var userModel UserModel

	// Find user by email
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(userModel.PasswordHash), []byte(password)); err != nil {
		return nil, user.ErrInvalidCredentials
	}

	// Convert to domain model
	domainUser := s.toDomainUser(&userModel)

	// Create auth result (token generation would be handled in a higher layer)
	authResult := &user.AuthResult{
		User: domainUser,
		// Token and ExpiresAt would be set by authentication service in a higher layer
	}

	return authResult, nil
}

// GetByID retrieves a user by ID
func (s *service) GetByID(ctx context.Context, id string) (*user.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, user.ErrUserNotFound
	}

	var userModel UserModel
	if err := s.db.WithContext(ctx).Where("id = ?", userID).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return s.toDomainUser(&userModel), nil
}

// UpdateProfile updates user profile information
func (s *service) UpdateProfile(ctx context.Context, id string, data user.UpdateProfileData) (*user.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, user.ErrUserNotFound
	}

	// Build update map
	updates := make(map[string]interface{})
	if data.FirstName != nil {
		updates["first_name"] = *data.FirstName
	}
	if data.LastName != nil {
		updates["last_name"] = *data.LastName
	}
	if data.Email != nil {
		updates["email"] = *data.Email
	}

	if len(updates) == 0 {
		// No updates to make, just return the existing user
		return s.GetByID(ctx, id)
	}

	// Update user
	if err := s.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) && data.Email != nil {
			return nil, user.ErrEmailAlreadyExists
		}
		return nil, err
	}

	// Return updated user
	return s.GetByID(ctx, id)
}

// GetPreferences retrieves user preferences
func (s *service) GetPreferences(ctx context.Context, userID string) (*user.UserPreferences, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, user.ErrUserNotFound
	}

	var prefsModel UserPreferencesModel
	if err := s.db.WithContext(ctx).Where("user_id = ?", parsedUserID).First(&prefsModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrPreferencesNotFound
		}
		return nil, err
	}

	return s.toDomainPreferences(&prefsModel)
}

// UpdatePreferences updates user preferences
func (s *service) UpdatePreferences(ctx context.Context, userID string, prefs user.UserPreferences) error {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return user.ErrUserNotFound
	}

	// Convert notification types to JSON
	notificationTypesJSON, err := json.Marshal(prefs.NotificationTypes)
	if err != nil {
		return err
	}

	// Update preferences
	updates := map[string]interface{}{
		"email_notifications": prefs.EmailNotifications,
		"push_notifications":  prefs.PushNotifications,
		"sms_notifications":   prefs.SMSNotifications,
		"theme":               prefs.Theme,
		"language":            prefs.Language,
		"timezone":            prefs.Timezone,
		"notification_types":  notificationTypesJSON,
	}

	if err := s.db.WithContext(ctx).Model(&UserPreferencesModel{}).Where("user_id = ?", parsedUserID).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// Helper methods for converting between GORM models and domain models
func (s *service) toDomainUser(model *UserModel) *user.User {
	return &user.User{
		ID:           model.ID,
		Email:        model.Email,
		PasswordHash: model.PasswordHash,
		FirstName:    model.FirstName,
		LastName:     model.LastName,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}

func (s *service) toDomainPreferences(model *UserPreferencesModel) (*user.UserPreferences, error) {
	var notificationTypes map[string]bool
	if err := json.Unmarshal(model.NotificationTypes, &notificationTypes); err != nil {
		return nil, err
	}

	return &user.UserPreferences{
		ID:                 model.ID,
		UserID:             model.UserID,
		EmailNotifications: model.EmailNotifications,
		PushNotifications:  model.PushNotifications,
		SMSNotifications:   model.SMSNotifications,
		Theme:              model.Theme,
		Language:           model.Language,
		Timezone:           model.Timezone,
		NotificationTypes:  notificationTypes,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
	}, nil
}
