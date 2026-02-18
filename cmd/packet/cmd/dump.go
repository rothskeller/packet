package cmd

import (
	"io"
	"os"

	"github.com/rothskeller/packet/cmd/packet/cmd/cmdutil"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:                   "dump {«msgID» | #«logentrynumber»}",
	Short:                 "Display a message in encoded form",
	Long:                  `Displays a message in its encoded form.  (If the argument is a log entry number, the message described by that log entry is displayed.)`,
	Args:                  cobra.ExactArgs(1),
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
			msg, _, err = cmdutil.MatchMessage(i, args[0], cmdutil.MMMessageOnly)
			return err
		}); err != nil {
			return err
		}
		io.WriteString(os.Stdout, msg.RFC5322())
		return nil
	},
}

func init() {
	RootCmd.AddCommand(dumpCmd)
}
