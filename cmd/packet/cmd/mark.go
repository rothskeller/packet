package cmd

import (
	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/msgifc"
	"github.com/spf13/pflag"
)

const (
	markSlug = `Change flags on a message or log entry`
	markHelp = `
usage: packet mark ⇥[-f] { «message-id» | «log-entry» } [not] «flag»...
       packet mark ⇥«message-id» sent
  -f, --force  ⇥Apply the change even if the message is invalid

The "packet mark" command sets flags on a message or log entry.  The flags are:
  - ⇥"delivered" (or "d") means that the message has been delivered to the recipient identified in this log entry (usually because we've received a delivery receipt for it).  The message must be a sent message.
  - ⇥"followup" (or "f") means that the log entry needs followup.
  - ⇥"ready" (or "r") means the message is ready to send at the next BBS connection.  The message must be an unsent outgoing message.
If the keyword "not" is used, the command clears the flags instead of setting them.

"packet mark ... ready" will not mark a message with validation errors ready unless the --force (or -f) flag is given.

"packet mark ... sent" is available only when the incident connection type is Manual, and is used to indicate that the operator has manually sent the identified message (usually after a "packet manual send" command to get the JNOS command with which to send it).  The message must be an unsent outgoing message, and will be marked as having been sent at the current time.  This is not reversible.
`
)

func cmdMark(args []string) (err error) {
	var (
		force bool
		f     byte
		not   bool
		msg   message.Message
		entry *incident.LogEntry
		c     = cio.Open()
	)
	flags := pflag.NewFlagSet("mark", pflag.ContinueOnError)
	flags.BoolVarP(&force, "force", "f", false, "Apply change even if message is invalid")
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"mark"})
	} else if err != nil {
		c.Error(err)
		return usage(markHelp)
	}
	if flags.NArg() < 2 || flags.NArg() > 3 {
		return usage(markHelp)
	}
	switch flags.Arg(1) {
	case "delivered", "d", "followup", "f", "ready", "r", "sent", "s":
		if flags.NArg() == 3 {
			usage(markHelp)
		}
		f = flags.Arg(1)[0]
	case "not":
		if flags.NArg() != 3 {
			usage(markHelp)
		}
		switch flags.Arg(2) {
		case "delivered", "d", "followup", "f", "ready", "r":
		default:
			usage(markHelp)
		}
		f, not = flags.Arg(2)[0], true
	}
	registerForms()
	if err = incWrite(false, func(i *incident.Incident) error {
		switch f {
		case 'r', 's':
			var dm *message.DraftMessage
			if msg, entry, err = matchMessage(i, flags.Arg(0), MMMessageOnly); err != nil {
				return err
			}
			if dm, _ = msg.(*message.DraftMessage); dm == nil {
				return errors.NewF("%q is not an unsent outgoing message.", flags.Arg(0))
			}
			switch {
			case f == 's':
				if err = i.MarkMessageSent(dm, entry); err != nil {
					return err
				}
				c.Confirm("%s marked sent.", entry.LocalMsgID)
			case f == 'r' && not:
				if dm.ReadyToSend() {
					dm.SetReadyToSend(false)
					if err = i.UpdateDraftMessage(entry.Ident, dm); err != nil {
						return err
					}
					c.Confirm("%s marked not ready to send.", entry.LocalMsgID)
				} else {
					err = errNoChange
					c.Confirm("No change: %s was already marked not ready to send.", entry.LocalMsgID)
				}
			case f == 'r' && !not:
				if dm.ReadyToSend() {
					err = errNoChange
					c.Confirm("No change: %s was already marked ready to send.", entry.LocalMsgID)
				} else {
					if err = message.ValidateMessage(msg, msgifc.VPacket); err != nil {
						if force {
							c := cio.Open()
							c.Error(err)
							c.Confirm("NOTE: marking ready anyway due to --force flag")
						} else {
							return err
						}
					}
					dm.SetReadyToSend(true)
					if err = i.UpdateDraftMessage(entry.Ident, dm); err != nil {
						return err
					}
					c.Confirm("%s marked ready to send.", entry.LocalMsgID)
				}
			}
		default:
			panic("not implemented")
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
