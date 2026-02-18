package cmd

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/rothskeller/packet/cmd/packet/cmd/cmdutil"
	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/spf13/cobra"
)

var pdfCmd = &cobra.Command{
	Use:                   "pdf {«msgID» | #«logentrynumber»}",
	Short:                 "Open a message in PDF form",
	Long:                  `Opens the system PDF viewer showing the identified message.  (If the argument is a log entry number, the message described by that log entry is shown.)`,
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	SilenceUsage:          true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			dir   string
			msg   message.Message
			entry *incident.LogEntry
			pdf   string
		)
		if dir, err = os.Getwd(); err != nil {
			return err
		}
		registerForms()
		if err = incident.Read(dir, func(i *incident.Incident) error {
			msg, entry, err = cmdutil.MatchMessage(i, args[0], cmdutil.MMMessageOnly)
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
		return cmdutil.ShowPDF(pdf)
	},
}

func init() {
	RootCmd.AddCommand(pdfCmd)
}
