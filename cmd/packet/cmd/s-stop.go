package cmd

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/rothskeller/packet/v4/cmd/packet/server"
	"github.com/spf13/pflag"
)

const (
	serverStopSlug = `Stops the internal web server`
	serverStopHelp = `
usage: packet server stop

The "packet server stop" command stops the internal web server if it is running.  It exits silently if the server is not running.
`
)

func cmdServerStop(args []string) (err error) {
	flags := pflag.NewFlagSet("s-stop", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdServerHelp([]string{"stop"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(serverStopHelp)
	}
	if len(args) != 0 {
		return usage(serverStopHelp)
	}
	if addr, err := server.GetAddress(false); err != nil {
		return err
	} else if addr == "" {
		return nil
	} else if resp, err := http.Post(addr+"/stop", "text/plain", nil); err != nil {
		slog.Error("POST /stop", "url", addr, "err", err)
		return err
	} else if resp.StatusCode >= 400 {
		slog.Error("POST /stop", "url", addr, "code", resp.StatusCode, "status", resp.Status)
		resp.Body.Close()
		return fmt.Errorf("stop request: %s", resp.Status)
	} else {
		slog.Info("sent /stop to server")
		resp.Body.Close()
		return nil
	}
}
