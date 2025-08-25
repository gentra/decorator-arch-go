package noop_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/encryption/noop"
)

func TestNewService_GivenNoParameters_WhenCreating_ThenReturnsService(t *testing.T) {
	service := noop.NewService()

	assert.NotNil(t, service)
}

func TestEncrypt_GivenPlaintext_WhenEncrypting_ThenReturnsPlaintextUnchanged(t *testing.T) {
	service := noop.NewService()

	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "simple text",
			plaintext: "hello world",
		},
		{
			name:      "empty string",
			plaintext: "",
		},
		{
			name:      "special characters",
			plaintext: "!@#$%^&*()_+-={}[]|\\:;\"'<>?,./",
		},
		{
			name:      "unicode text",
			plaintext: "こんにちは世界",
		},
		{
			name:      "long text",
			plaintext: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Encrypt(ctx, tt.plaintext)

			assert.NoError(t, err)
			assert.Equal(t, tt.plaintext, result)
		})
	}
}

func TestDecrypt_GivenCiphertext_WhenDecrypting_ThenReturnsCiphertextUnchanged(t *testing.T) {
	service := noop.NewService()

	tests := []struct {
		name       string
		ciphertext string
	}{
		{
			name:       "simple text",
			ciphertext: "hello world",
		},
		{
			name:       "empty string",
			ciphertext: "",
		},
		{
			name:       "special characters",
			ciphertext: "!@#$%^&*()_+-={}[]|\\:;\"'<>?,./",
		},
		{
			name:       "unicode text",
			ciphertext: "こんにちは世界",
		},
		{
			name:       "long text",
			ciphertext: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Decrypt(ctx, tt.ciphertext)

			assert.NoError(t, err)
			assert.Equal(t, tt.ciphertext, result)
		})
	}
}

func TestEncryptWithPurpose_GivenPlaintextAndPurpose_WhenEncrypting_ThenReturnsPlaintextUnchanged(t *testing.T) {
	service := noop.NewService()

	tests := []struct {
		name      string
		plaintext string
		purpose   string
	}{
		{
			name:      "encrypt with email purpose",
			plaintext: "user@example.com",
			purpose:   "user_email",
		},
		{
			name:      "encrypt with name purpose",
			plaintext: "John Doe",
			purpose:   "user_name",
		},
		{
			name:      "encrypt with empty purpose",
			plaintext: "secret data",
			purpose:   "",
		},
		{
			name:      "encrypt with unknown purpose",
			plaintext: "secret data",
			purpose:   "unknown_purpose",
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.EncryptWithPurpose(ctx, tt.plaintext, tt.purpose)

			assert.NoError(t, err)
			assert.Equal(t, tt.plaintext, result)
		})
	}
}

func TestDecryptWithPurpose_GivenCiphertextAndPurpose_WhenDecrypting_ThenReturnsCiphertextUnchanged(t *testing.T) {
	service := noop.NewService()

	tests := []struct {
		name       string
		ciphertext string
		purpose    string
	}{
		{
			name:       "decrypt with email purpose",
			ciphertext: "user@example.com",
			purpose:    "user_email",
		},
		{
			name:       "decrypt with name purpose",
			ciphertext: "John Doe",
			purpose:    "user_name",
		},
		{
			name:       "decrypt with empty purpose",
			ciphertext: "secret data",
			purpose:    "",
		},
		{
			name:       "decrypt with unknown purpose",
			ciphertext: "secret data",
			purpose:    "unknown_purpose",
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.DecryptWithPurpose(ctx, tt.ciphertext, tt.purpose)

			assert.NoError(t, err)
			assert.Equal(t, tt.ciphertext, result)
		})
	}
}

func TestEncryptBatch_GivenDataMap_WhenEncrypting_ThenReturnsDataUnchanged(t *testing.T) {
	service := noop.NewService()

	tests := []struct {
		name    string
		data    map[string]string
		purpose string
	}{
		{
			name: "encrypt user data",
			data: map[string]string{
				"email":      "user@example.com",
				"first_name": "John",
				"last_name":  "Doe",
			},
			purpose: "user_email",
		},
		{
			name:    "encrypt empty data",
			data:    map[string]string{},
			purpose: "user_email",
		},
		{
			name: "encrypt single field",
			data: map[string]string{
				"secret": "confidential info",
			},
			purpose: "secret_api_key",
		},
		{
			name: "encrypt with nil data map",
			data: nil,
			purpose: "test",
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.EncryptBatch(ctx, tt.data, tt.purpose)

			assert.NoError(t, err)
			if tt.data == nil {
				assert.NotNil(t, result)
				assert.Empty(t, result)
			} else {
				assert.Equal(t, len(tt.data), len(result))
				for key, value := range tt.data {
					assert.Equal(t, value, result[key])
				}
			}
		})
	}
}

func TestDecryptBatch_GivenDataMap_WhenDecrypting_ThenReturnsDataUnchanged(t *testing.T) {
	service := noop.NewService()

	tests := []struct {
		name    string
		data    map[string]string
		purpose string
	}{
		{
			name: "decrypt user data",
			data: map[string]string{
				"email":      "user@example.com",
				"first_name": "John",
				"last_name":  "Doe",
			},
			purpose: "user_email",
		},
		{
			name:    "decrypt empty data",
			data:    map[string]string{},
			purpose: "user_email",
		},
		{
			name: "decrypt single field",
			data: map[string]string{
				"secret": "confidential info",
			},
			purpose: "secret_api_key",
		},
		{
			name: "decrypt with nil data map",
			data: nil,
			purpose: "test",
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.DecryptBatch(ctx, tt.data, tt.purpose)

			assert.NoError(t, err)
			if tt.data == nil {
				assert.NotNil(t, result)
				assert.Empty(t, result)
			} else {
				assert.Equal(t, len(tt.data), len(result))
				for key, value := range tt.data {
					assert.Equal(t, value, result[key])
				}
			}
		})
	}
}

func TestGenerateKey_GivenNoParameters_WhenGenerating_ThenReturns32Bytes(t *testing.T) {
	service := noop.NewService()

	key, err := service.GenerateKey()

	assert.NoError(t, err)
	assert.Len(t, key, 32)

	// Generate another key and ensure they're different (should be random)
	key2, err := service.GenerateKey()
	assert.NoError(t, err)
	assert.NotEqual(t, key, key2)
}

func TestGenerateKeyForPurpose_GivenPurpose_WhenGenerating_ThenReturns32Bytes(t *testing.T) {
	service := noop.NewService()

	tests := []struct {
		name    string
		purpose string
	}{
		{
			name:    "generate key for email purpose",
			purpose: "user_email",
		},
		{
			name:    "generate key for empty purpose",
			purpose: "",
		},
		{
			name:    "generate key for unknown purpose",
			purpose: "unknown_purpose",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := service.GenerateKeyForPurpose(tt.purpose)

			assert.NoError(t, err)
			assert.Len(t, key, 32)
		})
	}
}

func TestRotateKeys_GivenNoParameters_WhenRotating_ThenSucceeds(t *testing.T) {
	service := noop.NewService()

	err := service.RotateKeys()

	assert.NoError(t, err)
}

func TestRotateKeyForPurpose_GivenPurpose_WhenRotating_ThenSucceeds(t *testing.T) {
	service := noop.NewService()

	tests := []struct {
		name    string
		purpose string
	}{
		{
			name:    "rotate key for email purpose",
			purpose: "user_email",
		},
		{
			name:    "rotate key for empty purpose",
			purpose: "",
		},
		{
			name:    "rotate key for unknown purpose",
			purpose: "unknown_purpose",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.RotateKeyForPurpose(tt.purpose)

			assert.NoError(t, err)
		})
	}
}

func TestNoOpService_GivenFullWorkflow_WhenExecuting_ThenMaintainsDataIntegrity(t *testing.T) {
	service := noop.NewService()
	ctx := context.Background()

	// Test complete encrypt/decrypt cycle
	plaintext := "sensitive user data"

	// Test basic encryption/decryption
	encrypted, err := service.Encrypt(ctx, plaintext)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, encrypted)

	decrypted, err := service.Decrypt(ctx, encrypted)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)

	// Test purpose-based encryption/decryption
	purpose := "user_email"
	encryptedWithPurpose, err := service.EncryptWithPurpose(ctx, plaintext, purpose)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, encryptedWithPurpose)

	decryptedWithPurpose, err := service.DecryptWithPurpose(ctx, encryptedWithPurpose, purpose)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decryptedWithPurpose)

	// Test batch operations
	batchData := map[string]string{
		"field1": "value1",
		"field2": "value2",
	}

	encryptedBatch, err := service.EncryptBatch(ctx, batchData, purpose)
	assert.NoError(t, err)
	assert.Equal(t, batchData, encryptedBatch)

	decryptedBatch, err := service.DecryptBatch(ctx, encryptedBatch, purpose)
	assert.NoError(t, err)
	assert.Equal(t, batchData, decryptedBatch)

	// Test key operations
	key, err := service.GenerateKey()
	assert.NoError(t, err)
	assert.Len(t, key, 32)

	keyForPurpose, err := service.GenerateKeyForPurpose(purpose)
	assert.NoError(t, err)
	assert.Len(t, keyForPurpose, 32)

	err = service.RotateKeys()
	assert.NoError(t, err)

	err = service.RotateKeyForPurpose(purpose)
	assert.NoError(t, err)
}