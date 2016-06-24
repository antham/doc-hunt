package ui

import (
	"fmt"

	"github.com/fatih/color"
)

// Prompt outputs a yellow prompt
func Prompt(prompt string) {
	color.Yellow(fmt.Sprintf("\n%s: ", prompt))
}

// Error outputs a red message error from an error
func Error(err error) {
	color.Red(err.Error())
}

// Success outputs a green successful message
func Success(message string) {
	color.Green(message)
}

// Info outputs a blue info message
func Info(message string) {
	color.Cyan(message)
}
