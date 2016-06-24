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
	"github.com/spf13/cobra"

	"github.com/antham/doc-hunt/file"
	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

var failOnError bool

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if documentation update could be needed",
	Run: func(cmd *cobra.Command, args []string) {
		status := file.FetchStatus()

		var hasErrors bool

		for _, s := range *status {
			if len(s.Status[file.Deleted]) > 0 || len(s.Status[file.Updated]) > 0 || len(s.Status[file.Failed]) > 0 {
				hasErrors = true

				break
			}
		}

		if hasErrors {
			renderCheck(status)
		} else {
			ui.Success("No changes found")
		}

		if failOnError {
			util.ErrorExit()
		} else {
			util.SuccessExit()
		}
	},
}

func init() {
	checkCmd.Flags().BoolVarP(&failOnError, "fail-on-error", "e", false, "return an error exit code (1)")

	RootCmd.AddCommand(checkCmd)
}
