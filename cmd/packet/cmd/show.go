package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/cmd/packet/pseudomsg"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/spf13/pflag"
)

const (
	showSlug = `Display a message, log entry, or incident configuration`
	showHelp = `
usage: packet show ⇥[-m] { «message-id» | «log-entry» | configuration} [«field»]
  -m, --message  ⇥Sets the value in the message, not the log entry

The "packet show" command displays a message, a log entry, or the incident configuration in the console in a two-column field/value table.  If the first argument is "configuration" (or any abbreviation), the incident configuration is shown.

If a second argument is given, it identifies a field of the item, and only that field is displayed.
`
)

func cmdShow(args []string) (err error) {
	var (
		asMessage bool
		msg       message.Message
	)
	flags := pflag.NewFlagSet("show", pflag.ContinueOnError)
	flags.BoolVarP(&asMessage, "message", "m", false, "Show message rather than log entry")
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"show"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(showHelp)
	}
	if n := flags.NArg(); n < 1 || n > 2 {
		return usage(showHelp)
	}
	registerForms()
	if err = incWrite(false, func(i *incident.Incident) error {
		var (
			flags matchMessageFlag
			entry *incident.LogEntry
		)
		if asMessage {
			flags |= MMMessageOnly
		}
		if msg, entry, err = matchMessage(i, args[0], flags); err != nil {
			return err
		} else if _, ok := msg.(*pseudomsg.LogEntryMessage); !ok && entry != nil && entry.Flags&incident.FUnread != 0 {
			entry.Flags &^= incident.FUnread
			return nil
		} else {
			return errNoChange
		}
	}); err != nil && err != errNoChange {
		return err
	}
	if len(args) == 2 {
		return showSingleField(msg, args[1])
	} else {
		showMessage(msg)
	}
	return nil
}

func showSingleField(msg message.Message, fieldname string) (err error) {
	c := cio.Open()
	field, err := matchField(msg, fieldname, c.OutputIsTerm)
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
