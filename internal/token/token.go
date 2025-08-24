package token

import (
	"context"
	"time"
)

// Service defines the token domain interface - the ONLY interface in this domain
type Service interface {
	// Token generation
	GenerateAuthToken(ctx context.Context, userID string, email string) (string, time.Time, error)
	GenerateRefreshToken(ctx context.Context, userID string) (string, error)
	GenerateAPIToken(ctx context.Context, userID string, scopes []string) (*APIToken, error)
	GeneratePasswordResetToken(ctx context.Context, userID string) (string, error)
	GenerateEmailVerificationToken(ctx context.Context, userID string) (string, error)

	// Token validation
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
	ValidateAPIToken(ctx context.Context, token string) (*APITokenClaims, error)
	ValidatePasswordResetToken(ctx context.Context, token string) (*TokenClaims, error)
	ValidateEmailVerificationToken(ctx context.Context, token string) (*TokenClaims, error)

	// Token management
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
	RevokeToken(ctx context.Context, token string) error
	RevokeAllTokensForUser(ctx context.Context, userID string) error

	// Token introspection
	GetTokenInfo(ctx context.Context, token string) (*TokenInfo, error)
	ListActiveTokens(ctx context.Context, userID string) ([]TokenInfo, error)
}

// Domain types and data structures

// TokenClaims represents the claims in a token
type TokenClaims struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	TokenType string    `json:"token_type"` // auth, refresh, reset, verification
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Issuer    string    `json:"issuer,omitempty"`
	Audience  string    `json:"audience,omitempty"`
	JTI       string    `json:"jti,omitempty"` // JWT ID
}

// APIToken represents an API token with scopes
type APIToken struct {
	ID        string     `json:"id"`
	Token     string     `json:"token"`
	UserID    string     `json:"user_id"`
	Name      string     `json:"name,omitempty"`
	Scopes    []string   `json:"scopes"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt time.Time  `json:"expires_at"`
	LastUsed  *time.Time `json:"last_used,omitempty"`
}

// APITokenClaims represents claims in an API token
type APITokenClaims struct {
	TokenClaims
	Scopes []string `json:"scopes"`
	Name   string   `json:"name,omitempty"`
}

// TokenPair represents an access token and refresh token pair
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"` // "bearer"
	ExpiresIn    int64     `json:"expires_in"` // seconds
	ExpiresAt    time.Time `json:"expires_at"`
	Scope        string    `json:"scope,omitempty"`
}

// TokenInfo contains information about a token
type TokenInfo struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	TokenType string     `json:"token_type"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt time.Time  `json:"expires_at"`
	LastUsed  *time.Time `json:"last_used,omitempty"`
	IsRevoked bool       `json:"is_revoked"`
	Scopes    []string   `json:"scopes,omitempty"`
	UserAgent string     `json:"user_agent,omitempty"`
	IPAddress string     `json:"ip_address,omitempty"`
}

// TokenConfig contains configuration for token service
type TokenConfig struct {
	// JWT configuration
	Secret          []byte        `json:"-"`                // Secret key for signing
	AccessTTL       time.Duration `json:"access_ttl"`       // Access token TTL
	RefreshTTL      time.Duration `json:"refresh_ttl"`      // Refresh token TTL
	ResetTTL        time.Duration `json:"reset_ttl"`        // Password reset token TTL
	VerificationTTL time.Duration `json:"verification_ttl"` // Email verification token TTL

	// Token settings
	Issuer    string `json:"issuer"`    // Token issuer
	Audience  string `json:"audience"`  // Token audience
	Algorithm string `json:"algorithm"` // Signing algorithm (HS256, RS256, etc.)

	// Security settings
	EnableRefresh    bool `json:"enable_refresh"`    // Enable refresh tokens
	EnableRevocation bool `json:"enable_revocation"` // Enable token revocation
	MaxActiveTokens  int  `json:"max_active_tokens"` // Max active tokens per user
}

// TokenError represents domain-specific token errors
type TokenError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func (e TokenError) Error() string {
	return e.Message
}

// Common token error codes
var (
	ErrInvalidToken      = TokenError{Code: "INVALID_TOKEN", Message: "Invalid or expired token"}
	ErrTokenExpired      = TokenError{Code: "TOKEN_EXPIRED", Message: "Token has expired"}
	ErrTokenRevoked      = TokenError{Code: "TOKEN_REVOKED", Message: "Token has been revoked"}
	ErrInvalidSignature  = TokenError{Code: "INVALID_SIGNATURE", Message: "Invalid token signature"}
	ErrMalformedToken    = TokenError{Code: "MALFORMED_TOKEN", Message: "Malformed token"}
	ErrTokenNotFound     = TokenError{Code: "TOKEN_NOT_FOUND", Message: "Token not found"}
	ErrInsufficientScope = TokenError{Code: "INSUFFICIENT_SCOPE", Message: "Insufficient token scope"}
)

// Helper methods for TokenClaims
func (c *TokenClaims) IsValid() bool {
	return c.UserID != "" && !c.ExpiresAt.IsZero()
}

func (c *TokenClaims) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

func (c *TokenClaims) IsAccessToken() bool {
	return c.TokenType == "access" || c.TokenType == "auth"
}

func (c *TokenClaims) IsRefreshToken() bool {
	return c.TokenType == "refresh"
}

func (c *TokenClaims) TimeUntilExpiry() time.Duration {
	return time.Until(c.ExpiresAt)
}

// Helper methods for APIToken
func (t *APIToken) IsValid() bool {
	return t.Token != "" && t.UserID != "" && !t.ExpiresAt.IsZero()
}

func (t *APIToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

func (t *APIToken) HasScope(scope string) bool {
	for _, s := range t.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// Helper methods for TokenPair
func (p *TokenPair) IsValid() bool {
	return p.AccessToken != "" && p.RefreshToken != ""
}

func (p *TokenPair) IsExpired() bool {
	return time.Now().After(p.ExpiresAt)
}

// Helper methods for TokenInfo
func (i *TokenInfo) IsActive() bool {
	return !i.IsRevoked && !i.IsExpired()
}

func (i *TokenInfo) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

// Helper methods for TokenConfig
func (c *TokenConfig) IsValid() bool {
	return len(c.Secret) > 0 && c.AccessTTL > 0 && c.Algorithm != ""
}

// Default token configuration
func DefaultTokenConfig() TokenConfig {
	return TokenConfig{
		AccessTTL:        time.Hour,
		RefreshTTL:       24 * time.Hour,
		ResetTTL:         30 * time.Minute,
		VerificationTTL:  24 * time.Hour,
		Issuer:           "decorator-arch-go",
		Audience:         "api",
		Algorithm:        "HS256",
		EnableRefresh:    true,
		EnableRevocation: true,
		MaxActiveTokens:  10,
	}
}
