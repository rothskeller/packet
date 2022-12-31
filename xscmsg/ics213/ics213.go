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

func create() xscmsg.Message {
	return xscform.CreateForm(formtype22, makeFields22())
}

func recognize(msg *pktmsg.Message, form *pktmsg.Form) xscmsg.Message {
	if form == nil || form.FormType != formtype22.HTML {
		return nil
	}
	if !xscmsg.OlderVersion(form.FormVersion, "2.2") {
		return xscform.AdoptForm(formtype22, makeFields22(), msg, form)
	}
	if xscmsg.OlderVersion(form.FormVersion, "2.0") {
		return nil
	}
	// We have a version 2.0 or 2.1 form.  The set of fields we choose
	// depends on whether it is a form we are receiving or sending.  If it
	// has a value in field 3, we'll consider it sent.  Otherwise, we'll
	// consider it received.
	if form.Get("3.") != "" {
		return xscform.AdoptForm(formtype20, makeFields20Tx(), msg, form)
	}
	return xscform.AdoptForm(formtype20, makeFields20Rx(), msg, form)
}

var formtype22 = &xscmsg.MessageType{
	Tag:     Tag,
	Name:    "ICS-213 general message form",
	Article: "an",
	HTML:    "form-ics213.html",
	Version: "2.2",
}

func makeFields22() []xscmsg.Field {
	return []xscmsg.Field{
		&xscform.MessageNumberField{Field: *xscform.NewField(originMessageNumberID, true)},
		&xscform.MessageNumberField{Field: *xscform.NewField(destinationMessageNumberID, false)},
		&xscform.DateFieldDefaultNow{DateField: xscform.DateField{Field: *xscform.NewField(dateID, true)}},
		&xscform.TimeFieldDefaultNow{TimeField: xscform.TimeField{Field: *xscform.NewField(timeID, true)}},
		&xscform.ChoicesField{Field: *xscform.NewField(handlingID22, true), Choices: handlingChoices},
		&xscform.ChoicesField{Field: *xscform.NewField(takeActionID, false), Choices: yesNoChoices},
		&xscform.ChoicesField{Field: *xscform.NewField(replyID, false), Choices: yesNoChoices},
		xscform.NewField(replyByID, false),
		xscform.NewField(toICSPositionID, true),
		xscform.NewField(fromICSPositionID, true),
		xscform.NewField(toLocationID, true),
		xscform.NewField(fromLocationID, true),
		xscform.NewField(toNameID, false),
		xscform.NewField(fromNameID, false),
		xscform.NewField(toTelID, false),
		xscform.NewField(fromTelID, false),
		xscform.NewField(subjectID, true),
		xscform.NewField(referenceID, false),
		xscform.NewField(messageID, true),
		xscform.FOpRelayRcvd(),
		xscform.FOpRelaySent(),
		&recSentField{xscform.ChoicesField{Field: *xscform.NewField(recSentID, false), Choices: recSentChoices}},
		xscform.FOpCall(),
		&methodField{xscform.ChoicesField{Field: *xscform.NewField(methodID, false), Choices: methodChoices}},
		xscform.FOpName(),
		&otherField{*xscform.NewField(otherID, false)},
		xscform.FOpDate(),
		xscform.FOpTime(),
	}
}

var formtype20 = &xscmsg.MessageType{
	Tag:     Tag,
	Name:    "ICS-213 general message form",
	Article: "an",
	HTML:    "form-ics213.html",
	Version: "2.1",
}

func makeFields20Rx() []xscmsg.Field {
	return []xscmsg.Field{
		&xscform.MessageNumberField{Field: *xscform.NewField(senderMessageNumberRxID, false)},
		&xscform.MessageNumberField{Field: *xscform.NewField(myMessageNumberRxID, true)},
		&xscform.MessageNumberField{Field: *xscform.NewField(receiverMessageNumberRxID, false)},
		&xscform.DateFieldDefaultNow{DateField: xscform.DateField{Field: *xscform.NewField(dateID, true)}},
		&xscform.TimeFieldDefaultNow{TimeField: xscform.TimeField{Field: *xscform.NewField(timeID, true)}},
		&xscform.ChoicesField{Field: *xscform.NewField(severityID, true), Choices: severityChoices},
		&xscform.ChoicesField{Field: *xscform.NewField(handlingID20, true), Choices: handlingChoices},
		&xscform.ChoicesField{Field: *xscform.NewField(takeActionID, false), Choices: yesNoChoices},
		&xscform.ChoicesField{Field: *xscform.NewField(replyID, false), Choices: yesNoChoices},
		&xscform.BooleanField{Field: *xscform.NewField(fyiID, false)},
		xscform.NewField(replyByID, false),
		xscform.NewField(toICSPositionID, true),
		xscform.NewField(fromICSPositionID, true),
		xscform.NewField(toLocationID, true),
		xscform.NewField(fromLocationID, true),
		xscform.NewField(toNameID, false),
		xscform.NewField(fromNameID, false),
		xscform.NewField(toTelID, false),
		xscform.NewField(fromTelID, false),
		xscform.NewField(subjectID, true),
		xscform.NewField(referenceID, false),
		xscform.NewField(messageID, true),
		xscform.FOpRelayRcvd(),
		xscform.FOpRelaySent(),
		&recSentField{xscform.ChoicesField{Field: *xscform.NewField(recSentID, false), Choices: recSentChoices}},
		xscform.FOpCall(),
		&methodField{xscform.ChoicesField{Field: *xscform.NewField(methodID, false), Choices: methodChoices}},
		xscform.FOpName(),
		&otherField{*xscform.NewField(otherID, false)},
		xscform.FOpDate(),
		xscform.FOpTime(),
	}
}

func makeFields20Tx() []xscmsg.Field {
	return []xscmsg.Field{
		&xscform.MessageNumberField{Field: *xscform.NewField(senderMessageNumberTxID, false)},
		&xscform.MessageNumberField{Field: *xscform.NewField(myMessageNumberTxID, true)},
		&xscform.MessageNumberField{Field: *xscform.NewField(receiverMessageNumberTxID, false)},
		&xscform.DateFieldDefaultNow{DateField: xscform.DateField{Field: *xscform.NewField(dateID, true)}},
		&xscform.TimeFieldDefaultNow{TimeField: xscform.TimeField{Field: *xscform.NewField(timeID, true)}},
		&xscform.ChoicesField{Field: *xscform.NewField(severityID, true), Choices: severityChoices},
		&xscform.ChoicesField{Field: *xscform.NewField(handlingID20, true), Choices: handlingChoices},
		&xscform.ChoicesField{Field: *xscform.NewField(takeActionID, false), Choices: yesNoChoices},
		&xscform.ChoicesField{Field: *xscform.NewField(replyID, false), Choices: yesNoChoices},
		&xscform.BooleanField{Field: *xscform.NewField(fyiID, false)},
		xscform.NewField(replyByID, false),
		xscform.NewField(toICSPositionID, true),
		xscform.NewField(fromICSPositionID, true),
		xscform.NewField(toLocationID, true),
		xscform.NewField(fromLocationID, true),
		xscform.NewField(toNameID, false),
		xscform.NewField(fromNameID, false),
		xscform.NewField(toTelID, false),
		xscform.NewField(fromTelID, false),
		xscform.NewField(subjectID, true),
		xscform.NewField(referenceID, false),
		xscform.NewField(messageID, true),
		xscform.FOpRelayRcvd(),
		xscform.FOpRelaySent(),
		&recSentField{xscform.ChoicesField{Field: *xscform.NewField(recSentID, false), Choices: recSentChoices}},
		xscform.FOpCall(),
		&methodField{xscform.ChoicesField{Field: *xscform.NewField(methodID, false), Choices: methodChoices}},
		xscform.FOpName(),
		&otherField{*xscform.NewField(otherID, false)},
		xscform.FOpDate(),
		xscform.FOpTime(),
	}
}

var (
	senderMessageNumberRxID = &xscmsg.FieldID{
		Tag:        "2.",
		Annotation: "txmsgno",
		Label:      "2. Sender's Msg #",
		Comment:    "message-number",
		Canonical:  xscmsg.FOriginMsgNo,
		ReadOnly:   true,
	}
	myMessageNumberRxID = &xscmsg.FieldID{
		Tag:       "MsgNo",
		Label:     "My Msg #",
		Comment:   "message-number",
		Canonical: xscmsg.FDestinationMsgNo,
	}
	receiverMessageNumberRxID = &xscmsg.FieldID{
		Tag:        "3.",
		Annotation: "rxmsgno",
		Label:      "3. Receiver's Msg #",
		Comment:    "message-number",
		ReadOnly:   true,
	}
	senderMessageNumberTxID = &xscmsg.FieldID{
		Tag:        "2.",
		Annotation: "txmsgno",
		Label:      "2. Sender's Msg #",
		Comment:    "message-number",
		ReadOnly:   true,
	}
	myMessageNumberTxID = &xscmsg.FieldID{
		Tag:       "MsgNo",
		Label:     "My Msg #",
		Comment:   "message-number",
		Canonical: xscmsg.FOriginMsgNo,
	}
	receiverMessageNumberTxID = &xscmsg.FieldID{
		Tag:        "3.",
		Annotation: "rxmsgno",
		Label:      "3. Receiver's Msg #",
		Comment:    "message-number",
		Canonical:  xscmsg.FDestinationMsgNo,
		ReadOnly:   true,
	}
	originMessageNumberID = &xscmsg.FieldID{
		Tag:       "MsgNo",
		Label:     "2. Origin Msg #",
		Comment:   "required message-number",
		Canonical: xscmsg.FOriginMsgNo,
	}
	destinationMessageNumberID = &xscmsg.FieldID{
		Tag:        "3.",
		Annotation: "rxmsgno",
		Label:      "3. Destination Msg #",
		Comment:    "message-number",
		Canonical:  xscmsg.FDestinationMsgNo,
		ReadOnly:   true,
	}
	dateID = &xscmsg.FieldID{
		Tag:        "1a.",
		Annotation: "date",
		Label:      "1. Date",
		Comment:    "required date",
		Canonical:  xscmsg.FMessageDate,
	}
	timeID = &xscmsg.FieldID{
		Tag:        "1b.",
		Annotation: "time",
		Label:      "1. Time (24hr)",
		Comment:    "required time",
		Canonical:  xscmsg.FMessageTime,
	}
	severityID = &xscmsg.FieldID{
		Tag:        "4.",
		Annotation: "severity",
		Label:      "4. Situation Severity",
		Comment:    "required: EMERGENCY, URGENT, OTHER",
	}
	severityChoices = []string{"EMERGENCY", "URGENT", "OTHER"}
	handlingID20    = &xscmsg.FieldID{
		Tag:        "5.",
		Annotation: "handling",
		Label:      "5. Message Handling Order",
		Comment:    "required: IMMEDIATE, PRIORITY, ROUTINE",
		Canonical:  xscmsg.FHandling,
	}
	handlingID22 = &xscmsg.FieldID{
		Tag:        "5.",
		Annotation: "handling",
		Label:      "5. Handling",
		Comment:    "required: IMMEDIATE, PRIORITY, ROUTINE",
		Canonical:  xscmsg.FHandling,
	}
	handlingChoices = []string{"IMMEDIATE", "PRIORITY", "ROUTINE"}
	takeActionID    = &xscmsg.FieldID{
		Tag:        "6a.",
		Annotation: "take-action",
		Label:      "6. Take Action",
		Comment:    "Yes, No",
	}
	yesNoChoices = []string{"Yes", "No"}
	replyID      = &xscmsg.FieldID{
		Tag:        "6b.",
		Annotation: "reply",
		Label:      "6. Reply",
		Comment:    "Yes, No",
	}
	fyiID = &xscmsg.FieldID{
		Tag:        "6c.",
		Annotation: "fyi",
		Label:      "6. For your information (no action required)",
		Comment:    "boolean",
	}
	replyByID = &xscmsg.FieldID{
		Tag:        "6d.",
		Annotation: "reply-by",
		Label:      "6. Reply by",
	}
	toICSPositionID = &xscmsg.FieldID{
		Tag:        "7.",
		Annotation: "to-ics-position",
		Label:      "7. To ICS Position",
		Comment:    "required",
		Canonical:  xscmsg.FToICSPosition,
	}
	fromICSPositionID = &xscmsg.FieldID{
		Tag:        "8.",
		Annotation: "from-ics-position",
		Label:      "8. From ICS Position",
		Comment:    "required",
	}
	toLocationID = &xscmsg.FieldID{
		Tag:        "9a.",
		Annotation: "to-location",
		Label:      "9. To Location",
		Comment:    "required",
		Canonical:  xscmsg.FToLocation,
	}
	fromLocationID = &xscmsg.FieldID{
		Tag:        "9b.",
		Annotation: "from-location",
		Label:      "9. From Location",
		Comment:    "required",
	}
	toNameID = &xscmsg.FieldID{
		Tag:   "ToName",
		Label: "To Name",
	}
	fromNameID = &xscmsg.FieldID{
		Tag:   "FmName",
		Label: "From Name",
	}
	toTelID = &xscmsg.FieldID{
		Tag:   "ToTel",
		Label: "To Telephone #",
	}
	fromTelID = &xscmsg.FieldID{
		Tag:   "FmTel",
		Label: "From Telephone #",
	}
	subjectID = &xscmsg.FieldID{
		Tag:        "10.",
		Annotation: "subject",
		Label:      "10. Subject",
		Comment:    "required",
		Canonical:  xscmsg.FSubject,
	}
	referenceID = &xscmsg.FieldID{
		Tag:        "11.",
		Annotation: "reference",
		Label:      "11. Reference",
	}
	messageID = &xscmsg.FieldID{
		Tag:        "12.",
		Annotation: "message",
		Label:      "12. Message",
		Comment:    "required",
	}
	recSentID = &xscmsg.FieldID{
		Tag:     "Rec-Sent",
		Label:   "Receiver or Sender",
		Comment: "receiver, sender",
	}
	recSentChoices = []string{"receiver", "sender"}
	methodID       = &xscmsg.FieldID{
		Tag:     "Method",
		Label:   "How Received or Sent",
		Comment: "Telephone, Dispatch Center, EOC Radio, FAX, Courier, Amateur Radio, Other",
	}
	methodChoices = []string{"Telephone", "Dispatch Center", "EOC Radio", "FAX", "Courier", "Amateur Radio", "Other"}
	otherID       = &xscmsg.FieldID{
		Tag:   "Other",
		Label: "Other",
	}
)

type recSentField struct{ xscform.ChoicesField }

func (f *recSentField) Default() string { return "sender" }

type methodField struct{ xscform.ChoicesField }

func (f *methodField) Default() string { return "Other" }

type otherField struct{ xscform.Field }

func (f *otherField) Default() string { return "Packet" }
