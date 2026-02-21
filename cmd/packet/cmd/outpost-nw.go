//go:build !windows

package cmd

import (
	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/spf13/pflag"
)

const (
	outpostSlug = `Subcommands related to Outpost integration`
	outpostHelp = `
The "packet outpost" command provides multiple subcommands related to its integration with the Outpost packet message manager on Windows.  These commands are not available on other operating systems.
`
)

func cmdOutpost(args []string) (err error) {
	var flags pflag.FlagSet

	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"outpost"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(outpostHelp)
	}
	return cmdHelp([]string{"outpost"})
}

func cmdOutpostHelp(args []string) (err error) {
	return cmdHelp([]string{"outpost"})
}
