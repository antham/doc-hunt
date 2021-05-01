package cmd

import (
	"github.com/spf13/cobra"

	"github.com/antham/doc-hunt/file"
	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

var failOnError bool

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if documentation update could be needed",
	Run: func(cmd *cobra.Command, args []string) {
		itemStatus, changesOccured, err := file.Container.GetManager().GetItemStatus()

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}

		if changesOccured {
			renderCheck(itemStatus)
		} else {
			ui.Success("No changes found")
		}

		if failOnError && changesOccured {
			util.ErrorExit()
		} else {
			util.SuccessExit()
		}
	},
}

func init() {
	checkCmd.Flags().BoolVarP(&failOnError, "fail-on-error", "e", false, "return an error exit code (1)")

	RootCmd.AddCommand(checkCmd)
}
