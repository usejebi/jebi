package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/jawahars16/jebi/internal/keystore"
)

type userService struct {
	workingDir string
	keystore   keystore.KeyStore
}

func NewUserService(workingDir string) *userService {
	return &userService{
		workingDir: workingDir,
		keystore:   keystore.NewDefault(workingDir),
	}
}

func NewUserServiceWithKeystore(workingDir string, ks keystore.KeyStore) *userService {
	return &userService{
		workingDir: workingDir,
		keystore:   ks,
	}
}

// Login authenticates with the server and returns a token
func (u *userService) Login(username, password, server string) (string, error) {
	// TODO: Implement actual authentication logic
	// For now, return a placeholder token
	return "placeholder-auth-token", fmt.Errorf("login not implemented yet")
}

// AuthenticateWithBrowser opens a browser for OAuth-style authentication
func (u *userService) AuthenticateWithBrowser(serverURL string) (*AuthResponse, error) {
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
		err := u.saveAuthResponse(authResult)
		if err != nil {
			return nil, fmt.Errorf("failed to save auth response: %w", err)
		}
		return authResult, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("authentication timeout. Please try again")
	}
}

func (u *userService) saveAuthResponse(resp *AuthResponse) error {
	// Save the full auth response
	if err := u.keystore.Set("auth_response", resp); err != nil {
		return fmt.Errorf("failed to save auth response: %w", err)
	}

	// Save individual tokens for easy access
	if err := u.keystore.Set("access_token", resp.Tokens.AccessToken); err != nil {
		return fmt.Errorf("failed to save access token: %w", err)
	}

	if resp.Tokens.RefreshToken != "" {
		if err := u.keystore.Set("refresh_token", resp.Tokens.RefreshToken); err != nil {
			return fmt.Errorf("failed to save refresh token: %w", err)
		}
	}

	// Save user information (without server field)
	user := User{
		ID:          resp.User.ID,
		Email:       resp.User.Email,
		Username:    resp.User.Username,
		DisplayName: resp.User.DisplayName,
	}
	if err := u.keystore.Set("current_user", user); err != nil {
		return fmt.Errorf("failed to save user info: %w", err)
	}

	return nil
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
func (u *userService) startCallbackServer() (int, <-chan *AuthResponse, error) {
	// Find an available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, nil, fmt.Errorf("failed to find available port: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// Create channel for authentication result
	authChan := make(chan *AuthResponse, 1)

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

		var response AuthResponse
		if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			authChan <- nil
			return
		}

		// Send success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})

		// Send result to channel
		authChan <- &response
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

// SaveAuthToken saves the authentication token securely
func (u *userService) SaveAuthToken(token string) error {
	return u.keystore.Set("access_token", token)
}

// LoadAuthToken loads the authentication token securely
func (u *userService) LoadAuthToken() (string, error) {
	var token string
	err := u.keystore.Get("access_token", &token)
	return token, err
}

// SaveCurrentUser saves the current user information securely
func (u *userService) SaveCurrentUser(user User) error {
	return u.keystore.Set("current_user", user)
}

// LoadCurrentUser loads the current user information securely
func (u *userService) LoadCurrentUser() (*User, error) {
	var user User
	err := u.keystore.Get("current_user", &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAuthResponse retrieves the full authentication response
func (u *userService) GetAuthResponse() (*AuthResponse, error) {
	var authResp AuthResponse
	err := u.keystore.Get("auth_response", &authResp)
	if err != nil {
		return nil, err
	}
	return &authResp, nil
}

// GetRefreshToken retrieves the refresh token if available
func (u *userService) GetRefreshToken() (string, error) {
	var token string
	err := u.keystore.Get("refresh_token", &token)
	return token, err
}

// IsAuthenticated checks if the user is currently authenticated
func (u *userService) IsAuthenticated() bool {
	return u.keystore.Exists("access_token") && u.keystore.Exists("current_user")
}

// GetUserInfo returns basic user information if authenticated
func (u *userService) GetUserInfo() (username string, err error) {
	user, err := u.LoadCurrentUser()
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

// RefreshAuthToken refreshes the authentication token using the refresh token
func (u *userService) RefreshAuthToken() error {
	// TODO: Implement token refresh logic
	// This would typically involve:
	// 1. Get refresh token
	// 2. Make API call to refresh endpoint
	// 3. Save new tokens
	return fmt.Errorf("token refresh not implemented yet")
}

// Logout clears the authentication token and user information
func (u *userService) Logout() error {
	// Clear all authentication data
	keys := []string{"access_token", "refresh_token", "current_user", "auth_response"}

	var errors []error
	for _, key := range keys {
		if err := u.keystore.Delete(key); err != nil {
			errors = append(errors, fmt.Errorf("failed to delete %s: %w", key, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("logout partially failed: %v", errors)
	}

	return nil
}
