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
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/antham/doc-hunt/model"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "List, add or delete configuration",
	Run: func(cmd *cobra.Command, args []string) {
		genErr := fmt.Errorf("Unvalid argument choose one those : list, add or delete")

		if len(args) == 0 {
			renderError(genErr)

			errorExit()
		}

		switch args[0] {
		case "add":
			docSource, fileSources, err := parseConfigAddArgs(args[1:])

			if err != nil {
				renderError(err)

				errorExit()
			}

			addConfig(docSource, fileSources)
		default:
			renderError(genErr)

			errorExit()
		}
	},
}

func parseConfigAddArgs(args []string) (string, []string, error) {
	if len(args) == 0 {
		return "", []string{}, fmt.Errorf("Missing file doc")
	}

	fileDoc := args[0]

	if _, err := os.Stat(fileDoc); os.IsNotExist(err) {
		return "", []string{}, fmt.Errorf("File doc %s doesn't exist", fileDoc)
	}

	if len(args) == 1 {
		return "", []string{}, fmt.Errorf("Missing file sources")
	}

	fileSources := strings.Split(args[1], ",")

	for _, fileSource := range fileSources {
		if _, err := os.Stat(fileSource); os.IsNotExist(err) {
			return "", []string{}, fmt.Errorf("File source %s doesn't exist", fileSource)
		}
	}

	return fileDoc, fileSources, nil
}

func addConfig(fileDoc string, fileSources []string) {
	doc := model.NewDoc(fileDoc)
	sources := model.NewSources(doc, fileSources)
	model.InsertConfig(doc, sources)

	renderSuccess("Config added")

	successExit()
}

func init() {
	RootCmd.AddCommand(configCmd)
}
