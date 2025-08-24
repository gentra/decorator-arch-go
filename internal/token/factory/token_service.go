package factory

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/gentra/decorator-arch-go/internal/token"
	"github.com/gentra/decorator-arch-go/internal/token/jwt"
)

// Config contains all configuration for building the token service
type Config struct {
	// Provider configuration
	Provider string // "jwt", "opaque", "custom"

	// JWT configuration
	JWTConfig token.TokenConfig

	// Key management
	AutoGenerateSecret bool
	SecretSize         int

	// RSA/ECDSA keys (for future implementation)
	PrivateKeyPath string
	PublicKeyPath  string

	// Token storage (for opaque tokens)
	StorageProvider string // "memory", "redis", "database"
	StorageConfig   map[string]interface{}

	// Security settings
	EnableBlacklist  bool
	BlacklistTTL     time.Duration
	EnableRotation   bool
	RotationInterval time.Duration

	// Feature flags
	Features FeatureFlags
}

// FeatureFlags controls token service behavior
type FeatureFlags struct {
	EnableJWTProvider        bool
	EnableOpaqueProvider     bool
	EnableHMACSignature      bool
	EnableRSASignature       bool
	EnableECDSASignature     bool
	EnableRefreshTokens      bool
	EnableAPITokens          bool
	EnablePasswordReset      bool
	EnableEmailVerification  bool
	EnableTokenRevocation    bool
	EnableTokenIntrospection bool
	EnableMetrics            bool
	EnableAuditLogging       bool
}

// DefaultFeatureFlags returns default feature flag configuration
func DefaultFeatureFlags() FeatureFlags {
	return FeatureFlags{
		EnableJWTProvider:        true,
		EnableOpaqueProvider:     false,
		EnableHMACSignature:      true,
		EnableRSASignature:       false,
		EnableECDSASignature:     false,
		EnableRefreshTokens:      true,
		EnableAPITokens:          true,
		EnablePasswordReset:      true,
		EnableEmailVerification:  true,
		EnableTokenRevocation:    true,
		EnableTokenIntrospection: true,
		EnableMetrics:            false,
		EnableAuditLogging:       false,
	}
}

// TokenServiceFactory creates and assembles the complete token service
type TokenServiceFactory struct {
	config Config
}

// NewFactory creates a new token service factory with the given configuration
func NewFactory(config Config) *TokenServiceFactory {
	return &TokenServiceFactory{
		config: config,
	}
}

// Build assembles and returns the complete token service based on configuration
func (f *TokenServiceFactory) Build() (token.Service, error) {
	// Prepare token configuration
	tokenConfig := f.config.JWTConfig

	// Auto-generate secret if needed
	if f.config.AutoGenerateSecret && len(tokenConfig.Secret) == 0 {
		secret, err := f.generateSecret()
		if err != nil {
			return nil, fmt.Errorf("failed to generate JWT secret: %w", err)
		}
		tokenConfig.Secret = secret
	}

	// Validate configuration
	if !tokenConfig.IsValid() {
		return nil, fmt.Errorf("invalid token configuration")
	}

	switch f.config.Provider {
	case "jwt":
		return f.buildJWTService(tokenConfig)
	case "opaque":
		return f.buildOpaqueService()
	default:
		// Default to JWT provider
		return f.buildJWTService(tokenConfig)
	}
}

// buildJWTService creates a JWT-based token service
func (f *TokenServiceFactory) buildJWTService(tokenConfig token.TokenConfig) (token.Service, error) {
	return jwt.NewService(tokenConfig)
}

// buildOpaqueService creates an opaque token service (placeholder)
func (f *TokenServiceFactory) buildOpaqueService() (token.Service, error) {
	// TODO: Implement opaque token service
	return nil, fmt.Errorf("opaque token provider not yet implemented")
}

// generateSecret generates a random secret for JWT signing
func (f *TokenServiceFactory) generateSecret() ([]byte, error) {
	secret := make([]byte, f.config.SecretSize)
	_, err := rand.Read(secret)
	return secret, err
}

// DefaultConfig returns a sensible default configuration for the token service
func DefaultConfig() Config {
	return Config{
		Provider:           "jwt",
		JWTConfig:          token.DefaultTokenConfig(),
		AutoGenerateSecret: true,
		SecretSize:         32, // 256 bits
		EnableBlacklist:    true,
		BlacklistTTL:       24 * time.Hour,
		EnableRotation:     false,
		RotationInterval:   30 * 24 * time.Hour, // 30 days
		StorageConfig:      make(map[string]interface{}),
		Features:           DefaultFeatureFlags(),
	}
}

// ConfigBuilder provides a fluent interface for building token configuration
type ConfigBuilder struct {
	config Config
}

// NewConfigBuilder creates a new configuration builder with defaults
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: DefaultConfig(),
	}
}

// WithProvider sets the token provider
func (b *ConfigBuilder) WithProvider(provider string) *ConfigBuilder {
	b.config.Provider = provider
	return b
}

// WithJWTConfig sets the JWT configuration
func (b *ConfigBuilder) WithJWTConfig(jwtConfig token.TokenConfig) *ConfigBuilder {
	b.config.JWTConfig = jwtConfig
	return b
}

// WithSecret sets the JWT signing secret
func (b *ConfigBuilder) WithSecret(secret []byte) *ConfigBuilder {
	b.config.JWTConfig.Secret = secret
	b.config.AutoGenerateSecret = false
	return b
}

// WithSecretString sets the JWT signing secret from string
func (b *ConfigBuilder) WithSecretString(secret string) *ConfigBuilder {
	b.config.JWTConfig.Secret = []byte(secret)
	b.config.AutoGenerateSecret = false
	return b
}

// WithTTLs sets token time-to-live values
func (b *ConfigBuilder) WithTTLs(accessTTL, refreshTTL, resetTTL, verificationTTL time.Duration) *ConfigBuilder {
	b.config.JWTConfig.AccessTTL = accessTTL
	b.config.JWTConfig.RefreshTTL = refreshTTL
	b.config.JWTConfig.ResetTTL = resetTTL
	b.config.JWTConfig.VerificationTTL = verificationTTL
	return b
}

// WithIssuerAndAudience sets the JWT issuer and audience
func (b *ConfigBuilder) WithIssuerAndAudience(issuer, audience string) *ConfigBuilder {
	b.config.JWTConfig.Issuer = issuer
	b.config.JWTConfig.Audience = audience
	return b
}

// WithAlgorithm sets the JWT signing algorithm
func (b *ConfigBuilder) WithAlgorithm(algorithm string) *ConfigBuilder {
	b.config.JWTConfig.Algorithm = algorithm
	return b
}

// WithRSAKeys sets RSA private and public key paths
func (b *ConfigBuilder) WithRSAKeys(privateKeyPath, publicKeyPath string) *ConfigBuilder {
	b.config.PrivateKeyPath = privateKeyPath
	b.config.PublicKeyPath = publicKeyPath
	b.config.JWTConfig.Algorithm = "RS256"
	b.config.Features.EnableRSASignature = true
	b.config.Features.EnableHMACSignature = false
	return b
}

// WithTokenRevocation enables token revocation/blacklisting
func (b *ConfigBuilder) WithTokenRevocation(enable bool, blacklistTTL time.Duration) *ConfigBuilder {
	b.config.EnableBlacklist = enable
	b.config.BlacklistTTL = blacklistTTL
	b.config.JWTConfig.EnableRevocation = enable
	b.config.Features.EnableTokenRevocation = enable
	return b
}

// WithKeyRotation enables automatic key rotation
func (b *ConfigBuilder) WithKeyRotation(enable bool, interval time.Duration) *ConfigBuilder {
	b.config.EnableRotation = enable
	b.config.RotationInterval = interval
	return b
}

// WithMaxActiveTokens sets the maximum active tokens per user
func (b *ConfigBuilder) WithMaxActiveTokens(max int) *ConfigBuilder {
	b.config.JWTConfig.MaxActiveTokens = max
	return b
}

// WithStorageProvider sets the storage provider for opaque tokens
func (b *ConfigBuilder) WithStorageProvider(provider string, config map[string]interface{}) *ConfigBuilder {
	b.config.StorageProvider = provider
	b.config.StorageConfig = config
	return b
}

// WithFeatures sets the feature flags
func (b *ConfigBuilder) WithFeatures(features FeatureFlags) *ConfigBuilder {
	b.config.Features = features
	return b
}

// EnableRefreshTokens enables refresh token functionality
func (b *ConfigBuilder) EnableRefreshTokens() *ConfigBuilder {
	b.config.JWTConfig.EnableRefresh = true
	b.config.Features.EnableRefreshTokens = true
	return b
}

// EnableAPITokens enables API token functionality
func (b *ConfigBuilder) EnableAPITokens() *ConfigBuilder {
	b.config.Features.EnableAPITokens = true
	return b
}

// EnableMetrics enables token metrics collection
func (b *ConfigBuilder) EnableMetrics() *ConfigBuilder {
	b.config.Features.EnableMetrics = true
	return b
}

// EnableAuditLogging enables audit logging for token operations
func (b *ConfigBuilder) EnableAuditLogging() *ConfigBuilder {
	b.config.Features.EnableAuditLogging = true
	return b
}

// ForDevelopment configures the service for development use
func (b *ConfigBuilder) ForDevelopment() *ConfigBuilder {
	b.config.Provider = "jwt"
	b.config.AutoGenerateSecret = true
	b.config.JWTConfig.AccessTTL = 1 * time.Hour
	b.config.JWTConfig.RefreshTTL = 24 * time.Hour
	b.config.EnableBlacklist = false
	b.config.Features.EnableMetrics = false
	b.config.Features.EnableAuditLogging = false
	return b
}

// ForProduction configures the service for production use
func (b *ConfigBuilder) ForProduction() *ConfigBuilder {
	b.config.AutoGenerateSecret = false // Should be explicitly provided
	b.config.JWTConfig.AccessTTL = 15 * time.Minute
	b.config.JWTConfig.RefreshTTL = 7 * 24 * time.Hour
	b.config.EnableBlacklist = true
	b.config.BlacklistTTL = 24 * time.Hour
	b.config.Features.EnableMetrics = true
	b.config.Features.EnableAuditLogging = true
	b.config.Features.EnableTokenRevocation = true
	return b
}

// Build returns the final configuration
func (b *ConfigBuilder) Build() Config {
	return b.config
}
