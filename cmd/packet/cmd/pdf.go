package cmd

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/rothskeller/packet/v4/cmd/packet/osdep"
	"github.com/rothskeller/packet/v4/errors"
	"github.com/rothskeller/packet/v4/incident"
	"github.com/rothskeller/packet/v4/message"
	"github.com/spf13/pflag"
)

const (
	pdfSlug = `Open a message in PDF form`
	pdfHelp = `
usage: packet pdf [flags] «message-id»
  -p, --practice  ⇥Practice message (no IDs or op info)
  -r, --rebuild   ⇥Rebuild the cached PDF

The "packet pdf" command opens the system PDF viewer showing the identified message.  With the --practice (or -p) flag, the message IDs and operator information can be left blank (useful for printing messages to be sent in exercises).
`
)

func cmdPDF(args []string) (err error) {
	var (
		msg      message.Message
		entry    *incident.LogEntry
		pdf      string
		practice bool
		rebuild  bool
	)
	flags := pflag.NewFlagSet("pdf", pflag.ContinueOnError)
	flags.BoolVarP(&practice, "practice", "p", false, "practice message (no IDs or op info)")
	flags.BoolVarP(&rebuild, "rebuild", "r", false, "rebuild the cached PDF")
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"pdf"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(pdfHelp)
	}
	if flags.NArg() != 1 {
		return usage(pdfHelp)
	}
	registerForms()
	if err = incWrite(false, func(i *incident.Incident) error {
		if msg, entry, err = matchMessage(i, flags.Arg(0), MMMessageOnly); err != nil {
			return err
		}
		pdf = filepath.Join(i.Dir, incident.ToPDF(entry.Filename()))
		if entry.Flags&incident.FUnread != 0 {
			entry.Flags &^= incident.FUnread
			return nil
		} else {
			return errNoChange
		}
	}); err != nil && err != errNoChange {
		return err
	}
	// It's possible that the PDF doesn't exist yet.  If so we need to
	// create it.
	if _, err = os.Stat(pdf); os.IsNotExist(err) || rebuild {
		if practice {
			nm := msg.Type().(message.EditableMType).NewDraft().(*message.DraftMessage)
			message.CopyFields(msg, nm)
			msg = nm
		}
		if err = msg.Type().RenderPDF(msg, pdf, ""); errors.IsType[message.Warning](err) {
			cio.Open().Warn("pdf rendering issue: %s", err)
			slog.Warn("RenderPDF", "f", pdf, "warn", err)
		} else if err != nil {
			slog.Error("RenderPDF", "f", pdf, "err", err)
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
