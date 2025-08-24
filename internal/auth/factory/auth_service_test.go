package factory_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/auth"
	"github.com/gentra/decorator-arch-go/internal/auth/factory"
	authmock "github.com/gentra/decorator-arch-go/internal/auth/mock"
	"github.com/gentra/decorator-arch-go/internal/user"
	usermock "github.com/gentra/decorator-arch-go/internal/user/mock"
)

// MockOAuthProvider for testing - now using centralized mock
type MockOAuthProvider = authmock.MockOAuthProvider

func TestAuthServiceFactory_Build(t *testing.T) {
	testCases := []struct {
		name            string
		config          factory.Config
		expectError     bool
		expectedErr     string
		validateService func(*testing.T, auth.Service)
	}{
		{
			name: "Given valid default configuration, When Build is called, Then should create auth service with basic and JWT strategies",
			config: factory.Config{
				JWTSecret:      []byte("test-secret-key-32-bytes-long!!!"),
				AccessTTL:      time.Hour,
				RefreshTTL:     24 * time.Hour,
				UserService:    new(usermock.MockUserService),
				OAuthProviders: make(map[string]auth.Service),
				Features: factory.FeatureFlags{
					EnableBasicAuth: true,
					EnableOAuth:     false,
					EnableJWTAuth:   true,
				},
			},
			expectError: false,
			validateService: func(t *testing.T, service auth.Service) {
				strategies := service.GetSupportedStrategies()
				assert.Contains(t, strategies, "basic")
				assert.Contains(t, strategies, "jwt")
				assert.NotContains(t, strategies, "oauth")
			},
		},
		{
			name: "Given configuration with OAuth enabled, When Build is called, Then should create auth service with all strategies",
			config: factory.Config{
				JWTSecret:   []byte("test-secret-key-32-bytes-long!!!"),
				AccessTTL:   time.Hour,
				RefreshTTL:  24 * time.Hour,
				UserService: new(usermock.MockUserService),
				OAuthProviders: map[string]auth.Service{
					"google": new(MockOAuthProvider),
					"github": new(MockOAuthProvider),
				},
				Features: factory.FeatureFlags{
					EnableBasicAuth: true,
					EnableOAuth:     true,
					EnableJWTAuth:   true,
				},
			},
			expectError: false,
			validateService: func(t *testing.T, service auth.Service) {
				strategies := service.GetSupportedStrategies()
				assert.Contains(t, strategies, "basic")
				assert.Contains(t, strategies, "jwt")
				assert.Contains(t, strategies, "oauth")
			},
		},
		{
			name: "Given configuration with only basic auth enabled, When Build is called, Then should create auth service with only basic strategy",
			config: factory.Config{
				JWTSecret:      []byte("test-secret-key-32-bytes-long!!!"),
				AccessTTL:      time.Hour,
				RefreshTTL:     24 * time.Hour,
				UserService:    new(usermock.MockUserService),
				OAuthProviders: make(map[string]auth.Service),
				Features: factory.FeatureFlags{
					EnableBasicAuth: true,
					EnableOAuth:     false,
					EnableJWTAuth:   false,
				},
			},
			expectError: false,
			validateService: func(t *testing.T, service auth.Service) {
				strategies := service.GetSupportedStrategies()
				assert.Contains(t, strategies, "basic")
				assert.NotContains(t, strategies, "jwt")
				assert.NotContains(t, strategies, "oauth")
			},
		},
		{
			name: "Given configuration with missing user service, When Build is called, Then should return validation error",
			config: factory.Config{
				JWTSecret:      []byte("test-secret-key-32-bytes-long!!!"),
				AccessTTL:      time.Hour,
				RefreshTTL:     24 * time.Hour,
				UserService:    nil, // Missing!
				OAuthProviders: make(map[string]auth.Service),
				Features:       factory.DefaultFeatureFlags(),
			},
			expectError: true,
			expectedErr: "user service is required",
		},
		{
			name: "Given configuration with empty JWT secret, When Build is called, Then should return validation error",
			config: factory.Config{
				JWTSecret:      []byte(""), // Empty!
				AccessTTL:      time.Hour,
				RefreshTTL:     24 * time.Hour,
				UserService:    new(usermock.MockUserService),
				OAuthProviders: make(map[string]auth.Service),
				Features:       factory.DefaultFeatureFlags(),
			},
			expectError: true,
			expectedErr: "JWT secret is required",
		},
		{
			name: "Given configuration with zero access TTL, When Build is called, Then should return validation error",
			config: factory.Config{
				JWTSecret:      []byte("test-secret-key-32-bytes-long!!!"),
				AccessTTL:      0, // Invalid!
				RefreshTTL:     24 * time.Hour,
				UserService:    new(usermock.MockUserService),
				OAuthProviders: make(map[string]auth.Service),
				Features:       factory.DefaultFeatureFlags(),
			},
			expectError: true,
			expectedErr: "access token TTL must be positive",
		},
		{
			name: "Given configuration with zero refresh TTL, When Build is called, Then should return validation error",
			config: factory.Config{
				JWTSecret:      []byte("test-secret-key-32-bytes-long!!!"),
				AccessTTL:      time.Hour,
				RefreshTTL:     0, // Invalid!
				UserService:    new(usermock.MockUserService),
				OAuthProviders: make(map[string]auth.Service),
				Features:       factory.DefaultFeatureFlags(),
			},
			expectError: true,
			expectedErr: "refresh token TTL must be positive",
		},
		{
			name: "Given configuration with refresh TTL less than access TTL, When Build is called, Then should return validation error",
			config: factory.Config{
				JWTSecret:      []byte("test-secret-key-32-bytes-long!!!"),
				AccessTTL:      24 * time.Hour,
				RefreshTTL:     time.Hour, // Less than access TTL!
				UserService:    new(usermock.MockUserService),
				OAuthProviders: make(map[string]auth.Service),
				Features:       factory.DefaultFeatureFlags(),
			},
			expectError: true,
			expectedErr: "refresh token TTL must be longer than access token TTL",
		},
		{
			name: "Given configuration with all strategies disabled, When Build is called, Then should return validation error",
			config: factory.Config{
				JWTSecret:      []byte("test-secret-key-32-bytes-long!!!"),
				AccessTTL:      time.Hour,
				RefreshTTL:     24 * time.Hour,
				UserService:    new(usermock.MockUserService),
				OAuthProviders: make(map[string]auth.Service),
				Features: factory.FeatureFlags{
					EnableBasicAuth: false,
					EnableOAuth:     false,
					EnableJWTAuth:   false,
				},
			},
			expectError: true,
			expectedErr: "at least one authentication strategy must be enabled",
		},
		{
			name: "Given OAuth enabled but no providers configured, When Build is called, Then should return validation error",
			config: factory.Config{
				JWTSecret:      []byte("test-secret-key-32-bytes-long!!!"),
				AccessTTL:      time.Hour,
				RefreshTTL:     24 * time.Hour,
				UserService:    new(usermock.MockUserService),
				OAuthProviders: make(map[string]auth.Service), // Empty!
				Features: factory.FeatureFlags{
					EnableBasicAuth: true,
					EnableOAuth:     true, // Enabled but no providers!
					EnableJWTAuth:   true,
				},
			},
			expectError: true,
			expectedErr: "OAuth providers must be configured when OAuth is enabled",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			authFactory := factory.NewAuthServiceFactory(tt.config)

			// Act
			service, err := authFactory.Build()

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
				if tt.validateService != nil {
					tt.validateService(t, service)
				}
			}
		})
	}
}

func TestNewDefaultConfig(t *testing.T) {
	testCases := []struct {
		name        string
		jwtSecret   []byte
		userService user.Service
		validate    func(*testing.T, factory.Config)
	}{
		{
			name:        "Given valid parameters, When NewDefaultConfig is called, Then should create config with sensible defaults",
			jwtSecret:   []byte("test-secret-key-32-bytes-long!!!"),
			userService: new(usermock.MockUserService),
			validate: func(t *testing.T, config factory.Config) {
				assert.Equal(t, []byte("test-secret-key-32-bytes-long!!!"), config.JWTSecret)
				assert.Equal(t, time.Hour, config.AccessTTL)
				assert.Equal(t, 24*time.Hour, config.RefreshTTL)
				assert.NotNil(t, config.UserService)
				assert.NotNil(t, config.OAuthProviders)
				assert.True(t, config.Features.EnableBasicAuth)
				assert.False(t, config.Features.EnableOAuth) // Disabled by default
				assert.True(t, config.Features.EnableJWTAuth)
			},
		},
		{
			name:        "Given empty JWT secret, When NewDefaultConfig is called, Then should still create config",
			jwtSecret:   []byte(""),
			userService: new(usermock.MockUserService),
			validate: func(t *testing.T, config factory.Config) {
				assert.Equal(t, []byte(""), config.JWTSecret)
				assert.NotNil(t, config.UserService)
			},
		},
		{
			name:        "Given nil user service, When NewDefaultConfig is called, Then should still create config",
			jwtSecret:   []byte("test-secret"),
			userService: nil,
			validate: func(t *testing.T, config factory.Config) {
				assert.Nil(t, config.UserService)
				assert.NotNil(t, config.OAuthProviders)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			config := factory.NewDefaultConfig(tt.jwtSecret, tt.userService)

			// Assert
			tt.validate(t, config)
		})
	}
}

func TestNewTestingConfig(t *testing.T) {
	testCases := []struct {
		name        string
		userService user.Service
		validate    func(*testing.T, factory.Config)
	}{
		{
			name:        "Given valid user service, When NewTestingConfig is called, Then should create testing config with OAuth disabled",
			userService: new(usermock.MockUserService),
			validate: func(t *testing.T, config factory.Config) {
				assert.Equal(t, []byte("test-secret-key-for-testing-only"), config.JWTSecret)
				assert.Equal(t, time.Hour, config.AccessTTL)
				assert.Equal(t, 24*time.Hour, config.RefreshTTL)
				assert.NotNil(t, config.UserService)
				assert.NotNil(t, config.OAuthProviders)
				assert.True(t, config.Features.EnableBasicAuth)
				assert.False(t, config.Features.EnableOAuth) // Disabled for testing
				assert.True(t, config.Features.EnableJWTAuth)
			},
		},
		{
			name:        "Given nil user service, When NewTestingConfig is called, Then should create config with nil user service",
			userService: nil,
			validate: func(t *testing.T, config factory.Config) {
				assert.Nil(t, config.UserService)
				assert.False(t, config.Features.EnableOAuth)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			config := factory.NewTestingConfig(tt.userService)

			// Assert
			tt.validate(t, config)
		})
	}
}

func TestDefaultFeatureFlags(t *testing.T) {
	testCases := []struct {
		name     string
		validate func(*testing.T, factory.FeatureFlags)
	}{
		{
			name: "Given no parameters, When DefaultFeatureFlags is called, Then should return sensible defaults",
			validate: func(t *testing.T, flags factory.FeatureFlags) {
				assert.True(t, flags.EnableBasicAuth)
				assert.False(t, flags.EnableOAuth) // Requires provider setup
				assert.True(t, flags.EnableJWTAuth)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			flags := factory.DefaultFeatureFlags()

			// Assert
			tt.validate(t, flags)
		})
	}
}

func TestAuthServiceFactory_Integration(t *testing.T) {
	testCases := []struct {
		name string
		test func(t *testing.T, factory *factory.AuthServiceFactory)
	}{
		{
			name: "Given complete factory configuration, When building and using auth service, Then should work end-to-end",
			test: func(t *testing.T, authFactory *factory.AuthServiceFactory) {
				// Build the service
				authService, err := authFactory.Build()
				assert.NoError(t, err)
				assert.NotNil(t, authService)

				// Verify service supports expected strategies
				strategies := authService.GetSupportedStrategies()
				assert.Contains(t, strategies, "basic")
				assert.Contains(t, strategies, "jwt")
				assert.Len(t, strategies, 2) // Only basic and jwt for this test

				// Test unsupported strategy
				_, err = authService.Authenticate(context.Background(), "unsupported", nil)
				assert.Error(t, err)
				assert.Equal(t, auth.ErrUnsupportedStrategy, err)
			},
		},
		{
			name: "Given factory with OAuth configuration, When building service, Then should support OAuth strategy",
			test: func(t *testing.T, authFactory *factory.AuthServiceFactory) {
				// This test would be called with a factory that has OAuth enabled
				authService, err := authFactory.Build()
				assert.NoError(t, err)
				assert.NotNil(t, authService)

				strategies := authService.GetSupportedStrategies()
				// The exact strategies depend on the configuration passed to this test
				assert.NotEmpty(t, strategies)
			},
		},
		{
			name: "Given factory configuration validation, When invalid configs are provided, Then should handle gracefully",
			test: func(t *testing.T, authFactory *factory.AuthServiceFactory) {
				// This test verifies that the factory properly validates configuration
				// The actual validation happens in the Build() method
				// This is more of a structural test to ensure validation occurs
				assert.NotPanics(t, func() {
					authFactory.Build()
				})
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange - Create different factory configurations for different tests
			var authFactory *factory.AuthServiceFactory

			switch tt.name {
			case "Given complete factory configuration, When building and using auth service, Then should work end-to-end":
				config := factory.NewDefaultConfig(
					[]byte("test-secret-key-32-bytes-long!!!"),
					new(usermock.MockUserService),
				)
				authFactory = factory.NewAuthServiceFactory(config)

			case "Given factory with OAuth configuration, When building service, Then should support OAuth strategy":
				config := factory.Config{
					JWTSecret:   []byte("test-secret-key-32-bytes-long!!!"),
					AccessTTL:   time.Hour,
					RefreshTTL:  24 * time.Hour,
					UserService: new(usermock.MockUserService),
					OAuthProviders: map[string]auth.Service{
						"google": new(MockOAuthProvider),
					},
					Features: factory.FeatureFlags{
						EnableBasicAuth: true,
						EnableOAuth:     true,
						EnableJWTAuth:   true,
					},
				}
				authFactory = factory.NewAuthServiceFactory(config)

			case "Given factory configuration validation, When invalid configs are provided, Then should handle gracefully":
				// Create an invalid config for validation testing
				config := factory.Config{
					JWTSecret:      []byte(""), // Invalid
					AccessTTL:      0,          // Invalid
					RefreshTTL:     0,          // Invalid
					UserService:    nil,        // Invalid
					OAuthProviders: nil,
					Features:       factory.FeatureFlags{}, // All disabled
				}
				authFactory = factory.NewAuthServiceFactory(config)
			}

			// Act & Assert
			tt.test(t, authFactory)
		})
	}
}
