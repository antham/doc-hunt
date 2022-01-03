package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/antham/doc-hunt/file"
)

var out io.Writer = os.Stdout

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
	output := fmt.Sprintf(`%s : %s`, color.CyanString("Document"), color.YellowString(doc.Identifier))
	output += "\n\n"

	for identifier, files := range *datas {
		output += fmt.Sprintf(`%s "%s" : `, color.CyanString("Files matching regexp"), color.YellowString(identifier))
		output += "\n"

		if len(*files) == 0 {
			output += fmt.Sprintf("    => %s\n", color.RedString("No files found"))
		}

		for _, file := range *files {
			output += fmt.Sprintf("    => %s\n", file)
		}

		output += "\n"
	}

	output = strings.TrimRight(output, "\n")

	fmt.Fprintf(out, "%s\n", output)
}
