package outpost

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/field"
)

var convertCmd = &cobra.Command{
	Use:   "convert msgfile state msgid subject opcall opname opdt spool copies",
	Short: "Convert a message to PDF for printing",
	Long: `
Command called by Outpost to convert a message to PDF for printing.  The
parameters (all required) are:

msgfile  Filename of the file containing the message
state    State of the message (new, ready, sent, unread, read, or draft)
msgid    Local message ID of the message
opcall   FCC call sign of the local operator
opname   Name of the local operator
opdt     Date and time that the message was received
spool    Directory into which to place the PDFs
copies   Name(s) to apply to the copies (separated by newlines)

This will make one PDF in the spool directory for each copy name (or one total,
if copies is empty).  The PDFs will be named with the subject of the message
and a sequence number.`,
	DisableFlagsInUseLine: true,
	SilenceUsage:          true,
	Args:                  cobra.ExactArgs(8),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			msgfile string
			msg     message.Message
			base    string
			seq     int
			spool   string
			ents    []os.DirEntry
		)
		// Read the message file.  Note that ReadNoHeader always returns
		// a DraftMessage.
		msgfile = decodeArg(args[0])
		if msg, err = message.ReadNoHeader(msgfile); msg == nil {
			slog.Error("can't read message file", "f", msgfile, "err", err)
			return fmt.Errorf("%s: %s", msgfile, err)
		}
		// If this is a received message, fill in the local operator
		// info and the destination message ID.
		switch args[1] {
		case "unread", "read", "received", "retrieved":
			for fd := range msg.Fields() {
				switch fd.Common() {
				case field.CDestinationMessageID:
					fd.SetValue(msg, decodeArg(args[2]))
				case field.COperatorCall:
					fd.SetValue(msg, decodeArg(args[3]))
				case field.COperatorName:
					fd.SetValue(msg, decodeArg(args[4]))
				case field.COperatorDate:
					date, _, _ := strings.Cut(decodeArg(args[5]), " ")
					fd.SetValue(msg, date)
				case field.COperatorTime:
					_, time, _ := strings.Cut(decodeArg(args[5]), " ")
					fd.SetValue(msg, time)
				}
			}
		}
		// Compute the base filename based on the subject.
		base = msg.Subject().EncodedSubject()
		base = strings.Map(func(r rune) rune {
			if strings.ContainsRune(`<>:"/\|?*^`, r) {
				return '~'
			}
			return r
		}, base)
		base += "^"
		// Find the highest-numbered file in the spool directory with
		// that basename.
		spool = decodeArg(args[6])
		if ents, err = os.ReadDir(spool); err != nil {
			slog.Error("os.ReadDir", "d", spool, "err", err)
			return fmt.Errorf("spool dir %s: %s", spool, err)
		}
		for _, ent := range ents {
			if !strings.HasPrefix(ent.Name(), base) || !strings.HasSuffix(ent.Name(), ".pdf") {
				continue
			}
			if n, err := strconv.Atoi(strings.TrimSuffix(ent.Name()[len(base):], ".pdf")); err != nil {
				continue
			} else if n >= seq {
				seq = n + 1
			}
		}
		// Create one PDF for each named copy.
		base = filepath.Join(spool, base)
		for copy := range strings.SplitSeq(decodeArg(args[7]), "\n") {
			if err = convertOne(msg, fmt.Sprintf("%s%03d.pdf", base, seq), copy); err != nil {
				return err
			}
			seq++
		}
		return nil
	},
}

func init() {
	Command.AddCommand(convertCmd)
}

func convertOne(msg message.Message, filename, copyname string) (err error) {
	err = msg.Type().RenderPDF(msg, filename, copyname)
	if _, ok := err.(message.Warning); ok {
		slog.Warn("RenderPDF", "warn", err)
	} else if err != nil {
		slog.Error("RenderPDF", "err", err)
		return fmt.Errorf("unable to create PDF: %s", err)
	}
	slog.Info("converted message to PDF", "f", filename, "copy", copyname, "s", msg.Subject().EncodedSubject())
	return nil
}
