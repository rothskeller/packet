package checkout

import (
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

// RenderTable renders the message as a set of field label / field value pairs,
// intended for read-only display to a human.
func (m *CheckOut) RenderTable() []message.LabelValue {
	return []message.LabelValue{
		{Label: "Origin Message Number", Value: m.OriginMsgID},
		{Label: "Handling", Value: m.Handling},
		{Label: "Message", Value: "Check-Out"},
		{Label: "Tactical Station", Value: common.SmartJoin(m.TacticalCallSign, m.TacticalStationName, " ")},
		{Label: "Operator", Value: common.SmartJoin(m.OperatorCallSign, m.OperatorName, " ")},
	}
}
