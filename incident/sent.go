package incident

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rothskeller/packet/v4/errors"
	"github.com/rothskeller/packet/v4/message"
	"github.com/rothskeller/packet/v4/message/address"
	"github.com/rothskeller/packet/v4/message/field"
	"github.com/rothskeller/packet/v4/message/receipt"
)

// MarkMessageSent marks a draft message as having been sent.  The parameters
// are the message and its log entry.
func (i *Incident) MarkMessageSent(dm *message.DraftMessage, le *LogEntry) (err error) {
	var (
		oldfname string
		sm       *message.SentMessage
		now      = time.Now()
	)
	// First, update the OpDate and OpTime fields in the message.
	for f := range dm.Fields() {
		switch f.Common() {
		case field.COperatorDate:
			f.SetValue(dm, now.Format("01/02/2006"))
		case field.COperatorTime:
			f.SetValue(dm, now.Format("15:04"))
		}
	}
	// Next, convert it to a SentMessage.
	sm = dm.ToSentMessage(i.Config.FromAddress(), now)
	// Update the log entry status and flags.
	oldfname = le.Filename()
	le.Seq = i.Seq
	le.Time = now
	le.Status = StatusSent
	if le.Flags&(FBulletin|FIsReceipt) == 0 {
		le.Flags |= FNeedsReceipt
	}
	// Save the message.
	if err = i.saveMessage(sm, le); err != nil {
		return err
	}
	// Update the To: address(es) in the log entry, adding additional log
	// entries if needed.
	if addrs, err := address.ParseList(sm.To()); err != nil {
		return err
	} else {
		for n, addr := range addrs {
			if n != 0 {
				c := *le
				le = &c
				le.Ident = i.nextLogIdent()
				le.Index = len(i.Log)
				i.Log = append(i.Log, le)
			}
			le.ToCall = addressToLogCall(addr)
		}
		i.sortLog()
	}
	// Remove the file associated with the unsent message.
	os.Remove(filepath.Join(i.Dir, oldfname))
	os.Remove(filepath.Join(i.Dir, strings.TrimSuffix(oldfname, ".txt")+".pdf"))
	return nil
}

// AddSentMessage takes an imported SentMessage and adds it in the incident,
// assigning it a local message ID along the way unless the message already has
// one.  Unless overwriteOK is true, the filename derived from the local
// message ID must not already exist.
func (i *Incident) AddSentMessage(msg *message.SentMessage, overwriteOK bool) (err error) {
	var (
		le       *LogEntry
		handling string
	)
	// The new message might be a receipt.  Those can't be imported.
	switch msg.Type() {
	case receipt.ReadReceipt, receipt.DeliveryReceipt:
		return errors.New("Outgoing receipts cannot be imported.")
	}
	// Create a log entry and assign a local message ID.
	le = &LogEntry{
		Ident:  i.nextLogIdent(),
		Index:  len(i.Log),
		Seq:    i.Seq,
		Status: StatusSent,
		Time:   msg.Date(),
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
					return err
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
	if le.LocalMsgID != "" {
		if _, err = os.Stat(filepath.Join(i.Dir, le.LocalMsgID+".txt")); !os.IsNotExist(err) {
			return errors.NewF("The file %s.txt already exists.", le.LocalMsgID)
		}
	} else {
		if le.LocalMsgID, err = i.nextMessageID(true); err != nil {
			return err
		}
		msg.Subject().SetSubjectMessageID(le.LocalMsgID)
		le.FromMsgID = le.LocalMsgID
	}
	le.Subject = msg.Subject().EncodedSubject()
	// Save the message.
	if err = i.saveMessage(msg, le); err != nil {
		return err
	}
	i.Log = append(i.Log, le)
	i.sortLog()
	slog.Info("import sent message", "id", le.Ident, "lid", le.LocalMsgID, "s", msg.Subject().EncodedSubject())
	return nil
}
