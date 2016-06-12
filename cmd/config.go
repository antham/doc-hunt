// Copyright © 2016 Anthony HAMON <hamon.anth@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "List, add or delete configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new configuration row",
	Run: func(cmd *cobra.Command, args []string) {
		doc, docCat, fileSources, err := parseConfigAddArgs(args)

		if err != nil {
			renderError(err)

			errorExit()
		}

		addConfig(doc, docCat, fileSources)
	},
}

var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete one or several configuration row",
	Run: func(cmd *cobra.Command, args []string) {
		list := file.ListConfig()

		if len(*list) == 0 {
			renderInfo("No config added yet")

			successExit()
		}

		renderList(list)

		configs, err := promptConfigToRemove(list)

		if err != nil {
			renderError(err)

			errorExit()
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

func parseConfigAddArgs(args []string) (string, file.DocCategory, []string, error) {
	var docCategory file.DocCategory

	if len(args) == 0 {
		return "", docCategory, []string{}, fmt.Errorf("Missing file doc")
	}

	doc := args[0]

	_, fileErr := os.Stat(doc)
	URL, URLErr := url.Parse(doc)

	if fileErr == nil {
		docCategory = file.FILE
	} else if URLErr == nil && URL.IsAbs() {
		docCategory = file.URL
	} else {
		return "", docCategory, []string{}, fmt.Errorf("Doc %s is not a valid existing file, nor a valid URL", doc)
	}

	if len(args) == 1 {
		return "", docCategory, []string{}, fmt.Errorf("Missing file sources")
	}

	fileSources := strings.Split(args[1], ",")

	for _, fileSource := range fileSources {
		if _, err := os.Stat(fileSource); os.IsNotExist(err) {
			return "", docCategory, []string{}, fmt.Errorf("File source %s doesn't exist", fileSource)
		}
	}

	return doc, docCategory, fileSources, nil
}

func listConfig() {
	list := file.ListConfig()

	if len(*list) == 0 {
		renderInfo("No config added yet")

		successExit()
	}

	renderList(list)

	successExit()
}

func addConfig(identifier string, docCat file.DocCategory, fileSources []string) {
	doc := file.NewDoc(identifier, docCat)
	sources := file.NewSources(doc, fileSources)
	file.InsertConfig(doc, sources)

	renderSuccess("Config added")

	successExit()
}

func delConfig(configs *[]file.Config) {
	file.RemoveConfigs(configs)
}

func promptConfigToRemove(configs *[]file.Config) (*[]file.Config, error) {
	renderPrompt()
	rl, err := readline.New(">> ")

	if err != nil {
		return nil, fmt.Errorf("Something wrong happened during argument fetching")
	}

	defer rl.Close()

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
