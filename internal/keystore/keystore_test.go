package keystore

import (
	"fmt"
	"testing"
)

// Example usage of the keystore
func ExampleKeyStore() {
	// Create a keystore (disk-only for testing)
	ks := NewDiskOnly("/tmp/jebi-keystore-example")

	// Store a simple string
	err := ks.Set("api_key", "secret-api-key-12345")
	if err != nil {
		fmt.Printf("Error storing api_key: %v\n", err)
		return
	}

	// Store a complex struct
	type UserProfile struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Token    string `json:"token"`
	}

	profile := UserProfile{
		Username: "johndoe",
		Email:    "john@example.com",
		Token:    "jwt-token-xyz",
	}

	err = ks.Set("user_profile", profile)
	if err != nil {
		fmt.Printf("Error storing user_profile: %v\n", err)
		return
	}

	// Retrieve the string
	var apiKey string
	err = ks.Get("api_key", &apiKey)
	if err != nil {
		fmt.Printf("Error retrieving api_key: %v\n", err)
		return
	}
	fmt.Printf("Retrieved API Key: %s\n", apiKey)

	// Retrieve the struct
	var retrievedProfile UserProfile
	err = ks.Get("user_profile", &retrievedProfile)
	if err != nil {
		fmt.Printf("Error retrieving user_profile: %v\n", err)
		return
	}
	fmt.Printf("Retrieved Profile: %+v\n", retrievedProfile)

	// Check if key exists
	exists := ks.Exists("api_key")
	fmt.Printf("API Key exists: %v\n", exists)

	// Delete a key
	err = ks.Delete("api_key")
	if err != nil {
		fmt.Printf("Error deleting api_key: %v\n", err)
		return
	}

	// Check if key still exists
	exists = ks.Exists("api_key")
	fmt.Printf("API Key exists after deletion: %v\n", exists)

	// Output:
	// Retrieved API Key: secret-api-key-12345
	// Retrieved Profile: {Username:johndoe Email:john@example.com Token:jwt-token-xyz}
	// API Key exists: true
	// API Key exists after deletion: false
}

// TestKeyStoreOperations demonstrates basic keystore operations
func TestKeyStoreOperations(t *testing.T) {
	// Use disk-only keystore for consistent testing
	ks := NewDiskOnly(t.TempDir())

	// Test Set and Get
	testValue := "test-secret-value"
	err := ks.Set("test_key", testValue)
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	var retrievedValue string
	err = ks.Get("test_key", &retrievedValue)
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}

	if retrievedValue != testValue {
		t.Errorf("Expected %s, got %s", testValue, retrievedValue)
	}

	// Test Exists
	if !ks.Exists("test_key") {
		t.Error("Key should exist")
	}

	// Test Delete
	err = ks.Delete("test_key")
	if err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}

	if ks.Exists("test_key") {
		t.Error("Key should not exist after deletion")
	}

	// Test Get non-existent key
	err = ks.Get("non_existent", &retrievedValue)
	if err == nil {
		t.Error("Expected error when getting non-existent key")
	}
}

// TestKeyStoreWithStruct demonstrates storing complex data types
func TestKeyStoreWithStruct(t *testing.T) {
	ks := NewDiskOnly(t.TempDir())

	type AuthTokens struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
	}

	original := AuthTokens{
		AccessToken:  "access-123",
		RefreshToken: "refresh-456",
		ExpiresIn:    3600,
	}

	// Store the struct
	err := ks.Set("auth_tokens", original)
	if err != nil {
		t.Fatalf("Failed to set auth tokens: %v", err)
	}

	// Retrieve the struct
	var retrieved AuthTokens
	err = ks.Get("auth_tokens", &retrieved)
	if err != nil {
		t.Fatalf("Failed to get auth tokens: %v", err)
	}

	// Verify all fields
	if retrieved.AccessToken != original.AccessToken {
		t.Errorf("AccessToken mismatch: expected %s, got %s", original.AccessToken, retrieved.AccessToken)
	}
	if retrieved.RefreshToken != original.RefreshToken {
		t.Errorf("RefreshToken mismatch: expected %s, got %s", original.RefreshToken, retrieved.RefreshToken)
	}
	if retrieved.ExpiresIn != original.ExpiresIn {
		t.Errorf("ExpiresIn mismatch: expected %d, got %d", original.ExpiresIn, retrieved.ExpiresIn)
	}
}
