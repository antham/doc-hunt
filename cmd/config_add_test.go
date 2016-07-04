package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/antham/doc-hunt/file"
)

func TestParseConfigAddArgsWithMissingFileDoc(t *testing.T) {
	_, _, _, _, err := parseConfigAddArgs([]string{})

	assert.EqualError(t, err, "Missing doc identifier", "Must return a missing file doc error")
}

func TestParseConfigAddArgsWithUnexistingFileDoc(t *testing.T) {
	_, _, _, _, err := parseConfigAddArgs([]string{"whatever", "test"})

	assert.EqualError(t, err, "Doc whatever is not a valid existing file, nor a valid URL", "Must return an unexisting doc identifier error")
}

func TestParseConfigAddArgsWithMissingFileSources(t *testing.T) {
	createMocks()

	_, _, _, _, err := parseConfigAddArgs([]string{"doc_file_to_track.txt"})

	assert.EqualError(t, err, "Missing source identifiers", "Must return a missing source identifier error")
}

func TestParseConfigAddArgsWithUnexistingFileSources(t *testing.T) {
	_, _, _, _, err := parseConfigAddArgs([]string{"doc_file_to_track.txt", "whatever"})

	assert.EqualError(t, err, "Source identifier whatever doesn't exist", "Must return a missing source identifier error")
}

func TestParseConfigAddArgsWithFile(t *testing.T) {
	createMocks()

	doc, docCat, folderSources, fileSources, err := parseConfigAddArgs([]string{"doc_file_to_track.txt", "source1.php,source2.php"})

	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, "doc_file_to_track.txt", doc, "Must return doc file path")
	assert.True(t, file.DFILE == docCat, "Must return a file doc category")
	assert.Equal(t, []string{"source1.php", "source2.php"}, fileSources, "Must return sources file path")
	assert.Equal(t, []string{}, folderSources, "Must return empty folder sources")
}

func TestParseConfigAddArgsWithURL(t *testing.T) {
	createMocks()

	doc, docCat, folderSources, fileSources, err := parseConfigAddArgs([]string{"http://google.com", "source1.php,source2.php"})

	assert.Equal(t, "http://google.com", doc, "Must return a doc url")
	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, []string{"source1.php", "source2.php"}, fileSources, "Must return sources file path")
	assert.True(t, file.DURL == docCat, "Must return an URL doc category")
	assert.Equal(t, []string{}, folderSources, "Must return empty folder sources")
}

func TestParseConfigAddArgsWithSourceFolder(t *testing.T) {
	createMocks()
	createSubTestDirectory("test2")

	doc, docCat, folderSources, fileSources, err := parseConfigAddArgs([]string{"doc_file_to_track.txt", "test2"})

	assert.Equal(t, "doc_file_to_track.txt", doc, "Must return a doc file")
	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, []string{}, fileSources, "Must return empty file sources")
	assert.True(t, file.DFILE == docCat, "Must return a file doc category")
	assert.Equal(t, []string{"test2"}, folderSources, "Must return a folder source")
}
