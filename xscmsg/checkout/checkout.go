// Package checkout defines the check-out message type.
package checkout

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/typedmsg"
	"github.com/rothskeller/packet/xscmsg"
)

// Type is the type definition for a check-out message.
var Type = typedmsg.MessageType{
	Tag:       "Check-Out",
	Name:      "check-out message",
	Article:   "a",
	Create:    create,
	Recognize: recognize,
}

// CheckOut is a check-out message.
type CheckOut struct {
	*xscmsg.Message
	TacCall string
	TacName string
	OpCall  string
	OpName  string
	fields  []xscmsg.Field
}

// NewCheckOut creates a new check-out message.
func NewCheckOut() *CheckOut {
	co := CheckOut{Message: &xscmsg.Message{Message: new(pktmsg.Message)}}
	co.Handling = "ROUTINE"
	return &co
}

func create() typedmsg.Message { return NewCheckOut() }

var checkInRE = regexp.MustCompile(`(?i)^Check-Out\s+([A-Z][A-Z0-9]{2,5})\s*,(.*)(?:\n([AKNW][A-Z0-9]{2,5})\s*,(.*))?`)

func recognize(base *pktmsg.Message) typedmsg.Message {
	if base.FormTag != "" || !strings.HasPrefix(strings.ToLower(base.Subject), "check-out ") {
		return nil
	}
	co := CheckOut{Message: &xscmsg.Message{Message: base}}
	if match := checkInRE.FindStringSubmatch(base.Body); match != nil {
		if match[3] != "" {
			co.TacCall = match[1]
			co.TacName = strings.TrimSpace(match[2])
			co.OpCall = match[3]
			co.OpName = strings.TrimSpace(match[4])
		} else {
			co.OpCall = match[1]
			co.OpName = strings.TrimSpace(match[2])
		}
	}
	return &co
}

// Type returns the type of the message.
func (m *CheckOut) Type() *typedmsg.MessageType { return &Type }

// View returns the set of viewable fields of the message.
func (m *CheckOut) View() []xscmsg.LabelValue {
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
func (m *CheckOut) Edit() []xscmsg.Field {
	if m.fields == nil {
		m.fields = []xscmsg.Field{
			xscmsg.NewToField(&m.To),
			xscmsg.NewOriginMsgIDField(&m.OriginMsgID),
			xscmsg.NewHandlingField(&m.Handling),
			&tacCallField{xscmsg.NewBaseField(&m.TacCall)},
			&tacNameField{xscmsg.NewBaseField(&m.TacName), &m.TacCall},
			&opCallField{xscmsg.NewBaseField(&m.OpCall)},
			&opNameField{xscmsg.NewBaseField(&m.OpName)},
		}
	}
	return m.fields
}

// GetBody is a no-op for check-out messages.
func (m *CheckOut) GetBody() string { return "" }

// SetBody is a no-op for check-out messages.
func (m *CheckOut) SetBody(string) {}

// SetSubject is a no-op for check-out messages.
func (m *CheckOut) SetSubject(string) {}

// GetOpCall returns the value of the operator call sign field.
func (m *CheckOut) GetOpCall() string { return m.OpCall }

// SetOperator sets the OpCall and OpName fields of the message (outgoing only).
func (m *CheckOut) SetOperator(received bool, call string, name string) {
	if !received {
		m.OpCall, m.OpName = call, name
	}
}

// SetTactical sets the TacCall and TacName fields of the message.
func (m *CheckOut) SetTactical(call string, name string) {
	m.TacCall, m.TacName = call, name
}

// Save renders the message for saving to local storage.
func (m *CheckOut) Save() string {
	m.encode()
	return m.Message.Save()
}

// Transmit renders the message for transmission.
func (m *CheckOut) Transmit() ([]string, string, string) {
	m.encode()
	return m.Message.Transmit()
}

func (m *CheckOut) encode() {
	if m.TacCall != "" {
		m.Subject = fmt.Sprintf("Check-Out %s, %s", m.TacCall, m.TacName)
		m.Body = fmt.Sprintf("Check-Out %s, %s\n%s, %s\n", m.TacCall, m.TacName, m.OpCall, m.OpName)
	} else {
		m.Subject = fmt.Sprintf("Check-Out %s, %s", m.OpCall, m.OpName)
		m.Body = fmt.Sprintf("Check-Out %s, %s\n", m.OpCall, m.OpName)
	}
}

// Validate checks the validity of the message and returns strings
// describing the issues.  This is used for received messages, and only
// checks for validity issues that Outpost and/or PackItForms check for.
func (m *CheckOut) Validate() []string {
	return xscmsg.ValidateFromFieldProblems(m)
}

type tacCallField struct{ *xscmsg.BaseField }

var tacCallSignRE = regexp.MustCompile(`^[A-Z][A-Z0-9]{2,5}$`)

func (f tacCallField) Label() string         { return "Tactical Call Sign" }
func (f tacCallField) Size() (w, h int)      { return xscmsg.CallSignWidth, 1 }
func (f tacCallField) SetValue(value string) { f.BaseField.SetValue(strings.ToUpper(value)) }
func (f tacCallField) Problem() string {
	if f.Value() != "" && !tacCallSignRE.MatchString(f.Value()) {
		return "The tactical call sign is not valid.  Valid tactical call signs contain three to six letters or digits, starting with a letter."
	}
	return ""
}
func (f tacCallField) Help() string {
	return "This is the tactical call sign of the station that is checking out."
}

type tacNameField struct {
	*xscmsg.BaseField
	tacCall *string
}

func (f tacNameField) Label() string    { return "Tactical Name" }
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
	return "This is the name of the station that is checking out.  If a tactical call sign is supplied, a tactical station name is required.  Otherwise, a tactical station name is not allowed."
}

type opCallField struct{ *xscmsg.BaseField }

func (f opCallField) Label() string         { return "Operator Call Sign" }
func (f opCallField) Size() (w, h int)      { return xscmsg.CallSignWidth, 1 }
func (f opCallField) SetValue(value string) { f.BaseField.SetValue(strings.ToUpper(value)) }
func (f opCallField) Problem() string {
	if f.Value() == "" {
		return "The operator call sign is required."
	}
	if !xscmsg.CallSignValid(f.Value()) {
		return "The operator call sign is not a valid FCC call sign."
	}
	return ""
}
func (f opCallField) Help() string {
	return "This is the FCC call sign of the operator who is checking out."
}

type opNameField struct{ *xscmsg.BaseField }

func (f opNameField) Label() string    { return "Operator Name" }
func (f opNameField) Size() (w, h int) { return 80, 1 }
func (f opNameField) Problem() string {
	if f.Value() == "" {
		return "The operator name is required."
	}
	return ""
}
func (f opNameField) Help() string {
	return "This is the name of the operator who is checking out.  It is required."
}
