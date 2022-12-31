package delivrcpt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
)

// Tag identifies delivery receipts.
const Tag = "DELIVERED"
const html = "delivery-receipt.html"

// deliveryReceiptRE matches the first lines of a delivery receipt message.  Its
// substrings are the local message ID, the delivery time, and the To address.
var deliveryReceiptRE = regexp.MustCompile(`^!LMI!([^!]+)!DR!(.+)\n.*\nTo: (.+)`)

func init() {
	xscmsg.RegisterCreate(Tag, create)
	xscmsg.RegisterType(recognize)
}

func create() *xscmsg.Message {
	return &xscmsg.Message{
		Type: messageType,
		Fields: []*xscmsg.Field{
			{Def: deliveredToDef},
			{Def: deliveredSubjectDef},
			{Def: localMessageIDDef},
			{Def: deliveredTimeDef},
		},
	}
}

func recognize(msg *pktmsg.Message, form *pktmsg.Form) *xscmsg.Message {
	var m = create()
	if form != nil {
		if form.FormType == html {
			m.Field("DeliveredTo").Value = form.Get("DeliveredTo")
			m.Field("DeliveredSubject").Value = form.Get("DeliveredSubject")
			m.Field("LocalMessageID").Value = form.Get("LocalMessageID")
			m.Field("DeliveredTime").Value = form.Get("DeliveredTime")
			return m
		}
		return nil
	}
	if subject := msg.Header.Get("Subject"); strings.HasPrefix(subject, "DELIVERED: ") {
		m.Field("DeliveredSubject").Value = subject[11:]
	} else {
		return nil
	}
	if match := deliveryReceiptRE.FindStringSubmatch(msg.Body); match != nil {
		m.Field("DeliveredTo").Value = match[3]
		m.Field("LocalMessageID").Value = match[1]
		m.Field("DeliveredTime").Value = match[2]
	}
	return m
}

var messageType = &xscmsg.MessageType{
	Tag:         Tag,
	Name:        "delivery receipt",
	Article:     "a",
	HTML:        html,
	SubjectFunc: encodeSubject,
	BodyFunc:    encodeBody,
}

func encodeSubject(m *xscmsg.Message) string {
	return "DELIVERED: " + m.Field("DeliveredSubject").Value
}

func encodeBody(m *xscmsg.Message, human bool) string {
	if human {
		form := &pktmsg.Form{
			FormType: html, PIFOVersion: "0", FormVersion: "0",
			Fields: []pktmsg.FormField{
				{Tag: "DeliveredTo", Value: m.Field("DeliveredTo").Value},
				{Tag: "DeliveredSubject", Value: m.Field("DeliveredSubject").Value},
				{Tag: "DeliveredTime", Value: m.Field("DeliveredTime").Value},
				{Tag: "LocalMessageID", Value: m.Field("LocalMessageID").Value},
			},
		}
		return form.Encode(nil, nil, true)
	}
	return fmt.Sprintf(
		"!LMI!%s!DR!%s\nYour Message\nTo: %s\nSubject: %s\nwas delivered on %[2]s\nRecipient's Local Message ID: %[1]s\n",
		m.Field("LocalMessageID").Value, m.Field("DeliveredTime").Value, m.Field("DeliveredTo").Value, m.Field("DeliveredSubject").Value,
	)
}

var (
	deliveredToDef = &xscmsg.FieldDef{
		Tag:        "DeliveredTo",
		Label:      "DeliveredTo",
		Validators: []xscmsg.Validator{xscmsg.ValidateRequired},
	}
	deliveredSubjectDef = &xscmsg.FieldDef{
		Tag:        "DeliveredSubject",
		Label:      "DeliveredSubject",
		Validators: []xscmsg.Validator{xscmsg.ValidateRequired},
	}
	localMessageIDDef = &xscmsg.FieldDef{
		Tag:        "LocalMessageID",
		Label:      "LocalMessageID",
		Validators: []xscmsg.Validator{xscmsg.ValidateRequired},
	}
	deliveredTimeDef = &xscmsg.FieldDef{
		Tag:        "DeliveredTime",
		Label:      "DeliveredTime",
		Validators: []xscmsg.Validator{xscmsg.ValidateRequired},
	}
)
