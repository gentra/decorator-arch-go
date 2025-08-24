package usecase

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gentra/decorator-arch-go/internal/auth"
)

// JWTTokenManager handles JWT token operations (moved from factory to usecase)
type JWTTokenManager struct {
	secret        []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
	revokedTokens map[string]time.Time // Simple in-memory revocation list
	mu            sync.RWMutex
}

// NewJWTTokenManager creates a new JWT token manager
func NewJWTTokenManager(secret []byte, accessTTL, refreshTTL time.Duration) *JWTTokenManager {
	return &JWTTokenManager{
		secret:        secret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		revokedTokens: make(map[string]time.Time),
	}
}

func (tm *JWTTokenManager) GenerateAuthToken(userID string, email string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(tm.accessTTL)

	claims := jwt.MapClaims{
		"user_id":    userID,
		"email":      email,
		"token_type": "access",
		"iat":        now.Unix(),
		"exp":        expiresAt.Unix(),
		"jti":        tm.generateJTI(userID, now, "access"),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tm.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

func (tm *JWTTokenManager) GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(tm.refreshTTL)

	claims := jwt.MapClaims{
		"user_id":    userID,
		"token_type": "refresh",
		"iat":        now.Unix(),
		"exp":        expiresAt.Unix(),
		"jti":        tm.generateJTI(userID, now, "refresh"),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tm.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

func (tm *JWTTokenManager) ValidateToken(tokenString string) (*auth.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tm.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, auth.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, auth.ErrInvalidToken
	}

	// Check if token is revoked
	if jti, ok := claims["jti"].(string); ok {
		if tm.isTokenRevoked(jti) {
			return nil, auth.ErrInvalidToken
		}
	}

	// Extract claims
	userID, _ := claims["user_id"].(string)
	email, _ := claims["email"].(string)
	tokenType, _ := claims["token_type"].(string)

	if userID == "" || tokenType == "" {
		return nil, auth.ErrInvalidToken
	}

	issuedAt := time.Unix(int64(claims["iat"].(float64)), 0)
	expiresAt := time.Unix(int64(claims["exp"].(float64)), 0)

	// Check if token is expired
	if time.Now().After(expiresAt) {
		return nil, auth.ErrTokenExpired
	}

	return &auth.TokenClaims{
		UserID:    userID,
		Email:     email,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		TokenType: tokenType,
		Strategy:  "jwt",
	}, nil
}

func (tm *JWTTokenManager) RevokeToken(tokenString string) error {
	// Parse token to get JTI
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tm.secret, nil
	})

	if err != nil {
		return fmt.Errorf("failed to parse token for revocation: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return auth.ErrInvalidToken
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return fmt.Errorf("token missing JTI claim")
	}

	expiresAt := time.Unix(int64(claims["exp"].(float64)), 0)

	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Add to revocation list with expiration time
	tm.revokedTokens[jti] = expiresAt

	// Clean up expired revoked tokens
	tm.cleanupExpiredRevokedTokens()

	return nil
}

func (tm *JWTTokenManager) isTokenRevoked(jti string) bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	expiresAt, exists := tm.revokedTokens[jti]
	if !exists {
		return false
	}

	// If the revoked token has expired, it's no longer relevant
	if time.Now().After(expiresAt) {
		return false
	}

	return true
}

func (tm *JWTTokenManager) cleanupExpiredRevokedTokens() {
	now := time.Now()
	for jti, expiresAt := range tm.revokedTokens {
		if now.After(expiresAt) {
			delete(tm.revokedTokens, jti)
		}
	}
}

func (tm *JWTTokenManager) generateJTI(userID string, issuedAt time.Time, tokenType string) string {
	return fmt.Sprintf("%s-%s-%d", userID, tokenType, issuedAt.Unix())
}
