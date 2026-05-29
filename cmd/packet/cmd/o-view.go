//go:build windows

package cmd

import (
	"fmt"
	"log/slog"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/cmd/packet/osdep"
	"github.com/rothskeller/packet/cmd/packet/server"
	"github.com/rothskeller/packet/form"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/field"
	"github.com/spf13/pflag"
)

const (
	outpostViewSlug = `Display a message in PDF format`
	outpostViewHelp = `
usage: packet outpost view view msgfile [msgid opcall opname opdt]

Command called by Outpost to display a message.  The required parameter is:

msgfile  Filename of the file containing the message

For a received message, additional parameters can be specified (all or none):

msgid    ⇥Local message ID of the message
opcall   ⇥FCC call sign of the local operator
opname   ⇥Name of the local operator
opdt     ⇥Date and time that the message was received
`
)

func cmdOutpostView(args []string) (err error) {
	var (
		msgfile string
		msg     message.Message
		address string
		u       *url.URL
		browse  *exec.Cmd
		out     string
	)
	flags := pflag.NewFlagSet("o-view", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdOutpostHelp([]string{"view"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(outpostViewHelp)
	}
	if len(args) != 1 && len(args) != 5 {
		return usage(outpostViewHelp)
	}
	registerForms()
	// Read the message file.  Note that since it has no headers,
	// it will be returned as a DraftMessage.
	msgfile = decodeOutpostArg(args[0])
	if msg, err = message.ReadNoHeader(msgfile); msg == nil {
		slog.Error("can't read message file", "f", msgfile, "err", err)
		return fmt.Errorf("%s: %s", msgfile, err)
	}
	if len(args) == 5 {
		// This is a received message.  Fill in the fields as
		// given.  (We're still leaving the type as
		// DraftMessage, though, because we don't have the
		// details to put into a ReceivedMessage.)
		var ft *form.FormType
		switch t := msg.Type().(type) {
		case *form.FormType:
			ft = t
		case form.EditableFormType:
			ft = t.FormType
		}
		if ft != nil {
			body := msg.Body().(*form.FormBody)
			body.SetField("RECEIVED", "RECEIVED")
			for fd := range ft.AllFields() {
				switch fd.Common {
				case field.CDestinationMessageID:
					body.SetField(fd.Tag, decodeOutpostArg(args[1]))
				case field.COperatorCall:
					body.SetField(fd.Tag, decodeOutpostArg(args[2]))
				case field.COperatorName:
					body.SetField(fd.Tag, decodeOutpostArg(args[3]))
				case field.COperatorDate:
					date, _, _ := strings.Cut(decodeOutpostArg(args[4]), " ")
					body.SetField(fd.Tag, date)
				case field.COperatorTime:
					_, time, _ := strings.Cut(decodeOutpostArg(args[4]), " ")
					body.SetField(fd.Tag, time)
				case field.CReceiverSender:
					body.SetField(fd.Tag, "receiver")
				case field.COperatorMethod:
					body.SetField(fd.Tag, "Other")
				case field.COperatorMethodOther:
					body.SetField(fd.Tag, "Packet")
				}
			}
		}
	}
	// Get the server address.  This also starts the server if not already
	// running.
	if address, err = server.OutpostGetAddress(); err != nil {
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
}
