package util

import (
	"fmt"
	"os"
	"path/filepath"
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

// GetFolderPath add trailing separator if it doesn't exist
func GetFolderPath(path string) string {
	return strings.TrimRight(path, fmt.Sprintf("%c", filepath.Separator)) + fmt.Sprintf("%c", filepath.Separator)
}

// ExtractFolderFiles fetch all files from a given relative folder
func ExtractFolderFiles(path string) *[]string {
	files := []string{}

	w := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, TrimAbsBasePath(path))
		}

		return nil
	}

	filepath.Walk(GetAbsPath(path), w)

	return &files
}
