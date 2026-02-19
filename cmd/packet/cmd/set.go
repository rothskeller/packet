package cmd

import (
	"os"

	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/spf13/cobra"
)

const (
	setSlug = `Set a value in a message, log entry, or incident configuration`
	setHelp = `
usage: packet set ⇥[flags] { «message-id» | «log-entry» | configuration } «field» [«value»]
  -f, --force    ⇥Sets the value even if it is invalid for the field
  -m, --message  ⇥Sets the value in the message, not the log entry

The "packet set" command sets the value of the named field in a message, a log entry, or the incident configuration.  If the first argument identifies a message, it must be for an unsent outgoing message, and the named field of that message is changed.  If it identifies a log entry, that log entry is changed.  If the first argument is "configuration" (or any abbreviation), the incident configuration is changed.

The second argument identifies the field of the item to be changed.

If a third argument is given, it is taken as the new value of the field; otherwise, the new value is read from the console or standard input.  It must be a valid value for the field unless the --force (or -f) flag is given.
`
)

func cmdSet(args []string) (err error) {
	panic("not implemented")
}

var setCmd = &cobra.Command{
	Use:   "set [-fm] {«msgID» | #«logentry» | configuration} «field» [«value»]",
	Short: "Set a value in a message, log entry, or incident configuration",
	Long: `Sets the value of the named field in a message, a log entry, or the incident configuration.  If the first argument is a message ID, it must be for an unsent outgoing message, and the named field of that message is changed.  If it is a log entry, that log entry is changed unless the --message (or -m) flag is given, in which case the message described by that log entry is changed.  If the first argument is "configuration" (or any abbreviation), the incident configuration is changed.

The second argument identifies the field of the item to be changed.  It can be either a field name or, in a forms message, the PackItForms tag for a field.  When standard output is a terminal, it can be a shortened version of a field name, such as "ocs" for "Operator Call Sign."

If a third argument is given, it is taken as the new value of the field; otherwise, the new value is read from the console or standard input.  It must be a valid value for the field unless the --force (or -f) flag is given.`,
	Args:                  cobra.RangeArgs(2, 3),
	DisableFlagsInUseLine: true,
	SilenceUsage:          true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			dir string
			msg message.Message
		)
		if dir, err = os.Getwd(); err != nil {
			return err
		}
		registerForms()
		if err = incident.Read(dir, func(i *incident.Incident) error {
			var flags matchMessageFlag
			if message, _ := cmd.Flags().GetBool("message"); message {
				flags |= MMMessageOnly
			}
			msg, _, err = matchMessage(i, args[0], flags)
			return err
		}); err != nil {
			return err
		}
		if len(args) == 2 {
			return showSingleField(msg, args[1])
		} else {
			showMessage(msg)
		}
		return nil
	},
}

func init() {
	setCmd.Flags().BoolP("force", "f", false, "Accept invalid value for the field")
	setCmd.Flags().BoolP("message", "m", false, "Change message rather than log entry")
}
