package incident

import (
	"cmp"
	"encoding/json/jsontext"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/address"
	"github.com/rothskeller/packet/message/subject"
)

// A LogEntry represents a single line in the ICS-309 log for the incident.
// There is one LogEntry for each received or draft message.  For sent
// messages, there is one LogEntry per recipient.  There can also be LogEntries
// unrelated to a packet message.
type LogEntry struct {
	// Ident is a unique identifier of the log entry within the incident.
	// It is immutable.
	Ident int `json:"id"`
	// Index identifies the location of the log entry in the sorted list of
	// log entries.  It can change whenever the list is sorted.  The
	// LogEntry with index N is always available at incident.Log[N].
	Index int `json:"idx"`
	// Seq is the sequence number of the incident at the time this log
	// entry was last changed (either in content or in placement).
	Seq int `json:"seq"`
	// Status is the status of the log entry.
	Status LogEntryStatus `json:"st"`
	// Flags is a bitmask of flags associated with the log entry.
	Flags LogEntryFlags `json:"f"`
	// Time is the timestamp of the log entry.
	Time time.Time `json:"t"`
	// FromCall is the call sign of the station that sent the message.  It
	// will be empty for an outgoing message.  It will be the bulletin area
	// name for a received bulletin.  It may be hand-edited.
	FromCall string `json:"fc,omitempty"`
	// FromMsgID is the message ID assigned to the message by the sending
	// station.  It may be hand-edited.
	FromMsgID string `json:"fi,omitempty"`
	// LocalMsgID is the message ID assigned to the message by the local
	// station (us).  For receipt messages, it is the local message ID of
	// the message being receipted.
	LocalMsgID string `json:"li,omitempty"`
	// ToCall is the call sign of a station to which the message was sent.
	// It will be empty for a received message.  It may be hand-edited.
	ToCall string `json:"tc,omitempty"`
	// ToMsgID is the message ID assigned to the message by the receiving
	// station identified in ToCall (or by us, for a received message).  It
	// may be hand-edited.
	ToMsgID string `json:"ti,omitempty"`
	// Subject is the message subject line, or log entry description.  It
	// may be hand-edited.
	Subject string `json:"s,omitempty"`
}

// Filename returns the filename of the message associated with the log entry,
// if any.
func (e *LogEntry) Filename() string {
	switch e.Status {
	case StatusDeleted, StatusHandEntered:
		return ""
	case StatusDraft, StatusQueued:
		return fmt.Sprintf(".unsent%04d.txt", e.Ident)
	}
	if e.Flags&FIsReceipt != 0 {
		return fmt.Sprintf(".receipt%04d.txt", e.Ident)
	}
	return e.LocalMsgID + ".txt"
}

// Linkname returns the name of the symbolic link to be created for the log
// entry, if any.
func (e *LogEntry) Linkname() string {
	switch e.Status {
	case StatusReceived:
		if e.Flags&FIsReceipt != 0 || e.FromMsgID == "" {
			return ""
		}
		return e.FromMsgID + ".txt"
	case StatusSent:
		if e.Flags&FIsReceipt != 0 || e.ToMsgID == "" {
			return ""
		}
		return e.ToMsgID + ".txt"
	default:
		return ""
	}
}

// ToPDF converts a .txt filename into a .pdf filename.
func ToPDF(txt string) string { return strings.TrimSuffix(txt, ".txt") + ".pdf" }

// CompareLogEntries is a comparison function for slices.SortFunc and similar,
// that puts log entries in order by timestamp and within that by index.
func CompareLogEntries(a, b *LogEntry) int {
	if c := a.Time.Compare(b.Time); c != 0 {
		return c
	}
	return cmp.Compare(a.Ident, b.Ident)
}

// LogEntryStatus is the status of a log entry, stored in LogEntry.Status.
type LogEntryStatus string

// Values for LogEntryStatus.
const (
	// StatusDraft is a log message for an outgoing message not yet ready
	// to send.
	StatusDraft LogEntryStatus = "D"
	// StatusHandEntered is a log entry not associated with a message.
	StatusHandEntered LogEntryStatus = "H"
	// StatusQueued is a log entry for an outgoing message ready to send
	// but not yet sent.
	StatusQueued LogEntryStatus = "Q"
	// StatusReceived is a log entry for a received message.
	StatusReceived LogEntryStatus = "R"
	// StatusSent is a log entry for a sent message.
	StatusSent LogEntryStatus = "S"
	// StatusDeleted is a log entry that has been deleted.
	StatusDeleted LogEntryStatus = "X"
)

// LogEntryFlags is a bitmask of flags associated with a log entry.
type LogEntryFlags uint

// Values for LogEntryFlags.
const (
	// FBulletin is set when the log entry is for a bulletin message.
	FBulletin LogEntryFlags = 1 << iota
	// FFollowup is set when the log entry needs followup.
	FFollowup
	// FHasReceipt is set when a receipt has been received for the message
	// and recipient identified in the log entry.
	FHasReceipt
	// FImmediate is set when the log entry is for a message with immediate
	// handling order.
	FImmediate
	// FIsReceipt is set when the log entry is for a receipt message.
	FIsReceipt
	// FNeedsReceipt is set when a receipt is expected, and has not yet
	// been received, for the message and recipient identified in the log
	// entry.
	FNeedsReceipt
	// FPriority is set when the log entry is for a message with priority
	// handling order.
	FPriority
	// FUnread is set on received messages (other than receipts) that
	// haven't been read.
	FUnread
	// FVoice is set when the (hand-entered) log entry is for a voice
	// message.
	FVoice
)

func (f LogEntryFlags) MarshalJSONTo(enc *jsontext.Encoder) (err error) {
	var sb strings.Builder

	if f&FBulletin != 0 {
		sb.WriteByte('B')
	}
	if f&FFollowup != 0 {
		sb.WriteByte('F')
	}
	if f&FImmediate != 0 {
		sb.WriteByte('I')
	}
	if f&FPriority != 0 {
		sb.WriteByte('P')
	}
	if f&FIsReceipt != 0 {
		sb.WriteByte('R')
	}
	if f&FUnread != 0 {
		sb.WriteByte('U')
	}
	if f&FVoice != 0 {
		sb.WriteByte('V')
	}
	if f&FNeedsReceipt != 0 {
		sb.WriteByte('-')
	}
	if f&FHasReceipt != 0 {
		sb.WriteByte('+')
	}
	return enc.WriteToken(jsontext.String(sb.String()))
}

func (f *LogEntryFlags) UnmarshalJSONFrom(dec *jsontext.Decoder) (err error) {
	if tok, err := dec.ReadToken(); err != nil {
		return err
	} else if tok.Kind() != '"' {
		return errors.New(`key "f" must map to a string`)
	} else {
		for _, r := range tok.String() {
			switch r {
			case 'B':
				*f |= FBulletin
			case 'F':
				*f |= FFollowup
			case 'I':
				*f |= FImmediate
			case 'P':
				*f |= FPriority
			case 'R':
				*f |= FIsReceipt
			case 'U':
				*f |= FUnread
			case 'V':
				*f |= FVoice
			case '-':
				*f |= FNeedsReceipt
			case '+':
				*f |= FHasReceipt
			default:
				return fmt.Errorf(`unknown flag "%s" in key "f"`, string(r))

			}
		}
	}
	return nil
}

// AddLogEntry adds a (manual) log entry to the incident.
func (i *Incident) AddLogEntry(le *LogEntry) {
	le.Index = len(i.Log)
	le.Ident = i.nextLogIdent()
	le.Status = StatusHandEntered
	le.Seq = i.Seq
	i.Log = append(i.Log, le)
	i.sortLog()
}

// UpdateLogEntry records updates to a log entry.
func (i *Incident) UpdateLogEntry(le *LogEntry) {
	le.Seq = i.Seq
	i.sortLog()
}

// ResetLogEntry resets the data in a log entry to match the contents of a
// message.
func (i *Incident) ResetLogEntry(le *LogEntry) (err error) {
	var msg message.Message

	if msg, err = i.GetMessageFromLogEntry(le); msg == nil && err != nil {
		return err
	} else if msg == nil {
		return errors.New("Only log entries corresponding to received or sent messages can be reset.")
	}
	switch msg := msg.(type) {
	case *message.DraftMessage, *message.JustReceivedMessage:
		return errors.New("Only log entries corresponding to received or sent messages can be reset.")
	case *message.ReceivedMessage:
		le.Time = msg.RxDate()
		if msg.Bulletin() {
			le.FromCall = strings.ToUpper(msg.RxArea())
		} else {
			le.FromCall, _, _ = strings.Cut(msg.To(), "@")
			le.FromCall = strings.ToUpper(strings.TrimSpace(le.FromCall))
		}
		le.FromMsgID = ""
		for f := range msg.Fields(msg) {
			switch f.Common() {
			case "originMessageID":
				le.FromMsgID = f.Value(msg)
			}
		}
		if le.FromMsgID == "" {
			if s, ok := msg.Subject().(*subject.SCCoSubject); ok {
				le.FromMsgID = s.SubjectMessageID()
			}
		}
		le.ToCall = ""
		le.ToMsgID = le.LocalMsgID
		le.Subject = msg.Subject().EncodedSubject()
	case *message.SentMessage:
		var preceders int

		le.Time = msg.Date()
		le.FromCall = ""
		le.FromMsgID = le.LocalMsgID
		// To compute ToCall and ToMsgID, we need to know which one of
		// the log entries for the message we are, i.e., how many other
		// log entries for the same message precede this one.
		for _, e := range i.Log {
			if e.LocalMsgID == le.LocalMsgID && e.Flags&FIsReceipt == 0 && e.Ident < le.Ident {
				preceders++
			}
		}
		if addrs, err := address.ParseList(msg.To()); err == nil && len(addrs) > preceders {
			le.ToCall, le.ToMsgID = addressToLogCall(addrs[preceders]), ""
			for r := range msg.Receipts() {
				if rcptaddr, err := address.Parse(r.ReceiverAddress); err == nil && strings.EqualFold(rcptaddr.Address, addrs[preceders].Address) && r.ReceiverMessageID != "" {
					le.ToMsgID = r.ReceiverMessageID
					break
				}
			}
		} else {
			le.ToCall, le.ToMsgID = "", ""
		}
		le.Subject = msg.Subject().EncodedSubject()
		if le.Flags&(FHasReceipt|FIsReceipt) == 0 {
			le.Flags |= FNeedsReceipt
		} else {
			le.Flags &^= FNeedsReceipt
		}
	}
	i.sortLog()
	return nil
}

// DeleteLogEntry removes a hand-entered log entry.
func (i *Incident) DeleteLogEntry(le *LogEntry) (err error) {
	if le.Status != StatusHandEntered {
		return errors.New("Only manually-entered log entries can be deleted.")
	}
	le.Status = StatusDeleted
	le.Seq = i.Seq
	le.ToCall, le.ToMsgID, le.FromCall, le.FromMsgID, le.Subject = "", "", "", "", ""
	le.Flags = 0
	return nil
}

// nextLogIdent returns the next unused log entry ident number.
func (i *Incident) nextLogIdent() (ident int) {
	for _, e := range i.Log {
		ident = max(ident, e.Ident)
	}
	return ident + 1
}

// sortLog sorts the log entries, and updates the sequence number of any that
// changed position.
func (i *Incident) sortLog() {
	if slices.IsSortedFunc(i.Log, CompareLogEntries) {
		return
	}
	oldindexes := make(map[int]int, len(i.Log))
	for _, e := range i.Log {
		oldindexes[e.Ident] = e.Index
	}
	slices.SortFunc(i.Log, CompareLogEntries)
	for idx, e := range i.Log {
		e.Index = idx
		if oldindexes[e.Ident] != e.Index {
			e.Seq = i.Seq
		}
	}
}

// addressToLogCall converts an address to a call sign as it should be displayed
// in the log.
func addressToLogCall(addr *address.Address) (call string) {
	call, _, _ = strings.Cut(addr.Address, "@")
	call = strings.ToUpper(strings.TrimSpace(call))
	if len(call) > 8 {
		call = call[:8]
	}
	return call
}
