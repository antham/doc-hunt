package cmd

import (
	"bytes"
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

func TestAddConfigWithDryRun(t *testing.T) {
	createMocks()
	err := file.Initialize()
	output := []byte{}
	out = bytes.NewBuffer(output)

	if err != nil {
		logrus.Fatal(err)
	}

	util.SuccessExit = func() {
		assert.Regexp(t, `doc_file_to_track.txt`, out, "Must render document")
		assert.Regexp(t, `Files matching regexp "source1.php"`, out, "Must render original regexp")
		assert.Regexp(t, `=> source1.php`, out, "Must render source")
		assert.Regexp(t, `Files matching regexp "source2.php"`, out, "Must render original regexp")
		assert.Regexp(t, `=> source2.php`, out, "Must render source")
	}

	os.Args = []string{"", "config", "add", "-n", "doc_file_to_track.txt", "source1.php,source2.php"}

	err = RootCmd.Execute()

	if err != nil {
		logrus.Fatal(err)
	}
}

func TestAddConfigWithDryRunAndMissingSource(t *testing.T) {
	createMocks()
	err := file.Initialize()
	output := []byte{}
	out = bytes.NewBuffer(output)

	if err != nil {
		logrus.Fatal(err)
	}

	util.SuccessExit = func() {
		assert.Regexp(t, `doc_file_to_track.txt`, out, "Must render document")
		assert.Regexp(t, `Files matching regexp "source1.php"`, out, "Must render original regexp")
		assert.Regexp(t, `=> source1.php`, out, "Must render source")
		assert.Regexp(t, `Files matching regexp "s.php"`, out, "Must render original regexp")
		assert.Regexp(t, `=> No files found`, out, "Must render source")
	}

	os.Args = []string{"", "config", "add", "-n", "doc_file_to_track.txt", "source1.php,s.php"}

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

func TestParseConfigAddArgsWithADocFile(t *testing.T) {
	createMocks()

	doc, sources, err := parseConfigAddArgs([]string{"doc_file_to_track.txt", "source1.php,source2.php"})

	assert.NoError(t, err, "Must return no error")
	assert.Regexp(t, "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", doc.ID, "Must return an id")
	assert.Equal(t, "doc_file_to_track.txt", doc.Identifier, "Must return doc file path")
	assert.EqualValues(t, file.DFILE, doc.Category, "Must return a file doc category")
	assert.Equal(t, "source1.php", (*sources)[0].Identifier, "Must return source file regexp")
	assert.EqualValues(t, file.SFILEREG, (*sources)[0].Category, "Must return regexp file type")
	assert.Equal(t, "source2.php", (*sources)[1].Identifier, "Must return source file regexp")
	assert.EqualValues(t, file.SFILEREG, (*sources)[1].Category, "Must return regexp file type")
}

func TestParseConfigAddArgsWithADocURL(t *testing.T) {
	createMocks()

	doc, sources, err := parseConfigAddArgs([]string{"http://google.com", "source1.php,source2.php"})

	assert.Equal(t, "http://google.com", doc.Identifier, "Must return a doc url")
	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, "source1.php", (*sources)[0].Identifier, "Must return source file regexp")
	assert.EqualValues(t, file.SFILEREG, (*sources)[0].Category, "Must return regexp file type")
	assert.Equal(t, "source2.php", (*sources)[1].Identifier, "Must return source file regexp")
	assert.EqualValues(t, file.SFILEREG, (*sources)[1].Category, "Must return regexp file type")
	assert.EqualValues(t, file.DURL, doc.Category, "Must return an URL doc category")
}

func TestParseConfigAddArgsWithADocFolder(t *testing.T) {
	createMocks()
	createSubTestDirectory("test2")

	doc, sources, err := parseConfigAddArgs([]string{"test2", "source1.php,source2.php"})

	assert.Equal(t, "test2", doc.Identifier, "Must return a doc folder")
	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, "source1.php", (*sources)[0].Identifier, "Must return source file regexp")
	assert.EqualValues(t, file.SFILEREG, (*sources)[0].Category, "Must return regexp file type")
	assert.Equal(t, "source2.php", (*sources)[1].Identifier, "Must return source file regexp")
	assert.EqualValues(t, file.SFILEREG, (*sources)[1].Category, "Must return regexp file type")
	assert.EqualValues(t, file.DFOLDER, doc.Category, "Must return a folder doc category")
}

func TestParseConfigAddArgsWithAFileSourceRegexp(t *testing.T) {
	createMocks()
	createSubTestDirectory("test2")

	doc, sources, err := parseConfigAddArgs([]string{"doc_file_to_track.txt", "test2"})

	assert.Equal(t, "doc_file_to_track.txt", doc.Identifier, "Must return a doc file")
	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, "test2", (*sources)[0].Identifier, "Must return sources path")
	assert.EqualValues(t, file.SFILEREG, (*sources)[0].Category, "Must return sources regexp type")
	assert.EqualValues(t, file.DFILE, doc.Category, "Must return a file doc category")
}
