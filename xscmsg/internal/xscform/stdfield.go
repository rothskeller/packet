package xscform

import (
	"fmt"

	"github.com/rothskeller/packet/xscmsg"
)

// FOriginMessageNumber creates an XSC-standard Origin Message Number field.
func FOriginMessageNumber() xscmsg.Field {
	return &MessageNumberField{Field: newField(originMessageNumberID, true)}
}

var originMessageNumberID = &xscmsg.FieldID{
	Tag:       "MsgNo",
	Label:     "Origin Message Number",
	Comment:   "required message-number",
	Canonical: xscmsg.FOriginMsgNo,
}

// FDestinationMessageNumber creates an XSC-standard Destination Message Number
// field.
func FDestinationMessageNumber() xscmsg.Field {
	return &MessageNumberField{Field: newField(destinationMessageNumberID, false)}
}

var destinationMessageNumberID = &xscmsg.FieldID{
	Tag:       "DestMsgNo",
	Label:     "Destination Message Number",
	Comment:   "message-number",
	Canonical: xscmsg.FDestinationMsgNo,
	ReadOnly:  true,
}

// FMessageDate creates an XSC-standard Message Date field.
func FMessageDate() xscmsg.Field {
	return &DateFieldDefaultNow{DateField{Field: newField(messageDateID, true)}}
}

var messageDateID = &xscmsg.FieldID{
	Tag:        "1a.",
	Annotation: "date",
	Label:      "Date",
	Comment:    "required date",
	Canonical:  xscmsg.FMessageDate,
}

// FMessageTime creates an XSC-standard Message Time field.
func FMessageTime() xscmsg.Field {
	return &TimeFieldDefaultNow{TimeField{Field: newField(messageTimeID, true)}}
}

var messageTimeID = &xscmsg.FieldID{
	Tag:        "1b.",
	Annotation: "time",
	Label:      "Time",
	Comment:    "required time",
	Canonical:  xscmsg.FMessageTime,
}

// FHandling creates an XSC-standard Handling field.
func FHandling() xscmsg.Field {
	return &ChoicesField{Field: newField(handlingID, true), Choices: handlingChoices}
}

var handlingID = &xscmsg.FieldID{
	Tag:        "5.",
	Annotation: "handling",
	Label:      "Handling",
	Comment:    "required: IMMEDIATE, PRIORITY, ROUTINE",
	Canonical:  xscmsg.FHandling,
}

var handlingChoices = []string{"IMMEDIATE", "PRIORITY", "ROUTINE"}

// FToICSPosition creates an XSC-standard To ICS Position field.
func FToICSPosition() xscmsg.Field {
	return NewField(toICSPositionID, true)
}

var toICSPositionID = &xscmsg.FieldID{
	Tag:        "7a.",
	Annotation: "to-ics-position",
	Label:      "To ICS Position",
	Comment:    "required",
	Canonical:  xscmsg.FToICSPosition,
}

// FFromICSPosition creates an XSC-standard From ICS Position field.
func FFromICSPosition() xscmsg.Field {
	return NewField(fromICSPositionID, true)
}

var fromICSPositionID = &xscmsg.FieldID{
	Tag:        "8a.",
	Annotation: "from-ics-position",
	Label:      "From ICS Position",
	Comment:    "required",
}

// FToLocation creates an XSC-standard To Location field.
func FToLocation() xscmsg.Field {
	return NewField(toLocationID, true)
}

var toLocationID = &xscmsg.FieldID{
	Tag:        "7b.",
	Annotation: "to-location",
	Label:      "To Location",
	Comment:    "required",
	Canonical:  xscmsg.FToLocation,
}

// FFromLocation creates an XSC-standard From Location field.
func FFromLocation() xscmsg.Field {
	return NewField(fromLocationID, true)
}

var fromLocationID = &xscmsg.FieldID{
	Tag:        "8b.",
	Annotation: "from-location",
	Label:      "From Location",
	Comment:    "required",
}

// FToName creates an XSC-standard To Name field.
func FToName() xscmsg.Field {
	return NewField(toNameID, false)
}

var toNameID = &xscmsg.FieldID{
	Tag:        "7c.",
	Annotation: "to-name",
	Label:      "To Name",
}

// FFromName creates an XSC-standard From Name field.
func FFromName() xscmsg.Field {
	return NewField(fromNameID, false)
}

var fromNameID = &xscmsg.FieldID{
	Tag:        "8c.",
	Annotation: "from-name",
	Label:      "From Name",
}

// FToContact creates an XSC-standard To Contact field.
func FToContact() xscmsg.Field {
	return NewField(toContactID, false)
}

var toContactID = &xscmsg.FieldID{
	Tag:        "7d.",
	Annotation: "to-contact",
	Label:      "To Contact Info",
}

// FFromContact creates an XSC-standard From Contact field.
func FFromContact() xscmsg.Field {
	return NewField(fromContactID, false)
}

var fromContactID = &xscmsg.FieldID{
	Tag:        "8d.",
	Annotation: "from-contact",
	Label:      "From Contact Info",
}

// FOpRelayRcvd creates an XSC-standard Operator Relay Received field.
func FOpRelayRcvd() xscmsg.Field {
	return NewField(opRelayRcvdID, false)
}

var opRelayRcvdID = &xscmsg.FieldID{
	Tag:   "OpRelayRcvd",
	Label: "Relay Rcvd",
}

// FOpRelaySent creates an XSC-standard Operator Relay Sent field.
func FOpRelaySent() xscmsg.Field {
	return NewField(opRelaySentID, false)
}

var opRelaySentID = &xscmsg.FieldID{
	Tag:   "OpRelaySent",
	Label: "Relay Sent",
}

// FOpName creates an XSC-standard Operator Name field.
func FOpName() xscmsg.Field {
	return NewField(opNameID, true)
}

var opNameID = &xscmsg.FieldID{
	Tag:       "OpName",
	Label:     "Operator Name",
	Canonical: xscmsg.FOpName,
}

// FOpCall creates an XSC-standard Operator Call field.
func FOpCall() xscmsg.Field {
	return &CallSignField{Field: newField(opCallID, true)}
}

var opCallID = &xscmsg.FieldID{
	Tag:       "OpCall",
	Label:     "Operator Call Sign",
	Comment:   "required call-sign",
	Canonical: xscmsg.FOpCall,
}

// FOpDate creates an XSC-standard Operator Date field.
func FOpDate() xscmsg.Field {
	return &DateFieldDefaultNow{DateField{Field: newField(opDateID, false)}}
}

var opDateID = &xscmsg.FieldID{
	Tag:       "OpDate",
	Label:     "Operator Date",
	Comment:   "date",
	Canonical: xscmsg.FOpDate,
}

// FOpTime creates an XSC-standard Operator Time field.
func FOpTime() xscmsg.Field {
	return &TimeFieldDefaultNow{TimeField{Field: newField(opTimeID, false)}}
}

var opTimeID = &xscmsg.FieldID{
	Tag:       "OpTime",
	Label:     "Operator Time",
	Comment:   "time",
	Canonical: xscmsg.FOpTime,
}

// unknownField is a field that was found while decoding a received form, that
// we did not expect to find in the form and have no field definition for.
type unknownField struct {
	tag   string
	value string
}

// ID returns the identification of the field.
func (f *unknownField) ID() *xscmsg.FieldID {
	return &xscmsg.FieldID{
		Tag:   f.tag,
		Label: f.tag,
	}
}

// Validate always fails for unknown fields.
func (f *unknownField) Validate(_ xscmsg.Message, _ bool) string {
	return fmt.Sprintf("form has a value for an unknown field %q", f.tag)
}

// Get returns the value of the field.
func (f *unknownField) Get() string { return f.value }

// Set sets the value of the field.
func (f *unknownField) Set(v string) { f.value = v }

// Default returns the default value of the field.
func (f *unknownField) Default() string { return "" }
