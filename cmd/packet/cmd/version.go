package cmd

import (
	"fmt"

	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/rothskeller/packet/v4/packetver"
	"github.com/spf13/pflag"
)

const (
	versionSlug = `Prints the current software version`
	versionHelp = `
usage: packet version

The "packet version" command displays the current software version number.  (For the versions of installed forms, use the "packet forms list" command.)
`
)

func cmdVersion(args []string) (err error) {
	flags := pflag.NewFlagSet("version", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"version"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(versionHelp)
	}
	if len(args) != 0 {
		return usage(versionHelp)
	}
	fmt.Println(packetver.Version)
	return nil
}
