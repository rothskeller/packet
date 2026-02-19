package cmd

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/cmd/packet/osdep"
	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/spf13/pflag"
)

const (
	pdfSlug = `Open a message in PDF form`
	pdfHelp = `
usage: packet pdf «message-id»

The "packet pdf" command opens the system PDF viewer showing the identified message.
`
)

func cmdPDF(args []string) (err error) {
	var (
		dir   string
		msg   message.Message
		entry *incident.LogEntry
		pdf   string
	)
	flags := pflag.NewFlagSet("pdf", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"pdf"})
	} else if err != nil {
		cio.Open().Error("%s", err.Error())
		return usage(pdfHelp)
	}
	if len(args) != 1 {
		return usage(pdfHelp)
	}
	registerForms()
	if dir, err = os.Getwd(); err != nil {
		return err
	}
	if err = incident.Read(dir, func(i *incident.Incident) error {
		msg, entry, err = matchMessage(i, args[0], MMMessageOnly)
		return err
	}); err != nil {
		return err
	}
	pdf = filepath.Join(dir, incident.ToPDF(entry.Filename()))
	// It's possible that the PDF doesn't exist yet.  If so we need to
	// create it.
	if _, err = os.Stat(pdf); os.IsNotExist(err) {
		if err = msg.Type().RenderPDF(msg, pdf, ""); err != nil {
			slog.Error("RenderPDF", "dir", dir, "id", entry.Ident, "err", err)
			return errors.NewF("Unable to create PDF: %s", err)
		}
		slog.Debug("Rendered missing PDF", "f", pdf)
	} else if err != nil {
		slog.Error("os.Stat", "f", pdf, "err", err)
		return errors.NewF("Unable to read PDF: %s", err)
	}
	var showcmd = osdep.OpenFileCommand(pdf)
	if err := showcmd.Start(); err != nil {
		return errors.NewF("Unable to start PDF viewer: %s", err)
	}
	go func() { showcmd.Wait() }()
	return nil
}
