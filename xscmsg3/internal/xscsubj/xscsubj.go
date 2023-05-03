package xscsubj

import (
	"strings"

	"github.com/rothskeller/packet/pktmsg"
)

// Message extends pktmsg.Message, adding methods for fields derived from
// decoding the XSC-standard subject line.
type Message interface {
	pktmsg.Message

	// OriginMsgID returns the OriginMessageID field of the message.
	// It contains the message ID assigned by the sender of the message.
	OriginMsgID() OriginMsgIDField
	// Severity returns the Severity field of the message.  It contains the
	// severity of the situation to which the message pertains.  This is an
	// obsolete field.
	Severity() SeverityField
	// Handling returns the Handling field of the message.  This contains
	// the message handling order.
	Handling() HandlingField
	// FormTag returns the FormTag field of the message.  For PackItForms
	// messages, this contains a tag identifying the form type.
	FormTag() FormTagField
	// SubjectHeader returns the SubjectHeader field of the message.  This
	// contains the unparsed Subject: header line of the message.
	SubjectHeader() SubjectHeaderField
}

// ParseSubject parses the subject of the supplied pktmsg.Message and wraps it
// in an object that has the parsed subject fields.
func ParseSubject(pm pktmsg.Message) Message {
	originMsgNo, severity, handling, formtag, subject := splitSubject(pm.Subject().Value())
	var sm = baseSubj{Message: pm}
	sm.originMsgID.SetValue(originMsgNo)
	sm.severity.SetValue(severity)
	sm.handling.SetValue(handling)
	sm.formtag.SetValue(formtag)
	sm.subject.SetValue(subject)
	return &sm
}

// baseSubj is a wrapper around a pktmsg.Message that parses the XSC-standard
// subject line format and splits it into message ID, severity, handling, form
// tag, and subject parts.
type baseSubj struct {
	pktmsg.Message
	originMsgID originMsgIDField
	severity    severityField
	handling    handlingField
	formtag     formTagField
	subject     subjectField
}

func (m *baseSubj) OriginMsgID() OriginMsgIDField     { return &m.originMsgID }
func (m *baseSubj) Severity() SeverityField           { return &m.severity }
func (m *baseSubj) Handling() HandlingField           { return &m.handling }
func (m *baseSubj) FormTag() FormTagField             { return &m.formtag }
func (m *baseSubj) Subject() pktmsg.SubjectField      { return &m.subject }
func (m *baseSubj) SubjectHeader() SubjectHeaderField { return m.Message.Subject() }

func (m *baseSubj) Save() string {
	m.encode()
	return m.Message.Save()
}

func (m *baseSubj) Transmit() (to []string, subject string, body string) {
	m.encode()
	return m.Message.Transmit()
}

func (m *baseSubj) encode() {
	// We might have received a message in which the Subject: header doesn't
	// match the form contents, or doesn't parse quite correctly.  We don't
	// want to rewrite the Subject: header when storing such messages.
	if pktmsg.Finalized(m) {
		return
	}
	// If we have none of the other fields that would go in a Subject:
	// header, return only the Subject field (i.e., don't prepend a few
	// empty fields with underscores).
	if m.originMsgID.Value() == "" && m.severity.Value() == "" && m.handling.Value() == "" && m.formtag.Value() == "" {
		m.Message.Subject().SetValue(m.subject.Value())
		return
	}
	// Build the proper Subject: header.
	var sb strings.Builder
	sb.WriteString(m.originMsgID.Value())
	sb.WriteByte('_')
	if v := m.severity.Value(); v != "" {
		sb.WriteString(encodeSeverity(v))
		sb.WriteByte('/')
	}
	sb.WriteString(encodeHandling(m.handling.Value()))
	sb.WriteByte('_')
	if v := m.formtag.Value(); v != "" {
		sb.WriteString(v)
		sb.WriteByte('_')
	}
	sb.WriteString(m.subject.Value())
	m.Subject().SetValue(sb.String())
}

func splitSubject(s string) (oid, sev, han, ftag, subj string) {
	var (
		codes string
		found bool
		parts []string
	)
	if codes, subj, found = strings.Cut(s, " "); found {
		subj = " " + subj
	}
	parts = strings.SplitN(codes, "_", 4)
	switch len(parts) {
	case 0, 1, 2:
		return "", "", "", "", s
	case 3:
		subj = parts[2] + subj
	case 4:
		ftag = parts[2]
		subj = parts[3] + subj
	}
	oid = parts[0]
	if idx := strings.IndexByte(parts[1], '/'); idx >= 0 {
		sev = decodeSeverity(parts[1][:idx])
		han = decodeHandling(parts[1][idx+1:])
	} else {
		han = decodeHandling(parts[1])
	}
	return
}

////////////////////////////////////////////////////////////////////////////////

// The FormTagField holds a tag identifying the form type.
type FormTagField interface{ pktmsg.SettableField }

// A formTagField is a read-only field holding the form tag from the subject
// line.  It does not need to be public.
type formTagField struct{ StringField }

////////////////////////////////////////////////////////////////////////////////

// The HandlingField holds the message handling order.
type HandlingField interface{ pktmsg.SettableField }
type handlingField struct{ StringField }

func decodeHandling(s string) string {
	switch s {
	case "I":
		return "IMMEDIATE"
	case "P":
		return "PRIORITY"
	case "R":
		return "ROUTINE"
	}
	return s
}

func encodeHandling(s string) string {
	switch s {
	case "IMMEDIATE":
		return "I"
	case "PRIORITY":
		return "P"
	case "ROUTINE":
		return "R"
	}
	return s
}

////////////////////////////////////////////////////////////////////////////////

// The OriginMsgIDField holds the message ID assigned by the sender of the
// message.
type OriginMsgIDField interface{ pktmsg.SettableField }
type originMsgIDField struct{ StringField }

////////////////////////////////////////////////////////////////////////////////

// The SeverityField holds the severity of the situation to which the message
// pertains.  This is an obsolete field.
type SeverityField interface{ pktmsg.SettableField }
type severityField struct{ StringField }

func decodeSeverity(s string) string {
	switch s {
	case "E":
		return "EMERGENCY"
	case "U":
		return "URGENT"
	case "O":
		return "OTHER"
	}
	return s
}

func encodeSeverity(s string) string {
	switch s {
	case "EMERGENCY":
		return "E"
	case "URGENT":
		return "U"
	case "OTHER":
		return "O"
	}
	return s
}

////////////////////////////////////////////////////////////////////////////////

// A StringField is a pktmsg.SettableField that contains a string.  It is the
// base implementation for nearly all Fields.
type StringField string

// Value returns the value of the Field as a string.
func (f StringField) Value() string { return string(f) }

// SetValue sets the string value of the Field.
func (f *StringField) SetValue(value string) { *f = StringField(strings.TrimSpace(value)) }

////////////////////////////////////////////////////////////////////////////////

type subjectField struct{ StringField }

////////////////////////////////////////////////////////////////////////////////

// The SubjectHeaderField holds the unparsed Subject: header line of the
// message.
type SubjectHeaderField interface{ pktmsg.Field }
