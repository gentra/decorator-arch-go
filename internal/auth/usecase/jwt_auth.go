package usecase

import (
	"context"
	"fmt"

	"github.com/gentra/decorator-arch-go/internal/auth"
	"github.com/gentra/decorator-arch-go/internal/user"
)

// JWTAuthStrategy implements auth.Service for JWT token authentication
type JWTAuthStrategy struct {
	userService  user.Service
	tokenManager *JWTTokenManager
}

// NewJWTAuthStrategy creates a new JWT authentication strategy
func NewJWTAuthStrategy(userService user.Service, tokenManager *JWTTokenManager) auth.Service {
	return &JWTAuthStrategy{
		userService:  userService,
		tokenManager: tokenManager,
	}
}

// Authenticate handles only "jwt" strategy
func (s *JWTAuthStrategy) Authenticate(ctx context.Context, strategy string, credentials interface{}) (*auth.AuthResult, error) {
	if strategy != "jwt" {
		return nil, auth.ErrUnsupportedStrategy
	}

	jwtCreds, ok := credentials.(auth.JWTCredentials)
	if !ok {
		return nil, fmt.Errorf("invalid credentials type for JWT auth")
	}

	// Validate token
	claims, err := s.tokenManager.ValidateToken(jwtCreds.Token)
	if err != nil {
		return nil, err
	}

	// Get user by ID
	userDomainUser, err := s.userService.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &auth.AuthResult{
		User:      convertUserDomainToAuth(userDomainUser),
		Token:     jwtCreds.Token,
		ExpiresAt: claims.ExpiresAt,
		Strategy:  "jwt",
	}, nil
}

// ValidateToken delegates to token manager
func (s *JWTAuthStrategy) ValidateToken(ctx context.Context, token string) (*auth.TokenClaims, error) {
	return s.tokenManager.ValidateToken(token)
}

// RefreshToken delegates to token manager
func (s *JWTAuthStrategy) RefreshToken(ctx context.Context, refreshToken string) (*auth.AuthResult, error) {
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
		Strategy:     "jwt",
	}, nil
}

// RevokeToken delegates to token manager
func (s *JWTAuthStrategy) RevokeToken(ctx context.Context, token string) error {
	return s.tokenManager.RevokeToken(token)
}

// GetSupportedStrategies returns jwt strategy
func (s *JWTAuthStrategy) GetSupportedStrategies() []string {
	return []string{"jwt"}
}
