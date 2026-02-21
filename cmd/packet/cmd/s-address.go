package cmd

import (
	"fmt"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/cmd/packet/server"
	"github.com/spf13/pflag"
)

const (
	serverAddressSlug = `Prints the URL of the internal web server`
	serverAddressHelp = `
usage: packet server address

The "packet server address" command prints the URL of the internal web server if it is running.
`
)

func cmdServerAddress(args []string) (err error) {
	flags := pflag.NewFlagSet("s-address", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdServerHelp([]string{"address"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(serverAddressHelp)
	}
	if len(args) != 0 {
		return usage(serverAddressHelp)
	}
	if addr, err := server.GetAddress(false); err != nil {
		return err
	} else {
		fmt.Println(addr)
	}
	return nil
}
