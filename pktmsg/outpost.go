package pktmsg

import (
	"encoding/base64"
	"strings"
)

type outpostMessage struct {
	*baseTx
	flags outpostFlagsField
	body  settableField
}

func newOutpostMessage() *outpostMessage {
	var m = outpostMessage{baseTx: newBaseTx()}
	return &m
}

func parseOutpostFlags(tx *baseTx) (*outpostMessage, error) {
	var (
		urg, rdr, rrr, found bool
		m                    = outpostMessage{baseTx: tx}
		body                 = tx.Body().Value()
	)
	m.body.SetValue(body)
	for {
		switch {
		case strings.HasPrefix(body, "\n"):
			// Remove newlines that might precede the first Outpost
			// code.  We'll keep this removal only if we actually
			// find an Outpost code.
			body = body[1:]
		case strings.HasPrefix(body, "!B64!"):
			// Message content has Base64 encoding.
			if dec, err := base64.StdEncoding.DecodeString(body[5:]); err == nil {
				body = string(dec)
				found = true
			} else {
				return &m, err
			}
		case strings.HasPrefix(body, "!RRR!"):
			rrr, found = true, true
			body = body[5:]
		case strings.HasPrefix(body, "!RDR!"):
			rdr, found = true, true
			body = body[5:]
		case strings.HasPrefix(body, "!URG!"):
			urg, found = true, true
			body = body[5:]
		default:
			m.flags.set(urg, rdr, rrr)
			if found {
				m.body.SetValue(body)
			} else {
				// If we didn't find any Outpost codes, leave
				// the body unchanged so as not to remove
				// initial newlines that might be significant.
			}
			return &m, nil
		}
	}
}

// Accessors for outpostMessage fields.
func (m *outpostMessage) Body() BodyField                 { return &m.body }
func (m *outpostMessage) OutpostFlags() OutpostFlagsField { return &m.flags }

func (m *outpostMessage) Save() string {
	m.encode()
	return m.baseTx.Save()
}

func (m *outpostMessage) Transmit() (to []string, subject, body string) {
	m.encode()
	return m.baseTx.Transmit()
}

func (m *outpostMessage) encode() {
	body := m.flags.Value() + m.body.Value()
	if strings.IndexFunc(body, nonASCII) >= 0 {
		body = "!B64!" + base64.StdEncoding.EncodeToString([]byte(body)) + "\n"
	}
	m.baseTx.Body().SetValue(body)
}
func nonASCII(r rune) bool {
	return r > 126 || (r < 32 && r != '\t' && r != '\n')
}

// The OutpostFlagsField holds the Outpost message flags.
type OutpostFlagsField interface {
	Field
	// Urgent returns the value of the Urgent flag.
	Urgent() bool
	// RequestDeliveryReceipt returns the value of the
	// RequestDeliveryReceipt flag.
	RequestDeliveryReceipt() bool
	// RequestReadReceipt returns the value of the RequestReadReceipt flag.
	RequestReadReceipt() bool
	// SetUrgent sets the value of the Urgent flag.
	SetUrgent(bool)
	// SetRequestDeliveryReceipt sets the value of the
	// RequestDeliveryReceipt flag.
	SetRequestDeliveryReceipt(bool)
	// SetRequestReadReceipt sets the value of the RequestReadReceipt flag.
	SetRequestReadReceipt(bool)
}

type outpostFlagsField struct{ field }

// SetValue blocks attempts to set the field value directly.  Use the SetUrgent,
// SetRequestDeliveryReceipt, and SetRequestReadReceipt methods instead.
func (f *outpostFlagsField) SetValue(value string) { panic("should not be called") }

// Urgent returns the value of the Urgent flag.
func (f outpostFlagsField) Urgent() bool {
	return strings.Contains(f.Value(), "!URG!")
}

// RequestDeliveryReceipt returns the value of the RequestDeliveryReceipt flag.
func (f outpostFlagsField) RequestDeliveryReceipt() bool {
	return strings.Contains(f.Value(), "!RDR!")
}

// RequestReadReceipt returns the value of the RequestReadReceipt flag.
func (f outpostFlagsField) RequestReadReceipt() bool {
	return strings.Contains(f.Value(), "!RRR!")
}

// SetUrgent sets the value of the Urgent flag.
func (f *outpostFlagsField) SetUrgent(want bool) {
	f.set(want, f.RequestDeliveryReceipt(), f.RequestReadReceipt())
}

// SetRequestDeliveryReceipt sets the value of the RequestDeliveryReceipt flag.
func (f *outpostFlagsField) SetRequestDeliveryReceipt(want bool) {
	f.set(f.Urgent(), want, f.RequestReadReceipt())
}

// SetRequestReadReceipt sets the value of the RequestReadReceipt flag.
func (f *outpostFlagsField) SetRequestReadReceipt(want bool) {
	f.set(f.Urgent(), f.RequestDeliveryReceipt(), want)
}

func (f *outpostFlagsField) set(urg, rdr, rrr bool) {
	var v string
	if urg {
		v += "!URG!"
	}
	if rdr {
		v += "!RDR!"
	}
	if rrr {
		v += "!RRR!"
	}
	f.field = field(v)
}
