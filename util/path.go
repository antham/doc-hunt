package util

import (
	"fmt"
	"path/filepath"
	"os"
	"strings"

	"github.com/antham/doc-hunt/ui"
)

// AppPath define path where app stands
var AppPath string

func init() {
	var err error

	AppPath, err = os.Getwd()

	if err != nil {
		ui.Error(err)

		ErrorExit()
	}
}

// GetAbsPath retrieve absolute path from relative file path
func GetAbsPath(relPath string) string {
	return fmt.Sprintf("%s/%s", AppPath, relPath)
}

// TrimAbsBasePath remove project folder
func TrimAbsBasePath(absPath string) string {
	return strings.TrimPrefix(absPath, fmt.Sprintf("%s%c", AppPath, filepath.Separator))
}

