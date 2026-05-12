package incident

import (
	"log/slog"
	"strings"
	"time"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/receipt"
)

// ReceiveMessage takes a JustReceivedMessage received from JNOS and saves it
// in the incident.  ReceiveMessage returns the delivery receipt that should be
// sent for the message, if any; it is up to the caller to queue the delivery
// receipt for sending.  It also returns the log entry for the received message.
func (i *Incident) ReceiveMessage(msg *message.JustReceivedMessage) (dr *message.DraftMessage, le *LogEntry, err error) {
	var handling string

	// The received message might be a receipt.  Those are handled
	// specially.
	switch msg.Type() {
	case receipt.DeliveryReceipt, receipt.ReadReceipt:
		err = i.receiveReceiptMessage(msg)
		return nil, nil, err
	}
	// Create a log entry and assign a local message ID.
	le = &LogEntry{
		Ident:   i.nextLogIdent(),
		Index:   len(i.Log),
		Seq:     i.Seq,
		Status:  StatusReceived,
		Flags:   FUnread,
		Time:    msg.RxDate(),
		Subject: msg.Subject().EncodedSubject(),
	}
	if le.LocalMsgID, err = i.nextMessageID(false); err != nil {
		return nil, nil, err
	}
	le.ToMsgID = le.LocalMsgID
	msg.SetLocalID(le.LocalMsgID)
	if msg.Bulletin() {
		le.Flags |= FBulletin
		le.FromCall = strings.ToUpper(msg.RxArea())
	} else {
		le.FromCall, _, _ = strings.Cut(msg.From(), "@")
		le.FromCall = strings.ToUpper(strings.TrimSpace(le.FromCall))
	}
	// Put the local message ID and the operator information into the
	// message fields if it has them.  Also extract the OMI and handling
	// from the message fields, if any, for use in the log entry.
	for f := range msg.Fields() {
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
	if le.FromMsgID == "" {
		le.FromMsgID = msg.Subject().SubjectMessageID()
	}
	if handling == "" {
		handling = msg.Subject().SubjectHandling()
	}
	switch handling {
	case "IMMEDIATE":
		le.Flags |= FImmediate
	case "PRIORITY":
		le.Flags |= FPriority
	}
	// Save the message.
	if err = i.saveMessage(msg, le); err != nil {
		return nil, nil, err
	}
	// Add the log entry to the log.
	i.Log = append(i.Log, le)
	i.sortLog()
	// Generate a delivery receipt if appropriate.  (It's up to the caller
	// to send it or not.)
	if !msg.Bulletin() && !msg.Autoresponse() && !i.Config.NoSendReceipts {
		dr, _ = i.MakeDeliveryReceipt(msg, le)
	}
	slog.Info("received message", "lid", le.LocalMsgID, "s", msg.Subject().EncodedSubject())
	return dr, le, nil
}

// MakeDeliveryReceipt makes a delivery receipt for the supplied received
// message.
func (i *Incident) MakeDeliveryReceipt(msg message.Message, le *LogEntry) (dr *message.DraftMessage, err error) {
	var date time.Time
	var from string

	switch msg := msg.(type) {
	case *message.ReceivedMessage:
		date = msg.RxDate()
		from = msg.From()
	case *message.JustReceivedMessage:
		if msg.Autoresponse() {
			return nil, errors.New("Delivery receipts are not appropriate for automatically generated messages.")
		}
		date = msg.RxDate()
		from = msg.From()
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
	if dr, err = receipt.NewDeliveryReceipt(from, msg.To(), msg.Subject().EncodedSubject(), le.LocalMsgID, date, ""); err != nil {
		return nil, err
	}
	return dr, nil
}

func (i *Incident) receiveReceiptMessage(rcpt *message.JustReceivedMessage) (err error) {
	// Create a log entry.
	var rcptle = LogEntry{
		Ident:   i.nextLogIdent(),
		Index:   len(i.Log),
		Seq:     i.Seq,
		Status:  StatusReceived,
		Flags:   FIsReceipt,
		Time:    rcpt.RxDate(),
		Subject: rcpt.Subject().EncodedSubject(),
	}
	var smr = message.SentMessageReceipt{ReceiverAddress: rcpt.From()}
	var subject string
	if rcpt.MType == receipt.DeliveryReceipt {
		drb := rcpt.Body().(*receipt.DeliveryReceiptBody)
		smr.ReceiverMessageID = drb.LocalMessageID()
		smr.ReceiptDate = drb.DeliveryTime()
		subject = drb.MessageSubject()
	} else {
		rrb := rcpt.Body().(*receipt.ReadReceiptBody)
		smr.ReceiptDate = rrb.ReadTime()
		smr.HasBeenRead = true
		subject = rrb.MessageSubject()
	}
	// Find the matching outgoing message, if any.
	for idx := len(i.Log) - 1; idx >= 0; idx-- {
		sentle := i.Log[idx]
		if sentle.Status != StatusSent || sentle.Subject != subject {
			continue
		}
		rcptle.LocalMsgID = sentle.LocalMsgID
		// Read the message.
		sent, err := i.GetMessageFromLogEntry(sentle)
		if err != nil {
			return err
		}
		// Add the receipt information to the header.
		sent.(*message.SentMessage).AddReceipt(smr)
		// Add the destination message ID to the message if appropriate.
		if smr.ReceiverMessageID != "" {
			for f := range sent.Fields() {
				if f.Common() == field.CDestinationMessageID {
					if f.Value(sent) == "" {
						f.SetValue(sent, smr.ReceiverMessageID)
					}
					break
				}
			}
		}
		sentle.Seq = i.Seq
		if err = i.saveMessage(sent, sentle); err != nil {
			return err
		}
		slog.Info("added receipt info to sent message", "lmi", sentle.LocalMsgID)
		// Find the log entry matching this particular recipient.
		fromCall, _, _ := strings.Cut(rcpt.From(), "@")
		for ; idx >= 0; idx-- {
			sentle := i.Log[idx]
			if sentle.Status != StatusSent || sentle.LocalMsgID != rcptle.LocalMsgID ||
				!strings.EqualFold(fromCall, sentle.ToCall) || sentle.Flags&FHasReceipt != 0 || sentle.ToMsgID != "" {
				continue
			}
			sentle.ToMsgID = smr.ReceiverMessageID
			sentle.Flags |= FHasReceipt
			sentle.Flags &^= FNeedsReceipt
			sentle.Seq = i.Seq
			break
		}
		if idx < 0 {
			slog.Warn("receipt not matched to any known addressee of sent message")
		}
		break
	}
	if rcptle.LocalMsgID == "" {
		slog.Warn("receipt not matched to any sent message")
	}
	// Save the receipt message.
	if err = i.saveMessage(rcpt, &rcptle); err != nil {
		return err
	}
	// Add the log entry to the log.
	i.Log = append(i.Log, &rcptle)
	i.sortLog()
	slog.Info("received receipt", "s", rcpt.Subject().EncodedSubject())
	return nil
}
