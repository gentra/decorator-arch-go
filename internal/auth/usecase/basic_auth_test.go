package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gentra/decorator-arch-go/internal/auth"
	authmock "github.com/gentra/decorator-arch-go/internal/auth/mock"
	"github.com/gentra/decorator-arch-go/internal/auth/usecase"
	"github.com/gentra/decorator-arch-go/internal/user"
)

func TestBasicAuthStrategy_Authenticate(t *testing.T) {
	testCases := []struct {
		name           string
		strategy       string
		credentials    interface{}
		setupMocks     func(*authmock.MockUserService)
		expectError    bool
		expectedErr    error
		validateResult func(*testing.T, *auth.AuthResult)
	}{
		{
			name:     "Given valid basic credentials, When Authenticate is called with basic strategy, Then should authenticate successfully",
			strategy: "basic",
			credentials: auth.BasicCredentials{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(mockUser *authmock.MockUserService) {
				// Mock successful login
				loginResult := &user.AuthResult{
					User: &user.User{
						ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
						Email:     "test@example.com",
						FirstName: "John",
						LastName:  "Doe",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
				mockUser.On("Login", mock.Anything, "test@example.com", "password123").Return(loginResult, nil)
			},
			expectError: false,
			validateResult: func(t *testing.T, result *auth.AuthResult) {
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.User.ID)
				assert.Equal(t, "test@example.com", result.User.Email)
				assert.Equal(t, "John", result.User.FirstName)
				assert.Equal(t, "Doe", result.User.LastName)
				assert.NotEmpty(t, result.Token)
				assert.NotEmpty(t, result.RefreshToken)
				assert.Equal(t, "basic", result.Strategy)
			},
		},
		{
			name:     "Given unsupported strategy, When Authenticate is called, Then should return unsupported strategy error",
			strategy: "oauth",
			credentials: auth.BasicCredentials{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(mockUser *authmock.MockUserService) {
				// No mocks needed - should fail before calling user service
			},
			expectError: true,
			expectedErr: auth.ErrUnsupportedStrategy,
		},
		{
			name:     "Given invalid credentials type, When Authenticate is called, Then should return credentials type error",
			strategy: "basic",
			credentials: auth.OAuthCredentials{
				Provider:    "google",
				AccessToken: "oauth-token",
			},
			setupMocks: func(mockUser *authmock.MockUserService) {
				// No mocks needed - should fail type assertion
			},
			expectError: true,
		},
		{
			name:     "Given invalid login credentials, When Authenticate is called, Then should return invalid credentials error",
			strategy: "basic",
			credentials: auth.BasicCredentials{
				Email:    "invalid@example.com",
				Password: "wrongpassword",
			},
			setupMocks: func(mockUser *authmock.MockUserService) {
				mockUser.On("Login", mock.Anything, "invalid@example.com", "wrongpassword").Return(nil, user.ErrInvalidCredentials)
			},
			expectError: true,
			expectedErr: auth.ErrInvalidCredentials,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserService := new(authmock.MockUserService)
			tt.setupMocks(mockUserService)

			// Create real JWT token manager for integration testing
			secret := []byte("test-secret-key-for-testing")
			tokenManager := usecase.NewJWTTokenManager(secret, time.Hour, 24*time.Hour)

			basicAuth := usecase.NewBasicAuthStrategy(mockUserService, tokenManager)

			// Act
			result, err := basicAuth.Authenticate(context.Background(), tt.strategy, tt.credentials)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
			}

			mockUserService.AssertExpectations(t)
		})
	}
}

func TestBasicAuthStrategy_ValidateToken(t *testing.T) {
	testCases := []struct {
		name        string
		setupToken  func(*usecase.JWTTokenManager) string
		expectError bool
	}{
		{
			name: "Given valid token, When ValidateToken is called, Then should return valid claims",
			setupToken: func(tokenManager *usecase.JWTTokenManager) string {
				token, _, _ := tokenManager.GenerateAuthToken("user-123", "test@example.com")
				return token
			},
			expectError: false,
		},
		{
			name: "Given invalid token, When ValidateToken is called, Then should return validation error",
			setupToken: func(tokenManager *usecase.JWTTokenManager) string {
				return "invalid-jwt-token"
			},
			expectError: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserService := new(authmock.MockUserService)
			secret := []byte("test-secret-key-for-testing")
			tokenManager := usecase.NewJWTTokenManager(secret, time.Hour, 24*time.Hour)

			testToken := tt.setupToken(tokenManager)
			basicAuth := usecase.NewBasicAuthStrategy(mockUserService, tokenManager)

			// Act
			claims, err := basicAuth.ValidateToken(context.Background(), testToken)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, "user-123", claims.UserID)
			}
		})
	}
}

func TestBasicAuthStrategy_GetSupportedStrategies(t *testing.T) {
	t.Run("Given BasicAuthStrategy, When GetSupportedStrategies is called, Then should return only basic strategy", func(t *testing.T) {
		// Arrange
		mockUserService := new(authmock.MockUserService)
		secret := []byte("test-secret-key-for-testing")
		tokenManager := usecase.NewJWTTokenManager(secret, time.Hour, 24*time.Hour)
		basicAuth := usecase.NewBasicAuthStrategy(mockUserService, tokenManager)

		// Act
		strategies := basicAuth.GetSupportedStrategies()

		// Assert
		assert.Equal(t, []string{"basic"}, strategies)
	})
}
