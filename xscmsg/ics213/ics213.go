package ics213

import (
	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/xscform"
)

// Tag identifies ICS-213 forms.
const Tag = "ICS213"

func init() {
	xscmsg.RegisterCreate(formtype, create)
	xscmsg.RegisterType(recognize)
}

func create() *xscmsg.Message {
	return xscform.CreateForm(formtype, fieldDefsV22)
}

func recognize(msg *pktmsg.Message, form *pktmsg.Form) *xscmsg.Message {
	if form == nil || form.FormType != formtype.HTML {
		return nil
	}
	if !xscmsg.OlderVersion(form.FormVersion, "2.2") {
		return xscform.AdoptForm(formtype, fieldDefsV22, msg, form)
	}
	if xscmsg.OlderVersion(form.FormVersion, "2.0") {
		return nil
	}
	// We have a version 2.0 or 2.1 form.  The set of fields we choose
	// depends on whether it is a form we are receiving or sending.  If it
	// has a value in field 3, we'll consider it sent.  Otherwise, we'll
	// consider it received.
	if form.Get("3.") != "" {
		return xscform.AdoptForm(formtype, fieldDefsV21Tx, msg, form)
	}
	return xscform.AdoptForm(formtype, fieldDefsV21Rx, msg, form)
}

var formtype = &xscmsg.MessageType{
	Tag:         Tag,
	Name:        "ICS-213 general message form",
	Article:     "an",
	HTML:        "form-ics213.html",
	Version:     "2.2",
	SubjectFunc: encodeSubject,
	BodyFunc:    xscform.EncodeBody,
}

func encodeSubject(m *xscmsg.Message) string {
	ho, _ := xscmsg.ParseHandlingOrder(m.Field("5.").Value)
	omsgno := m.Field("MsgNo").Value
	if f := m.Field("2."); f != nil && f.Value != "" {
		omsgno = f.Value
	}
	subject := m.Field("10.").Value
	return xscmsg.EncodeSubject(omsgno, ho, m.Type.Tag, subject)
}

var fieldDefsV22 = []*xscmsg.FieldDef{
	originMessageNumberDef, destinationMessageNumberDef,
	dateDef, timeDef /* no severity */, handlingDefV22,
	takeActionDef, replyDef /* no fyi */, replyByDef,
	toICSPositionDef, fromICSPositionDef, toLocationDef, fromLocationDef, toNameDef, fromNameDef, toTelDef, fromTelDef,
	subjectDef, referenceDef, messageDef, xscform.OpRelayRcvdDef, xscform.OpRelaySentDef, recSentDef, xscform.OpCallDef,
	xscform.OpNameDef, methodDef, otherDef, xscform.OpDateDef, xscform.OpTimeDef,
}

// Versions 2.0 and 2.1 had three message number fields, only two of which were
// filled in for any given form.  This code has different sets of field
// definitions for transmitted vs. received messages, because different message
// number fields get marked as "ORIGIN" and "DESTINATION" depending on the flow
// direction.

var fieldDefsV21Rx = []*xscmsg.FieldDef{
	senderMessageNumberRxDef, myMessageNumberRxDef, receiverMessageNumberRxDef,
	dateDef, timeDef, severityDef, handlingDefV21,
	takeActionDef, replyDef, fyiDef, replyByDef,
	toICSPositionDef, fromICSPositionDef, toLocationDef, fromLocationDef, toNameDef, fromNameDef, toTelDef, fromTelDef,
	subjectDef, referenceDef, messageDef, xscform.OpRelayRcvdDef, xscform.OpRelaySentDef, recSentDef, xscform.OpCallDef,
	xscform.OpNameDef, methodDef, otherDef, xscform.OpDateDef, xscform.OpTimeDef,
}

var fieldDefsV21Tx = []*xscmsg.FieldDef{
	senderMessageNumberTxDef, myMessageNumberTxDef, receiverMessageNumberTxDef,
	dateDef, timeDef, severityDef, handlingDefV21,
	takeActionDef, replyDef, fyiDef, replyByDef,
	toICSPositionDef, fromICSPositionDef, toLocationDef, fromLocationDef, toNameDef, fromNameDef, toTelDef, fromTelDef,
	subjectDef, referenceDef, messageDef, xscform.OpRelayRcvdDef, xscform.OpRelaySentDef, recSentDef, xscform.OpCallDef,
	xscform.OpNameDef, methodDef, otherDef, xscform.OpDateDef, xscform.OpTimeDef,
}

var (
	senderMessageNumberRxDef = &xscmsg.FieldDef{
		Tag:        "2.",
		Label:      "2. Sender's Msg #",
		KeyFunc:    msgNoKeyFunc,
		Flags:      xscmsg.Readonly,
		Validators: []xscmsg.Validator{xscform.ValidateMessageNumber},
	}
	senderMessageNumberTxDef = &xscmsg.FieldDef{
		Tag:        "2.",
		Label:      "2. Sender's Msg #",
		KeyFunc:    msgNoKeyFunc,
		Flags:      xscmsg.Readonly,
		Validators: []xscmsg.Validator{xscform.ValidateMessageNumber},
	}
	myMessageNumberRxDef = &xscmsg.FieldDef{
		Tag:        "MsgNo",
		Label:      "My Msg #",
		KeyFunc:    msgNoKeyFunc,
		Validators: []xscmsg.Validator{xscform.ValidateMessageNumber},
		Flags:      xscmsg.Required,
	}
	myMessageNumberTxDef = &xscmsg.FieldDef{
		Tag:        "MsgNo",
		Label:      "My Msg #",
		KeyFunc:    msgNoKeyFunc,
		Validators: []xscmsg.Validator{xscform.ValidateMessageNumber},
		Flags:      xscmsg.Required,
	}
	receiverMessageNumberRxDef = &xscmsg.FieldDef{
		Tag:        "3.",
		Label:      "3. Receiver's Msg #",
		KeyFunc:    msgNoKeyFunc,
		Flags:      xscmsg.Readonly,
		Validators: []xscmsg.Validator{xscform.ValidateMessageNumber},
	}
	receiverMessageNumberTxDef = &xscmsg.FieldDef{
		Tag:        "3.",
		Label:      "3. Receiver's Msg #",
		KeyFunc:    msgNoKeyFunc,
		Flags:      xscmsg.Readonly,
		Validators: []xscmsg.Validator{xscform.ValidateMessageNumber},
	}
	originMessageNumberDef = &xscmsg.FieldDef{
		Tag:        "MsgNo",
		Label:      "2. Origin Msg #",
		Key:        xscmsg.FOriginMsgNo,
		Validators: []xscmsg.Validator{xscform.ValidateMessageNumber},
		Flags:      xscmsg.Required,
	}
	destinationMessageNumberDef = &xscmsg.FieldDef{
		Tag:        "3.",
		Label:      "3. Destination Msg #",
		Key:        xscmsg.FDestinationMsgNo,
		Flags:      xscmsg.Readonly,
		Validators: []xscmsg.Validator{xscform.ValidateMessageNumber},
	}
	dateDef = &xscmsg.FieldDef{
		Tag:         "1a.",
		Label:       "1. Date",
		Comment:     "MM/DD/YYYY",
		DefaultFunc: xscform.DefaultDate,
		Validators:  []xscmsg.Validator{xscform.ValidateDate},
		Flags:       xscmsg.Required,
	}
	timeDef = &xscmsg.FieldDef{
		Tag:         "1b.",
		Label:       "1. Time (24hr)",
		Comment:     "HH:MM",
		DefaultFunc: xscform.DefaultTime,
		Validators:  []xscmsg.Validator{xscform.ValidateTime},
		Flags:       xscmsg.Required,
	}
	severityDef = &xscmsg.FieldDef{
		Tag:        "4.",
		Label:      "4. Situation Severity",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"EMERGENCY", "URGENT", "OTHER"},
		Flags:      xscmsg.Required,
	}
	handlingDefV21 = &xscmsg.FieldDef{
		Tag:        "5.",
		Label:      "5. Message Handling Order",
		Key:        xscmsg.FHandling,
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"IMMEDIATE", "PRIORITY", "ROUTINE"},
		Flags:      xscmsg.Required,
	}
	handlingDefV22 = &xscmsg.FieldDef{
		Tag:        "5.",
		Label:      "5. Handling",
		Key:        xscmsg.FHandling,
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"IMMEDIATE", "PRIORITY", "ROUTINE"},
		Flags:      xscmsg.Required,
	}
	takeActionDef = &xscmsg.FieldDef{
		Tag:        "6a.",
		Label:      "6. Take Action",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Yes", "No"},
	}
	replyDef = &xscmsg.FieldDef{
		Tag:        "6b.",
		Label:      "6. Reply",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Yes", "No"},
	}
	fyiDef = &xscmsg.FieldDef{
		Tag:        "6c.",
		Label:      "6. For your information",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	replyByDef = &xscmsg.FieldDef{
		Tag:   "6d.",
		Label: "6. Reply by",
	}
	toICSPositionDef = &xscmsg.FieldDef{
		Tag:   "7.",
		Label: "7. To ICS Position",
		Key:   xscmsg.FToICSPosition,
		Flags: xscmsg.Required,
	}
	fromICSPositionDef = &xscmsg.FieldDef{
		Tag:   "8.",
		Label: "8. From ICS Position",
		Flags: xscmsg.Required,
	}
	toLocationDef = &xscmsg.FieldDef{
		Tag:   "9a.",
		Label: "9. To Location",
		Key:   xscmsg.FToLocation,
		Flags: xscmsg.Required,
	}
	fromLocationDef = &xscmsg.FieldDef{
		Tag:   "9b.",
		Label: "9. From Location",
		Flags: xscmsg.Required,
	}
	toNameDef = &xscmsg.FieldDef{
		Tag:   "ToName",
		Label: "To Name",
	}
	fromNameDef = &xscmsg.FieldDef{
		Tag:   "FmName",
		Label: "From Name",
	}
	toTelDef = &xscmsg.FieldDef{
		Tag:   "ToTel",
		Label: "To Telephone #",
	}
	fromTelDef = &xscmsg.FieldDef{
		Tag:   "FmTel",
		Label: "From Telephone #",
	}
	subjectDef = &xscmsg.FieldDef{
		Tag:   "10.",
		Label: "10. Subject",
		Key:   xscmsg.FSubject,
		Flags: xscmsg.Required,
	}
	referenceDef = &xscmsg.FieldDef{
		Tag:   "11.",
		Label: "11. Reference",
	}
	messageDef = &xscmsg.FieldDef{
		Tag:   "12.",
		Key:   xscmsg.FBody,
		Label: "12. Message",
		Flags: xscmsg.Required | xscmsg.Multiline,
	}
	recSentDef = &xscmsg.FieldDef{
		Tag:          "Rec-Sent",
		Label:        "Receiver or Sender",
		DefaultValue: "sender",
		Validators:   []xscmsg.Validator{xscform.ValidateChoices},
		Choices:      []string{"receiver", "sender"},
	}
	methodDef = &xscmsg.FieldDef{
		Tag:          "Method",
		Label:        "How Received or Sent",
		DefaultValue: "Other",
		Validators:   []xscmsg.Validator{xscform.ValidateChoices},
		Choices:      []string{"Telephone", "Dispatch Center", "EOC Radio", "FAX", "Courier", "Amateur Radio", "Other"},
	}
	otherDef = &xscmsg.FieldDef{
		Tag:          "Other",
		Label:        "Other",
		DefaultValue: "Packet",
	}
)

// msgNoKeyFunc computes the field keys for message numbers in a v2.0 or v2.1
// ICS-213 form.
func msgNoKeyFunc(f *xscmsg.Field, m *xscmsg.Message) xscmsg.FieldKey {
	// There are three cases, depending on which of the three message number
	// fields contain values.  If the "2.txmsgno" field has a value, then it
	// is the origin, the "MsgNo" field is the destination, and the
	// "3.rxmsgno" field is unused.
	if m.Field("2.").Value != "" {
		switch f.Def.Tag {
		case "2.":
			return xscmsg.FOriginMsgNo
		case "MsgNo":
			return xscmsg.FDestinationMsgNo
		default: // "3."
			return ""
		}
	}
	// If the "3.rxmsgno" field has a value, then it is the destination,
	// "MsgNo" is the origin, and "2.txmsgno" is unused.
	if m.Field("3.").Value != "" {
		switch f.Def.Tag {
		case "3.":
			return xscmsg.FDestinationMsgNo
		case "MsgNo":
			return xscmsg.FOriginMsgNo
		default: // "2."
			return ""
		}
	}
	// If neither of them has a value, then "2.txmsgno" is the "old ICS213
	// sender", "MsgNo" is the origin, and "3.rxmsgno" is the destination.
	// Code that is generating new messages can simply use origin and
	// destination as usual.  Code that is receiving messages needs to know
	// that if a "old ICS213 sender" field exists, the origin message number
	// should be moved into it before setting the destination message
	// number.
	switch f.Def.Tag {
	case "2.":
		return xscmsg.FOldICS213TxMsgNo
	case "MsgNo":
		return xscmsg.FOriginMsgNo
	default: // "3."
		return xscmsg.FDestinationMsgNo
	}
}
