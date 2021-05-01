package cmd

import (
	"github.com/spf13/cobra"

	"github.com/antham/doc-hunt/file"
	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update documentation references",
	Run: func(cmd *cobra.Command, args []string) {
		err := file.Container.GetManager().UpdateFingerprints()

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}

		ui.Success("Update succeeded")
		util.SuccessExit()
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
