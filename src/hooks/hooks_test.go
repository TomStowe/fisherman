package hooks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

	if !strings.HasPrefix(string(content), "#!/bin/sh\n# Pre-commit hook disabled by Fisherman\nexit 0") {
		t.Fatalf("Expected pre-commit hook to be disabled.")
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

// TestParseDuration verifies that parseDuration correctly parses shorthand duration strings.
func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
		hasError bool
	}{
		{"1h", time.Hour, false},
		{"2d", 48 * time.Hour, false},
		{"30m", 30 * time.Minute, false},
		{"5", 0, true},                 // Invalid format
		{"2x", 0, true},                // Unsupported unit
		{"", 0, true},                  // Empty string
		{"0d", 0, false},               // Zero duration
		{"10h", 10 * time.Hour, false}, // Larger duration
	}

	for _, test := range tests {
		result, err := parseDuration(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("parseDuration(%s) error = %v, wantErr %v", test.input, err, test.hasError)
			continue
		}
		if result != test.expected {
			t.Errorf("parseDuration(%s) = %v, want %v", test.input, result, test.expected)
		}
	}
}

// TestRemoveTimeoutCheck verifies that removeTimeoutCheck correctly removes the timeout check from the script.
func TestRemoveTimeoutCheck(t *testing.T) {
	inputScript := `#!/bin/sh
# Temporarily disabled by Fisherman
current_date=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
if [ "$current_date" \< "2099-12-31T23:59:59Z" ]; then
    exit 0
fi

# Actual pre-commit script content
echo "Running pre-commit checks"
`
	expectedOutput := `#!/bin/sh

# Actual pre-commit script content
echo "Running pre-commit checks"
`

	result := removeTimeoutCheck(inputScript)
	if result != expectedOutput {
		t.Errorf("removeTimeoutCheck() = %v, want %v", result, expectedOutput)
	}
}

// TestTemporarilyDisableHook verifies that TemporarilyDisableHook injects a timeout check in the hook.
func TestTemporarilyDisableHook(t *testing.T) {
	hooksDir := filepath.Join(os.TempDir(), "test_hooks")
	os.MkdirAll(hooksDir, 0755)
	defer os.RemoveAll(hooksDir)

	hookPath := filepath.Join(hooksDir, "pre-commit")
	originalContent := "echo \"Running pre-commit checks\"\n"
	fullContent := "#!/bin/sh\n" + originalContent
	os.WriteFile(hookPath, []byte(fullContent), 0755)

	// Temporarily disable the hook for 1 hour
	err := TemporarilyDisableHook(hookPath, "1h")
	if err != nil {
		t.Fatalf("Failed to disable hook temporarily: %v", err)
	}

	// Verify the injected timeout check
	data, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("Failed to read hook file: %v", err)
	}
	content := string(data)

	// Check if the timeout check and original content are present
	if !strings.Contains(content, `current_date=$(date -u +"%Y-%m-%dT%H:%M:%SZ")`) ||
		!strings.Contains(content, "exit 0") ||
		!strings.Contains(content, originalContent) {
		t.Errorf("Hook could not be temporarily disabled")
	}
}
