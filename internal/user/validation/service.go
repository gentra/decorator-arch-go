package validation

import (
	"context"

	"github.com/gentra/decorator-arch-go/internal/user"
	"github.com/gentra/decorator-arch-go/internal/validation"
)

// service implements user.Service with validation capabilities
type service struct {
	next              user.Service
	validationService validation.Service
}

// NewService creates a new validation-enabled user service
func NewService(next user.Service, validationService validation.Service) user.Service {
	return &service{
		next:              next,
		validationService: validationService,
	}
}

// Register validates registration data before creating a user
func (s *service) Register(ctx context.Context, data user.RegisterData) (*user.User, error) {
	// Validate registration data using the validation domain service
	if err := s.validationService.ValidateUserRegistration(ctx, data); err != nil {
		return nil, err
	}

	// Call next service if validation passes
	return s.next.Register(ctx, data)
}

// Login validates login credentials before authentication
func (s *service) Login(ctx context.Context, email, password string) (*user.AuthResult, error) {
	// Validate email format
	if err := s.validationService.ValidateEmail(ctx, email); err != nil {
		return nil, err
	}

	// Validate password
	if err := s.validationService.ValidatePassword(ctx, password); err != nil {
		return nil, err
	}

	// Call next service if validation passes
	return s.next.Login(ctx, email, password)
}

// GetByID validates the user ID before retrieval
func (s *service) GetByID(ctx context.Context, id string) (*user.User, error) {
	// Validate user ID format
	if err := s.validationService.ValidateUserID(ctx, id); err != nil {
		return nil, err
	}

	// Call next service if validation passes
	return s.next.GetByID(ctx, id)
}

// UpdateProfile validates profile update data before updating
func (s *service) UpdateProfile(ctx context.Context, id string, data user.UpdateProfileData) (*user.User, error) {
	// Validate user ID
	if err := s.validationService.ValidateUserID(ctx, id); err != nil {
		return nil, err
	}

	// Validate update data
	if err := s.validationService.ValidateUserUpdate(ctx, data); err != nil {
		return nil, err
	}

	// Call next service if validation passes
	return s.next.UpdateProfile(ctx, id, data)
}

// GetPreferences validates user ID before retrieving preferences
func (s *service) GetPreferences(ctx context.Context, userID string) (*user.UserPreferences, error) {
	// Validate user ID
	if err := s.validationService.ValidateUserID(ctx, userID); err != nil {
		return nil, err
	}

	// Call next service if validation passes
	return s.next.GetPreferences(ctx, userID)
}

// UpdatePreferences validates data before updating preferences
func (s *service) UpdatePreferences(ctx context.Context, userID string, prefs user.UserPreferences) error {
	// Validate user ID
	if err := s.validationService.ValidateUserID(ctx, userID); err != nil {
		return err
	}

	// Validate preferences data
	if err := s.validationService.ValidateUserPreferences(ctx, prefs); err != nil {
		return err
	}

	// Call next service if validation passes
	return s.next.UpdatePreferences(ctx, userID, prefs)
}
