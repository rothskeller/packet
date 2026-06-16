package cmd

import (
	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/rothskeller/packet/v4/errors"
	"github.com/spf13/pflag"
)

const (
	manualSlug = `Subcommands related to manual JNOS operations`
	manualHelp = `
The "packet manual" (or "man") command provides multiple subcommands related to manual operations (where the operator interacts with JNOS directly rather than this packet software).

Available subcommands include:
  dr       ⇥` + manualDRSlug + `
  help     ⇥` + manualHelpSlug + `
  receive  ⇥` + manualReceiveSlug + `
  send     ⇥` + manualSendSlug + `
For help on a subcommand, run "packet manual help «command»".
`
)

func cmdManual(args []string) (err error) {
	var flags pflag.FlagSet

	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"manual"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(manualHelp)
	}
	if len(args) == 0 {
		return cmdHelp([]string{"manual"})
	}
	switch args[0] {
	case "dr", "receipt":
		return cmdManualDR(args[1:])
	case "help", "h":
		return cmdManualHelp(args[1:])
	case "receive", "r":
		return cmdManualReceive(args[1:])
	case "send", "s":
		return cmdManualSend(args[1:])
	default:
		return errors.NewF("No such command %q.", "manual "+args[0])
	}
}
