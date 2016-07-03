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
		updateDeletedAndAddedFiles()

		file.UpdateItemsFingeprint()

		ui.Success("Update configuration succeeded")
		util.SuccessExit()
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}

func updateDeletedAndAddedFiles() {
	deleted := map[string]bool{}
	added := []file.Item{}

	for _, result := range *file.BuildStatus() {
		for _, filename := range result.Status[file.IDELETED] {
			deleted[filename] = true
		}

		if _, ok := result.Status[file.IADDED]; ok == true {
			items := result.Status[file.IADDED]

			for _, item := range *file.NewItems(&items, &result.Source) {
				added = append(added, item)
			}
		}
	}

	extractDeletedFiles := func(filenames *map[string]bool) *[]string {
		results := make([]string, len(*filenames))

		for filename := range *filenames {
			results = append(results, filename)
		}

		return &results
	}

	file.InsertItems(&added)
	file.DeleteItems(extractDeletedFiles(&deleted))
}
