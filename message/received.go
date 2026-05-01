package message

import (
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"net/mail"
	"net/textproto"
	"regexp"
	"strings"
	"time"

	"github.com/rothskeller/packet/message/field"
)

// receivedRE is the regular expression for the "Received: " line that this
// package generates when saving a received message.
var receivedRE = regexp.MustCompile(`^FROM (\S*)\.(?:ampr|scc-ares-races)\.org BY (?:packet|pktmsg).local(?: ID (\S+))?(?: FOR (\S+))?; (\w\w\w, \d\d \w\w\w \d\d\d\d \d\d:\d\d:\d\d [-+]\d\d\d\d)$`)

// A ReceivedMessage is for a message that the local system received from a
// BBS at some point in the past, was stored locally on disk, and has been read
// from local storage.
type ReceivedMessage struct {
	*common
	rxBBS   string
	rxArea  string
	rxDate  time.Time
	from    string
	date    time.Time
	localID string
}

var _ Message = (*ReceivedMessage)(nil)

func (m *ReceivedMessage) receivedMessage() *ReceivedMessage { return m }

// RxBBS returns the name of the BBS from which the message was retrieved, if
// known.  (It may not be known for a manually received message.)
func (m *ReceivedMessage) RxBBS() string { return m.rxBBS }

// RxArea returns the name of the bulletin area from which the (bulletin)
// message was retrieved, or an empty string if the message is not a bulletin.
func (m *ReceivedMessage) RxArea() string { return m.rxArea }

// RxDate returns the timestamp when we retrieved the message from the BBS.
func (m *ReceivedMessage) RxDate() time.Time { return m.rxDate }

// From returns the list of senders of the message (which is virtually always
// just one sender).  Use ParseAddressList to decode it (but note that From:
// lines in received messages might not be syntactically correct).
func (m *ReceivedMessage) From() string { return m.from }

// Date returns the timestamp when the message was sent (i.m., the Date:
// header).
func (m *ReceivedMessage) Date() time.Time { return m.date }

// LocalID returns the local message ID as recorded in its Received header.
func (m *ReceivedMessage) LocalID() string { return m.localID }

// Bulletin returns whether the message is a BBS bulletin (as opposed to a
// private message).  This is determined by whether it was received from a
// bulletin area on the BBS.
func (m *ReceivedMessage) Bulletin() bool { return m.rxArea != "" }

// RFC5322 returns the message encoded in RFC-5322 format for storage or email
// transmission.
func (m *ReceivedMessage) RFC5322() string {
	hdr := make(textproto.MIMEHeader)

	rcvd := fmt.Sprintf("FROM %s.scc-ares-races.org BY packet.local", m.rxBBS)
	if m.localID != "" {
		rcvd += " ID " + m.localID
	}
	if m.rxArea != "" {
		rcvd += " FOR " + m.rxArea
	}
	rcvd += ";\n\t" + m.rxDate.Format(time.RFC1123Z)
	hdr.Set("Received", rcvd)
	if m.from != "" {
		hdr.Set("From", m.from)
	}
	hdr.Set("Date", m.date.Format(time.RFC1123Z))
	return m.rfc5322(hdr)
}

// readReceivedMessage creates a ReceivedMessage with the supplied common
// message fields and based on the supplied headers.
func readReceivedMessage(filename string, hdr mail.Header, cm *common) (_ Message, err error) {
	m := ReceivedMessage{common: cm}
	m.OnDirty(func(reason string) {
		if reason != "envelope.common.subject" {
			panic("ReceivedMessage should not change: " + reason)
		}
	})

	if match := receivedRE.FindStringSubmatch(hdr.Get("Received")); match != nil {
		m.rxBBS = match[1]
		m.localID = match[2]
		m.rxArea = match[3]
		m.rxDate, _ = mail.ParseDate(match[4])
	} else {
		// This shouldn't happen:  stored messages with a Received: header
		// should always have our Received: header format
		slog.Error("incorrect Received header", "f", filename, "Received", hdr.Get("Received"))
		return nil, errors.New("incorrect Received: header format for stored received message")
	}
	m.from = strings.Join(hdr["From"], ", ")
	if t, err := mail.ParseDate(hdr.Get("Date")); err == nil {
		m.date = t
	}
	return &m, nil
}

func (m *ReceivedMessage) Fields() iter.Seq[field.Field] {
	return func(yield func(field.Field) bool) {
		for _, f := range receivedFields {
			if !yield(f) {
				return
			}
		}
		for f := range m.Subject().Fields() {
			if !yield(f) {
				return
			}
		}
		for f := range m.Body().Fields() {
			if !yield(f) {
				return
			}
		}
	}
}

var receivedFields = []field.Field{
	field.NewField("", "From").
		Common(field.CHeaderFrom).
		ValueFunc(func(m Message) string { return rm(m).from }).
		MakeField(),
	field.NewField("", "To").
		Common(field.CHeaderTo).
		ValueFunc(func(m Message) string { return m.To() }).
		MakeField(),
	field.NewField("", "Sent").
		Common(field.CHeaderDate).
		ValueFunc(func(m Message) string { return rm(m).date.Format("01/02/2006 15:04") }).
		MakeField(),
	field.NewField("", "Received").
		Common(field.CHeaderReceived).
		ValueFunc(func(msg Message) string {
			var val string
			m := rm(msg)
			if m.rxArea != "" {
				val = "in " + m.rxArea + " "
			}
			val += "at " + m.rxDate.Format("01/02/2006 15:04")
			if m.localID != "" {
				val += " as " + m.localID
			}
			return val
		}).MakeField(),
}

func rm(m Message) *ReceivedMessage {
	return m.(interface{ receivedMessage() *ReceivedMessage }).receivedMessage()
}
