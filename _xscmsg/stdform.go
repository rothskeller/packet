package xscmsg

import (
	"time"

	"github.com/rothskeller/packet/pktmsg"
)

// NewStdForm creates a new StdForm implementation based on the supplied
// message.  It is not usable without being embedded in a specific form type.
// If the supplied message is nil, a new blank message will be created.
func NewStdForm(msg *Message) *StdForm {
	if msg == nil {
		msg = NewBaseMessage(nil)
	}
	return &StdForm{Message: msg}
}

// StdForm is the base class for all XSC-standard forms.  It contains all of the
// fields they have in common.
type StdForm struct {
	*Message
	DestinationMsgID string
	MessageDate      string
	MessageTime      string
	ToICSPosition    string
	ToLocation       string
	ToName           string
	ToContact        string
	FromICSPosition  string
	FromLocation     string
	FromName         string
	FromContact      string
	OpRelayRcvd      string
	OpRelaySent      string
	OpName           string
	OpCall           string
	OpDate           string
	OpTime           string
}

// GetOpCall returns the value of the operator call sign field.
func (m *StdForm) GetOpCall() string { return m.OpCall }

// GetToICSPosition returns the value of the "To ICS Position" field.
func (m *StdForm) GetToICSPosition() string { return m.ToICSPosition }

// GetToLocation returns the value of the "To Location" field.
func (m *StdForm) GetToLocation() string { return m.ToLocation }

// SetDestinationMsgID sets the destination message ID for the message.
func (m *StdForm) SetDestinationMsgID(value string) { m.DestinationMsgID = value }

// SetOperator sets the OpCall, OpName, OpDate, and OpTime fields of the
// message.
func (m *StdForm) SetOperator(received bool, call string, name string) {
	m.OpCall, m.OpName = call, name
	m.OpDate = time.Now().Format("01/02/2006")
	m.OpTime = time.Now().Format("15:04")
}

// ReadTaggedFields reads the tagged fields of the enclosed message and assigns
// the values of the standard form fields from them.  It removes the tagged
// fields that correspond to standard form fields from the argument slice, and
// returns the resulting slice.
func (m *StdForm) ReadTaggedFields(fields []pktmsg.TaggedField) []pktmsg.TaggedField {
	j := 0
	for _, tf := range fields {
		switch tf.Tag {
		case "MsgNo":
			m.OriginMsgID = tf.Value
		case "DestMsgNo":
			m.DestinationMsgID = tf.Value
		case "1a.":
			m.MessageDate = tf.Value
		case "1b.":
			m.MessageTime = tf.Value
		case "5.":
			m.Handling = tf.Value
		case "7a.":
			m.ToICSPosition = tf.Value
		case "7b.":
			m.ToLocation = tf.Value
		case "7c.":
			m.ToName = tf.Value
		case "7d.":
			m.ToContact = tf.Value
		case "8a.":
			m.FromICSPosition = tf.Value
		case "8b.":
			m.FromLocation = tf.Value
		case "8c.":
			m.FromName = tf.Value
		case "8d.":
			m.FromContact = tf.Value
		case "OpRelayRcvd":
			m.OpRelayRcvd = tf.Value
		case "OpRelaySent":
			m.OpRelaySent = tf.Value
		case "OpName":
			m.OpName = tf.Value
		case "OpCall":
			m.OpCall = tf.Value
		case "OpDate":
			m.OpDate = tf.Value
		case "OpTime":
			m.OpTime = tf.Value
		default:
			fields[j], j = tf, j+1
		}
	}
	return fields[:j]
}

// MakeLeadingTaggedFields returns the set of tagged fields for the standard
// form fields that precede the form-specific fields.
func (m *StdForm) MakeLeadingTaggedFields() []pktmsg.TaggedField {
	return []pktmsg.TaggedField{
		{Tag: "MsgNo", Value: m.OriginMsgID},
		{Tag: "DestMsgNo", Value: m.DestinationMsgID},
		{Tag: "1a.", Value: m.MessageDate},
		{Tag: "1b.", Value: m.MessageTime},
		{Tag: "5.", Value: m.Handling},
		{Tag: "7a.", Value: m.ToICSPosition},
		{Tag: "8a.", Value: m.FromICSPosition},
		{Tag: "7b.", Value: m.ToLocation},
		{Tag: "8b.", Value: m.FromLocation},
		{Tag: "7c.", Value: m.ToName},
		{Tag: "8c.", Value: m.FromName},
		{Tag: "7d.", Value: m.ToContact},
		{Tag: "8d.", Value: m.FromContact},
	}
}

// MakeTrailingTaggedFields returns the set of tagged fields for the standard
// form fields that follow the form-specific fields.
func (m *StdForm) MakeTrailingTaggedFields() []pktmsg.TaggedField {
	return []pktmsg.TaggedField{
		{Tag: "OpRelayRcvd", Value: m.OpRelayRcvd},
		{Tag: "OpRelaySent", Value: m.OpRelaySent},
		{Tag: "OpName", Value: m.OpName},
		{Tag: "OpCall", Value: m.OpCall},
		{Tag: "OpDate", Value: m.OpDate},
		{Tag: "OpTime", Value: m.OpTime},
	}
}

// Validate checks the validity of the message and returns strings describing
// the issues.  This is used for received messages, and only checks for validity
// issues that Outpost and/or PackItForms check for.
func (m *StdForm) Validate() (problems []string) {
	if m.MessageDate == "" {
		problems = append(problems, "The message date is required.")
	} else if !ValidDate(m.MessageDate) {
		problems = append(problems, "The message date is not a valid date.  Dates must be in MM/DD/YYYY format.")
	}
	if m.MessageTime == "" {
		problems = append(problems, "The message time is required.")
	} else if !ValidTime(m.MessageTime) {
		problems = append(problems, "The message time is not a valid time.  Times must be in HH:MM format (24-hour clock).")
	}
	if m.Handling == "" {
		problems = append(problems, "The message handling order is required.")
	} else if !ValidRestrictedValue(m.Handling, handlingChoices) {
		problems = append(problems, "The message handling order must be ROUTINE, PRIORITY, or IMMEDIATE.")
	}
	if m.ToICSPosition == "" {
		problems = append(problems, "The To ICS Position is required.")
	}
	if m.ToLocation == "" {
		problems = append(problems, "The To Location is required.")
	}
	if m.FromICSPosition == "" {
		problems = append(problems, "The From ICS Position is required.")
	}
	if m.FromLocation == "" {
		problems = append(problems, "The From Location is required.")
	}
	if m.OpCall == "" {
		problems = append(problems, "The operator call sign is required.")
	} else if !ValidFCCCallSign(m.OpCall) {
		problems = append(problems, "The operator call sign is not a valid FCC call sign.")
	}
	if m.OpDate != "" && !ValidDate(m.OpDate) {
		problems = append(problems, "The operator date is not a valid date.  Dates must be in MM/DD/YYYY format.")
	}
	if m.OpTime != "" && !ValidTime(m.OpTime) {
		problems = append(problems, "The operator time is not a valid time.  Times must be in HH:MM format (24-hour clock).")
	}
	return problems
}

// MakeLeadingViewFields returns the set of viewable fields for the standard
// form fields that precede the form-specific fields.
func (m *StdForm) MakeLeadingViewFields() (lvs []LabelValue) {
	lvs = m.ViewHeaders()
	lvs = append(lvs, []LabelValue{
		{"Origin Message Number", m.OriginMsgID},
		{"Destination Message Number", m.DestinationMsgID},
		{"Message Date/Time", smartJoin(m.MessageDate, m.MessageTime, " ")},
		{"Handling Order", m.Handling},
		{"To ICS Position", m.ToICSPosition},
		{"To Location", m.ToLocation},
		{"To Name", m.ToName},
		{"To Contact", m.ToContact},
		{"From ICS Position", m.FromICSPosition},
		{"From Location", m.FromLocation},
		{"From Name", m.FromName},
		{"From Contact", m.FromContact},
	}...)
	return lvs
}

// MakeTrailingViewFields returns the set of viewable fields for the standard
// form fields that follow the form-specific fields.
func (m *StdForm) MakeTrailingViewFields() (lvs []LabelValue) {
	lvs = []LabelValue{
		{"Operator Relay Received", m.OpRelayRcvd},
		{"Operator Relay Sent", m.OpRelaySent},
	}
	if m.RxBBS != "" {
		lvs = append(lvs, []LabelValue{
			{"Received by", smartJoin(m.OpName, m.OpCall, " ")},
			{"Received at", smartJoin(m.OpDate, m.OpTime, " ")},
		}...)
	} else {
		lvs = append(lvs, []LabelValue{
			{"Sent by", smartJoin(m.OpName, m.OpCall, " ")},
			{"Sent at", smartJoin(m.OpDate, m.OpTime, " ")},
		}...)
	}
	return lvs
}

func smartJoin(a, b, sep string) string {
	if a == "" || b == "" {
		return a + b
	}
	return a + sep + b
}

// MakeLeadingEditFields returns the set of edit Fields for the standard form
// fields that precede the form-specific fields.
func (m *StdForm) MakeLeadingEditFields() []Field {
	return []Field{
		NewToField(&m.To),
		NewOriginMsgIDField(&m.OriginMsgID),
		WrapRequiredField(WrapDateField(&messageDateField{&baseField{&m.MessageDate}})),
		WrapRequiredField(WrapTimeField(&messageTimeField{&baseField{&m.MessageTime}})),
		NewHandlingField(&m.Handling),
		WrapRequiredField(&toICSPositionField{&baseField{&m.ToICSPosition}}),
		WrapRequiredField(&toLocationField{&baseField{&m.ToLocation}}),
		&toNameField{&baseField{&m.ToName}},
		&toContactField{&baseField{&m.ToContact}},
		WrapRequiredField(&fromICSPositionField{&baseField{&m.FromICSPosition}}),
		WrapRequiredField(&fromLocationField{&baseField{&m.FromLocation}}),
		&fromNameField{&baseField{&m.FromName}},
		&fromContactField{&baseField{&m.FromContact}},
	}
}

// MakeTrailingEditFields returns the set of edit Fields for the standard form
// fields that follow the form-specific fields.
func (m *StdForm) MakeTrailingEditFields() []Field {
	return []Field{
		&opRelayRcvdField{&baseField{&m.OpRelayRcvd}},
		&opRelaySentField{&baseField{&m.OpRelaySent}},
	}
}

type messageDateField struct{ Field }

func (f messageDateField) Label() string { return "Message Date" }
func (f messageDateField) Help() string {
	return "This is the date the message was written."
}

type messageTimeField struct{ Field }

func (f messageTimeField) Label() string { return "Message Time" }
func (f messageTimeField) Help() string {
	return "This is the time the message was written."
}

type toICSPositionField struct{ Field }

func (f toICSPositionField) Label() string { return "To ICS Position" }
func (f toICSPositionField) Help() string {
	return "This is the ICS position to which the message is addressed."
}

type toLocationField struct{ Field }

func (f toLocationField) Label() string { return "To Location" }
func (f toLocationField) Help() string {
	return "This is the location of the recipient ICS position."
}

type toNameField struct{ Field }

func (f toNameField) Label() string { return "To Name" }
func (f toNameField) Help() string {
	return "This is the name of the person holding the recipient ICS position.  It is optional and rarely provided."
}

type toContactField struct{ Field }

func (f toContactField) Label() string { return "To Contact Info" }
func (f toContactField) Help() string {
	return "This is contact information (phone number, email, etc.) for the receipient.  It is optional and rarely provided."
}

type fromICSPositionField struct{ Field }

func (f fromICSPositionField) Label() string { return "From ICS Position" }
func (f fromICSPositionField) Help() string {
	return "This is the ICS position of the message author."
}

type fromLocationField struct{ Field }

func (f fromLocationField) Label() string { return "From Location" }
func (f fromLocationField) Help() string {
	return "This is the location of the message author."
}

type fromNameField struct{ Field }

func (f fromNameField) Label() string { return "From Name" }
func (f fromNameField) Help() string {
	return "This is the name of the message author.  It is optional and rarely provided."
}

type fromContactField struct{ Field }

func (f fromContactField) Label() string { return "From Contact Info" }
func (f fromContactField) Help() string {
	return "This is contact information (phone number, email, etc.) for the message author.  It is optional and rarely provided."
}

type opRelayRcvdField struct{ Field }

func (f opRelayRcvdField) Label() string { return "(Op) Relay Received" }
func (f opRelayRcvdField) Help() string {
	return "This is the name of the station from which this message was directly received.  It is filled in for messages that go through a relay station."
}

type opRelaySentField struct{ Field }

func (f opRelaySentField) Label() string { return "(Op) Relay Sent" }
func (f opRelaySentField) Help() string {
	return "This is the name of the station to which this message was directly sent.  It is filled in for messages that go through a relay station."
}
