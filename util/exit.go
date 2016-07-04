package util

import (
	"os"
)

// ErrorExit exit with error status code 1
var ErrorExit = func() {
	os.Exit(1)
}

// SuccessExit exit with success status code 0
var SuccessExit = func() {
	os.Exit(0)
}
