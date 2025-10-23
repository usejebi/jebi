package crypt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jawahars16/jebi/internal/core"
)

// GenerateKey creates a 32-byte random AES key and returns it in base64 form.
func (s *cryptService) GenerateKey() (encoded string, err error) {
	raw := make([]byte, core.KeyLengthBytes) // AES-256 = 32 bytes
	if _, err = rand.Read(raw); err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}
	encoded = base64.StdEncoding.EncodeToString(raw)
	return encoded, nil
}

func (s *cryptService) SaveKey(encodedKey string) error {
	// Try to save to keystore first
	if err := s.keystore.Set("encryption_key", encodedKey); err != nil {
		// Fallback to file storage if keystore fails
		if err := os.MkdirAll(filepath.Dir(s.keyFilePath), 0700); err != nil {
			return fmt.Errorf("failed to create key directory: %w", err)
		}
		if err := os.WriteFile(s.keyFilePath, []byte(encodedKey), 0600); err != nil {
			return fmt.Errorf("failed to save key to both keystore and file: keystore error: %v, file error: %w", err, err)
		}
	}
	return nil
}

func (s *cryptService) LoadKey() ([]byte, error) {
	// Try to load from keystore first
	var encodedKey string
	err := s.keystore.Get("encryption_key", &encodedKey)
	if err != nil {
		// Fallback to file storage if keystore fails
		encodedKey, err = s.readFromFile()
		if err != nil {
			return nil, fmt.Errorf("failed to load key from both keystore and file: %w", err)
		}
	}

	decodedKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}
	return decodedKey, nil
}

func (s *cryptService) readFromFile() (string, error) {
	encodedKey, err := os.ReadFile(s.keyFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read key file: %w", err)
	}
	return string(encodedKey), nil
}
