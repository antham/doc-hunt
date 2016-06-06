package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func createTestDirectory() {
	os.RemoveAll("/tmp/doc-hunt")

	err := os.Mkdir("/tmp/doc-hunt", 0777)

	if err != nil && !os.IsExist(err) {
		logrus.Fatal(err)
	}
}

func createDocFile() {
	content := []byte("")
	err := ioutil.WriteFile("/tmp/doc-hunt/doc_test_1", content, 0644)

	if err != nil {
		logrus.Fatal(err)
	}
}

func createSourceFiles() {
	content := []byte("")

	for _, name := range []string{"/tmp/doc-hunt/source_test_1", "/tmp/doc-hunt/source_test_2"} {
		err := ioutil.WriteFile(name, content, 0644)

		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func TestParseConfigAddArgsWithMissingFileDoc(t *testing.T) {
	_, _, err := parseConfigAddArgs([]string{})

	assert.EqualError(t, err, "Missing file doc", "Must return a missing file doc error")
}

func TestParseConfigAddArgsWithUnexistingFileDoc(t *testing.T) {
	_, _, err := parseConfigAddArgs([]string{"/tmp/doc-hunt/whatever"})

	assert.EqualError(t, err, "File doc /tmp/doc-hunt/whatever doesn't exist", "Must return an unexisting file doc error")
}

func TestParseConfigAddArgsWithMissingFileSources(t *testing.T) {
	createTestDirectory()
	createDocFile()

	_, _, err := parseConfigAddArgs([]string{"/tmp/doc-hunt/doc_test_1"})

	assert.EqualError(t, err, "Missing file sources", "Must return a missing file sources error")
}

func TestParseConfigAddArgsWithUnexistingFileSources(t *testing.T) {
	_, _, err := parseConfigAddArgs([]string{"/tmp/doc-hunt/doc_test_1", "/tmp/doc-hunt/whatever"})

	assert.EqualError(t, err, "File source /tmp/doc-hunt/whatever doesn't exist", "Must return a unexisting source file error")
}

func TestParseConfigAddArgs(t *testing.T) {
	createTestDirectory()
	createDocFile()
	createSourceFiles()

	docFile, sourceFiles, err := parseConfigAddArgs([]string{"/tmp/doc-hunt/doc_test_1", "/tmp/doc-hunt/source_test_1,/tmp/doc-hunt/source_test_2"})

	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, "/tmp/doc-hunt/doc_test_1", docFile, "Must return doc file path")
	assert.Equal(t, []string{"/tmp/doc-hunt/source_test_1", "/tmp/doc-hunt/source_test_2"}, sourceFiles, "Must return sources file path")
}
