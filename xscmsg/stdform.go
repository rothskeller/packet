package xscmsg

import "time"

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

// Validate checks the validity of the message and returns strings describing
// the issues.  This is used for received messages, and only checks for validity
// issues that Outpost and/or PackItForms check for.
func (m *StdForm) Validate() (problems []string) {
	if m.MessageDate == "" {
		problems = append(problems, "The message date is required.")
	} else if !DateValid(m.MessageDate) {
		problems = append(problems, "The message date is not a valid date.  Dates must be in MM/DD/YYYY format.")
	}
	if m.MessageTime == "" {
		problems = append(problems, "The message time is required.")
	} else if !TimeValid(m.MessageTime) {
		problems = append(problems, "The message time is not a valid time.  Times must be in HH:MM format (24-hour clock).")
	}
	if m.Handling == "" {
		problems = append(problems, "The message handling order is required.")
	} else if !ChoicesValid(m.Handling, handlingChoices) {
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
	} else if !CallSignValid(m.OpCall) {
		problems = append(problems, "The operator call sign is not a valid FCC call sign.")
	}
	if m.OpDate != "" && !DateValid(m.OpDate) {
		problems = append(problems, "The operator date is not a valid date.  Dates must be in MM/DD/YYYY format.")
	}
	if m.OpTime != "" && !TimeValid(m.OpTime) {
		problems = append(problems, "The operator time is not a valid time.  Times must be in HH:MM format (24-hour clock).")
	}
	return problems
}

// View returns the set of viewable fields in the message.
func (m *StdForm) View() (lvs1, lvs2 []LabelValue) {
	lvs1 = m.ViewHeaders()
	lvs1 = append(lvs1, []LabelValue{
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
	lvs2 = []LabelValue{
		{"Operator Relay Received", m.OpRelayRcvd},
		{"Operator Relay Sent", m.OpRelaySent},
	}
	if m.RxBBS != "" {
		lvs2 = append(lvs2, []LabelValue{
			{"Received by", smartJoin(m.OpName, m.OpCall, " ")},
			{"Received at", smartJoin(m.OpDate, m.OpTime, " ")},
		}...)
	} else {
		lvs2 = append(lvs2, []LabelValue{
			{"Sent by", smartJoin(m.OpName, m.OpCall, " ")},
			{"Sent at", smartJoin(m.OpDate, m.OpTime, " ")},
		}...)
	}
	return
}
func smartJoin(a, b, sep string) string {
	if a == "" || b == "" {
		return a + b
	}
	return a + sep + b
}
