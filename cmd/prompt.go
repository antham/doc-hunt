package cmd

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"

	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
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
		ui.Error(fmt.Errorf("Something went wrong when initializing prompt"))

		util.ErrorExit()
	}
}

var basePrompt = func(prompt string, callback checker) string {
	fmt.Print(prompt + "\n")

	defer func() {
		if err := rl.Close(); err != nil {
			ui.Error(err)

			util.ErrorExit()
		}
	}()

	for {
		line, err := rl.Readline()

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}

		err = callback(line)

		if err != nil {
			ui.Error(err)

			continue
		}

		return line
	}
}

var globalPrompt = func(filename string, deleted *map[string]bool, moved *map[string]string) checker {
	return func(line string) error {
		if line == "d" {
			(*deleted)[filename] = true

			return nil
		}

		if line == "r" {
			basePrompt(fmt.Sprintf(`Write an existing file name to rename "%s"`, filename), renamePrompt(filename, moved))

			return nil
		}

		return fmt.Errorf(`This action "%s" : doesn't exist`, line)
	}
}

var renamePrompt = func(path string, moved *map[string]string) checker {
	return func(line string) error {
		filename := util.GetAbsPath(line)
		_, err := os.Stat(filename)

		if err != nil && !os.IsExist(err) {
			return fmt.Errorf(`File "%s" doesn't exist, please enter an existing filename`, line)
		}

		(*moved)[path] = line

		return nil
	}
}
