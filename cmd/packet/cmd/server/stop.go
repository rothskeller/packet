package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rothskeller/packet/cmd/packet/server"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:          "stop",
	Short:        "Stops the internal web server",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
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
	},
}

func init() {
	Command.AddCommand(stopCmd)
}
