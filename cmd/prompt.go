package cmd

import (
	"fmt"

	"github.com/chzyer/readline"

	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

type terminalReader interface {
	Readline() (string, error)
	Close() error
}

var rl terminalReader

func init() {
	var err error

	rl, err = readline.New(">> ")

	if err != nil {
		ui.Error(fmt.Errorf("Something went wrong when initializing prompt"))

		util.ErrorExit()
	}
}
