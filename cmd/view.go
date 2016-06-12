package cmd

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/antham/doc-hunt/file"
)

func renderList(list *[]file.Config) {
	color.Magenta("----")

	for i, config := range *list {
		out := color.CyanString(fmt.Sprintf("%d", i))
		out += color.YellowString(" - document : ")
		out += config.Doc.Identifier
		out += color.YellowString(" => sources : ")

		for j, source := range config.Sources {
			if j != 0 {
				out += color.YellowString(", ")
			}

			out += source.Path
		}

		fmt.Printf("%s\n", out)
		color.Magenta("----")
	}
}

func renderPrompt() {
	color.Yellow(fmt.Sprintf("\nChoose configurations number to remove, each separated with a comma : "))
}

func renderError(err error) {
	color.Red(err.Error())
}

func renderSuccess(str string) {
	color.Green(str)
}

func renderInfo(str string) {
	color.Cyan(str)
}
