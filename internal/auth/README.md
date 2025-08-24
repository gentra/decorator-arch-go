# Auth Domain Implementation

This directory contains a complete implementation of the Authentication domain using the Decorator Architecture pattern. It demonstrates how authentication functionality can be separated into its own domain while being composed with other domains through clean interfaces.

## 🏗️ Architecture Overview

The auth domain is designed as a standalone service that can be:
- Used as a library within a monolith
- Extracted as a separate microservice  
- Composed with other domains through adapters
- Tested independently

## 📁 Directory Structure

```
internal/auth/
├── auth.go                 # Domain interface and types (ONLY auth.Service interface)
├── factory/                # Service factory and strategy implementations
│   ├── auth_service.go     # Main factory and auth service implementation
│   ├── jwt_manager.go      # JWT token management
│   └── strategies.go       # Strategy implementations (all implement auth.Service)
├── usecase/                # Business logic layer
│   └── service.go          # Usecase implementation (implements auth.Service)
└── README.md              # This file
```

## 🎯 Architecture Compliance

### Single Interface Rule
The auth domain strictly follows the architectural rule of having **only ONE interface**:

```go
// The ONLY interface in the auth domain
type Service interface {
    Authenticate(ctx context.Context, strategy string, credentials interface{}) (*AuthResult, error)
    ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
    RefreshToken(ctx context.Context, refreshToken string) (*AuthResult, error)
    RevokeToken(ctx context.Context, token string) error
    GetSupportedStrategies() []string
}
```

### Strategy Pattern Implementation
All authentication strategies (basic, OAuth, JWT) implement the same `auth.Service` interface:
- **Basic Auth Strategy**: Handles username/password authentication
- **OAuth Strategy**: Handles external provider authentication  
- **JWT Strategy**: Handles token-based authentication
- **OAuth Providers**: Each provider (Google, GitHub, etc.) implements `auth.Service`

The factory uses strategy pattern to route requests to the appropriate implementation based on the `strategy` parameter.

## 🔐 Authentication Strategies

### 1. Basic Authentication
Username/password authentication with JWT tokens:
```go
credentials := auth.BasicCredentials{
    Email:    "user@example.com",
    Password: "password123",
}
result, err := authService.Authenticate(ctx, "basic", credentials)
```

### 2. OAuth Authentication  
External provider authentication (Google, GitHub, etc.):
```go
credentials := auth.OAuthCredentials{
    Provider:    "google",
    AccessToken: "oauth-access-token",
}
result, err := authService.Authenticate(ctx, "oauth", credentials)
```

### 3. JWT Token Authentication
Direct token-based authentication:
```go
credentials := auth.JWTCredentials{
    Token: "existing-jwt-token",
}
result, err := authService.Authenticate(ctx, "jwt", credentials)
```

## 🔧 Service Configuration

### Basic Configuration
```go
userService := createUserService() // user.Service implementation
config := factory.NewDefaultConfig(jwtSecret, userService)
authFactory := factory.NewAuthServiceFactory(config)
authService, _ := authFactory.Build()
```

### Production Configuration with OAuth
```go
// Create OAuth provider services (each implements auth.Service)
googleOAuth := createGoogleOAuthService()
githubOAuth := createGitHubOAuthService()

config := factory.Config{
    JWTSecret:   jwtSecret,
    AccessTTL:   time.Hour,
    RefreshTTL:  24 * time.Hour,
    UserService: userService,
    OAuthProviders: map[string]auth.Service{
        "google": googleOAuth,
        "github": githubOAuth,
    },
    Features: factory.FeatureFlags{
        EnableBasicAuth: true,
        EnableOAuth:     true,
        EnableJWTAuth:   true,
    },
}
authFactory := factory.NewAuthServiceFactory(config)
authService, _ := authFactory.Build()
```

### Testing Configuration
```go
config := factory.NewTestingConfig(userService)
authFactory := factory.NewAuthServiceFactory(config)
authService, _ := authFactory.Build()
```

## 🎫 Token Management

### JWT Token Service
- **Access Tokens**: Short-lived (default: 1 hour)
- **Refresh Tokens**: Long-lived (default: 24 hours)  
- **Token Revocation**: In-memory revocation list
- **Token Validation**: Signature verification + revocation check

### Token Lifecycle
```go
// Authenticate and get tokens
authResult, _ := authService.Authenticate(ctx, "basic", credentials)

// Validate access token
claims, _ := authService.ValidateToken(ctx, authResult.Token)

// Refresh expired token  
newResult, _ := authService.RefreshToken(ctx, authResult.RefreshToken)

// Revoke token
_ = authService.RevokeToken(ctx, authResult.Token)
```

## 🔗 Domain Integration

### Clean Domain Separation
The auth domain integrates with other domains through clean interfaces:

```go
// Create dependencies (other domain services)
userService := createUserService()           // user.Service
notificationService := createNotificationService() // notification.Service

// Create auth service
config := factory.NewDefaultConfig(jwtSecret, userService)
authService, _ := factory.NewAuthServiceFactory(config).Build()

// Use auth service in user domain through adapter
authEnabledUserService := userAuth.NewService(userService, authService)
```

### Architecture Benefits
- **Single Interface Rule**: Each domain has exactly one Service interface
- **Strategy Pattern**: All auth strategies implement auth.Service  
- **Clean Dependencies**: Auth depends on user.Service, not implementation details
- **Microservice Ready**: Domains can be extracted to separate services with zero code changes

## 🧪 Testing

### Unit Testing
Each strategy can be tested independently since they all implement auth.Service:
```go
// Test basic auth strategy
basicAuth := &basicAuthStrategy{
    userService:  mockUserService,
    tokenManager: mockTokenManager,
}
result, err := basicAuth.Authenticate(ctx, "basic", credentials)

// Test OAuth strategy
oauthAuth := &oauthAuthStrategy{
    userService:    mockUserService,
    tokenManager:   mockTokenManager,
    oauthProviders: map[string]auth.Service{"google": mockGoogleService},
}
result, err := oauthAuth.Authenticate(ctx, "oauth", credentials)
```

### Integration Testing
```go
// Test complete auth service
userService := createTestUserService()
config := factory.NewTestingConfig(userService)
authService, _ := factory.NewAuthServiceFactory(config).Build()

// Test end-to-end authentication flow
authResult, err := authService.Authenticate(ctx, "basic", credentials)
```

## 📊 Benefits Demonstrated

### Architectural Compliance
- **Single Interface Rule**: Strictly follows "one interface per domain" rule
- **No Interface Pollution**: Zero leaked interfaces - only auth.Service exists
- **Strategy Implementation**: All strategies implement the same interface
- **Factory Composition**: All logic assembly handled in factory layer

### Domain Separation
- **Single Responsibility**: Auth domain focuses only on authentication
- **Clear Boundaries**: Clean interface between auth and user domains
- **Independent Evolution**: Auth and user domains evolve separately
- **Microservice Extraction**: Zero-code-change extraction to separate service

### Strategy Pattern Benefits
- **Uniform Interface**: All strategies implement auth.Service
- **Runtime Configuration**: Enable/disable strategies via feature flags
- **Easy Extension**: Add new strategies by implementing auth.Service
- **Provider Flexibility**: OAuth providers are just auth.Service implementations

### Testing Benefits
- **Mock-Friendly**: Single interface makes mocking straightforward
- **Strategy Isolation**: Each strategy can be tested independently
- **Factory Testing**: Complete service assembly can be tested
- **Integration Testing**: End-to-end flows with real dependencies

## 🚀 Usage Examples

### Basic Usage
```go
// Create auth service
userService := createUserService() // user.Service implementation
config := factory.NewDefaultConfig(jwtSecret, userService)
authService, _ := factory.NewAuthServiceFactory(config).Build()

// Authenticate user
result, err := authService.Authenticate(ctx, "basic", auth.BasicCredentials{
    Email:    "user@example.com", 
    Password: "password",
})

// Use token for subsequent requests
claims, err := authService.ValidateToken(ctx, result.Token)
```

### With User Domain Integration
```go
// Create services
userService := createUserService()
authService := createAuthService(userService)

// Create user service with auth capabilities
authEnabledUserService := userAuth.NewService(userService, authService)

// Login through user domain (delegates to auth domain)
authResult, err := authEnabledUserService.Login(ctx, email, password)

// The user service now has auth integration
// but still implements only user.Service interface
```

### Microservice Architecture
```bash
# Each domain can be a separate microservice with zero code changes

auth-service:
  - Implements auth.Service interface
  - JWT token generation/validation
  - Multiple authentication strategies  
  - Token revocation management
  - OAuth provider integrations

user-service:
  - Implements user.Service interface
  - User CRUD operations
  - Profile management
  - Preferences management
  - Uses auth-service via HTTP/gRPC (same auth.Service interface)
```

## 🔮 Future Enhancements

All future enhancements follow the single interface rule - each new capability implements `auth.Service`:

- **Redis Token Storage**: Replace in-memory revocation with Redis (implement auth.Service)
- **LDAP Strategy**: Add LDAP authentication strategy (implement auth.Service)  
- **Multi-factor Auth**: 2FA/MFA as auth strategies (implement auth.Service)
- **Session Management**: Session-based auth strategy (implement auth.Service)
- **Biometric Auth**: Fingerprint/face recognition (implement auth.Service)
- **Social Auth**: Additional OAuth providers (implement auth.Service)
- **API Key Auth**: API key authentication strategy (implement auth.Service)
- **Certificate Auth**: X.509 certificate authentication (implement auth.Service)

### Integration Enhancements
- **Audit Integration**: Use `audit.Service` for authentication logging
- **Rate Limit Integration**: Use `ratelimit.Service` for auth protection
- **Notification Integration**: Use `notification.Service` for auth alerts

## 📚 Related Documentation

- [Main README](../../README.md) - Project overview and architectural rules
- [User Domain](../user/README.md) - User domain implementation  
- [Technical Plan](../../docs/technical-plan.md) - Architecture plan

## 🎯 Summary

This auth domain demonstrates **perfect compliance** with our architectural rules:

✅ **Single Interface**: Only `auth.Service` exists  
✅ **Strategy Pattern**: All auth methods implement `auth.Service`  
✅ **Factory Composition**: All assembly logic in factory layer  
✅ **Clean Dependencies**: Depends on `user.Service`, not implementations  
✅ **Microservice Ready**: Zero-code-change extraction possible  

The auth domain shows how authentication can be completely separated while maintaining clean integration with other domains through the strict single-interface architecture.
