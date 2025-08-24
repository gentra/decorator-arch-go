package encryption

import (
	"context"
)

// Service defines the encryption domain interface - the ONLY interface in this domain
type Service interface {
	// General encryption operations
	Encrypt(ctx context.Context, plaintext string) (string, error)
	Decrypt(ctx context.Context, ciphertext string) (string, error)

	// Purpose-based encryption for different data types across modules
	EncryptWithPurpose(ctx context.Context, plaintext, purpose string) (string, error)
	DecryptWithPurpose(ctx context.Context, ciphertext, purpose string) (string, error)

	// Batch operations for efficiency
	EncryptBatch(ctx context.Context, data map[string]string, purpose string) (map[string]string, error)
	DecryptBatch(ctx context.Context, data map[string]string, purpose string) (map[string]string, error)

	// Key management
	GenerateKey() ([]byte, error)
	GenerateKeyForPurpose(purpose string) ([]byte, error)
	RotateKeys() error
	RotateKeyForPurpose(purpose string) error
}

// Domain types and data structures

// EncryptedData represents encrypted data with metadata
type EncryptedData struct {
	Data      string `json:"data"`
	Algorithm string `json:"algorithm"`
	KeyID     string `json:"key_id,omitempty"`
}

// EncryptionConfig contains configuration for encryption service
type EncryptionConfig struct {
	Algorithm      string            `json:"algorithm"`       // AES-256, etc.
	KeySize        int               `json:"key_size"`        // Key size in bytes
	PurposeKeys    map[string]string `json:"purpose_keys"`    // Key IDs for different purposes
	RotationPeriod string            `json:"rotation_period"` // Key rotation period
	DefaultPurpose string            `json:"default_purpose"` // Default purpose when none specified
}

// EncryptionError represents domain-specific encryption errors
type EncryptionError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func (e EncryptionError) Error() string {
	return e.Message
}

// Common encryption error codes
var (
	ErrInvalidKey       = EncryptionError{Code: "INVALID_KEY", Message: "Invalid encryption key"}
	ErrEncryptionFailed = EncryptionError{Code: "ENCRYPTION_FAILED", Message: "Encryption operation failed"}
	ErrDecryptionFailed = EncryptionError{Code: "DECRYPTION_FAILED", Message: "Decryption operation failed"}
	ErrKeyNotFound      = EncryptionError{Code: "KEY_NOT_FOUND", Message: "Encryption key not found"}
	ErrInvalidData      = EncryptionError{Code: "INVALID_DATA", Message: "Invalid data format"}
)

// Helper methods for EncryptedData
func (e *EncryptedData) IsValid() bool {
	return e.Data != "" && e.Algorithm != ""
}

func (e *EncryptedData) IsEmpty() bool {
	return e.Data == ""
}

// Helper methods for EncryptionConfig
func (c *EncryptionConfig) IsValid() bool {
	return c.Algorithm != "" && c.KeySize > 0
}

// DefaultEncryptionConfig returns default encryption configuration
func DefaultEncryptionConfig() EncryptionConfig {
	return EncryptionConfig{
		Algorithm: "AES-256-GCM",
		KeySize:   32, // 256 bits
		PurposeKeys: map[string]string{
			"user.email":       "user-email-key-v1",
			"user.name":        "user-name-key-v1",
			"user.phone":       "user-phone-key-v1",
			"payment.card":     "payment-card-key-v1",
			"document.content": "document-content-key-v1",
			"secret.api_key":   "secret-apikey-key-v1",
			"default":          "default-key-v1",
		},
		RotationPeriod: "30d",
		DefaultPurpose: "default",
	}
}

// Common purpose constants for encryption
const (
	PurposeUserEmail       = "user.email"
	PurposeUserName        = "user.name"
	PurposeUserPhone       = "user.phone"
	PurposePaymentCard     = "payment.card"
	PurposeDocumentContent = "document.content"
	PurposeSecretAPIKey    = "secret.api_key"
	PurposeDefault         = "default"
)
