package cmd

import (
	"fmt"

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
		resolveDeletedAndMovedFiles()

		file.UpdateSourcesFingeprint()

		ui.Success("Update configuration succeeded")
		util.SuccessExit()
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}

func resolveDeletedAndMovedFiles() {
	deleted := map[string]bool{}
	moved := map[string]string{}

	for _, s := range *file.FetchStatus() {
		for _, filename := range s.Status[file.Deleted] {
			if _, ok := deleted[filename]; ok == true {
				continue
			}

			if _, ok := moved[filename]; ok == true {
				continue
			}

			basePrompt(fmt.Sprintf(`File "%s" is removed : rename (r) or delete (d) ?`, filename), globalPrompt(filename, &deleted, &moved))
		}
	}

	extractDeletedFiles := func(filenames *map[string]bool) *[]string {
		results := make([]string, len(*filenames))

		for filename := range *filenames {
			results = append(results, filename)
		}

		return &results
	}

	file.DeleteSources(extractDeletedFiles(&deleted))
	file.UpdateFilenameSources(&moved)
}
