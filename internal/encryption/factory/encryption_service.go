package factory

import (
	"crypto/rand"
	"fmt"

	"github.com/gentra/decorator-arch-go/internal/encryption"
	"github.com/gentra/decorator-arch-go/internal/encryption/aes"
	"github.com/gentra/decorator-arch-go/internal/encryption/noop"
)

// Config contains all configuration for building the encryption service
type Config struct {
	// Encryption algorithm
	Algorithm string // "aes", "noop"

	// Key management
	DefaultKey  []byte
	PurposeKeys map[string][]byte

	// Key generation settings
	KeySize int // For AES: 16, 24, or 32 bytes (AES-128, AES-192, AES-256)

	// Security settings
	AutoGenerateKeys bool
	KeyRotationDays  int

	// Feature flags
	Features FeatureFlags
}

// FeatureFlags controls encryption service behavior
type FeatureFlags struct {
	EnableAESEncryption    bool
	EnableNoOpEncryption   bool
	EnableKeyRotation      bool
	EnablePurposeKeys      bool
	EnableBatchOperations  bool
	EnableKeyDerivation    bool
	EnableCompressionFirst bool
}

// DefaultFeatureFlags returns default feature flag configuration
func DefaultFeatureFlags() FeatureFlags {
	return FeatureFlags{
		EnableAESEncryption:    true,
		EnableNoOpEncryption:   false,
		EnableKeyRotation:      true,
		EnablePurposeKeys:      true,
		EnableBatchOperations:  true,
		EnableKeyDerivation:    false,
		EnableCompressionFirst: false,
	}
}

// EncryptionServiceFactory creates and assembles the complete encryption service
type EncryptionServiceFactory struct {
	config Config
}

// NewFactory creates a new encryption service factory with the given configuration
func NewFactory(config Config) *EncryptionServiceFactory {
	return &EncryptionServiceFactory{
		config: config,
	}
}

// Build assembles and returns the complete encryption service based on configuration
func (f *EncryptionServiceFactory) Build() (encryption.Service, error) {
	switch f.config.Algorithm {
	case "aes":
		return f.buildAESService()
	case "noop":
		return f.buildNoOpService()
	default:
		// Default to AES if algorithm not specified or invalid
		return f.buildAESService()
	}
}

// buildAESService creates an AES-based encryption service
func (f *EncryptionServiceFactory) buildAESService() (encryption.Service, error) {
	// Generate default key if not provided and auto-generation is enabled
	defaultKey := f.config.DefaultKey
	if len(defaultKey) == 0 && f.config.AutoGenerateKeys {
		var err error
		defaultKey, err = f.generateKey()
		if err != nil {
			return nil, fmt.Errorf("failed to generate default key: %w", err)
		}
	}

	// Validate key size
	if len(defaultKey) != f.config.KeySize {
		return nil, fmt.Errorf("default key size must be %d bytes, got %d", f.config.KeySize, len(defaultKey))
	}

	// Handle purpose keys
	purposeKeys := f.config.PurposeKeys
	if purposeKeys == nil {
		purposeKeys = make(map[string][]byte)
	}

	// Auto-generate purpose keys if enabled and not provided
	if f.config.AutoGenerateKeys && f.config.Features.EnablePurposeKeys {
		if err := f.generatePurposeKeys(purposeKeys); err != nil {
			return nil, fmt.Errorf("failed to generate purpose keys: %w", err)
		}
	}

	// Create AES service
	if f.config.Features.EnablePurposeKeys && len(purposeKeys) > 0 {
		return aes.NewService(purposeKeys, defaultKey)
	}

	// Fallback to default key only
	singleKeyMap := map[string][]byte{
		encryption.PurposeDefault: defaultKey,
	}
	return aes.NewService(singleKeyMap, defaultKey)
}

// buildNoOpService creates a no-operation encryption service (for development/testing)
func (f *EncryptionServiceFactory) buildNoOpService() (encryption.Service, error) {
	return noop.NewService(), nil
}

// generateKey generates a random encryption key of the configured size
func (f *EncryptionServiceFactory) generateKey() ([]byte, error) {
	key := make([]byte, f.config.KeySize)
	_, err := rand.Read(key)
	return key, err
}

// generatePurposeKeys generates keys for standard purposes if they don't exist
func (f *EncryptionServiceFactory) generatePurposeKeys(purposeKeys map[string][]byte) error {
	standardPurposes := []string{
		encryption.PurposeUserEmail,
		encryption.PurposeUserName,
		encryption.PurposeUserPhone,
		encryption.PurposePaymentCard,
		encryption.PurposeDocumentContent,
		encryption.PurposeSecretAPIKey,
	}

	for _, purpose := range standardPurposes {
		if _, exists := purposeKeys[purpose]; !exists {
			key, err := f.generateKey()
			if err != nil {
				return fmt.Errorf("failed to generate key for purpose %s: %w", purpose, err)
			}
			purposeKeys[purpose] = key
		}
	}

	return nil
}

// DefaultConfig returns a sensible default configuration for the encryption service
func DefaultConfig() Config {
	return Config{
		Algorithm:        "aes",
		KeySize:          32, // AES-256
		AutoGenerateKeys: true,
		KeyRotationDays:  90,
		PurposeKeys:      make(map[string][]byte),
		Features:         DefaultFeatureFlags(),
	}
}

// ConfigBuilder provides a fluent interface for building encryption configuration
type ConfigBuilder struct {
	config Config
}

// NewConfigBuilder creates a new configuration builder with defaults
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: DefaultConfig(),
	}
}

// WithAlgorithm sets the encryption algorithm
func (b *ConfigBuilder) WithAlgorithm(algorithm string) *ConfigBuilder {
	b.config.Algorithm = algorithm
	return b
}

// WithDefaultKey sets the default encryption key
func (b *ConfigBuilder) WithDefaultKey(key []byte) *ConfigBuilder {
	b.config.DefaultKey = key
	return b
}

// WithPurposeKey adds a purpose-specific encryption key
func (b *ConfigBuilder) WithPurposeKey(purpose string, key []byte) *ConfigBuilder {
	if b.config.PurposeKeys == nil {
		b.config.PurposeKeys = make(map[string][]byte)
	}
	b.config.PurposeKeys[purpose] = key
	return b
}

// WithKeySize sets the encryption key size
func (b *ConfigBuilder) WithKeySize(size int) *ConfigBuilder {
	b.config.KeySize = size
	return b
}

// WithAutoGenerateKeys enables automatic key generation
func (b *ConfigBuilder) WithAutoGenerateKeys(enable bool) *ConfigBuilder {
	b.config.AutoGenerateKeys = enable
	return b
}

// WithKeyRotationDays sets the key rotation period in days
func (b *ConfigBuilder) WithKeyRotationDays(days int) *ConfigBuilder {
	b.config.KeyRotationDays = days
	return b
}

// WithFeatures sets the feature flags
func (b *ConfigBuilder) WithFeatures(features FeatureFlags) *ConfigBuilder {
	b.config.Features = features
	return b
}

// EnableNoOpMode switches to no-operation encryption (for development)
func (b *ConfigBuilder) EnableNoOpMode() *ConfigBuilder {
	b.config.Algorithm = "noop"
	b.config.Features.EnableNoOpEncryption = true
	b.config.Features.EnableAESEncryption = false
	return b
}

// EnableProductionMode switches to AES encryption with secure defaults
func (b *ConfigBuilder) EnableProductionMode() *ConfigBuilder {
	b.config.Algorithm = "aes"
	b.config.KeySize = 32 // AES-256
	b.config.AutoGenerateKeys = true
	b.config.Features.EnableAESEncryption = true
	b.config.Features.EnableNoOpEncryption = false
	b.config.Features.EnableKeyRotation = true
	b.config.Features.EnablePurposeKeys = true
	return b
}

// Build returns the final configuration
func (b *ConfigBuilder) Build() Config {
	return b.config
}
