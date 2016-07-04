package cmd

import (
	"github.com/spf13/cobra"

	"github.com/antham/doc-hunt/file"
	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

var listConfigCmd = &cobra.Command{
	Use:   "list",
	Short: "List all recorded configurations",
	Run: func(cmd *cobra.Command, args []string) {
		listConfig()
	},
}

func listConfig() {
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

	util.SuccessExit()
}

func init() {
	configCmd.AddCommand(listConfigCmd)
}
