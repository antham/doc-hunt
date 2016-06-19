package cmd

import (
	"fmt"

	"github.com/chzyer/readline"
)

type checker func(line string) error

type terminalReader interface {
	Readline() (string, error)
	Close() error
}

var rl terminalReader

func init() {
	var err error

	rl, err = readline.New(">> ")

	if err != nil {
		renderError(fmt.Errorf("Something went wrong when initializing prompt"))

		errorExit()
	}
}

var basePrompt = func(prompt string, callback checker) string {
	fmt.Print(prompt + "\n")

	defer func() {
		if err := rl.Close(); err != nil {
			renderError(err)

			errorExit()
		}
	}()

	for {
		line, err := rl.Readline()

		if err != nil {
			renderError(err)

			errorExit()
		}

		err = callback(line)

		if err != nil {
			renderError(err)

			continue
		}

		return line
	}
}
