package encryption_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/encryption"
)

func TestEncryptedData_IsValid(t *testing.T) {
	tests := []struct {
		name          string
		encryptedData encryption.EncryptedData
		expected      bool
	}{
		{
			name: "Given encrypted data with data and algorithm, When IsValid is called, Then should return true",
			encryptedData: encryption.EncryptedData{
				Data:      "encrypted_content",
				Algorithm: "AES-256-GCM",
			},
			expected: true,
		},
		{
			name: "Given encrypted data with empty data, When IsValid is called, Then should return false",
			encryptedData: encryption.EncryptedData{
				Data:      "",
				Algorithm: "AES-256-GCM",
			},
			expected: false,
		},
		{
			name: "Given encrypted data with empty algorithm, When IsValid is called, Then should return false",
			encryptedData: encryption.EncryptedData{
				Data:      "encrypted_content",
				Algorithm: "",
			},
			expected: false,
		},
		{
			name: "Given encrypted data with both data and algorithm empty, When IsValid is called, Then should return false",
			encryptedData: encryption.EncryptedData{
				Data:      "",
				Algorithm: "",
			},
			expected: false,
		},
		{
			name: "Given encrypted data with data, algorithm and key ID, When IsValid is called, Then should return true",
			encryptedData: encryption.EncryptedData{
				Data:      "encrypted_content",
				Algorithm: "AES-256-GCM",
				KeyID:     "key-123",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.encryptedData.IsValid()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEncryptedData_IsEmpty(t *testing.T) {
	tests := []struct {
		name          string
		encryptedData encryption.EncryptedData
		expected      bool
	}{
		{
			name: "Given encrypted data with empty data, When IsEmpty is called, Then should return true",
			encryptedData: encryption.EncryptedData{
				Data:      "",
				Algorithm: "AES-256-GCM",
			},
			expected: true,
		},
		{
			name: "Given encrypted data with data, When IsEmpty is called, Then should return false",
			encryptedData: encryption.EncryptedData{
				Data:      "encrypted_content",
				Algorithm: "AES-256-GCM",
			},
			expected: false,
		},
		{
			name: "Given encrypted data with data and empty algorithm, When IsEmpty is called, Then should return false",
			encryptedData: encryption.EncryptedData{
				Data:      "encrypted_content",
				Algorithm: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.encryptedData.IsEmpty()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEncryptionConfig_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		config   encryption.EncryptionConfig
		expected bool
	}{
		{
			name: "Given encryption config with algorithm and key size, When IsValid is called, Then should return true",
			config: encryption.EncryptionConfig{
				Algorithm: "AES-256-GCM",
				KeySize:   32,
			},
			expected: true,
		},
		{
			name: "Given encryption config with empty algorithm, When IsValid is called, Then should return false",
			config: encryption.EncryptionConfig{
				Algorithm: "",
				KeySize:   32,
			},
			expected: false,
		},
		{
			name: "Given encryption config with zero key size, When IsValid is called, Then should return false",
			config: encryption.EncryptionConfig{
				Algorithm: "AES-256-GCM",
				KeySize:   0,
			},
			expected: false,
		},
		{
			name: "Given encryption config with negative key size, When IsValid is called, Then should return false",
			config: encryption.EncryptionConfig{
				Algorithm: "AES-256-GCM",
				KeySize:   -1,
			},
			expected: false,
		},
		{
			name: "Given encryption config with both algorithm empty and zero key size, When IsValid is called, Then should return false",
			config: encryption.EncryptionConfig{
				Algorithm: "",
				KeySize:   0,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.config.IsValid()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultEncryptionConfig(t *testing.T) {
	t.Run("Given default encryption config call, When DefaultEncryptionConfig is called, Then should return valid default configuration", func(t *testing.T) {
		// Act
		config := encryption.DefaultEncryptionConfig()

		// Assert
		assert.Equal(t, "AES-256-GCM", config.Algorithm)
		assert.Equal(t, 32, config.KeySize)
		assert.Equal(t, "30d", config.RotationPeriod)
		assert.Equal(t, "default", config.DefaultPurpose)
		
		// Check purpose keys
		assert.NotNil(t, config.PurposeKeys)
		expectedKeys := map[string]string{
			"user.email":       "user-email-key-v1",
			"user.name":        "user-name-key-v1",
			"user.phone":       "user-phone-key-v1",
			"payment.card":     "payment-card-key-v1",
			"document.content": "document-content-key-v1",
			"secret.api_key":   "secret-apikey-key-v1",
			"default":          "default-key-v1",
		}
		
		for purpose, expectedKeyID := range expectedKeys {
			keyID, exists := config.PurposeKeys[purpose]
			assert.True(t, exists, "Purpose %s should exist", purpose)
			assert.Equal(t, expectedKeyID, keyID, "Purpose %s should have correct key ID", purpose)
		}
		
		// Validate the config
		assert.True(t, config.IsValid())
	})
}

func TestEncryptionError_Error(t *testing.T) {
	tests := []struct {
		name     string
		encErr   encryption.EncryptionError
		expected string
	}{
		{
			name: "Given encryption error with message, When Error is called, Then should return message",
			encErr: encryption.EncryptionError{
				Code:    "TEST_ERROR",
				Message: "Test encryption error",
			},
			expected: "Test encryption error",
		},
		{
			name: "Given encryption error with empty message, When Error is called, Then should return empty string",
			encErr: encryption.EncryptionError{
				Code:    "TEST_ERROR",
				Message: "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.encErr.Error()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEncryptionErrors_Constants(t *testing.T) {
	tests := []struct {
		name         string
		err          encryption.EncryptionError
		expectedCode string
	}{
		{
			name:         "Given ErrInvalidKey, When accessing code, Then should have correct code",
			err:          encryption.ErrInvalidKey,
			expectedCode: "INVALID_KEY",
		},
		{
			name:         "Given ErrEncryptionFailed, When accessing code, Then should have correct code",
			err:          encryption.ErrEncryptionFailed,
			expectedCode: "ENCRYPTION_FAILED",
		},
		{
			name:         "Given ErrDecryptionFailed, When accessing code, Then should have correct code",
			err:          encryption.ErrDecryptionFailed,
			expectedCode: "DECRYPTION_FAILED",
		},
		{
			name:         "Given ErrKeyNotFound, When accessing code, Then should have correct code",
			err:          encryption.ErrKeyNotFound,
			expectedCode: "KEY_NOT_FOUND",
		},
		{
			name:         "Given ErrInvalidData, When accessing code, Then should have correct code",
			err:          encryption.ErrInvalidData,
			expectedCode: "INVALID_DATA",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assert
			assert.Equal(t, tt.expectedCode, tt.err.Code)
			assert.NotEmpty(t, tt.err.Message)
		})
	}
}

func TestPurposeConstants(t *testing.T) {
	tests := []struct {
		name         string
		purpose      string
		expectedStr  string
	}{
		{
			name:         "Given PurposeUserEmail constant, When accessing string value, Then should have correct value",
			purpose:      encryption.PurposeUserEmail,
			expectedStr:  "user.email",
		},
		{
			name:         "Given PurposeUserName constant, When accessing string value, Then should have correct value",
			purpose:      encryption.PurposeUserName,
			expectedStr:  "user.name",
		},
		{
			name:         "Given PurposeUserPhone constant, When accessing string value, Then should have correct value",
			purpose:      encryption.PurposeUserPhone,
			expectedStr:  "user.phone",
		},
		{
			name:         "Given PurposePaymentCard constant, When accessing string value, Then should have correct value",
			purpose:      encryption.PurposePaymentCard,
			expectedStr:  "payment.card",
		},
		{
			name:         "Given PurposeDocumentContent constant, When accessing string value, Then should have correct value",
			purpose:      encryption.PurposeDocumentContent,
			expectedStr:  "document.content",
		},
		{
			name:         "Given PurposeSecretAPIKey constant, When accessing string value, Then should have correct value",
			purpose:      encryption.PurposeSecretAPIKey,
			expectedStr:  "secret.api_key",
		},
		{
			name:         "Given PurposeDefault constant, When accessing string value, Then should have correct value",
			purpose:      encryption.PurposeDefault,
			expectedStr:  "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assert
			assert.Equal(t, tt.expectedStr, tt.purpose)
		})
	}
}

func TestEncryptionConfig_CompleteStructure(t *testing.T) {
	t.Run("Given encryption config with all fields, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		purposeKeys := map[string]string{
			"user.email": "email-key-v1",
			"user.phone": "phone-key-v1",
		}
		
		config := encryption.EncryptionConfig{
			Algorithm:      "AES-256-GCM",
			KeySize:        32,
			PurposeKeys:    purposeKeys,
			RotationPeriod: "30d",
			DefaultPurpose: "default",
		}

		// Assert
		assert.Equal(t, "AES-256-GCM", config.Algorithm)
		assert.Equal(t, 32, config.KeySize)
		assert.Equal(t, purposeKeys, config.PurposeKeys)
		assert.Equal(t, "30d", config.RotationPeriod)
		assert.Equal(t, "default", config.DefaultPurpose)
		assert.True(t, config.IsValid())
	})
}

func TestEncryptedData_CompleteStructure(t *testing.T) {
	t.Run("Given encrypted data with all fields, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		encData := encryption.EncryptedData{
			Data:      "encrypted_content_base64",
			Algorithm: "AES-256-GCM",
			KeyID:     "key-v1-123",
		}

		// Assert
		assert.Equal(t, "encrypted_content_base64", encData.Data)
		assert.Equal(t, "AES-256-GCM", encData.Algorithm)
		assert.Equal(t, "key-v1-123", encData.KeyID)
		assert.True(t, encData.IsValid())
		assert.False(t, encData.IsEmpty())
	})
}