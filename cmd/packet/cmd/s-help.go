package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/rothskeller/packet/v4/cmd/packet/cio"
)

const serverHelpSlug = `Print help for "server" subcommands`

func cmdServerHelp(args []string) (err error) {
	var (
		helpText string
		c        = cio.Open()
	)
	if len(args) != 0 {
		switch args[0] {
		case "address":
			helpText = serverAddressHelp
		case "start":
			helpText = serverStartHelp
		case "stop":
			helpText = serverStopHelp
		default:
			c.ErrorF("There is no command or help topic %q.", "server "+args[0])
		}
	}
	if helpText == "" {
		helpText = serverHelp
	}
	helpText = strings.TrimLeft(helpText, "\n") // Allows newline after `
	io.WriteString(os.Stdout, c.WrapText(helpText))
	return nil
}
