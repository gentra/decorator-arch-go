package auth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/auth"
)

func TestUser_GetFullName(t *testing.T) {
	tests := []struct {
		name     string
		user     auth.User
		expected string
	}{
		{
			name: "Given user with first and last name, When GetFullName is called, Then should return concatenated name",
			user: auth.User{
				FirstName: "John",
				LastName:  "Doe",
			},
			expected: "John Doe",
		},
		{
			name: "Given user with empty first name, When GetFullName is called, Then should return space and last name",
			user: auth.User{
				FirstName: "",
				LastName:  "Doe",
			},
			expected: " Doe",
		},
		{
			name: "Given user with empty last name, When GetFullName is called, Then should return first name and space",
			user: auth.User{
				FirstName: "John",
				LastName:  "",
			},
			expected: "John ",
		},
		{
			name: "Given user with both names empty, When GetFullName is called, Then should return single space",
			user: auth.User{
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

func TestUser_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		user     auth.User
		expected bool
	}{
		{
			name: "Given user with ID and email, When IsValid is called, Then should return true",
			user: auth.User{
				ID:    "user-123",
				Email: "test@example.com",
			},
			expected: true,
		},
		{
			name: "Given user with empty ID, When IsValid is called, Then should return false",
			user: auth.User{
				ID:    "",
				Email: "test@example.com",
			},
			expected: false,
		},
		{
			name: "Given user with empty email, When IsValid is called, Then should return false",
			user: auth.User{
				ID:    "user-123",
				Email: "",
			},
			expected: false,
		},
		{
			name: "Given user with both ID and email empty, When IsValid is called, Then should return false",
			user: auth.User{
				ID:    "",
				Email: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.user.IsValid()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthResult_IsValid(t *testing.T) {
	tests := []struct {
		name       string
		authResult auth.AuthResult
		expected   bool
	}{
		{
			name: "Given auth result with user, token and expires at, When IsValid is called, Then should return true",
			authResult: auth.AuthResult{
				User: &auth.User{
					ID:    "user-123",
					Email: "test@example.com",
				},
				Token:     "jwt-token",
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: true,
		},
		{
			name: "Given auth result with nil user, When IsValid is called, Then should return false",
			authResult: auth.AuthResult{
				User:      nil,
				Token:     "jwt-token",
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "Given auth result with empty token, When IsValid is called, Then should return false",
			authResult: auth.AuthResult{
				User: &auth.User{
					ID:    "user-123",
					Email: "test@example.com",
				},
				Token:     "",
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "Given auth result with zero expires at, When IsValid is called, Then should return false",
			authResult: auth.AuthResult{
				User: &auth.User{
					ID:    "user-123",
					Email: "test@example.com",
				},
				Token:     "jwt-token",
				ExpiresAt: time.Time{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.authResult.IsValid()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthResult_IsExpired(t *testing.T) {
	tests := []struct {
		name       string
		authResult auth.AuthResult
		expected   bool
	}{
		{
			name: "Given auth result with future expires at, When IsExpired is called, Then should return false",
			authResult: auth.AuthResult{
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "Given auth result with past expires at, When IsExpired is called, Then should return true",
			authResult: auth.AuthResult{
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			expected: true,
		},
		{
			name: "Given auth result with current time expires at, When IsExpired is called, Then should return true",
			authResult: auth.AuthResult{
				ExpiresAt: time.Now(),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.authResult.IsExpired()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTokenClaims_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		claims   auth.TokenClaims
		expected bool
	}{
		{
			name: "Given token claims with user ID, email and expires at, When IsValid is called, Then should return true",
			claims: auth.TokenClaims{
				UserID:    "user-123",
				Email:     "test@example.com",
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: true,
		},
		{
			name: "Given token claims with empty user ID, When IsValid is called, Then should return false",
			claims: auth.TokenClaims{
				UserID:    "",
				Email:     "test@example.com",
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "Given token claims with empty email, When IsValid is called, Then should return false",
			claims: auth.TokenClaims{
				UserID:    "user-123",
				Email:     "",
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "Given token claims with zero expires at, When IsValid is called, Then should return false",
			claims: auth.TokenClaims{
				UserID:    "user-123",
				Email:     "test@example.com",
				ExpiresAt: time.Time{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.claims.IsValid()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTokenClaims_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		claims   auth.TokenClaims
		expected bool
	}{
		{
			name: "Given token claims with future expires at, When IsExpired is called, Then should return false",
			claims: auth.TokenClaims{
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "Given token claims with past expires at, When IsExpired is called, Then should return true",
			claims: auth.TokenClaims{
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			expected: true,
		},
		{
			name: "Given token claims with current time expires at, When IsExpired is called, Then should return true",
			claims: auth.TokenClaims{
				ExpiresAt: time.Now(),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.claims.IsExpired()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTokenClaims_IsAccessToken(t *testing.T) {
	tests := []struct {
		name     string
		claims   auth.TokenClaims
		expected bool
	}{
		{
			name: "Given token claims with access token type, When IsAccessToken is called, Then should return true",
			claims: auth.TokenClaims{
				TokenType: "access",
			},
			expected: true,
		},
		{
			name: "Given token claims with refresh token type, When IsAccessToken is called, Then should return false",
			claims: auth.TokenClaims{
				TokenType: "refresh",
			},
			expected: false,
		},
		{
			name: "Given token claims with empty token type, When IsAccessToken is called, Then should return false",
			claims: auth.TokenClaims{
				TokenType: "",
			},
			expected: false,
		},
		{
			name: "Given token claims with other token type, When IsAccessToken is called, Then should return false",
			claims: auth.TokenClaims{
				TokenType: "other",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.claims.IsAccessToken()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTokenClaims_IsRefreshToken(t *testing.T) {
	tests := []struct {
		name     string
		claims   auth.TokenClaims
		expected bool
	}{
		{
			name: "Given token claims with refresh token type, When IsRefreshToken is called, Then should return true",
			claims: auth.TokenClaims{
				TokenType: "refresh",
			},
			expected: true,
		},
		{
			name: "Given token claims with access token type, When IsRefreshToken is called, Then should return false",
			claims: auth.TokenClaims{
				TokenType: "access",
			},
			expected: false,
		},
		{
			name: "Given token claims with empty token type, When IsRefreshToken is called, Then should return false",
			claims: auth.TokenClaims{
				TokenType: "",
			},
			expected: false,
		},
		{
			name: "Given token claims with other token type, When IsRefreshToken is called, Then should return false",
			claims: auth.TokenClaims{
				TokenType: "other",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.claims.IsRefreshToken()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthError_Error(t *testing.T) {
	tests := []struct {
		name     string
		authErr  auth.AuthError
		expected string
	}{
		{
			name: "Given auth error with message, When Error is called, Then should return message",
			authErr: auth.AuthError{
				Code:    "TEST_ERROR",
				Message: "Test error message",
			},
			expected: "Test error message",
		},
		{
			name: "Given auth error with empty message, When Error is called, Then should return empty string",
			authErr: auth.AuthError{
				Code:    "TEST_ERROR",
				Message: "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.authErr.Error()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthErrors_Constants(t *testing.T) {
	tests := []struct {
		name         string
		err          auth.AuthError
		expectedCode string
	}{
		{
			name:         "Given ErrInvalidCredentials, When accessing code, Then should have correct code",
			err:          auth.ErrInvalidCredentials,
			expectedCode: "INVALID_CREDENTIALS",
		},
		{
			name:         "Given ErrUserNotFound, When accessing code, Then should have correct code",
			err:          auth.ErrUserNotFound,
			expectedCode: "USER_NOT_FOUND",
		},
		{
			name:         "Given ErrInvalidToken, When accessing code, Then should have correct code",
			err:          auth.ErrInvalidToken,
			expectedCode: "INVALID_TOKEN",
		},
		{
			name:         "Given ErrTokenExpired, When accessing code, Then should have correct code",
			err:          auth.ErrTokenExpired,
			expectedCode: "TOKEN_EXPIRED",
		},
		{
			name:         "Given ErrUnsupportedStrategy, When accessing code, Then should have correct code",
			err:          auth.ErrUnsupportedStrategy,
			expectedCode: "UNSUPPORTED_STRATEGY",
		},
		{
			name:         "Given ErrInvalidRefreshToken, When accessing code, Then should have correct code",
			err:          auth.ErrInvalidRefreshToken,
			expectedCode: "INVALID_REFRESH_TOKEN",
		},
		{
			name:         "Given ErrUserAlreadyExists, When accessing code, Then should have correct code",
			err:          auth.ErrUserAlreadyExists,
			expectedCode: "USER_EXISTS",
		},
		{
			name:         "Given ErrOAuthProviderNotFound, When accessing code, Then should have correct code",
			err:          auth.ErrOAuthProviderNotFound,
			expectedCode: "OAUTH_PROVIDER_NOT_FOUND",
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

func TestBasicCredentials_Structure(t *testing.T) {
	t.Run("Given basic credentials with email and password, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		credentials := auth.BasicCredentials{
			Email:    "test@example.com",
			Password: "password123",
		}

		// Assert
		assert.Equal(t, "test@example.com", credentials.Email)
		assert.Equal(t, "password123", credentials.Password)
	})
}

func TestOAuthCredentials_Structure(t *testing.T) {
	t.Run("Given OAuth credentials with all fields, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		credentials := auth.OAuthCredentials{
			Provider:     "google",
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			ExpiresIn:    3600,
		}

		// Assert
		assert.Equal(t, "google", credentials.Provider)
		assert.Equal(t, "access-token", credentials.AccessToken)
		assert.Equal(t, "refresh-token", credentials.RefreshToken)
		assert.Equal(t, 3600, credentials.ExpiresIn)
	})
}

func TestJWTCredentials_Structure(t *testing.T) {
	t.Run("Given JWT credentials with token, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		credentials := auth.JWTCredentials{
			Token: "jwt-token",
		}

		// Assert
		assert.Equal(t, "jwt-token", credentials.Token)
	})
}

func TestOAuthUserInfo_Structure(t *testing.T) {
	t.Run("Given OAuth user info with all fields, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		userInfo := auth.OAuthUserInfo{
			ID:        "oauth-user-123",
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Verified:  true,
		}

		// Assert
		assert.Equal(t, "oauth-user-123", userInfo.ID)
		assert.Equal(t, "test@example.com", userInfo.Email)
		assert.Equal(t, "John", userInfo.FirstName)
		assert.Equal(t, "Doe", userInfo.LastName)
		assert.True(t, userInfo.Verified)
	})
}

func TestCreateUserData_Structure(t *testing.T) {
	t.Run("Given create user data with all fields, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		userData := auth.CreateUserData{
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
		}

		// Assert
		assert.Equal(t, "test@example.com", userData.Email)
		assert.Equal(t, "password123", userData.Password)
		assert.Equal(t, "John", userData.FirstName)
		assert.Equal(t, "Doe", userData.LastName)
	})
}