package factory_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/encryption"
	"github.com/gentra/decorator-arch-go/internal/encryption/factory"
)

func TestDefaultFeatureFlags_GivenNoParameters_WhenCreating_ThenReturnsDefaults(t *testing.T) {
	flags := factory.DefaultFeatureFlags()

	assert.True(t, flags.EnableAESEncryption)
	assert.False(t, flags.EnableNoOpEncryption)
	assert.True(t, flags.EnableKeyRotation)
	assert.True(t, flags.EnablePurposeKeys)
	assert.True(t, flags.EnableBatchOperations)
	assert.False(t, flags.EnableKeyDerivation)
	assert.False(t, flags.EnableCompressionFirst)
}

func TestNewFactory_GivenConfig_WhenCreating_ThenReturnsFactory(t *testing.T) {
	config := factory.DefaultConfig()
	fact := factory.NewFactory(config)

	assert.NotNil(t, fact)
}

func TestBuild_GivenAESConfig_WhenBuilding_ThenReturnsAESService(t *testing.T) {
	tests := []struct {
		name      string
		algorithm string
		keySize   int
		autoGen   bool
		expectErr bool
	}{
		{
			name:      "valid AES configuration",
			algorithm: "aes",
			keySize:   32,
			autoGen:   true,
			expectErr: false,
		},
		{
			name:      "valid AES with manual keys",
			algorithm: "aes",
			keySize:   32,
			autoGen:   false,
			expectErr: false,
		},
		{
			name:      "invalid key size",
			algorithm: "aes",
			keySize:   16,
			autoGen:   false,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := factory.Config{
				Algorithm:        tt.algorithm,
				KeySize:          tt.keySize,
				AutoGenerateKeys: tt.autoGen,
				Features:         factory.DefaultFeatureFlags(),
			}

			if !tt.autoGen {
				config.DefaultKey = make([]byte, tt.keySize)
			}

			fact := factory.NewFactory(config)
			service, err := fact.Build()

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)

				// Test that the service works
				ctx := context.Background()
				plaintext := "test data"
				encrypted, err := service.Encrypt(ctx, plaintext)
				assert.NoError(t, err)

				decrypted, err := service.Decrypt(ctx, encrypted)
				assert.NoError(t, err)
				assert.Equal(t, plaintext, decrypted)
			}
		})
	}
}

func TestBuild_GivenNoOpConfig_WhenBuilding_ThenReturnsNoOpService(t *testing.T) {
	config := factory.Config{
		Algorithm: "noop",
		Features:  factory.DefaultFeatureFlags(),
	}

	fact := factory.NewFactory(config)
	service, err := fact.Build()

	assert.NoError(t, err)
	assert.NotNil(t, service)

	// Test that the service is no-op (returns plaintext as-is)
	ctx := context.Background()
	plaintext := "test data"
	encrypted, err := service.Encrypt(ctx, plaintext)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, encrypted)

	decrypted, err := service.Decrypt(ctx, encrypted)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestBuild_GivenUnknownAlgorithm_WhenBuilding_ThenDefaultsToAES(t *testing.T) {
	config := factory.Config{
		Algorithm:        "unknown",
		KeySize:          32,
		AutoGenerateKeys: true,
		Features:         factory.DefaultFeatureFlags(),
	}

	fact := factory.NewFactory(config)
	service, err := fact.Build()

	assert.NoError(t, err)
	assert.NotNil(t, service)

	// Should behave like AES service (encrypted data is different from plaintext)
	ctx := context.Background()
	plaintext := "test data"
	encrypted, err := service.Encrypt(ctx, plaintext)
	assert.NoError(t, err)
	assert.NotEqual(t, plaintext, encrypted)

	decrypted, err := service.Decrypt(ctx, encrypted)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestBuild_GivenPurposeKeysEnabled_WhenBuilding_ThenSupportsPurposeEncryption(t *testing.T) {
	config := factory.Config{
		Algorithm:        "aes",
		KeySize:          32,
		AutoGenerateKeys: true,
		Features: factory.FeatureFlags{
			EnableAESEncryption: true,
			EnablePurposeKeys:   true,
		},
	}

	fact := factory.NewFactory(config)
	service, err := fact.Build()

	assert.NoError(t, err)
	assert.NotNil(t, service)

	// Test purpose-based encryption
	ctx := context.Background()
	plaintext := "test data"
	purpose := encryption.PurposeUserEmail

	encrypted, err := service.EncryptWithPurpose(ctx, plaintext, purpose)
	assert.NoError(t, err)
	assert.NotEqual(t, plaintext, encrypted)

	decrypted, err := service.DecryptWithPurpose(ctx, encrypted, purpose)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestDefaultConfig_GivenNoParameters_WhenCreating_ThenReturnsValidConfig(t *testing.T) {
	config := factory.DefaultConfig()

	assert.Equal(t, "aes", config.Algorithm)
	assert.Equal(t, 32, config.KeySize)
	assert.True(t, config.AutoGenerateKeys)
	assert.Equal(t, 90, config.KeyRotationDays)
	assert.NotNil(t, config.PurposeKeys)
	assert.NotNil(t, config.Features)
}

func TestConfigBuilder_GivenFluentInterface_WhenBuilding_ThenReturnsCustomConfig(t *testing.T) {
	customKey := make([]byte, 32)
	copy(customKey, []byte("12345678901234567890123456789012"))

	config := factory.NewConfigBuilder().
		WithAlgorithm("aes").
		WithKeySize(32).
		WithDefaultKey(customKey).
		WithPurposeKey(encryption.PurposeUserEmail, customKey).
		WithAutoGenerateKeys(false).
		WithKeyRotationDays(30).
		Build()

	assert.Equal(t, "aes", config.Algorithm)
	assert.Equal(t, 32, config.KeySize)
	assert.Equal(t, customKey, config.DefaultKey)
	assert.Equal(t, customKey, config.PurposeKeys[encryption.PurposeUserEmail])
	assert.False(t, config.AutoGenerateKeys)
	assert.Equal(t, 30, config.KeyRotationDays)
}

func TestConfigBuilder_GivenNoOpMode_WhenBuilding_ThenConfiguresNoOp(t *testing.T) {
	config := factory.NewConfigBuilder().
		EnableNoOpMode().
		Build()

	assert.Equal(t, "noop", config.Algorithm)
	assert.True(t, config.Features.EnableNoOpEncryption)
	assert.False(t, config.Features.EnableAESEncryption)
}

func TestConfigBuilder_GivenProductionMode_WhenBuilding_ThenConfiguresSecureDefaults(t *testing.T) {
	config := factory.NewConfigBuilder().
		EnableProductionMode().
		Build()

	assert.Equal(t, "aes", config.Algorithm)
	assert.Equal(t, 32, config.KeySize)
	assert.True(t, config.AutoGenerateKeys)
	assert.True(t, config.Features.EnableAESEncryption)
	assert.False(t, config.Features.EnableNoOpEncryption)
	assert.True(t, config.Features.EnableKeyRotation)
	assert.True(t, config.Features.EnablePurposeKeys)
}

func TestConfigBuilder_GivenCustomFeatures_WhenBuilding_ThenAppliesFeatures(t *testing.T) {
	customFeatures := factory.FeatureFlags{
		EnableAESEncryption:   false,
		EnableNoOpEncryption:  true,
		EnableKeyRotation:     false,
		EnablePurposeKeys:     false,
		EnableBatchOperations: false,
	}

	config := factory.NewConfigBuilder().
		WithFeatures(customFeatures).
		Build()

	assert.Equal(t, customFeatures, config.Features)
}

func TestBuild_GivenBatchOperationsDisabled_WhenBuilding_ThenServiceStillSupportsBatch(t *testing.T) {
	config := factory.Config{
		Algorithm:        "noop",
		Features: factory.FeatureFlags{
			EnableBatchOperations: false, // This flag is informational
		},
	}

	fact := factory.NewFactory(config)
	service, err := fact.Build()

	assert.NoError(t, err)
	assert.NotNil(t, service)

	// Even with batch operations "disabled", the service should still support it
	// (the flag is more about configuration/optimization hints)
	ctx := context.Background()
	data := map[string]string{
		"field1": "value1",
		"field2": "value2",
	}

	encrypted, err := service.EncryptBatch(ctx, data, "test")
	assert.NoError(t, err)
	assert.Equal(t, data, encrypted) // no-op service returns data as-is

	decrypted, err := service.DecryptBatch(ctx, encrypted, "test")
	assert.NoError(t, err)
	assert.Equal(t, data, decrypted)
}

func TestBuild_GivenManualKeyConfiguration_WhenBuilding_ThenUsesProvidedKeys(t *testing.T) {
	defaultKey := make([]byte, 32)
	copy(defaultKey, []byte("defaultkey1234567890123456789012"))

	purposeKey := make([]byte, 32)
	copy(purposeKey, []byte("purposekey1234567890123456789012"))

	config := factory.Config{
		Algorithm:        "aes",
		KeySize:          32,
		DefaultKey:       defaultKey,
		AutoGenerateKeys: false,
		PurposeKeys: map[string][]byte{
			encryption.PurposeUserEmail: purposeKey,
		},
		Features: factory.DefaultFeatureFlags(),
	}

	fact := factory.NewFactory(config)
	service, err := fact.Build()

	assert.NoError(t, err)
	assert.NotNil(t, service)

	// Test that the service uses the provided keys by encrypting/decrypting
	ctx := context.Background()
	plaintext := "test data"

	// Test default key
	encrypted, err := service.Encrypt(ctx, plaintext)
	assert.NoError(t, err)
	assert.NotEqual(t, plaintext, encrypted)

	decrypted, err := service.Decrypt(ctx, encrypted)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)

	// Test purpose-specific key
	purposeEncrypted, err := service.EncryptWithPurpose(ctx, plaintext, encryption.PurposeUserEmail)
	assert.NoError(t, err)
	assert.NotEqual(t, plaintext, purposeEncrypted)
	assert.NotEqual(t, encrypted, purposeEncrypted) // Should be different from default encryption

	purposeDecrypted, err := service.DecryptWithPurpose(ctx, purposeEncrypted, encryption.PurposeUserEmail)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, purposeDecrypted)
}

func TestBuild_GivenKeyGenerationFailure_WhenBuilding_ThenReturnsError(t *testing.T) {
	config := factory.Config{
		Algorithm:        "aes",
		KeySize:          0, // Invalid key size to force error
		AutoGenerateKeys: true,
		Features:         factory.DefaultFeatureFlags(),
	}

	fact := factory.NewFactory(config)
	service, err := fact.Build()

	assert.Error(t, err)
	assert.Nil(t, service)
}