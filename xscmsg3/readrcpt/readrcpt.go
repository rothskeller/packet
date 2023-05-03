package readrcpt

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
)

// Tag identifies read receipts.
const Tag = "READ"
const name = "read receipt"

// Field keys for read receipt fields:
const (
	FReadTo      pktmsg.FieldKey = "READTO"
	FReadSubject pktmsg.FieldKey = "READSUBJECT"
	FReadTime    pktmsg.FieldKey = "READTIME"
)

// readReceiptRE matches the first lines of a read receipt message.  Its
// substrings are the read time and the To address.
var readReceiptRE = regexp.MustCompile(`^!RR!(.+)\n.*\n\nTo: (.+)`)

func init() {
	xscmsg.Register(Tag, name, create, recognize)
}

type readReceipt struct {
	pktmsg.Message
	readTo      readToField
	readSubject readSubjectField
	readTime    readTimeField
}

func create() xscmsg.Message {
	var m = readReceipt{Message: pktmsg.NewMessage()}
	m.readTime.SetValue(time.Now().Format("01/02/2006 15:04"))
	return &m
}

func recognize(pm pktmsg.Message) xscmsg.Message {
	var m = readReceipt{Message: pm}
	if subject := pktmsg.KeyValue(pm, pktmsg.FSubject); strings.HasPrefix(subject, "READ: ") {
		m.readSubject.SetValue(subject[6:])
	} else {
		return nil
	}
	if match := readReceiptRE.FindStringSubmatch(pktmsg.KeyValue(pm, pktmsg.FBody)); match != nil {
		m.readTo.SetValue(match[2])
		m.readTime.SetValue(match[1])
	}
	return m
}

func (readReceipt) Type() xscmsg.MessageType {
	return xscmsg.MessageType{Tag: Tag, Name: name, Article: "a"}
}
func (m readReceipt) Iterate(fn func(pktmsg.Field)) {
	m.Message.Iterate(func(f pktmsg.Field) {
		if f, ok := f.(pktmsg.KeyedField); ok {
			if key := f.Key(); key == pktmsg.FSubject || key == pktmsg.FBody {
				return
			}
		}
		fn(f)
	})
	fn(&m.readTo)
	fn(&m.readSubject)
	fn(&m.readTime)
}
func (m readReceipt) KeyedField(key pktmsg.FieldKey) pktmsg.KeyedField {
	switch key {
	case pktmsg.FSubject, pktmsg.FBody:
		return nil
	case FReadTo:
		return &m.readTo
	case FReadSubject:
		return &m.readSubject
	case FReadTime:
		return &m.readTime
	}
	return m.Message.KeyedField(key)
}
func (m readReceipt) Save() string {
	m.encode()
	return m.Message.Save()
}
func (m readReceipt) Transmit() (to []string, subject string, body string) {
	m.encode()
	return m.Message.Transmit()
}
func (m readReceipt) encode() {
	m.KeyedField(pktmsg.FSubject).SetValue("READ: " + string(m.readSubject))
	m.KeyedField(pktmsg.FBody).SetValue(fmt.Sprintf(
		"!RR!%s\nYour Message\n\nTo: %s\nSubject: %s\n\nwas read on %[1]s\n",
		m.readTime, m.readTo, m.readSubject,
	))
}

type readToField string

func (f readToField) Key() pktmsg.FieldKey   { return FReadTo }
func (f readToField) Value() string          { return string(f) }
func (f *readToField) SetValue(value string) { *f = readToField(value) }

type readSubjectField string

func (f readSubjectField) Key() pktmsg.FieldKey   { return FReadSubject }
func (f readSubjectField) Value() string          { return string(f) }
func (f *readSubjectField) SetValue(value string) { *f = readSubjectField(value) }

type readTimeField string

func (f readTimeField) Key() pktmsg.FieldKey   { return FReadTime }
func (f readTimeField) Value() string          { return string(f) }
func (f *readTimeField) SetValue(value string) { *f = readTimeField(value) }
