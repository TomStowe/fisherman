package hooks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnableHook(t *testing.T) {
	hooksDir := filepath.Join(os.TempDir(), "test_hooks")
	os.MkdirAll(hooksDir, 0755)
	defer os.RemoveAll(hooksDir)

	hookPath := filepath.Join(hooksDir, "pre-commit")

	// Create a dummy hook file to enable
	err := os.WriteFile(hookPath, []byte("#!/bin/sh\necho 'Pre-commit hook'"), 0755)
	if err != nil {
		t.Fatalf("Failed to create dummy pre-commit hook: %v", err)
	}

	// Test enabling the hook
	err = EnableHook(hookPath, "")
	if err != nil {
		t.Fatalf("Failed to enable hook: %v", err)
	}

	// Check if the exit 0 line is removed
	content, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("Failed to read pre-commit hook: %v", err)
	}

	if strings.HasPrefix(string(content), "#!/bin/sh\nexit 0") {
		t.Fatalf("Expected pre-commit hook file to not have 'exit 0' at the start.")
	}
}

func TestDisableHook(t *testing.T) {
	hooksDir := filepath.Join(os.TempDir(), "test_hooks")
	os.MkdirAll(hooksDir, 0755)
	defer os.RemoveAll(hooksDir)

	hookPath := filepath.Join(hooksDir, "pre-commit")

	// Create a dummy hook file to disable
	err := os.WriteFile(hookPath, []byte("#!/bin/sh\necho 'Pre-commit hook'"), 0755)
	if err != nil {
		t.Fatalf("Failed to create dummy pre-commit hook: %v", err)
	}

	// Test disabling the hook
	err = DisableHook(hookPath)
	if err != nil {
		t.Fatalf("Failed to disable hook: %v", err)
	}

	// Check if the exit 0 line is added
	content, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("Failed to read pre-commit hook: %v", err)
	}

	if !strings.HasPrefix(string(content), "#!/bin/sh\nexit 0") {
		t.Fatalf("Expected pre-commit hook file to have 'exit 0' at the start.")
	}
}

func TestEnableHookFromFile(t *testing.T) {
	hooksDir := filepath.Join(os.TempDir(), "test_hooks")
	os.MkdirAll(hooksDir, 0755)
	defer os.RemoveAll(hooksDir)

	hookPath := filepath.Join(hooksDir, "pre-commit")
	testHookFile := filepath.Join(hooksDir, "test_hook.sh")

	// Create a test hook file
	err := os.WriteFile(testHookFile, []byte("#!/bin/sh\necho 'This is a test hook'"), 0755)
	if err != nil {
		t.Fatalf("Failed to create test hook file: %v", err)
	}

	// Test enabling with the specified hook file
	err = EnableHook(hookPath, testHookFile)
	if err != nil {
		t.Fatalf("Failed to enable hook from file: %v", err)
	}

	// Check if the content matches the test hook file
	content, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("Failed to read pre-commit hook: %v", err)
	}

	expectedContent, err := os.ReadFile(testHookFile)
	if err != nil {
		t.Fatalf("Failed to read test hook file: %v", err)
	}

	if string(content) != string(expectedContent) {
		t.Fatalf("Expected pre-commit hook file content to match the test hook file.")
	}
}

func TestIsGitRepository(t *testing.T) {
	// Create a temporary directory and a .git directory inside it
	tempDir := t.TempDir()
	os.Mkdir(filepath.Join(tempDir, ".git"), 0755)

	// Change the current working directory to the tempDir
	os.Chdir(tempDir)

	if !IsGitRepository() {
		t.Fatal("Expected to be a Git repository, but it was not.")
	}

	// Remove .git directory to test non-repository case
	os.RemoveAll(filepath.Join(tempDir, ".git"))

	if IsGitRepository() {
		t.Fatal("Expected to not be a Git repository, but it was.")
	}
}
