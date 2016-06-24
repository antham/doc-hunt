package cmd

import (
	"fmt"
	"strings"

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

			out += source.Path
		}

		fmt.Printf("%s\n", out)
		color.Magenta("----")
	}
}

func renderCheck(list *[]file.Result) {
	output := ""

	for _, r := range *list {
		if len(r.Status[file.Deleted]) != 0 || len(r.Status[file.Updated]) != 0 || len(r.Status[file.Failed]) != 0 {
			output += color.CyanString(r.Doc.Identifier + "\n")

			for status, sources := range r.Status {
				switch status {
				case file.Updated:
					output += fmt.Sprintf("\n  %s \n\n", color.YellowString(strings.ToLower(status.String())))
				case file.Deleted, file.Failed:
					output += fmt.Sprintf("\n  %s \n\n", color.RedString(strings.ToLower(status.String())))
				}

				if status != file.Untouched {

					for _, s := range sources {
						output += fmt.Sprintf("    => %s\n", s)
					}
				}
			}

			output += color.MagentaString("----\n")
		}
	}

	color.Magenta("----")

	fmt.Print(output)
}
