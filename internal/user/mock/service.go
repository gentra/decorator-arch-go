package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/gentra/decorator-arch-go/internal/user"
	"github.com/gentra/decorator-arch-go/internal/validationrule"
)

// MockUserService is a mock implementation of user.Service
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(ctx context.Context, data user.RegisterData) (*user.User, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, email, password string) (*user.AuthResult, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.AuthResult), args.Error(1)
}

func (m *MockUserService) GetByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) UpdateProfile(ctx context.Context, id string, data user.UpdateProfileData) (*user.User, error) {
	args := m.Called(ctx, id, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) GetPreferences(ctx context.Context, userID string) (*user.UserPreferences, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.UserPreferences), args.Error(1)
}

func (m *MockUserService) UpdatePreferences(ctx context.Context, userID string, prefs user.UserPreferences) error {
	args := m.Called(ctx, userID, prefs)
	return args.Error(0)
}

// MockValidationService is a mock implementation of validation.Service
type MockValidationService struct {
	mock.Mock
}

func (m *MockValidationService) ValidateStruct(ctx context.Context, data interface{}) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockValidationService) ValidateField(ctx context.Context, field string, value interface{}, rules string) error {
	args := m.Called(ctx, field, value, rules)
	return args.Error(0)
}

func (m *MockValidationService) ValidateUserRegistration(ctx context.Context, data interface{}) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockValidationService) ValidateUserUpdate(ctx context.Context, data interface{}) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockValidationService) ValidateUserPreferences(ctx context.Context, data interface{}) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockValidationService) ValidateUserID(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockValidationService) ValidateEmail(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockValidationService) ValidatePassword(ctx context.Context, password string) error {
	args := m.Called(ctx, password)
	return args.Error(0)
}

func (m *MockValidationService) AddCustomRule(name string, rule validationrule.Service) error {
	args := m.Called(name, rule)
	return args.Error(0)
}

func (m *MockValidationService) RemoveCustomRule(name string) error {
	args := m.Called(name)
	return args.Error(0)
}
