package message

import (
	"fmt"
	"io"
	"maps"
	"net/textproto"
	"slices"
	"strings"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message/address"
	"github.com/rothskeller/packet/message/body"
	"github.com/rothskeller/packet/message/cachetrack"
	"github.com/rothskeller/packet/message/msgifc"
	"github.com/rothskeller/packet/message/payload"
	"github.com/rothskeller/packet/message/subject"
	"k8s.io/apimachinery/pkg/util/sets"
)

// Message is the interface satisfied by all messages.
type Message = msgifc.Message
type ValidateFlags = msgifc.ValidateFlags

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
		m.Tracker.MarkDirty("envelope.common.Subject")
	}
	if m.payload.Dirty() {
		m.Tracker.MarkDirty("envelope.common.Payload")
	}
	m.subject.OnDirty(m.Tracker.MarkDirty)
	m.payload.OnDirty(m.Tracker.MarkDirty)
}

func (m *common) Type() MType              { return m.MType }
func (m *common) SetType(t MType)          { m.MType = t }
func (m *common) Subject() subject.Subject { return m.subject }
func (m *common) To() string               { return m.to }
func (m *common) Payload() payload.Payload { return m.payload }
func (m *common) Body() body.Body          { return m.payload.Body() }

func (m *common) SetSubject(s subject.Subject) {
	m.subject = s
	if m.subject.Dirty() {
		m.Tracker.MarkDirty("envelope.common.subject")
	}
	m.subject.OnDirty(m.Tracker.MarkDirty)
}

func (m *common) RFC5322() string { return m.rfc5322(nil) }

func (m *common) rfc5322(headers textproto.MIMEHeader) string {
	var (
		sb     strings.Builder
		hnames = sets.New(slices.Collect(maps.Keys(headers))...)
	)
	if hnames.Has("Received") {
		fmt.Fprintf(&sb, "Received: %s\n", headers.Get("Received"))
		hnames.Delete("Received")
	}
	if hnames.Has("From") {
		fmt.Fprintf(&sb, "From: %s\n", strings.Join(headers["From"], ",\n\t"))
		hnames.Delete("From")
	}
	if m.to != "" {
		fmt.Fprintf(&sb, "To: %s\n", rfc5322AddressList(m.to))
	}
	if s := m.subject.EncodedSubject(); s != "" {
		fmt.Fprintf(&sb, "Subject: %s\n", s)
	}
	if hnames.Has("Date") {
		fmt.Fprintf(&sb, "Date: %s\n", headers.Get("Date"))
		hnames.Delete("Date")
	}
	for key := range hnames {
		fmt.Fprintf(&sb, "%s: %s\n", key, strings.Join(headers[key], ",\n\t"))
	}
	io.WriteString(&sb, "\n")
	io.WriteString(&sb, m.payload.Encode())
	return sb.String()
}

// common embeds three different interfaces that satisfy CacheTracker, so we
// need explicit methods to direct calls to those functions to the correct one.

func (m *common) Dirty() bool             { return m.Tracker.Dirty() }
func (m *common) MarkClean()              { m.Tracker.MarkClean() }
func (m *common) OnDirty(fn func(string)) { m.Tracker.OnDirty(fn) }
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

// Validate validates an entire message.  It returns a join of all validation
// errors, or nil if there are none.
func ValidateMessage(msg Message, flags msgifc.ValidateFlags) (err error) {
	for f := range msg.Fields() {
		err = errors.Join(err, f.Validate(msg, f, flags))
	}
	return err
}
