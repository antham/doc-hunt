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

var dryRun bool

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

		if dryRun {
			execDryRun(doc, sources)
		}

		createConfig(doc, sources)
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

	for _, identifier := range identifiers {
		parsedIdentifier, cat := file.ParseIdentifier(identifier)

		if cat == file.SERROR {
			return &[]file.Source{}, fmt.Errorf("Source identifier %s is not correct", identifier)
		}

		source := file.NewSource(doc, parsedIdentifier, cat)
		sources = append(sources, *source)
	}

	return &sources, nil
}

func execDryRun(doc *file.Doc, sources *[]file.Source) {
	datas := map[string]*[]string{}

	for _, source := range *sources {
		files, err := util.ExtractFilesMatchingReg(source.Identifier)

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}

		datas[source.Identifier] = files
	}

	renderDryRun(doc, &datas)
	util.SuccessExit()
}

func createConfig(doc *file.Doc, sources *[]file.Source) {
	c := file.Container.GetConfigRepository()
	err := c.CreateFromDocAndSources(doc, sources)

	if err != nil {
		ui.Error(err)

		util.ErrorExit()
	}

	ui.Success("Config added")
	util.SuccessExit()
}

func init() {
	configCmd.AddCommand(addConfigCmd)
	addConfigCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "simulate what a config would record")
}
