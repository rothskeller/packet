package cmd

import (
	"os"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/incident"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list [-fnr]",
	Aliases: []string{"l", "li", "lis", "ls", "lo", "log"},
	Short:   "Lists messages and log entries",
	Long: `Displays a list of the messages and other log entries in the current incident.
When standard output is a terminal, they are displayed in tabular form.  By
default, they are displayed in a five-column format: time, from, local ID, to,
and summary.  The --full (-f) flag switches to the full ICS-309 six-column
format (too wide for most console windows).  The --numbers (-n) flag adds a
left-hand column giving the log entry numbers for each entry, making it easier
to address them with other commands.  Receipt messages are not shown unless the
--receipts (-r) flag is given.

When standard output is not a terminal, the output is a CSV file which always
includes all columns and all messages.  The above flags are ignored.

Five markers can appear in the displayed log entries:
  - DRAFT indicates an unsent message that is not ready to send.
  - READY indicates an unsent message that is ready to send.
  - NO RCPT indicates a sent message for which we have not received a receipt.
  - NEW indicates a received message that has not been read.
  - VOICE indicates a log entry for a voice message, to be emitted on a
    separate voice ICS-309 form.
In the tabular output, DRAFT and READY appear in the TIME column, NO RCPT
appears in the TO MSG # column (6-column layout) or FROM column (5-column
layout), and NEW and VOICE appear in the MESSAGE column.  In the CSV output,
all of these appear in a dedicated FLAG column.

Log entries can also have a "Needs Followup" marker.  This appears as a red
star in the right margin in tabular output, and as an "X" in the "NF" column in
CSV output.

Log entries are color-coded in the list, when displayed on a capable terminal.
Bulletins are in cyan, immediate messages are in red, priority messages are in
yellow, and receipts (if shown) are in grey.`,
	Args:                  cobra.NoArgs,
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if dir, err := os.Getwd(); err != nil {
			return err
		} else {
			return incident.Read(dir, func(i *incident.Incident) error {
				full, _ := cmd.Flags().GetBool("full")
				numbers, _ := cmd.Flags().GetBool("numbers")
				receipts, _ := cmd.Flags().GetBool("receipts")
				cio.Open().EmitLogList(i.Log, full, numbers, receipts)
				return nil
			})
		}
	},
}

func init() {
	listCmd.Flags().BoolP("full", "f", false, "Full ICS-309 layout")
	listCmd.Flags().BoolP("numbers", "n", false, "Include log entry numbers")
	listCmd.Flags().BoolP("receipts", "r", false, "Include log entries for receipt messages")
	RootCmd.AddCommand(listCmd)
}
