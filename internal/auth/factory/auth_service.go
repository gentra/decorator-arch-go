package factory

import (
	"fmt"
	"time"

	"github.com/gentra/decorator-arch-go/internal/auth"
	"github.com/gentra/decorator-arch-go/internal/auth/usecase"
	"github.com/gentra/decorator-arch-go/internal/user"
)

// Config contains all configuration for building the auth service
type Config struct {
	// JWT configuration
	JWTSecret  []byte
	AccessTTL  time.Duration
	RefreshTTL time.Duration

	// User integration (from user domain)
	UserService user.Service

	// OAuth providers (now auth.Service implementations)
	OAuthProviders map[string]auth.Service

	// Feature flags
	Features FeatureFlags
}

// FeatureFlags controls which authentication strategies are enabled
type FeatureFlags struct {
	EnableBasicAuth bool
	EnableOAuth     bool
	EnableJWTAuth   bool
}

// DefaultFeatureFlags returns default feature flag configuration
func DefaultFeatureFlags() FeatureFlags {
	return FeatureFlags{
		EnableBasicAuth: true,
		EnableOAuth:     false, // Disabled by default as it requires provider setup
		EnableJWTAuth:   true,
	}
}

// AuthServiceFactory creates and assembles the complete auth service
type AuthServiceFactory struct {
	config Config
}

// NewAuthServiceFactory creates a new factory with the given configuration
func NewAuthServiceFactory(config Config) *AuthServiceFactory {
	return &AuthServiceFactory{
		config: config,
	}
}

// Build assembles and returns the complete auth service with enabled strategies
func (f *AuthServiceFactory) Build() (auth.Service, error) {
	// Validate required configuration
	if err := f.validateConfig(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Create JWT token manager (from usecase)
	tokenManager := usecase.NewJWTTokenManager(f.config.JWTSecret, f.config.AccessTTL, f.config.RefreshTTL)

	// Create the auth orchestrator (business logic layer)
	orchestrator := usecase.NewAuthOrchestrator(tokenManager)

	// Register enabled strategies
	if f.config.Features.EnableBasicAuth {
		basicStrategy := usecase.NewBasicAuthStrategy(f.config.UserService, tokenManager)
		orchestrator.RegisterStrategy("basic", basicStrategy)
	}

	if f.config.Features.EnableOAuth && len(f.config.OAuthProviders) > 0 {
		oauthStrategy := usecase.NewOAuthAuthStrategy(f.config.UserService, tokenManager, f.config.OAuthProviders)
		orchestrator.RegisterStrategy("oauth", oauthStrategy)
	}

	if f.config.Features.EnableJWTAuth {
		jwtStrategy := usecase.NewJWTAuthStrategy(f.config.UserService, tokenManager)
		orchestrator.RegisterStrategy("jwt", jwtStrategy)
	}

	// Return the orchestrator - pure composition, no business logic in factory
	return orchestrator, nil
}

// validateConfig validates the factory configuration
func (f *AuthServiceFactory) validateConfig() error {
	if f.config.UserService == nil {
		return fmt.Errorf("user service is required")
	}

	if len(f.config.JWTSecret) == 0 {
		return fmt.Errorf("JWT secret is required")
	}

	if f.config.AccessTTL <= 0 {
		return fmt.Errorf("access token TTL must be positive")
	}

	if f.config.RefreshTTL <= 0 {
		return fmt.Errorf("refresh token TTL must be positive")
	}

	if f.config.RefreshTTL <= f.config.AccessTTL {
		return fmt.Errorf("refresh token TTL must be longer than access token TTL")
	}

	// Validate that at least one strategy is enabled
	if !f.config.Features.EnableBasicAuth && !f.config.Features.EnableOAuth && !f.config.Features.EnableJWTAuth {
		return fmt.Errorf("at least one authentication strategy must be enabled")
	}

	// Validate OAuth configuration if enabled
	if f.config.Features.EnableOAuth && len(f.config.OAuthProviders) == 0 {
		return fmt.Errorf("OAuth providers must be configured when OAuth is enabled")
	}

	return nil
}

// Helper methods for creating common configurations

// NewDefaultConfig creates a default configuration for the auth service factory
func NewDefaultConfig(jwtSecret []byte, userService user.Service) Config {
	return Config{
		JWTSecret:      jwtSecret,
		AccessTTL:      time.Hour,
		RefreshTTL:     24 * time.Hour,
		UserService:    userService,
		OAuthProviders: make(map[string]auth.Service),
		Features:       DefaultFeatureFlags(),
	}
}

// NewTestingConfig creates a configuration suitable for testing
func NewTestingConfig(userService user.Service) Config {
	return Config{
		JWTSecret:      []byte("test-secret-key-for-testing-only"),
		AccessTTL:      time.Hour,
		RefreshTTL:     24 * time.Hour,
		UserService:    userService,
		OAuthProviders: make(map[string]auth.Service),
		Features: FeatureFlags{
			EnableBasicAuth: true,
			EnableOAuth:     false, // Disable OAuth for simpler testing
			EnableJWTAuth:   true,
		},
	}
}

// Factory is now a pure composition root - all business logic moved to usecase
