package jwt

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/gentra/decorator-arch-go/internal/token"
)

// service implements token.Service interface using JWT
type service struct {
	config        token.TokenConfig
	revokedTokens map[string]time.Time // Simple in-memory revocation list
	mu            sync.RWMutex
}

// NewService creates a new JWT-based token service
func NewService(config token.TokenConfig) (token.Service, error) {
	if !config.IsValid() {
		return nil, fmt.Errorf("invalid token configuration")
	}

	return &service{
		config:        config,
		revokedTokens: make(map[string]time.Time),
	}, nil
}

// GenerateAuthToken generates an authentication token
func (s *service) GenerateAuthToken(ctx context.Context, userID string, email string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(s.config.AccessTTL)
	jti := s.generateJTI(userID, now)

	claims := jwt.MapClaims{
		"user_id":    userID,
		"email":      email,
		"token_type": "auth",
		"iat":        now.Unix(),
		"exp":        expiresAt.Unix(),
		"iss":        s.config.Issuer,
		"aud":        s.config.Audience,
		"jti":        jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.config.Secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// GenerateRefreshToken generates a refresh token
func (s *service) GenerateRefreshToken(ctx context.Context, userID string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.config.RefreshTTL)
	jti := s.generateJTI(userID, now)

	claims := jwt.MapClaims{
		"user_id":    userID,
		"token_type": "refresh",
		"iat":        now.Unix(),
		"exp":        expiresAt.Unix(),
		"iss":        s.config.Issuer,
		"aud":        s.config.Audience,
		"jti":        jti,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(s.config.Secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

// GenerateAPIToken generates an API token with scopes
func (s *service) GenerateAPIToken(ctx context.Context, userID string, scopes []string) (*token.APIToken, error) {
	now := time.Now()
	expiresAt := now.Add(s.config.AccessTTL * 24) // API tokens last longer
	id := uuid.New().String()
	jti := s.generateJTI(userID, now)

	claims := jwt.MapClaims{
		"user_id":    userID,
		"token_type": "api",
		"scopes":     scopes,
		"iat":        now.Unix(),
		"exp":        expiresAt.Unix(),
		"iss":        s.config.Issuer,
		"aud":        s.config.Audience,
		"jti":        jti,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(s.config.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign API token: %w", err)
	}

	return &token.APIToken{
		ID:        id,
		Token:     tokenString,
		UserID:    userID,
		Scopes:    scopes,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}, nil
}

// GeneratePasswordResetToken generates a password reset token
func (s *service) GeneratePasswordResetToken(ctx context.Context, userID string) (string, error) {
	return s.generateSpecialToken(userID, "reset", s.config.ResetTTL)
}

// GenerateEmailVerificationToken generates an email verification token
func (s *service) GenerateEmailVerificationToken(ctx context.Context, userID string) (string, error) {
	return s.generateSpecialToken(userID, "verification", s.config.VerificationTTL)
}

// ValidateToken validates a token and returns claims
func (s *service) ValidateToken(ctx context.Context, tokenString string) (*token.TokenClaims, error) {
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.config.Secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !jwtToken.Valid {
		return nil, token.ErrInvalidToken
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, token.ErrMalformedToken
	}

	// Check if token is revoked
	if jti, ok := claims["jti"].(string); ok {
		if s.isTokenRevoked(jti) {
			return nil, token.ErrTokenRevoked
		}
	}

	// Extract claims
	userID, _ := claims["user_id"].(string)
	email, _ := claims["email"].(string)
	tokenType, _ := claims["token_type"].(string)
	issuer, _ := claims["iss"].(string)
	audience, _ := claims["aud"].(string)
	jti, _ := claims["jti"].(string)

	if userID == "" || tokenType == "" {
		return nil, token.ErrMalformedToken
	}

	issuedAt := time.Unix(int64(claims["iat"].(float64)), 0)
	expiresAt := time.Unix(int64(claims["exp"].(float64)), 0)

	// Check if token is expired
	if time.Now().After(expiresAt) {
		return nil, token.ErrTokenExpired
	}

	return &token.TokenClaims{
		UserID:    userID,
		Email:     email,
		TokenType: tokenType,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		Issuer:    issuer,
		Audience:  audience,
		JTI:       jti,
	}, nil
}

// ValidateAPIToken validates an API token
func (s *service) ValidateAPIToken(ctx context.Context, tokenString string) (*token.APITokenClaims, error) {
	claims, err := s.ValidateToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "api" {
		return nil, token.ErrInvalidToken
	}

	// Parse the token again to get scopes
	jwtToken, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.config.Secret, nil
	})

	jwtClaims := jwtToken.Claims.(jwt.MapClaims)
	scopes, _ := jwtClaims["scopes"].([]interface{})
	scopeStrings := make([]string, len(scopes))
	for i, scope := range scopes {
		scopeStrings[i] = scope.(string)
	}

	return &token.APITokenClaims{
		TokenClaims: *claims,
		Scopes:      scopeStrings,
	}, nil
}

// ValidatePasswordResetToken validates a password reset token
func (s *service) ValidatePasswordResetToken(ctx context.Context, tokenString string) (*token.TokenClaims, error) {
	claims, err := s.ValidateToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "reset" {
		return nil, token.ErrInvalidToken
	}

	return claims, nil
}

// ValidateEmailVerificationToken validates an email verification token
func (s *service) ValidateEmailVerificationToken(ctx context.Context, tokenString string) (*token.TokenClaims, error) {
	claims, err := s.ValidateToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "verification" {
		return nil, token.ErrInvalidToken
	}

	return claims, nil
}

// RefreshToken generates a new access token from a refresh token
func (s *service) RefreshToken(ctx context.Context, refreshToken string) (*token.TokenPair, error) {
	claims, err := s.ValidateToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if !claims.IsRefreshToken() {
		return nil, token.ErrInvalidToken
	}

	// Generate new access token
	accessToken, expiresAt, err := s.GenerateAuthToken(ctx, claims.UserID, claims.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &token.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken, // Keep the same refresh token
		TokenType:    "bearer",
		ExpiresIn:    int64(s.config.AccessTTL.Seconds()),
		ExpiresAt:    expiresAt,
	}, nil
}

// RevokeToken revokes a token
func (s *service) RevokeToken(ctx context.Context, tokenString string) error {
	// Parse token to get JTI
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.config.Secret, nil
	})

	if err != nil {
		return fmt.Errorf("failed to parse token for revocation: %w", err)
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return token.ErrMalformedToken
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return fmt.Errorf("token missing JTI claim")
	}

	expiresAt := time.Unix(int64(claims["exp"].(float64)), 0)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Add to revocation list with expiration time
	s.revokedTokens[jti] = expiresAt

	// Clean up expired revoked tokens
	s.cleanupExpiredRevokedTokens()

	return nil
}

// RevokeAllTokensForUser revokes all tokens for a user (placeholder)
func (s *service) RevokeAllTokensForUser(ctx context.Context, userID string) error {
	// In a real implementation, this would query all tokens for the user and revoke them
	// For this JWT implementation, we'd need to maintain a list of active tokens per user
	return nil
}

// GetTokenInfo returns information about a token
func (s *service) GetTokenInfo(ctx context.Context, tokenString string) (*token.TokenInfo, error) {
	claims, err := s.ValidateToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	return &token.TokenInfo{
		ID:        claims.JTI,
		UserID:    claims.UserID,
		TokenType: claims.TokenType,
		CreatedAt: claims.IssuedAt,
		ExpiresAt: claims.ExpiresAt,
		IsRevoked: s.isTokenRevoked(claims.JTI),
	}, nil
}

// ListActiveTokens lists active tokens for a user (placeholder)
func (s *service) ListActiveTokens(ctx context.Context, userID string) ([]token.TokenInfo, error) {
	// In a real implementation, this would query active tokens for the user
	// For this JWT implementation, we'd need to maintain a registry of active tokens
	return []token.TokenInfo{}, nil
}

// Helper methods

func (s *service) generateSpecialToken(userID, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	expiresAt := now.Add(ttl)
	jti := s.generateJTI(userID, now)

	claims := jwt.MapClaims{
		"user_id":    userID,
		"token_type": tokenType,
		"iat":        now.Unix(),
		"exp":        expiresAt.Unix(),
		"iss":        s.config.Issuer,
		"aud":        s.config.Audience,
		"jti":        jti,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(s.config.Secret)
}

func (s *service) generateJTI(userID string, issuedAt time.Time) string {
	return fmt.Sprintf("%s-%d", userID, issuedAt.Unix())
}

func (s *service) isTokenRevoked(jti string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expiresAt, exists := s.revokedTokens[jti]
	if !exists {
		return false
	}

	// If the revoked token has expired, it's no longer relevant
	if time.Now().After(expiresAt) {
		return false
	}

	return true
}

func (s *service) cleanupExpiredRevokedTokens() {
	now := time.Now()
	for jti, expiresAt := range s.revokedTokens {
		if now.After(expiresAt) {
			delete(s.revokedTokens, jti)
		}
	}
}
