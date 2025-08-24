package aes

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/gentra/decorator-arch-go/internal/encryption"
)

// service implements encryption.Service interface using AES encryption
type service struct {
	purposeKeys map[string][]byte
	defaultKey  []byte
}

// NewService creates a new AES-based encryption service with purpose-specific keys
func NewService(purposeKeys map[string][]byte, defaultKey []byte) (encryption.Service, error) {
	if len(defaultKey) != 32 {
		return nil, fmt.Errorf("default encryption key must be 32 bytes for AES-256")
	}

	// Validate all purpose keys
	for purpose, key := range purposeKeys {
		if len(key) != 32 {
			return nil, fmt.Errorf("encryption key for purpose '%s' must be 32 bytes for AES-256", purpose)
		}
	}

	return &service{
		purposeKeys: purposeKeys,
		defaultKey:  defaultKey,
	}, nil
}

// NewServiceWithDefaults creates a service with default purpose keys for common use cases
func NewServiceWithDefaults() (encryption.Service, error) {
	// Generate random keys for each purpose
	purposeKeys := make(map[string][]byte)
	purposes := []string{
		encryption.PurposeUserEmail,
		encryption.PurposeUserName,
		encryption.PurposeUserPhone,
		encryption.PurposePaymentCard,
		encryption.PurposeDocumentContent,
		encryption.PurposeSecretAPIKey,
	}

	for _, purpose := range purposes {
		key := make([]byte, 32)
		if _, err := rand.Read(key); err != nil {
			return nil, fmt.Errorf("failed to generate key for purpose '%s': %w", purpose, err)
		}
		purposeKeys[purpose] = key
	}

	// Generate default key
	defaultKey := make([]byte, 32)
	if _, err := rand.Read(defaultKey); err != nil {
		return nil, fmt.Errorf("failed to generate default key: %w", err)
	}

	return NewService(purposeKeys, defaultKey)
}

// Encrypt encrypts plaintext using AES-GCM with the default key
func (s *service) Encrypt(ctx context.Context, plaintext string) (string, error) {
	return s.encrypt(plaintext, s.defaultKey)
}

// Decrypt decrypts ciphertext using AES-GCM with the default key
func (s *service) Decrypt(ctx context.Context, ciphertext string) (string, error) {
	return s.decrypt(ciphertext, s.defaultKey)
}

// EncryptWithPurpose encrypts data for a specific purpose using the appropriate key
func (s *service) EncryptWithPurpose(ctx context.Context, plaintext, purpose string) (string, error) {
	key := s.getKeyForPurpose(purpose)
	return s.encrypt(plaintext, key)
}

// DecryptWithPurpose decrypts data for a specific purpose using the appropriate key
func (s *service) DecryptWithPurpose(ctx context.Context, ciphertext, purpose string) (string, error) {
	key := s.getKeyForPurpose(purpose)
	return s.decrypt(ciphertext, key)
}

// EncryptBatch encrypts multiple data items for a specific purpose
func (s *service) EncryptBatch(ctx context.Context, data map[string]string, purpose string) (map[string]string, error) {
	key := s.getKeyForPurpose(purpose)
	result := make(map[string]string)

	for field, plaintext := range data {
		encrypted, err := s.encrypt(plaintext, key)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt field '%s': %w", field, err)
		}
		result[field] = encrypted
	}

	return result, nil
}

// DecryptBatch decrypts multiple data items for a specific purpose
func (s *service) DecryptBatch(ctx context.Context, data map[string]string, purpose string) (map[string]string, error) {
	key := s.getKeyForPurpose(purpose)
	result := make(map[string]string)

	for field, ciphertext := range data {
		decrypted, err := s.decrypt(ciphertext, key)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt field '%s': %w", field, err)
		}
		result[field] = decrypted
	}

	return result, nil
}

// GenerateKey generates a random 32-byte key for AES-256
func (s *service) GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	return key, err
}

// GenerateKeyForPurpose generates a new key for a specific purpose
func (s *service) GenerateKeyForPurpose(purpose string) ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate key for purpose '%s': %w", purpose, err)
	}

	// Store the new key for the purpose
	s.purposeKeys[purpose] = key
	return key, nil
}

// RotateKeys rotates all encryption keys
func (s *service) RotateKeys() error {
	// Rotate default key
	newDefaultKey, err := s.GenerateKey()
	if err != nil {
		return fmt.Errorf("failed to rotate default key: %w", err)
	}
	s.defaultKey = newDefaultKey

	// Rotate all purpose keys
	for purpose := range s.purposeKeys {
		if _, err := s.GenerateKeyForPurpose(purpose); err != nil {
			return fmt.Errorf("failed to rotate key for purpose '%s': %w", purpose, err)
		}
	}

	return nil
}

// RotateKeyForPurpose rotates the key for a specific purpose
func (s *service) RotateKeyForPurpose(purpose string) error {
	_, err := s.GenerateKeyForPurpose(purpose)
	return err
}

// encrypt encrypts a plaintext string using AES-GCM
func (s *service) encrypt(plaintext string, key []byte) (string, error) {
	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the plaintext
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode to base64 for storage
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt decrypts a base64-encoded ciphertext using AES-GCM
func (s *service) decrypt(ciphertext string, key []byte) (string, error) {
	// Decode from base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extract nonce and ciphertext
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, cipherData := data[:nonceSize], data[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// getKeyForPurpose returns the appropriate key for a given purpose
func (s *service) getKeyForPurpose(purpose string) []byte {
	if key, exists := s.purposeKeys[purpose]; exists {
		return key
	}
	// Fall back to default key if purpose not found
	return s.defaultKey
}
