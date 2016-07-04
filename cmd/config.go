package cmd

import (
	"github.com/spf13/cobra"

	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "List, add or delete configuration",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
}
