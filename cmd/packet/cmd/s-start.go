package cmd

import (
	"os"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/cmd/packet/server"
	"github.com/spf13/pflag"
)

const (
	serverStartSlug = `Starts the internal web server`
	serverStartHelp = `
usage: packet server start

The "packet server start" command starts the internal web server used to serve forms and other web pages for the packet software.  This command is invoked automatically when needed and should not be invoked manually.  If the server is already running, this command exits silently.  Otherwise, this command does not exit until the server is stopped (by idle timeout, POST /stop request, touch of the stop file, or signal).
`
)

func cmdServerStart(args []string) (err error) {
	flags := pflag.NewFlagSet("s-start", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdServerHelp([]string{"start"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(serverStartHelp)
	}
	if len(args) != 0 {
		return usage(serverStartHelp)
	}
	registerForms()
	server.Start(os.Stdout)
	return nil // not reachable
}
