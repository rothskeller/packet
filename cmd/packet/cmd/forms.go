package cmd

import (
	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/rothskeller/packet/v4/errors"
	"github.com/spf13/pflag"
)

const (
	formsSlug = `Subcommands related to message forms`
	formsHelp = `
The "packet forms" command provides multiple subcommands related to packet message forms.

Available subcommands include:
  help     ⇥` + formsHelpSlug + `
  install  ⇥` + formsInstallSlug + `
  list     ⇥` + formsListSlug + `
  update   ⇥` + formsUpdateSlug + `
For help on a subcommand, run "packet forms help «command»".
` // Note, "reset" is not documented intentionally.
)

func cmdForms(args []string) (err error) {
	var flags pflag.FlagSet

	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"forms"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(formsHelp)
	}
	if len(args) == 0 {
		return cmdHelp([]string{"forms"})
	}
	switch args[0] {
	case "help", "h":
		return cmdFormsHelp(args[1:])
	case "install":
		return cmdFormsInstall(args[1:])
	case "list", "l":
		return cmdFormsList(args[1:])
	case "reset":
		return cmdFormsReset(args[1:])
	case "update":
		return cmdFormsUpdate(args[1:])
	default:
		return errors.NewF("No such command %q.", "forms "+args[0])
	}
}
