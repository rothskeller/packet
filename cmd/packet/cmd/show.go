package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/cmd/packet/cmd/cmdutil"
	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/messageid"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [-m] {message|logentry|configuration} [field]",
	Short: "Display a message, log entry, or incident configuration",
	Long: `Displays a message, a log entry, or the incident configuration in the console in a two-column field/value table.  The first argument may be:
  - ⇥A local message ID
  - ⇥A remote message ID, if it is unambiguous
  - ⇥The numeric part of a local or remote message ID, if it is unambiguous
  - ⇥A pound sign followed by a log entry number (e.g., "#23")
  - ⇥The word "configuration" or any abbreviation of it
The first three cases will display the fields of the identified message.  A log entry number will display the fields of the log entry.  The word "configuration" will display the fields of the incident configuration.

The --message (-m) flag is relevant only when a log entry number is given.  With this flag, the command will display the message described by the log entry rather than the log entry itself.

If a second argument is given, it identifies a field of the item, and only that field is displayed.  The second argument can be either a field name or, in a forms message, the PackItForms tag for a field.  When standard output is a terminal, it can be a shortened version of the field name, such as "ocs" for "Operator Call Sign."`,
	Args:                  cobra.RangeArgs(1, 2),
	DisableFlagsInUseLine: true,
	SilenceUsage:          true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			dir     string
			item    message.Message
			showmsg bool
		)
		if dir, err = os.Getwd(); err != nil {
			return err
		}
		showmsg, _ = cmd.Flags().GetBool("message")
		if strings.HasPrefix("configuration", args[0]) {
			item, err = configMessage(dir)
		} else if strings.HasPrefix(args[0], "#") {
			var num int
			if num, err = strconv.Atoi(args[0][1:]); err != nil || num < 1 {
				return errors.NewF("%q is not a valid log entry number", args[0])
			} else {
				item, err = logEntryMessage(dir, num, showmsg)
			}
		} else {
			item, err = getMessage(dir, args[0])
		}
		if err != nil {
			return err
		}
		if len(args) == 2 {
			return showSingleField(item, args[1])
		} else {
			showMessage(item)
		}
		return nil
	},
}

func init() {
	showCmd.Flags().BoolP("message", "m", false, "Show message rather than log entry")
	RootCmd.AddCommand(showCmd)
}

func configMessage(dir string) (msg message.Message, err error) {
	var config *incident.Config

	if err = incident.Read(dir, func(i *incident.Incident) error {
		config = i.Config
		return nil
	}); err != nil {
		return nil, err
	}
	return &cmdutil.ConfigMessage{Config: config}, nil
}

func logEntryMessage(dir string, num int, showmsg bool) (msg message.Message, err error) {
	var entry *incident.LogEntry

	if err = incident.Read(dir, func(i *incident.Incident) error {
		for _, e := range i.Log {
			if e.Ident == num {
				entry = e
				break
			}
		}
		if entry != nil && showmsg {
			msg, err = i.GetMessageFromLogEntry(entry)
		}
		return err
	}); err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, errors.NewF("There is no log entry #%d.", num)
	} else if entry.Status == incident.StatusDeleted {
		return nil, errors.NewF("The log entry #%d has been deleted and is not recoverable.", num)
	} else if msg != nil {
		return msg, nil
	}
	return &cmdutil.LogEntryMessage{LogEntry: entry}, nil
}

func getMessage(dir, msgid string) (msg message.Message, err error) {
	registerForms()
	num, _ := strconv.Atoi(msgid)
	err = incident.Read(dir, func(i *incident.Incident) error {
		var entry *incident.LogEntry
		var ambiguous bool

		for _, e := range i.Log {
			var local string

			if e.Status == incident.StatusDeleted || e.Status == incident.StatusHandEntered {
				continue
			}
			if local = e.LocalMsgID; e.Flags&incident.FIsReceipt != 0 {
				local = ""
			}
			if strings.EqualFold(local, msgid) {
				entry, ambiguous = e, false
				break
			}
			if strings.EqualFold(e.FromMsgID, msgid) || strings.EqualFold(e.ToMsgID, msgid) {
				ambiguous = ambiguous || entry != nil
				entry = e
				continue
			}
			if num != 0 {
				if _, n, _, err := messageid.Decode(e.FromMsgID, true, false); err == nil && n == num {
					ambiguous = ambiguous || (entry != nil && entry.LocalMsgID != local)
					entry = e
				} else if _, n, _, err = messageid.Decode(local, true, false); err == nil && n == num {
					ambiguous = ambiguous || (entry != nil && entry.LocalMsgID != local)
					entry = e
				} else if _, n, _, err = messageid.Decode(e.ToMsgID, true, false); err == nil && n == num {
					ambiguous = ambiguous || (entry != nil && entry.LocalMsgID != local)
					entry = e
				}
			}
		}
		if entry == nil {
			return errors.NewF("There is no message %q.", msgid)
		} else if ambiguous {
			return errors.NewF("The string %q is ambiguous: it identifies multiple messages.", msgid)
		}
		msg, err = i.GetMessageFromLogEntry(entry)
		return err
	})
	return msg, err
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
