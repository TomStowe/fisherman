package logger

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	successColour = color.New(color.FgGreen).SprintFunc()
	warningColour = color.New(color.FgHiMagenta).SprintFunc()
	errorColour   = color.New(color.FgRed).SprintFunc()
)

// LogSuccess prints a success message
func LogSuccess(message string) {
	fmt.Println(successColour(message))
}

// LogWarning prints a warning message
func LogWarning(message string) {
	fmt.Println(warningColour(message))
}

// LogError prints an error message
func LogError(message string) {
	fmt.Println(errorColour(message))
}
