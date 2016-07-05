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

func renderCheck(list *map[file.Doc]map[file.ItemStatus]map[string]bool) {
	output := ""

	for doc, status := range *list {
		if isCheckRenderable(&status) {
			output += color.CyanString(doc.Identifier + "\n")

			output += renderStatus(&status)

			output += color.MagentaString("----\n")
		}
	}

	color.Magenta("----")

	fmt.Print(output)
}

func isCheckRenderable(status *map[file.ItemStatus]map[string]bool) bool {
	return (len((*status)[file.IDELETED]) != 0 || len((*status)[file.IUPDATED]) != 0 || len((*status)[file.IFAILED]) != 0 || len((*status)[file.IADDED]) != 0)
}

func renderStatus(status *map[file.ItemStatus]map[string]bool) string {
	var output string

	for _, s := range []file.ItemStatus{file.IADDED, file.IUPDATED, file.IDELETED} {
		if len((*status)[s]) == 0 {
			continue
		}

		switch s {
		case file.IADDED:
			output += fmt.Sprintf("\n  %s \n\n", color.GreenString("Added"))
		case file.IUPDATED:
			output += fmt.Sprintf("\n  %s \n\n", color.YellowString("Updated"))
		case file.IFAILED:
			output += fmt.Sprintf("\n  %s \n\n", color.RedString("An error occured"))
		case file.IDELETED:
			output += fmt.Sprintf("\n  %s \n\n", color.RedString("Deleted"))
		}

		for filename := range (*status)[s] {
			output += fmt.Sprintf("    => %s\n", filename)
		}

	}

	return output
}
