package incident

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/delivrcpt"
	"github.com/rothskeller/packet/xscmsg/readrcpt"
)

var ErrDuplicateReceipt = errors.New("duplicate receipt")

// ReceiveMessage takes a raw message received from JNOS and saves it in the
// incident. "bbs" and "area" are the names of the BBS from which the message
// was retrieved and, if the message is a bulletin, the bulletin area from which
// it was retrieved. "msgid" is the message ID pattern; a unique ID will be
// given to the new message using this pattern if it needs one.  "opcall" and
// "opname" are the operator call sign and name to put into the received
// message.
//
// If the received message is a human-originated message, it will be returned as
// "lmi", "env", and "msg".  If a delivery receipt should be sent for it, the
// receipt is returned as "oenv" and "omsg".  (Note that the delivery receipt
// has not been saved; the caller should do that after sending it.)
//
// If the received message is a receipt, it will be returned as "env" and "msg".
// If the receipt can be matched against a previously sent message in the
// incident, that previously sent message is returned as "lmi", "oenv", and
// "omsg" are the message for which it is a receipt.  Note: in the case where
// a receipt has already been received for the message, all five return values
// are filled in, and "err" will be ErrDuplicateReceipt.
//
// If the received message has an error, "err" will be set and the other return
// values will be zero.
func ReceiveMessage(raw, bbs, area, msgid, opcall, opname string) (
	lmi string, env *envelope.Envelope, msg message.Message, oenv *envelope.Envelope, omsg message.Message, err error,
) {
	// Parse the message.
	var body string
	env, body, err = envelope.ParseRetrieved(raw, bbs, area)
	if err != nil {
		err = fmt.Errorf("parse retrieved message: %s", err)
		return
	}
	if env.Autoresponse {
		env = nil
		return // autoresponses are ignored
	}
	msg = message.Decode(env.SubjectLine, body)
	// If it's a receipt, it's handled specially.
	switch msg.(type) {
	case *delivrcpt.DeliveryReceipt, *readrcpt.ReadReceipt:
		lmi, oenv, omsg, err = recordReceipt(env, msg)
		return
	}
	// Assign a local message ID.  Put it, and the opcall/opname, into the
	// message if it has fields for it.
	lmi = UniqueMessageID(msgid)
	if mb := msg.Base(); mb.FDestinationMsgID != nil {
		*mb.FDestinationMsgID = lmi
	}
	msg.SetOperator(opcall, opname, true)
	// Save the message.
	var rmi string
	if b := msg.Base(); b.FOriginMsgID != nil {
		rmi = *b.FOriginMsgID
	}
	if err = SaveMessage(lmi, rmi, env, msg, false); err != nil {
		err = fmt.Errorf("save received %s: %s", lmi, err)
		return
	}
	if area != "" { // bulletin: no delivery receipt
		return
	}
	// Return delivery receipt.
	var dr = delivrcpt.New()
	dr.LocalMessageID = lmi
	dr.DeliveredTime = time.Now().Format("01/02/2006 15:04")
	dr.MessageSubject = env.SubjectLine
	dr.MessageTo = env.To
	var denv = new(envelope.Envelope)
	denv.SubjectLine = dr.EncodeSubject()
	denv.To = env.From
	return lmi, env, msg, denv, dr, nil
}

// recordReceipt matches a received receipt with the corresponding outgoing
// message.
func recordReceipt(env *envelope.Envelope, msg message.Message) (
	lmi string, oenv *envelope.Envelope, omsg message.Message, err error,
) {
	var (
		subject string
		dtstamp string
		to      string
		ext     string
		rmi     string
	)
	switch msg := msg.(type) {
	case *delivrcpt.DeliveryReceipt:
		subject, to, dtstamp = msg.MessageSubject, msg.MessageTo, msg.DeliveredTime
		ext = ".DR"
		rmi = msg.LocalMessageID
	case *readrcpt.ReadReceipt:
		subject, to = msg.MessageSubject, msg.MessageTo
		ext = ".RR"
	}
	if subject != "" {
		if lmi, err = subjectToLMI(subject); err != nil {
			return "", nil, nil, err
		}
	}
	if lmi == "" {
		if lmi, err = makeFakeSentMessage(subject, to, env); err != nil {
			return "", nil, nil, err
		}
	}
	if lmi == "" {
		return
	}
	if oenv, omsg, err = ReadMessage(lmi); err != nil {
		err = fmt.Errorf("read message %s for receipt: %s", lmi, err)
		return
	}
	if _, err = os.Stat(lmi + ext); !errors.Is(err, os.ErrNotExist) {
		err = ErrDuplicateReceipt
		return
	}
	if err = SaveReceipt(lmi, env, msg); err != nil {
		err = fmt.Errorf("save receipt for %s: %s", lmi, err)
		return
	}
	if rmi == "" {
		return // read receipt, nothing more to do
	}
	oenv.DeliveredDate, oenv.DeliveredRMI = dtstamp, rmi
	if mb := msg.Base(); mb.FDestinationMsgID != nil {
		*mb.FDestinationMsgID = rmi
	}
	if err = SaveMessage(lmi, rmi, oenv, omsg, false); err != nil {
		err = fmt.Errorf("add RMI: save message %s: %s", lmi, err)
		return
	}
	return
}

// subjectToLMI scans all sent messages in reverse chronological order looking
// for one with the specified subject.  If found, it returns the LMI.
func subjectToLMI(subject string) (lmi string, err error) {
	lmis, err := AllLMIs()
	if err != nil {
		return "", err
	}
	for i := len(lmis) - 1; i >= 0; i-- {
		lmi = lmis[i]
		if env, _, err := readEnvelope(lmi, ""); err == nil &&
			!env.IsReceived() && env.IsFinal() && env.SubjectLine == subject {
			return lmi, nil
		}
	}
	return "", nil
}

func makeFakeSentMessage(subject, to string, rcptenv *envelope.Envelope) (lmi string, err error) {
	// Can we discern an LMI from the subject line of the message being
	// receipted?
	if lmi, _, _, _, _ = message.DecodeSubject(subject); !MsgIDRE.MatchString(lmi) {
		return "", nil
	}
	// Is that LMI available?  We don't already have something named that?
	if _, err = os.Stat(lmi + ".txt"); !errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	// Create a fake sent message.
	var env = new(envelope.Envelope)
	env.Date = rcptenv.Date
	env.From = rcptenv.To
	env.SubjectLine = subject
	env.To = to
	var content = env.RenderSaved(`**** MESSAGE CONTENTS UNKNOWN ****

A receipt was received for a message with this ID, but that message was sent
in a different incident or by different software.
`)
	// Save the message to its text file.
	if err = os.WriteFile(lmi+".txt", []byte(content), 0666); err != nil {
		return "", err
	}
	// Set the modification time of the text file based on the envelope.
	// (This allows AllLMIs to sort by file modification time.)
	if !env.Date.IsZero() {
		os.Chtimes(lmi+".txt", env.Date, env.Date) // error ignored
	}
	return lmi, nil
}
