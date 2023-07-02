// Package checkin defines the check-in message type.
package checkin

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/typedmsg"
	"github.com/rothskeller/packet/xscmsg"
)

// Type is the type definition for a check-in message.
var Type = typedmsg.MessageType{
	Tag:       "Check-In",
	Name:      "check-in message",
	Article:   "a",
	Create:    create,
	Recognize: recognize,
}

// CheckIn is a check-in message.
type CheckIn struct {
	*xscmsg.Message
	TacCall string
	TacName string
	OpCall  string
	OpName  string
	fields  []xscmsg.Field
}

// NewCheckIn creates a new check-in message.
func NewCheckIn() *CheckIn {
	ci := CheckIn{Message: &xscmsg.Message{Message: new(pktmsg.Message)}}
	ci.Handling = "ROUTINE"
	return &ci
}

func create() typedmsg.Message { return NewCheckIn() }

var checkInRE = regexp.MustCompile(`(?i)^Check-In\s+([A-Z][A-Z0-9]{2,5})\s*,(.*)(?:\n([AKNW][A-Z0-9]{2,5})\s*,(.*))?`)

func recognize(base *pktmsg.Message) typedmsg.Message {
	if base.FormTag != "" || !strings.HasPrefix(strings.ToLower(base.Subject), "check-in ") {
		return nil
	}
	ci := CheckIn{Message: &xscmsg.Message{Message: base}}
	if match := checkInRE.FindStringSubmatch(base.Body); match != nil {
		if match[3] != "" {
			ci.TacCall = match[1]
			ci.TacName = strings.TrimSpace(match[2])
			ci.OpCall = match[3]
			ci.OpName = strings.TrimSpace(match[4])
		} else {
			ci.OpCall = match[1]
			ci.OpName = strings.TrimSpace(match[2])
		}
	}
	return &ci
}

// Type returns the type of the message.
func (m *CheckIn) Type() *typedmsg.MessageType { return &Type }

// View returns the set of viewable fields of the message.
func (m *CheckIn) View() []xscmsg.LabelValue {
	var lvs = m.ViewHeaders()
	lvs = append(lvs,
		xscmsg.LV("Subject", m.Subject),
		xscmsg.LV("Tactical Station", smartJoin(m.TacName, m.TacCall, " ")),
		xscmsg.LV("Operator", smartJoin(m.OpName, m.OpCall, " ")),
	)
	return lvs
}
func smartJoin(a, b, sep string) string {
	if a == "" || b == "" {
		return a + b
	}
	return a + sep + b
}

// Edit returns the set of editable fields of the message.
func (m *CheckIn) Edit() []xscmsg.Field {
	if m.fields == nil {
		m.fields = []xscmsg.Field{
			xscmsg.NewToField(&m.To),
			xscmsg.NewOriginMsgIDField(&m.OriginMsgID),
			xscmsg.NewHandlingField(&m.Handling),
			xscmsg.WrapTacticalCallSignField(&tacCallField{xscmsg.BaseField(&m.TacCall)}),
			&tacNameField{xscmsg.BaseField(&m.TacName), &m.TacCall},
			xscmsg.WrapRequiredField(xscmsg.WrapFCCCallSignField(&opCallField{xscmsg.BaseField(&m.OpCall)})),
			xscmsg.WrapRequiredField(&opNameField{xscmsg.BaseField(&m.OpName)}),
		}
	}
	return m.fields
}

// GetBody is a no-op for check-in messages.
func (m *CheckIn) GetBody() string { return "" }

// SetBody is a no-op for check-in messages.
func (m *CheckIn) SetBody(string) {}

// SetSubject is a no-op for check-in messages.
func (m *CheckIn) SetSubject(string) {}

// GetOpCall returns the value of the operator call sign field.
func (m *CheckIn) GetOpCall() string { return m.OpCall }

// SetOperator sets the OpCall and OpName fields of the message (outgoing only).
func (m *CheckIn) SetOperator(received bool, call string, name string) {
	if !received {
		m.OpCall, m.OpName = call, name
	}
}

// SetTactical sets the TacCall and TacName fields of the message.
func (m *CheckIn) SetTactical(call string, name string) {
	m.TacCall, m.TacName = call, name
}

// Save renders the message for saving to local storage.
func (m *CheckIn) Save() string {
	m.encode()
	return m.Message.Save()
}

// Transmit renders the message for transmission.
func (m *CheckIn) Transmit() ([]string, string, string) {
	m.encode()
	return m.Message.Transmit()
}

func (m *CheckIn) encode() {
	if m.TacCall != "" {
		m.Subject = fmt.Sprintf("Check-In %s, %s", m.TacCall, m.TacName)
		m.Body = fmt.Sprintf("Check-In %s, %s\n%s, %s\n", m.TacCall, m.TacName, m.OpCall, m.OpName)
	} else {
		m.Subject = fmt.Sprintf("Check-In %s, %s", m.OpCall, m.OpName)
		m.Body = fmt.Sprintf("Check-In %s, %s\n", m.OpCall, m.OpName)
	}
}

// Validate checks the validity of the message and returns strings
// describing the issues.  This is used for received messages, and only
// checks for validity issues that Outpost and/or PackItForms check for.
func (m *CheckIn) Validate() []string {
	return xscmsg.ValidateFromFieldProblems(m)
}

type tacCallField struct{ xscmsg.Field }

func (f tacCallField) Label() string { return "Tactical Call Sign" }
func (f tacCallField) Help() string {
	return "This is the tactical call sign of the station."
}

type tacNameField struct {
	xscmsg.Field
	tacCall *string
}

func (f tacNameField) Label() string    { return "Tactical Station Name" }
func (f tacNameField) Size() (w, h int) { return 80, 1 }
func (f tacNameField) Problem() string {
	if f.Value() == "" && *f.tacCall != "" {
		return "The tactical station name is required when a tactical call sign is provided."
	}
	if f.Value() != "" && *f.tacCall == "" {
		return "A tactical station name is not allowed unless a tactical call sign is also provided."
	}
	return ""
}
func (f tacNameField) Help() string {
	return "This is the name of the station.  If a tactical call sign is supplied, a tactical station name is required.  Otherwise, a tactical station name is not allowed."
}

type opCallField struct{ xscmsg.Field }

func (f opCallField) Label() string { return "Operator Call Sign" }
func (f opCallField) Help() string {
	return "This is the FCC call sign of the station operator."
}

type opNameField struct{ xscmsg.Field }

func (f opNameField) Label() string    { return "Operator Name" }
func (f opNameField) Size() (w, h int) { return 80, 1 }
func (f opNameField) Help() string {
	return "This is the name of the station operator."
}
