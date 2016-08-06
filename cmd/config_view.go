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

			out += source.Identifier
		}

		fmt.Printf("%s\n", out)
		color.Magenta("----")
	}
}

func renderDryRun(doc *file.Doc, datas *map[string]*[]string) {
	out := fmt.Sprintf(`%s : %s`, color.CyanString("Document"), color.YellowString(doc.Identifier))
	out += fmt.Sprintf("\n\n")

	for identifier, files := range *datas {
		out += fmt.Sprintf(`%s "%s" : `, color.CyanString("Files matching regexp"), color.YellowString(identifier))
		out += fmt.Sprintf("\n")

		if len(*files) == 0 {
			out += fmt.Sprintf("    => %s\n", color.RedString("No files found"))
		}

		for _, file := range *files {
			out += fmt.Sprintf("    => %s\n", file)
		}

		out += fmt.Sprintf("\n")
	}

	out = strings.TrimRight(out, "\n")

	fmt.Printf("%s\n", out)
}
