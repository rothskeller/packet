package message

import (
	"net/mail"
	"net/textproto"
	"strings"
	"time"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message/address"
	"github.com/rothskeller/packet/message/payload"
	"github.com/rothskeller/packet/message/subject"
)

var ErrInvalidTo = errors.New("The \"To\" address list is not valid.")

// A DraftMessage is an envelope for a message that the local system is
// preparing to send but has not yet sent.
type DraftMessage struct {
	*common
	bulletin    bool
	readyToSend bool
}

var _ Message = (*DraftMessage)(nil)

// NewDraftMessage creates a new DraftMessage with the specified MType, Subject,
// and Payload.  It is marked as a bulletin if the bulletin flag is set.  It is
// not initially marked ReadyToSend, and has no initial To address.
func NewDraftMessage(mtype MType, subject subject.Subject, payload payload.Payload, bulletin bool) (m *DraftMessage) {
	m = &DraftMessage{
		common: &common{
			MType:   mtype,
			subject: subject,
			payload: payload,
		},
		bulletin: bulletin,
	}
	m.init()
	return m
}

// Draft returns whether the message is a draft.
func (m *DraftMessage) Draft() bool { return true }

// Received returns whether the message has been received by the local system.
// It returns false for a message that has been sent or is being prepared to be
// sent.
func (m *DraftMessage) Received() bool { return false }

// Bulletin returns whether the message will be a BBS bulletin (as opposed to a
// private message).
func (m *DraftMessage) Bulletin() bool { return m.bulletin }

// ReadyToSend returns whether the message is ready to be sent during the next
// BBS connection.
func (m *DraftMessage) ReadyToSend() bool { return m.readyToSend }

// SetReadyToSend sets whether the draft message is ready to be sent during the
// next BBS connection.
func (m *DraftMessage) SetReadyToSend(ready bool) {
	if m.readyToSend != ready {
		m.readyToSend = ready
		m.MarkDirty("envelope.DraftMessage.ReadyToSend")
	}
}

// SetTo sets the list of recipients for the message.
func (m *common) SetTo(to string) (err error) {
	to = strings.ReplaceAll(to, "\n", " ")
	if _, err = address.ParseList(to); err != nil {
		err = ErrInvalidTo
	}
	m.to = to
	m.MarkDirty("envelope.DraftMessage.To")
	return err
}

// RFC5322 returns the message encoded in RFC-5322 format for storage or email
// transmission.
func (m *DraftMessage) RFC5322() string {
	hdr := make(textproto.MIMEHeader)

	if m.bulletin {
		hdr.Set("X-Packet-Bulletin", "true")
	}
	if m.readyToSend {
		hdr.Set("X-Packet-Queued", "true")
	}
	return m.rfc5322(hdr)
}

// readDraftMessage creates a DraftMessage with the supplied common message
// fields and based on the supplied headers.
func readDraftMessage(hdr mail.Header, c *common) (_ Message, err error) {
	m := DraftMessage{common: c}

	m.bulletin = hdr.Get("X-Packet-Bulletin") != ""
	m.readyToSend = hdr.Get("X-Packet-Queued") != ""
	return &m, nil
}

// ToSentMessage converts a DraftMessage to a SentMessage.  It would be
// called when the message has been sent.
func (m *DraftMessage) ToSentMessage(from string, date time.Time) (s *SentMessage) {
	s = &SentMessage{
		common:   m.common,
		from:     from,
		date:     date,
		bulletin: m.bulletin,
	}
	s.init()
	return s
}
