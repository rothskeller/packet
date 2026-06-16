package cmd

import (
	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/rothskeller/packet/v4/errors"
	"github.com/spf13/pflag"
)

const (
	serverSlug = `Subcommands to control the built-in web server`
	serverHelp = `
The "packet server" command provides multiple subcommands to control the packet software's internal web server.  These commands are not generally used by humans.

Available subcommands include:
  address  ⇥` + serverAddressSlug + `
  help     ⇥` + serverHelpSlug + `
  start    ⇥` + serverStartSlug + `
  stop     ⇥` + serverStopSlug + `
For help on a subcommand, run "packet server help «command»".
`
)

func cmdServer(args []string) (err error) {
	var flags pflag.FlagSet

	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"server"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(serverHelp)
	}
	if len(args) == 0 {
		return cmdHelp([]string{"server"})
	}
	switch args[0] {
	case "address":
		return cmdServerAddress(args[1:])
	case "help", "h":
		return cmdServerHelp(args[1:])
	case "start":
		return cmdServerStart(args[1:])
	case "stop":
		return cmdServerStop(args[1:])
	default:
		return errors.NewF("No such command %q.", "server "+args[0])
	}
}
