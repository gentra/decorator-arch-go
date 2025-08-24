# User Domain Implementation

This directory contains a complete implementation of the User domain using the Decorator Architecture pattern in Go. It demonstrates how to build flexible, testable, and maintainable services that can easily evolve from monoliths to microservices.

## 🏗️ Architecture Overview

The user domain is implemented using the Decorator pattern, where each layer adds specific functionality while maintaining the same interface. This allows for:

- **Flexibility**: Add or remove layers without affecting others
- **Testability**: Each layer can be tested in isolation
- **Separation of Concerns**: Each layer has a single responsibility
- **Easy Microservice Extraction**: Layers can be moved to separate services

## 📁 Directory Structure

```
internal/user/
├── user.go                 # Domain interface and types
├── factory/                # Service factory for assembling decorator chain
│   └── factory.go
├── gorm/                   # Database persistence layer (GORM)
│   ├── models.go
│   └── service.go
├── redis/                  # Caching layer (Redis)
│   ├── service.go
│   └── service_test.go
├── audit/                  # Audit logging layer
│   └── service.go
├── ratelimit/              # Rate limiting layer
│   └── service.go
├── encryption/             # Data encryption layer
│   └── service.go
├── validation/             # Input validation layer
│   ├── service.go
│   └── service_test.go
├── usecase/                # Business logic layer
│   └── service.go
├── auth/                   # Authentication strategies
│   └── strategies.go
└── README.md              # This file
```

## 🎯 Service Layers

The service is built using the following layers (from top to bottom):

### 1. UseCase Layer
- **Purpose**: Business logic and orchestration
- **Responsibilities**: 
  - Domain business rules
  - Event publishing
  - External service coordination
  - Token generation
- **Always enabled**: Yes

### 2. Validation Layer
- **Purpose**: Input validation and business rules
- **Responsibilities**:
  - Email format validation
  - Password strength validation
  - Name validation
  - User preferences validation
- **Configuration**: Can be disabled for testing

### 3. Encryption Layer
- **Purpose**: Data encryption for sensitive fields
- **Responsibilities**:
  - Email encryption/decryption
  - Name encryption/decryption
  - AES-256 encryption
- **Configuration**: Optional (NoOp encryptor for development)

### 4. Rate Limiting Layer
- **Purpose**: API protection against abuse
- **Responsibilities**:
  - Per-user rate limiting
  - Per-operation rate limits
  - Configurable limits and windows
- **Implementation**: In-memory sliding window (Redis-based for production)

### 5. Audit Layer
- **Purpose**: Activity logging and audit trail
- **Responsibilities**:
  - Log all operations with metadata
  - Track user actions
  - Security audit trail
  - Performance metrics
- **Implementation**: Console logger (database/external service for production)

### 6. Cache Layer
- **Purpose**: Performance optimization through caching
- **Responsibilities**:
  - Cache-aside pattern
  - TTL management
  - Cache invalidation
  - Graceful fallback on cache failures
- **Implementation**: Redis

### 7. Storage Layer (GORM)
- **Purpose**: Database persistence
- **Responsibilities**:
  - CRUD operations
  - Transaction management
  - Data consistency
- **Always enabled**: Yes

## 🚀 Usage Examples

### Basic Service (Minimal Configuration)
```go
config := factory.NewTestingConfig(db, dependencies)
serviceFactory := factory.NewUserServiceFactory(config)
userService, _ := serviceFactory.BuildMinimal()
```

### Full Service (All Layers)
```go
config := factory.NewDefaultConfig(db, redisClient, dependencies)
serviceFactory := factory.NewUserServiceFactory(config)
userService, _ := serviceFactory.Build()
```

### Production Service (with Encryption)
```go
emailKey, _ := encryption.GenerateAESKey()
nameKey, _ := encryption.GenerateAESKey()
config := factory.NewProductionConfig(db, redisClient, emailKey, nameKey, dependencies)
serviceFactory := factory.NewUserServiceFactory(config)
userService, _ := serviceFactory.Build()
```

### Custom Configuration
```go
config := factory.Config{
    DB:          db,
    RedisClient: redisClient,
    Features: factory.FeatureFlags{
        EnableCache:      true,
        EnableAudit:      true,
        EnableRateLimit:  false, // Disable for testing
        EnableEncryption: false,
        EnableValidation: true,
    },
    Dependencies: dependencies,
}
serviceFactory := factory.NewUserServiceFactory(config)
userService, _ := serviceFactory.Build()
```

## 🔐 Authentication Strategies

The domain includes multiple authentication strategies using the Strategy pattern:

### Basic Authentication
- Username/password authentication
- JWT token generation
- Refresh token support

### OAuth Authentication
- External provider integration (Google, GitHub, etc.)
- Automatic user creation
- JWT token generation

### Usage
```go
jwtSecret := []byte("your-secret-key")
tokenTTL := time.Hour
authFactory := auth.NewAuthStrategyFactory(userService, jwtSecret, tokenTTL)

// Basic auth
basicAuth, _ := authFactory.CreateStrategy("basic")
authResult, _ := basicAuth.Authenticate(ctx, auth.BasicCredentials{
    Email:    "user@example.com",
    Password: "password",
})

// OAuth
oauthAuth, _ := authFactory.CreateStrategy("oauth")
authResult, _ := oauthAuth.Authenticate(ctx, auth.OAuthCredentials{
    Provider:    "google",
    AccessToken: "oauth-token",
})
```

## 🧪 Testing

The implementation includes comprehensive unit tests using table-driven tests with Gherkin syntax:

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific layer tests
go test ./redis/
go test ./validation/
```

### Test Examples
```go
func TestUserCacheService_GetByID(t *testing.T) {
    tests := []struct {
        name           string
        // ... test configuration
    }{
        {
            name: "Given user not in cache, When GetByID is called, Then should fetch from next service and cache result",
            // ... test implementation
        },
        {
            name: "Given user exists in cache, When GetByID is called, Then should return cached result without calling next service",
            // ... test implementation
        },
    }
    // ... test execution
}
```

## 🎨 Demo Application

Run the included demo application to see the architecture in action:

```bash
# Build the demo
go build -o user_demo ./examples/user_service_demo.go

# Run the demo
./user_demo
```

The demo showcases:
- Minimal service configuration
- Full service with all layers
- Production configuration with encryption
- Authentication strategies
- Performance comparison (caching effects)
- Audit logging with context

## 🔧 Configuration

### Environment Variables
```bash
# Database
DATABASE_URL=postgres://user:pass@localhost/userdb
REDIS_URL=redis://localhost:6379

# Feature Flags
ENABLE_CACHE=true
ENABLE_AUDIT=true
ENABLE_RATE_LIMIT=true
ENABLE_ENCRYPTION=false
ENABLE_VALIDATION=true

# Layer Configuration
CACHE_TTL=300s
RATE_LIMIT_PER_MINUTE=60
AUDIT_LOG_LEVEL=info
ENCRYPTION_KEY=base64-encoded-key
```

### Feature Flags
```go
type FeatureFlags struct {
    EnableCache      bool // Redis caching
    EnableAudit      bool // Audit logging
    EnableRateLimit  bool // Rate limiting
    EnableEncryption bool // Data encryption
    EnableValidation bool // Input validation
}
```

## 📊 Benefits Demonstrated

### Flexibility
- Layers can be added/removed via configuration
- Feature flags enable/disable functionality
- No code changes required for different environments

### Testability
- Each layer tested in isolation
- Mock interfaces for dependencies
- Table-driven tests with clear scenarios

### Performance
- Caching layer improves response times
- Rate limiting protects against abuse
- Benchmarks measure overhead of each layer

### Security
- Encryption for sensitive data
- Audit logging for compliance
- Rate limiting for DoS protection
- Input validation for data integrity

### Maintainability
- Single responsibility per layer
- Clear interfaces between layers
- Consistent patterns across the codebase

## 🚀 Future Enhancements

- Event sourcing layer for complete audit trail
- CQRS implementation with separate read/write models
- Circuit breaker pattern for external services
- Distributed tracing integration
- Health check endpoints for each layer
- Configuration hot-reloading
- A/B testing framework integration

## 📚 Related Documentation

- [Main README](../../README.md) - Project overview
- [Technical Plan](../../docs/technical-plan.md) - Detailed implementation plan
- [Factory Pattern](./factory/README.md) - Service factory documentation
- [Auth Strategies](./auth/README.md) - Authentication strategy patterns

This implementation serves as a comprehensive example of the Decorator Architecture pattern in Go, demonstrating how to build scalable, maintainable, and testable services that can evolve with your needs.
