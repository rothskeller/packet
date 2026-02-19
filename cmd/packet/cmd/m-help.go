package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/cio"
)

const manualHelpSlug = `Print help for "manual" subcommands`

func cmdManualHelp(args []string) (err error) {
	var (
		helpText string
		c        = cio.Open()
	)
	if len(args) != 0 {
		switch args[0] {
		case "dr", "receipt":
			helpText = manualDRHelp
		case "receive", "r":
			helpText = manualReceiveHelp
		case "send", "s":
			helpText = manualSendHelp
		default:
			c.Error("There is no command or help topic %q.", "manual "+args[0])
		}
	}
	if helpText == "" {
		helpText = manualHelp
	}
	helpText = strings.TrimLeft(helpText, "\n") // Allows newline after `
	io.WriteString(os.Stdout, c.WrapText(helpText))
	return nil
}
