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

func TestUpdate(t *testing.T) {
	createMocks()
	err := file.Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	err = file.CreateConfig("doc_file_to_track.txt", file.DFILE, []string{}, []string{"source1.php"})

	if err != nil {
		logrus.Fatal(err)
	}

	ui.Success = func(msg string) {
		assert.Equal(t, "Update succeeded", msg, "Must display a success message")
	}

	util.SuccessExit = func() {
		t.SkipNow()
	}

	os.Args = []string{"", "update"}

	RootCmd.Execute()
}
