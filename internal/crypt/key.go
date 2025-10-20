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
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(s.keyFilePath), 0700); err != nil {
		return fmt.Errorf("failed to create key directory: %w", err)
	}
	return os.WriteFile(s.keyFilePath, []byte(encodedKey), 0600)
}

func (s *cryptService) LoadKey() ([]byte, error) {
	encodedKey, err := s.readFromFile()
	if err != nil {
		return nil, err
	}
	decodedKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key from file: %w", err)
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
