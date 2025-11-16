package outpost

import (
	"fmt"
	"log/slog"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rothskeller/packet/cmd/packet/osdep"
	"github.com/rothskeller/packet/cmd/packet/server"
	"github.com/rothskeller/packet/form"
	"github.com/rothskeller/packet/message"
)

var viewCmd = &cobra.Command{
	Use:   "view msgfile [msgid opcall opname opdt]",
	Short: "Display a message in PDF format",
	Long: `Command called by Outpost to display a message.  The required parameter is:

msgfile  Filename of the file containing the message

For a received message, additional parameters can be specified (all or none):

msgid    Local message ID of the message
opcall   FCC call sign of the local operator
opname   Name of the local operator
opdt     Date and time that the message was received`,
	DisableFlagsInUseLine: true,
	SilenceUsage:          true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 && len(args) != 5 {
			return fmt.Errorf("accepts 1 or 5 arguments, received %d", len(args))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			msgfile string
			msg     message.Message
			address string
			u       *url.URL
			browse  *exec.Cmd
			out     string
		)
		// Verify that the message file is reachable.
		msgfile = decodeArg(args[0])
		if msg, err = message.ReadNoHeader(msgfile); msg == nil {
			slog.Error("can't read message file", "f", msgfile, "err", err)
			return fmt.Errorf("%s: %s", msgfile, err)
		}
		if len(args) == 5 {
			if ft, ok := msg.Type().(form.FormType); ok {
				body := msg.Body().(*form.FormBody)
				body.SetField("RECEIVED", "RECEIVED")
				for fd := range ft.AllFields() {
					switch fd.Common {
					case "destinationMessageID":
						body.SetField(fd.Tag, decodeArg(args[1]))
					case "operatorCall":
						body.SetField(fd.Tag, decodeArg(args[2]))
					case "operatorName":
						body.SetField(fd.Tag, decodeArg(args[3]))
					case "operatorDate":
						date, _, _ := strings.Cut(decodeArg(args[4]), " ")
						body.SetField(fd.Tag, date)
					case "operatorTime":
						_, time, _ := strings.Cut(decodeArg(args[4]), " ")
						body.SetField(fd.Tag, time)
					}
				}
			}
		}
		// Get the server address.  This also starts the server if not already
		// running.
		if address, err = server.GetAddress(true); err != nil {
			return fmt.Errorf("starting server: %s", err)
		}
		// Create a temp file for the PDF.
		if out, err = server.CreateTempPDF(msg); err != nil {
			return err
		}
		// Build the request URL.
		if u, err = url.Parse(address); err != nil {
			slog.Error("url.Parse", "url", address, "err", err)
			return fmt.Errorf("parse server address %q: %s", address, err)
		}
		u.Path = "/pdf/" + filepath.Base(out)
		// Open that URL in a browser.
		browse = osdep.OpenURLCommand(u.String())
		if err = browse.Run(); err != nil {
			slog.Error("cmd.Run", "url", u, "err", err)
			return fmt.Errorf("open PDF in browser: %s", err)
		}
		slog.Info("opened browser to view message")
		return nil
	},
}

func init() {
	Command.AddCommand(viewCmd)
}
