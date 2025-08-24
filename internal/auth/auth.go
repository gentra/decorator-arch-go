package auth

import (
	"context"
	"time"
)

// Service defines the authentication domain interface - the ONLY interface in this domain
type Service interface {
	// Authentication operations
	Authenticate(ctx context.Context, strategy string, credentials interface{}) (*AuthResult, error)
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResult, error)
	RevokeToken(ctx context.Context, token string) error

	// Service capabilities
	GetSupportedStrategies() []string
}

// Domain types and data structures

// AuthResult contains authentication result data
type AuthResult struct {
	User         *User     `json:"user"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
	Strategy     string    `json:"strategy"`
}

// TokenClaims represents the claims in an authentication token
type TokenClaims struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
	TokenType string    `json:"token_type"` // "access" or "refresh"
	Strategy  string    `json:"strategy"`   // "basic", "oauth", etc.
}

// User represents a user for authentication purposes
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Credentials for different authentication methods

// BasicCredentials for username/password authentication
type BasicCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// OAuthCredentials for OAuth authentication
type OAuthCredentials struct {
	Provider     string `json:"provider"` // "google", "github", etc.
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// JWTCredentials for JWT token authentication
type JWTCredentials struct {
	Token string `json:"token"`
}

// OAuth provider data structures

// OAuthUserInfo contains user information from OAuth provider
type OAuthUserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Verified  bool   `json:"verified"`
}

// User provider data structures (for integration with user domain)

// CreateUserData contains data for creating a new user
type CreateUserData struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// AuthError represents domain-specific authentication errors
type AuthError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func (e AuthError) Error() string {
	return e.Message
}

// Common authentication error codes
var (
	ErrInvalidCredentials    = AuthError{Code: "INVALID_CREDENTIALS", Message: "Invalid email or password"}
	ErrUserNotFound          = AuthError{Code: "USER_NOT_FOUND", Message: "User not found"}
	ErrInvalidToken          = AuthError{Code: "INVALID_TOKEN", Message: "Invalid or expired token"}
	ErrTokenExpired          = AuthError{Code: "TOKEN_EXPIRED", Message: "Token has expired"}
	ErrUnsupportedStrategy   = AuthError{Code: "UNSUPPORTED_STRATEGY", Message: "Authentication strategy not supported"}
	ErrInvalidRefreshToken   = AuthError{Code: "INVALID_REFRESH_TOKEN", Message: "Invalid refresh token"}
	ErrUserAlreadyExists     = AuthError{Code: "USER_EXISTS", Message: "User already exists"}
	ErrOAuthProviderNotFound = AuthError{Code: "OAUTH_PROVIDER_NOT_FOUND", Message: "OAuth provider not configured"}
)

// Helper methods for domain types

// Helper methods for User
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) IsValid() bool {
	return u.ID != "" && u.Email != ""
}

// Helper methods for AuthResult
func (r *AuthResult) IsValid() bool {
	return r.User != nil && r.Token != "" && !r.ExpiresAt.IsZero()
}

func (r *AuthResult) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// Helper methods for TokenClaims
func (c *TokenClaims) IsValid() bool {
	return c.UserID != "" && c.Email != "" && !c.ExpiresAt.IsZero()
}

func (c *TokenClaims) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

func (c *TokenClaims) IsAccessToken() bool {
	return c.TokenType == "access"
}

func (c *TokenClaims) IsRefreshToken() bool {
	return c.TokenType == "refresh"
}
