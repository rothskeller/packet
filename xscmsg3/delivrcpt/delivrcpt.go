package delivrcpt

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
)

// Tag identifies delivery receipts.
const Tag = "DELIVERED"
const name = "delivery receipt"

// Field keys for delivery receipt fields:
const (
	FDeliveredTo      pktmsg.FieldKey = "DELIVEREDTO"
	FDeliveredSubject pktmsg.FieldKey = "DELIVEREDSUBJECT"
	FLocalMessageID   pktmsg.FieldKey = "LOCALMESSAGEID"
	FDeliveredTime    pktmsg.FieldKey = "DELIVEREDTIME"
)

// deliveryReceiptRE matches the first lines of a delivery receipt message.  Its
// substrings are the local message ID, the delivery time, and the To address.
var deliveryReceiptRE = regexp.MustCompile(`^!LMI!([^!]+)!DR!(.+)\n.*\nTo: (.+)`)

func init() {
	xscmsg.Register(Tag, name, create, recognize)
}

type deliveryReceipt struct {
	pktmsg.Message
	deliveredTo      deliveredToField
	deliveredSubject deliveredSubjectField
	localMessageID   localMessageIDField
	deliveredTime    deliveredTimeField
}

func create() xscmsg.Message {
	var m = deliveryReceipt{Message: pktmsg.NewMessage()}
	m.deliveredTime.SetValue(time.Now().Format("01/02/2006 15:04"))
	return &m
}

func recognize(pm pktmsg.Message) xscmsg.Message {
	var m = deliveryReceipt{Message: pm}
	if subject := pktmsg.KeyValue(pm, pktmsg.FSubject); strings.HasPrefix(subject, "DELIVERED: ") {
		m.deliveredSubject.SetValue(subject[11:])
	} else {
		return nil
	}
	if match := deliveryReceiptRE.FindStringSubmatch(pktmsg.KeyValue(pm, pktmsg.FBody)); match != nil {
		m.deliveredTo.SetValue(match[3])
		m.localMessageID.SetValue(match[1])
		m.deliveredTime.SetValue(match[2])
	}
	return m
}

func (deliveryReceipt) Type() xscmsg.MessageType {
	return xscmsg.MessageType{Tag: Tag, Name: name, Article: "a"}
}
func (m deliveryReceipt) Iterate(fn func(pktmsg.Field)) {
	m.Message.Iterate(func(f pktmsg.Field) {
		if f, ok := f.(pktmsg.KeyedField); ok {
			if key := f.Key(); key == pktmsg.FSubject || key == pktmsg.FBody {
				return
			}
		}
		fn(f)
	})
	fn(&m.deliveredTo)
	fn(&m.deliveredSubject)
	fn(&m.localMessageID)
	fn(&m.deliveredTime)
}
func (m deliveryReceipt) KeyedField(key pktmsg.FieldKey) pktmsg.KeyedField {
	switch key {
	case pktmsg.FSubject, pktmsg.FBody:
		return nil
	case FDeliveredTo:
		return &m.deliveredTo
	case FDeliveredSubject:
		return &m.deliveredSubject
	case FLocalMessageID:
		return &m.localMessageID
	case FDeliveredTime:
		return &m.deliveredTime
	}
	return m.Message.KeyedField(key)
}
func (m deliveryReceipt) Save() string {
	m.encode()
	return m.Message.Save()
}
func (m deliveryReceipt) Transmit() (to []string, subject string, body string) {
	m.encode()
	return m.Message.Transmit()
}
func (m deliveryReceipt) encode() {
	m.KeyedField(pktmsg.FSubject).SetValue("DELIVERED: " + string(m.deliveredSubject))
	m.KeyedField(pktmsg.FBody).SetValue(fmt.Sprintf(
		"!LMI!%s!DR!%s\nYour Message\nTo: %s\nSubject: %s\nwas delivered on %[2]s\nRecipient's Local Message ID: %[1]s\n",
		m.localMessageID, m.deliveredTime, m.deliveredTo, m.deliveredSubject,
	))
}

type deliveredToField string

func (f deliveredToField) Key() pktmsg.FieldKey   { return FDeliveredTo }
func (f deliveredToField) Value() string          { return string(f) }
func (f *deliveredToField) SetValue(value string) { *f = deliveredToField(value) }

type deliveredSubjectField string

func (f deliveredSubjectField) Key() pktmsg.FieldKey   { return FDeliveredSubject }
func (f deliveredSubjectField) Value() string          { return string(f) }
func (f *deliveredSubjectField) SetValue(value string) { *f = deliveredSubjectField(value) }

type localMessageIDField string

func (f localMessageIDField) Key() pktmsg.FieldKey   { return FLocalMessageID }
func (f localMessageIDField) Value() string          { return string(f) }
func (f *localMessageIDField) SetValue(value string) { *f = localMessageIDField(value) }

type deliveredTimeField string

func (f deliveredTimeField) Key() pktmsg.FieldKey   { return FDeliveredTime }
func (f deliveredTimeField) Value() string          { return string(f) }
func (f *deliveredTimeField) SetValue(value string) { *f = deliveredTimeField(value) }
