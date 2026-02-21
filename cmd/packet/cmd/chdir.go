package cmd

import (
	"os"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/spf13/pflag"
)

const (
	chdirSlug = `Switch to a different directory`
	chdirHelp = `
usage: cd «directory»

The "cd" (or "chdir") command switches to the named directory.  The directory does not have to contain incident data.

This command is only useful from within the packet shell, since it does not affect the working directory of the process that invoked the shell.
`
)

func cmdChdir(args []string) (err error) {
	flags := pflag.NewFlagSet("chdir", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"chdir"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(chdirHelp)
	}
	if len(args) != 1 {
		return usage(chdirHelp)
	}
	return os.Chdir(args[0])
}
