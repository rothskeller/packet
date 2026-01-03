package incident

import (
	"log/slog"
	"strings"
	"time"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/receipt"
	"github.com/rothskeller/packet/message/subject"
)

// ReceiveMessage takes a JustReceivedMessage received from JNOS and saves it
// in the incident.  ReceiveMessage returns the delivery receipt that should be
// sent for the message, if any; it is up to the caller to queue the delivery
// receipt for sending.
func (i *Incident) ReceiveMessage(msg *message.JustReceivedMessage) (dr *message.DraftMessage, err error) {
	var (
		le       LogEntry
		handling string
	)
	// The received message might be a receipt.  Those are handled
	// specially.
	switch msg.Type() {
	case receipt.DeliveryReceipt, receipt.ReadReceipt:
		err = i.receiveReceiptMessage(msg)
		return nil, err
	}
	// Create a log entry and assign a local message ID.
	le = LogEntry{
		Ident:   i.nextLogIdent(),
		Index:   len(i.Log),
		Seq:     i.Seq,
		Status:  StatusReceived,
		Flags:   FUnread,
		Time:    msg.RxDate(),
		Subject: msg.Subject().EncodedSubject(),
	}
	if le.LocalMsgID, err = i.nextMessageID(false); err != nil {
		return nil, err
	}
	le.ToMsgID = le.LocalMsgID
	if msg.Bulletin() {
		le.Flags |= FBulletin
		le.FromCall = strings.ToUpper(msg.RxArea())
	} else {
		le.FromCall, _, _ = strings.Cut(msg.To(), "@")
		le.FromCall = strings.ToUpper(strings.TrimSpace(le.FromCall))
	}
	// Put the local message ID and the operator information into the
	// message fields if it has them.  Also extract the OMI and handling
	// from the message fields, if any, for use in the log entry.
	for f := range msg.Fields(msg) {
		switch f.Common() {
		case "originMessageID":
			le.FromMsgID = f.Value(msg)
		case "destinationMessageID":
			f.SetValue(msg, le.LocalMsgID)
		case "handling":
			handling = f.Value(msg)
		case "operatorName":
			f.SetValue(msg, i.Config.OpName)
		case "operatorCall":
			f.SetValue(msg, i.Config.OpCall)
		case "operatorDate":
			f.SetValue(msg, msg.RxDate().Format("01/02/2006"))
		case "operatorTime":
			f.SetValue(msg, msg.RxDate().Format("15:04"))
		case "receiverSender":
			f.SetValue(msg, "receiver")
		case "operatorMethod":
			f.SetValue(msg, "Other")
		case "operatorMethodOther":
			f.SetValue(msg, "Packet")
		}
	}
	// Check the subject line for OMI and handling that we didn't get from
	// the message body.  Then set log flags for the handling.
	if s, ok := msg.Subject().(*subject.SCCoSubject); ok {
		if le.FromMsgID == "" {
			le.FromMsgID = s.SubjectMessageID()
		}
		if handling == "" {
			handling = s.SubjectHandling()
		}
	}
	switch handling {
	case "IMMEDIATE":
		le.Flags |= FImmediate
	case "PRIORITY":
		le.Flags |= FPriority
	}
	// Save the message.
	if err = i.saveMessage(msg, &le); err != nil {
		return nil, err
	}
	// Add the log entry to the log.
	i.Log = append(i.Log, &le)
	i.sortLog()
	// Generate a delivery receipt if appropriate.  (It's up to the caller
	// to send it or not.)
	if !msg.Bulletin() && !msg.Autoresponse() {
		dr, _ = i.MakeDeliveryReceipt(msg, &le)
	}
	slog.Info("received message", "lid", le.LocalMsgID, "s", msg.Subject().EncodedSubject())
	return dr, nil
}

// MakeDeliveryReceipt makes a delivery receipt for the supplied received
// message.
func (i *Incident) MakeDeliveryReceipt(msg message.Message, le *LogEntry) (dr *message.DraftMessage, err error) {
	var date time.Time

	switch msg := msg.(type) {
	case *message.ReceivedMessage:
		date = msg.RxDate()
	case *message.JustReceivedMessage:
		if msg.Autoresponse() {
			return nil, errors.New("Delivery receipts are not appropriate for automatically generated messages.")
		}
		date = msg.RxDate()
	case *message.DraftMessage, *message.SentMessage:
		return nil, errors.New("Delivery receipts cannot be sent for outgoing messages.")
	}
	if msg.Bulletin() {
		return nil, errors.New("Delivery receipts are not appropriate for bulletins and notices.")
	}
	switch msg.Type() {
	case receipt.DeliveryReceipt, receipt.ReadReceipt:
		return nil, errors.New("Delivery receipts are not appropriate for receipt messages.")
	}
	for _, e := range i.Log {
		if e.LocalMsgID == le.LocalMsgID && e.Flags&FIsReceipt != 0 {
			return nil, errors.New("A delivery receipt has already been generated for this message.")
		}
	}
	if dr, err = receipt.NewDeliveryReceipt(msg.To(), msg.Subject().EncodedSubject(), le.LocalMsgID, date, ""); err != nil {
		return nil, err
	}
	return dr, nil
}

func (i *Incident) receiveReceiptMessage(msg *message.JustReceivedMessage) (err error) {
	panic("not implemented")
}

/*
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
	if err2 := SaveMessage(lmi, rmi, env, msg, false, true); err2 != nil {
		err = fmt.Errorf("save received %s: %s", lmi, err2)
		return
	}
	if area != "" || env.Autoresponse { // bulletin, bounce: no delivery receipt
		return
	}
	// Return delivery receipt.
	dr := delivrcpt.New()
	dr.LocalMessageID = lmi
	dr.DeliveredTime = time.Now().Format("01/02/2006 15:04")
	dr.MessageSubject = env.SubjectLine
	dr.MessageTo = env.To
	denv := new(envelope.Envelope)
	denv.SubjectLine = dr.EncodeSubject()
	denv.To = env.From
	return lmi, env, msg, denv, dr, err
}

// recordReceipt matches a received receipt with the corresponding outgoing
// message.
func recordReceipt(env *envelope.Envelope, msg message.Message) (
	lmi string, oenv *envelope.Envelope, omsg message.Message, err error,
) {
	var (
		subject string
		to      string
		rmi     string
	)
	switch msg := msg.(type) {
	case *delivrcpt.DeliveryReceipt:
		subject, to = msg.MessageSubject, msg.MessageTo
		rmi = msg.LocalMessageID
	case *readrcpt.ReadReceipt:
		subject, to = msg.MessageSubject, msg.MessageTo
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
	if err = SaveReceipt(lmi, env, msg); err != nil {
		err = fmt.Errorf("save receipt for %s: %s", lmi, err)
		return
	}
	if rmi == "" {
		return // read receipt, nothing more to do
	}
	if mb := msg.Base(); mb.FDestinationMsgID != nil && *mb.FDestinationMsgID == "" {
		*mb.FDestinationMsgID = rmi
	}
	if err = SaveMessage(lmi, rmi, oenv, omsg, false, false); err != nil {
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
	env := new(envelope.Envelope)
	env.Date = rcptenv.Date
	env.From = rcptenv.To
	env.SubjectLine = subject
	env.To = to
	content := env.RenderSaved(`**** MESSAGE CONTENTS UNKNOWN ****

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

*/
