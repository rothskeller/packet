package incident

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/address"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/messageid"
	"github.com/rothskeller/packet/message/receipt"
	"github.com/rothskeller/packet/message/subject"
)

// AddDraftMessage takes a DraftMessage and saves it in the incident, assigning
// it a local message ID along the way unless the message already has one.  If
// defaults is true, default values from the incident configuration override
// values provided in the message; otherwise, it's the other way around.  The
// function returns the ident of the corresponding new log entry.
func (i *Incident) AddDraftMessage(msg *message.DraftMessage, defaults bool) (ident int, err error) {
	var (
		le       LogEntry
		handling string
	)
	// The new message might be a receipt.  Those are handled specially.
	switch msg.Type() {
	case receipt.ReadReceipt:
		return 0, errors.New("This software does not support sending read receipts, only receiving them.")
	case receipt.DeliveryReceipt:
		return i.addDraftDeliveryReceipt(msg)
	}
	// Create a log entry and assign a local message ID.
	le = LogEntry{
		Ident:  i.nextLogIdent(),
		Index:  len(i.Log),
		Seq:    i.Seq,
		Status: StatusDraft,
		Time:   time.Now(),
	}
	if msg.Bulletin() {
		le.Flags |= FBulletin
	}
	if addrs, err := address.ParseList(msg.To()); err == nil && len(addrs) != 0 {
		le.ToCall = addressToLogCall(addrs[0])
	}
	// Put the local message ID, incident defaults, and operator
	// information into the message fields if it has them.  Also extract
	// the handling from the message fields, if any, for use in the log
	// entry.
	for f := range msg.Fields(msg) {
		switch f.Common() {
		case field.CDefaultBody:
			maybeSetValue(msg, f, i.Config.DefaultBody, defaults)
		case field.CFromICSPosition:
			maybeSetValue(msg, f, i.Config.DefaultFromPos, defaults)
		case field.CFromLocation:
			maybeSetValue(msg, f, i.Config.DefaultFromLoc, defaults)
		case field.CHandling:
			handling = f.Value(msg)
		case field.CMessageDate, field.CFormDate:
			maybeSetValue(msg, f, time.Now().Format("01/02/2006"), defaults)
		case field.COperatorCall:
			f.SetValue(msg, i.Config.OpCall)
		case field.COperatorDate, field.COperatorTime:
			f.SetValue(msg, "")
		case field.COperatorMethod:
			f.SetValue(msg, "Other")
		case field.COperatorMethodOther:
			f.SetValue(msg, "Packet")
		case field.COperatorName:
			f.SetValue(msg, i.Config.OpName)
		case field.COriginMessageID:
			if le.LocalMsgID = f.Value(msg); le.LocalMsgID == "" {
				if le.LocalMsgID, err = i.nextMessageID(true); err != nil {
					return 0, err
				}
				f.SetValue(msg, le.LocalMsgID)
			}
			le.FromMsgID = le.LocalMsgID
		case field.CReceiverSender:
			f.SetValue(msg, "sender")
		case field.CTacticalCall:
			if i.Config.TacCall != "" {
				f.SetValue(msg, i.Config.TacCall)
			}
		case field.CTacticalName:
			if i.Config.TacName != "" {
				f.SetValue(msg, i.Config.TacName)
			}
		case field.CToICSPosition:
			maybeSetValue(msg, f, i.Config.DefaultToPos, defaults)
		case field.CToLocation:
			maybeSetValue(msg, f, i.Config.DefaultToLoc, defaults)
		case field.CUseTactical:
			if i.Config.TacCall != "" {
				f.SetValue(msg, "checked")
			}
		}
	}
	// Check the subject line for info that we didn't get from the
	// message body.  Then set log flags for the handling.
	if s, ok := msg.Subject().(*subject.SCCoSubject); ok {
		if le.LocalMsgID == "" {
			le.LocalMsgID = s.SubjectMessageID()
			le.FromMsgID = le.LocalMsgID
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
	// If we still don't have a local message ID, assign one.
	if le.LocalMsgID == "" {
		if le.LocalMsgID, err = i.nextMessageID(true); err != nil {
			return 0, err
		}
		if s, ok := msg.Subject().(*subject.SCCoSubject); ok {
			s.SetSubjectMessageID(le.LocalMsgID)
		}
		le.FromMsgID = le.LocalMsgID
	}
	le.Subject = msg.Subject().EncodedSubject()
	// Save the message.
	if err = i.saveMessage(msg, &le); err != nil {
		return 0, err
	}
	i.Log = append(i.Log, &le)
	i.sortLog()
	slog.Info("create draft message", "id", le.Ident, "lid", le.LocalMsgID, "s", msg.Subject().EncodedSubject())
	return le.Ident, nil
}

func maybeSetValue(msg *message.DraftMessage, f field.Field, value string, defaults bool) {
	if (defaults && value != "") || f.Value(msg) == "" {
		f.SetValue(msg, value)
	}
}

func (i *Incident) addDraftDeliveryReceipt(dr *message.DraftMessage) (ident int, err error) {
	// Create a log entry and assign a local message ID.
	var le = LogEntry{
		Ident:      i.nextLogIdent(),
		Index:      len(i.Log),
		Seq:        i.Seq,
		Status:     StatusQueued,
		Time:       time.Now(),
		Subject:    dr.Subject().EncodedSubject(),
		Flags:      FIsReceipt,
		LocalMsgID: dr.Body().(*receipt.DeliveryReceiptBody).LocalMessageID(),
	}
	if addrs, err := address.ParseList(dr.To()); err == nil && len(addrs) != 0 {
		le.ToCall = addressToLogCall(addrs[0])
	}
	// Save the message.
	if err = i.saveMessage(dr, &le); err != nil {
		return 0, err
	}
	i.Log = append(i.Log, &le)
	i.sortLog()
	slog.Info("create draft receipt", "id", le.Ident, "s", dr.Subject().EncodedSubject())
	return le.Ident, nil
}

// UpdateDraftMessage saves changes to an existing DraftMessage in the incident.
func (i *Incident) UpdateDraftMessage(ident int, msg *message.DraftMessage) (err error) {
	var (
		le       *LogEntry
		omi      string
		handling string
	)
	// Find the existing log entry for the message.
	if le = i.GetLogEntryByIdent(ident); le == nil {
		slog.Error("no such message", "dir", i.Dir, "id", ident)
		return fmt.Errorf("no such message %d", ident)
	}
	if (le.Status != StatusDraft && le.Status != StatusQueued) || le.Flags&FIsReceipt != 0 {
		slog.Error("message is not a draft", "dir", i.Dir, "id", ident)
		return fmt.Errorf("message %d is not a draft", ident)
	}
	// Get the origin message ID and handling order from the message.
	for f := range msg.Fields(msg) {
		switch f.Common() {
		case field.COriginMessageID:
			omi = f.Value(msg)
		case field.CHandling:
			handling = f.Value(msg)
		}
	}
	if s, ok := msg.Subject().(*subject.SCCoSubject); ok {
		if omi == "" {
			omi = s.SubjectMessageID()
		}
		if handling == "" {
			handling = s.SubjectHandling()
		}
	}
	// If the origin message ID in the message is different from the LMI,
	// we'll need to change the LMI.  But first, make sure it's valid and
	// not in use.
	if omi != "" && omi != le.LocalMsgID {
		if omi, err = messageid.Cleanup(omi, false); err != nil {
			slog.Error("invalid origin message ID", "omi", omi)
			return fmt.Errorf("invalid origin message ID %q", omi)
		}
		if i.logEntryForLMI(omi) != nil {
			slog.Error("origin message ID already in use", "omi", omi)
			return fmt.Errorf("another message has origin message ID %q", omi)
		}
		le.LocalMsgID, le.FromMsgID = omi, omi
	}
	// Save the message.
	if err = i.saveMessage(msg, le); err != nil {
		return err
	}
	// Update the log entry.
	le.Seq = i.Seq
	le.Subject = msg.Subject().EncodedSubject()
	if msg.ReadyToSend() {
		le.Status = StatusQueued
	} else {
		le.Status = StatusDraft
	}
	le.Flags &^= FBulletin | FImmediate | FPriority
	if msg.Bulletin() {
		le.Flags |= FBulletin
	}
	switch handling {
	case "IMMEDIATE":
		le.Flags |= FImmediate
	case "PRIORITY":
		le.Flags |= FPriority
	}
	if addrs, err := address.ParseList(msg.To()); err == nil && len(addrs) != 0 {
		le.ToCall = addressToLogCall(addrs[0])
	} else {
		le.ToCall = ""
	}
	slog.Info("update draft message", "dir", i.Dir, "id", ident, "lmi", le.LocalMsgID, "s", msg.Subject().EncodedSubject())
	return nil
}

// DeleteMessage deletes an unsent message from the incident.
func (i *Incident) DeleteMessage(le *LogEntry) (err error) {
	switch le.Status {
	case StatusDraft, StatusQueued, StatusHandEntered: // ok
	default:
		return errors.New("Only unsent messages and non-message log entries can be deleted.")
	}
	if fname := le.Filename(); fname != "" {
		os.Remove(filepath.Join(i.Dir, fname))
		os.Remove(filepath.Join(i.Dir, ToPDF(fname)))
	}
	le.Seq = i.Seq
	le.Status = StatusDeleted
	le.FromCall, le.FromMsgID, le.ToCall, le.ToMsgID, le.Subject, le.Flags = "", "", "", "", "", 0
	return nil
}
