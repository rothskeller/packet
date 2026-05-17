package cmd

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/cmd/packet/osdep"
	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/spf13/pflag"
)

const (
	printSlug = `Sends a message to the system printer`
	printHelp = `
usage: packet print «message-id»

The "packet print" command sends the identified message to the system default printer.
`
)

func cmdPrint(args []string) (err error) {
	var (
		msg message.Message
		le  *incident.LogEntry
		dir string
		pdf string
		cmd *exec.Cmd
	)
	flags := pflag.NewFlagSet("print", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"print"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(printHelp)
	}
	if len(args) != 1 {
		return usage(printHelp)
	}
	registerForms()
	if err = incWrite(false, func(i *incident.Incident) error {
		dir = i.Dir
		if msg, le, err = matchMessage(i, args[0], MMMessageOnly); err != nil {
			return err
		}
		pdf = incident.ToPDF(le.Filename())
		if le.Flags&incident.FUnread != 0 {
			le.Flags &^= incident.FUnread
			return nil
		} else {
			return errNoChange
		}
	}); err != nil && err != errNoChange {
		return err
	}
	// It's possible that the PDF doesn't exist yet.  If so we need to
	// create it.
	pdf = filepath.Join(dir, pdf)
	if _, err := os.Stat(pdf); os.IsNotExist(err) {
		if err = msg.Type().RenderPDF(msg, pdf, ""); errors.IsType[message.Warning](err) {
			cio.Open().Warn("pdf rendering issue: %s", err)
			slog.Warn("RenderPDF", "dir", dir, "id", le.Ident, "warn", err)
		} else if err != nil {
			slog.Error("RenderPDF", "dir", dir, "id", le.Ident, "err", err)
			return errors.NewF("Unable to create PDF: %s", err)
		}
	} else if err != nil {
		return err
	}
	if cmd = osdep.PrintPDFCommand(pdf); cmd == nil {
		slog.Error("server print not supported")
		return errors.New("No command was found on this system that can print PDF files.")
	}
	if err = cmd.Run(); err == nil {
		return errors.NewF("The print command failed: %s", errNoChange)
	}
	return nil
}
