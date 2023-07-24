package shell

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/checkin"
	"github.com/rothskeller/packet/message/checkout"
	"github.com/rothskeller/packet/message/ics213"
)

// cmdNew implements the new command.
func cmdNew(args []string) bool {
	var (
		tag string
		msg message.IEdit
	)
	// Check arguments.  For convenience, aliases are available for the
	// Check-In and Check-Out message types, which would otherwise have to
	// be fully typed out because they are so similar.
	if len(args) == 1 {
		tag = args[0]
		if tag == "ci" {
			tag = checkin.Type.Tag
		}
		if tag == "co" {
			tag = checkout.Type.Tag
		}
	} else {
		io.WriteString(os.Stderr, "usage: packet new <message-type>\n")
		return false
	}
	// Translate the tag (prefix) into a message type.
	if msg = msgForTag(tag); msg == nil {
		return false
	}
	// If the message has a default body field, put the default text in it.
	if msg, ok := msg.(message.ISetBody); ok {
		msg.SetBody(config.DefBody)
	}
	return newAndReply(new(envelope.Envelope), msg)
}

var aliases = map[string]string{
	checkin.Type.Tag:  "ci",
	checkout.Type.Tag: "co",
}

func helpNew() {
	io.WriteString(os.Stdout, `The "new" command creates a draft outgoing message.
    usage: packet new <message-type>
<message-type> must be an unambiguous abbreviation of one of these types:
`)
	var tags = make([]string, 0, len(message.RegisteredTypes))
	var taglen int
	for tag := range message.RegisteredTypes {
		if _, ok := message.Create(tag).(message.IEdit); !ok {
			continue
		}
		tags = append(tags, tag)
		aliaslen := len(aliases[tag])
		if aliaslen != 0 {
			aliaslen += 3
		}
		if len(tag)+aliaslen > taglen {
			taglen = len(tag) + aliaslen
		}
	}
	sort.Strings(tags)
	for _, tag := range tags {
		if alias := aliases[tag]; alias != "" {
			fmt.Printf("    %-*s  (%s)\n", taglen, tag+" ("+alias+")", message.RegisteredTypes[tag].Name)
		} else {
			fmt.Printf("    %-*s  (%s)\n", taglen, tag, message.RegisteredTypes[tag].Name)
		}
	}
	io.WriteString(os.Stdout, `The new message is opened in an editor; run "help edit" for details.
`)
}

func cmdReply(args []string) bool {
	var (
		recvid string
		tag    string
		renv   *envelope.Envelope
		rmsg   message.Message
		senv   envelope.Envelope
		smsg   message.IEdit
		err    error
	)
	// Parse arguments.
	switch len(args) {
	case 1:
		recvid = args[0]
	case 2:
		recvid, tag = args[0], args[1]
	default:
		fmt.Fprintf(os.Stderr, "usage: packet reply <message-id> [<message-type>]\n")
		return false
	}
	// Find the received message we're replying to.
	switch lmis := expandMessageID(recvid, true); len(lmis) {
	case 0:
		fmt.Fprintf(os.Stderr, "ERROR: no such message %q\n", recvid)
		return false
	case 1:
		if renv, rmsg, err = readMessage(lmis[0]); err != nil {
			return false
		}
		if !renv.IsReceived() {
			fmt.Fprintf(os.Stderr, "ERROR: message %s is not a received message\n", lmis[0])
			return false
		}
		if renv.ReceivedArea != "" {
			io.WriteString(os.Stderr, "ERROR: cannot reply to a bulletin message\n")
			return false
		}
	default:
		fmt.Fprintf(os.Stderr, "ERROR: %q is ambiguous (%s)\n", recvid, strings.Join(lmis, ", "))
		return false
	}
	// Set up the envelope for the new message.
	senv.To = []string{renv.From}
	// Create the empty message of the appropriate type.
	if tag != "" {
		if smsg = msgForTag(tag); smsg == nil {
			return false
		}
	} else if _, ok := rmsg.(message.IEdit); ok {
		smsg = message.Create(rmsg.Type().Tag).(message.IEdit)
	} else {
		fmt.Fprintf(os.Stderr, "ERROR: can't create %s %s; specify some other message type\n",
			rmsg.Type().Article, rmsg.Type().Name)
		return false
	}
	// If the message has a default body field, put the default text in it.
	// (We may overwrite this below if the received message has a body.)
	if smsg, ok := smsg.(message.ISetBody); ok {
		smsg.SetBody(config.DefBody)
	}
	// Copy over the data from the received message to the reply.
	if rmsg, ok := rmsg.(message.HumanMessage); ok {
		if smsg, ok := smsg.(message.HumanMessage); ok {
			smsg.SetHandling(rmsg.GetHandling())
		}
		if smsg, ok := smsg.(message.ISetSubject); ok {
			smsg.SetSubject(rmsg.GetSubject())
		}
		if smsg, ok := smsg.(*ics213.ICS213); ok {
			smsg.Reference = rmsg.GetOriginID()
		}
	}
	if rmsg, ok := rmsg.(message.IGetBody); ok {
		if smsg, ok := smsg.(message.ISetBody); ok {
			if body := rmsg.GetBody(); body != "" {
				smsg.SetBody(rmsg.GetBody())
			}
		}
	}
	return newAndReply(&senv, smsg)
}

func helpReply() {
	io.WriteString(os.Stdout, `
The "reply" command starts a reply message.
    usage: packet reply <message-id> [<message-type>]
The "reply" command creates a new draft message with the same handling order
and subject as the received message identified by <message-id>.  The reply's
"To" address is set to the received message's "From" address.
    The reply is the same type of message as the received message, unless a
different <message-type> is specified.  (Enter "help new" for a list of
message types.)
    If the received message is a body-centric type (e.g., plain text or
ICS-213), its body is copied into the reply message.
    If the reply message type has a "Reference" field, it is set to the origin
message ID of the received message.
    The new message is opened in an editor; run "help edit" for details.
`)
}

// msgForTag returns a created message of the type specified by the tag.  It
// returns nil if the tag is invalid.
func msgForTag(tag string) (msg message.IEdit) {
	for rt := range message.RegisteredTypes {
		if len(rt) < len(tag) || !strings.EqualFold(tag, rt[:len(tag)]) {
			continue
		}
		if m, ok := message.Create(rt).(message.IEdit); !ok {
			continue
		} else if msg != nil {
			fmt.Fprintf(os.Stderr, "ERROR: message type %q is ambiguous\n", tag)
			return nil
		} else {
			msg = m
		}
	}
	if msg == nil {
		fmt.Fprintf(os.Stderr, "ERROR: no such message type %q\n", tag)
	}
	return msg
}

// newAndReply is the common code shared by new and reply.
func newAndReply(env *envelope.Envelope, msg message.IEdit) bool {
	// If we have a message ID pattern, give the new message an ID.
	// Otherwise we'll leave it blank and the editor will force the user to
	// provide one.
	if config.MessageID != "" {
		msg.SetOriginID(getNextMessageID(config.MessageID))
	}
	// Special case for check-in and check-out messages: fill in the call
	// signs and names from the incident/activation settings.  If we don't
	// have any yet, no harm done.
	switch msg := msg.(type) {
	case *checkin.CheckIn:
		msg.TacticalCallSign, msg.TacticalStationName = config.TacCall, config.TacName
		msg.OperatorCallSign, msg.OperatorName = config.OpCall, config.OpName
	case *checkout.CheckOut:
		msg.TacticalCallSign, msg.TacticalStationName = config.TacCall, config.TacName
		msg.OperatorCallSign, msg.OperatorName = config.OpCall, config.OpName
	}
	// Run the editor on the message.  Note that the editor is told there
	// is no LMI, even if we put one in the Origin Message ID field.  That
	// gives the editor the correct title bar, and disables the LMI change
	// behavior in case the user edits the Origin Message ID field.
	if !editMessage("", env, msg) {
		return false
	}
	// If we don't have a message ID model in the incident/activation
	// settings yet, use the Origin Message ID of the new message as our
	// model going forward.
	if config.MessageID == "" {
		setSetting([]string{"msgid", msg.GetOriginID()})
	}
	// If this was a check-in message, and we don't already have operator
	// and tactical information in the incident/activation settings, store
	// the values from the message into the settings.
	if msg, ok := msg.(*checkin.CheckIn); ok {
		if config.OpCall == "" && fccCallSignRE.MatchString(msg.OperatorCallSign) && msg.OperatorName != "" {
			setSetting([]string{"operator", msg.OperatorCallSign, msg.OperatorName})
		}
		if !config.TacRequested && config.TacCall == "" {
			if tacCallSignRE.MatchString(msg.TacticalCallSign) && msg.TacticalStationName != "" {
				setSetting([]string{"tactical", msg.TacticalCallSign, msg.TacticalStationName})
			} else if msg.TacticalCallSign == "" {
				setSetting([]string{"tactical"}) // so we don't ask for it
			}
		}
	}
	return true
}
