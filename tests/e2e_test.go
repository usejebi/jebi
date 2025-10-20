package cmd

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// findRepoRoot walks upward until it finds go.mod — ensures path independence.
func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("go.mod not found in any parent directory")
		}
		dir = parent
	}
}

// buildBinary builds a **fresh binary** for each test run.
func buildBinary(t *testing.T) string {
	t.Helper()
	repoRoot := findRepoRoot(t)

	// Create a unique filename for this test (safe for parallel runs).
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "jebi-"+t.Name())

	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = repoRoot
	cmd.Env = os.Environ()

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build binary: %v\nOutput:\n%s", err, out)
	}

	// Ensure the binary is executable (macOS/Linux safety)
	if err := os.Chmod(binPath, 0755); err != nil {
		t.Fatalf("failed to mark binary executable: %v", err)
	}

	// remove .jebi folder from tmpDir in case it exists from previous runs
	jebiDir := filepath.Join(tmpDir, ".jebi")
	if _, err := os.Stat(jebiDir); err == nil {
		if err := os.RemoveAll(jebiDir); err != nil {
			t.Fatalf("failed to remove existing .jebi directory: %v", err)
		}
	}

	return binPath
}

// runCLI runs the compiled CLI binary in a given working directory.
func runCLI(ctx context.Context, t *testing.T, binPath, workDir string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), "NO_COLOR=1")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	return out.String(), err
}

func TestHappyPath(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	bin := buildBinary(t) // fresh binary per test

	// Step 1: init
	out, err := runCLI(ctx, t, bin, tmpDir,
		"init",
		"-n", "TestProject",
		"-d", "E2E project for testing",
		"-e", "dev",
	)
	assert.NoError(t, err, out)
	assert.Contains(t, out, "initialized successfully", out)

	// Step 2: add secret
	out, err = runCLI(ctx, t, bin, tmpDir, "add", "API_KEY", "super-secret-value")
	assert.NoError(t, err, out)
	assert.Contains(t, out, "'API_KEY' added successfully", out)

	// Step 3: commit
	out, err = runCLI(ctx, t, bin, tmpDir, "commit", "-m", "Add API key")
	assert.NoError(t, err, out)
	assert.Contains(t, out, "Commit created for environment 'dev'", out)

	// Step 4: status
	out, err = runCLI(ctx, t, bin, tmpDir, "status")
	assert.NoError(t, err, out)
	assert.Contains(t, out, "No pending changes", out)

	// Step 5: check .jebi folder exists
	appDir := filepath.Join(tmpDir, ".jebi")
	_, err = os.Stat(appDir)
	assert.NoError(t, err)
	assert.False(t, os.IsNotExist(err))

	t.Log("✅ Happy path E2E test passed successfully")
}
