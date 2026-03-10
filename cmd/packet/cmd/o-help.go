//go:build windows

package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/cio"
)

const outpostHelpSlug = `Print help for "outpost" subcommands`

func cmdOutpostHelp(args []string) (err error) {
	var (
		helpText string
		c        = cio.Open()
	)
	if len(args) != 0 {
		switch args[0] {
		case "convert":
			helpText = outpostConvertHelp
		case "edit":
			helpText = outpostEditHelp
		case "install":
			helpText = outpostInstallHelp
		case "new":
			helpText = outpostNewHelp
		case "uninstall":
			helpText = outpostUninstallHelp
		case "view":
			helpText = outpostViewHelp
		default:
			c.ErrorF("There is no command or help topic %q.", "outpost "+args[0])
		}
	}
	if helpText == "" {
		helpText = outpostHelp
	}
	helpText = strings.TrimLeft(helpText, "\n") // Allows newline after `
	io.WriteString(os.Stdout, c.WrapText(helpText))
	return nil
}
