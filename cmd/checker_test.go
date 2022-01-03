package cmd

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/antham/doc-hunt/file"

	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

func TestCheck(t *testing.T) {
	createMocks()
	err := file.Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	doc := file.Doc{
		Identifier: "doc_file_to_track.txt",
		Category:   file.DFILE,
	}

	sources := []file.Source{
		file.Source{
			Identifier: "source1.php",
			Category:   file.SFILEREG,
		},
	}

	err = file.CreateConfig(&doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	ui.Success = func(msg string) {
		assert.Equal(t, "No changes found", msg, "Must display a message showing no changes occured")
	}

	util.SuccessExit = func() {
		t.SkipNow()
	}

	os.Args = []string{"", "check"}

	err = RootCmd.Execute()

	if err != nil {
		logrus.Fatal(err)
	}
}
