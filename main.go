package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/TomStowe/fisherman/src/hooks"
	"github.com/TomStowe/fisherman/src/logger"
)

func main() {
	// Define CLI flags
	enable := flag.Bool("enable", false, "Enable pre-commit hook")
	disable := flag.Bool("disable", false, "Disable pre-commit hook")
	hookFile := flag.String("file", "", "File containing the pre-commit hook script")

	flag.Parse()

	// Verify if the current directory is a Git repository
	if !hooks.IsGitRepository() {
		logger.LogError("This is not a Git repository")
		os.Exit(1)
	}

	// Get the Git hooks directory
	gitHooksDir, err := hooks.GetGitHooksDir()
	if err != nil {
		logger.LogError(fmt.Sprintf("Error getting Git hooks directory: %v", err))
		os.Exit(1)
	}

	// Define the path to the pre-commit hook and the backup file
	hookPath := filepath.Join(gitHooksDir, "pre-commit")

	// If a file is specified, enable the hook by default
	if *hookFile != "" {
		*enable = true
	}

	// Handle enabling the hook
	if *enable {
		if err := hooks.EnableHook(hookPath, *hookFile); err != nil {
			logger.LogError(fmt.Sprintf("Error enabling pre-commit hook: %v", err))
			os.Exit(1)
		}
		logger.LogSuccess("Pre-commit hook enabled")
	}

	// Handle disabling the hook
	if *disable {
		if err := hooks.DisableHook(hookPath); err != nil {
			logger.LogError(fmt.Sprintf("Error disabling pre-commit hook: %v", err))
			os.Exit(1)
		}
		logger.LogSuccess("Pre-commit hook disabled")
	}

	// If no flags are provided, print usage
	if !*enable && !*disable {
		fmt.Println("Please specify either -enable or -disable flag")
		flag.Usage()
	}
}
