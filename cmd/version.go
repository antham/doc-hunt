package cmd

import (
	"github.com/antham/doc-hunt/file"
	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "App version",
	Run: func(cmd *cobra.Command, args []string) {
		ui.Info("v" + file.Container.GetVersion().Get())

		util.SuccessExit()
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
