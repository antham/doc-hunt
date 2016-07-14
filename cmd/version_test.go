package cmd

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

func TestVersion(t *testing.T) {

	ui.Info = func(msg string) {
		assert.Equal(t, "v"+version, msg, "Must output version")
	}

	util.SuccessExit = func() {
	}

	os.Args = []string{"", "version"}

	err := RootCmd.Execute()

	if err != nil {
		logrus.Fatal(err)
	}
}
