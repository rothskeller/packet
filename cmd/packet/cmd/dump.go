package cmd

import (
	"io"
	"os"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/spf13/pflag"
)

const (
	dumpSlug = `Show a message in encoded form`
	dumpHelp = `
usage: packet dump «message-id»

The "dump" command displays a message in its PackItForms- and RFC-5322-encoded format, as it would be transmitted over the air.
`
)

var errNoChange = errors.New("don't write incident: no change was made")

func cmdDump(args []string) (err error) {
	var (
		msg message.Message
	)
	flags := pflag.NewFlagSet("dump", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"dump"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(dumpHelp)
	}
	if len(args) != 1 {
		return usage(dumpHelp)
	}
	registerForms()
	if err = incWrite(false, func(i *incident.Incident) error {
		var le *incident.LogEntry

		if msg, le, err = matchMessage(i, args[0], MMMessageOnly); err != nil {
			return err
		} else if le.Flags&incident.FUnread != 0 {
			le.Flags &^= incident.FUnread
			return nil
		} else {
			return errNoChange
		}
	}); err != nil && err != errNoChange {
		return err
	}
	io.WriteString(os.Stdout, msg.RFC5322())
	return nil
}
