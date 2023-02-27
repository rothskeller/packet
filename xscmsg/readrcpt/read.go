package readrcpt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
)

// Tag identifies read receipts.
const Tag = "READ"
const html = "read-receipt.html"

// readReceiptRE matches the first lines of a read receipt message.  Its
// substrings are the read time and the To address.
var readReceiptRE = regexp.MustCompile(`^!RR!(.+)\n.*\n\nTo: (.+)`)

func init() {
	xscmsg.RegisterCreate(messageType, create)
	xscmsg.RegisterType(recognize)
}

func create() *xscmsg.Message {
	return &xscmsg.Message{
		Type: messageType,
		Fields: []*xscmsg.Field{
			{Def: readToDef},
			{Def: readSubjectDef},
			{Def: readTimeDef},
		},
	}
}

func recognize(msg *pktmsg.Message, form *pktmsg.Form) *xscmsg.Message {
	var m = create()
	if subject := msg.Header.Get("Subject"); strings.HasPrefix(subject, "READ: ") {
		m.Field("ReadSubject").Value = subject[11:]
	} else {
		return nil
	}
	if match := readReceiptRE.FindStringSubmatch(msg.Body); match != nil {
		m.Field("ReadTo").Value = match[2]
		m.Field("ReadTime").Value = match[1]
	}
	return m
}

var messageType = &xscmsg.MessageType{
	Tag:         Tag,
	Name:        "read receipt",
	Article:     "a",
	HTML:        html,
	SubjectFunc: encodeSubject,
	BodyFunc:    encodeBody,
}

func encodeSubject(m *xscmsg.Message) string {
	return "Read: " + m.Field("ReadSubject").Value
}

func encodeBody(m *xscmsg.Message) string {
	return fmt.Sprintf("!RR!%s\nYour Message\n\nTo: %s\nSubject: %s\n\nwas read on %[1]s\n",
		m.Field("ReadTime").Value, m.Field("ReadTo").Value, m.Field("ReadSubject").Value)
}

var (
	readToDef = &xscmsg.FieldDef{
		Tag:   "ReadTo",
		Label: "ReadTo",
		Flags: xscmsg.Required,
	}
	readSubjectDef = &xscmsg.FieldDef{
		Tag:   "ReadSubject",
		Label: "ReadSubject",
		Flags: xscmsg.Required,
	}
	readTimeDef = &xscmsg.FieldDef{
		Tag:   "ReadTime",
		Label: "ReadTime",
		Flags: xscmsg.Required,
	}
)
