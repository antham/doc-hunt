package util

import (
	"os"
)

// ErrorExit exit with error status code 1
func ErrorExit() {
	os.Exit(1)
}

// SuccessExit exit with success status code 0
func SuccessExit() {
	os.Exit(0)
}
