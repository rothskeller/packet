package incident

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/address"
	"github.com/rothskeller/packet/message/field"
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
	return nil
}
