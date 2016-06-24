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
		status := file.FetchStatus()

		var hasErrors bool

		for _, s := range *status {
			if len(s.Status[file.Deleted]) > 0 || len(s.Status[file.Updated]) > 0 || len(s.Status[file.Failed]) > 0 {
				hasErrors = true

				break
			}
		}

		if hasErrors {
			renderCheck(status)
		} else {
			ui.Success("No changes found")
		}

		if failOnError {
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
