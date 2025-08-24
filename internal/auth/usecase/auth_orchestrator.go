package usecase

import (
	"context"
	"fmt"

	"github.com/gentra/decorator-arch-go/internal/auth"
)

// AuthOrchestrator implements auth.Service and orchestrates different authentication strategies
// This contains the core business logic for authentication management
type AuthOrchestrator struct {
	tokenManager    *JWTTokenManager
	strategyManager *StrategyManager
}

// NewAuthOrchestrator creates a new authentication orchestrator
func NewAuthOrchestrator(tokenManager *JWTTokenManager) *AuthOrchestrator {
	return &AuthOrchestrator{
		tokenManager:    tokenManager,
		strategyManager: NewStrategyManager(),
	}
}

// RegisterStrategy registers an authentication strategy
func (s *AuthOrchestrator) RegisterStrategy(name string, strategy auth.Service) {
	s.strategyManager.RegisterStrategy(name, strategy)
}

// Authenticate handles authentication by delegating to the appropriate strategy
func (s *AuthOrchestrator) Authenticate(ctx context.Context, strategy string, credentials interface{}) (*auth.AuthResult, error) {
	return s.strategyManager.Authenticate(ctx, strategy, credentials)
}

// ValidateToken validates an authentication token
func (s *AuthOrchestrator) ValidateToken(ctx context.Context, token string) (*auth.TokenClaims, error) {
	return s.tokenManager.ValidateToken(token)
}

// RefreshToken generates a new access token using a refresh token
func (s *AuthOrchestrator) RefreshToken(ctx context.Context, refreshToken string) (*auth.AuthResult, error) {
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

	// Create user from claims (simplified)
	authUser := &auth.User{
		ID:    claims.UserID,
		Email: claims.Email,
	}

	return &auth.AuthResult{
		User:         authUser,
		Token:        accessToken,
		RefreshToken: refreshToken, // Keep the same refresh token
		ExpiresAt:    expiresAt,
		Strategy:     claims.Strategy,
	}, nil
}

// RevokeToken revokes an authentication token
func (s *AuthOrchestrator) RevokeToken(ctx context.Context, token string) error {
	return s.tokenManager.RevokeToken(token)
}

// GetSupportedStrategies returns the list of supported authentication strategies
func (s *AuthOrchestrator) GetSupportedStrategies() []string {
	return s.strategyManager.GetSupportedStrategies()
}

// StrategyManager manages authentication strategies - this is core business logic
type StrategyManager struct {
	strategies map[string]auth.Service
}

// NewStrategyManager creates a new strategy manager
func NewStrategyManager() *StrategyManager {
	return &StrategyManager{
		strategies: make(map[string]auth.Service),
	}
}

// RegisterStrategy registers an authentication strategy
func (sm *StrategyManager) RegisterStrategy(name string, strategy auth.Service) {
	sm.strategies[name] = strategy
}

// Authenticate handles authentication using the specified strategy
func (sm *StrategyManager) Authenticate(ctx context.Context, strategyName string, credentials interface{}) (*auth.AuthResult, error) {
	strategy, exists := sm.strategies[strategyName]
	if !exists {
		return nil, auth.ErrUnsupportedStrategy
	}

	return strategy.Authenticate(ctx, strategyName, credentials)
}

// GetSupportedStrategies returns all registered strategy names
func (sm *StrategyManager) GetSupportedStrategies() []string {
	strategies := make([]string, 0, len(sm.strategies))
	for strategyName := range sm.strategies {
		strategies = append(strategies, strategyName)
	}
	return strategies
}
