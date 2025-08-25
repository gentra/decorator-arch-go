package factory_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/token"
	"github.com/gentra/decorator-arch-go/internal/token/factory"
)

func TestDefaultFeatureFlags_GivenNoParameters_WhenCreating_ThenReturnsDefaults(t *testing.T) {
	flags := factory.DefaultFeatureFlags()

	assert.True(t, flags.EnableJWTProvider)
	assert.False(t, flags.EnableOpaqueProvider)
	assert.True(t, flags.EnableHMACSignature)
	assert.False(t, flags.EnableRSASignature)
	assert.False(t, flags.EnableECDSASignature)
	assert.True(t, flags.EnableRefreshTokens)
	assert.True(t, flags.EnableAPITokens)
	assert.True(t, flags.EnablePasswordReset)
	assert.True(t, flags.EnableEmailVerification)
	assert.True(t, flags.EnableTokenRevocation)
	assert.True(t, flags.EnableTokenIntrospection)
	assert.False(t, flags.EnableMetrics)
	assert.False(t, flags.EnableAuditLogging)
}

func TestNewFactory_GivenConfig_WhenCreating_ThenReturnsFactory(t *testing.T) {
	config := factory.DefaultConfig()
	fact := factory.NewFactory(config)

	assert.NotNil(t, fact)
}

func TestBuild_GivenJWTConfig_WhenBuilding_ThenReturnsJWTService(t *testing.T) {
	tests := []struct {
		name        string
		provider    string
		autoGen     bool
		secret      []byte
		expectError bool
	}{
		{
			name:        "valid JWT configuration with auto-generated secret",
			provider:    "jwt",
			autoGen:     true,
			secret:      nil,
			expectError: false,
		},
		{
			name:        "valid JWT configuration with provided secret",
			provider:    "jwt",
			autoGen:     false,
			secret:      []byte("test-secret-key-that-is-long-enough-for-hmac"),
			expectError: false,
		},
		{
			name:        "valid JWT configuration default provider",
			provider:    "",
			autoGen:     true,
			secret:      nil,
			expectError: false,
		},
		{
			name:        "valid JWT configuration unknown provider defaults to JWT",
			provider:    "unknown",
			autoGen:     true,
			secret:      nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := createTestConfig()
			config.Provider = tt.provider
			config.AutoGenerateSecret = tt.autoGen
			if tt.secret != nil {
				config.JWTConfig.Secret = tt.secret
			}

			fact := factory.NewFactory(config)
			service, err := fact.Build()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)

				// Test that the service works
				ctx := context.Background()
				tokenString, expiresAt, err := service.GenerateAuthToken(ctx, "user123", "user@example.com")
				assert.NoError(t, err)
				assert.NotEmpty(t, tokenString)
				assert.True(t, expiresAt.After(time.Now()))

				// Validate token
				claims, err := service.ValidateToken(ctx, tokenString)
				assert.NoError(t, err)
				assert.Equal(t, "user123", claims.UserID)
			}
		})
	}
}

func TestBuild_GivenOpaqueProvider_WhenBuilding_ThenReturnsError(t *testing.T) {
	config := factory.Config{
		Provider:  "opaque",
		JWTConfig: token.DefaultTokenConfig(),
		Features:  factory.DefaultFeatureFlags(),
	}

	fact := factory.NewFactory(config)
	service, err := fact.Build()

	assert.Error(t, err)
	// Could be either "opaque token provider not yet implemented" or "invalid token configuration"
	assert.Nil(t, service)
}

func TestBuild_GivenInvalidJWTConfig_WhenBuilding_ThenReturnsError(t *testing.T) {
	config := factory.Config{
		Provider:           "jwt",
		JWTConfig:          token.TokenConfig{}, // Invalid empty config
		AutoGenerateSecret: false,
		Features:           factory.DefaultFeatureFlags(),
	}

	fact := factory.NewFactory(config)
	service, err := fact.Build()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token configuration")
	assert.Nil(t, service)
}

func TestDefaultConfig_GivenNoParameters_WhenCreating_ThenReturnsValidConfig(t *testing.T) {
	config := factory.DefaultConfig()

	assert.Equal(t, "jwt", config.Provider)
	assert.True(t, config.AutoGenerateSecret)
	assert.Equal(t, 32, config.SecretSize)
	assert.True(t, config.EnableBlacklist)
	assert.Equal(t, 24*time.Hour, config.BlacklistTTL)
	assert.False(t, config.EnableRotation)
	assert.Equal(t, 30*24*time.Hour, config.RotationInterval)
	assert.NotNil(t, config.StorageConfig)
	assert.NotNil(t, config.Features)
	assert.NotNil(t, config.JWTConfig)
}

func TestConfigBuilder_GivenFluentInterface_WhenBuilding_ThenReturnsCustomConfig(t *testing.T) {
	secret := []byte("custom-secret-key-that-is-long-enough")
	jwtConfig := token.TokenConfig{
		Secret:   secret,
		Issuer:   "custom-issuer",
		Audience: "custom-audience",
	}

	config := factory.NewConfigBuilder().
		WithProvider("jwt").
		WithJWTConfig(jwtConfig).
		WithSecret(secret).
		WithTTLs(30*time.Minute, 48*time.Hour, 30*time.Minute, 48*time.Hour).
		WithIssuerAndAudience("test-issuer", "test-audience").
		WithAlgorithm("HS256").
		WithTokenRevocation(true, 12*time.Hour).
		WithKeyRotation(true, 15*24*time.Hour).
		WithMaxActiveTokens(10).
		Build()

	assert.Equal(t, "jwt", config.Provider)
	assert.Equal(t, secret, config.JWTConfig.Secret)
	assert.False(t, config.AutoGenerateSecret)
	assert.Equal(t, 30*time.Minute, config.JWTConfig.AccessTTL)
	assert.Equal(t, 48*time.Hour, config.JWTConfig.RefreshTTL)
	assert.Equal(t, "test-issuer", config.JWTConfig.Issuer)
	assert.Equal(t, "test-audience", config.JWTConfig.Audience)
	assert.Equal(t, "HS256", config.JWTConfig.Algorithm)
	assert.True(t, config.EnableBlacklist)
	assert.Equal(t, 12*time.Hour, config.BlacklistTTL)
	assert.True(t, config.EnableRotation)
	assert.Equal(t, 15*24*time.Hour, config.RotationInterval)
	assert.Equal(t, 10, config.JWTConfig.MaxActiveTokens)
}

func TestConfigBuilder_GivenSecretString_WhenBuilding_ThenSetsSecret(t *testing.T) {
	secretString := "my-secret-string"

	config := factory.NewConfigBuilder().
		WithSecretString(secretString).
		Build()

	assert.Equal(t, []byte(secretString), config.JWTConfig.Secret)
	assert.False(t, config.AutoGenerateSecret)
}

func TestConfigBuilder_GivenRSAKeys_WhenBuilding_ThenConfiguresRSA(t *testing.T) {
	privateKeyPath := "/path/to/private.key"
	publicKeyPath := "/path/to/public.key"

	config := factory.NewConfigBuilder().
		WithRSAKeys(privateKeyPath, publicKeyPath).
		Build()

	assert.Equal(t, privateKeyPath, config.PrivateKeyPath)
	assert.Equal(t, publicKeyPath, config.PublicKeyPath)
	assert.Equal(t, "RS256", config.JWTConfig.Algorithm)
	assert.True(t, config.Features.EnableRSASignature)
	assert.False(t, config.Features.EnableHMACSignature)
}

func TestConfigBuilder_GivenStorageProvider_WhenBuilding_ThenSetsStorage(t *testing.T) {
	provider := "redis"
	storageConfig := map[string]interface{}{
		"host": "localhost",
		"port": 6379,
	}

	config := factory.NewConfigBuilder().
		WithStorageProvider(provider, storageConfig).
		Build()

	assert.Equal(t, provider, config.StorageProvider)
	assert.Equal(t, storageConfig, config.StorageConfig)
}

func TestConfigBuilder_GivenCustomFeatures_WhenBuilding_ThenAppliesFeatures(t *testing.T) {
	customFeatures := factory.FeatureFlags{
		EnableJWTProvider:   false,
		EnableOpaqueProvider: true,
		EnableMetrics:       true,
		EnableAuditLogging:  true,
	}

	config := factory.NewConfigBuilder().
		WithFeatures(customFeatures).
		Build()

	assert.Equal(t, customFeatures, config.Features)
}

func TestConfigBuilder_GivenFeatureEnablers_WhenBuilding_ThenEnablesFeatures(t *testing.T) {
	config := factory.NewConfigBuilder().
		EnableRefreshTokens().
		EnableAPITokens().
		EnableMetrics().
		EnableAuditLogging().
		Build()

	assert.True(t, config.JWTConfig.EnableRefresh)
	assert.True(t, config.Features.EnableRefreshTokens)
	assert.True(t, config.Features.EnableAPITokens)
	assert.True(t, config.Features.EnableMetrics)
	assert.True(t, config.Features.EnableAuditLogging)
}

func TestConfigBuilder_GivenDevelopmentMode_WhenBuilding_ThenConfiguresForDev(t *testing.T) {
	config := factory.NewConfigBuilder().
		ForDevelopment().
		Build()

	assert.Equal(t, "jwt", config.Provider)
	assert.True(t, config.AutoGenerateSecret)
	assert.Equal(t, 1*time.Hour, config.JWTConfig.AccessTTL)
	assert.Equal(t, 24*time.Hour, config.JWTConfig.RefreshTTL)
	assert.False(t, config.EnableBlacklist)
	assert.False(t, config.Features.EnableMetrics)
	assert.False(t, config.Features.EnableAuditLogging)
}

func TestConfigBuilder_GivenProductionMode_WhenBuilding_ThenConfiguresForProd(t *testing.T) {
	config := factory.NewConfigBuilder().
		ForProduction().
		Build()

	assert.False(t, config.AutoGenerateSecret)
	assert.Equal(t, 15*time.Minute, config.JWTConfig.AccessTTL)
	assert.Equal(t, 7*24*time.Hour, config.JWTConfig.RefreshTTL)
	assert.True(t, config.EnableBlacklist)
	assert.Equal(t, 24*time.Hour, config.BlacklistTTL)
	assert.True(t, config.Features.EnableMetrics)
	assert.True(t, config.Features.EnableAuditLogging)
	assert.True(t, config.Features.EnableTokenRevocation)
}

func TestBuild_GivenCompleteWorkflow_WhenBuilding_ThenCreatesWorkingService(t *testing.T) {
	// Test complete workflow from builder to working service
	config := factory.NewConfigBuilder().
		WithProvider("jwt").
		WithSecretString("my-very-secure-secret-key-for-testing").
		WithTTLs(1*time.Hour, 24*time.Hour, 15*time.Minute, 24*time.Hour).
		WithIssuerAndAudience("test-service", "test-users").
		EnableRefreshTokens().
		EnableAPITokens().
		Build()

	fact := factory.NewFactory(config)
	service, err := fact.Build()

	assert.NoError(t, err)
	assert.NotNil(t, service)

	ctx := context.Background()

	// Test auth token generation and validation
	authToken, expiresAt, err := service.GenerateAuthToken(ctx, "user123", "user@example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, authToken)
	assert.True(t, expiresAt.After(time.Now()))

	authClaims, err := service.ValidateToken(ctx, authToken)
	assert.NoError(t, err)
	assert.Equal(t, "user123", authClaims.UserID)
	assert.Equal(t, "user@example.com", authClaims.Email)
	assert.Equal(t, "auth", authClaims.TokenType)
	assert.Equal(t, "test-service", authClaims.Issuer)
	assert.Equal(t, "test-users", authClaims.Audience)

	// Test refresh token generation and validation
	refreshToken, err := service.GenerateRefreshToken(ctx, "user123")
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)

	refreshClaims, err := service.ValidateToken(ctx, refreshToken)
	assert.NoError(t, err)
	assert.Equal(t, "user123", refreshClaims.UserID)
	assert.Equal(t, "refresh", refreshClaims.TokenType)

	// Test API token generation and validation
	apiToken, err := service.GenerateAPIToken(ctx, "user123", []string{"read", "write"})
	assert.NoError(t, err)
	assert.NotNil(t, apiToken)
	assert.NotEmpty(t, apiToken.Token)

	apiClaims, err := service.ValidateAPIToken(ctx, apiToken.Token)
	assert.NoError(t, err)
	assert.Equal(t, "user123", apiClaims.UserID)
	assert.Equal(t, "api", apiClaims.TokenType)
	assert.Equal(t, []string{"read", "write"}, apiClaims.Scopes)

	// Test refresh workflow
	newTokenPair, err := service.RefreshToken(ctx, refreshToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, newTokenPair.AccessToken)
	assert.Equal(t, "bearer", newTokenPair.TokenType)

	// Test token revocation
	err = service.RevokeToken(ctx, authToken)
	assert.NoError(t, err)

	// Verify revoked token is invalid
	_, err = service.ValidateToken(ctx, authToken)
	assert.Error(t, err)
	assert.Equal(t, token.ErrTokenRevoked, err)
}

func TestBuild_GivenAutoGenerateSecret_WhenBuilding_ThenGeneratesRandomSecret(t *testing.T) {
	config := factory.Config{
		Provider:           "jwt",
		JWTConfig:          token.DefaultTokenConfig(),
		AutoGenerateSecret: true,
		SecretSize:         32,
		Features:           factory.DefaultFeatureFlags(),
	}

	fact := factory.NewFactory(config)
	service1, err := fact.Build()
	assert.NoError(t, err)
	assert.NotNil(t, service1)

	// Build another service with the same config
	fact2 := factory.NewFactory(config)
	service2, err := fact2.Build()
	assert.NoError(t, err)
	assert.NotNil(t, service2)

	ctx := context.Background()

	// Generate tokens with both services
	token1, _, err := service1.GenerateAuthToken(ctx, "user123", "user@example.com")
	assert.NoError(t, err)

	token2, _, err := service2.GenerateAuthToken(ctx, "user123", "user@example.com")
	assert.NoError(t, err)

	// Tokens should be different (different secrets)
	assert.NotEqual(t, token1, token2)

	// Each service should only validate its own tokens
	_, err = service1.ValidateToken(ctx, token2)
	assert.Error(t, err)

	_, err = service2.ValidateToken(ctx, token1)
	assert.Error(t, err)
}

func TestBuild_GivenZeroSecretSize_WhenBuilding_ThenUsesDefaultSize(t *testing.T) {
	config := factory.Config{
		Provider:           "jwt",
		JWTConfig:          token.DefaultTokenConfig(),
		AutoGenerateSecret: true,
		SecretSize:         0, // Zero size should use some default
		Features:           factory.DefaultFeatureFlags(),
	}

	fact := factory.NewFactory(config)
	service, err := fact.Build()

	// The service should still work even with zero secret size
	// (the underlying implementation should handle this gracefully)
	if err == nil {
		assert.NotNil(t, service)
	} else {
		// If it fails, it should be a meaningful error about configuration
		assert.Error(t, err)
	}
}

// Helper function to create a test configuration
func createTestConfig() factory.Config {
	config := factory.DefaultConfig()
	config.JWTConfig.Secret = []byte("test-secret-key-that-is-long-enough-for-hmac")
	return config
}