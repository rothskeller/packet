package incident

import (
	"cmp"
	"encoding/json/jsontext"
	"errors"
	"fmt"
	"strings"
	"time"
)

// A LogEntry represents a single line in the ICS-309 log for the incident.
// There is one LogEntry for each received or draft message.  For sent
// sent messages, there is one LogEntry per recipient.  There can also be
// LogEntries unrelated to a packet message.
type LogEntry struct {
	// Index is a unique identifier of the log entry within the incident.
	// Indexes are positive integers that grow monotonically and are never
	// reused.  The index space may be sparse since log entries can be
	// deleted.
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
	// Filename is the name of the file, within the incident directory,
	// containing the message described by the log entry.  It is empty for
	// hand-edited log entries.
	Filename string `json:"fn,omitempty"`
	// Linkname is the name of the symbolic link, within the incident
	// directory, pointing to the message described by the log entry, named
	// with the remote message ID.  It is empty if there is no remote
	// message ID.
	Linkname string `json:"ln,omitempty"`
}

// CompareLogEntries is a comparison function for slices.SortFunc and similar,
// that puts log entries in order by timestamp and within that by index.
func CompareLogEntries(a, b *LogEntry) int {
	if c := a.Time.Compare(b.Time); c != 0 {
		return c
	}
	return cmp.Compare(a.Index, b.Index)
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
