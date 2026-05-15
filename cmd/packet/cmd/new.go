package cmd

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/cmd/packet/pseudomsg"
	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/messageid"
	"github.com/spf13/pflag"
)

const (
	newSlug = `Create a new outgoing message or log entry`
	newHelp = `
usage: packet new ⇥[-flags] [«new-message-type»] [«new-message-id»]
  -c, --copy «message-id»    ⇥Send a copy of an existing message
  -l, --log                  ⇥Create a log entry not associated with a message
  -r, --reply «message-id»   ⇥Send a reply to an existing received message
  -x, --resend «message-id»  ⇥Correct and resend an already-sent message

The "packet new" (or "n") command creates a new outgoing message or log entry.  If standard input and output are terminals, the new message or log entry will be opened for editing (see "packet help edit" for details).  Otherwise, the local message ID for the new message (or log entry number for the new log entry) will be printed to standard output, and subsequent "set" commands can be used to populate it.

When the --copy (or -c) flag is given, the new message will be an exact copy of the named source message except for being given a new local message ID.  A «new-message-type» argument is not allowed.

When the --log (or -l) flag is given, a new log entry is created that is not associated with any message.  No arguments are allowed.

When the --reply (or -r) flag is given, the new message will have the same handling order, subject line, and body as the named source message, which must be a received message.  Its "To" address will be set to the "From" address of the source message.  The To and From ICS Position, Location, Name, and Contact fields, if any, will be swapped.  The new message will have the same message type as the source message unless a «new-message-type» is given on the command line.  If the message type has a "Reference" field, it will be filled with the source message's origin message ID.

When the --resend (or -x) flag is given, the new message will be identical to the source message except with the suffix of the message number changed to "R" (or the next available letter if "R" has already been used).  The source message must be an already-sent message.

When no flags are given, an empty message of «new-message-type» is created.  «new-message-type» must be the tag name or create key of one of the supported message types; if omitted, a plain text message is created.  Use "packet forms list" to get a list of supported forms.

If a «new-message-id» is provided on the command line, the new message is created with that local message ID.  The sequence number in it will be incremented as needed to make it unique.  The «new-message-id» may be just an integer, in which case the message number and prefix in the incident configuration are used (see "packet help config").  If no «new-message-id» is given, one will be automatically assigned based on the incident configuration.
`
)

func cmdNew(args []string) (err error) {
	var (
		copyOf  string
		newLog  bool
		replyTo string
		resend  string
		mtname  string
		msgtype message.EditableMType
		msgid   string
		mpfx    string
		mseq    int
		msfx    string
		newmsg  *message.DraftMessage
		newle   *incident.LogEntry
		c       = cio.Open()
	)
	flags := pflag.NewFlagSet("new", pflag.ContinueOnError)
	flags.StringVarP(&copyOf, "copy", "c", "", "Send a copy of an existing message")
	flags.BoolVarP(&newLog, "log", "l", false, "Create a log entry not associated with a message")
	flags.StringVarP(&replyTo, "reply", "r", "", "Send a reply to an existing received message")
	flags.StringVarP(&resend, "resend", "x", "", "Correct and resend an already-sent message")
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"new"})
	}
	if err == nil {
		err = gaveMutuallyExclusiveFlags(flags, "copy", "log", "reply", "resend")
	}
	if err != nil {
		c.Error(err)
		return usage(newHelp)
	}
	registerForms() // needed to recognize a message type
	switch flags.NArg() {
	case 0:
		// nothing
	case 1:
		arg := flags.Arg(0)
		for mt := range message.AllTypes() {
			if emt, ok := mt.(message.EditableMType); ok {
				if strings.EqualFold(arg, emt.CreateTag()) || strings.EqualFold(arg, emt.CreateKey()) {
					msgtype = emt
					break
				}
			}
		}
		if msgtype == nil {
			msgid = arg
		}
	case 2:
		mtname, msgid = flags.Arg(0), flags.Arg(1)
	default:
		return usage(newHelp)
	}
	if copyOf != "" && (msgtype != nil || mtname != "") {
		c.ErrorF(`When using the --copy flag, do not specify a message type.  The copy will have the same message type as the original.`)
		return usage(newHelp)
	}
	if newLog && (msgtype != nil || mtname != "" || msgid != "") {
		c.ErrorF(`When using the --log flag, do not specify a message type or message ID.`)
		return usage(newHelp)
	}
	if resend != "" && (msgtype != nil || mtname != "" || msgid != "") {
		c.ErrorF(`When using the --resend flag, do not specify a message type or new message ID.`)
		return usage(newHelp)
	}
	if mtname != "" {
		for mt := range message.AllTypes() {
			if emt, ok := mt.(message.EditableMType); ok {
				if strings.EqualFold(mtname, emt.CreateTag()) || strings.EqualFold(mtname, emt.CreateKey()) {
					msgtype = emt
					break
				}
			}
		}
		if msgtype == nil {
			c.ErrorF(`There is no editable message type %q.  Use "packet forms list" to get a list of message types.`, mtname)
			return usage(newHelp)
		}
	}
	if msgid != "" {
		if mseq, err = strconv.Atoi(msgid); err != nil || mseq < 1 {
			if mpfx, mseq, msfx, err = messageid.Decode(msgid, true, true); err != nil {
				c.ErrorF(`%q is not a valid new message ID: %s`, msgid, err)
				return usage(newHelp)
			}
		}
	}
	if err = incWrite(true, func(i *incident.Incident) error {
		var (
			srcmsg message.Message
			srcle  *incident.LogEntry
		)
		// Creating a new log entry is a special case.
		if newLog {
			newle = &incident.LogEntry{Time: time.Now()}
			i.AddLogEntry(newle)
			if c.InputIsTerm && c.OutputIsTerm {
				return doEdit(c, i, newle, pseudomsg.NewLogEntryMessage(newle), nil, false, true)
			}
			return nil
		}
		// For real messages, we need some config.
		if err = requiredConfig(i, "OpCall", "OpName", "TxMessageID"); err != nil {
			return err
		}
		// Get the message we're supposed to copy, reply to, or resend,
		// if any.
		if copyOf != "" {
			if srcmsg, srcle, err = matchMessage(i, copyOf, MMMessageOnly); err != nil {
				return err
			}
		}
		if replyTo != "" {
			if srcmsg, srcle, err = matchMessage(i, replyTo, MMMessageOnly); err != nil {
				return err
			} else if _, ok := srcmsg.(*message.ReceivedMessage); !ok {
				return errors.NewF(`Message %s is not a received message.`, srcle.LocalMsgID)
			}
		}
		if resend != "" {
			if srcmsg, srcle, err = matchMessage(i, resend, MMMessageOnly); err != nil {
				return err
			} else if _, ok := srcmsg.(*message.SentMessage); !ok {
				return errors.NewF(`Message %s is not a sent message.`, srcle.LocalMsgID)
			}
		}
		// If we have a new message ID, sanitize it and be sure it's not
		// already in use.
		if mseq != 0 {
			if mpfx == "" {
				mpfx, _, msfx, _ = messageid.Decode(i.Config.TxMessageID, true, false)
			}
			msgid, _ = messageid.Encode(mpfx, mseq, msfx)
			if slices.IndexFunc(i.Log, func(le *incident.LogEntry) bool {
				return le.LocalMsgID == msgid
			}) >= 0 {
				return errors.NewF(`The message ID %s is already in use.`, msgid)
			}
		} else if resend != "" {
			if msgid, err = i.ResendMessageID(srcle.LocalMsgID); err != nil {
				return err
			}
		}
		// If we don't already have a message type, but we have a source
		// message, use that type.
		if msgtype == nil && srcmsg != nil {
			if emt, ok := srcmsg.Type().(message.EditableMType); ok {
				msgtype = emt
			} else {
				return errors.NewF(`New %ss cannot be created.  Use a current message type instead.`,
					strings.Join(strings.Fields(srcmsg.Type().Name())[1:], " "))
			}
		}
		// If we still don't have a message type, make a plain text
		// message.
		if msgtype == nil {
			msgtype = message.PlainMessage
		}
		// Create a new draft message of that type.
		newmsg = msgtype.NewDraft().(*message.DraftMessage)
		// Set the fields of the new message based on the source.
		if copyOf != "" {
			message.CopyFields(srcmsg, newmsg)
		} else if replyTo != "" {
			i.ApplyDefaults(newmsg)
			message.MakeReply(srcmsg.(*message.ReceivedMessage), newmsg)
		} else if resend != "" {
			message.CopyFields(srcmsg.(*message.SentMessage), newmsg)
			newmsg.SetTo(srcmsg.To())
			for f := range newmsg.Fields() {
				switch f.Common() {
				case field.COriginMessageID, field.CSubjectMessageID:
					f.SetValue(newmsg, msgid)
				}
			}
		} else {
			i.ApplyDefaults(newmsg)
		}
		if newle, err = i.AddDraftMessage(newmsg); err != nil {
			return err
		}
		if c.InputIsTerm && c.OutputIsTerm {
			return doEdit(c, i, newle, newmsg, nil, false, false)
		}
		return nil
	}); err != nil {
		return err
	}
	if !c.InputIsTerm || !c.OutputIsTerm {
		if newLog {
			fmt.Println(newle.Ident)
		} else {
			fmt.Println(newle.LocalMsgID)
		}
	}
	return nil
}
