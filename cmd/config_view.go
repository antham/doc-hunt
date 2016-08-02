package cmd

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/antham/doc-hunt/file"
)

func renderConfig(list *[]file.Config) {
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

			out += source.Identifier
		}

		fmt.Printf("%s\n", out)
		color.Magenta("----")
	}
}
