package pktmsg

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// baseRx extends outpostMessage by adding the fields for a stored-received message.
type baseRx struct {
	*outpostMessage
	rxBBS  rxBBSField
	rxArea rxAreaField
	rxDate rxDateField
}

// receivedRE is the regular expression for the "Received: " line that this
// package generates when saving a received message.
var receivedRE = regexp.MustCompile(`^FROM (\S+)\.ampr\.org BY pktmsg.local(?: FOR (\S+))?; (\w\w\w, \d\d \w\w\w \d\d\d\d \d\d:\d\d:\d\d [-+]\d\d\d\d)$`)

func parseBaseRx(om *outpostMessage, recv string) (Message, error) {
	if match := receivedRE.FindStringSubmatch(recv); match != nil {
		var m = baseRx{outpostMessage: om}
		m.rxBBS = field(match[1])
		m.rxArea = field(match[2])
		m.rxDate = dateField{field(match[3])}
		return &m, nil
	}
	// This shouldn't happen:  stored messages with a Received: header
	// should always have our Received: header format
	return om, errors.New("incorrect Received: header format for stored received message")
}

// Accessors for baseRx fields.
func (m *baseRx) RxArea() RxAreaField {
	if m.rxArea.Value() != "" {
		return m.rxArea
	}
	return nil
}
func (m *baseRx) RxBBS() RxBBSField   { return m.rxBBS }
func (m *baseRx) RxDate() RxDateField { return m.rxDate }

// Save returns the message, formatted for saving to local storage.
// Note that this can be a lossy operation; some Fields are not
// preserved in local storage.
func (m *baseRx) Save() string {
	var sb strings.Builder

	sb.WriteString("Received: FROM ")
	sb.WriteString(m.rxBBS.Value())
	sb.WriteString(".ampr.org BY pktmsg.local")
	linelen := 15 + len(m.rxBBS.Value()) + 25 + 2 + len(time.RFC1123Z)
	if m.rxArea.Value() != "" {
		sb.WriteString(" FOR ")
		sb.WriteString(m.rxArea.Value())
		linelen += 5 + len(m.rxArea.Value())
	}
	if linelen > 78 {
		sb.WriteString(";\n\t")
	} else {
		sb.WriteString("; ")
	}
	sb.WriteString(m.rxDate.Value())
	sb.WriteByte('\n')
	sb.WriteString(m.outpostMessage.Save())
	return sb.String()
}

// A dateField is a field containing a date/time string.
type dateField struct{ field }

func (f dateField) Time() (t time.Time) {
	t, _ = time.Parse(time.RFC1123Z, f.Value())
	return t
}

// The RxAreaField holds the BBS bulletin area from which the message was
// retrieved.  It is present only on received bulletin messages.
type RxAreaField interface{ Field }
type rxAreaField = field

// The RxBBSField holds the name of the BBS from which the message was
// retrieved.  It is present only on received messages.
type RxBBSField interface{ Field }
type rxBBSField = field

// The RxDateField holds the time the message was received locally.  It is
// present only on received messages.
type RxDateField interface {
	Field
	// Time returns the time encoded in the field.
	Time() time.Time
}
type rxDateField = dateField
