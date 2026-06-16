package cmd

import (
	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/rothskeller/packet/v4/incident"
	"github.com/spf13/pflag"
)

const (
	listSlug = `List messages and log entries`
	listHelp = `
usage: packet list [flags]
  -f, --full      ⇥Use full 6-column ICS-309 layout
  -n, --numbers   ⇥Include log entry numbers
  -r, --receipts  ⇥Include receipt messages

The "packet list" (or "l", "ls", or "log") command lists messages and log entries in the incident.  Sent or received receipt messages are omitted unless --receipts (or -r) is given.

When standard output is a terminal, the log can be emitted in either of two formats.  When the --full (or -f) flag is given, the log is displayed in the county-standard ICS-309 tabular format.  Note that this format is too wide to be useful in most terminal windows.  Without that flag, the log is displayed in a more compact five-column format with "Time", "From", "Msg ID", "To", and "Message" columns.  In either case, if the --numbers (or -n) flag is given, an additional left-hand column displays the log entry number for each entry.  This is useful when issuing commands to edit those log entries.

Five markers can appear in the displayed log entries:
  - ⇥DRAFT indicates an unsent message that is not ready to send.
  - ⇥READY indicates an unsent message that is ready to send.
  - ⇥NO RCPT indicates a sent message for which we have not received a receipt.
  - ⇥NEW indicates a received message that has not been read.
  - ⇥VOICE indicates a log entry for a voice message, to be emitted on a separate voice ICS-309 form.
DRAFT and READY appear in the TIME column, NO RCPT appears in the TO MSG # column (6-column layout) or FROM column (5-column layout), and NEW and VOICE appear in the MESSAGE column.  Log entries can also have a "Needs Followup" marker.  This appears as a red star in the right margin.

Log entries are color-coded in the list, when displayed on a capable terminal.  Bulletins are in cyan, immediate messages are in red, priority messages are in yellow, and receipts (if shown) are in grey.

When standard output is not a terminal, the log is emitted in CSV format with 9 columns: log entry number (regardless of whether --numbers was given), flags (any of those listed above), the six standard ICS-309 columns, and the needs followup column (either empty or "X").
`
)

func cmdList(args []string) (err error) {
	var full, numbers, receipts bool

	flags := pflag.NewFlagSet("list", pflag.ContinueOnError)
	flags.BoolVarP(&full, "full", "f", false, "Full ICS-309 layout")
	flags.BoolVarP(&numbers, "numbers", "n", false, "Include log entry numbers")
	flags.BoolVarP(&receipts, "receipts", "r", false, "Include log entries for receipt messages")
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"list"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(listHelp)
	}
	if flags.NArg() != 0 {
		return usage(listHelp)
	}
	return incRead(func(i *incident.Incident) error {
		cio.Open().EmitLogList(i.Log, full, numbers, receipts, i.Config.ActiveCall())
		return nil
	})
}
