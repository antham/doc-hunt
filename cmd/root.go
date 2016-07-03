package cmd

import (
	"github.com/spf13/cobra"

	"github.com/antham/doc-hunt/file"
	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "doc-hunt",
	Short: "Ensure your documentation is up to date",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := file.Initialize(); err != nil {
		ui.Error(err)

		util.ErrorExit()
	}

	if err := RootCmd.Execute(); err != nil {
		ui.Error(err)
		util.ErrorExit()
	}
}
