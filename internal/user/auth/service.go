package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/gentra/decorator-arch-go/internal/auth"
	"github.com/gentra/decorator-arch-go/internal/user"
)

// service implements user.Service interface but delegates authentication to auth domain
type service struct {
	next        user.Service
	authService auth.Service
}

// NewService creates a new user auth service that wraps user operations with auth capabilities
func NewService(next user.Service, authService auth.Service) user.Service {
	return &service{
		next:        next,
		authService: authService,
	}
}

// Register creates a new user (delegates to next service)
func (s *service) Register(ctx context.Context, data user.RegisterData) (*user.User, error) {
	return s.next.Register(ctx, data)
}

// Login authenticates a user using the auth domain with basic strategy
func (s *service) Login(ctx context.Context, email, password string) (*user.AuthResult, error) {
	// Use auth domain for authentication
	authCredentials := auth.BasicCredentials{
		Email:    email,
		Password: password,
	}

	authResult, err := s.authService.Authenticate(ctx, "basic", authCredentials)
	if err != nil {
		// Convert auth domain errors back to user domain errors
		if err == auth.ErrInvalidCredentials {
			return nil, user.ErrInvalidCredentials
		}
		if err == auth.ErrUserNotFound {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	// Convert auth domain result to user domain result
	userAuthResult := &user.AuthResult{
		User:         s.convertAuthUserToUserDomain(authResult.User),
		Token:        authResult.Token,
		RefreshToken: authResult.RefreshToken,
		ExpiresAt:    authResult.ExpiresAt,
	}

	return userAuthResult, nil
}

// GetByID retrieves a user by ID (delegates to next service)
func (s *service) GetByID(ctx context.Context, id string) (*user.User, error) {
	return s.next.GetByID(ctx, id)
}

// UpdateProfile updates user profile (delegates to next service)
func (s *service) UpdateProfile(ctx context.Context, id string, data user.UpdateProfileData) (*user.User, error) {
	return s.next.UpdateProfile(ctx, id, data)
}

// GetPreferences retrieves user preferences (delegates to next service)
func (s *service) GetPreferences(ctx context.Context, userID string) (*user.UserPreferences, error) {
	return s.next.GetPreferences(ctx, userID)
}

// UpdatePreferences updates user preferences (delegates to next service)
func (s *service) UpdatePreferences(ctx context.Context, userID string, prefs user.UserPreferences) error {
	return s.next.UpdatePreferences(ctx, userID, prefs)
}

// This auth adapter only implements user.Service interface
// All authentication logic is handled by the auth domain service internally

// convertAuthUserToUserDomain converts auth domain user to user domain user
func (s *service) convertAuthUserToUserDomain(authUser *auth.User) *user.User {
	if authUser == nil {
		return nil
	}

	// Parse UUID from string
	userID, err := uuid.Parse(authUser.ID)
	if err != nil {
		// Fallback to nil UUID if parsing fails
		userID = uuid.Nil
	}

	return &user.User{
		ID:           userID,
		Email:        authUser.Email,
		FirstName:    authUser.FirstName,
		LastName:     authUser.LastName,
		PasswordHash: authUser.PasswordHash,
		CreatedAt:    authUser.CreatedAt,
		UpdatedAt:    authUser.UpdatedAt,
	}
}

// No additional interfaces defined here - only implements user.Service
// Following the architectural rule: "Each implementation folder should only implement the Service interface"
