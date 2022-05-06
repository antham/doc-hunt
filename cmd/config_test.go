package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/antham/doc-hunt/util"
)

func init() {
	util.AppPath = filepath.Clean(util.AppPath + "/../testing-directory/" + uuid.NewV4().String())
}

func createTestDirectory() {
	if err := os.MkdirAll(util.AppPath, 0777); err != nil && !os.IsExist(err) {
		logrus.Fatal(err)
	}
}

func createSubTestDirectory(path string) {
	if err := os.MkdirAll(fmt.Sprintf("%s%c%s", util.AppPath, filepath.Separator, path), 0777); err != nil && !os.IsExist(err) {
		logrus.Fatal(err)
	}
}

func createADocFile() {
	content := []byte("This is a doc file")
	if err := ioutil.WriteFile(util.GetAbsPath("/doc_file_to_track.txt"), content, 0644); err != nil {
		logrus.Fatal(err)
	}
}

func createSourceFiles() {
	for i := 1; i <= 10; i++ {
		createSourceFile([]byte("<?php echo 'A source file';"), fmt.Sprintf("source%d.php", i))
	}
}

func createSourceFile(content []byte, filename string) {
	if err := ioutil.WriteFile(util.GetAbsPath(filename), content, 0644); err != nil {
		logrus.Fatal(err)
	}
}

func createMocks() {
	createTestDirectory()
	createADocFile()
	createSourceFiles()
}

func TestConfig(t *testing.T) {
	os.Args = []string{"", "config"}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	RootCmd.SetOutput(writer)
	err := RootCmd.Execute()

	if err != nil {
		logrus.Fatal(err)
	}

	err = writer.Flush()

	if err != nil {
		logrus.Fatal(err)
	}

	assert.Contains(t, b.String(), "List, add or delete", "Must display config help")
}
