package ui

import (
	"fmt"

	"github.com/fatih/color"
)

// Prompt outputs a yellow prompt
var Prompt = func(prompt string) {
	color.Yellow(fmt.Sprintf("\n%s: ", prompt))
}

// Error outputs a red message error from an error
var Error = func(err error) {
	color.Red(err.Error())
}

// Success outputs a green successful message
var Success = func(message string) {
	color.Green(message)
}

// Info outputs a blue info message
var Info = func(message string) {
	color.Cyan(message)
}
