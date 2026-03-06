package incident

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/address"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/messageid"
	"github.com/rothskeller/packet/message/receipt"
)

// AddDraftMessage takes a DraftMessage and saves it in the incident, assigning
// it a local message ID along the way unless the message already has one.  The
// function returns the ident of the corresponding new log entry.
func (i *Incident) AddDraftMessage(msg *message.DraftMessage) (le *LogEntry, err error) {
	var (
		handling string
	)
	// The new message might be a receipt.  Those are handled specially.
	switch msg.Type() {
	case receipt.ReadReceipt:
		return nil, errors.New("This software does not support sending read receipts, only receiving them.")
	case receipt.DeliveryReceipt:
		return i.addDraftDeliveryReceipt(msg)
	}
	// Create a log entry and assign a local message ID.
	le = &LogEntry{
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
	// Put the local message ID into the message field if it has one.  Also
	// extract the handling from the message field, if any, for use in the
	// log entry.
	for f := range msg.Fields() {
		switch f.Common() {
		case field.CHandling:
			handling = f.Value(msg)
		case field.COriginMessageID:
			if le.LocalMsgID = f.Value(msg); le.LocalMsgID == "" {
				if le.LocalMsgID, err = i.nextMessageID(true); err != nil {
					return nil, err
				}
				f.SetValue(msg, le.LocalMsgID)
			}
			le.FromMsgID = le.LocalMsgID
		}
	}
	// Check the subject line for info that we didn't get from the
	// message body.  Then set log flags for the handling.
	if le.LocalMsgID == "" {
		le.LocalMsgID = msg.Subject().SubjectMessageID()
		le.FromMsgID = le.LocalMsgID
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
	// If we still don't have a local message ID, assign one.
	if le.LocalMsgID == "" {
		if le.LocalMsgID, err = i.nextMessageID(true); err != nil {
			return nil, err
		}
		msg.Subject().SetSubjectMessageID(le.LocalMsgID)
		le.FromMsgID = le.LocalMsgID
	}
	le.Subject = msg.Subject().EncodedSubject()
	// Save the message.
	if err = i.saveMessage(msg, le); err != nil {
		return nil, err
	}
	i.Log = append(i.Log, le)
	i.sortLog()
	slog.Info("create draft message", "id", le.Ident, "lid", le.LocalMsgID, "s", msg.Subject().EncodedSubject())
	return le, nil
}

func (i *Incident) addDraftDeliveryReceipt(dr *message.DraftMessage) (le *LogEntry, err error) {
	// Create a log entry and assign a local message ID.
	le = &LogEntry{
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
	if err = i.saveMessage(dr, le); err != nil {
		return nil, err
	}
	i.Log = append(i.Log, le)
	i.sortLog()
	slog.Info("create draft receipt", "id", le.Ident, "s", dr.Subject().EncodedSubject())
	return le, nil
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
	for f := range msg.Fields() {
		switch f.Common() {
		case field.COriginMessageID:
			omi = f.Value(msg)
		case field.CHandling:
			handling = f.Value(msg)
		}
	}
	if omi == "" {
		omi = msg.Subject().SubjectMessageID()
	}
	if handling == "" {
		handling = msg.Subject().SubjectHandling()
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

// ApplyDefaults applies the default field values from the incident
// configuration to the supplied draft message, overriding any values already
// contained in those fields.
func (i *Incident) ApplyDefaults(dm *message.DraftMessage) {
	// Put the local message ID, incident defaults, and operator
	// information into the message fields if it has them.  Also extract
	// the handling from the message fields, if any, for use in the log
	// entry.
	for f := range dm.Fields() {
		switch f.Common() {
		case field.CDefaultBody:
			maybeSetValue(dm, f, i.Config.DefaultBody)
		case field.CFromICSPosition:
			maybeSetValue(dm, f, i.Config.DefaultFromPos)
		case field.CFromLocation:
			maybeSetValue(dm, f, i.Config.DefaultFromLoc)
		case field.CMessageDate, field.CFormDate:
			maybeSetValue(dm, f, time.Now().Format("01/02/2006"))
		case field.COperatorCall:
			f.SetValue(dm, i.Config.OpCall)
		case field.COperatorDate, field.COperatorTime:
			f.SetValue(dm, "")
		case field.COperatorMethod:
			f.SetValue(dm, "Other")
		case field.COperatorMethodOther:
			f.SetValue(dm, "Packet")
		case field.COperatorName:
			f.SetValue(dm, i.Config.OpName)
		case field.CReceiverSender:
			f.SetValue(dm, "sender")
		case field.CTacticalCall:
			if i.Config.TacCall != "" {
				f.SetValue(dm, i.Config.TacCall)
			}
		case field.CTacticalName:
			if i.Config.TacName != "" {
				f.SetValue(dm, i.Config.TacName)
			}
		case field.CToICSPosition:
			maybeSetValue(dm, f, i.Config.DefaultToPos)
		case field.CToLocation:
			maybeSetValue(dm, f, i.Config.DefaultToLoc)
		case field.CUseTactical:
			if i.Config.TacCall != "" {
				f.SetValue(dm, "checked")
			}
		}
	}
}

func maybeSetValue(msg *message.DraftMessage, f field.Field, value string) {
	if value != "" {
		f.SetValue(msg, value)
	}
}

// ResendMessageID sets the local message ID of the supplied draft message to
// the local message ID of the supplied sent message, with the suffix changed to
// 'R'.
func (i *Incident) ResendMessageID(sentID string) (resendID string, err error) {
	var pfx, seq, _, _ = messageid.Decode(sentID, true, false)
	resendID, _ = messageid.Encode(pfx, seq, "R")

	if slices.IndexFunc(i.Log, func(le *LogEntry) bool {
		return le.LocalMsgID == resendID
	}) >= 0 {
		return "", errors.NewF("The local message ID %s is already in use.", resendID)
	}
	return resendID, nil
}
