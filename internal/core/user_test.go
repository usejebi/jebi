package core

import (
	"testing"

	"github.com/jawahars16/jebi/internal/keystore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserServiceWithKeystore(t *testing.T) {
	// Create a user service with disk-only keystore for testing
	ks := keystore.NewDiskOnly(t.TempDir())
	userSvc := NewUserServiceWithKeystore(t.TempDir(), ks)

	// Test saving and loading auth token
	testToken := "test-access-token-123"
	err := userSvc.SaveAuthToken(testToken)
	require.NoError(t, err, "Failed to save auth token")

	loadedToken, err := userSvc.LoadAuthToken()
	require.NoError(t, err, "Failed to load auth token")
	assert.Equal(t, testToken, loadedToken, "Loaded token should match saved token")

	// Test saving and loading user
	testUser := User{Email: "testuser"}
	err = userSvc.SaveCurrentUser(testUser)
	require.NoError(t, err, "Failed to save user")

	loadedUser, err := userSvc.LoadCurrentUser()
	require.NoError(t, err, "Failed to load user")
	assert.Equal(t, testUser.Email, loadedUser.Email, "Loaded user email should match saved email")

	// Test authentication status
	assert.True(t, userSvc.IsAuthenticated(), "User should be authenticated")

	// Test getting user info
	email, err := userSvc.GetUserInfo()
	require.NoError(t, err, "Failed to get user info")
	assert.Equal(t, testUser.Email, email, "User info email should match saved email")

	// Test logout
	err = userSvc.Logout()
	require.NoError(t, err, "Failed to logout")

	// User should no longer be authenticated
	assert.False(t, userSvc.IsAuthenticated(), "User should not be authenticated after logout")

	// Loading user info should fail after logout
	_, err = userSvc.GetUserInfo()
	assert.Error(t, err, "Expected error when getting user info after logout")
}

func TestSaveAuthResponse(t *testing.T) {
	// Create a user service with disk-only keystore for testing
	ks := keystore.NewDiskOnly(t.TempDir())
	userSvc := NewUserServiceWithKeystore(t.TempDir(), ks)

	// Create a mock auth response
	authResp := &AuthResponse{
		User: User{
			ID:          "user-123",
			Email:       "johndoe@example.com",
			Username:    "johndoe",
			DisplayName: "John Doe",
		},
		Tokens: Tokens{
			AccessToken:  "access-token-123",
			RefreshToken: "refresh-token-456",
			ExpiresIn:    3600,
		},
	}

	// Save the auth response
	err := userSvc.saveAuthResponse(authResp)
	require.NoError(t, err, "Failed to save auth response")

	user := authResp.User
	tokens := authResp.Tokens

	// Verify access token was saved
	token, err := userSvc.LoadAuthToken()
	require.NoError(t, err, "Failed to load access token")
	assert.Equal(t, tokens.AccessToken, token, "Loaded access token should match saved token")

	// Verify refresh token was saved
	refreshToken, err := userSvc.GetRefreshToken()
	require.NoError(t, err, "Failed to load refresh token")
	assert.Equal(t, tokens.RefreshToken, refreshToken, "Loaded refresh token should match saved token")

	// Verify user was saved
	userInfo, err := userSvc.LoadCurrentUser()
	require.NoError(t, err, "Failed to load user")
	assert.Equal(t, user.Email, userInfo.Email, "Loaded user email should match saved email")
}
