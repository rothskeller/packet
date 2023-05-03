package xscmsg

import (
	"strings"

	"github.com/rothskeller/packet/pktmsg"
)

// parseSubject parses the subject of the supplied pktmsg.Message and wraps it
// in an xscMessage with the parsed subject fields.
func parseSubject(pm pktmsg.Message) *xscMessage {
	originMsgNo, severity, handling, formtag, subject := splitSubject(pm.Subject().Value())
	var xm = xscMessage{Message: pm}
	xm.originMsgID.SetValue(originMsgNo)
	xm.severity.SetValue(severity, severityValues)
	xm.handling.SetValue(handling)
	xm.formtag = formTagField(formtag)
	xm.subject.SetValue(subject)
	return &xm
}

// xscMessage is a wrapper around a pktmsg.Message that parses the XSC-standard
// subject line format and splits it into message ID, severity, handling, form
// tag, and subject parts.
type xscMessage struct {
	pktmsg.Message
	originMsgID OriginMsgIDField
	severity    SeverityField
	handling    HandlingField
	formtag     formTagField
	subject     subjectField
}

func (m *xscMessage) Save() string {
	m.encode()
	return m.Message.Save()
}

func (m *xscMessage) Transmit() (to []string, subject string, body string) {
	m.encode()
	return m.Message.Transmit()
}

func (m *xscMessage) encode() {
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

// A formTagField is a read-only field holding the form tag from the subject
// line.  It does not need to be public.
type formTagField string

func (f formTagField) Value() string { return string(f) }

////////////////////////////////////////////////////////////////////////////////

type subjectField struct{ StringField }
