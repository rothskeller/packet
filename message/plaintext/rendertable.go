package plaintext

import "github.com/rothskeller/packet/message"

// RenderTable renders the message as a set of field label / field value pairs,
// intended for read-only display to a human.
func (m *PlainText) RenderTable() []message.LabelValue {
	return []message.LabelValue{
		{Label: "Origin Message Number", Value: m.OriginMsgID},
		{Label: "Handling Order", Value: m.Handling},
		{Label: "Subject", Value: m.Subject},
		{Label: "Body", Value: m.Body},
	}
}
