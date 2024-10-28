package logger

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen).SprintFunc()
	errorColor   = color.New(color.FgRed).SprintFunc()
)

// LogSuccess prints a success message
func LogSuccess(message string) {
	fmt.Println(successColor(message))
}

// LogError prints an error message
func LogError(message string) {
	fmt.Println(errorColor(message))
}
