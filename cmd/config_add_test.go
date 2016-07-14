package cmd

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/antham/doc-hunt/file"

	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

func TestAddConfigWithMissingFileDoc(t *testing.T) {
	ui.Error = func(err error) {
		assert.EqualError(t, err, "Missing doc identifier", "Must return a missing file doc error")
	}

	util.ErrorExit = func() {
		t.SkipNow()
	}

	os.Args = []string{"", "config", "add"}

	err := RootCmd.Execute()

	if err != nil {
		logrus.Fatal(err)
	}
}

func TestAddConfigWithMissingFileSources(t *testing.T) {
	ui.Error = func(err error) {
		assert.EqualError(t, err, "Missing source identifiers", "Must return a missing source identifier error")
	}

	util.ErrorExit = func() {
		t.SkipNow()
	}

	os.Args = []string{"", "config", "add", "test"}

	err := RootCmd.Execute()

	if err != nil {
		logrus.Fatal(err)
	}
}

func TestAddConfigWithMoreThanTwoArguments(t *testing.T) {
	ui.Error = func(err error) {
		assert.EqualError(t, err, "No more than 2 arguments expected", "Must return an overflow argument error")
	}

	util.ErrorExit = func() {
		t.SkipNow()
	}

	os.Args = []string{"", "config", "add", "test", "test", "test"}

	err := RootCmd.Execute()

	if err != nil {
		logrus.Fatal(err)
	}
}

func TestAddConfig(t *testing.T) {
	createMocks()
	err := file.Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	ui.Success = func(msg string) {
		assert.Equal(t, "Config added", msg, "Must display a success message")
	}

	util.SuccessExit = func() {
		t.SkipNow()
	}

	os.Args = []string{"", "config", "add", "doc_file_to_track.txt", "source1.php,source2.php"}

	err = RootCmd.Execute()

	if err != nil {
		logrus.Fatal(err)
	}
}

func TestParseConfigAddArgsWithUnexistingFileDoc(t *testing.T) {
	createMocks()

	_, _, err := parseConfigAddArgs([]string{"whatever", "test"})

	assert.EqualError(t, err, "Doc whatever is not a valid existing file, nor a valid existing folder, nor a valid URL", "Must return an unexisting doc identifier error")
}

func TestParseConfigAddArgsWithUnexistingFileSources(t *testing.T) {
	createMocks()

	_, _, err := parseConfigAddArgs([]string{"doc_file_to_track.txt", "whatever"})

	assert.EqualError(t, err, "Source identifier whatever doesn't exist", "Must return a missing source identifier error")
}

func TestParseConfigAddArgsWithFile(t *testing.T) {
	createMocks()

	doc, sources, err := parseConfigAddArgs([]string{"doc_file_to_track.txt", "source1.php,source2.php"})

	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, "doc_file_to_track.txt", doc.Identifier, "Must return doc file path")
	assert.EqualValues(t, file.DFILE, doc.Category, "Must return a file doc category")
	assert.Equal(t, "source1.php", (*sources)[0].Identifier, "Must return sources file path")
	assert.EqualValues(t, file.SFILE, (*sources)[0].Category, "Must return sources file type")
	assert.Equal(t, "source2.php", (*sources)[1].Identifier, "Must return sources file path")
	assert.EqualValues(t, file.SFILE, (*sources)[1].Category, "Must return sources file type")
}

func TestParseConfigAddArgsWithURL(t *testing.T) {
	createMocks()

	doc, sources, err := parseConfigAddArgs([]string{"http://google.com", "source1.php,source2.php"})

	assert.Equal(t, "http://google.com", doc.Identifier, "Must return a doc url")
	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, "source1.php", (*sources)[0].Identifier, "Must return sources file path")
	assert.EqualValues(t, file.SFILE, (*sources)[0].Category, "Must return sources file type")
	assert.Equal(t, "source2.php", (*sources)[1].Identifier, "Must return sources file path")
	assert.EqualValues(t, file.SFILE, (*sources)[1].Category, "Must return sources file type")
	assert.EqualValues(t, file.DURL, doc.Category, "Must return an URL doc category")
}

func TestParseConfigAddArgsWithAFolder(t *testing.T) {
	createMocks()
	createSubTestDirectory("test2")

	doc, sources, err := parseConfigAddArgs([]string{"test2", "source1.php,source2.php"})

	assert.Equal(t, "test2", doc.Identifier, "Must return a doc folder")
	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, "source1.php", (*sources)[0].Identifier, "Must return sources file path")
	assert.EqualValues(t, file.SFILE, (*sources)[0].Category, "Must return sources file type")
	assert.Equal(t, "source2.php", (*sources)[1].Identifier, "Must return sources file path")
	assert.EqualValues(t, file.SFILE, (*sources)[1].Category, "Must return sources file type")
	assert.EqualValues(t, file.DFOLDER, doc.Category, "Must return a folder doc category")
}

func TestParseConfigAddArgsWithSourceFolder(t *testing.T) {
	createMocks()
	createSubTestDirectory("test2")

	doc, sources, err := parseConfigAddArgs([]string{"doc_file_to_track.txt", "test2"})

	assert.Equal(t, "doc_file_to_track.txt", doc.Identifier, "Must return a doc file")
	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, "test2", (*sources)[0].Identifier, "Must return sources path")
	assert.EqualValues(t, file.SFOLDER, (*sources)[0].Category, "Must return sources folder type")
	assert.EqualValues(t, file.DFILE, doc.Category, "Must return a file doc category")
}
