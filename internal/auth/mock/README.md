# Auth Domain Mocks

This directory contains centralized mock implementations for the auth domain, consolidating all auth-related mocks that were previously scattered throughout the usecase directory.

## 📁 Mock Implementations

### MockUserService
Alias to the centralized user mock service from `internal/user/mock`. Provides mock implementations for user domain operations.

### MockJWTTokenManager
Mock implementation of JWT token management functionality. Implements:
- `GenerateAuthToken(userID, email)` - Generates JWT access tokens
- `GenerateRefreshToken(userID)` - Generates JWT refresh tokens  
- `ValidateToken(token)` - Validates JWT tokens
- `RevokeToken(token)` - Revokes JWT tokens

### MockAuthStrategy
Mock implementation of `auth.Service` interface for testing authentication strategies. Implements:
- `Authenticate(ctx, strategy, credentials)` - Mock authentication
- `ValidateToken(ctx, token)` - Mock token validation
- `RefreshToken(ctx, refreshToken)` - Mock token refresh
- `RevokeToken(ctx, token)` - Mock token revocation
- `GetSupportedStrategies()` - Mock strategy list

### MockOAuthProvider
Mock implementation of OAuth provider functionality. Implements the same interface as `MockAuthStrategy` but specifically for OAuth testing scenarios.

## 🔧 Usage

### Import the mocks
```go
import (
    authmock "github.com/gentra/decorator-arch-go/internal/auth/mock"
)
```

### Create mock instances
```go
// Create mock user service
mockUserService := new(authmock.MockUserService)

// Create mock JWT token manager
mockTokenManager := new(authmock.MockJWTTokenManager)

// Create mock auth strategy
mockStrategy := new(authmock.MockAuthStrategy)

// Create mock OAuth provider
mockOAuthProvider := new(authmock.MockOAuthProvider)
```

### Set up mock expectations
```go
// Set up user service mock
mockUserService.On("Login", mock.Anything, "test@example.com", "password123").Return(loginResult, nil)

// Set up auth strategy mock
mockStrategy.On("Authenticate", mock.Anything, "basic", mock.Anything).Return(authResult, nil)

// Verify expectations
mockUserService.AssertExpectations(t)
mockStrategy.AssertExpectations(t)
```

## 🎯 Benefits of Centralization

1. **Single Source of Truth**: All auth mocks are defined in one location
2. **Consistency**: Ensures all tests use the same mock implementations
3. **Maintainability**: Changes to mock interfaces only need to be made in one place
4. **Reusability**: Mocks can be easily shared across different test files
5. **Clean Architecture**: Follows the project's architectural principles

## 📚 Related Documentation

- [Auth Domain README](../README.md) - Main auth domain documentation
- [User Domain Mocks](../../user/mock/README.md) - User domain mocks
- [Main Project README](../../../README.md) - Project overview and architecture
