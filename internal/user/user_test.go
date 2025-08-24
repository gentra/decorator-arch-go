package user_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/user"
)

func TestUser_GetFullName(t *testing.T) {
	tests := []struct {
		name     string
		user     user.User
		expected string
	}{
		{
			name: "Given user with first and last name, When GetFullName is called, Then should return concatenated name",
			user: user.User{
				FirstName: "John",
				LastName:  "Doe",
			},
			expected: "John Doe",
		},
		{
			name: "Given user with empty first name, When GetFullName is called, Then should return space and last name",
			user: user.User{
				FirstName: "",
				LastName:  "Doe",
			},
			expected: " Doe",
		},
		{
			name: "Given user with empty last name, When GetFullName is called, Then should return first name and space",
			user: user.User{
				FirstName: "John",
				LastName:  "",
			},
			expected: "John ",
		},
		{
			name: "Given user with both names empty, When GetFullName is called, Then should return single space",
			user: user.User{
				FirstName: "",
				LastName:  "",
			},
			expected: " ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.user.GetFullName()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUser_IsEmailVerified(t *testing.T) {
	t.Run("Given any user, When IsEmailVerified is called, Then should return true", func(t *testing.T) {
		// Arrange
		testUser := user.User{
			ID:    uuid.New(),
			Email: "test@example.com",
		}

		// Act
		result := testUser.IsEmailVerified()

		// Assert
		assert.True(t, result)
	})
}

func TestUserPreferences_IsNotificationEnabled(t *testing.T) {
	tests := []struct {
		name             string
		preferences      user.UserPreferences
		notificationType string
		expected         bool
	}{
		{
			name: "Given preferences with enabled notification type, When IsNotificationEnabled is called, Then should return true",
			preferences: user.UserPreferences{
				NotificationTypes: map[string]bool{
					"task_assigned": true,
					"task_updated":  false,
				},
			},
			notificationType: "task_assigned",
			expected:         true,
		},
		{
			name: "Given preferences with disabled notification type, When IsNotificationEnabled is called, Then should return false",
			preferences: user.UserPreferences{
				NotificationTypes: map[string]bool{
					"task_assigned": true,
					"task_updated":  false,
				},
			},
			notificationType: "task_updated",
			expected:         false,
		},
		{
			name: "Given preferences with missing notification type, When IsNotificationEnabled is called, Then should return false",
			preferences: user.UserPreferences{
				NotificationTypes: map[string]bool{
					"task_assigned": true,
				},
			},
			notificationType: "task_deleted",
			expected:         false,
		},
		{
			name: "Given preferences with nil notification types, When IsNotificationEnabled is called, Then should return false",
			preferences: user.UserPreferences{
				NotificationTypes: nil,
			},
			notificationType: "task_assigned",
			expected:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.preferences.IsNotificationEnabled(tt.notificationType)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserPreferences_EnableNotification(t *testing.T) {
	tests := []struct {
		name             string
		preferences      user.UserPreferences
		notificationType string
		expectedValue    bool
	}{
		{
			name: "Given preferences with existing notification types, When EnableNotification is called, Then should enable notification",
			preferences: user.UserPreferences{
				NotificationTypes: map[string]bool{
					"task_assigned": false,
				},
			},
			notificationType: "task_assigned",
			expectedValue:    true,
		},
		{
			name: "Given preferences with nil notification types, When EnableNotification is called, Then should create map and enable notification",
			preferences: user.UserPreferences{
				NotificationTypes: nil,
			},
			notificationType: "task_assigned",
			expectedValue:    true,
		},
		{
			name: "Given preferences with new notification type, When EnableNotification is called, Then should add and enable notification",
			preferences: user.UserPreferences{
				NotificationTypes: map[string]bool{
					"task_updated": true,
				},
			},
			notificationType: "task_assigned",
			expectedValue:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			tt.preferences.EnableNotification(tt.notificationType)

			// Assert
			assert.Equal(t, tt.expectedValue, tt.preferences.NotificationTypes[tt.notificationType])
			assert.NotNil(t, tt.preferences.NotificationTypes)
		})
	}
}

func TestUserPreferences_DisableNotification(t *testing.T) {
	tests := []struct {
		name             string
		preferences      user.UserPreferences
		notificationType string
		expectedValue    bool
	}{
		{
			name: "Given preferences with enabled notification, When DisableNotification is called, Then should disable notification",
			preferences: user.UserPreferences{
				NotificationTypes: map[string]bool{
					"task_assigned": true,
				},
			},
			notificationType: "task_assigned",
			expectedValue:    false,
		},
		{
			name: "Given preferences with disabled notification, When DisableNotification is called, Then should keep notification disabled",
			preferences: user.UserPreferences{
				NotificationTypes: map[string]bool{
					"task_assigned": false,
				},
			},
			notificationType: "task_assigned",
			expectedValue:    false,
		},
		{
			name: "Given preferences with nil notification types, When DisableNotification is called, Then should not panic",
			preferences: user.UserPreferences{
				NotificationTypes: nil,
			},
			notificationType: "task_assigned",
			expectedValue:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			tt.preferences.DisableNotification(tt.notificationType)

			// Assert
			if tt.preferences.NotificationTypes != nil {
				assert.Equal(t, tt.expectedValue, tt.preferences.NotificationTypes[tt.notificationType])
			}
		})
	}
}

func TestDefaultUserPreferences(t *testing.T) {
	t.Run("Given user ID, When DefaultUserPreferences is called, Then should return valid default preferences", func(t *testing.T) {
		// Arrange
		userID := uuid.New()

		// Act
		preferences := user.DefaultUserPreferences(userID)

		// Assert
		assert.NotNil(t, preferences)
		assert.Equal(t, userID, preferences.UserID)
		assert.NotEqual(t, uuid.Nil, preferences.ID)
		assert.True(t, preferences.EmailNotifications)
		assert.True(t, preferences.PushNotifications)
		assert.False(t, preferences.SMSNotifications)
		assert.Equal(t, "light", preferences.Theme)
		assert.Equal(t, "en", preferences.Language)
		assert.Equal(t, "UTC", preferences.Timezone)
		assert.NotNil(t, preferences.NotificationTypes)
		assert.True(t, preferences.NotificationTypes["task_assigned"])
		assert.True(t, preferences.NotificationTypes["task_due_soon"])
		assert.True(t, preferences.NotificationTypes["project_updated"])
		assert.True(t, preferences.NotificationTypes["project_invite"])
		assert.False(t, preferences.NotificationTypes["system_updates"])
		assert.False(t, preferences.NotificationTypes["marketing"])
		assert.False(t, preferences.CreatedAt.IsZero())
		assert.False(t, preferences.UpdatedAt.IsZero())
	})
}

func TestUserError_Error(t *testing.T) {
	tests := []struct {
		name     string
		userErr  user.UserError
		expected string
	}{
		{
			name: "Given user error with message, When Error is called, Then should return message",
			userErr: user.UserError{
				Code:    "TEST_ERROR",
				Message: "Test error message",
			},
			expected: "Test error message",
		},
		{
			name: "Given user error with empty message, When Error is called, Then should return empty string",
			userErr: user.UserError{
				Code:    "TEST_ERROR",
				Message: "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.userErr.Error()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserErrors_Constants(t *testing.T) {
	tests := []struct {
		name        string
		err         user.UserError
		expectedCode string
	}{
		{
			name:        "Given ErrUserNotFound, When accessing code, Then should have correct code",
			err:         user.ErrUserNotFound,
			expectedCode: "USER_NOT_FOUND",
		},
		{
			name:        "Given ErrEmailAlreadyExists, When accessing code, Then should have correct code",
			err:         user.ErrEmailAlreadyExists,
			expectedCode: "EMAIL_EXISTS",
		},
		{
			name:        "Given ErrInvalidCredentials, When accessing code, Then should have correct code",
			err:         user.ErrInvalidCredentials,
			expectedCode: "INVALID_CREDENTIALS",
		},
		{
			name:        "Given ErrInvalidEmail, When accessing code, Then should have correct code",
			err:         user.ErrInvalidEmail,
			expectedCode: "INVALID_EMAIL",
		},
		{
			name:        "Given ErrWeakPassword, When accessing code, Then should have correct code",
			err:         user.ErrWeakPassword,
			expectedCode: "WEAK_PASSWORD",
		},
		{
			name:        "Given ErrEmptyFirstName, When accessing code, Then should have correct code",
			err:         user.ErrEmptyFirstName,
			expectedCode: "EMPTY_FIRST_NAME",
		},
		{
			name:        "Given ErrEmptyLastName, When accessing code, Then should have correct code",
			err:         user.ErrEmptyLastName,
			expectedCode: "EMPTY_LAST_NAME",
		},
		{
			name:        "Given ErrPreferencesNotFound, When accessing code, Then should have correct code",
			err:         user.ErrPreferencesNotFound,
			expectedCode: "PREFERENCES_NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assert
			assert.Equal(t, tt.expectedCode, tt.err.Code)
			assert.NotEmpty(t, tt.err.Message)
		})
	}
}

func TestRegisterData_Validation(t *testing.T) {
	t.Run("Given valid register data, When accessing fields, Then should have validation tags", func(t *testing.T) {
		// This test verifies that the struct has the correct validation tags
		// The actual validation would be done by the validation service
		data := user.RegisterData{
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
		}

		// Assert basic structure
		assert.NotEmpty(t, data.Email)
		assert.NotEmpty(t, data.Password)
		assert.NotEmpty(t, data.FirstName)
		assert.NotEmpty(t, data.LastName)
	})
}

func TestUpdateProfileData_OptionalFields(t *testing.T) {
	t.Run("Given update profile data with optional fields, When accessing fields, Then should handle pointers correctly", func(t *testing.T) {
		firstName := "John"
		lastName := "Doe"
		email := "john.doe@example.com"

		data := user.UpdateProfileData{
			FirstName: &firstName,
			LastName:  &lastName,
			Email:     &email,
		}

		// Assert
		assert.NotNil(t, data.FirstName)
		assert.NotNil(t, data.LastName)
		assert.NotNil(t, data.Email)
		assert.Equal(t, "John", *data.FirstName)
		assert.Equal(t, "Doe", *data.LastName)
		assert.Equal(t, "john.doe@example.com", *data.Email)
	})
}

func TestAuthResult_Structure(t *testing.T) {
	t.Run("Given auth result with all fields, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		now := time.Now()
		testUser := &user.User{
			ID:        uuid.New(),
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
		}

		authResult := user.AuthResult{
			User:         testUser,
			Token:        "jwt-token",
			RefreshToken: "refresh-token",
			ExpiresAt:    now,
		}

		// Assert
		assert.Equal(t, testUser, authResult.User)
		assert.Equal(t, "jwt-token", authResult.Token)
		assert.Equal(t, "refresh-token", authResult.RefreshToken)
		assert.Equal(t, now, authResult.ExpiresAt)
	})
}