package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/cio"
)

const formsHelpSlug = `Print help for "forms" subcommands`

func cmdFormsHelp(args []string) (err error) {
	var (
		helpText string
		c        = cio.Open()
	)
	if len(args) != 0 {
		switch args[0] {
		case "install":
			helpText = formsInstallHelp
		case "list", "l":
			helpText = formsListHelp
		default:
			c.Error("There is no command or help topic %q.", "forms "+args[0])
		}
	}
	if helpText == "" {
		helpText = formsHelp
	}
	helpText = strings.TrimLeft(helpText, "\n") // Allows newline after `
	io.WriteString(os.Stdout, c.WrapText(helpText))
	return nil
}
