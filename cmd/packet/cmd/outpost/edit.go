package outpost

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit msgfile msgindex",
	Short: "Start an editor on an existing unsent message",
	Long: `Command called by Outpost to start editing an existing (unsent) message. The parameters are:

msgfile   ⇥Filename of the file containing the existing message
msgindex  ⇥Index of the message in the Outpost database`,
	DisableFlagsInUseLine: true,
	SilenceUsage:          true,
	Args:                  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			msgfile string
			values  = make(url.Values)
		)
		// Verify that the message file is reachable.
		msgfile = decodeArg(args[0])
		if _, err := os.Stat(msgfile); err != nil {
			slog.Error("os.Stat", "f", msgfile, "err", err)
			return fmt.Errorf("%s: %s", msgfile, err)
		}
		values.Set("msgfile", msgfile)
		values.Set("index", decodeArg(args[1]))
		if err := neCommon(values, "/outpost-edit"); err != nil {
			return err
		}
		slog.Info("opened browser to edit message", "f", msgfile)
		return nil
	},
}

func init() {
	Command.AddCommand(editCmd)
}

var argEncodeRE = regexp.MustCompile(`(?i)~[0-9a-f][0-9a-f]`)

// decodeArg reverses Outpost's argument encoding.
func decodeArg(s string) string {
	return argEncodeRE.ReplaceAllStringFunc(s, func(e string) string {
		code, _ := strconv.ParseInt(e[1:], 16, 64)
		return string(byte(code))
	})
}
