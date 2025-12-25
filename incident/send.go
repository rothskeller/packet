package incident

import (
	"log/slog"
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/address"
	"github.com/rothskeller/packet/message/receipt"
	"github.com/rothskeller/packet/message/subject"
)

// AddDraftMessage takes a DraftMessage and saves it in the incident, assigning
// it a local message ID along the way.  It returns the assigned LMI, which may
// be empty if the draft message is a receipt.
func (i *Incident) AddDraftMessage(msg *message.DraftMessage) (lmi string, err error) {
	var (
		le       LogEntry
		handling string
	)
	// The new message might be a receipt.  Those are handled specially.
	switch msg.Type() {
	case receipt.DeliveryReceipt, receipt.ReadReceipt:
		err = i.addDraftReceiptMessage(msg)
		return "", err
	}
	// Create a log entry and assign a local message ID.
	le = LogEntry{
		Seq:     i.Seq,
		Status:  StatusDraft,
		Time:    time.Now(),
		Subject: msg.Subject().EncodedSubject(),
	}
	if le.LocalMsgID, err = i.nextMessageID(true); err != nil {
		return "", err
	}
	le.FromMsgID = le.LocalMsgID
	le.Filename = le.LocalMsgID + ".txt"
	if msg.Bulletin() {
		le.Flags |= FBulletin
	} else {
		le.Flags |= FNeedsReceipt
	}
	// Put the local message ID and the operator information into the
	// message fields if it has them.  Also extract the handling from the
	// message fields, if any, for use in the log entry.
	for f := range msg.Fields() {
		switch f.Common() {
		case "handling":
			handling = f.Value(msg)
		case "messageDate", "formDate":
			f.SetValue(msg, time.Now().Format("01/02/2006"))
		case "operatorName":
			f.SetValue(msg, i.Config.OpName)
		case "operatorCall":
			f.SetValue(msg, i.Config.OpCall)
		case "receiverSender":
			f.SetValue(msg, "sender")
		case "operatorMethod":
			f.SetValue(msg, "Other")
		case "operatorMethodOther":
			f.SetValue(msg, "Packet")
		case "originMessageID":
			f.SetValue(msg, le.LocalMsgID)
		case "tacticalCall":
			if i.Config.TacCall != "" {
				f.SetValue(msg, i.Config.TacCall)
			}
		case "tacticalName":
			if i.Config.TacName != "" {
				f.SetValue(msg, i.Config.TacName)
			}
		case "useTactical":
			if i.Config.TacCall != "" {
				f.SetValue(msg, "checked")
			}
		}
	}
	// Check the subject line for handling that we didn't get from the
	// message body.  Then set log flags for the handling.
	if handling == "" {
		if s, ok := msg.Subject().(*subject.SCCoSubject); ok {
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
		return "", err
	}
	// Add log entries to the log, one per recipient (or exactly one if
	// there is no recipient yet, which will usually be the case).
	if addrs, err := address.ParseList(msg.To()); err != nil {
		return "", err
	} else if len(addrs) == 0 {
		i.Log = append(i.Log, &le)
	} else {
		for _, addr := range addrs {
			call, _, _ := strings.Cut(addr.Address, "@")
			call = strings.ToUpper(strings.TrimSpace(call))
			if len(call) > 8 {
				call = call[:8]
			}
			clone := le
			clone.ToCall = call
			clone.Index = len(i.Log)
			i.Log = append(i.Log, &clone)
		}
	}
	slog.Info("create draft message", "lid", le.LocalMsgID, "s", msg.Subject().EncodedSubject())
	return le.LocalMsgID, nil
}

func (inc *Incident) addDraftReceiptMessage(msg *message.DraftMessage) error {
	panic("unimplemented")
}
