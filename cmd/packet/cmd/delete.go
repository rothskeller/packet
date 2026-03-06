package cmd

import (
	"fmt"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/spf13/pflag"
)

const (
	deleteSlug = `Delete an unsent message or manual log entry`
	deleteHelp = `
usage: packet delete { «message-id» | «log-entry» }

Given a message ID or log entry of an outgoing message that has not been sent, the "delete" command deletes that message.  For safety, the usual shorthands for message IDs are not accepted by the "delete" command; it must be spelled out fully.  Given a manual log entry, the "delete" command deletes that log entry.  These deletions are not reversible.`
)

func cmdDelete(args []string) (err error) {
	var confirm string

	flags := pflag.NewFlagSet("delete", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"delete"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(deleteHelp)
	}
	if len(args) != 1 {
		return usage(deleteHelp)
	}
	if err := incWrite(false, func(i *incident.Incident) error {
		var msg message.Message
		var le *incident.LogEntry
		if msg, le, err = matchMessage(i, args[0], MMNoAbbrev|MMNoRemote); err != nil {
			return err
		}
		if le.Status == incident.StatusHandEntered {
			if err = i.DeleteLogEntry(le); err != nil {
				return err
			}
			confirm = fmt.Sprintf("Log entry #%d deleted.", le.Ident)
		} else if _, ok := msg.(*message.DraftMessage); ok {
			if err = i.DeleteMessage(le); err != nil {
				return err
			}
			confirm = fmt.Sprintf("Draft message %s deleted.", le.LocalMsgID)
		} else if strings.HasPrefix(args[0], "#") {
			return errors.New("cannot delete a log entry for a real message")
		} else {
			return errors.New("cannot delete a message once it is sent or received")
		}
		return nil
	}); err != nil {
		return err
	}
	cio.Open().Confirm("%s", confirm)
	return nil
}
