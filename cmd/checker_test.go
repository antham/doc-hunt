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

func TestCheck(t *testing.T) {
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
		assert.Equal(t, "No changes found", msg, "Must display a message showing no changes occured")
	}

	util.SuccessExit = func() {
		t.SkipNow()
	}

	os.Args = []string{"", "check"}

	RootCmd.Execute()
}
