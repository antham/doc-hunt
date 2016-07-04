package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/antham/doc-hunt/util"
)

func TestParseConfigAddArgsWithMissingFileDoc(t *testing.T) {
	_, _, _, _, err := parseConfigAddArgs([]string{})

	assert.EqualError(t, err, "Missing file doc", "Must return a missing file doc error")
}

func TestParseConfigAddArgsWithUnexistingFileDoc(t *testing.T) {
	_, _, _, _, err := parseConfigAddArgs([]string{"whatever"})

	assert.EqualError(t, err, "Doc "+util.GetAbsPath("whatever")+" is not a valid existing file, nor a valid URL", "Must return an unexisting file doc error")
}

func TestParseConfigAddArgsWithMissingFileSources(t *testing.T) {
	createTestDirectory()
	createDocFile()

	_, _, _, _, err := parseConfigAddArgs([]string{"doc_test_1"})

	assert.EqualError(t, err, "Missing file/folder sources", "Must return a missing file sources error")
}

func TestParseConfigAddArgsWithUnexistingFileSources(t *testing.T) {
	_, _, _, _, err := parseConfigAddArgs([]string{"doc_test_1", "whatever"})

	assert.EqualError(t, err, "File/folder source whatever doesn't exist", "Must return a unexisting source file error")
}

func TestParseConfigAddArgsWithFile(t *testing.T) {
	createTestDirectory()
	createDocFile()
	createSourceFiles()

	doc, docCat, folderSources, fileSources, err := parseConfigAddArgs([]string{"doc_test_1", "source_test_1,source_test_2"})

	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, "doc_test_1", doc, "Must return doc file path")
	assert.True(t, 0 == docCat, "Must return a file doc category")
	assert.Equal(t, []string{"source_test_1", "source_test_2"}, fileSources, "Must return sources file path")
	assert.Equal(t, []string{}, folderSources, "Must return empty folder sources")
}

func TestParseConfigAddArgsWithURL(t *testing.T) {
	createTestDirectory()
	createDocFile()
	createSourceFiles()

	doc, docCat, folderSources, fileSources, err := parseConfigAddArgs([]string{"http://google.com", "source_test_1,source_test_2"})

	assert.Equal(t, "http://google.com", doc, "Must return a doc url")
	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, []string{"source_test_1", "source_test_2"}, fileSources, "Must return sources file path")
	assert.True(t, 1 == docCat, "Must return an URL doc category")
	assert.Equal(t, []string{}, folderSources, "Must return empty folder sources")
}
