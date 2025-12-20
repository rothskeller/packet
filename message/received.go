package message

import (
	"errors"
	"fmt"
	"net/mail"
	"net/textproto"
	"regexp"
	"strings"
	"time"
)

// receivedRE is the regular expression for the "Received: " line that this
// package generates when saving a received message.
var receivedRE = regexp.MustCompile(`^FROM (\S*)\.ampr\.org BY pktmsg.local(?: FOR (\S+))?; (\w\w\w, \d\d \w\w\w \d\d\d\d \d\d:\d\d:\d\d [-+]\d\d\d\d)$`)

// A ReceivedMessage is for a message that the local system received from a
// BBS at some point in the past, was stored locally on disk, and has been read
// from local storage.
type ReceivedMessage struct {
	*common
	rxBBS  string
	rxArea string
	rxDate time.Time
	from   string
	date   time.Time
}

var _ Message = (*ReceivedMessage)(nil)

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

// Bulletin returns whether the message is a BBS bulletin (as opposed to a
// private message).  This is determined by whether it was received from a
// bulletin area on the BBS.
func (m *ReceivedMessage) Bulletin() bool { return m.rxArea != "" }

// Draft returns whether the message is a draft.
func (m *ReceivedMessage) Draft() bool { return false }

// Received returns whether the message has been received by the local system.
// It returns false for a message that has been sent or is being prepared to be
// sent.
func (m *ReceivedMessage) Received() bool { return true }

// RFC5322 returns the message encoded in RFC-5322 format for storage or email
// transmission.
func (m *ReceivedMessage) RFC5322() string {
	hdr := make(textproto.MIMEHeader)

	if m.rxArea != "" {
		hdr.Set("Received", fmt.Sprintf("FROM %s.ampr.org BY pktmsg.local FOR %s;\n\t%s",
			m.rxBBS, m.rxArea, m.rxDate.Format(time.RFC1123Z)))
	} else {
		hdr.Set("Received", fmt.Sprintf("FROM %s.ampr.org BY pktmsg.local; %s",
			m.rxBBS, m.rxDate.Format(time.RFC1123Z)))
	}
	if m.from != "" {
		hdr.Set("From", m.from)
	}
	hdr.Set("Date", m.date.Format(time.RFC1123Z))
	return m.rfc5322(hdr)
}

// readReceivedMessage creates a ReceivedMessage with the supplied common
// message fields and based on the supplied headers.
func readReceivedMessage(hdr mail.Header, cm *common) (_ Message, err error) {
	m := ReceivedMessage{common: cm}
	m.OnDirty(func(reason string) {
		panic("ReceivedMessage should not change: " + reason)
	})

	if match := receivedRE.FindStringSubmatch(hdr.Get("Received")); match != nil {
		m.rxBBS = match[1]
		m.rxArea = match[2]
		m.rxDate, _ = time.Parse(time.RFC1123Z, match[3])
	} else {
		// This shouldn't happen:  stored messages with a Received: header
		// should always have our Received: header format
		return nil, errors.New("incorrect Received: header format for stored received message")
	}
	m.from = strings.Join(hdr["From"], ", ")
	if t, err := time.Parse(time.RFC1123Z, hdr.Get("Date")); err == nil {
		m.date = t
	}
	return &m, nil
}
