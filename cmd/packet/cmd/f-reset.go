package cmd

import (
	"os"

	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/rothskeller/packet/v4/form/formdefs"
	"github.com/spf13/pflag"
)

const (
	// formsResetSlug = `Reset all forms to the built-in versions`
	formsResetHelp = `
usage: packet forms reset

The "packet forms reset" command removes all reseted forms and replaces them with the versions built into the packet executable.

(This command is for development use and may be removed in future versions.)
`
)

func cmdFormsReset(args []string) (err error) {
	flags := pflag.NewFlagSet("f-reset", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdFormsHelp([]string{"reset"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(formsResetHelp)
	}
	if len(args) != 0 {
		return usage(formsResetHelp)
	}
	os.RemoveAll(formdefs.FormsDir())
	registerForms()
	return nil
}
