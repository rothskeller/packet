package server

import (
	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/cmd/packet/server"
	"github.com/rothskeller/packet/form/formdefs"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the internal web server",
	Long:  `Starts the internal web server used to serve forms and other web pages for the packet software.  This command is invoked automatically when needed and should not be invoked manually.  If the server is already running, this command exits silently.  Otherwise, this command does not exit until the server is stopped (by idle timeout, POST /stop request, touch of the stop file, or signal).`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := formdefs.RegisterForms(); err != nil {
			cio.Open().Warn("%s", err)
		}
		server.Start()
	},
}

func init() {
	Command.AddCommand(startCmd)
}
