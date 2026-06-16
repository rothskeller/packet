package cmd

import (
	"flag"

	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/rothskeller/packet/v4/cmd/packet/server"
	"github.com/spf13/pflag"
)

const (
	serverStartSlug = `Starts the internal web server`
	serverStartHelp = `
usage: packet server start [-o] [-p port]
  --outpost, -o   ⇥Stop server when Outpost closes
  --port, -p      ⇥Port number for server

The "packet server start" command starts the internal web server used to serve forms and other web pages for the packet software.  This command is invoked automatically when needed and should not be invoked manually.  If the server is already running, this command prints the existing server address and exits.  Otherwise, this command prints the new server address, and does not exit until the server is stopped (by POST /stop request, touch of the stop file, or signal).

If the --outpost (or -o) flag is specified, the server will exit when it can no longer connect to opdirect (i.e., when Outpost is closed).

If a port number is specified with the --port (or -p) flag, any newly started server will operate on that port.  This is for debugging; normally a random port number is used.
`
)

func cmdServerStart(args []string) (err error) {
	var (
		outpost bool
		port    int
	)
	flags := pflag.NewFlagSet("s-start", pflag.ContinueOnError)
	flags.BoolVarP(&outpost, "outpost", "o", false, "stop when Outpost stops")
	flags.IntVarP(&port, "port", "p", 0, "port number for server")
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdServerHelp([]string{"start"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(serverStartHelp)
	}
	if flag.NArg() != 0 {
		return usage(serverStartHelp)
	}
	if port != 0 && (port < 1024 || port > 65535) {
		cio.Open().ErrorF("invalid port number")
		return usage(serverStartHelp)
	}
	registerForms()
	return server.Start(port, outpost)
}
