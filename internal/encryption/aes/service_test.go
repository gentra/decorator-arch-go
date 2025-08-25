package aes_test

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/encryption"
	"github.com/gentra/decorator-arch-go/internal/encryption/aes"
)

func TestNewService_GivenValidKeys_WhenCreating_ThenReturnsService(t *testing.T) {
	tests := []struct {
		name         string
		purposeKeys  map[string][]byte
		defaultKey   []byte
		expectError  bool
		errorMessage string
	}{
		{
			name:        "valid default key and purpose keys",
			purposeKeys: map[string][]byte{encryption.PurposeUserEmail: make([]byte, 32)},
			defaultKey:  make([]byte, 32),
			expectError: false,
		},
		{
			name:         "invalid default key size",
			purposeKeys:  map[string][]byte{},
			defaultKey:   make([]byte, 16),
			expectError:  true,
			errorMessage: "default encryption key must be 32 bytes for AES-256",
		},
		{
			name:         "invalid purpose key size",
			purposeKeys:  map[string][]byte{encryption.PurposeUserEmail: make([]byte, 24)},
			defaultKey:   make([]byte, 32),
			expectError:  true,
			errorMessage: "encryption key for purpose 'user.email' must be 32 bytes for AES-256",
		},
		{
			name:        "empty purpose keys map",
			purposeKeys: map[string][]byte{},
			defaultKey:  make([]byte, 32),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := aes.NewService(tt.purposeKeys, tt.defaultKey)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
			}
		})
	}
}

func TestNewServiceWithDefaults_GivenNoParameters_WhenCreating_ThenGeneratesAllKeys(t *testing.T) {
	service, err := aes.NewServiceWithDefaults()

	assert.NoError(t, err)
	assert.NotNil(t, service)

	// Test that the service can encrypt and decrypt
	ctx := context.Background()
	plaintext := "test data"

	ciphertext, err := service.Encrypt(ctx, plaintext)
	assert.NoError(t, err)
	assert.NotEqual(t, plaintext, ciphertext)

	decrypted, err := service.Decrypt(ctx, ciphertext)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestEncrypt_GivenPlaintext_WhenEncrypting_ThenReturnsEncryptedData(t *testing.T) {
	service, err := createTestService()
	assert.NoError(t, err)

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
			ciphertext, err := service.Encrypt(ctx, tt.plaintext)

			assert.NoError(t, err)
			assert.NotEmpty(t, ciphertext)
			assert.NotEqual(t, tt.plaintext, ciphertext)

			// Verify it's valid base64
			_, err = base64.StdEncoding.DecodeString(ciphertext)
			assert.NoError(t, err)

			// Verify we can decrypt it back
			decrypted, err := service.Decrypt(ctx, ciphertext)
			assert.NoError(t, err)
			assert.Equal(t, tt.plaintext, decrypted)
		})
	}
}

func TestEncryptWithPurpose_GivenPurpose_WhenEncrypting_ThenUsesCorrectKey(t *testing.T) {
	purposeKeys := map[string][]byte{
		encryption.PurposeUserEmail: make([]byte, 32),
		encryption.PurposeUserName:  make([]byte, 32),
	}
	// Generate different keys
	copy(purposeKeys[encryption.PurposeUserEmail], []byte("12345678901234567890123456789012"))
	copy(purposeKeys[encryption.PurposeUserName], []byte("abcdefghijklmnopqrstuvwxyz123456"))

	defaultKey := make([]byte, 32)
	copy(defaultKey, []byte("defaultkey1234567890123456789012"))

	service, err := aes.NewService(purposeKeys, defaultKey)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		plaintext string
		purpose   string
	}{
		{
			name:      "encrypt with email purpose",
			plaintext: "user@example.com",
			purpose:   encryption.PurposeUserEmail,
		},
		{
			name:      "encrypt with name purpose",
			plaintext: "John Doe",
			purpose:   encryption.PurposeUserName,
		},
		{
			name:      "encrypt with unknown purpose uses default key",
			plaintext: "secret data",
			purpose:   "unknown_purpose",
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ciphertext, err := service.EncryptWithPurpose(ctx, tt.plaintext, tt.purpose)

			assert.NoError(t, err)
			assert.NotEmpty(t, ciphertext)
			assert.NotEqual(t, tt.plaintext, ciphertext)

			// Verify we can decrypt it back with the same purpose
			decrypted, err := service.DecryptWithPurpose(ctx, ciphertext, tt.purpose)
			assert.NoError(t, err)
			assert.Equal(t, tt.plaintext, decrypted)
		})
	}
}

func TestEncryptBatch_GivenMultipleFields_WhenEncrypting_ThenEncryptsAll(t *testing.T) {
	service, err := createTestService()
	assert.NoError(t, err)

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
			purpose: encryption.PurposeUserEmail,
		},
		{
			name:    "encrypt empty data",
			data:    map[string]string{},
			purpose: encryption.PurposeUserEmail,
		},
		{
			name: "encrypt single field",
			data: map[string]string{
				"secret": "confidential info",
			},
			purpose: encryption.PurposeSecretAPIKey,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := service.EncryptBatch(ctx, tt.data, tt.purpose)

			assert.NoError(t, err)
			assert.Equal(t, len(tt.data), len(encrypted))

			// Verify all fields are encrypted differently than original
			for field, original := range tt.data {
				encryptedValue, exists := encrypted[field]
				assert.True(t, exists)
				assert.NotEqual(t, original, encryptedValue)
			}

			// Verify we can decrypt all fields back
			decrypted, err := service.DecryptBatch(ctx, encrypted, tt.purpose)
			assert.NoError(t, err)
			assert.Equal(t, tt.data, decrypted)
		})
	}
}

func TestDecrypt_GivenInvalidData_WhenDecrypting_ThenReturnsError(t *testing.T) {
	service, err := createTestService()
	assert.NoError(t, err)

	tests := []struct {
		name       string
		ciphertext string
	}{
		{
			name:       "invalid base64",
			ciphertext: "invalid-base64!@#",
		},
		{
			name:       "empty string",
			ciphertext: "",
		},
		{
			name:       "valid base64 but too short",
			ciphertext: base64.StdEncoding.EncodeToString([]byte("short")),
		},
		{
			name:       "random data",
			ciphertext: base64.StdEncoding.EncodeToString([]byte("this is not encrypted data but long enough to pass length check")),
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Decrypt(ctx, tt.ciphertext)
			assert.Error(t, err)
		})
	}
}

func TestGenerateKey_GivenNoParameters_WhenGenerating_ThenReturns32Bytes(t *testing.T) {
	service, err := createTestService()
	assert.NoError(t, err)

	key, err := service.GenerateKey()

	assert.NoError(t, err)
	assert.Len(t, key, 32)

	// Generate another key and ensure they're different
	key2, err := service.GenerateKey()
	assert.NoError(t, err)
	assert.NotEqual(t, key, key2)
}

func TestGenerateKeyForPurpose_GivenPurpose_WhenGenerating_ThenStoresAndReturnsKey(t *testing.T) {
	service, err := createTestService()
	assert.NoError(t, err)

	purpose := "test_purpose"
	key, err := service.GenerateKeyForPurpose(purpose)

	assert.NoError(t, err)
	assert.Len(t, key, 32)

	// Test that the new key is used for encryption
	ctx := context.Background()
	plaintext := "test data"

	ciphertext, err := service.EncryptWithPurpose(ctx, plaintext, purpose)
	assert.NoError(t, err)

	decrypted, err := service.DecryptWithPurpose(ctx, ciphertext, purpose)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestRotateKeys_GivenExistingService_WhenRotating_ThenUpdatesAllKeys(t *testing.T) {
	purposeKeys := map[string][]byte{
		encryption.PurposeUserEmail: make([]byte, 32),
	}
	defaultKey := make([]byte, 32)

	service, err := aes.NewService(purposeKeys, defaultKey)
	assert.NoError(t, err)

	// Encrypt some data with old keys
	ctx := context.Background()
	plaintext := "test data"
	oldCiphertext, err := service.Encrypt(ctx, plaintext)
	assert.NoError(t, err)

	// Rotate keys
	err = service.RotateKeys()
	assert.NoError(t, err)

	// Old ciphertext should no longer decrypt correctly
	_, err = service.Decrypt(ctx, oldCiphertext)
	assert.Error(t, err)

	// New encryption should work
	newCiphertext, err := service.Encrypt(ctx, plaintext)
	assert.NoError(t, err)
	assert.NotEqual(t, oldCiphertext, newCiphertext)

	decrypted, err := service.Decrypt(ctx, newCiphertext)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestRotateKeyForPurpose_GivenSpecificPurpose_WhenRotating_ThenUpdatesOnlyThatKey(t *testing.T) {
	purposeKeys := map[string][]byte{
		encryption.PurposeUserEmail: make([]byte, 32),
		encryption.PurposeUserName:  make([]byte, 32),
	}
	defaultKey := make([]byte, 32)

	service, err := aes.NewService(purposeKeys, defaultKey)
	assert.NoError(t, err)

	ctx := context.Background()
	plaintext := "test data"

	// Encrypt with both purposes
	emailCiphertext, err := service.EncryptWithPurpose(ctx, plaintext, encryption.PurposeUserEmail)
	assert.NoError(t, err)

	nameCiphertext, err := service.EncryptWithPurpose(ctx, plaintext, encryption.PurposeUserName)
	assert.NoError(t, err)

	// Rotate only email key
	err = service.RotateKeyForPurpose(encryption.PurposeUserEmail)
	assert.NoError(t, err)

	// Email ciphertext should no longer decrypt
	_, err = service.DecryptWithPurpose(ctx, emailCiphertext, encryption.PurposeUserEmail)
	assert.Error(t, err)

	// Name ciphertext should still decrypt
	decrypted, err := service.DecryptWithPurpose(ctx, nameCiphertext, encryption.PurposeUserName)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

// Helper function to create a test service
func createTestService() (encryption.Service, error) {
	purposeKeys := map[string][]byte{
		encryption.PurposeUserEmail: make([]byte, 32),
	}
	copy(purposeKeys[encryption.PurposeUserEmail], []byte("12345678901234567890123456789012"))

	defaultKey := make([]byte, 32)
	copy(defaultKey, []byte("defaultkey1234567890123456789012"))

	return aes.NewService(purposeKeys, defaultKey)
}