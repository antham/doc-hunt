package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/antham/doc-hunt/file"
	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

var addConfigCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new configuration row",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 2 {
			ui.Error(fmt.Errorf("No more than 2 arguments expected"))

			util.ErrorExit()
		}

		doc, docCat, folderSources, fileSources, err := parseConfigAddArgs(args)

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}

		addConfig(doc, docCat, folderSources, fileSources)
	},
}

func addConfig(docIdentifier string, docCat file.DocCategory, folderSources []string, fileSources []string) {

	err := file.CreateConfig(docIdentifier, docCat, folderSources, fileSources)

	if err != nil {
		ui.Error(err)

		util.ErrorExit()
	}

	ui.Success("Config added")

	util.SuccessExit()
}

func parseConfigAddArgs(args []string) (string, file.DocCategory, []string, []string, error) {
	var docCategory file.DocCategory
	folderSources := []string{}
	fileSources := []string{}

	if len(args) == 0 {
		return "", docCategory, folderSources, fileSources, fmt.Errorf("Missing file doc")
	}

	doc := args[0]
	docFilename := util.GetAbsPath(doc)

	_, fileErr := os.Stat(docFilename)
	URL, URLErr := url.Parse(doc)

	if fileErr == nil {
		docCategory = file.DFILE
	} else if URLErr == nil && URL.IsAbs() {
		docCategory = file.DURL
	} else {
		return "", docCategory, folderSources, fileSources, fmt.Errorf("Doc %s is not a valid existing file, nor a valid URL", docFilename)
	}

	if len(args) == 1 {
		return "", docCategory, folderSources, fileSources, fmt.Errorf("Missing file/folder sources")
	}

	for _, source := range strings.Split(args[1], ",") {
		path := util.GetAbsPath(source)

		f, err := os.Stat(path)

		if os.IsNotExist(err) {
			return "", docCategory, folderSources, fileSources, fmt.Errorf("File/folder source %s doesn't exist", source)
		}

		if f.IsDir() {
			folderSources = append(folderSources, source)
		} else {
			fileSources = append(fileSources, source)
		}
	}

	return doc, docCategory, folderSources, fileSources, nil
}

func init() {
	configCmd.AddCommand(addConfigCmd)
}
