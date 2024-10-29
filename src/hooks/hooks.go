package hooks

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/TomStowe/fisherman/src/logger"
)

const (
	permenantlyDisabledMarker = "# Pre-commit hook disabled by Fisherman\nexit 0\n"
	temporaryDisabledComment  = "# Temporarily disabled by Fisherman"
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

		stringContent := string(content)

		// Remove any permenantly disabled markers
		stringContent = removePermenantDisabledCheck(stringContent)

		// Remove any temporary disabled markers
		stringContent = removeTimeoutCheck(stringContent)
		content = []byte(stringContent)
	}

	// Write the content to the pre-commit hook file
	if err := os.WriteFile(hookPath, content, 0755); err != nil {
		return err
	}

	logger.LogSuccess("Pre-commit hook enabled")
	return nil
}

// DisableHook disables the pre-commit hook by adding an exit 0 line
func DisableHook(hookPath string) error {
	if _, err := os.Stat(hookPath); err == nil {
		// Read the current hook content
		content, err := os.ReadFile(hookPath)
		if err != nil {
			return err
		}

		// Convert content to a string for manipulation
		contentStr := string(content)

		// Check if the hook is already disabled
		if strings.Contains(contentStr, permenantlyDisabledMarker) {
			logger.LogWarning("Pre-commit hook is already disabled")
			return nil
		}

		// Remove the temporary disabled option if needed
		contentStr = removeTimeoutCheck(contentStr)

		// Check if the shebang exists
		if strings.HasPrefix(contentStr, "#!/bin/sh") {
			// Insert exit 0 and a comment after the shebang
			contentStr = strings.Replace(contentStr, "#!/bin/sh\n", "#!/bin/sh\n"+permenantlyDisabledMarker, 1)

			// Write the modified content back to the hook file
			err = os.WriteFile(hookPath, []byte(contentStr), 0755)
			if err != nil {
				return fmt.Errorf("failed to write hook file: %v", err)
			}

			// Log the exact time until which the hook is disabled
			logger.LogSuccess("Pre-commit hook disabled")
			return nil
		}

		return errors.New("no valid shebang found to disable the pre-commit hook")
	}
	return errors.New("no pre-commit hook to disable")
}

// TemporarilyDisableHook injects a timeout check at the beginning of the hook file
func TemporarilyDisableHook(hookPath string, duration string) error {
	// Parse the duration and calculate the timeout end time
	dur, err := parseDuration(duration)
	if err != nil {
		return err
	}
	endTime := time.Now().Add(dur)
	timeout := endTime.Format(time.RFC3339) // Using RFC3339 for standard formatting

	// Read the existing hook content
	contentBytes, err := os.ReadFile(hookPath)
	if err != nil {
		return fmt.Errorf("failed to read hook file: %v", err)
	}
	contentStr := string(contentBytes)

	// Remove the permenant and temporary disabled if present
	contentStr = removePermenantDisabledCheck(contentStr)
	contentStr = removeTimeoutCheck(contentStr)

	// Insert after the shebang
	if strings.HasPrefix(contentStr, "#!/bin/sh") {
		timeoutCheck := fmt.Sprintf(`#!/bin/sh
%s
current_date=$(date -u +"%%Y-%%m-%%dT%%H:%%M:%%SZ")
if [ "$current_date" \< "%s" ]; then
	exit 0
fi
`, temporaryDisabledComment, timeout)

		contentStr = strings.Replace(contentStr, "#!/bin/sh\n", timeoutCheck, 1)
	}

	// Write the modified content back to the hook file
	err = os.WriteFile(hookPath, []byte(contentStr), 0755)
	if err != nil {
		return fmt.Errorf("failed to write hook file with timeout check: %v", err)
	}

	// Log the exact time until which the hook is disabled
	logger.LogSuccess(fmt.Sprintf("Pre-commit hook disabled until %s", endTime.Format("2006-01-02 15:04:05 UTC")))

	return nil
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

// removeTimeoutCheck removes the timeout-checking block from the script content.
func removeTimeoutCheck(content string) string {
	lines := strings.Split(content, "\n")
	startIndex, endIndex := -1, -1

	// Find the start and end of the timeout-check block
	for i, line := range lines {
		if strings.Contains(line, temporaryDisabledComment) {
			startIndex = i
		}
		if startIndex != -1 && line == "fi" {
			endIndex = i + 1 // Include the "fi" line
			break
		}
	}

	// If we found a timeout block, remove it
	if startIndex != -1 && endIndex != -1 {
		lines = append(lines[:startIndex], lines[endIndex:]...)
	}

	return strings.Join(lines, "\n")
}

func removePermenantDisabledCheck(content string) string {
	return strings.Replace(content, permenantlyDisabledMarker, "", -1)
}

// parseDuration parses a shorthand duration string (e.g., "1h" for 1 hour, "2d" for 2 days, "30m" for 30 minutes).
func parseDuration(duration string) (time.Duration, error) {
	if len(duration) < 2 {
		return 0, errors.New("invalid duration format: too short, expected format like '1h' or '2d'")
	}

	// Separate the numeric part and the unit part
	numericPart := duration[:len(duration)-1]
	unit := duration[len(duration)-1]

	// Convert the numeric part to an integer
	amount, err := strconv.Atoi(numericPart)
	if err != nil {
		return 0, errors.New("invalid duration format: unable to parse numeric part")
	}

	// Determine the time.Duration based on the unit
	switch unit {
	case 'w':
		return time.Duration(amount) * 7 * 24 * time.Hour, nil
	case 'd':
		return time.Duration(amount) * 24 * time.Hour, nil
	case 'h':
		return time.Duration(amount) * time.Hour, nil
	case 'm':
		return time.Duration(amount) * time.Minute, nil
	case 's':
		return time.Duration(amount) * time.Second, nil
	default:
		return 0, errors.New("invalid duration format: must end with 'h' (hours), 'd' (days), or 'm' (minutes)")
	}
}
