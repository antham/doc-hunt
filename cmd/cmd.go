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

	"github.com/chzyer/readline"
)

type checker func(line string) error

type terminalReader interface {
	Readline() (string, error)
	Close() error
}

var rl terminalReader

func init() {
	var err error

	rl, err = readline.New(">> ")

	if err != nil {
		renderError(fmt.Errorf("Something went wrong when initializing prompt"))

		errorExit()
	}
}

func errorExit() {
	os.Exit(1)
}

func successExit() {
	os.Exit(0)
}

func basePrompt(prompt string, callback checker) string {
	fmt.Print(prompt + "\n")

	defer func() {
		if err := rl.Close(); err != nil {
			renderError(err)

			errorExit()
		}
	}()

	for {
		line, err := rl.Readline()

		if err != nil {
			renderError(err)

			errorExit()
		}

		err = callback(line)

		if err != nil {
			renderError(err)

			continue
		}

		return line
	}
}
