# KeyStore Package

A generic, secure keystore for storing sensitive data with platform-specific keyring support and disk fallback.

## Features

- **Platform Keyring Support**: Uses the system keyring/keychain on supported platforms (macOS, Linux, Windows)
- **Disk Fallback**: Automatically falls back to encrypted disk storage if keyring is unavailable
- **Generic Interface**: Store any JSON-serializable data type
- **Secure Permissions**: Files stored with restricted permissions (0600)
- **Cross-Platform**: Works on all major platforms

## Usage

### Basic Usage

```go
import "github.com/jawahars16/jebi/internal/keystore"

// Create a keystore with default settings (keyring + disk fallback)
ks := keystore.NewDefault("/path/to/working/dir")

// Store a simple value
err := ks.Set("api_key", "secret-api-key-12345")

// Retrieve the value
var apiKey string
err = ks.Get("api_key", &apiKey)

// Check if key exists
exists := ks.Exists("api_key")

// Delete a key
err = ks.Delete("api_key")
```

### Storing Complex Data

```go
type AuthTokens struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
}

tokens := AuthTokens{
    AccessToken:  "access-123",
    RefreshToken: "refresh-456",
    ExpiresIn:    3600,
}

// Store the struct
err := ks.Set("auth_tokens", tokens)

// Retrieve the struct
var retrieved AuthTokens
err = ks.Get("auth_tokens", &retrieved)
```

### Configuration Options

```go
// Use only disk storage (no keyring)
ks := keystore.NewDiskOnly("/path/to/working/dir")

// Custom configuration
ks := keystore.NewKeyStore(keystore.Config{
    ServiceName: "my-app",
    Username:    "my-user",
    WorkingDir:  "/path/to/working/dir",
    UseKeyring:  true,
})
```

## How it Works

1. **Primary Storage**: When `UseKeyring` is true, the keystore first attempts to use the platform's secure keyring:

   - **macOS**: Keychain Access
   - **Linux**: Secret Service API (GNOME Keyring, KDE Wallet)
   - **Windows**: Windows Credential Manager

2. **Fallback Storage**: If keyring is unavailable or fails, data is stored on disk:

   - Location: `{WorkingDir}/.jebi/keystore/{key}.json`
   - Permissions: 0600 (owner read/write only)
   - Format: JSON-encoded data

3. **Data Serialization**: All data is JSON-encoded before storage, allowing complex data types to be stored and retrieved.

## Security Considerations

- **Keyring Security**: When using keyring storage, data security depends on the platform's keyring implementation
- **Disk Security**: Disk storage uses file permissions (0600) but is not encrypted - consider this for highly sensitive data
- **JSON Serialization**: Sensitive data may be temporarily visible in memory during JSON encoding/decoding

## Use Cases in Jebi

This keystore can be used throughout the Jebi CLI for:

- **Authentication Tokens**: Store OAuth tokens, JWT tokens, API keys
- **Encryption Keys**: Store AES encryption keys for secret management
- **User Preferences**: Store user configuration and preferences
- **Session Data**: Store temporary session information

## Future Enhancements

- Disk encryption for fallback storage
- Key rotation and expiration
- Audit logging
- Integration with external secret management services
