package incident

/*
import (
	"encoding/json/jsontext"
	"errors"
	"iter"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
)

// An MData is a record describing a message (draft, sent, or received).
type MData struct {
	// Index is a unique identifier of the message within this incident.
	// It is immutable and allocated monotonically.
	Index uint `json:"i"`
	// Flags is a bitmask of flags defining the state of the message.
	Flags MDataFlag `json:"f,omitempty"`
	// Timestamp is the time at which this message was sent or received by
	// the local station.
	Timestamp time.Time `json:"ts,omitzero,format='2006-01-02T15:04'"`
	// FromCall is the callsign of the station that sent the message.
	FromCall string `json:"fc,omitempty"`
	// FromMsgID is the message ID that the sending station assigned to the
	// message.
	FromMsgID string `json:"fm,omitempty"`
	// ToCalls is the list of callsigns of the stations to which the
	// message was sent.
	ToCalls []string `json:"tcs,omitempty"`
	// ToMsgIDs is a list of message IDs assigned to the message by the
	// stations that received it.  It is parallel to ToCalls and must have
	// the same number of elements, even if some of them are empty strings.
	ToMsgIDs []string `json:"tms,omitempty"`
	// Subject is the subject of the message.
	Subject string `json:"s,omitempty"`
	// Filename is the name of the file in which the message is stored.
	Filename string `json:"fn,omitempty"`
	// Message is the parsed message content.  It is not always filled in.
	Message message.Message
}

// Clone creates a copy of the message data.
func (m *MData) Clone() (o *MData) {
	if m == nil {
		return nil
	}
	o = new(MData)
	*o = *m
	o.ToCalls = slices.Clone(m.ToCalls)
	o.ToMsgIDs = slices.Clone(m.ToMsgIDs)
	return o
}

// LocalMessageID returns the local message ID of the message: the FromMsgID
// if we sent it or the first/only ToMsgID if we received it.
func (m *MData) LocalMsgID() string {
	if m.Flags&MsgReceived == 0 {
		return m.FromMsgID
	} else if len(m.ToMsgIDs) != 0 {
		return m.ToMsgIDs[0]
	} else {
		return ""
	}
}

// RemoteMsgIDs returns the remote message IDs of the message: the ToMsgIDs if
// we sent it or the FromMsgID if we received it.
func (m *MData) RemoteMsgIDs() []string {
	if m.Flags&MsgReceived == 0 {
		return m.ToMsgIDs
	} else if m.FromMsgID != "" {
		return []string{m.FromMsgID}
	} else {
		return nil
	}
}

// GetMDataByIndex returns the message with the specified index, or nil if
// there is none.
func (inc *Incident) GetMDataByIndex(idx int) (m *MData) {
	if idx > 0 && idx < len(inc.Messages) {
		return inc.Messages[idx]
	}
	return nil
}

// GetMDatasByID returns the messages with the specified ID (which should be in
// all caps).  If localOnly is true, only messages with the specified *local*
// ID are returned.
func (inc *Incident) GetMDatasByID(msgID string, localOnly bool) (ms []*MData) {
	for _, m := range inc.Messages {
		if m == nil {
			continue
		}
		if l := m.LocalMsgID(); l == msgID {
			ms = append(ms, m)
		}
		if localOnly {
			continue
		}
		if slices.Contains(m.RemoteMsgIDs(), msgID) {
			ms = append(ms, m)
		}
	}
	return ms
}

// GetMDatasByIDSeq returns the messages with message IDs whose sequence number
// is the specified number.  Leading zeors are not significant.  If localOnly
// is true, only the local message IDs are checked.
func (inc *Incident) GetMDatasByIDSeq(seq int, localOnly bool) (ms []*MData) {
	for _, m := range inc.Messages {
		if m == nil {
			continue
		}
		if l := m.LocalMsgID(); MessageIDSeqNum(l) == seq {
			ms = append(ms, m)
			continue
		}
		if localOnly {
			continue
		}
		for _, r := range m.RemoteMsgIDs() {
			if MessageIDSeqNum(r) == seq {
				ms = append(ms, m)
				break
			}
		}
	}
	return ms
}

// MessageIDSeqNum returns the sequence number extracted from the message ID.
// It returns -1 if the message ID doesn't have a sequence number in it or has
// an invalid format.
func MessageIDSeqNum(msgID string) int {
	if _, msgID, ok := strings.Cut(msgID, "-"); !ok {
		return -1
	} else {
		if len(msgID) != 0 {
			if last := msgID[len(msgID)-1]; last >= 'A' && last <= 'Z' {
				msgID = msgID[:len(msgID)-1]
			}
		}
		if seq, err := strconv.Atoi(msgID); err != nil || seq < 1 {
			return -1
		} else {
			return seq
		}
	}
}

// AllMDatas returns an iterator on all messages in the incident in index
// order.
func (inc *Incident) AllMDatas() iter.Seq[*MData] {
	return func(yield func(*MData) bool) {
		for i := 0; i < len(inc.Messages); i++ {
			if m := inc.Messages[i]; m != nil {
				if !yield(m) {
					return
				}
			}
		}
	}
}

// AddMData adds a message to the incident, giving it the next available
// index.
func (inc *Incident) AddMData(m *MData) (err error) {
	if len(inc.Messages) == 0 {
		inc.Messages = []*MData{nil}
	}
	m.Index = uint(len(inc.Messages))
	inc.Messages = append(inc.Messages, m)
	slog.Info("recorded new message", "inc", inc.Dir, "idx", m.Index, "id", m.LocalMsgID())
	return nil
}

// UpdateMData updates an existing message in the incident.
func (inc *Incident) UpdateMData(m *MData) (err error) {
	if m.Index < 1 {
		return errors.New("no such message")
	}
	if int(m.Index) >= len(inc.Messages) {
		return errors.New("no such message")
	}
	inc.Messages[m.Index] = m
	slog.Info("recorded update of message", "inc", inc.Dir, "idx", m.Index, "id", m.LocalMsgID())
	return nil
}

// DeleteMData deletes a message from the incident.
func (inc *Incident) DeleteMData(m *MData) (err error) {
	if m.Index < 1 {
		return errors.New("no such message")
	}
	if int(m.Index) >= len(inc.Messages) || inc.Messages[m.Index] == nil {
		return errors.New("no such message")
	}
	inc.Messages[m.Index] = nil
	slog.Info("recorded deletion of message", "inc", inc.Dir, "idx", m.Index, "id", m.LocalMsgID())
	return nil
}

type MDataFlag uint

const (
	// MsgBulletin marks the message as being a bulletin.
	MsgBulletin MDataFlag = 1 << iota
	// MsgDelivered marks the (sent) message as having been delivered,
	// i.e.,, we received a delivery receipt for it.
	MsgDelivered
	// MsgFinal indicates that the message has been sent or received (i.e.,
	// is not a draft).
	MsgFinal
	// MsgImmediate marks the message as being at the IMMEDIATE handle
	// order.
	MsgImmediate
	// MsgPriority marks the message as being at the PRIORITY handling
	// order.
	MsgPriority
	// MsgQueued indicates that the (draft) message should be sent on next
	// BBS connection.
	MsgQueued
	// MsgReceived indicates that the message has been received (i.e., not
	// an outbound message).
	MsgReceived
	// MsgVoice indicates that the "message" is an ICS-309 entry for a voice
	// communication, not a packet message.
	MsgVoice
)

func (f *MDataFlag) UnmarshalJSONFrom(d *jsontext.Decoder) (err error) {
	if t, err := d.ReadToken(); err != nil {
		return err
	} else if t.Kind() != '"' {
		return errors.New("MessageFlag must be encoded as JSON string")
	} else {
		for _, r := range t.String() {
			switch r {
			case 'B':
				*f |= MsgBulletin
			case 'D':
				*f |= MsgDelivered
			case 'F':
				*f |= MsgFinal
			case 'I':
				*f |= MsgImmediate
			case 'P':
				*f |= MsgPriority
			case 'Q':
				*f |= MsgQueued
			case 'R':
				*f |= MsgReceived
			case 'V':
				*f |= MsgVoice
			default:
				return errors.New("MessageFlag contains unknown character")
			}
		}
	}
	return nil
}

func (f MDataFlag) MarshalJSONTo(e *jsontext.Encoder) (err error) {
	var s string

	if f&MsgBulletin != 0 {
		s += "B"
	}
	if f&MsgDelivered != 0 {
		s += "D"
	}
	if f&MsgFinal != 0 {
		s += "F"
	}
	if f&MsgImmediate != 0 {
		s += "I"
	}
	if f&MsgPriority != 0 {
		s += "P"
	}
	if f&MsgQueued != 0 {
		s += "Q"
	}
	if f&MsgReceived != 0 {
		s += "R"
	}
	if f&MsgVoice != 0 {
		s += "V"
	}
	return e.WriteToken(jsontext.String(s))
}
*/
