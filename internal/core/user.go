package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/jawahars16/jebi/internal/io"
)

type userService struct {
	workingDir string
}

func NewUserService(workingDir string) *userService {
	return &userService{
		workingDir: workingDir,
	}
}

func (u *userService) authTokenPath() string {
	return filepath.Join(u.workingDir, fmt.Sprintf(".%s", AppName), "auth.token")
}

func (u *userService) userConfigPath() string {
	return filepath.Join(u.workingDir, fmt.Sprintf(".%s", AppName), "user.config")
}

// Login authenticates with the server and returns a token
func (u *userService) Login(username, password, server string) (string, error) {
	// TODO: Implement actual authentication logic
	// For now, return a placeholder token
	return "placeholder-auth-token", fmt.Errorf("login not implemented yet")
}

// AuthenticateWithBrowser opens a browser for OAuth-style authentication
func (u *userService) AuthenticateWithBrowser(serverURL string) (*AuthResult, error) {
	// Start local callback server
	port, authChan, err := u.startCallbackServer()
	if err != nil {
		return nil, fmt.Errorf("failed to start callback server: %w", err)
	}

	// Build authentication URL with cli_port param
	authURL := fmt.Sprintf("%s?source=cli&callbackURL=http://localhost:%s/auth/callback", serverURL, strconv.Itoa(port))

	// Open browser to authentication page
	if err := u.openBrowser(authURL); err != nil {
		return nil, fmt.Errorf("failed to open browser: %w", err)
	}

	// Wait for authentication result with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	select {
	case authResult := <-authChan:
		if authResult == nil {
			return nil, fmt.Errorf("authentication failed")
		}
		return authResult, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("authentication timeout. Please try again")
	}
}

// openBrowser opens the specified URL in the user's default browser
func (u *userService) openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return exec.Command(cmd, args...).Start()
}

// startCallbackServer starts a local HTTP server to receive the auth callback
func (u *userService) startCallbackServer() (int, <-chan *AuthResult, error) {
	// Find an available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, nil, fmt.Errorf("failed to find available port: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// Create channel for authentication result
	authChan := make(chan *AuthResult, 1)

	// Create HTTP server
	mux := http.NewServeMux()

	// Handle the callback endpoint (now /auth/callback)
	mux.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload struct {
			Tokens AuthResult `json:"tokens"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			authChan <- nil
			return
		}

		// Send success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})

		// Send result to channel
		authChan <- &payload.Tokens
	})

	// Start server in background
	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: mux,
	}

	go func() {
		server.ListenAndServe()
	}()

	// Automatically shutdown server after timeout
	go func() {
		time.Sleep(5 * time.Minute)
		server.Close()
	}()

	return port, authChan, nil
}

// SaveAuthToken saves the authentication token to disk
func (u *userService) SaveAuthToken(token string) error {
	return io.WriteJSONToFile(u.authTokenPath(), map[string]string{"token": token})
}

// LoadAuthToken loads the authentication token from disk
func (u *userService) LoadAuthToken() (string, error) {
	data, err := io.ReadJSONFile[map[string]string](u.authTokenPath())
	if err != nil {
		return "", err
	}
	return data["token"], nil
}

// SaveCurrentUser saves the current user information
func (u *userService) SaveCurrentUser(user User) error {
	return io.WriteJSONToFile(u.userConfigPath(), user)
}

// LoadCurrentUser loads the current user information
func (u *userService) LoadCurrentUser() (*User, error) {
	user, err := io.ReadJSONFile[User](u.userConfigPath())
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Logout clears the authentication token and user information
func (u *userService) Logout() error {
	// TODO: Implement logout logic (clear token, notify server, etc.)
	return fmt.Errorf("logout not implemented yet")
}
