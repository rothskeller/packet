package message

import (
	"fmt"
	"io"
	"maps"
	"net/textproto"
	"slices"
	"strings"

	"github.com/rothskeller/packet/message/address"
	"github.com/rothskeller/packet/message/body"
	"github.com/rothskeller/packet/message/cachetrack"
	"github.com/rothskeller/packet/message/payload"
	"github.com/rothskeller/packet/message/subject"
	"k8s.io/apimachinery/pkg/util/sets"
)

// Message is the interface satisfied by all messages.
type Message interface {
	cachetrack.CacheTracker
	MType
	// Type returns the message type for the message.
	Type() MType
	// SetType sets the message type for the message.
	SetType(MType)
	// RFC5322 returns the message encoded in RFC-5322 format for storage
	// or email transmission.
	RFC5322() string
	// To returns the list of recipients for the message.  Use
	// ParseAddressList to decode it (but note that To: lines in received
	// messages might not be syntactically correct).
	To() string
	// Subject is the subject of the message.
	Subject() subject.Subject
	// Bulletin returns whether the message is a BBS bulletin (as opposed
	// to a private message).
	Bulletin() bool
	// Draft returns whether the message is a draft.
	Draft() bool
	// Received returns whether the message has been received by the local
	// system.  It returns false for a message that has been sent or is
	// being prepared to be sent.
	Received() bool
	// Payload is the payload of the message.
	Payload() payload.Payload
	// Body is the body of the message.
	Body() body.Body
}

//-----------------------------------------------------------------------------

// common is the common parts of the message that are the same for all four
// implementations.
type common struct {
	subject subject.Subject
	payload payload.Payload
	to      string
	cachetrack.Tracker
	MType
}

func (m *common) init() {
	if m.subject.Dirty() {
		m.Tracker.MarkDirty("envelope.common.subject")
	}
	if m.payload.Dirty() {
		m.Tracker.MarkDirty("envelope.common.Payload")
	}
	m.subject.OnDirty(m.Tracker.MarkDirty)
	m.payload.OnDirty(m.Tracker.MarkDirty)
}

// Type returns the message type for the message.
func (m *common) Type() MType { return m.MType }

// SetType sets the message type for the message.
func (m *common) SetType(t MType) { m.MType = t }

// RFC5322 returns the message encoded in RFC-5322 format for storage or email
// transmission.
func (m *common) RFC5322() string { return m.rfc5322(nil) }

func (m *common) rfc5322(headers textproto.MIMEHeader) string {
	var (
		sb     strings.Builder
		hnames = sets.New(slices.Collect(maps.Keys(headers))...)
	)
	if hnames.Has("Received") {
		fmt.Fprintf(&sb, "Received: %s\r\n", headers.Get("Received"))
		hnames.Delete("Received")
	}
	if hnames.Has("From") {
		fmt.Fprintf(&sb, "From: %s\r\n", strings.Join(headers["From"], ",\r\n\t"))
		hnames.Delete("From")
	}
	if m.to != "" {
		fmt.Fprintf(&sb, "To: %s\r\n", rfc5322AddressList(m.to))
	}
	if s := m.subject.EncodedSubject(); s != "" {
		fmt.Fprintf(&sb, "Subject: %s\r\n", s)
	}
	if hnames.Has("Date") {
		fmt.Fprintf(&sb, "Date: %s\r\n", headers.Get("Date"))
		hnames.Delete("Date")
	}
	for key := range hnames {
		fmt.Fprintf(&sb, "%s: %s\r\n", key, strings.Join(headers[key], ",\r\n\t"))
	}
	io.WriteString(&sb, "\r\n")
	io.WriteString(&sb, m.payload.Encode())
	return sb.String()
}

// Subject returns the Subject of the message.
func (m *common) Subject() subject.Subject { return m.subject }

// To returns the list of recipients for the message.  Use ParseAddressList to
// decode it (but note that To: lines in received messages might not be
// syntactically correct).
func (m *common) To() string { return m.to }

// Payload returns the payload of the message.
func (m *common) Payload() payload.Payload { return m.payload }

// Body returns the body of the message.
func (m *common) Body() body.Body { return m.payload.Body() }

// common embeds three different interfaces that satisfy CacheTracker, so we
// need explicit methods to direct calls to those functions to the correct one.

// Dirty returns whether the cache is dirty (i.m., invalid).
func (m *common) Dirty() bool { return m.Tracker.Dirty() }

// MarkClean marks the cache as clean.
func (m *common) MarkClean() { m.Tracker.MarkClean() }

// OnDirty registers a function to be called when the cache becomes dirty.
func (m *common) OnDirty(fn func(string)) { m.Tracker.OnDirty(fn) }

// MarkDirty marks the cache as dirty.
func (m *common) MarkDirty(reason string) { m.Tracker.MarkDirty(reason) }

// rfc5322AddressList parses the provided string as an address list and, if
// successful, returns it reformatted into canonical format for inclusion in an
// RFC-5322 header.  If the string cannot be parsed successfully, it is
// returned unmodified.
func rfc5322AddressList(s string) string {
	if addrs, err := address.ParseList(s); err == nil {
		list := make([]string, len(addrs))
		for i, a := range addrs {
			list[i] = a.String()
		}
		return strings.Join(list, ",\r\n\t")
	}
	return s
}
