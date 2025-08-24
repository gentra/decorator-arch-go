package usecase

import (
	"context"
	"fmt"

	"github.com/gentra/decorator-arch-go/internal/auth"
	"github.com/gentra/decorator-arch-go/internal/user"
)

// OAuthAuthStrategy implements auth.Service for OAuth authentication
type OAuthAuthStrategy struct {
	userService    user.Service
	tokenManager   *JWTTokenManager
	oauthProviders map[string]auth.Service // OAuth providers implement auth.Service
}

// NewOAuthAuthStrategy creates a new OAuth authentication strategy
func NewOAuthAuthStrategy(userService user.Service, tokenManager *JWTTokenManager, oauthProviders map[string]auth.Service) auth.Service {
	return &OAuthAuthStrategy{
		userService:    userService,
		tokenManager:   tokenManager,
		oauthProviders: oauthProviders,
	}
}

// Authenticate handles only "oauth" strategy
func (s *OAuthAuthStrategy) Authenticate(ctx context.Context, strategy string, credentials interface{}) (*auth.AuthResult, error) {
	if strategy != "oauth" {
		return nil, auth.ErrUnsupportedStrategy
	}

	oauthCreds, ok := credentials.(auth.OAuthCredentials)
	if !ok {
		return nil, fmt.Errorf("invalid credentials type for OAuth")
	}

	provider, exists := s.oauthProviders[oauthCreds.Provider]
	if !exists {
		return nil, auth.ErrOAuthProviderNotFound
	}

	// Delegate to OAuth provider service
	return provider.Authenticate(ctx, "oauth", credentials)
}

// ValidateToken delegates to token manager
func (s *OAuthAuthStrategy) ValidateToken(ctx context.Context, token string) (*auth.TokenClaims, error) {
	return s.tokenManager.ValidateToken(token)
}

// RefreshToken delegates to token manager
func (s *OAuthAuthStrategy) RefreshToken(ctx context.Context, refreshToken string) (*auth.AuthResult, error) {
	// Validate refresh token
	claims, err := s.tokenManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if !claims.IsRefreshToken() {
		return nil, auth.ErrInvalidRefreshToken
	}

	// Generate new access token
	accessToken, expiresAt, err := s.tokenManager.GenerateAuthToken(claims.UserID, claims.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Create user from claims
	authUser := &auth.User{
		ID:    claims.UserID,
		Email: claims.Email,
	}

	return &auth.AuthResult{
		User:         authUser,
		Token:        accessToken,
		RefreshToken: refreshToken, // Keep the same refresh token
		ExpiresAt:    expiresAt,
		Strategy:     "oauth",
	}, nil
}

// RevokeToken delegates to token manager
func (s *OAuthAuthStrategy) RevokeToken(ctx context.Context, token string) error {
	return s.tokenManager.RevokeToken(token)
}

// GetSupportedStrategies returns oauth strategy
func (s *OAuthAuthStrategy) GetSupportedStrategies() []string {
	return []string{"oauth"}
}
