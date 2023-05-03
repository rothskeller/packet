package xscmsg

import (
	"fmt"

	"github.com/rothskeller/packet/pktmsg"
)

// Message is the base implementation of xscmsg.IMessage, suitable for embedding
// in implementations of XSC message types.  Note that it is not a complete
// implementation; it is missing the Type and View methods.
type Message struct{ *pktmsg.Message }

// GetOpCall returns the value of the operator call sign field if the message
// type has such a field, or an empty string otherwise.
func (m Message) GetOpCall() string { return "" }

// GetToICSPosition returns the value of the "To ICS Position" field if the
// message type has such a field, or an empty string otherwise.
func (m Message) GetToICSPosition() string { return "" }

// GetToLocation returns the value of the "To Location" field if the message
// type has such a field, or an empty string otherwise.
func (m Message) GetToLocation() string { return "" }

// SetDestinationMsgID sets the destination message ID for the message, if the
// message type has such a field.  It is a no-op otherwise.
func (m *Message) SetDestinationMsgID(string) {}

// SetOperator sets the OpCall, OpName, OpDate, and OpTime fields of the message
// if it has them.  For messages without these fields, this is a no-op.
func (m *Message) SetOperator(bool, string, string) {}

// SetReference sets the reference field of the message, if it has one.
// It's a no-op otherwise.
func (m *Message) SetReference(string) {}

// SetTactical sets the TacCall and TacName fields of the message if it has
// them; otherwise, this is a no-op.
func (m *Message) SetTactical(string, string) {}

// Validate checks the validity of the message and returns strings describing
// the issues.  This is used for received messages, and only checks for validity
// issues that Outpost and/or PackItForms check for.
func (m Message) Validate() []string { return nil }

// ValidateFromFieldProblems implements Validate by calling Problem on all
// fields returned by Edit.  It is the most common implementation of Validate.
func ValidateFromFieldProblems(m Editable) (problems []string) {
	for _, f := range m.Edit() {
		if problem := f.Problem(); problem != "" {
			problems = append(problems, problem)
		}
	}
	return problems
}

// ViewHeaders returns the set of viewable fields for the message envelope, not
// including the subject.
func (m *Message) ViewHeaders() (lvs []LabelValue) {
	lvs = []LabelValue{
		{"From", m.From.String()},
		{"To", m.To.String()},
	}
	if m.SentDate != "" {
		lvs = append(lvs, LabelValue{"Sent", m.SentDate.Time().Format("01/02/2006 15:04")})
	}
	if m.RxBBS != "" {
		var received string
		if m.RxArea != "" {
			received = fmt.Sprintf("%s from %s on %s", m.RxDate.Time().Format("01/02/2006 15:04"), m.RxArea, m.RxBBS)
		} else {
			received = fmt.Sprintf("%s from %s", m.RxDate.Time().Format("01/02/2006 15:04"), m.RxBBS)
		}
		lvs = append(lvs, LabelValue{"Received", received})
	}
	return lvs
}
