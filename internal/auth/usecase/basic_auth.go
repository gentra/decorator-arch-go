package usecase

import (
	"context"
	"fmt"

	"github.com/gentra/decorator-arch-go/internal/auth"
	"github.com/gentra/decorator-arch-go/internal/user"
)

// BasicAuthStrategy implements auth.Service for basic username/password authentication
type BasicAuthStrategy struct {
	userService  user.Service
	tokenManager *JWTTokenManager // Will move this to usecase package
}

// NewBasicAuthStrategy creates a new basic authentication strategy
func NewBasicAuthStrategy(userService user.Service, tokenManager *JWTTokenManager) auth.Service {
	return &BasicAuthStrategy{
		userService:  userService,
		tokenManager: tokenManager,
	}
}

// Authenticate handles only "basic" strategy
func (s *BasicAuthStrategy) Authenticate(ctx context.Context, strategy string, credentials interface{}) (*auth.AuthResult, error) {
	if strategy != "basic" {
		return nil, auth.ErrUnsupportedStrategy
	}

	basicCreds, ok := credentials.(auth.BasicCredentials)
	if !ok {
		return nil, fmt.Errorf("invalid credentials type for basic auth")
	}

	// Use user service to validate credentials
	authResult, err := s.userService.Login(ctx, basicCreds.Email, basicCreds.Password)
	if err != nil {
		return nil, auth.ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, expiresAt, err := s.tokenManager.GenerateAuthToken(authResult.User.ID.String(), authResult.User.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(authResult.User.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &auth.AuthResult{
		User:         convertUserDomainToAuth(authResult.User),
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		Strategy:     "basic",
	}, nil
}

// ValidateToken delegates to token manager
func (s *BasicAuthStrategy) ValidateToken(ctx context.Context, token string) (*auth.TokenClaims, error) {
	return s.tokenManager.ValidateToken(token)
}

// RefreshToken delegates to token manager
func (s *BasicAuthStrategy) RefreshToken(ctx context.Context, refreshToken string) (*auth.AuthResult, error) {
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
		Strategy:     "basic",
	}, nil
}

// RevokeToken delegates to token manager
func (s *BasicAuthStrategy) RevokeToken(ctx context.Context, token string) error {
	return s.tokenManager.RevokeToken(token)
}

// GetSupportedStrategies returns only basic auth
func (s *BasicAuthStrategy) GetSupportedStrategies() []string {
	return []string{"basic"}
}

// Helper function to convert user domain to auth domain
func convertUserDomainToAuth(userDomainUser *user.User) *auth.User {
	if userDomainUser == nil {
		return nil
	}

	return &auth.User{
		ID:           userDomainUser.ID.String(),
		Email:        userDomainUser.Email,
		FirstName:    userDomainUser.FirstName,
		LastName:     userDomainUser.LastName,
		PasswordHash: userDomainUser.PasswordHash,
		CreatedAt:    userDomainUser.CreatedAt,
		UpdatedAt:    userDomainUser.UpdatedAt,
	}
}
