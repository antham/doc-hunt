package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/antham/doc-hunt/ui"
	"github.com/dlclark/regexp2"
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

// ExtractFilesMatchingReg extract files matching given regexp
func ExtractFilesMatchingReg(exp string) (*[]string, error) {
	files := []string{}
	flags := (regexp2.RegexOptions)(regexp2.ExplicitCapture)

	re, perr := regexp2.Compile(exp, flags)

	if perr != nil {
		return &files, fmt.Errorf(perr.Error())
	}

	w := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		pathTrimmed := TrimAbsBasePath(path)

		ok, err := re.MatchString(pathTrimmed)

		if err != nil {
			return err
		}

		if !info.IsDir() && ok {
			files = append(files, pathTrimmed)
		}

		return nil
	}

	err := filepath.Walk(AppPath, w)

	return &files, err
}
