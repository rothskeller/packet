package pktmsg

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Save returns the message encoded in a format suitable for saving to local
// storage.
func (m *Message) Save() string {
	var (
		sb   strings.Builder
		body string
	)

	if m.RxBBS != "" {
		if m.RxArea != "" {
			fmt.Fprintf(&sb, "Received: FROM %s.ampr.org BY pktmsg.local FOR %s;\n\t%s\n", m.RxBBS, m.RxArea, m.RxDate)
		} else {
			fmt.Fprintf(&sb, "Received: FROM %s.ampr.org BY pktmsg.local; %s\n", m.RxBBS, m.RxDate)
		}
	}
	if m.From != "" {
		fmt.Fprintf(&sb, "From: %s\n", m.From)
	}
	if m.To != "" {
		fmt.Fprintf(&sb, "To: %s\n", m.To)
	}
	m.encodeSubject()
	if m.SubjectHeader != "" {
		fmt.Fprintf(&sb, "Subject: %s\n", m.SubjectHeader)
	}
	if m.SentDate != "" {
		fmt.Fprintf(&sb, "Date: %s\n", m.SentDate)
	}
	sb.WriteByte('\n')
	if m.PIFOVersion != "" {
		body = m.encodeForm()
	} else {
		body = m.Body
	}
	sb.WriteString(m.encodeOutpost(body))
	return sb.String()
}

// Transmit returns the destination addresses, subject header, and body
// of the message, suitable for transmission through JNOS.
func (m *Message) Transmit() (to []string, subject string, body string) {
	m.encodeSubject()
	if m.PIFOVersion != "" {
		body = m.encodeForm()
	} else {
		body = m.Body
	}
	return m.To.Addresses(), m.SubjectHeader, m.encodeOutpost(body)
}

func (m *Message) encodeForm() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "!SCCoPIFO!\n#T: %s\n#V: %s-%s\n", m.FormHTML, m.PIFOVersion, m.FormVersion)
	for _, f := range m.TaggedFields {
		sb.WriteString(f.pifoEncode())
	}
	sb.WriteString("!/ADDON!\n")
	return sb.String()
}

var quoteSCCoPIFO = strings.NewReplacer(`\`, `\\`, "\n", `\n`, "]", "`]")

func (t TaggedField) pifoEncode() string {
	if t.Value == "" {
		return "" // omit empty fields
	}
	value := quoteSCCoPIFO.Replace(t.Value)
	if strings.HasSuffix(value, "`") {
		value += "]]"
	}
	enc := fmt.Sprintf("%s: [%s]", t.Tag, value)
	var wrapped string
	for len(enc) > 128 {
		wrapped += enc[:128] + "\n"
		enc = enc[128:]
	}
	return wrapped + enc + "\n"
}

func (m *Message) encodeOutpost(body string) string {
	needB64 := strings.IndexFunc(body, nonASCII) >= 0
	if !needB64 && !m.OutpostUrgent && !m.RequestDeliveryReceipt && !m.RequestReadReceipt {
		return body
	}
	var sb strings.Builder
	if m.OutpostUrgent {
		sb.WriteString("!URG!")
	}
	if m.RequestDeliveryReceipt {
		sb.WriteString("!RDR!")
	}
	if m.RequestReadReceipt {
		sb.WriteString("!RRR!")
	}
	sb.WriteString(body)
	body = sb.String()
	if needB64 {
		return "!B64!" + base64.StdEncoding.EncodeToString([]byte(body)) + "\n"
	}
	return body
}
func nonASCII(r rune) bool {
	return r > 126 || (r < 32 && r != '\t' && r != '\n')
}

func (m *Message) encodeSubject() {
	// We might have received a message in which the Subject: header doesn't
	// match the form contents, or doesn't parse quite correctly.  We don't
	// want to rewrite the Subject: header when storing such messages.
	if m.Finalized() {
		return
	}
	// If we have none of the other fields that would go in a Subject:
	// header, return only the Subject field (i.e., don't prepend a few
	// empty fields with underscores).
	if m.OriginMsgID == "" && m.Severity == "" && m.Handling == "" && m.FormTag == "" {
		m.SubjectHeader = m.Subject
		return
	}
	// Build the proper Subject: header.
	var sb strings.Builder
	sb.WriteString(m.OriginMsgID)
	sb.WriteByte('_')
	if m.Severity != "" {
		sb.WriteString(encodeSeverity(m.Severity))
		sb.WriteByte('/')
	}
	sb.WriteString(encodeHandling(m.Handling))
	sb.WriteByte('_')
	if m.FormTag != "" {
		sb.WriteString(m.FormTag)
		sb.WriteByte('_')
	}
	sb.WriteString(m.Subject)
	m.SubjectHeader = sb.String()
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
