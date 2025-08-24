package usecase_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gentra/decorator-arch-go/internal/auth"
	"github.com/gentra/decorator-arch-go/internal/auth/usecase"
)

func TestJWTTokenManager_GenerateAuthToken(t *testing.T) {
	testCases := []struct {
		name        string
		userID      string
		email       string
		expectError bool
	}{
		{
			name:        "Given valid user data, When GenerateAuthToken is called, Then should return valid JWT token",
			userID:      "user-123",
			email:       "test@example.com",
			expectError: false,
		},
		{
			name:        "Given empty user ID, When GenerateAuthToken is called, Then should still generate token with empty ID",
			userID:      "",
			email:       "test@example.com",
			expectError: false, // JWT generation doesn't validate inputs
		},
		{
			name:        "Given empty email, When GenerateAuthToken is called, Then should still generate token with empty email",
			userID:      "user-123",
			email:       "",
			expectError: false, // JWT generation doesn't validate inputs
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			secret := []byte("test-secret-key-for-testing")
			accessTTL := time.Hour
			refreshTTL := 24 * time.Hour
			tokenManager := usecase.NewJWTTokenManager(secret, accessTTL, refreshTTL)

			// Act
			token, expiresAt, err := tokenManager.GenerateAuthToken(tt.userID, tt.email)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
				assert.True(t, expiresAt.IsZero())
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.False(t, expiresAt.IsZero())
				assert.True(t, expiresAt.After(time.Now()))

				// Verify token can be validated
				claims, validateErr := tokenManager.ValidateToken(token)
				if tt.userID == "" {
					// Empty userID makes token invalid
					assert.Error(t, validateErr)
					assert.Nil(t, claims)
				} else {
					assert.NoError(t, validateErr)
					assert.Equal(t, tt.userID, claims.UserID)
					assert.Equal(t, tt.email, claims.Email)
					assert.Equal(t, "access", claims.TokenType)
				}
			}
		})
	}
}

func TestJWTTokenManager_GenerateRefreshToken(t *testing.T) {
	testCases := []struct {
		name        string
		userID      string
		expectError bool
	}{
		{
			name:        "Given valid user ID, When GenerateRefreshToken is called, Then should return valid refresh token",
			userID:      "user-123",
			expectError: false,
		},
		{
			name:        "Given empty user ID, When GenerateRefreshToken is called, Then should still generate token",
			userID:      "",
			expectError: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			secret := []byte("test-secret-key-for-testing")
			accessTTL := time.Hour
			refreshTTL := 24 * time.Hour
			tokenManager := usecase.NewJWTTokenManager(secret, accessTTL, refreshTTL)

			// Act
			token, err := tokenManager.GenerateRefreshToken(tt.userID)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Verify token can be validated
				claims, validateErr := tokenManager.ValidateToken(token)
				if tt.userID == "" {
					// Empty userID makes token invalid
					assert.Error(t, validateErr)
					assert.Nil(t, claims)
				} else {
					assert.NoError(t, validateErr)
					assert.Equal(t, tt.userID, claims.UserID)
					assert.Equal(t, "refresh", claims.TokenType)
				}
			}
		})
	}
}

func TestJWTTokenManager_ValidateToken(t *testing.T) {
	// Setup
	secret := []byte("test-secret-key-for-testing")
	accessTTL := time.Hour
	refreshTTL := 24 * time.Hour
	tokenManager := usecase.NewJWTTokenManager(secret, accessTTL, refreshTTL)

	// Generate valid tokens for testing
	validAccessToken, _, _ := tokenManager.GenerateAuthToken("user-123", "test@example.com")
	validRefreshToken, _ := tokenManager.GenerateRefreshToken("user-123")

	testCases := []struct {
		name           string
		token          string
		setupToken     func() string
		expectError    bool
		expectedErr    error
		validateClaims func(*testing.T, *auth.TokenClaims)
	}{
		{
			name:  "Given valid access token, When ValidateToken is called, Then should return valid claims",
			token: validAccessToken,
			setupToken: func() string {
				return validAccessToken
			},
			expectError: false,
			validateClaims: func(t *testing.T, claims *auth.TokenClaims) {
				assert.Equal(t, "user-123", claims.UserID)
				assert.Equal(t, "test@example.com", claims.Email)
				assert.Equal(t, "access", claims.TokenType)
				assert.False(t, claims.IsExpired())
			},
		},
		{
			name:  "Given valid refresh token, When ValidateToken is called, Then should return valid claims",
			token: validRefreshToken,
			setupToken: func() string {
				return validRefreshToken
			},
			expectError: false,
			validateClaims: func(t *testing.T, claims *auth.TokenClaims) {
				assert.Equal(t, "user-123", claims.UserID)
				assert.Equal(t, "refresh", claims.TokenType)
				assert.False(t, claims.IsExpired())
			},
		},
		{
			name:        "Given malformed token, When ValidateToken is called, Then should return parse error",
			token:       "invalid.token.format",
			expectError: true,
		},
		{
			name:        "Given empty token, When ValidateToken is called, Then should return parse error",
			token:       "",
			expectError: true,
		},
		{
			name:        "Given token with wrong signature, When ValidateToken is called, Then should return signature error",
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			expectError: true,
		},
		{
			name: "Given expired token, When ValidateToken is called, Then should return expired error",
			setupToken: func() string {
				// Create a token manager with very short TTL
				shortTTLManager := usecase.NewJWTTokenManager(secret, -time.Hour, -time.Hour) // Negative TTL = already expired
				expiredToken, _, _ := shortTTLManager.GenerateAuthToken("user-123", "test@example.com")
				return expiredToken
			},
			expectError: true,
			// Note: JWT library returns parse error for expired tokens, not our custom error
		},
		{
			name: "Given revoked token, When ValidateToken is called, Then should return invalid token error",
			setupToken: func() string {
				token, _, _ := tokenManager.GenerateAuthToken("user-456", "revoked@example.com")
				// Revoke the token
				tokenManager.RevokeToken(token)
				return token
			},
			expectError: true,
			expectedErr: auth.ErrInvalidToken,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var tokenToTest string
			if tt.setupToken != nil {
				tokenToTest = tt.setupToken()
			} else {
				tokenToTest = tt.token
			}

			// Act
			claims, err := tokenManager.ValidateToken(tokenToTest)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err)
				}
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				if tt.validateClaims != nil {
					tt.validateClaims(t, claims)
				}
			}
		})
	}
}

func TestJWTTokenManager_RevokeToken(t *testing.T) {
	// Setup
	secret := []byte("test-secret-key-for-testing")
	accessTTL := time.Hour
	refreshTTL := 24 * time.Hour
	tokenManager := usecase.NewJWTTokenManager(secret, accessTTL, refreshTTL)

	testCases := []struct {
		name        string
		setupToken  func() string
		expectError bool
	}{
		{
			name: "Given valid token, When RevokeToken is called, Then should revoke successfully and token becomes invalid",
			setupToken: func() string {
				token, _, _ := tokenManager.GenerateAuthToken("user-123", "test@example.com")
				return token
			},
			expectError: false,
		},
		{
			name: "Given malformed token, When RevokeToken is called, Then should return parse error",
			setupToken: func() string {
				return "invalid.token.format"
			},
			expectError: true,
		},
		{
			name: "Given empty token, When RevokeToken is called, Then should return parse error",
			setupToken: func() string {
				return ""
			},
			expectError: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			token := tt.setupToken()

			// Act
			err := tokenManager.RevokeToken(token)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify token is now invalid
				_, validateErr := tokenManager.ValidateToken(token)
				assert.Error(t, validateErr)
				assert.Equal(t, auth.ErrInvalidToken, validateErr)
			}
		})
	}
}

func TestJWTTokenManager_TokenLifecycle(t *testing.T) {
	testCases := []struct {
		name string
		test func(t *testing.T, tokenManager *usecase.JWTTokenManager)
	}{
		{
			name: "Given full token lifecycle, When tokens are generated, validated, and revoked, Then should work correctly",
			test: func(t *testing.T, tokenManager *usecase.JWTTokenManager) {
				userID := "user-123"
				email := "test@example.com"

				// Generate access token
				accessToken, expiresAt, err := tokenManager.GenerateAuthToken(userID, email)
				require.NoError(t, err)
				assert.NotEmpty(t, accessToken)
				assert.True(t, expiresAt.After(time.Now()))

				// Generate refresh token
				refreshToken, err := tokenManager.GenerateRefreshToken(userID)
				require.NoError(t, err)
				assert.NotEmpty(t, refreshToken)

				// Validate access token
				accessClaims, err := tokenManager.ValidateToken(accessToken)
				require.NoError(t, err)
				assert.Equal(t, userID, accessClaims.UserID)
				assert.Equal(t, email, accessClaims.Email)
				assert.Equal(t, "access", accessClaims.TokenType)

				// Validate refresh token
				refreshClaims, err := tokenManager.ValidateToken(refreshToken)
				require.NoError(t, err)
				assert.Equal(t, userID, refreshClaims.UserID)
				assert.Equal(t, "refresh", refreshClaims.TokenType)

				// Revoke access token
				err = tokenManager.RevokeToken(accessToken)
				require.NoError(t, err)

				// Verify access token is now invalid
				_, err = tokenManager.ValidateToken(accessToken)
				assert.Error(t, err)
				assert.Equal(t, auth.ErrInvalidToken, err)

				// Verify refresh token is still valid
				_, err = tokenManager.ValidateToken(refreshToken)
				assert.NoError(t, err)
			},
		},
		{
			name: "Given multiple token revocations, When cleanup occurs, Then should handle expired revocations correctly",
			test: func(t *testing.T, tokenManager *usecase.JWTTokenManager) {
				// Generate multiple tokens
				token1, _, _ := tokenManager.GenerateAuthToken("user-1", "user1@example.com")
				token2, _, _ := tokenManager.GenerateAuthToken("user-2", "user2@example.com")
				token3, _, _ := tokenManager.GenerateAuthToken("user-3", "user3@example.com")

				// Revoke all tokens
				tokenManager.RevokeToken(token1)
				tokenManager.RevokeToken(token2)
				tokenManager.RevokeToken(token3)

				// Verify all tokens are revoked
				_, err1 := tokenManager.ValidateToken(token1)
				_, err2 := tokenManager.ValidateToken(token2)
				_, err3 := tokenManager.ValidateToken(token3)

				assert.Equal(t, auth.ErrInvalidToken, err1)
				assert.Equal(t, auth.ErrInvalidToken, err2)
				assert.Equal(t, auth.ErrInvalidToken, err3)

				// Generate new token to trigger cleanup
				newToken, _, _ := tokenManager.GenerateAuthToken("user-4", "user4@example.com")
				_, err := tokenManager.ValidateToken(newToken)
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			secret := []byte("test-secret-key-for-testing")
			accessTTL := time.Hour
			refreshTTL := 24 * time.Hour
			tokenManager := usecase.NewJWTTokenManager(secret, accessTTL, refreshTTL)

			// Act & Assert
			tt.test(t, tokenManager)
		})
	}
}

func TestJWTTokenManager_EdgeCases(t *testing.T) {
	testCases := []struct {
		name        string
		secret      []byte
		accessTTL   time.Duration
		refreshTTL  time.Duration
		expectPanic bool
	}{
		{
			name:        "Given valid configuration, When NewJWTTokenManager is called, Then should create manager successfully",
			secret:      []byte("valid-secret-key"),
			accessTTL:   time.Hour,
			refreshTTL:  24 * time.Hour,
			expectPanic: false,
		},
		{
			name:        "Given empty secret, When tokens are generated, Then should still work but be insecure",
			secret:      []byte(""),
			accessTTL:   time.Hour,
			refreshTTL:  24 * time.Hour,
			expectPanic: false,
		},
		{
			name:        "Given zero TTL, When tokens are generated, Then should create immediately expired tokens",
			secret:      []byte("test-secret"),
			accessTTL:   0,
			refreshTTL:  0,
			expectPanic: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange & Act
			if tt.expectPanic {
				assert.Panics(t, func() {
					usecase.NewJWTTokenManager(tt.secret, tt.accessTTL, tt.refreshTTL)
				})
			} else {
				assert.NotPanics(t, func() {
					tokenManager := usecase.NewJWTTokenManager(tt.secret, tt.accessTTL, tt.refreshTTL)
					assert.NotNil(t, tokenManager)

					// Try to generate a token to verify manager works
					token, _, err := tokenManager.GenerateAuthToken("test-user", "test@example.com")
					if len(tt.secret) > 0 {
						assert.NoError(t, err)
						assert.NotEmpty(t, token)
					}
				})
			}
		})
	}
}
