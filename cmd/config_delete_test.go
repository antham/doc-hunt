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

func TestParseConfigDelArgsWithArgumentNotANumber(t *testing.T) {
	configs := []file.Config{
		file.Config{},
		file.Config{},
		file.Config{},
		file.Config{},
	}

	_, err := parseConfigDelArgs(&configs, "1,2,3,a")

	assert.EqualError(t, err, "a is not a number", "Must return an error")
}

func TestParseConfigDelArgsWithArgumentNotInRange(t *testing.T) {
	configs := []file.Config{
		file.Config{},
		file.Config{},
		file.Config{},
	}

	_, err := parseConfigDelArgs(&configs, "3,4")

	assert.EqualError(t, err, "Value 3 is out of bounds", "Must return an error")
}

func TestParseConfigDelArgs(t *testing.T) {
	configs := []file.Config{
		file.Config{Doc: file.Doc{Identifier: "doc0.txt"}},
		file.Config{Doc: file.Doc{Identifier: "doc1.txt"}},
		file.Config{Doc: file.Doc{Identifier: "doc2.txt"}},
		file.Config{Doc: file.Doc{Identifier: "doc3.txt"}},
		file.Config{Doc: file.Doc{Identifier: "doc4.txt"}},
	}

	expected := &[]file.Config{
		file.Config{Doc: file.Doc{Identifier: "doc3.txt"}},
		file.Config{Doc: file.Doc{Identifier: "doc4.txt"}},
	}

	results, err := parseConfigDelArgs(&configs, "3,4")

	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, expected, results, "Must return configs")
}

func TestPromptConfigToRemove(t *testing.T) {
	rl = &mockTerminalReader{"1,3,0", nil}

	c := []file.Config{
		file.Config{Doc: file.Doc{Identifier: "doc0.txt"}},
		file.Config{Doc: file.Doc{Identifier: "doc1.txt"}},
		file.Config{Doc: file.Doc{Identifier: "doc2.txt"}},
		file.Config{Doc: file.Doc{Identifier: "doc3.txt"}},
		file.Config{Doc: file.Doc{Identifier: "doc4.txt"}},
	}

	configs, err := promptConfigToRemove(&c)

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *configs, 3, "Must return 3 configs")
	assert.Equal(t, "doc1.txt", (*configs)[0].Doc.Identifier, "Must have chosen correct document")
	assert.Equal(t, "doc3.txt", (*configs)[1].Doc.Identifier, "Must have chosen correct document")
	assert.Equal(t, "doc0.txt", (*configs)[2].Doc.Identifier, "Must have chosen correct document")
}

func TestDeleteConfig(t *testing.T) {
	createMocks()
	err := file.Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	doc := file.NewDoc("doc_file_to_track.txt", file.DFILE)

	sources := []file.Source{
		*file.NewSource(doc, "source1.php", file.SFILEREG),
	}

	err = file.Container.GetConfigRepository().CreateFromDocAndSources(doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	rl = &mockTerminalReader{"0", nil}

	ui.Success = func(msg string) {
		assert.Equal(t, "Config deleted", msg, "Must delete config successfully and display a message")
	}

	util.SuccessExit = func() {
		t.SkipNow()
	}

	os.Args = []string{"", "config", "del"}

	err = RootCmd.Execute()

	if err != nil {
		logrus.Fatal(err)
	}
}
