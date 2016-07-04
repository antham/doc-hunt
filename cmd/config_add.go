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
	var err error
	folderSources := []string{}
	fileSources := []string{}

	switch len(args) {
	case 0:
		err = fmt.Errorf("Missing doc identifier")
	case 1:
		err = fmt.Errorf("Missing source identifiers")
	}

	if err != nil {
		return "", docCategory, folderSources, fileSources, err
	}

	docIdentifier := args[0]
	docCategory, err = parseDocCategory(docIdentifier)

	if err != nil {
		return "", docCategory, folderSources, fileSources, err
	}

	folderSources, fileSources, err = parseSources(strings.Split(args[1], ","))

	if err != nil {
		return "", docCategory, folderSources, fileSources, err
	}

	return docIdentifier, docCategory, folderSources, fileSources, nil
}

func parseDocCategory(docIdentifier string) (file.DocCategory, error) {
	docFilename := util.GetAbsPath(docIdentifier)

	_, fileErr := os.Stat(docFilename)
	URL, URLErr := url.Parse(docIdentifier)

	if fileErr == nil {
		return file.DFILE, nil
	} else if URLErr == nil && URL.IsAbs() {
		return file.DURL, nil
	} else {
		return file.DERROR, fmt.Errorf("Doc %s is not a valid existing file, nor a valid URL", docIdentifier)
	}
}

func parseSources(sources []string) ([]string, []string, error) {
	folderSources := []string{}
	fileSources := []string{}

	for _, source := range sources {
		path := util.GetAbsPath(source)

		f, err := os.Stat(path)

		if os.IsNotExist(err) {
			return folderSources, fileSources, fmt.Errorf("Source identifier %s doesn't exist", source)
		}

		if f.IsDir() {
			folderSources = append(folderSources, source)
		} else {
			fileSources = append(fileSources, source)
		}
	}

	return folderSources, fileSources, nil
}

func init() {
	configCmd.AddCommand(addConfigCmd)
}
