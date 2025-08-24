package encryption

import (
	"context"
	"fmt"

	"github.com/gentra/decorator-arch-go/internal/encryption"
	"github.com/gentra/decorator-arch-go/internal/user"
)

// service implements user.Service with encryption capabilities
// This decorator wraps another user.Service and encrypts/decrypts sensitive data
type service struct {
	next              user.Service
	encryptionService encryption.Service
}

// NewService creates a new encryption decorator for user service
func NewService(next user.Service, encryptionService encryption.Service) user.Service {
	return &service{
		next:              next,
		encryptionService: encryptionService,
	}
}

// Register creates a new user with sensitive data encryption
func (s *service) Register(ctx context.Context, data user.RegisterData) (*user.User, error) {
	// Encrypt sensitive fields before storing
	if data.Email != "" {
		encryptedEmail, err := s.encryptionService.EncryptWithPurpose(ctx, data.Email, encryption.PurposeUserEmail)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt email: %w", err)
		}
		data.Email = encryptedEmail
	}

	if data.FirstName != "" {
		encryptedFirstName, err := s.encryptionService.EncryptWithPurpose(ctx, data.FirstName, encryption.PurposeUserName)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt first name: %w", err)
		}
		data.FirstName = encryptedFirstName
	}

	if data.LastName != "" {
		encryptedLastName, err := s.encryptionService.EncryptWithPurpose(ctx, data.LastName, encryption.PurposeUserName)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt last name: %w", err)
		}
		data.LastName = encryptedLastName
	}

	// Call next service with encrypted data
	result, err := s.next.Register(ctx, data)
	if err != nil {
		return nil, err
	}

	// Decrypt sensitive fields after retrieval
	if result.Email != "" {
		decryptedEmail, err := s.encryptionService.DecryptWithPurpose(ctx, result.Email, encryption.PurposeUserEmail)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt email: %w", err)
		}
		result.Email = decryptedEmail
	}

	if result.FirstName != "" {
		decryptedFirstName, err := s.encryptionService.DecryptWithPurpose(ctx, result.FirstName, encryption.PurposeUserName)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt first name: %w", err)
		}
		result.FirstName = decryptedFirstName
	}

	if result.LastName != "" {
		decryptedLastName, err := s.encryptionService.DecryptWithPurpose(ctx, result.LastName, encryption.PurposeUserName)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt last name: %w", err)
		}
		result.LastName = decryptedLastName
	}

	return result, nil
}

// Login authenticates a user (encrypt email for lookup)
func (s *service) Login(ctx context.Context, email, password string) (*user.AuthResult, error) {
	// Encrypt email for lookup in the database
	encryptedEmail, err := s.encryptionService.EncryptWithPurpose(ctx, email, encryption.PurposeUserEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt email for login: %w", err)
	}

	// Call next service with encrypted email
	result, err := s.next.Login(ctx, encryptedEmail, password)
	if err != nil {
		return nil, err
	}

	// Decrypt user data in the result if present
	if result.User != nil {
		if result.User.Email != "" {
			decryptedEmail, err := s.encryptionService.DecryptWithPurpose(ctx, result.User.Email, encryption.PurposeUserEmail)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt email: %w", err)
			}
			result.User.Email = decryptedEmail
		}

		if result.User.FirstName != "" {
			decryptedFirstName, err := s.encryptionService.DecryptWithPurpose(ctx, result.User.FirstName, encryption.PurposeUserName)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt first name: %w", err)
			}
			result.User.FirstName = decryptedFirstName
		}

		if result.User.LastName != "" {
			decryptedLastName, err := s.encryptionService.DecryptWithPurpose(ctx, result.User.LastName, encryption.PurposeUserName)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt last name: %w", err)
			}
			result.User.LastName = decryptedLastName
		}
	}

	return result, nil
}

// GetByID retrieves a user by ID and decrypts sensitive data
func (s *service) GetByID(ctx context.Context, id string) (*user.User, error) {
	// Call next service
	result, err := s.next.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	// Decrypt sensitive fields after retrieval
	if result.Email != "" {
		decryptedEmail, err := s.encryptionService.DecryptWithPurpose(ctx, result.Email, encryption.PurposeUserEmail)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt email: %w", err)
		}
		result.Email = decryptedEmail
	}

	if result.FirstName != "" {
		decryptedFirstName, err := s.encryptionService.DecryptWithPurpose(ctx, result.FirstName, encryption.PurposeUserName)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt first name: %w", err)
		}
		result.FirstName = decryptedFirstName
	}

	if result.LastName != "" {
		decryptedLastName, err := s.encryptionService.DecryptWithPurpose(ctx, result.LastName, encryption.PurposeUserName)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt last name: %w", err)
		}
		result.LastName = decryptedLastName
	}

	return result, nil
}

// UpdateProfile updates user profile with encryption
func (s *service) UpdateProfile(ctx context.Context, id string, data user.UpdateProfileData) (*user.User, error) {
	// Encrypt sensitive fields before updating
	if data.Email != nil && *data.Email != "" {
		encryptedEmail, err := s.encryptionService.EncryptWithPurpose(ctx, *data.Email, encryption.PurposeUserEmail)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt email: %w", err)
		}
		data.Email = &encryptedEmail
	}

	if data.FirstName != nil && *data.FirstName != "" {
		encryptedFirstName, err := s.encryptionService.EncryptWithPurpose(ctx, *data.FirstName, encryption.PurposeUserName)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt first name: %w", err)
		}
		data.FirstName = &encryptedFirstName
	}

	if data.LastName != nil && *data.LastName != "" {
		encryptedLastName, err := s.encryptionService.EncryptWithPurpose(ctx, *data.LastName, encryption.PurposeUserName)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt last name: %w", err)
		}
		data.LastName = &encryptedLastName
	}

	// Call next service with encrypted data
	result, err := s.next.UpdateProfile(ctx, id, data)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	// Decrypt sensitive fields after retrieval
	if result.Email != "" {
		decryptedEmail, err := s.encryptionService.DecryptWithPurpose(ctx, result.Email, encryption.PurposeUserEmail)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt email: %w", err)
		}
		result.Email = decryptedEmail
	}

	if result.FirstName != "" {
		decryptedFirstName, err := s.encryptionService.DecryptWithPurpose(ctx, result.FirstName, encryption.PurposeUserName)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt first name: %w", err)
		}
		result.FirstName = decryptedFirstName
	}

	if result.LastName != "" {
		decryptedLastName, err := s.encryptionService.DecryptWithPurpose(ctx, result.LastName, encryption.PurposeUserName)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt last name: %w", err)
		}
		result.LastName = decryptedLastName
	}

	return result, nil
}

// GetPreferences retrieves user preferences (no encryption needed for preferences)
func (s *service) GetPreferences(ctx context.Context, userID string) (*user.UserPreferences, error) {
	// Preferences don't contain sensitive data that needs encryption
	// Just pass through to next service
	return s.next.GetPreferences(ctx, userID)
}

// UpdatePreferences updates user preferences (no encryption needed for preferences)
func (s *service) UpdatePreferences(ctx context.Context, userID string, prefs user.UserPreferences) error {
	// Preferences don't contain sensitive data that needs encryption
	// Just pass through to next service
	return s.next.UpdatePreferences(ctx, userID, prefs)
}
