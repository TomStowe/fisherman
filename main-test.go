package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestMain checks the main functionality of the Fisherman CLI
func TestMain(t *testing.T) {
	hooksDir := filepath.Join(os.TempDir(), "test_hooks")
	os.MkdirAll(hooksDir, 0755)
	defer os.RemoveAll(hooksDir)

	// Set the environment variable to point to the test hooks directory
	os.Setenv("GIT_DIR", hooksDir)
	os.Setenv("GIT_HOOKS_PATH", hooksDir)

	// Test enabling the hook
	cmd := exec.Command("go", "run", "fisherman.go", "--enable", "--file", "hooks_test.go")
	cmd.Dir = hooksDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to enable pre-commit hook: %v", err)
	}

	// Check if the pre-commit hook was created
	hookPath := filepath.Join(hooksDir, "pre-commit")
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		t.Fatalf("Expected pre-commit hook file to be created, but it was not.")
	}

	// Test disabling the hook
	cmd = exec.Command("go", "run", "fisherman.go", "--disable")
	cmd.Dir = hooksDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to disable pre-commit hook: %v", err)
	}

	// Check if the exit 0 line was added
	content, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("Failed to read pre-commit hook: %v", err)
	}

	if !containsExit0(string(content)) {
		t.Fatalf("Expected pre-commit hook file to have 'exit 0' at the start.")
	}
}

// Helper function to check for 'exit 0' in the hook content
func containsExit0(content string) bool {
	return content == "#!/bin/sh\nexit 0\n" || content == "#!/bin/sh\nexit 0\n# other contents"
}
