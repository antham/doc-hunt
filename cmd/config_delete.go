package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"

	"github.com/antham/doc-hunt/file"
	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

var delConfigCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete one or several configuration row",
	Run: func(cmd *cobra.Command, args []string) {
		list, err := file.ListConfig()

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}

		if len(*list) == 0 {
			ui.Info("No config added yet")

			util.SuccessExit()
		}

		renderConfig(list)

		configs, err := promptConfigToRemove(list)

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}

		delConfig(configs)
	},
}

func delConfig(configs *[]file.Config) {
	err := file.RemoveConfigs(configs)

	if err != nil {
		ui.Error(err)

		util.ErrorExit()
	}
}

func promptConfigToRemove(configs *[]file.Config) (*[]file.Config, error) {
	ui.Prompt("Choose configurations number to remove, each separated with a comma")
	rl, err := readline.New(">> ")

	if err != nil {
		return nil, fmt.Errorf("Something wrong happened during argument fetching")
	}

	defer func() {
		if err := rl.Close(); err != nil {
			ui.Error(err)

			util.ErrorExit()
		}
	}()

	line, _ := rl.Readline()

	return parseConfigDelArgs(configs, line)
}

func parseConfigDelArgs(configs *[]file.Config, line string) (*[]file.Config, error) {
	results := []file.Config{}

	for _, sel := range strings.Split(line, ",") {
		strings.TrimSpace(sel)
		n, err := strconv.Atoi(sel)

		if err != nil {
			return nil, fmt.Errorf("%s is not a number", sel)
		}

		if n < 0 || n >= len(*configs) {
			return nil, fmt.Errorf("Value %d is out of bounds", n)
		}

		results = append(results, (*configs)[n])
	}

	return &results, nil
}

func init() {
	configCmd.AddCommand(delConfigCmd)
}