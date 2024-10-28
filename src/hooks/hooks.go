package hooks

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// EnableHook enables the pre-commit hook by writing the specified content to the hook file
func EnableHook(hookPath, hookFile string) error {
	var content []byte
	var err error

	if hookFile != "" {
		// Read the content from the specified file
		content, err = os.ReadFile(hookFile)
		if err != nil {
			return err
		}
	} else {
		// Read the current hook content
		content, err = os.ReadFile(hookPath)
		if err != nil {
			return err
		}

		// Remove the exit code if it exists
		contentStr := string(content)
		if strings.HasPrefix(contentStr, "#!/bin/sh\nexit 0") {
			contentStr = strings.Replace(contentStr, "#!/bin/sh\nexit 0", "#!/bin/sh", 1)
			content = []byte(contentStr)
		}
	}

	// Write the content to the pre-commit hook file
	return os.WriteFile(hookPath, content, 0755)
}

// DisableHook disables the pre-commit hook by adding an exit 0 line
func DisableHook(hookPath string) error {
	if _, err := os.Stat(hookPath); err == nil {
		// Read the current hook content
		content, err := os.ReadFile(hookPath)
		if err != nil {
			return err
		}

		// Add exit 0 at the start
		contentStr := string(content)
		if !strings.HasPrefix(contentStr, "#!/bin/sh\nexit 0") {
			contentStr = "#!/bin/sh\nexit 0\n" + contentStr
			return os.WriteFile(hookPath, []byte(contentStr), 0755)
		}
	}
	return errors.New("no pre-commit hook to disable")
}

// IsGitRepository checks if the current directory is a Git repository
func IsGitRepository() bool {
	_, err := os.Stat(".git")
	return !os.IsNotExist(err)
}

// GetGitHooksDir returns the Git hooks directory based on core.hooksPath or defaults to .git/hooks
func GetGitHooksDir() (string, error) {
	cmd := exec.Command("git", "config", "--get", "core.hooksPath")
	output, err := cmd.Output()
	if err != nil {
		// If core.hooksPath is not set, default to .git/hooks
		defaultHooksDir := filepath.Join(".git", "hooks")
		if _, err := os.Stat(defaultHooksDir); os.IsNotExist(err) {
			return "", err
		}
		return defaultHooksDir, nil
	}

	hooksPath := strings.TrimSpace(string(output))
	if hooksPath == "" {
		return "", err
	}

	return hooksPath, nil
}
