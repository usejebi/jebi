package keystore

import (
	"os/user"
	"path/filepath"
)

// NewDefault creates a keystore with sensible defaults
func NewDefault(workingDir string) KeyStore {
	// Get current user for keyring username
	currentUser, err := user.Current()
	username := "jebi-user"
	if err == nil && currentUser.Username != "" {
		username = currentUser.Username
	}

	return NewKeyStore(Config{
		ServiceName: "jebi-cli",
		Username:    username,
		WorkingDir:  workingDir,
		UseKeyring:  true, // Enable keyring by default
	})
}

// NewDiskOnly creates a keystore that only uses disk storage
func NewDiskOnly(workingDir string) KeyStore {
	return NewKeyStore(Config{
		ServiceName: "jebi-cli",
		Username:    "local",
		WorkingDir:  workingDir,
		UseKeyring:  false, // Disable keyring
	})
}

// GetKeystoreDir returns the directory used for disk storage
func GetKeystoreDir(workingDir string) string {
	return filepath.Join(workingDir, ".jebi", "keystore")
}
