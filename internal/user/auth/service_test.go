package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gentra/decorator-arch-go/internal/auth"
	"github.com/gentra/decorator-arch-go/internal/user"
	userAuth "github.com/gentra/decorator-arch-go/internal/user/auth"
)

// Mock implementations for testing
type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) Register(ctx context.Context, data user.RegisterData) (*user.User, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserService) Login(ctx context.Context, email, password string) (*user.AuthResult, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.AuthResult), args.Error(1)
}

func (m *mockUserService) GetByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserService) UpdateProfile(ctx context.Context, id string, data user.UpdateProfileData) (*user.User, error) {
	args := m.Called(ctx, id, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserService) GetPreferences(ctx context.Context, userID string) (*user.UserPreferences, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.UserPreferences), args.Error(1)
}

func (m *mockUserService) UpdatePreferences(ctx context.Context, userID string, prefs user.UserPreferences) error {
	args := m.Called(ctx, userID, prefs)
	return args.Error(0)
}

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Authenticate(ctx context.Context, strategy string, credentials interface{}) (*auth.AuthResult, error) {
	args := m.Called(ctx, strategy, credentials)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.AuthResult), args.Error(1)
}

func (m *mockAuthService) CreateUser(ctx context.Context, data auth.CreateUserData) (*auth.User, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.User), args.Error(1)
}

func (m *mockAuthService) ValidateToken(ctx context.Context, token string) (*auth.TokenClaims, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenClaims), args.Error(1)
}

func (m *mockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*auth.AuthResult, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.AuthResult), args.Error(1)
}

func (m *mockAuthService) RevokeToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *mockAuthService) GetSupportedStrategies() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func TestNewService_GivenValidDependencies_WhenCreating_ThenReturnsService(t *testing.T) {
	mockNext := &mockUserService{}
	mockAuth := &mockAuthService{}

	service := userAuth.NewService(mockNext, mockAuth)

	assert.NotNil(t, service)
}

func TestRegister_GivenValidData_WhenRegistering_ThenDelegatesToNext(t *testing.T) {
	tests := []struct {
		name        string
		data        user.RegisterData
		nextResult  *user.User
		nextError   error
		expectError bool
	}{
		{
			name: "successful registration",
			data: user.RegisterData{
				Email:     "user@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Password:  "password123",
			},
			nextResult: &user.User{
				ID:        uuid.New(),
				Email:     "user@example.com",
				FirstName: "John",
				LastName:  "Doe",
			},
			nextError:   nil,
			expectError: false,
		},
		{
			name: "registration failure",
			data: user.RegisterData{
				Email:     "invalid@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Password:  "weak",
			},
			nextResult:  nil,
			nextError:   errors.New("registration failed"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNext := &mockUserService{}
			mockAuth := &mockAuthService{}

			// Setup expectations
			mockNext.On("Register", mock.Anything, tt.data).Return(tt.nextResult, tt.nextError)

			service := userAuth.NewService(mockNext, mockAuth)

			// Execute
			ctx := context.Background()
			result, err := service.Register(ctx, tt.data)

			// Verify
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.nextResult, result)
			}

			mockNext.AssertExpectations(t)
			mockAuth.AssertNotCalled(t, "Authenticate")
		})
	}
}

func TestLogin_GivenValidCredentials_WhenLoggingIn_ThenUsesAuthDomain(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		password      string
		authResult    *auth.AuthResult
		authError     error
		expectError   bool
		expectedError error
	}{
		{
			name:     "successful login",
			email:    "user@example.com",
			password: "password123",
			authResult: &auth.AuthResult{
				User: &auth.User{
					ID:           "user-123",
					Email:        "user@example.com",
					FirstName:    "John",
					LastName:     "Doe",
					PasswordHash: "hashed",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				},
				Token:        "jwt-token",
				RefreshToken: "refresh-token",
				ExpiresAt:    time.Now().Add(time.Hour),
			},
			authError:   nil,
			expectError: false,
		},
		{
			name:          "invalid credentials",
			email:         "user@example.com",
			password:      "wrongpassword",
			authResult:    nil,
			authError:     auth.ErrInvalidCredentials,
			expectError:   true,
			expectedError: user.ErrInvalidCredentials,
		},
		{
			name:          "user not found",
			email:         "nonexistent@example.com",
			password:      "password123",
			authResult:    nil,
			authError:     auth.ErrUserNotFound,
			expectError:   true,
			expectedError: user.ErrUserNotFound,
		},
		{
			name:          "generic auth error",
			email:         "user@example.com",
			password:      "password123",
			authResult:    nil,
			authError:     errors.New("generic auth error"),
			expectError:   true,
			expectedError: errors.New("generic auth error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNext := &mockUserService{}
			mockAuth := &mockAuthService{}

			// Setup expectations
			expectedCredentials := auth.BasicCredentials{
				Email:    tt.email,
				Password: tt.password,
			}
			mockAuth.On("Authenticate", mock.Anything, "basic", expectedCredentials).Return(tt.authResult, tt.authError)

			service := userAuth.NewService(mockNext, mockAuth)

			// Execute
			ctx := context.Background()
			result, err := service.Login(ctx, tt.email, tt.password)

			// Verify
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.expectedError != nil {
					if tt.expectedError == user.ErrInvalidCredentials || tt.expectedError == user.ErrUserNotFound {
						assert.Equal(t, tt.expectedError, err)
					} else {
						assert.Contains(t, err.Error(), tt.expectedError.Error())
					}
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.authResult.Token, result.Token)
				assert.Equal(t, tt.authResult.RefreshToken, result.RefreshToken)
				assert.Equal(t, tt.authResult.ExpiresAt, result.ExpiresAt)

				// Verify user conversion
				assert.NotNil(t, result.User)
				assert.Equal(t, tt.authResult.User.Email, result.User.Email)
				assert.Equal(t, tt.authResult.User.FirstName, result.User.FirstName)
				assert.Equal(t, tt.authResult.User.LastName, result.User.LastName)
				assert.Equal(t, tt.authResult.User.PasswordHash, result.User.PasswordHash)
			}

			mockAuth.AssertExpectations(t)
			mockNext.AssertNotCalled(t, "Login")
		})
	}
}

func TestLogin_GivenAuthUserWithInvalidUUID_WhenLoggingIn_ThenUsesNilUUID(t *testing.T) {
	mockNext := &mockUserService{}
	mockAuth := &mockAuthService{}

	authResult := &auth.AuthResult{
		User: &auth.User{
			ID:           "invalid-uuid",
			Email:        "user@example.com",
			FirstName:    "John",
			LastName:     "Doe",
			PasswordHash: "hashed",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		Token:        "jwt-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	expectedCredentials := auth.BasicCredentials{
		Email:    "user@example.com",
		Password: "password123",
	}
	mockAuth.On("Authenticate", mock.Anything, "basic", expectedCredentials).Return(authResult, nil)

	service := userAuth.NewService(mockNext, mockAuth)

	// Execute
	ctx := context.Background()
	result, err := service.Login(ctx, "user@example.com", "password123")

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.User)
	assert.Equal(t, uuid.Nil, result.User.ID) // Should fallback to nil UUID

	mockAuth.AssertExpectations(t)
}

func TestGetByID_GivenUserID_WhenGetting_ThenDelegatesToNext(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		nextResult  *user.User
		nextError   error
		expectError bool
	}{
		{
			name:   "successful get by ID",
			userID: "user123",
			nextResult: &user.User{
				ID:    uuid.New(),
				Email: "user@example.com",
			},
			nextError:   nil,
			expectError: false,
		},
		{
			name:        "user not found",
			userID:      "nonexistent",
			nextResult:  nil,
			nextError:   user.ErrUserNotFound,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNext := &mockUserService{}
			mockAuth := &mockAuthService{}

			// Setup expectations
			mockNext.On("GetByID", mock.Anything, tt.userID).Return(tt.nextResult, tt.nextError)

			service := userAuth.NewService(mockNext, mockAuth)

			// Execute
			ctx := context.Background()
			result, err := service.GetByID(ctx, tt.userID)

			// Verify
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.nextResult, result)
			}

			mockNext.AssertExpectations(t)
			mockAuth.AssertNotCalled(t, "Authenticate")
		})
	}
}

func TestUpdateProfile_GivenProfileData_WhenUpdating_ThenDelegatesToNext(t *testing.T) {
	firstName := "John"
	lastName := "Doe"

	tests := []struct {
		name        string
		userID      string
		data        user.UpdateProfileData
		nextResult  *user.User
		nextError   error
		expectError bool
	}{
		{
			name:   "successful profile update",
			userID: "user123",
			data: user.UpdateProfileData{
				FirstName: &firstName,
				LastName:  &lastName,
			},
			nextResult: &user.User{
				ID:        uuid.New(),
				FirstName: firstName,
				LastName:  lastName,
			},
			nextError:   nil,
			expectError: false,
		},
		{
			name:   "profile update failure",
			userID: "user123",
			data: user.UpdateProfileData{
				FirstName: &firstName,
			},
			nextResult:  nil,
			nextError:   errors.New("update failed"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNext := &mockUserService{}
			mockAuth := &mockAuthService{}

			// Setup expectations
			mockNext.On("UpdateProfile", mock.Anything, tt.userID, tt.data).Return(tt.nextResult, tt.nextError)

			service := userAuth.NewService(mockNext, mockAuth)

			// Execute
			ctx := context.Background()
			result, err := service.UpdateProfile(ctx, tt.userID, tt.data)

			// Verify
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.nextResult, result)
			}

			mockNext.AssertExpectations(t)
			mockAuth.AssertNotCalled(t, "Authenticate")
		})
	}
}

func TestGetPreferences_GivenUserID_WhenGetting_ThenDelegatesToNext(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		nextResult  *user.UserPreferences
		nextError   error
		expectError bool
	}{
		{
			name:   "successful get preferences",
			userID: "user123",
			nextResult: &user.UserPreferences{
				Theme:    "dark",
				Language: "en",
				Timezone: "UTC",
			},
			nextError:   nil,
			expectError: false,
		},
		{
			name:        "preferences not found",
			userID:      "nonexistent",
			nextResult:  nil,
			nextError:   errors.New("preferences not found"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNext := &mockUserService{}
			mockAuth := &mockAuthService{}

			// Setup expectations
			mockNext.On("GetPreferences", mock.Anything, tt.userID).Return(tt.nextResult, tt.nextError)

			service := userAuth.NewService(mockNext, mockAuth)

			// Execute
			ctx := context.Background()
			result, err := service.GetPreferences(ctx, tt.userID)

			// Verify
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.nextResult, result)
			}

			mockNext.AssertExpectations(t)
			mockAuth.AssertNotCalled(t, "Authenticate")
		})
	}
}

func TestUpdatePreferences_GivenPreferences_WhenUpdating_ThenDelegatesToNext(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		prefs       user.UserPreferences
		nextError   error
		expectError bool
	}{
		{
			name:   "successful preferences update",
			userID: "user123",
			prefs: user.UserPreferences{
				Theme:    "dark",
				Language: "en",
				Timezone: "UTC",
			},
			nextError:   nil,
			expectError: false,
		},
		{
			name:   "preferences update failure",
			userID: "user123",
			prefs: user.UserPreferences{
				Theme:    "invalid",
				Language: "invalid",
				Timezone: "invalid",
			},
			nextError:   errors.New("invalid preferences"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNext := &mockUserService{}
			mockAuth := &mockAuthService{}

			// Setup expectations
			mockNext.On("UpdatePreferences", mock.Anything, tt.userID, tt.prefs).Return(tt.nextError)

			service := userAuth.NewService(mockNext, mockAuth)

			// Execute
			ctx := context.Background()
			err := service.UpdatePreferences(ctx, tt.userID, tt.prefs)

			// Verify
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockNext.AssertExpectations(t)
			mockAuth.AssertNotCalled(t, "Authenticate")
		})
	}
}

func TestConvertAuthUserToUserDomain_GivenNilAuthUser_WhenConverting_ThenReturnsNil(t *testing.T) {
	mockNext := &mockUserService{}
	mockAuth := &mockAuthService{}

	authResult := &auth.AuthResult{
		User:         nil, // Nil user
		Token:        "jwt-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	expectedCredentials := auth.BasicCredentials{
		Email:    "user@example.com",
		Password: "password123",
	}
	mockAuth.On("Authenticate", mock.Anything, "basic", expectedCredentials).Return(authResult, nil)

	service := userAuth.NewService(mockNext, mockAuth)

	// Execute
	ctx := context.Background()
	result, err := service.Login(ctx, "user@example.com", "password123")

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.User) // Should be nil

	mockAuth.AssertExpectations(t)
}

func TestUserAuthService_GivenCompleteWorkflow_WhenExecuting_ThenWorksCorrectly(t *testing.T) {
	mockNext := &mockUserService{}
	mockAuth := &mockAuthService{}

	// Test data
	registerData := user.RegisterData{
		Email:     "user@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password123",
	}

	userResult := &user.User{
		ID:        uuid.New(),
		Email:     "user@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	authResult := &auth.AuthResult{
		User: &auth.User{
			ID:           userResult.ID.String(),
			Email:        "user@example.com",
			FirstName:    "John",
			LastName:     "Doe",
			PasswordHash: "hashed",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		Token:        "jwt-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	prefs := user.UserPreferences{
		Theme:    "dark",
		Language: "en",
		Timezone: "UTC",
	}

	// Setup expectations
	mockNext.On("Register", mock.Anything, registerData).Return(userResult, nil)
	mockNext.On("GetByID", mock.Anything, userResult.ID.String()).Return(userResult, nil)
	mockNext.On("UpdateProfile", mock.Anything, userResult.ID.String(), mock.Anything).Return(userResult, nil)
	mockNext.On("GetPreferences", mock.Anything, userResult.ID.String()).Return(&prefs, nil)
	mockNext.On("UpdatePreferences", mock.Anything, userResult.ID.String(), prefs).Return(nil)

	expectedCredentials := auth.BasicCredentials{
		Email:    "user@example.com",
		Password: "password123",
	}
	mockAuth.On("Authenticate", mock.Anything, "basic", expectedCredentials).Return(authResult, nil)

	service := userAuth.NewService(mockNext, mockAuth)
	ctx := context.Background()

	// Test registration (delegates to next)
	regResult, err := service.Register(ctx, registerData)
	assert.NoError(t, err)
	assert.Equal(t, userResult, regResult)

	// Test login (uses auth domain)
	loginResult, err := service.Login(ctx, "user@example.com", "password123")
	assert.NoError(t, err)
	assert.NotNil(t, loginResult)
	assert.Equal(t, authResult.Token, loginResult.Token)

	// Test get by ID (delegates to next)
	getResult, err := service.GetByID(ctx, userResult.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, userResult, getResult)

	// Test update profile (delegates to next)
	firstName := "Jane"
	updateData := user.UpdateProfileData{FirstName: &firstName}
	updateResult, err := service.UpdateProfile(ctx, userResult.ID.String(), updateData)
	assert.NoError(t, err)
	assert.Equal(t, userResult, updateResult)

	// Test get preferences (delegates to next)
	prefsResult, err := service.GetPreferences(ctx, userResult.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, &prefs, prefsResult)

	// Test update preferences (delegates to next)
	err = service.UpdatePreferences(ctx, userResult.ID.String(), prefs)
	assert.NoError(t, err)

	mockNext.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}