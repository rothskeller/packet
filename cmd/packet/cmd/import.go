package cmd

import (
	"path/filepath"

	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/rothskeller/packet/v4/incident"
	"github.com/rothskeller/packet/v4/message"
	"github.com/spf13/pflag"
)

const (
	importSlug = `Import a message text file into the incident`
	importHelp = `
usage: packet import «text-file»

The "import" reads a text file containing a packet message and adds it to the incident as if it had been sent or received through the packet software.  The message is imported as a received message if it has a Received header; otherwise as a sent message if it has a Date header; otherwise as a draft message.
`
)

func cmdImport(args []string) (err error) {
	var msg message.Message

	flags := pflag.NewFlagSet("import", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"import"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(importHelp)
	}
	if len(args) != 1 {
		return usage(importHelp)
	}
	registerForms()
	if msg, err = message.Read(args[0]); msg == nil {
		return err
	} else if err != nil {
		cio.Open().Warn("%s", err)
	}
	if err = incWrite(false, func(i *incident.Incident) error {
		switch msg := msg.(type) {
		case *message.DraftMessage:
			_, err = i.AddDraftMessage(msg)
		case *message.SentMessage:
			err = i.AddSentMessage(msg, true)
		case *message.ReceivedMessage:
			overwriteOK := msg.LocalID() != "" && filepath.Base(args[0]) == msg.LocalID()+".txt"
			err = i.AddReceivedMessage(msg, overwriteOK)
		}
		return err
	}); err != nil && err != errNoChange {
		return err
	}
	return nil
}
