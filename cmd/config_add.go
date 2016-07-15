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

		var err error

		switch true {
		case len(args) == 0:
			err = fmt.Errorf("Missing doc identifier")
		case len(args) == 1:
			err = fmt.Errorf("Missing source identifiers")
		case len(args) > 2:
			err = fmt.Errorf("No more than 2 arguments expected")
		}

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}

		doc, sources, err := parseConfigAddArgs(args)

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}

		err = file.CreateConfig(doc, sources)

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}

		ui.Success("Config added")
		util.SuccessExit()
	},
}

func parseConfigAddArgs(args []string) (*file.Doc, *[]file.Source, error) {
	identifier := args[0]
	category, err := parseDocCategory(identifier)

	doc := file.NewDoc(identifier, category)

	if err != nil {
		return &file.Doc{}, &[]file.Source{}, err
	}

	sources, err := parseSources(doc, strings.Split(args[1], ","))

	if err != nil {
		return &file.Doc{}, &[]file.Source{}, err
	}

	return doc, sources, nil
}

func parseDocCategory(docIdentifier string) (file.DocCategory, error) {
	docFilename := util.GetAbsPath(docIdentifier)

	fileInfo, fileErr := os.Stat(docFilename)
	URL, URLErr := url.Parse(docIdentifier)

	if fileErr == nil && !fileInfo.IsDir() {
		return file.DFILE, nil
	} else if fileErr == nil && fileInfo.IsDir() {
		return file.DFOLDER, nil
	} else if URLErr == nil && URL.IsAbs() {
		return file.DURL, nil
	}

	return file.DERROR, fmt.Errorf("Doc %s is not a valid existing file, nor a valid existing folder, nor a valid URL", docIdentifier)
}

func parseSources(doc *file.Doc, identifiers []string) (*[]file.Source, error) {
	sources := []file.Source{}
	var cat file.SourceCategory

	for _, identifier := range identifiers {
		path := util.GetAbsPath(identifier)

		f, err := os.Stat(path)

		if os.IsNotExist(err) {
			return &sources, fmt.Errorf("Source identifier %s doesn't exist", identifier)
		}

		if f.IsDir() {
			cat = file.SFOLDER
		} else {
			cat = file.SFILE
		}

		source := file.NewSource(doc, identifier, cat)
		sources = append(sources, *source)
	}

	return &sources, nil
}

func init() {
	configCmd.AddCommand(addConfigCmd)
}
