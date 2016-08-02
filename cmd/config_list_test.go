package cmd

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

func TestListConfigWithNoConfig(t *testing.T) {
	createMocks()

	ui.Info = func(msg string) {
		assert.Equal(t, "No config added yet", msg, "Must return a message showing that no config exists yet")
	}

	util.SuccessExit = func() {
		t.SkipNow()
	}

	os.Args = []string{"", "config", "list"}

	err := RootCmd.Execute()

	if err != nil {
		logrus.Fatal(err)
	}
}
