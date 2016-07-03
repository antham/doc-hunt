package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"

	"github.com/antham/doc-hunt/file"
	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "List, add or delete configuration",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}
	},
}

var addCmd = &cobra.Command{
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

var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete one or several configuration row",
	Run: func(cmd *cobra.Command, args []string) {
		list := file.ListConfig()

		if len(*list) == 0 {
			ui.Info("No config added yet")

			util.SuccessExit()
		}

		renderConfig(list)

		configs, err := promptConfigToRemove(list)

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}

		delConfig(configs)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all recorded configurations",
	Run: func(cmd *cobra.Command, args []string) {
		listConfig()
	},
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

func listConfig() {
	list := file.ListConfig()

	if len(*list) == 0 {
		ui.Info("No config added yet")

		util.SuccessExit()
	}

	renderConfig(list)

	util.SuccessExit()
}

func addConfig(identifier string, docCat file.DocCategory, folderSources []string, fileSources []string) {
	doc := file.NewDoc(identifier, docCat)
	file.InsertDoc(doc)

	for _, identifier := range folderSources {
		source := file.NewSource(doc, identifier, file.SFOLDER)
		file.InsertSource(source)

		items := file.NewItems(util.ExtractFolderFiles(identifier), source)
		file.InsertItems(items)
	}

	for _, identifier := range fileSources {
		source := file.NewSource(doc, identifier, file.SFILE)
		file.InsertSource(source)

		items := file.NewItems(&[]string{identifier}, source)
		file.InsertItems(items)
	}

	ui.Success("Config added")

	util.SuccessExit()
}

func delConfig(configs *[]file.Config) {
	file.RemoveConfigs(configs)
}

func promptConfigToRemove(configs *[]file.Config) (*[]file.Config, error) {
	ui.Prompt("Choose configurations number to remove, each separated with a comma")
	rl, err := readline.New(">> ")

	if err != nil {
		return nil, fmt.Errorf("Something wrong happened during argument fetching")
	}

	defer func() {
		if err := rl.Close(); err != nil {
			ui.Error(err)

			util.ErrorExit()
		}
	}()

	line, _ := rl.Readline()

	return parseConfigDelArgs(configs, line)
}

func parseConfigDelArgs(configs *[]file.Config, line string) (*[]file.Config, error) {
	results := []file.Config{}

	for _, sel := range strings.Split(line, ",") {
		strings.TrimSpace(sel)
		n, err := strconv.Atoi(sel)

		if err != nil {
			return nil, fmt.Errorf("%s is not a number", sel)
		}

		if n < 0 || n >= len(*configs) {
			return nil, fmt.Errorf("Value %d is out of bounds", n)
		}

		results = append(results, (*configs)[n])
	}

	return &results, nil
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(addCmd)
	configCmd.AddCommand(delCmd)
	configCmd.AddCommand(listCmd)
}
