package token_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/token"
)

func TestTokenClaims_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		claims   token.TokenClaims
		expected bool
	}{
		{
			name: "Given token claims with user ID and expires at, When IsValid is called, Then should return true",
			claims: token.TokenClaims{
				UserID:    "user-123",
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: true,
		},
		{
			name: "Given token claims with empty user ID, When IsValid is called, Then should return false",
			claims: token.TokenClaims{
				UserID:    "",
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "Given token claims with zero expires at, When IsValid is called, Then should return false",
			claims: token.TokenClaims{
				UserID:    "user-123",
				ExpiresAt: time.Time{},
			},
			expected: false,
		},
		{
			name: "Given token claims with both user ID and expires at empty, When IsValid is called, Then should return false",
			claims: token.TokenClaims{
				UserID:    "",
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
		claims   token.TokenClaims
		expected bool
	}{
		{
			name: "Given token claims with future expires at, When IsExpired is called, Then should return false",
			claims: token.TokenClaims{
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "Given token claims with past expires at, When IsExpired is called, Then should return true",
			claims: token.TokenClaims{
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			expected: true,
		},
		{
			name: "Given token claims with current time expires at, When IsExpired is called, Then should return true",
			claims: token.TokenClaims{
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
		claims   token.TokenClaims
		expected bool
	}{
		{
			name: "Given token claims with access token type, When IsAccessToken is called, Then should return true",
			claims: token.TokenClaims{
				TokenType: "access",
			},
			expected: true,
		},
		{
			name: "Given token claims with auth token type, When IsAccessToken is called, Then should return true",
			claims: token.TokenClaims{
				TokenType: "auth",
			},
			expected: true,
		},
		{
			name: "Given token claims with refresh token type, When IsAccessToken is called, Then should return false",
			claims: token.TokenClaims{
				TokenType: "refresh",
			},
			expected: false,
		},
		{
			name: "Given token claims with empty token type, When IsAccessToken is called, Then should return false",
			claims: token.TokenClaims{
				TokenType: "",
			},
			expected: false,
		},
		{
			name: "Given token claims with other token type, When IsAccessToken is called, Then should return false",
			claims: token.TokenClaims{
				TokenType: "reset",
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
		claims   token.TokenClaims
		expected bool
	}{
		{
			name: "Given token claims with refresh token type, When IsRefreshToken is called, Then should return true",
			claims: token.TokenClaims{
				TokenType: "refresh",
			},
			expected: true,
		},
		{
			name: "Given token claims with access token type, When IsRefreshToken is called, Then should return false",
			claims: token.TokenClaims{
				TokenType: "access",
			},
			expected: false,
		},
		{
			name: "Given token claims with auth token type, When IsRefreshToken is called, Then should return false",
			claims: token.TokenClaims{
				TokenType: "auth",
			},
			expected: false,
		},
		{
			name: "Given token claims with empty token type, When IsRefreshToken is called, Then should return false",
			claims: token.TokenClaims{
				TokenType: "",
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

func TestTokenClaims_TimeUntilExpiry(t *testing.T) {
	tests := []struct {
		name     string
		claims   token.TokenClaims
		expected bool // true if positive, false if negative/zero
	}{
		{
			name: "Given token claims with future expires at, When TimeUntilExpiry is called, Then should return positive duration",
			claims: token.TokenClaims{
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: true,
		},
		{
			name: "Given token claims with past expires at, When TimeUntilExpiry is called, Then should return negative duration",
			claims: token.TokenClaims{
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.claims.TimeUntilExpiry()

			// Assert
			if tt.expected {
				assert.True(t, result > 0)
			} else {
				assert.True(t, result <= 0)
			}
		})
	}
}

func TestAPIToken_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		apiToken token.APIToken
		expected bool
	}{
		{
			name: "Given API token with token, user ID and expires at, When IsValid is called, Then should return true",
			apiToken: token.APIToken{
				Token:     "api-token-123",
				UserID:    "user-123",
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: true,
		},
		{
			name: "Given API token with empty token, When IsValid is called, Then should return false",
			apiToken: token.APIToken{
				Token:     "",
				UserID:    "user-123",
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "Given API token with empty user ID, When IsValid is called, Then should return false",
			apiToken: token.APIToken{
				Token:     "api-token-123",
				UserID:    "",
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "Given API token with zero expires at, When IsValid is called, Then should return false",
			apiToken: token.APIToken{
				Token:     "api-token-123",
				UserID:    "user-123",
				ExpiresAt: time.Time{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.apiToken.IsValid()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAPIToken_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		apiToken token.APIToken
		expected bool
	}{
		{
			name: "Given API token with future expires at, When IsExpired is called, Then should return false",
			apiToken: token.APIToken{
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "Given API token with past expires at, When IsExpired is called, Then should return true",
			apiToken: token.APIToken{
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			expected: true,
		},
		{
			name: "Given API token with current time expires at, When IsExpired is called, Then should return true",
			apiToken: token.APIToken{
				ExpiresAt: time.Now(),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.apiToken.IsExpired()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAPIToken_HasScope(t *testing.T) {
	tests := []struct {
		name     string
		apiToken token.APIToken
		scope    string
		expected bool
	}{
		{
			name: "Given API token with multiple scopes and existing scope, When HasScope is called, Then should return true",
			apiToken: token.APIToken{
				Scopes: []string{"read", "write", "admin"},
			},
			scope:    "read",
			expected: true,
		},
		{
			name: "Given API token with multiple scopes and non-existing scope, When HasScope is called, Then should return false",
			apiToken: token.APIToken{
				Scopes: []string{"read", "write"},
			},
			scope:    "admin",
			expected: false,
		},
		{
			name: "Given API token with empty scopes, When HasScope is called, Then should return false",
			apiToken: token.APIToken{
				Scopes: []string{},
			},
			scope:    "read",
			expected: false,
		},
		{
			name: "Given API token with nil scopes, When HasScope is called, Then should return false",
			apiToken: token.APIToken{
				Scopes: nil,
			},
			scope:    "read",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.apiToken.HasScope(tt.scope)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTokenPair_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		tokenPair token.TokenPair
		expected  bool
	}{
		{
			name: "Given token pair with access and refresh tokens, When IsValid is called, Then should return true",
			tokenPair: token.TokenPair{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
			},
			expected: true,
		},
		{
			name: "Given token pair with empty access token, When IsValid is called, Then should return false",
			tokenPair: token.TokenPair{
				AccessToken:  "",
				RefreshToken: "refresh-token",
			},
			expected: false,
		},
		{
			name: "Given token pair with empty refresh token, When IsValid is called, Then should return false",
			tokenPair: token.TokenPair{
				AccessToken:  "access-token",
				RefreshToken: "",
			},
			expected: false,
		},
		{
			name: "Given token pair with both tokens empty, When IsValid is called, Then should return false",
			tokenPair: token.TokenPair{
				AccessToken:  "",
				RefreshToken: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.tokenPair.IsValid()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTokenPair_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		tokenPair token.TokenPair
		expected  bool
	}{
		{
			name: "Given token pair with future expires at, When IsExpired is called, Then should return false",
			tokenPair: token.TokenPair{
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "Given token pair with past expires at, When IsExpired is called, Then should return true",
			tokenPair: token.TokenPair{
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			expected: true,
		},
		{
			name: "Given token pair with current time expires at, When IsExpired is called, Then should return true",
			tokenPair: token.TokenPair{
				ExpiresAt: time.Now(),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.tokenPair.IsExpired()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTokenInfo_IsActive(t *testing.T) {
	tests := []struct {
		name      string
		tokenInfo token.TokenInfo
		expected  bool
	}{
		{
			name: "Given token info not revoked and not expired, When IsActive is called, Then should return true",
			tokenInfo: token.TokenInfo{
				IsRevoked: false,
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: true,
		},
		{
			name: "Given token info revoked but not expired, When IsActive is called, Then should return false",
			tokenInfo: token.TokenInfo{
				IsRevoked: true,
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "Given token info not revoked but expired, When IsActive is called, Then should return false",
			tokenInfo: token.TokenInfo{
				IsRevoked: false,
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			expected: false,
		},
		{
			name: "Given token info revoked and expired, When IsActive is called, Then should return false",
			tokenInfo: token.TokenInfo{
				IsRevoked: true,
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.tokenInfo.IsActive()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTokenInfo_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		tokenInfo token.TokenInfo
		expected  bool
	}{
		{
			name: "Given token info with future expires at, When IsExpired is called, Then should return false",
			tokenInfo: token.TokenInfo{
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "Given token info with past expires at, When IsExpired is called, Then should return true",
			tokenInfo: token.TokenInfo{
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			expected: true,
		},
		{
			name: "Given token info with current time expires at, When IsExpired is called, Then should return true",
			tokenInfo: token.TokenInfo{
				ExpiresAt: time.Now(),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.tokenInfo.IsExpired()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTokenConfig_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		config token.TokenConfig
		expected bool
	}{
		{
			name: "Given token config with secret, access TTL and algorithm, When IsValid is called, Then should return true",
			config: token.TokenConfig{
				Secret:    []byte("secret-key"),
				AccessTTL: time.Hour,
				Algorithm: "HS256",
			},
			expected: true,
		},
		{
			name: "Given token config with empty secret, When IsValid is called, Then should return false",
			config: token.TokenConfig{
				Secret:    []byte{},
				AccessTTL: time.Hour,
				Algorithm: "HS256",
			},
			expected: false,
		},
		{
			name: "Given token config with zero access TTL, When IsValid is called, Then should return false",
			config: token.TokenConfig{
				Secret:    []byte("secret-key"),
				AccessTTL: 0,
				Algorithm: "HS256",
			},
			expected: false,
		},
		{
			name: "Given token config with empty algorithm, When IsValid is called, Then should return false",
			config: token.TokenConfig{
				Secret:    []byte("secret-key"),
				AccessTTL: time.Hour,
				Algorithm: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.config.IsValid()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultTokenConfig(t *testing.T) {
	t.Run("Given default token config call, When DefaultTokenConfig is called, Then should return valid default configuration", func(t *testing.T) {
		// Act
		config := token.DefaultTokenConfig()

		// Assert
		assert.Equal(t, time.Hour, config.AccessTTL)
		assert.Equal(t, 24*time.Hour, config.RefreshTTL)
		assert.Equal(t, 30*time.Minute, config.ResetTTL)
		assert.Equal(t, 24*time.Hour, config.VerificationTTL)
		assert.Equal(t, "decorator-arch-go", config.Issuer)
		assert.Equal(t, "api", config.Audience)
		assert.Equal(t, "HS256", config.Algorithm)
		assert.True(t, config.EnableRefresh)
		assert.True(t, config.EnableRevocation)
		assert.Equal(t, 10, config.MaxActiveTokens)
	})
}

func TestTokenError_Error(t *testing.T) {
	tests := []struct {
		name     string
		tokenErr token.TokenError
		expected string
	}{
		{
			name: "Given token error with message, When Error is called, Then should return message",
			tokenErr: token.TokenError{
				Code:    "TEST_ERROR",
				Message: "Test error message",
			},
			expected: "Test error message",
		},
		{
			name: "Given token error with empty message, When Error is called, Then should return empty string",
			tokenErr: token.TokenError{
				Code:    "TEST_ERROR",
				Message: "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.tokenErr.Error()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTokenErrors_Constants(t *testing.T) {
	tests := []struct {
		name         string
		err          token.TokenError
		expectedCode string
	}{
		{
			name:         "Given ErrInvalidToken, When accessing code, Then should have correct code",
			err:          token.ErrInvalidToken,
			expectedCode: "INVALID_TOKEN",
		},
		{
			name:         "Given ErrTokenExpired, When accessing code, Then should have correct code",
			err:          token.ErrTokenExpired,
			expectedCode: "TOKEN_EXPIRED",
		},
		{
			name:         "Given ErrTokenRevoked, When accessing code, Then should have correct code",
			err:          token.ErrTokenRevoked,
			expectedCode: "TOKEN_REVOKED",
		},
		{
			name:         "Given ErrInvalidSignature, When accessing code, Then should have correct code",
			err:          token.ErrInvalidSignature,
			expectedCode: "INVALID_SIGNATURE",
		},
		{
			name:         "Given ErrMalformedToken, When accessing code, Then should have correct code",
			err:          token.ErrMalformedToken,
			expectedCode: "MALFORMED_TOKEN",
		},
		{
			name:         "Given ErrTokenNotFound, When accessing code, Then should have correct code",
			err:          token.ErrTokenNotFound,
			expectedCode: "TOKEN_NOT_FOUND",
		},
		{
			name:         "Given ErrInsufficientScope, When accessing code, Then should have correct code",
			err:          token.ErrInsufficientScope,
			expectedCode: "INSUFFICIENT_SCOPE",
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