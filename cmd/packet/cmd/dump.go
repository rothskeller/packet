package cmd

import (
	"io"
	"os"

	"github.com/rothskeller/packet/cmd/packet/cio"
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

func cmdDump(args []string) (err error) {
	var (
		msg    message.Message
		dir, _ = os.Getwd()
	)
	flags := pflag.NewFlagSet("dump", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"dump"})
	} else if err != nil {
		cio.Open().Error("%s", err.Error())
		return usage(dumpHelp)
	}
	if len(args) != 1 {
		return usage(dumpHelp)
	}
	registerForms()
	if err = incident.Read(dir, func(i *incident.Incident) error {
		msg, _, err = matchMessage(i, args[0], MMMessageOnly)
		return err
	}); err != nil {
		return err
	}
	io.WriteString(os.Stdout, msg.RFC5322())
	return nil
}
