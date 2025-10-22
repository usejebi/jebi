package keystore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/zalando/go-keyring"
)

// KeyStore provides secure storage for sensitive data
type KeyStore interface {
	Set(key string, value interface{}) error
	Get(key string, target interface{}) error
	Delete(key string) error
	Exists(key string) bool
}

// keyStore implements KeyStore interface
type keyStore struct {
	serviceName string
	workingDir  string
	useKeyring  bool
}

// Config holds configuration for the keystore
type Config struct {
	ServiceName string // e.g., "jebi-cli"
	WorkingDir  string // fallback directory for disk storage
	UseKeyring  bool   // whether to attempt keyring first
}

// NewKeyStore creates a new KeyStore instance
func NewKeyStore(config Config) KeyStore {
	return &keyStore{
		serviceName: config.ServiceName,
		workingDir:  config.WorkingDir,
		useKeyring:  config.UseKeyring,
	}
}

// Set stores a value securely
func (k *keyStore) Set(key string, value interface{}) error {
	// Serialize value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to serialize value: %w", err)
	}

	// Try keyring first if enabled
	if k.useKeyring {
		if err := k.setKeyring(key, string(data)); err == nil {
			return nil
		}
		// Fall back to disk if keyring fails
	}

	// Fallback to disk storage
	return k.setDisk(key, data)
}

// Get retrieves a value securely
func (k *keyStore) Get(key string, target interface{}) error {
	var data string
	var err error

	// Try keyring first if enabled
	if k.useKeyring {
		data, err = k.getKeyring(key)
		if err == nil {
			return json.Unmarshal([]byte(data), target)
		}
		// Fall back to disk if keyring fails
	}

	// Fallback to disk storage
	diskData, err := k.getDisk(key)
	if err != nil {
		return err
	}

	return json.Unmarshal(diskData, target)
}

// Delete removes a value securely
func (k *keyStore) Delete(key string) error {
	var keyringErr, diskErr error

	// Try to delete from keyring if enabled
	if k.useKeyring {
		keyringErr = k.deleteKeyring(key)
	}

	// Try to delete from disk
	diskErr = k.deleteDisk(key)

	// Return error only if both failed
	if keyringErr != nil && diskErr != nil {
		return fmt.Errorf("failed to delete from keyring: %v, failed to delete from disk: %v", keyringErr, diskErr)
	}

	return nil
}

// Exists checks if a key exists
func (k *keyStore) Exists(key string) bool {
	// Check keyring first if enabled
	if k.useKeyring {
		if _, err := k.getKeyring(key); err == nil {
			return true
		}
	}

	// Check disk storage
	_, err := k.getDisk(key)
	return err == nil
}

// setKeyring stores data in platform keyring
func (k *keyStore) setKeyring(key, value string) error {
	if !k.isKeyringSupportedPlatform() {
		return fmt.Errorf("keyring not supported on platform: %s", runtime.GOOS)
	}

	return keyring.Set(k.serviceName, key, value)
}

// getKeyring retrieves data from platform keyring
func (k *keyStore) getKeyring(key string) (string, error) {
	if !k.isKeyringSupportedPlatform() {
		return "", fmt.Errorf("keyring not supported on platform: %s", runtime.GOOS)
	}

	return keyring.Get(k.serviceName, key)
}

// deleteKeyring removes data from platform keyring
func (k *keyStore) deleteKeyring(key string) error {
	if !k.isKeyringSupportedPlatform() {
		return fmt.Errorf("keyring not supported on platform: %s", runtime.GOOS)
	}

	return keyring.Delete(k.serviceName, key)
}

// setDisk stores data on disk with restricted permissions
func (k *keyStore) setDisk(key string, data []byte) error {
	filePath := k.getDiskPath(key)

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create keystore directory: %w", err)
	}

	// Write file with restricted permissions
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write keystore file: %w", err)
	}

	return nil
}

// getDisk retrieves data from disk
func (k *keyStore) getDisk(key string) ([]byte, error) {
	filePath := k.getDiskPath(key)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("key not found: %s", key)
		}
		return nil, fmt.Errorf("failed to read keystore file: %w", err)
	}
	return data, nil
}

// deleteDisk removes data from disk
func (k *keyStore) deleteDisk(key string) error {
	filePath := k.getDiskPath(key)
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete keystore file: %w", err)
	}
	return nil
}

// getDiskPath returns the file path for disk storage
func (k *keyStore) getDiskPath(key string) string {
	return filepath.Join(k.workingDir, ".jebi", "keystore", fmt.Sprintf("%s.json", key))
}

// isKeyringSupportedPlatform checks if keyring is supported on current platform
func (k *keyStore) isKeyringSupportedPlatform() bool {
	switch runtime.GOOS {
	case "darwin", "linux", "windows":
		return true
	default:
		return false
	}
}
