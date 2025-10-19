package crypt

import (
	"crypto/rand"
	"testing"
)

func Test_EncryptDecrypt(t *testing.T) {
	cryptService := NewService("/tmp") // workingDir is not used in this test
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("failed to generate random key: %v", err)
	}

	plaintext := "my super secret value"
	ciphertextB64, nonceB64, err := cryptService.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := cryptService.Decrypt(key, ciphertextB64, nonceB64)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Fatalf("decrypted value mismatch: got %q, want %q", decrypted, plaintext)
	}
}
