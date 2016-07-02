package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	AppPath = filepath.Clean(AppPath + "/../" + "test")
}

func removeTestDirectory() {
	err := os.RemoveAll(AppPath)

	if err != nil {
		logrus.Fatal(err)
	}
}

func createTestDirectory() {
	err := os.Mkdir(AppPath, 0777)

	if err != nil && !os.IsExist(err) {
		logrus.Fatal(err)
	}
}

func createSubTestDirectory(path string) {
	err := os.Mkdir(fmt.Sprintf("%s%c%s", AppPath, filepath.Separator, path), 0777)

	if err != nil && !os.IsExist(err) {
		logrus.Fatal(err)
	}
}

func TestGetAbsPath(t *testing.T) {
	assert.Equal(t, AppPath+"/file", GetAbsPath("file"), "Must return an absolute path")
}

func TestTrimAbsBasePath(t *testing.T) {
	assert.Equal(t, "test", TrimAbsBasePath(AppPath+"/test"), "Must return relative path")
}

func TestGetFolderPathAddTrailingSeparator(t *testing.T) {
	assert.Equal(t, "test/", GetFolderPath("test"), "Must add a trailing seperator")
}

func TestGetFolderPathAddTrailingSeparatorWithMultipleTrailingSeparator(t *testing.T) {
	assert.Equal(t, "test/", GetFolderPath("test////////////"), "Must add a trailing seperator")
}

func TestExtractFolderFiles(t *testing.T) {
	removeTestDirectory()
	createTestDirectory()
	createSubTestDirectory("test1")
	createSubTestDirectory("test1/test2")
	createSubTestDirectory("test1/test2/test3")

	for _, filename := range []string{"file1", "test1/file2", "test1/test2/file3", "test1/test2/test3/file4", "test1/test2/test3/file5"} {
		content := []byte("This is a file")
		err := ioutil.WriteFile(GetAbsPath(filename), content, 0644)

		if err != nil {
			logrus.Fatal(err)
		}
	}

	assert.Equal(t, &[]string{"file1", "test1/file2", "test1/test2/file3", "test1/test2/test3/file4", "test1/test2/test3/file5"}, ExtractFolderFiles("/"), "Must extract all files in folder")
}
