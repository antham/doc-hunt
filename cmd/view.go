package cmd

import (
	"github.com/fatih/color"
)

func renderError(err error) {
	color.Red(err.Error())
}

func renderSuccess(str string) {
	color.Green(str)
}
