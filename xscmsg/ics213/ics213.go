package ics213

import (
	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/internal/xscform"
)

// Tag identifies ICS-213 forms.
const Tag = "ICS213"

func init() {
	xscmsg.RegisterCreate(Tag, create)
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
	SubjectFunc: xscform.EncodeSubject,
	BodyFunc:    xscform.EncodeBody,
}

var fieldDefsV22 = []*xscmsg.FieldDef{
	originMessageNumberDef, destinationMessageNumberDef,
	dateDef, timeDef /* no severity */, handlingDefV22,
	takeActionDef, replyDef /* no fyi */, replyByDef,
	toICSPositionDef, fromICSPositionDef, toLocationDef, fromLocationDef, toNameDef, fromNameDef, toTelDef, fromTelDef,
	subjectDef, referenceDef, messageDef, xscform.OpRelayRcvdDef, xscform.OpRelaySentDef, recSentDef, xscform.OpCallDef,
	methodDef, xscform.OpNameDef, otherDef, xscform.OpDateDef, xscform.OpTimeDef,
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
	methodDef, xscform.OpNameDef, otherDef, xscform.OpDateDef, xscform.OpTimeDef,
}

var fieldDefsV21Tx = []*xscmsg.FieldDef{
	senderMessageNumberTxDef, myMessageNumberTxDef, receiverMessageNumberTxDef,
	dateDef, timeDef, severityDef, handlingDefV21,
	takeActionDef, replyDef, fyiDef, replyByDef,
	toICSPositionDef, fromICSPositionDef, toLocationDef, fromLocationDef, toNameDef, fromNameDef, toTelDef, fromTelDef,
	subjectDef, referenceDef, messageDef, xscform.OpRelayRcvdDef, xscform.OpRelaySentDef, recSentDef, xscform.OpCallDef,
	methodDef, xscform.OpNameDef, otherDef, xscform.OpDateDef, xscform.OpTimeDef,
}

var (
	senderMessageNumberRxDef = &xscmsg.FieldDef{
		Tag:        "2.",
		Annotation: "txmsgno",
		Label:      "2. Sender's Msg #",
		Comment:    "message-number",
		Canonical:  xscmsg.FOriginMsgNo,
		ReadOnly:   true,
		Validators: []xscmsg.Validator{xscform.ValidateMessageNumber},
	}
	senderMessageNumberTxDef = &xscmsg.FieldDef{
		Tag:        "2.",
		Annotation: "txmsgno",
		Label:      "2. Sender's Msg #",
		Comment:    "message-number",
		ReadOnly:   true,
		Validators: []xscmsg.Validator{xscform.ValidateMessageNumber},
	}
	myMessageNumberRxDef = &xscmsg.FieldDef{
		Tag:        "MsgNo",
		Label:      "My Msg #",
		Comment:    "required message-number",
		Canonical:  xscmsg.FDestinationMsgNo,
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateMessageNumber},
	}
	myMessageNumberTxDef = &xscmsg.FieldDef{
		Tag:        "MsgNo",
		Label:      "My Msg #",
		Comment:    "required message-number",
		Canonical:  xscmsg.FOriginMsgNo,
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateMessageNumber},
	}
	receiverMessageNumberRxDef = &xscmsg.FieldDef{
		Tag:        "3.",
		Annotation: "rxmsgno",
		Label:      "3. Receiver's Msg #",
		Comment:    "message-number",
		ReadOnly:   true,
		Validators: []xscmsg.Validator{xscform.ValidateMessageNumber},
	}
	receiverMessageNumberTxDef = &xscmsg.FieldDef{
		Tag:        "3.",
		Annotation: "rxmsgno",
		Label:      "3. Receiver's Msg #",
		Comment:    "message-number",
		Canonical:  xscmsg.FDestinationMsgNo,
		ReadOnly:   true,
		Validators: []xscmsg.Validator{xscform.ValidateMessageNumber},
	}
	originMessageNumberDef = &xscmsg.FieldDef{
		Tag:        "MsgNo",
		Label:      "2. Origin Msg #",
		Comment:    "required message-number",
		Canonical:  xscmsg.FOriginMsgNo,
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateMessageNumber},
	}
	destinationMessageNumberDef = &xscmsg.FieldDef{
		Tag:        "3.",
		Annotation: "rxmsgno",
		Label:      "3. Destination Msg #",
		Comment:    "message-number",
		Canonical:  xscmsg.FDestinationMsgNo,
		ReadOnly:   true,
		Validators: []xscmsg.Validator{xscform.ValidateMessageNumber},
	}
	dateDef = &xscmsg.FieldDef{
		Tag:         "1a.",
		Annotation:  "date",
		Label:       "1. Date",
		Comment:     "required date",
		Canonical:   xscmsg.FMessageDate,
		DefaultFunc: xscform.DefaultDate,
		Validators:  []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateDate},
	}
	timeDef = &xscmsg.FieldDef{
		Tag:         "1b.",
		Annotation:  "time",
		Label:       "1. Time (24hr)",
		Comment:     "required time",
		Canonical:   xscmsg.FMessageTime,
		DefaultFunc: xscform.DefaultTime,
		Validators:  []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateTime},
	}
	severityDef = &xscmsg.FieldDef{
		Tag:        "4.",
		Annotation: "severity",
		Label:      "4. Situation Severity",
		Comment:    "required: EMERGENCY, URGENT, OTHER",
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateChoices},
		Choices:    []string{"EMERGENCY", "URGENT", "OTHER"},
	}
	handlingDefV21 = &xscmsg.FieldDef{
		Tag:        "5.",
		Annotation: "handling",
		Label:      "5. Message Handling Order",
		Comment:    "required: IMMEDIATE, PRIORITY, ROUTINE",
		Canonical:  xscmsg.FHandling,
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateChoices},
		Choices:    []string{"IMMEDIATE", "PRIORITY", "ROUTINE"},
	}
	handlingDefV22 = &xscmsg.FieldDef{
		Tag:        "5.",
		Annotation: "handling",
		Label:      "5. Handling",
		Comment:    "required: IMMEDIATE, PRIORITY, ROUTINE",
		Canonical:  xscmsg.FHandling,
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateChoices},
		Choices:    []string{"IMMEDIATE", "PRIORITY", "ROUTINE"},
	}
	takeActionDef = &xscmsg.FieldDef{
		Tag:        "6a.",
		Annotation: "take-action",
		Label:      "6. Take Action",
		Comment:    "Yes, No",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Yes", "No"},
	}
	replyDef = &xscmsg.FieldDef{
		Tag:        "6b.",
		Annotation: "reply",
		Label:      "6. Reply",
		Comment:    "Yes, No",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Yes", "No"},
	}
	fyiDef = &xscmsg.FieldDef{
		Tag:        "6c.",
		Annotation: "fyi",
		Label:      "6. For your information (no action required)",
		Comment:    "boolean",
		Validators: []xscmsg.Validator{xscform.ValidateBoolean},
	}
	replyByDef = &xscmsg.FieldDef{
		Tag:        "6d.",
		Annotation: "reply-by",
		Label:      "6. Reply by",
	}
	toICSPositionDef = &xscmsg.FieldDef{
		Tag:        "7.",
		Annotation: "to-ics-position",
		Label:      "7. To ICS Position",
		Comment:    "required",
		Canonical:  xscmsg.FToICSPosition,
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	fromICSPositionDef = &xscmsg.FieldDef{
		Tag:        "8.",
		Annotation: "from-ics-position",
		Label:      "8. From ICS Position",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	toLocationDef = &xscmsg.FieldDef{
		Tag:        "9a.",
		Annotation: "to-location",
		Label:      "9. To Location",
		Comment:    "required",
		Canonical:  xscmsg.FToLocation,
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	fromLocationDef = &xscmsg.FieldDef{
		Tag:        "9b.",
		Annotation: "from-location",
		Label:      "9. From Location",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
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
		Tag:        "10.",
		Annotation: "subject",
		Label:      "10. Subject",
		Comment:    "required",
		Canonical:  xscmsg.FSubject,
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	referenceDef = &xscmsg.FieldDef{
		Tag:        "11.",
		Annotation: "reference",
		Label:      "11. Reference",
	}
	messageDef = &xscmsg.FieldDef{
		Tag:        "12.",
		Annotation: "message",
		Label:      "12. Message",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	recSentDef = &xscmsg.FieldDef{
		Tag:          "Rec-Sent",
		Label:        "Receiver or Sender",
		Comment:      "receiver, sender",
		DefaultValue: "sender",
		Validators:   []xscmsg.Validator{xscform.ValidateChoices},
		Choices:      []string{"receiver", "sender"},
	}
	methodDef = &xscmsg.FieldDef{
		Tag:          "Method",
		Label:        "How Received or Sent",
		Comment:      "Telephone, Dispatch Center, EOC Radio, FAX, Courier, Amateur Radio, Other",
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
