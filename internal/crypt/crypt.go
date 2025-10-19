package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jawahars16/jebi/internal/core"
	jio "github.com/jawahars16/jebi/internal/io"
)

type cryptService struct {
	workingDir string
}

func NewService(workingDir string) *cryptService {
	return &cryptService{
		workingDir: workingDir,
	}
}

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
	path := filepath.Join(s.workingDir, fmt.Sprintf(".%s", core.AppName), core.KeyFileName)
	return os.WriteFile(path, []byte(encodedKey), 0600)
}

// Encrypt encrypts plaintext with AES-GCM using the given 32-byte key.
func (s *cryptService) Encrypt(key []byte, plaintext string) (ciphertextB64, nonceB64 string, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", fmt.Errorf("invalid AES key: %w", err)
	}

	nonce := make([]byte, 12) // 96-bit nonce recommended for GCM
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", fmt.Errorf("failed to read nonce: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", fmt.Errorf("failed to create GCM: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext),
		base64.StdEncoding.EncodeToString(nonce), nil
}

// Decrypt decrypts a base64-encoded AES-GCM ciphertext using the given key and nonce.
func (s *cryptService) Decrypt(key []byte, ciphertextB64, nonceB64 string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", fmt.Errorf("invalid ciphertext encoding: %w", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(nonceB64)
	if err != nil {
		return "", fmt.Errorf("invalid nonce encoding: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("invalid AES key: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(plaintext), nil
}

func (s *cryptService) GetKey() ([]byte, error) {
	path := filepath.Join(s.workingDir, fmt.Sprintf(".%s", core.AppName), core.KeyFileName)
	encodedKey, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	rawKey, err := base64.StdEncoding.DecodeString(string(encodedKey))
	if err != nil {
		return nil, fmt.Errorf("invalid key encoding: %w", err)
	}

	if len(rawKey) != core.KeyLengthBytes {
		return nil, fmt.Errorf("invalid key length: expected %d bytes, got %d", core.KeyLengthBytes, len(rawKey))
	}

	return rawKey, nil
}

func (s *cryptService) LoadSecrets(env string) (map[string]string, error) {
	keyData, err := os.ReadFile(filepath.Join(s.workingDir, fmt.Sprintf(".%s", core.AppName), core.KeyFileName))
	if err != nil {
		return nil, fmt.Errorf("missing encryption key: %w", err)
	}

	rawKey, err := base64.StdEncoding.DecodeString(string(keyData))
	if err != nil {
		return nil, fmt.Errorf("invalid key encoding: %w", err)
	}

	file := filepath.Join(s.workingDir, fmt.Sprintf(".%s", core.AppName), env, core.SecretFileName)
	enc, err := jio.ReadJSONFile[map[string]core.Secret](file)
	if err != nil {
		return nil, fmt.Errorf("failed to read secrets: %w", err)
	}

	out := make(map[string]string)
	for k, v := range enc {
		plaintext := v.Value
		if v.Nonce != "" {
			plaintext, err = s.Decrypt(rawKey, v.Value, v.Nonce)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt %s: %w", k, err)
			}
		}
		out[k] = plaintext
	}

	return out, nil
}
