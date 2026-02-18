package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/cmd/packet/cmd/cmdutil"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [-m] {«msgID» | #«logentrynumber» | configuration} [«field»]",
	Short: "Display a message, log entry, or incident configuration",
	Long: `Displays a message, a log entry, or the incident configuration in the console in a two-column field/value table.  If the first argument is a message ID, the message is shown.  If it is a log entry, that log entry is shown unless the --message (or -m) flag is given, in which case the message described by that log entry is shown.  If the first argument is "configuration" (or any abbreviation), the incident configuration is shown.

If a second argument is given, it identifies a field of the item, and only that field is displayed.  The second argument can be either a field name or, in a forms message, the PackItForms tag for a field.  When standard output is a terminal, it can be a shortened version of the field name, such as "ocs" for "Operator Call Sign."`,
	Args:                  cobra.RangeArgs(1, 2),
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
			var flags cmdutil.MatchMessageFlag
			if message, _ := cmd.Flags().GetBool("message"); message {
				flags |= cmdutil.MMMessageOnly
			}
			msg, _, err = cmdutil.MatchMessage(i, args[0], flags)
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
	showCmd.Flags().BoolP("message", "m", false, "Show message rather than log entry")
	RootCmd.AddCommand(showCmd)
}

func showSingleField(msg message.Message, fieldname string) (err error) {
	c := cio.Open()
	field, err := cmdutil.MatchFieldName(msg, fieldname, c.OutputIsTerm)
	if err != nil {
		return err
	}
	value := field.ToHuman(field.Value(msg))
	if strings.HasSuffix(value, "\n") {
		io.WriteString(os.Stdout, value)
	} else if value != "" {
		fmt.Println(value)
	}
	return nil
}

func showMessage(msg message.Message) {
	nv := cio.Open().NewNameValueList()
	for f := range msg.Fields() {
		if !f.Visible(msg) {
			continue
		}
		value := f.ToHuman(f.Value(msg))
		if value == "" {
			continue
		}
		nv.ShowNVPair(f.Label(), value)
	}
	nv.Close()
}
