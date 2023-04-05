package pktmsg

import (
	"net/mail"
	"strings"
	"time"
)

// baseTx is the base implementation of a pktmsg.Message.
type baseTx struct {
	fromAddr fromAddrField
	toAddrs  toAddrsField
	subject  settableField
	sentDate sentDateField
	body     settableField
}

// newBaseTx returns a new baseTx message.
func newBaseTx() *baseTx { return new(baseTx) }

// parseBaseTx returns a new baseTx message, initialized with the data from the
// supplied header and body.
func parseBaseTx(h mail.Header, body string) *baseTx {
	var m baseTx

	m.fromAddr.SetValue(h.Get("From"))
	m.parseRecipients(h)
	m.subject.SetValue(h.Get("Subject"))
	if t, err := mail.ParseDate(h.Get("Date")); err == nil {
		m.sentDate.SetTime(t)
	}
	m.body.SetValue(body)
	return &m
}

// Accessors for baseTx fields.
func (m *baseTx) Body() BodyField         { return &m.body }
func (m *baseTx) FromAddr() FromAddrField { return &m.fromAddr }
func (m *baseTx) SentDate() SentDateField { return &m.sentDate }
func (m *baseTx) Subject() SubjectField   { return &m.subject }
func (m *baseTx) ToAddrs() ToAddrsField   { return &m.toAddrs }

// baseTx messages don't have these fields, so their accessors return nil.
func (baseTx) BBSRxDate() BBSRxDateField       { return nil }
func (baseTx) FormHTML() FormHTMLField         { return nil }
func (baseTx) FormVersion() FormVersionField   { return nil }
func (baseTx) NotPlainText() NotPlainTextField { return nil }
func (baseTx) OutpostFlags() OutpostFlagsField { return nil }
func (baseTx) PIFOVersion() PIFOVersionField   { return nil }
func (baseTx) ReturnAddr() ReturnAddrField     { return nil }
func (baseTx) RxArea() RxAreaField             { return nil }
func (baseTx) RxBBS() RxBBSField               { return nil }
func (baseTx) RxDate() RxDateField             { return nil }

// parseRecipients extracts the recipient addresses from the To, Cc, and Bcc
// headers.
func (m *baseTx) parseRecipients(h mail.Header) {
	var v string
	v = strings.Join(h["To"], ", ")
	if cc := h["Cc"]; len(cc) != 0 {
		if v != "" {
			v += ", "
		}
		v += strings.Join(cc, ", ")
	}
	if bcc := h["Bcc"]; len(bcc) != 0 {
		if v != "" {
			v += ", "
		}
		v += strings.Join(bcc, ", ")
	}
	m.toAddrs.SetValue(v)
}

func (baseTx) TaggedField(string) Field         { return nil }
func (baseTx) TaggedFields(func(string, Field)) {}

// Save returns the message, formatted for saving to local storage.
// Note that this can be a lossy operation; some Fields are not
// preserved in local storage.
func (m *baseTx) Save() string {
	var sb strings.Builder

	if v := m.fromAddr.Value(); v != "" {
		sb.WriteString("From: ")
		sb.WriteString(v)
		sb.WriteByte('\n')
	}
	if v := m.toAddrs.Value(); v != "" {
		sb.WriteString("To: ")
		sb.WriteString(v)
		sb.WriteByte('\n')
	}
	if v := m.subject.Value(); v != "" {
		sb.WriteString("Subject: ")
		sb.WriteString(v)
		sb.WriteByte('\n')
	}
	if v := m.sentDate.Value(); v != "" {
		sb.WriteString("Date: ")
		sb.WriteString(v)
		sb.WriteByte('\n')
	}
	sb.WriteByte('\n')
	sb.WriteString(m.body.Value())
	if v := m.body.Value(); v != "" && v[len(v)-1] != '\n' {
		sb.WriteByte('\n')
	}
	return sb.String()
}

// Transmit returns the destination addresses, subject header, and body
// of the message, suitable for transmission through JNOS.
func (m *baseTx) Transmit() (to []string, subject string, body string) {
	return m.toAddrs.Addrs(), m.subject.Value(), m.body.Value()
}

// field is the base implementation of Field.
type field string

func (f field) Value() string { return string(f) }

// settableField is the base implementation of SettableField.
type settableField string

// A settableDateField is a field containing a date/time string.
type settableDateField struct{ settableField }

func (f settableDateField) Time() (t time.Time) {
	t, _ = time.Parse(time.RFC1123Z, f.Value())
	return t
}

func (f *settableDateField) SetTime(t time.Time) {
	if t.IsZero() {
		f.SetValue("")
	} else {
		f.SetValue(t.Format(time.RFC1123Z))
	}
}

func (f settableField) Value() string { return string(f) }
func (f *settableField) SetValue(value string) {
	*f = settableField(value)
}

// The BodyField holds the body of the message.
type BodyField interface{ SettableField }

// The FromAddrField holds the origin address (From: header).  It may contain a
// name as well as an address.
type FromAddrField interface {
	SettableField
	// Addr returns the bare address from the field, omitting any name
	// comment or angle brackets.
	Addr() string
}

type fromAddrField struct{ settableField }

func (f fromAddrField) Addr() string {
	if addr, err := mail.ParseAddress(f.Value()); err == nil {
		return addr.Address
	}
	return f.Value()
}

// The SentDateField holds the time the message was sent (Date: header).  It is
// empty for outgoing messages that have not yet been sent.
type SentDateField interface {
	SettableField
	// Time returns the time encoded in the field.
	Time() time.Time
	// SetTime sets the time encoded in the field.
	SetTime(time.Time)
}

type sentDateField = settableDateField

// The SubjectField holds the message subject (Subject: header).
type SubjectField interface{ SettableField }

// The ToAddrsField holds the list of destination addresses (To: header).
type ToAddrsField interface {
	SettableField
	// Addrs returns the bare addresses from the field.  Name comments and angle
	// brackets are omitted.
	Addrs() []string
}

type toAddrsField struct{ settableField }

func (f toAddrsField) Addrs() (addrs []string) {
	if alist, err := mail.ParseAddressList(f.Value()); err == nil {
		addrs = make([]string, len(alist))
		for i, a := range alist {
			addrs[i] = a.Address
		}
	} else {
		addrs = strings.Split(f.Value(), ",")
		for i, a := range addrs {
			addrs[i] = strings.TrimSpace(a)
		}
	}
	return addrs
}
