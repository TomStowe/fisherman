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
	enable := flag.Bool("enable", false, "Enable the pre-commit hook")
	disable := flag.Bool("disable", false, "Disable the pre-commit hook")
	disabledFor := flag.String("disabled-for", "", "Temporarily disable the hook for a duration (e.g. 30s, 5m, 1h, 2d, 10w)")
	file := flag.String("file", "", "File to set as pre-commit hook content")
	flag.Parse()

	hookPath := filepath.Join(".git", "hooks", "pre-commit")
	if *disabledFor != "" {
		err := hooks.TemporarilyDisableHook(hookPath, *disabledFor)
		if err != nil {
			logger.LogError(fmt.Sprintf("Failed to disable hook temporarily: %v", err))
			os.Exit(1)
		}

		return
	}

	if *disable {
		if err := hooks.DisableHook(hookPath); err != nil {
			logger.LogError(fmt.Sprintf("Failed to disable hook: %v", err))
			os.Exit(1)
		}
	} else if *enable || *file != "" {
		if err := hooks.EnableHook(hookPath, *file); err != nil {
			logger.LogError(fmt.Sprintf("Failed to enable hook: %v", err))
			os.Exit(1)
		}
	}
}
