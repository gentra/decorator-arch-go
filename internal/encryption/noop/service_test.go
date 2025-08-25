package noop

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/encryption"
)

func TestNewService_GivenNoParameters_WhenCreating_ThenReturnsService(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "creates noop encryption service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			service := NewService()

			// Then
			assert.NotNil(t, service)
			assert.Implements(t, (*encryption.Service)(nil), service)
		})
	}
}

func TestService_GivenPlaintext_WhenEncrypt_ThenReturnsUnchangedData(t *testing.T) {
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
			name:      "long text",
			plaintext: "This is a very long piece of text that would normally be encrypted but with noop service returns unchanged",
		},
		{
			name:      "unicode text",
			plaintext: "Hello 世界! 🌍 Test émojis and ñoñó characters",
		},
		{
			name:      "special characters",
			plaintext: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			service := NewService()
			ctx := context.Background()

			// When
			result, err := service.Encrypt(ctx, tt.plaintext)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, tt.plaintext, result)
		})
	}
}

func TestService_GivenCiphertext_WhenDecrypt_ThenReturnsUnchangedData(t *testing.T) {
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
			name:       "encoded-looking text",
			ciphertext: "U29tZXRoaW5nIHRoYXQgbG9va3MgbGlrZSBiYXNlNjQ=",
		},
		{
			name:       "random string",
			ciphertext: "random-cipher-text-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			service := NewService()
			ctx := context.Background()

			// When
			result, err := service.Decrypt(ctx, tt.ciphertext)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, tt.ciphertext, result)
		})
	}
}

func TestService_GivenPlaintextAndPurpose_WhenEncryptWithPurpose_ThenReturnsUnchangedData(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
		purpose   string
	}{
		{
			name:      "user email purpose",
			plaintext: "user@example.com",
			purpose:   encryption.PurposeUserEmail,
		},
		{
			name:      "user name purpose",
			plaintext: "John Doe",
			purpose:   encryption.PurposeUserName,
		},
		{
			name:      "custom purpose",
			plaintext: "sensitive data",
			purpose:   "custom-purpose",
		},
		{
			name:      "empty purpose",
			plaintext: "some data",
			purpose:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			service := NewService()
			ctx := context.Background()

			// When
			result, err := service.EncryptWithPurpose(ctx, tt.plaintext, tt.purpose)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, tt.plaintext, result)
		})
	}
}

func TestService_GivenCiphertextAndPurpose_WhenDecryptWithPurpose_ThenReturnsUnchangedData(t *testing.T) {
	tests := []struct {
		name       string
		ciphertext string
		purpose    string
	}{
		{
			name:       "user email purpose",
			ciphertext: "encrypted-email-data",
			purpose:    encryption.PurposeUserEmail,
		},
		{
			name:       "payment card purpose",
			ciphertext: "encrypted-card-data",
			purpose:    encryption.PurposePaymentCard,
		},
		{
			name:       "custom purpose",
			ciphertext: "encrypted-custom-data",
			purpose:    "custom-purpose",
		},
		{
			name:       "empty purpose",
			ciphertext: "some encrypted data",
			purpose:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			service := NewService()
			ctx := context.Background()

			// When
			result, err := service.DecryptWithPurpose(ctx, tt.ciphertext, tt.purpose)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, tt.ciphertext, result)
		})
	}
}

func TestService_GivenDataMap_WhenEncryptBatch_ThenReturnsUnchangedData(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]string
		purpose string
	}{
		{
			name: "multiple fields",
			data: map[string]string{
				"email":     "user@example.com",
				"firstName": "John",
				"lastName":  "Doe",
				"phone":     "555-1234",
			},
			purpose: encryption.PurposeUserEmail,
		},
		{
			name:    "empty map",
			data:    map[string]string{},
			purpose: encryption.PurposeUserName,
		},
		{
			name: "single field",
			data: map[string]string{
				"secret": "top-secret-data",
			},
			purpose: encryption.PurposeSecretAPIKey,
		},
		{
			name: "fields with empty values",
			data: map[string]string{
				"field1": "value1",
				"field2": "",
				"field3": "value3",
			},
			purpose: "custom-purpose",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			service := NewService()
			ctx := context.Background()

			// When
			result, err := service.EncryptBatch(ctx, tt.data, tt.purpose)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, tt.data, result)
		})
	}
}

func TestService_GivenEncryptedDataMap_WhenDecryptBatch_ThenReturnsUnchangedData(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]string
		purpose string
	}{
		{
			name: "multiple encrypted fields",
			data: map[string]string{
				"email":    "encrypted-email",
				"password": "encrypted-password",
				"token":    "encrypted-token",
			},
			purpose: encryption.PurposeSecretAPIKey,
		},
		{
			name:    "empty map",
			data:    map[string]string{},
			purpose: encryption.PurposeUserEmail,
		},
		{
			name: "single encrypted field",
			data: map[string]string{
				"data": "encrypted-data",
			},
			purpose: encryption.PurposeDocumentContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			service := NewService()
			ctx := context.Background()

			// When
			result, err := service.DecryptBatch(ctx, tt.data, tt.purpose)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, tt.data, result)
		})
	}
}

func TestService_GivenKeyGeneration_WhenGenerateKey_ThenReturnsRandomKey(t *testing.T) {
	tests := []struct {
		name         string
		expectedSize int
	}{
		{
			name:         "generates 32-byte key",
			expectedSize: 32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			service := NewService()

			// When
			key, err := service.GenerateKey()

			// Then
			assert.NoError(t, err)
			assert.Len(t, key, tt.expectedSize)
			assert.NotEmpty(t, key)
		})
	}
}

func TestService_GivenPurpose_WhenGenerateKeyForPurpose_ThenReturnsRandomKey(t *testing.T) {
	tests := []struct {
		name    string
		purpose string
	}{
		{
			name:    "generate key for email purpose",
			purpose: encryption.PurposeUserEmail,
		},
		{
			name:    "generate key for custom purpose",
			purpose: "custom-purpose",
		},
		{
			name:    "generate key for empty purpose",
			purpose: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			service := NewService()

			// When
			key, err := service.GenerateKeyForPurpose(tt.purpose)

			// Then
			assert.NoError(t, err)
			assert.Len(t, key, 32)
			assert.NotEmpty(t, key)
		})
	}
}

func TestService_GivenKeys_WhenRotateKeys_ThenSucceedsWithoutAction(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "rotate all keys",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			service := NewService()

			// When
			err := service.RotateKeys()

			// Then
			assert.NoError(t, err)
		})
	}
}

func TestService_GivenPurpose_WhenRotateKeyForPurpose_ThenSucceedsWithoutAction(t *testing.T) {
	tests := []struct {
		name    string
		purpose string
	}{
		{
			name:    "rotate key for existing purpose",
			purpose: encryption.PurposeUserEmail,
		},
		{
			name:    "rotate key for custom purpose",
			purpose: "custom-purpose",
		},
		{
			name:    "rotate key for empty purpose",
			purpose: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			service := NewService()

			// When
			err := service.RotateKeyForPurpose(tt.purpose)

			// Then
			assert.NoError(t, err)
		})
	}
}

func TestService_GivenRoundTripOperations_WhenEncryptAndDecrypt_ThenDataUnchanged(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
		purpose   string
	}{
		{
			name:      "round trip with default encryption",
			plaintext: "test data",
			purpose:   "",
		},
		{
			name:      "round trip with purpose encryption",
			plaintext: "sensitive data",
			purpose:   encryption.PurposeUserEmail,
		},
		{
			name:      "round trip with empty data",
			plaintext: "",
			purpose:   encryption.PurposeUserName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			service := NewService()
			ctx := context.Background()

			// When - Encrypt
			var encrypted string
			var err error
			if tt.purpose == "" {
				encrypted, err = service.Encrypt(ctx, tt.plaintext)
			} else {
				encrypted, err = service.EncryptWithPurpose(ctx, tt.plaintext, tt.purpose)
			}
			assert.NoError(t, err)

			// When - Decrypt
			var decrypted string
			if tt.purpose == "" {
				decrypted, err = service.Decrypt(ctx, encrypted)
			} else {
				decrypted, err = service.DecryptWithPurpose(ctx, encrypted, tt.purpose)
			}

			// Then
			assert.NoError(t, err)
			assert.Equal(t, tt.plaintext, decrypted)
			
			// In noop service, encrypted should equal original plaintext
			assert.Equal(t, tt.plaintext, encrypted)
		})
	}
}