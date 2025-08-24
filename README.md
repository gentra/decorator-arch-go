# Decorator Architecture in Go

## Overview

This repository is an example of applying the Decorator pattern as a code architecture in go.

## Goals

The goals for this code architecture are:

- Flexibility to start with 1 monolithic service but with the ability to easily pluck out the individual components into separate microservices when the scale calls for it.
- Flexibility to add feature-flags to the codebase without making a mess in the codebase.
- Ability to easily test the codebase.

## Architectural Rules

This project follows strict architectural rules to ensure clean separation of concerns and maintainability:

### Core Domain Rules

1. **Single Interface Per Domain**: Inside each domain folder (e.g., `internal/user`, `internal/auth`) there should only be **ONE** interface called `Service`. This rule is absolute with no exceptions.

2. **Implementation Folder Constraint**: Each implementation folder (e.g., `internal/user/usecase`, `internal/user/encryption`) should **ONLY** contain files that implement the domain's `Service` interface. No other interfaces are allowed in implementation folders.

3. **Interface Extraction Rule**: If you need a different interface, it **MUST** become its own separate domain folder `internal/<new-domain>`. Do **NOT** put interfaces in factory folders or implementation folders.

4. **Factory-Managed Composition**: The strategy and decorator patterns should all be handled inside the domain's factory folder `internal/<domain>/factory`. Factories act as composition roots that orchestrate domain services.

### Interface Design Rules

5. **Generic Service Design**: Domain services should be designed to be reusable across many modules/components. Avoid domain-specific method names when possible (e.g., use `EncryptWithPurpose()` instead of `EncryptEmail()`).

6. **No Cross-Domain Interfaces**: Implementation folders cannot define interfaces that belong to other domains. Extract them into separate domains instead.

7. **Clean Dependencies**: Domains should depend on other domain interfaces, not on implementation details. Use dependency injection through factory configurations.

### Enforcement Rules

8. **Zero Tolerance**: These rules are strictly enforced. Any violation should be immediately refactored.

9. **Interface Ownership**: Each domain owns exactly one interface. If you think you need multiple interfaces in one domain, create separate domains instead.

10. **Adapter Pattern**: When integrating domains (like `internal/user/auth` using `internal/auth`), use the adapter pattern where the adapter implements the target domain's interface.

## Directory Structure

Here's the directory structure of the codebase showing our multi-domain architecture:

```
.
├── cmd/                    # Application entry points grouped by delivery mechanisms
│   └── rest/              # REST API entry point
├── internal/              # Domain-driven architecture with strict separation
│   ├── user/              # User domain (main business domain)
│   │   ├── user.go        # ONLY the user.Service interface and types
│   │   ├── factory/       # Composition root for user service decorators
│   │   ├── gorm/          # Database persistence layer
│   │   ├── redis/         # Caching decorator layer
│   │   ├── audit/         # Audit logging decorator (uses audit domain)
│   │   ├── encryption/    # Data encryption decorator (uses encryption domain)
│   │   ├── ratelimit/     # Rate limiting decorator (uses ratelimit domain)
│   │   ├── validation/    # Input validation decorator (uses validation domain)
│   │   ├── usecase/       # Business logic layer (uses notification, token, events domains)
│   │   └── auth/          # Auth integration adapter (uses auth domain)
│   ├── auth/              # Authentication domain
│   │   ├── auth.go        # ONLY the auth.Service interface and types
│   │   ├── factory/       # JWT management and auth strategies
│   │   └── usecase/       # Auth business logic implementation
│   ├── audit/             # Audit logging domain
│   │   ├── audit.go       # ONLY the audit.Service interface and types
│   │   └── console/       # Console logging implementation
│   ├── encryption/        # Generic encryption domain
│   │   ├── encryption.go  # ONLY the encryption.Service interface and types
│   │   ├── aes/           # AES encryption implementation
│   │   └── noop/          # No-op encryption implementation
│   ├── ratelimit/         # Rate limiting domain
│   │   ├── ratelimit.go   # ONLY the ratelimit.Service interface and types
│   │   └── memory/        # In-memory rate limiter implementation
│   ├── validation/        # Input validation domain
│   │   ├── validation.go  # ONLY the validation.Service interface and types
│   │   └── standard/      # Standard validation rules implementation
│   ├── validationrule/    # Validation rule domain
│   │   └── validationrule.go # ONLY the validationrule.Service interface and types
│   ├── notification/      # Notification domain
│   │   ├── notification.go # ONLY the notification.Service interface and types
│   │   └── mock/          # Mock notification implementation
│   ├── token/             # Token management domain
│   │   ├── token.go       # ONLY the token.Service interface and types
│   │   └── jwt/           # JWT token implementation
│   ├── events/            # Event publishing domain
│   │   ├── events.go      # ONLY the events.Service interface and types
│   │   └── memory/        # In-memory event publisher implementation
│   └── eventhandler/      # Event handler domain
│       └── eventhandler.go # ONLY the eventhandler.Service interface and types
├── examples/              # Demo applications showing the architecture
├── docs/                  # Technical documentation
├── migrations/            # Database migrations
├── .env.example          # Environment configuration template
├── go.mod                # Go module dependencies
└── README.md             # This file
```

### Key Architecture Principles Demonstrated:

- **One Interface Per Domain**: Each `internal/<domain>/` folder contains exactly one `Service` interface
- **Clean Separation**: No interfaces leak between domains
- **Decorator Pattern**: User domain shows how decorators wrap core functionality
- **Adapter Pattern**: `user/auth` shows how domains integrate via adapters
- **Generic Design**: `encryption`, `validation`, etc. are reusable across domains
- **Factory Composition**: Each domain's factory handles strategy and composition logic

## Decorator Pattern

This project extensively uses the Decorator pattern for service implementations. This pattern allows us to add behavior to services by wrapping them in decorator layers, each with a specific responsibility.

### Example of the implementation

1. **Base Interface**
For each domain, there should only be one interface that's used across the layers (either using decorator pattern or strategy pattern).

```go
// Example: user.Service interface
type Service interface {
    Create(ctx context.Context, data CreateUserData) (*User, error)
    GetByID(ctx context.Context, id string) (*User, error)
    // ... other methods
}
```

2. **Service Chain**
Services are chained in layers, where each layer implements the same interface and adds specific functionality:
```go
// Example for factory.go
userPg := userPostgres.NewService(db)                    // Storage layer
userCache := userCache.NewService(userPg, 5*time.Minute) // Cache layer
userEnc := userEncryption.NewService(userCache, key)     // Encryption layer
userService := userCore.NewService(userEnc)              // Business logic layer
```

3. **Layer Responsibilities**

Each decorator has a specific responsibility. There are no rules in how many layers or specific names, it's all up to the developer. For example:

- **Postgres Layer** (`postgres`): Handles postgres database operations
- **Redis Layer** (`redis`): Adds redis cache before hitting the database
- **Encryption Layer** (`encryption`): Handles encryption/decryption of sensitive data
- **UseCase Layer** (`usecase`): Implements business logic and validation

4. **Strategy and Decorator Patterns in Factory**

Strategy and decorator patterns are managed entirely within the factory package. For example:

```go
// Example factory implementing strategy pattern
package factory

type Config struct {
    Features FeatureFlags
}

type FeatureFlags struct {
    EnableBasicAuth bool
    EnableOAuth     bool
    EnableJWTAuth   bool
}

func (f *Factory) Build() (domain.Service, error) {
    // All strategy logic handled here
    if f.config.Features.EnableBasicAuth {
        // Add basic auth strategy
    }
    
    if f.config.Features.EnableOAuth {
        // Add OAuth strategy
    }
    
    // Return assembled service with all strategies/decorators
}
```

5. **Code Versioning and Feature Flags**

Code versioning and feature flags can be implemented by simply adding a new file that implements the interface. For example:

```
├── internal/              # Main folder where each domain is grouped together and can be plucked out into separate microservices.
│   ├── domain1/          # Domain-specific package where there's a domain interface and types, factory methods, then any kind of implementation-folders necessary for the domain.
│   │   ├── domain1.go    # Domain interface and types
│   │   ├── factory/      # Factory methods for creating domain objects, handling all strategy logic
│   │   ├── usecase/      # Usecase layer where the business logic is implemented.
│   │   ├──── usecase.go      # Implementation of the interface for v1 of the code
│   │   └──── usecase_v2.go      # Implementation of the interface for version 2 of the business logic
...
```

Then, the feature flag can be implemented by using the strategy pattern in the factory method. For example:

```go
// Example for factory.go
if (featureFlag.IsEnabled("v2")) {
    userService := usecase.NewService()
} else {
    userService := userCore.NewServiceV2()
}
```


### Adding a New Decorator

To add a new decorator:

1. Create a new package with a service struct:
```go
type service struct {
    next user.Service  // Reference to next service
    // Additional fields needed
}
```

2. Implement the interface methods:
```go
func (s *service) MethodName(ctx context.Context, ...) (..., error) {
    // Add your logic before calling next
    result, err := s.next.MethodName(ctx, ...)
    // Add your logic after calling next
    return result, err
}
```

3. Add a constructor:
```go
func NewService(next user.Service, ...) user.Service {
    return &service{
        next: next,
        // Initialize other fields
    }
}
```

## Domain Examples

### User Domain (Main Business Domain)
Demonstrates the full Decorator Architecture with cross-domain dependencies:
- **Storage Layer** (`gorm`): Database operations using GORM
- **Caching Layer** (`redis`): Performance optimization with Redis
- **Audit Layer** (`audit`): Uses `audit.Service` for operation logging
- **Rate Limiting Layer** (`ratelimit`): Uses `ratelimit.Service` for API protection
- **Encryption Layer** (`encryption`): Uses `encryption.Service` for data security
- **Validation Layer** (`validation`): Uses `validation.Service` for input validation
- **UseCase Layer** (`usecase`): Business logic with `notification.Service`, `token.Service`, `events.Service`
- **Auth Adapter** (`auth`): Adapter that uses `auth.Service` for authentication

### Supporting Domains (Single-Purpose Services)

**Auth Domain**: Authentication and authorization
- **Single Responsibility**: User authentication, token management, strategy handling
- **Clean Interface**: Only `auth.Service` with auth-specific methods
- **Strategy Pattern**: Multiple auth strategies (basic, OAuth, JWT) in factory

**Encryption Domain**: Generic encryption service
- **Reusable Design**: Purpose-based encryption (`EncryptWithPurpose`)
- **Multiple Implementations**: AES encryption, no-op for development
- **Cross-Domain Usage**: Used by user domain and potentially others

**Audit Domain**: System-wide audit logging
- **Generic Interface**: Can log any domain's operations
- **Flexible Implementation**: Console logger with future database/external service support
- **Compliance Ready**: Structured audit entries with metadata

**Validation Domain**: Input validation service
- **Reusable Validators**: Email, password, UUID, user-specific validations
- **Domain Agnostic**: Can validate any domain's input data
- **Error Handling**: Structured validation errors with field details

**Rate Limiting Domain**: API protection service
- **Configurable Limits**: Per-operation, per-user rate limiting
- **Algorithm Choice**: Sliding window implementation
- **Cross-Domain Protection**: Any domain can use rate limiting

**Notification Domain**: Communication service
- **Multi-Channel**: Email, push, SMS notification support
- **Async Operations**: Non-blocking notification sending
- **Template Support**: Welcome emails, profile updates, etc.

**Token Domain**: Token management service
- **JWT Implementation**: Auth tokens, refresh tokens
- **Configurable TTL**: Different expiration times per token type
- **Secure Generation**: Cryptographically secure token creation

**Events Domain**: Event publishing service
- **Domain Events**: User registered, logged in, profile updated
- **Async Processing**: In-memory publisher with future message queue support
- **Event Sourcing Ready**: Structured events with aggregate information

**Event Handler Domain**: Event processing service
- **Decoupled Handlers**: Clean separation of event handling logic
- **Multi-type Support**: Handlers can process multiple event types
- **Configurable Processing**: Retry, timeout, and concurrency control

**Validation Rule Domain**: Custom validation logic service
- **Extensible Rules**: Define custom validation logic independently
- **Reusable Components**: Rules can be shared across multiple validation contexts
- **Metadata Support**: Rich configuration and parameter handling

## Benefits

### Design Benefits
- **Strict Separation of Concerns**: Each domain handles exactly one responsibility
- **Single Interface Guarantee**: Every domain has exactly one public interface, eliminating confusion
- **Predictable Structure**: Consistent patterns across all domains make the codebase easy to navigate
- **Zero Interface Pollution**: No leaked interfaces between domains ensures clean boundaries

### Development Benefits  
- **Easy Domain Extraction**: Any domain can be moved to a microservice without code changes
- **Reusable Components**: Generic domains (encryption, validation, etc.) work across multiple business domains
- **Clear Dependencies**: Factory injection makes all dependencies explicit and testable
- **Rapid Feature Development**: New features can be added as decorators without modifying existing code

### Testing Benefits
- **Complete Isolation**: Each layer can be tested independently by mocking its dependencies
- **Mock-Friendly Design**: Single interfaces per domain make mocking straightforward
- **Behavioral Testing**: Decorator pattern enables testing of behavior combinations
- **Domain-Specific Testing**: Each domain can have focused, relevant tests

### Maintenance Benefits
- **Refactoring Safety**: Strict interfaces prevent breaking changes from propagating
- **Performance Optimization**: Individual layers can be optimized without affecting others
- **Feature Flags**: Decorators can be enabled/disabled via configuration
- **Debugging Clarity**: Issues can be traced to specific layers in the decorator chain

### Scalability Benefits
- **Microservice Migration**: Domains can be extracted to separate services with minimal effort
- **Independent Scaling**: Different domains can be scaled based on their specific needs
- **Team Organization**: Different teams can own different domains with clear boundaries
- **Technology Flexibility**: Each domain can use the most appropriate technology stack

## Running the code

```bash
go run examples/user_service_demo.go
```