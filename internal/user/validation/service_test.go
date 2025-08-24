package validation_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gentra/decorator-arch-go/internal/user"
	usermock "github.com/gentra/decorator-arch-go/internal/user/mock"
	"github.com/gentra/decorator-arch-go/internal/user/validation"
	validationDomain "github.com/gentra/decorator-arch-go/internal/validation"
)

func TestUserValidationService_Register(t *testing.T) {
	tests := []struct {
		name                string
		setupMocks          func(*usermock.MockUserService)
		setupValidator      func(*usermock.MockValidationService)
		registerData        user.RegisterData
		expectedUser        *user.User
		expectedError       error
		expectNextCalled    bool
		expectedFieldErrors []string
	}{
		{
			name: "Given valid registration data, When Register is called, Then should validate and pass to next service",
			setupMocks: func(mockNext *usermock.MockUserService) {
				validData := user.RegisterData{
					Email:     "valid@example.com",
					Password:  "SecurePass123!",
					FirstName: "John",
					LastName:  "Doe",
				}

				createdUser := &user.User{
					ID:        uuid.New(),
					Email:     "valid@example.com",
					FirstName: "John",
					LastName:  "Doe",
				}

				mockNext.On("Register", mock.Anything, validData).Return(createdUser, nil)
			},
			setupValidator: func(mockValidator *usermock.MockValidationService) {
				validData := user.RegisterData{
					Email:     "valid@example.com",
					Password:  "SecurePass123!",
					FirstName: "John",
					LastName:  "Doe",
				}
				mockValidator.On("ValidateUserRegistration", mock.Anything, validData).Return(nil)
			},
			registerData: user.RegisterData{
				Email:     "valid@example.com",
				Password:  "SecurePass123!",
				FirstName: "John",
				LastName:  "Doe",
			},
			expectedUser: &user.User{
				Email:     "valid@example.com",
				FirstName: "John",
				LastName:  "Doe",
			},
			expectedError:       nil,
			expectNextCalled:    true,
			expectedFieldErrors: nil,
		},
		{
			name: "Given invalid email format, When Register is called, Then should return validation error and not call next service",
			setupMocks: func(mockNext *usermock.MockUserService) {
				// Next service should not be called
			},
			setupValidator: func(mockValidator *usermock.MockValidationService) {
				validationError := validationDomain.ValidationErrors{
					Errors: []validationDomain.ValidationError{
						{Field: "email", Message: "must be a valid email address"},
					},
				}
				mockValidator.On("ValidateUserRegistration", mock.Anything, mock.Anything).Return(validationError)
			},
			registerData: user.RegisterData{
				Email:     "invalid-email",
				Password:  "SecurePass123!",
				FirstName: "John",
				LastName:  "Doe",
			},
			expectedUser:        nil,
			expectedError:       validationDomain.ValidationErrors{},
			expectNextCalled:    false,
			expectedFieldErrors: []string{"email"},
		},
		{
			name: "Given weak password, When Register is called, Then should return validation error and not call next service",
			setupMocks: func(mockNext *usermock.MockUserService) {
				// Next service should not be called
			},
			setupValidator: func(mockValidator *usermock.MockValidationService) {
				validationError := validationDomain.ValidationErrors{
					Errors: []validationDomain.ValidationError{
						{Field: "password", Message: "password does not meet security requirements"},
					},
				}
				mockValidator.On("ValidateUserRegistration", mock.Anything, mock.Anything).Return(validationError)
			},
			registerData: user.RegisterData{
				Email:     "valid@example.com",
				Password:  "weak", // Too short, no uppercase, no special chars
				FirstName: "John",
				LastName:  "Doe",
			},
			expectedUser:        nil,
			expectedError:       validationDomain.ValidationErrors{},
			expectNextCalled:    false,
			expectedFieldErrors: []string{"password"},
		},
		{
			name: "Given empty first name, When Register is called, Then should return validation error and not call next service",
			setupMocks: func(mockNext *usermock.MockUserService) {
				// Next service should not be called
			},
			setupValidator: func(mockValidator *usermock.MockValidationService) {
				validationError := validationDomain.ValidationErrors{
					Errors: []validationDomain.ValidationError{
						{Field: "first_name", Message: "field is required"},
					},
				}
				mockValidator.On("ValidateUserRegistration", mock.Anything, mock.Anything).Return(validationError)
			},
			registerData: user.RegisterData{
				Email:     "valid@example.com",
				Password:  "SecurePass123!",
				FirstName: "", // Empty first name
				LastName:  "Doe",
			},
			expectedUser:        nil,
			expectedError:       validationDomain.ValidationErrors{},
			expectNextCalled:    false,
			expectedFieldErrors: []string{"first_name"},
		},
		{
			name: "Given multiple validation errors, When Register is called, Then should return all validation errors and not call next service",
			setupMocks: func(mockNext *usermock.MockUserService) {
				// Next service should not be called
			},
			setupValidator: func(mockValidator *usermock.MockValidationService) {
				validationError := validationDomain.ValidationErrors{
					Errors: []validationDomain.ValidationError{
						{Field: "email", Message: "must be a valid email address"},
						{Field: "password", Message: "password does not meet security requirements"},
						{Field: "first_name", Message: "field is required"},
						{Field: "last_name", Message: "field is required"},
					},
				}
				mockValidator.On("ValidateUserRegistration", mock.Anything, mock.Anything).Return(validationError)
			},
			registerData: user.RegisterData{
				Email:     "invalid-email", // Invalid email
				Password:  "weak",          // Weak password
				FirstName: "J",             // Too short
				LastName:  "",              // Empty
			},
			expectedUser:        nil,
			expectedError:       validationDomain.ValidationErrors{},
			expectNextCalled:    false,
			expectedFieldErrors: []string{"email", "password", "first_name", "last_name"},
		},
		{
			name: "Given valid data but next service fails, When Register is called, Then should validate successfully but return next service error",
			setupMocks: func(mockNext *usermock.MockUserService) {
				validData := user.RegisterData{
					Email:     "existing@example.com",
					Password:  "SecurePass123!",
					FirstName: "John",
					LastName:  "Doe",
				}

				mockNext.On("Register", mock.Anything, validData).Return(nil, user.ErrEmailAlreadyExists)
			},
			setupValidator: func(mockValidator *usermock.MockValidationService) {
				validData := user.RegisterData{
					Email:     "existing@example.com",
					Password:  "SecurePass123!",
					FirstName: "John",
					LastName:  "Doe",
				}
				mockValidator.On("ValidateUserRegistration", mock.Anything, validData).Return(nil)
			},
			registerData: user.RegisterData{
				Email:     "existing@example.com",
				Password:  "SecurePass123!",
				FirstName: "John",
				LastName:  "Doe",
			},
			expectedUser:        nil,
			expectedError:       user.ErrEmailAlreadyExists,
			expectNextCalled:    true,
			expectedFieldErrors: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockNext := new(usermock.MockUserService)
			mockValidator := new(usermock.MockValidationService)
			validationService := validation.NewService(mockNext, mockValidator)

			tt.setupMocks(mockNext)
			if tt.setupValidator != nil {
				tt.setupValidator(mockValidator)
			}

			// Act
			result, err := validationService.Register(context.Background(), tt.registerData)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)

				// Check if it's a validation error with expected fields
				if len(tt.expectedFieldErrors) > 0 {
					var validationErrors validationDomain.ValidationErrors
					assert.ErrorAs(t, err, &validationErrors, "Error should be ValidationErrors type")

					// Verify that all expected field errors are present
					errorFields := make(map[string]bool)
					for _, validationErr := range validationErrors.Errors {
						errorFields[validationErr.Field] = true
					}

					for _, expectedField := range tt.expectedFieldErrors {
						assert.True(t, errorFields[expectedField], "Expected validation error for field: %s", expectedField)
					}
				} else {
					// Check for specific error type
					assert.ErrorIs(t, err, tt.expectedError)
				}

				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedUser.Email, result.Email)
				assert.Equal(t, tt.expectedUser.FirstName, result.FirstName)
				assert.Equal(t, tt.expectedUser.LastName, result.LastName)
			}

			// Verify mock expectations
			if tt.expectNextCalled {
				mockNext.AssertExpectations(t)
			} else {
				mockNext.AssertNotCalled(t, "Register")
			}
		})
	}
}

func TestUserValidationService_GetByID(t *testing.T) {
	tests := []struct {
		name             string
		setupMocks       func(*usermock.MockUserService)
		setupValidator   func(*usermock.MockValidationService)
		userID           string
		expectedUser     *user.User
		expectedError    error
		expectNextCalled bool
	}{
		{
			name: "Given valid UUID, When GetByID is called, Then should validate and pass to next service",
			setupMocks: func(mockNext *usermock.MockUserService) {
				validID := "550e8400-e29b-41d4-a716-446655440000"
				testUser := &user.User{
					ID:        uuid.MustParse(validID),
					Email:     "test@example.com",
					FirstName: "Test",
					LastName:  "User",
				}

				mockNext.On("GetByID", mock.Anything, validID).Return(testUser, nil)
			},
			setupValidator: func(mockValidator *usermock.MockValidationService) {
				validID := "550e8400-e29b-41d4-a716-446655440000"
				mockValidator.On("ValidateUserID", mock.Anything, validID).Return(nil)
			},
			userID: "550e8400-e29b-41d4-a716-446655440000",
			expectedUser: &user.User{
				ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				Email:     "test@example.com",
				FirstName: "Test",
				LastName:  "User",
			},
			expectedError:    nil,
			expectNextCalled: true,
		},
		{
			name: "Given invalid UUID format, When GetByID is called, Then should return validation error and not call next service",
			setupMocks: func(mockNext *usermock.MockUserService) {
				// Next service should not be called
			},
			setupValidator: func(mockValidator *usermock.MockValidationService) {
				validationError := validationDomain.ValidationError{
					Field:   "user_id",
					Message: "must be a valid UUID",
					Value:   "invalid-uuid",
					Rule:    "uuid",
				}
				mockValidator.On("ValidateUserID", mock.Anything, "invalid-uuid").Return(validationError)
			},
			userID:           "invalid-uuid",
			expectedUser:     nil,
			expectedError:    validationDomain.ValidationError{},
			expectNextCalled: false,
		},
		{
			name: "Given empty string as ID, When GetByID is called, Then should return validation error and not call next service",
			setupMocks: func(mockNext *usermock.MockUserService) {
				// Next service should not be called
			},
			setupValidator: func(mockValidator *usermock.MockValidationService) {
				validationError := validationDomain.ValidationError{
					Field:   "user_id",
					Message: "must be a valid UUID",
					Value:   "",
					Rule:    "uuid",
				}
				mockValidator.On("ValidateUserID", mock.Anything, "").Return(validationError)
			},
			userID:           "",
			expectedUser:     nil,
			expectedError:    validationDomain.ValidationError{},
			expectNextCalled: false,
		},
		{
			name: "Given valid UUID but user not found, When GetByID is called, Then should validate successfully but return not found error",
			setupMocks: func(mockNext *usermock.MockUserService) {
				validID := "550e8400-e29b-41d4-a716-446655440999"
				mockNext.On("GetByID", mock.Anything, validID).Return(nil, user.ErrUserNotFound)
			},
			setupValidator: func(mockValidator *usermock.MockValidationService) {
				validID := "550e8400-e29b-41d4-a716-446655440999"
				mockValidator.On("ValidateUserID", mock.Anything, validID).Return(nil)
			},
			userID:           "550e8400-e29b-41d4-a716-446655440999",
			expectedUser:     nil,
			expectedError:    user.ErrUserNotFound,
			expectNextCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockNext := new(usermock.MockUserService)
			mockValidator := new(usermock.MockValidationService)
			validationService := validation.NewService(mockNext, mockValidator)

			tt.setupMocks(mockNext)
			if tt.setupValidator != nil {
				tt.setupValidator(mockValidator)
			}

			// Act
			result, err := validationService.GetByID(context.Background(), tt.userID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)

				// Check error type
				if _, ok := tt.expectedError.(validationDomain.ValidationError); ok {
					var validationError validationDomain.ValidationError
					assert.ErrorAs(t, err, &validationError, "Error should be ValidationError type")
					assert.Equal(t, "user_id", validationError.Field)
				} else {
					assert.ErrorIs(t, err, tt.expectedError)
				}

				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedUser.ID, result.ID)
				assert.Equal(t, tt.expectedUser.Email, result.Email)
				assert.Equal(t, tt.expectedUser.FirstName, result.FirstName)
				assert.Equal(t, tt.expectedUser.LastName, result.LastName)
			}

			// Verify mock expectations
			if tt.expectNextCalled {
				mockNext.AssertExpectations(t)
			} else {
				mockNext.AssertNotCalled(t, "GetByID")
			}
		})
	}
}

func TestUserValidationService_UpdatePreferences(t *testing.T) {
	tests := []struct {
		name                string
		setupMocks          func(*usermock.MockUserService)
		setupValidator      func(*usermock.MockValidationService)
		userID              string
		preferences         user.UserPreferences
		expectedError       error
		expectNextCalled    bool
		expectedFieldErrors []string
	}{
		{
			name: "Given valid preferences, When UpdatePreferences is called, Then should validate and pass to next service",
			setupMocks: func(mockNext *usermock.MockUserService) {
				validID := "550e8400-e29b-41d4-a716-446655440000"
				mockNext.On("UpdatePreferences", mock.Anything, validID, mock.Anything).Return(nil)
			},
			setupValidator: func(mockValidator *usermock.MockValidationService) {
				validID := "550e8400-e29b-41d4-a716-446655440000"
				mockValidator.On("ValidateUserID", mock.Anything, validID).Return(nil)
				mockValidator.On("ValidateUserPreferences", mock.Anything, mock.Anything).Return(nil)
			},
			userID: "550e8400-e29b-41d4-a716-446655440000",
			preferences: user.UserPreferences{
				ID:                 uuid.New(),
				UserID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				EmailNotifications: true,
				PushNotifications:  false,
				SMSNotifications:   false,
				Theme:              "dark",
				Language:           "en",
				Timezone:           "UTC",
				NotificationTypes: map[string]bool{
					"task_assigned":   true,
					"project_updated": false,
				},
			},
			expectedError:       nil,
			expectNextCalled:    true,
			expectedFieldErrors: nil,
		},
		{
			name: "Given invalid theme, When UpdatePreferences is called, Then should return validation error and not call next service",
			setupMocks: func(mockNext *usermock.MockUserService) {
				// Next service should not be called
			},
			setupValidator: func(mockValidator *usermock.MockValidationService) {
				validID := "550e8400-e29b-41d4-a716-446655440000"
				mockValidator.On("ValidateUserID", mock.Anything, validID).Return(nil)
				validationError := validationDomain.ValidationErrors{
					Errors: []validationDomain.ValidationError{
						{Field: "theme", Message: "must be one of: light, dark, auto"},
					},
				}
				mockValidator.On("ValidateUserPreferences", mock.Anything, mock.Anything).Return(validationError)
			},
			userID: "550e8400-e29b-41d4-a716-446655440000",
			preferences: user.UserPreferences{
				Theme:    "invalid-theme", // Invalid theme
				Language: "en",
				Timezone: "UTC",
			},
			expectedError:       validationDomain.ValidationErrors{},
			expectNextCalled:    false,
			expectedFieldErrors: []string{"theme"},
		},
		{
			name: "Given invalid language code, When UpdatePreferences is called, Then should return validation error and not call next service",
			setupMocks: func(mockNext *usermock.MockUserService) {
				// Next service should not be called
			},
			setupValidator: func(mockValidator *usermock.MockValidationService) {
				validID := "550e8400-e29b-41d4-a716-446655440000"
				mockValidator.On("ValidateUserID", mock.Anything, validID).Return(nil)
				validationError := validationDomain.ValidationErrors{
					Errors: []validationDomain.ValidationError{
						{Field: "language", Message: "must be a 2-letter language code"},
					},
				}
				mockValidator.On("ValidateUserPreferences", mock.Anything, mock.Anything).Return(validationError)
			},
			userID: "550e8400-e29b-41d4-a716-446655440000",
			preferences: user.UserPreferences{
				Theme:    "light",
				Language: "invalid", // Invalid language code (should be 2 chars)
				Timezone: "UTC",
			},
			expectedError:       validationDomain.ValidationErrors{},
			expectNextCalled:    false,
			expectedFieldErrors: []string{"language"},
		},
		{
			name: "Given invalid notification type, When UpdatePreferences is called, Then should return validation error and not call next service",
			setupMocks: func(mockNext *usermock.MockUserService) {
				// Next service should not be called
			},
			setupValidator: func(mockValidator *usermock.MockValidationService) {
				validID := "550e8400-e29b-41d4-a716-446655440000"
				mockValidator.On("ValidateUserID", mock.Anything, validID).Return(nil)
				validationError := validationDomain.ValidationErrors{
					Errors: []validationDomain.ValidationError{
						{Field: "notification_types", Message: "invalid notification type"},
					},
				}
				mockValidator.On("ValidateUserPreferences", mock.Anything, mock.Anything).Return(validationError)
			},
			userID: "550e8400-e29b-41d4-a716-446655440000",
			preferences: user.UserPreferences{
				Theme:    "light",
				Language: "en",
				Timezone: "UTC",
				NotificationTypes: map[string]bool{
					"invalid_notification_type": true, // Invalid notification type
				},
			},
			expectedError:       validationDomain.ValidationErrors{},
			expectNextCalled:    false,
			expectedFieldErrors: []string{"notification_types"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockNext := new(usermock.MockUserService)
			mockValidator := new(usermock.MockValidationService)
			validationService := validation.NewService(mockNext, mockValidator)

			tt.setupMocks(mockNext)
			if tt.setupValidator != nil {
				tt.setupValidator(mockValidator)
			}

			// Act
			err := validationService.UpdatePreferences(context.Background(), tt.userID, tt.preferences)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)

				// Check if it's a validation error with expected fields
				if len(tt.expectedFieldErrors) > 0 {
					var validationErrors validationDomain.ValidationErrors
					assert.ErrorAs(t, err, &validationErrors, "Error should be ValidationErrors type")

					// Verify that all expected field errors are present
					errorFields := make(map[string]bool)
					for _, validationErr := range validationErrors.Errors {
						errorFields[validationErr.Field] = true
					}

					for _, expectedField := range tt.expectedFieldErrors {
						assert.True(t, errorFields[expectedField], "Expected validation error for field: %s", expectedField)
					}
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify mock expectations
			if tt.expectNextCalled {
				mockNext.AssertExpectations(t)
			} else {
				mockNext.AssertNotCalled(t, "UpdatePreferences")
			}
		})
	}
}
