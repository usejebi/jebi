package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"path/filepath"

	"github.com/jawahars16/jebi/internal/core"
	jio "github.com/jawahars16/jebi/internal/io"
	"github.com/jawahars16/jebi/internal/keystore"
)

type cryptService struct {
	workingDir  string
	keyFilePath string
	keystore    keystore.KeyStore
}

func NewService(workingDir string) *cryptService {
	return &cryptService{
		workingDir:  workingDir,
		keyFilePath: filepath.Join(workingDir, fmt.Sprintf(".%s", core.AppName), core.KeyFilePath),
		keystore:    keystore.NewDefault(workingDir),
	}
}

func NewServiceWithKeystore(workingDir string, ks keystore.KeyStore) *cryptService {
	return &cryptService{
		workingDir:  workingDir,
		keyFilePath: filepath.Join(workingDir, fmt.Sprintf(".%s", core.AppName), core.KeyFilePath),
		keystore:    ks,
	}
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

func (s *cryptService) LoadSecrets(env string) (map[string]string, error) {
	rawKey, err := s.LoadKey()
	if err != nil {
		return nil, fmt.Errorf("failed to load encryption key: %w", err)
	}

	file := filepath.Join(s.workingDir, fmt.Sprintf(".%s", core.AppName), core.EnvDirPath, env, core.SecretFileName)
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
