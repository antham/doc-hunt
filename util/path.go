package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/antham/doc-hunt/ui"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
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
	flags := pcre.ANCHORED

	re, perr := pcre.Compile(exp, flags)

	if perr != nil {
		return &files, fmt.Errorf(perr.String())
	}

	w := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		pathTrimmed := TrimAbsBasePath(path)

		if !info.IsDir() && re.MatcherString(pathTrimmed, flags).Matches() {
			files = append(files, pathTrimmed)
		}

		return nil
	}

	err := filepath.Walk(AppPath, w)

	return &files, err
}
