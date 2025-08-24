package noop

import (
	"context"
	"crypto/rand"

	"github.com/gentra/decorator-arch-go/internal/encryption"
)

// service implements encryption.Service interface with no-op operations for testing/development
type service struct{}

// NewService creates a new no-op encryption service
func NewService() encryption.Service {
	return &service{}
}

// Encrypt returns the plaintext as-is (no encryption)
func (s *service) Encrypt(ctx context.Context, plaintext string) (string, error) {
	return plaintext, nil
}

// Decrypt returns the ciphertext as-is (no decryption)
func (s *service) Decrypt(ctx context.Context, ciphertext string) (string, error) {
	return ciphertext, nil
}

// EncryptWithPurpose returns the plaintext as-is (no encryption)
func (s *service) EncryptWithPurpose(ctx context.Context, plaintext, purpose string) (string, error) {
	return plaintext, nil
}

// DecryptWithPurpose returns the ciphertext as-is (no decryption)
func (s *service) DecryptWithPurpose(ctx context.Context, ciphertext, purpose string) (string, error) {
	return ciphertext, nil
}

// EncryptBatch returns all data as-is (no encryption)
func (s *service) EncryptBatch(ctx context.Context, data map[string]string, purpose string) (map[string]string, error) {
	result := make(map[string]string)
	for key, value := range data {
		result[key] = value
	}
	return result, nil
}

// DecryptBatch returns all data as-is (no decryption)
func (s *service) DecryptBatch(ctx context.Context, data map[string]string, purpose string) (map[string]string, error) {
	result := make(map[string]string)
	for key, value := range data {
		result[key] = value
	}
	return result, nil
}

// GenerateKey generates a random 32-byte key
func (s *service) GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	return key, err
}

// GenerateKeyForPurpose generates a key for a specific purpose (no-op)
func (s *service) GenerateKeyForPurpose(purpose string) ([]byte, error) {
	return s.GenerateKey()
}

// RotateKeys is a no-op for the no-op service
func (s *service) RotateKeys() error {
	return nil
}

// RotateKeyForPurpose is a no-op for the no-op service
func (s *service) RotateKeyForPurpose(purpose string) error {
	return nil
}
