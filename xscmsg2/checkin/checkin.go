package checkin

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/xscform"
)

// Tag identifies check-in messages.
const Tag = "Check-In"
const html = "check-in.html"

var checkInRE = regexp.MustCompile(`(?i)^Check-In\s+([A-Z][A-Z0-9]{2,5})\s*,(.*)(?:\n([AKNW][A-Z0-9]{2,5})\s*,(.*))?`)

func init() {
	xscmsg.RegisterCreate(formtype, create)
	xscmsg.RegisterType(recognize)
}

func create() *xscmsg.Message {
	return xscform.CreateForm(formtype, fieldDefs)
}

func recognize(msg *pktmsg.Message, form *pktmsg.Form) *xscmsg.Message {
	subject := xscmsg.ParseSubject(msg.Header.Get("Subject"))
	if subject == nil || subject.FormTag != "" || !strings.HasPrefix(strings.ToLower(subject.Subject), "check-in ") {
		return nil
	}
	var m = create()
	m.RawMessage = msg
	m.Field("MsgNo").Value = subject.MessageNumber
	if match := checkInRE.FindStringSubmatch(msg.Body); match != nil {
		if match[3] != "" {
			m.Field("TacCall").Value = match[1]
			m.Field("TacName").Value = strings.TrimSpace(match[2])
			m.Field("OpCall").Value = match[3]
			m.Field("OpName").Value = strings.TrimSpace(match[4])
		} else {
			m.Field("OpCall").Value = match[1]
			m.Field("OpName").Value = strings.TrimSpace(match[2])
		}
	}
	return m
}

var formtype = &xscmsg.MessageType{
	Tag:         Tag,
	Name:        "check-in message",
	Article:     "a",
	HTML:        html,
	Version:     "1.0",
	SubjectFunc: encodeSubject,
	BodyFunc:    encodeBody,
}

func encodeBody(m *xscmsg.Message) string {
	opname := m.Field("OpName").Value
	opcall := m.Field("OpCall").Value
	taccall := m.Field("TacCall").Value
	if taccall != "" {
		tacname := m.Field("TacName").Value
		return fmt.Sprintf("Check-In %s, %s\n%s, %s\n", taccall, tacname, opcall, opname)
	}
	return fmt.Sprintf("Check-In %s, %s\n", opcall, opname)
}

func encodeSubject(m *xscmsg.Message) string {
	var name, call string
	call = m.Field("TacCall").Value
	if call != "" {
		name = m.Field("TacName").Value
	} else {
		name = m.Field("OpName").Value
		call = m.Field("OpCall").Value
	}
	msgno := m.Field("MsgNo").Value
	return xscmsg.EncodeSubject(msgno, xscmsg.HandlingRoutine, "", fmt.Sprintf("Check-In %s, %s", call, name))
}

var fieldDefs = []*xscmsg.FieldDef{
	msgNoDef, tacCallDef, tacNameDef, opCallDef, opNameDef,
}

var (
	msgNoDef = &xscmsg.FieldDef{
		Tag:        "MsgNo",
		Label:      "Message Number",
		Key:        xscmsg.FOriginMsgNo,
		Validators: []xscmsg.Validator{xscform.ValidateMessageNumber},
		Flags:      xscmsg.Required,
	}
	tacCallDef = &xscmsg.FieldDef{
		Tag:   "TacCall",
		Label: "Tactical Call Sign",
		Key:   xscmsg.FTacCall,
	}
	tacNameDef = &xscmsg.FieldDef{
		Tag:        "TacName",
		Label:      "Tactical Station Name",
		Key:        xscmsg.FTacName,
		Validators: []xscmsg.Validator{validateBothOrNone},
	}
	opCallDef = &xscmsg.FieldDef{
		Tag:        "OpCall",
		Label:      "Operator Call Sign",
		Comment:    "call sign",
		Key:        xscmsg.FOpCall,
		Validators: []xscmsg.Validator{xscform.ValidateCallSign},
		Flags:      xscmsg.Required,
	}
	opNameDef = &xscmsg.FieldDef{
		Tag:   "OpName",
		Label: "Operator Name",
		Key:   xscmsg.FOpName,
		Flags: xscmsg.Required,
	}
)

func validateBothOrNone(f *xscmsg.Field, msg *xscmsg.Message, _ bool) string {
	taccall := msg.Field("TacCall").Value
	if (taccall == "") != (f.Value == "") {
		if taccall == "" {
			return "The TacName field is set but the TacCall field is not."
		}
		return "The TacCall field is set but the TacName field is not."
	}
	return ""
}
