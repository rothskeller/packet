package ics213

import (
	"time"

	"github.com/rothskeller/packet/xscmsg"
)

func createICS213(old, tx bool, c xscmsg.FieldContainer) (pifo, human []xscmsg.FormField) {
	dateField := newDateField(c)
	timeField := newTimeField(c)
	takeActionField := newTakeActionField(c)
	replyField := newReplyField(c)
	replyByField := newReplyByField(c)
	toICSPositionField := newToICSPositionField(c)
	fromICSPositionField := newFromICSPositionField(c)
	toLocationField := newToLocationField(c)
	fromLocationField := newFromLocationField(c)
	toNameField := newToNameField(c)
	fromNameField := newFromNameField(c)
	toTelField := newToTelField(c)
	fromTelField := newFromTelField(c)
	subjectField := newSubjectField(c)
	referenceField := newReferenceField(c)
	messageField := newMessageField(c)
	opRelayRcvdField := newOpRelayRcvdField(c)
	opRelaySentField := newOpRelaySentField(c)
	recSentField := newRecSentField(c)
	opCallField := newOpCallField(c)
	opNameField := newOpNameField(c)
	methodField := newMethodField(c)
	otherField := newOtherField(c)
	opDateField := newOpDateField(c)
	opTimeField := newOpTimeField(c)
	if old { // <= v2.1
		var senderMessageNumberField, myMessageNumberField, receiverMessageNumberField xscmsg.FormField
		if tx {
			senderMessageNumberField = newSenderMessageNumberTxField(c)
			myMessageNumberField = newMyMessageNumberTxField(c)
			receiverMessageNumberField = newReceiverMessageNumberTxField(c)
		} else {
			senderMessageNumberField = newSenderMessageNumberRxField(c)
			myMessageNumberField = newMyMessageNumberRxField(c)
			receiverMessageNumberField = newReceiverMessageNumberRxField(c)
		}
		severityField := newSeverityField(c)
		handlingOrderField := newHandlingOrderField(c)
		fyiField := newFYIField(c)
		pifo = []xscmsg.FormField{
			senderMessageNumberField, myMessageNumberField, receiverMessageNumberField, dateField, timeField,
			severityField, handlingOrderField, takeActionField, replyField, replyByField, fyiField, toICSPositionField,
			fromICSPositionField, toLocationField, fromLocationField, toNameField, fromNameField, toTelField,
			fromTelField, subjectField, referenceField, messageField, opRelayRcvdField, opRelaySentField, recSentField,
			opCallField, opNameField, methodField, otherField, opDateField, opTimeField,
		}
		human = []xscmsg.FormField{
			senderMessageNumberField, myMessageNumberField, receiverMessageNumberField, dateField, timeField,
			severityField, handlingOrderField, takeActionField, replyField, replyByField, fyiField, toICSPositionField,
			toLocationField, toNameField, toTelField, fromICSPositionField, fromLocationField, fromNameField,
			fromTelField, subjectField, referenceField, messageField, opRelayRcvdField, opRelaySentField, recSentField,
			methodField, otherField, opCallField, opNameField, opDateField, opTimeField,
		}
	} else { // >= v2.2
		// Message numbers simplified, severity removed, handling label
		// changed, FYI removed.
		originMessageNumberField := newOriginMessageNumberField(c)
		destinationMessageNumberField := newDestinationMessageNumberField(c)
		handlingField := newHandlingField(c)
		pifo = []xscmsg.FormField{
			originMessageNumberField, destinationMessageNumberField, dateField, timeField,
			handlingField, takeActionField, replyField, replyByField, toICSPositionField,
			fromICSPositionField, toLocationField, fromLocationField, toNameField, fromNameField, toTelField,
			fromTelField, subjectField, referenceField, messageField, opRelayRcvdField, opRelaySentField, recSentField,
			opCallField, opNameField, methodField, otherField, opDateField, opTimeField,
		}
		human = []xscmsg.FormField{
			originMessageNumberField, destinationMessageNumberField, dateField, timeField,
			handlingField, takeActionField, replyField, replyByField, toICSPositionField,
			toLocationField, toNameField, toTelField, fromICSPositionField, fromLocationField, fromNameField,
			fromTelField, subjectField, referenceField, messageField, opRelayRcvdField, opRelaySentField, recSentField,
			methodField, otherField, opCallField, opNameField, opDateField, opTimeField,
		}
	}
	return pifo, human
}

type v = xscmsg.Value // shortcut used in the definitions below

func newSenderMessageNumberRxField(c xscmsg.FieldContainer) (f xscmsg.FormField) { // <= v2.1, received
	f = xscmsg.BaseField(c, xscmsg.FOriginMsgNo, "2.", "2. Sender's Msg #", 10, 1,
		`This is the origin message number that the sender of this message assigned to it.`)
	f = xscmsg.MessageNumberField(true, f)
	f = xscmsg.RequiredField(f)
	f = xscmsg.ReadOnlyField(f)
	return f
}

func newMyMessageNumberRxField(c xscmsg.FieldContainer) (f xscmsg.FormField) { // <= v2.1, received
	f = xscmsg.BaseField(c, xscmsg.FDestinationMsgNo, "MsgNo", "My Msg #", 10, 1,
		`This is the destination message number that we assigned to this message when we received it.`)
	f = xscmsg.MessageNumberField(true, f)
	f = xscmsg.RequiredField(f)
	return f
}

func newReceiverMessageNumberRxField(c xscmsg.FieldContainer) (f xscmsg.FormField) { // <= v2.1, received
	f = xscmsg.BaseField(c, "", "3.", "3. Receiver's Msg #", 10, 1,
		`This field is not used on received messages and must be empty.`)
	f = xscmsg.ReadOnlyField(f)
	return f
}

func newSenderMessageNumberTxField(c xscmsg.FieldContainer) (f xscmsg.FormField) { // <= v2.1, sent
	f = xscmsg.BaseField(c, xscmsg.FOriginMsgNo, "2.", "2. Sender's Msg #", 10, 1,
		`This field is not used on sent messages and must be empty.`)
	f = xscmsg.ReadOnlyField(f)
	return f
}

func newMyMessageNumberTxField(c xscmsg.FieldContainer) (f xscmsg.FormField) { // <= v2.1, sent
	f = xscmsg.BaseField(c, xscmsg.FOriginMsgNo, "MsgNo", "My Msg #", 10, 1,
		`This is the origin message number that we assigned to this outgoing message when we created it.`)
	f = xscmsg.MessageNumberField(true, f)
	f = xscmsg.RequiredField(f)
	return f
}

func newReceiverMessageNumberTxField(c xscmsg.FieldContainer) (f xscmsg.FormField) { // <= v2.1, sent
	f = xscmsg.BaseField(c, xscmsg.FDestinationMsgNo, "3.", "3. Receiver's Msg #", 10, 1,
		`This is the destination message number that the receiver of this message assigned to it.`)
	f = xscmsg.MessageNumberField(true, f)
	f = xscmsg.RequiredField(f)
	f = xscmsg.ReadOnlyField(f)
	return f
}

func newOriginMessageNumberField(c xscmsg.FieldContainer) (f xscmsg.FormField) { // >= v2.2
	f = xscmsg.BaseField(c, xscmsg.FOriginMsgNo, "MsgNo", "2. Origin Msg #", 10, 1,
		`This is the message number for this message assigned to it by its origin station.`)
	f = xscmsg.MessageNumberField(true, f)
	f = xscmsg.RequiredField(f)
	return f
}

func newDestinationMessageNumberField(c xscmsg.FieldContainer) (f xscmsg.FormField) { // >= v2.2
	f = xscmsg.BaseField(c, xscmsg.FDestinationMsgNo, "3.", "3. Destination Msg #", 10, 1,
		`This is the message number for this message assigned to it by its destination station.`)
	f = xscmsg.MessageNumberField(true, f)
	f = xscmsg.ReadOnlyField(f)
	return f
}

func newDateField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, "", "1a.", "1. Date", 10, 1,
		`This is the date when the message was written (not necessarily when it was sent).`)
	f = xscmsg.DateField(f)
	f = dateField{f}
	f = xscmsg.RequiredField(f)
	return f
}

type dateField struct{ xscmsg.FormField }

func (df dateField) Default() xscmsg.Value { return df.FromHuman(time.Now().Format("01/02/2006")) }

func newTimeField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, "", "1b.", "1. Time (24hr)", 5, 1,
		`This is the time when the message was written (not when it was sent).`)
	f = xscmsg.TimeField(f)
	f = timeField{f}
	f = xscmsg.RequiredField(f)
	return f
}

type timeField struct{ xscmsg.FormField }

func (tf timeField) Default() xscmsg.Value { return tf.FromHuman(time.Now().Format("15:04")) }

func newSeverityField(c xscmsg.FieldContainer) (f xscmsg.FormField) { // <= v2.1 only
	f = xscmsg.BaseField(c, "", "4.", "4. Situation Severity", 1, 1,
		`This is the severity of the situation to which this message is related.`)
	f = xscmsg.ChoicesField(f, v("EMERGENCY"), v("URGENT"), v("OTHER"))
	f = xscmsg.RequiredField(f)
	return f
}

func newHandlingOrderField(c xscmsg.FieldContainer) (f xscmsg.FormField) { // <= v2.1 only
	f = xscmsg.BaseField(c, xscmsg.FHandling, "5.", "5. Message Handling Order", 1, 1,
		`This is the handling order with which this message should be processed.`)
	f = xscmsg.ChoicesField(f, v("IMMEDIATE"), v("PRIORITY"), v("ROUTINE"))
	f = xscmsg.RequiredField(f)
	return f
}

func newHandlingField(c xscmsg.FieldContainer) (f xscmsg.FormField) { // >= v2.2 only
	f = xscmsg.BaseField(c, xscmsg.FHandling, "5.", "5. Handling", 1, 1,
		`This is the handling order with which this message should be processed.`)
	f = xscmsg.ChoicesField(f, v("IMMEDIATE"), v("PRIORITY"), v("ROUTINE"))
	f = xscmsg.RequiredField(f)
	return f
}

func newTakeActionField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, "", "6a.", "6. Take Action", 1, 1,
		`This indicates whether the recipient of the message is expected to take action on it.`)
	f = xscmsg.ChoicesField(f, v("Yes"), v("No"))
	f = xscmsg.NoValidatePIFO(xscmsg.RequiredField, f)
	return f
}

func newReplyField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, "", "6b.", "6. Reply", 1, 1,
		`This indicates whether the recipient of the message is expected to reply to it.`)
	f = xscmsg.ChoicesField(f, v("Yes"), v("No"))
	f = xscmsg.NoValidatePIFO(xscmsg.RequiredField, f)
	return f
}

func newFYIField(c xscmsg.FieldContainer) (f xscmsg.FormField) { // <= v2.1
	f = xscmsg.BaseField(c, "", "6c.", "6. For your information", 1, 1,
		`This indicates that the message is informational only.`)
	f = xscmsg.ChoicesField(f, v("checked"), "Yes", v(""), "No")
	return f
}

func newReplyByField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, "", "6d.", "6. Reply by", 5, 1,
		"This is the time by which a reply is requested.")
	f = xscmsg.TimeField(f)
	f = replyByField{f}
	return f
}

type replyByField struct {
	xscmsg.FormField
}

func (f replyByField) Editable() bool {
	return f.Container().FieldValue("6b.") == "Yes"
}

func (f replyByField) Validate() string {
	if f.Value() != "" && f.Container().FieldValue("6b.") != "Yes" {
		return `A value cannot be specified for the "Reply by" field unless the "Reply" field is set to "Yes".`
	}
	return f.FormField.Validate()
}

// Note there is no ValidatePIFO for replyByField; the restriction is not
// enforced by PackItForms.

func newToICSPositionField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, xscmsg.FToICSPosition, "7.", "7. To ICS Position", 40, 1,
		`This is the ICS Position (section, branch, unit, etc.) to which the message is addressed.`)
	f = xscmsg.RequiredField(f)
	return f
}

func newFromICSPositionField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, "", "8.", "8. From ICS Position", 40, 1,
		`This is the ICS Position (section, branch, unit, etc.) from which the message originated.`)
	f = xscmsg.RequiredField(f)
	return f
}

func newToLocationField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, xscmsg.FToLocation, "9a.", "9. To Location", 40, 1,
		`This is the location (EOC, shelter, etc.) to which the message is addressed.`)
	f = xscmsg.RequiredField(f)
	return f
}

func newFromLocationField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, "", "9b.", "9b. From Location", 40, 1,
		`This is the location (EOC, shelter, etc.) from which the message originated.`)
	f = xscmsg.RequiredField(f)
	return f
}

func newToNameField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, "", "ToName", "To Name", 40, 1,
		`This is the name of the person to whom the message is addressed.  This is optional and rarely provided; the "To ICS Position" field is the critical one.`)
	return f
}

func newFromNameField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, "", "FmName", "From Name", 40, 1,
		`This is the name of the person who wrote the message.  This is optional and rarely provided; the "From ICS Position" field is the critical one.`)
	return f
}

func newToTelField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, "", "ToTel", "To Telephone #", 40, 1,
		`This is the telephone number of the person to whom the message is addressed.  This is optional and rarely provided.`)
	return f
}

func newFromTelField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, "", "FmTel", "From Telephone #", 40, 1,
		`This is the telephone number of the person who wrote the message.  This is optional and rarely provided.`)
	return f
}

func newSubjectField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, xscmsg.FSubject, "10.", "10. Subject", 80, 1,
		`This is the subject of the message.`)
	f = xscmsg.RequiredField(f)
	return f
}

func newReferenceField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, xscmsg.FReference, "11.", "11. Reference", 55, 1,
		`This is the origin message number of the previous message, if any, to which this is a response.`)
	f = xscmsg.MessageNumberField(false, f)
	return f
}

func newMessageField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, xscmsg.FBody, "12.", "12. Message", 80, 8,
		`This is the body of the message.`)
	f = xscmsg.RequiredField(f)
	return f
}

func newOpRelayRcvdField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, "", "OpRelayRcvd", "Relay Rcvd", 32, 1,
		`This is the name of the relay station from which this message was received, if it was not received directly from the origin station.`)
	return f
}

func newOpRelaySentField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, "", "OpRelaySent", "Relay Sent", 32, 1,
		`This is the name of the relay station to which this message was sent, if it was not sent directly to the destination station.`)
	return f
}

func newRecSentField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, "", "Rec-Sent", "Receiver or Sender", 1, 1,
		`This indicates whether the message was received by the local station or sent by it.`)
	f = xscmsg.ChoicesField(f, v("receiver"), v("sender"))
	f = xscmsg.RequiredField(f)
	f = xscmsg.ReadOnlyField(f)
	return f
}

func newOpCallField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, xscmsg.FOpCall, "OpCall", "Operator Call Sign", 28, 1,
		`This is the FCC call sign of the local radio operator who handled this message.`)
	f = xscmsg.RequiredField(xscmsg.CallSignField(true, f))
	f = xscmsg.ReadOnlyField(f)
	return f
}

func newOpNameField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, xscmsg.FOpName, "OpName", "Operator Name", 29, 1,
		`This is the name of the local radio operator who handled this message.`)
	f = xscmsg.RequiredField(f)
	f = xscmsg.ReadOnlyField(f)
	return f
}

func newMethodField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, "", "Method", "How Received or Sent", 1, 1,
		`This indicates the method by which the message was received or sent.`)
	f = xscmsg.ChoicesField(f, v("Telephone"), v("Dispatch Center"), v("EOC Radio"), v("FAX"), v("Courier"), v("Amateur Radio"), v("Other"))
	f = methodField{f}
	f = xscmsg.RequiredField(f)
	f = xscmsg.ReadOnlyField(f)
	return f
}

type methodField struct{ xscmsg.FormField }

func (methodField) Default() xscmsg.Value { return v("Other") }

func newOtherField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, "", "Other", "How Received or Sent: Other", 18, 1,
		`This is the "Other" method by which this message was received or sent, if the "Other" option is selected.`)
	f = otherField{f}
	return f
}

type otherField struct{ xscmsg.FormField }

func (of otherField) Default() xscmsg.Value { return of.FromHuman("Packet") }

func (of otherField) Editable() bool { return false }

func (of otherField) Validate() string {
	if of.Container().FieldValue("Method") == "Other" {
		if of.Value() == "" {
			return `A value for the "Other" field is required when the "Other" option is selected.`
		}
	} else if of.Value() != "" {
		return `A value cannot be specified for the "Other" field unless the "Other" option is selected.`
	}
	return ""
}

// Note there is no ValidatePIFO for otherField; the restrictions are not
// enforced by PackItForms.

func newOpDateField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, xscmsg.FOpDate, "OpDate", "Operator Date", 10, 1,
		`This is the date when the message was sent or received by the local station.`)
	f = xscmsg.DateField(f)
	f = xscmsg.RequiredField(f)
	f = xscmsg.ReadOnlyField(f)
	return f
}

func newOpTimeField(c xscmsg.FieldContainer) (f xscmsg.FormField) {
	f = xscmsg.BaseField(c, xscmsg.FOpTime, "OpTime", "Operator Time", 5, 1,
		`This is the time when the message was sent or received by the local station.`)
	f = xscmsg.TimeField(f)
	f = xscmsg.RequiredField(f)
	f = xscmsg.ReadOnlyField(f)
	return f
}
