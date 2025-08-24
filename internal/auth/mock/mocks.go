package mock

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/gentra/decorator-arch-go/internal/auth"
	usermock "github.com/gentra/decorator-arch-go/internal/user/mock"
)

// MockUserService alias to centralized mock for backward compatibility
type MockUserService = usermock.MockUserService

// MockJWTTokenManager for testing auth strategies
type MockJWTTokenManager struct {
	mock.Mock
}

func (m *MockJWTTokenManager) GenerateAuthToken(userID string, email string) (string, time.Time, error) {
	args := m.Called(userID, email)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockJWTTokenManager) GenerateRefreshToken(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockJWTTokenManager) ValidateToken(tokenString string) (*auth.TokenClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenClaims), args.Error(1)
}

func (m *MockJWTTokenManager) RevokeToken(tokenString string) error {
	args := m.Called(tokenString)
	return args.Error(0)
}

// MockAuthStrategy is a mock implementation of auth.Service for testing
type MockAuthStrategy struct {
	mock.Mock
}

func (m *MockAuthStrategy) Authenticate(ctx context.Context, strategy string, credentials interface{}) (*auth.AuthResult, error) {
	args := m.Called(ctx, strategy, credentials)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.AuthResult), args.Error(1)
}

func (m *MockAuthStrategy) ValidateToken(ctx context.Context, token string) (*auth.TokenClaims, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenClaims), args.Error(1)
}

func (m *MockAuthStrategy) RefreshToken(ctx context.Context, refreshToken string) (*auth.AuthResult, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.AuthResult), args.Error(1)
}

func (m *MockAuthStrategy) RevokeToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAuthStrategy) GetSupportedStrategies() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

// MockOAuthProvider implements auth.Service for testing OAuth functionality
type MockOAuthProvider struct {
	mock.Mock
}

func (m *MockOAuthProvider) Authenticate(ctx context.Context, strategy string, credentials interface{}) (*auth.AuthResult, error) {
	args := m.Called(ctx, strategy, credentials)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.AuthResult), args.Error(1)
}

func (m *MockOAuthProvider) ValidateToken(ctx context.Context, token string) (*auth.TokenClaims, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenClaims), args.Error(1)
}

func (m *MockOAuthProvider) RefreshToken(ctx context.Context, refreshToken string) (*auth.AuthResult, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.AuthResult), args.Error(1)
}

func (m *MockOAuthProvider) RevokeToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockOAuthProvider) GetSupportedStrategies() []string {
	args := m.Called()
	return args.Get(0).([]string)
}
