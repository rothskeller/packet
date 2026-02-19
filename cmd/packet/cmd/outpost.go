//go:build windows

package cmd

import (
	"regexp"
	"strconv"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/errors"
	"github.com/spf13/pflag"
)

const (
	outpostSlug = `Subcommands related to Outpost integration`
	outpostHelp = `
The "packet outpost" command provides multiple subcommands related to its integration with the Outpost packet message manager.  These commands are not generally used by humans.

Available subcommands include:
  convert    ⇥` + outpostConvertSlug + `
  edit       ⇥` + outpostEditSlug + `
  help       ⇥` + outpostHelpSlug + `
  install    ⇥` + outpostInstallSlug + `
  new        ⇥` + outpostNewSlug + `
  uninstall  ⇥` + outpostUninstallSlug + `
  view       ⇥` + outpostViewSlug + `
For help on a subcommand, run "packet outpost help «command»".
`
)

func cmdOutpost(args []string) (err error) {
	var flags pflag.FlagSet

	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"outpost"})
	} else if err != nil {
		cio.Open().Error("%s", err.Error())
		return usage(outpostHelp)
	}
	if len(args) == 0 {
		return cmdHelp([]string{"outpost"})
	}
	switch args[0] {
	case "convert":
		return cmdOutpostConvert(args[1:])
	case "edit":
		return cmdOutpostEdit(args[1:])
	case "help", "h":
		return cmdOutpostHelp(args[1:])
	case "install":
		return cmdOutpostInstall(args[1:])
	case "new":
		return cmdOutpostNew(args[1:])
	case "uninstall":
		return cmdOutpostUninstall(args[1:])
	case "view":
		return cmdOutpostView(args[1:])
	default:
		return errors.NewF("No such command %q.", "outpost "+args[0])
	}
}

var outpostArgEncodeRE = regexp.MustCompile(`(?i)~[0-9a-f][0-9a-f]`)

// decodeOutpostArg reverses Outpost's argument encoding.
func decodeOutpostArg(s string) string {
	return outpostArgEncodeRE.ReplaceAllStringFunc(s, func(e string) string {
		code, _ := strconv.ParseInt(e[1:], 16, 64)
		return string(byte(code))
	})
}
