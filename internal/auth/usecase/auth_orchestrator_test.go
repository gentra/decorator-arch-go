package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gentra/decorator-arch-go/internal/auth"
	authmock "github.com/gentra/decorator-arch-go/internal/auth/mock"
	"github.com/gentra/decorator-arch-go/internal/auth/usecase"
)

func TestAuthOrchestrator_Authenticate(t *testing.T) {
	testCases := []struct {
		name        string
		strategy    string
		credentials interface{}
		setupMocks  func(*authmock.MockAuthStrategy)
		expectError bool
		expectedErr error
	}{
		{
			name:        "Given valid basic strategy, When Authenticate is called, Then should delegate to strategy and return success",
			strategy:    "basic",
			credentials: auth.BasicCredentials{Email: "test@example.com", Password: "password123"},
			setupMocks: func(mockStrategy *authmock.MockAuthStrategy) {
				expectedResult := &auth.AuthResult{
					User: &auth.User{
						ID:    "user-123",
						Email: "test@example.com",
					},
					Token:     "jwt-token",
					ExpiresAt: time.Now().Add(time.Hour),
					Strategy:  "basic",
				}
				mockStrategy.On("Authenticate", mock.Anything, "basic", mock.Anything).Return(expectedResult, nil)
			},
			expectError: false,
		},
		{
			name:        "Given unsupported strategy, When Authenticate is called, Then should return unsupported strategy error",
			strategy:    "unknown",
			credentials: auth.BasicCredentials{Email: "test@example.com", Password: "password123"},
			setupMocks:  func(mockStrategy *authmock.MockAuthStrategy) {},
			expectError: true,
			expectedErr: auth.ErrUnsupportedStrategy,
		},
		{
			name:        "Given valid strategy but invalid credentials, When Authenticate is called, Then should return credentials error",
			strategy:    "basic",
			credentials: auth.BasicCredentials{Email: "invalid@example.com", Password: "wrong"},
			setupMocks: func(mockStrategy *authmock.MockAuthStrategy) {
				mockStrategy.On("Authenticate", mock.Anything, "basic", mock.Anything).Return(nil, auth.ErrInvalidCredentials)
			},
			expectError: true,
			expectedErr: auth.ErrInvalidCredentials,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			secret := []byte("test-secret-key-for-testing")
			tokenManager := usecase.NewJWTTokenManager(secret, time.Hour, 24*time.Hour)
			orchestrator := usecase.NewAuthOrchestrator(tokenManager)

			if tt.strategy == "basic" {
				mockStrategy := new(authmock.MockAuthStrategy)
				tt.setupMocks(mockStrategy)
				orchestrator.RegisterStrategy("basic", mockStrategy)
			}

			// Act
			result, err := orchestrator.Authenticate(context.Background(), tt.strategy, tt.credentials)

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
				assert.Equal(t, tt.strategy, result.Strategy)
			}
		})
	}
}

func TestAuthOrchestrator_ValidateToken(t *testing.T) {
	testCases := []struct {
		name        string
		setupToken  func(*usecase.JWTTokenManager) string
		expectError bool
		expectedErr error
	}{
		{
			name: "Given valid token, When ValidateToken is called, Then should return token claims",
			setupToken: func(tokenManager *usecase.JWTTokenManager) string {
				token, _, _ := tokenManager.GenerateAuthToken("user-123", "test@example.com")
				return token
			},
			expectError: false,
		},
		{
			name: "Given invalid token, When ValidateToken is called, Then should return invalid token error",
			setupToken: func(tokenManager *usecase.JWTTokenManager) string {
				return "invalid-jwt-token"
			},
			expectError: true,
		},
		{
			name: "Given malformed token, When ValidateToken is called, Then should return parse error",
			setupToken: func(tokenManager *usecase.JWTTokenManager) string {
				return "malformed.token.format"
			},
			expectError: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			secret := []byte("test-secret-key-for-testing")
			tokenManager := usecase.NewJWTTokenManager(secret, time.Hour, 24*time.Hour)
			orchestrator := usecase.NewAuthOrchestrator(tokenManager)

			testToken := tt.setupToken(tokenManager)

			// Act
			claims, err := orchestrator.ValidateToken(context.Background(), testToken)

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

func TestAuthOrchestrator_RefreshToken(t *testing.T) {
	testCases := []struct {
		name        string
		setupToken  func(*usecase.JWTTokenManager) string
		expectError bool
		expectedErr error
	}{
		{
			name: "Given valid refresh token, When RefreshToken is called, Then should generate new access token",
			setupToken: func(tokenManager *usecase.JWTTokenManager) string {
				refreshToken, _ := tokenManager.GenerateRefreshToken("user-123")
				return refreshToken
			},
			expectError: false,
		},
		{
			name: "Given access token instead of refresh token, When RefreshToken is called, Then should return invalid refresh token error",
			setupToken: func(tokenManager *usecase.JWTTokenManager) string {
				accessToken, _, _ := tokenManager.GenerateAuthToken("user-123", "test@example.com")
				return accessToken
			},
			expectError: true,
			expectedErr: auth.ErrInvalidRefreshToken,
		},
		{
			name: "Given invalid refresh token, When RefreshToken is called, Then should return validation error",
			setupToken: func(tokenManager *usecase.JWTTokenManager) string {
				return "invalid-refresh-token"
			},
			expectError: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			secret := []byte("test-secret-key-for-testing")
			tokenManager := usecase.NewJWTTokenManager(secret, time.Hour, 24*time.Hour)
			orchestrator := usecase.NewAuthOrchestrator(tokenManager)

			testToken := tt.setupToken(tokenManager)

			// Act
			result, err := orchestrator.RefreshToken(context.Background(), testToken)

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
				assert.NotEmpty(t, result.Token)
				assert.Equal(t, testToken, result.RefreshToken) // Should keep same refresh token
			}
		})
	}
}

func TestAuthOrchestrator_RevokeToken(t *testing.T) {
	testCases := []struct {
		name        string
		setupToken  func(*usecase.JWTTokenManager) string
		expectError bool
	}{
		{
			name: "Given valid token, When RevokeToken is called, Then should revoke successfully",
			setupToken: func(tokenManager *usecase.JWTTokenManager) string {
				token, _, _ := tokenManager.GenerateAuthToken("user-123", "test@example.com")
				return token
			},
			expectError: false,
		},
		{
			name: "Given invalid token, When RevokeToken is called, Then should return error",
			setupToken: func(tokenManager *usecase.JWTTokenManager) string {
				return "invalid-token"
			},
			expectError: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			secret := []byte("test-secret-key-for-testing")
			tokenManager := usecase.NewJWTTokenManager(secret, time.Hour, 24*time.Hour)
			orchestrator := usecase.NewAuthOrchestrator(tokenManager)

			testToken := tt.setupToken(tokenManager)

			// Act
			err := orchestrator.RevokeToken(context.Background(), testToken)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthOrchestrator_GetSupportedStrategies(t *testing.T) {
	testCases := []struct {
		name               string
		registerStrategies func(*usecase.AuthOrchestrator)
		expectedStrategies []string
	}{
		{
			name: "Given no strategies registered, When GetSupportedStrategies is called, Then should return empty list",
			registerStrategies: func(orchestrator *usecase.AuthOrchestrator) {
				// No strategies registered
			},
			expectedStrategies: []string{},
		},
		{
			name: "Given basic strategy registered, When GetSupportedStrategies is called, Then should return basic strategy",
			registerStrategies: func(orchestrator *usecase.AuthOrchestrator) {
				mockStrategy := new(authmock.MockAuthStrategy)
				orchestrator.RegisterStrategy("basic", mockStrategy)
			},
			expectedStrategies: []string{"basic"},
		},
		{
			name: "Given multiple strategies registered, When GetSupportedStrategies is called, Then should return all strategies",
			registerStrategies: func(orchestrator *usecase.AuthOrchestrator) {
				mockBasic := new(authmock.MockAuthStrategy)
				mockOAuth := new(authmock.MockAuthStrategy)
				mockJWT := new(authmock.MockAuthStrategy)

				orchestrator.RegisterStrategy("basic", mockBasic)
				orchestrator.RegisterStrategy("oauth", mockOAuth)
				orchestrator.RegisterStrategy("jwt", mockJWT)
			},
			expectedStrategies: []string{"basic", "oauth", "jwt"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			secret := []byte("test-secret-key-for-testing")
			tokenManager := usecase.NewJWTTokenManager(secret, time.Hour, 24*time.Hour)
			orchestrator := usecase.NewAuthOrchestrator(tokenManager)
			tt.registerStrategies(orchestrator)

			// Act
			strategies := orchestrator.GetSupportedStrategies()

			// Assert
			assert.Len(t, strategies, len(tt.expectedStrategies))
			for _, expected := range tt.expectedStrategies {
				assert.Contains(t, strategies, expected)
			}
		})
	}
}
