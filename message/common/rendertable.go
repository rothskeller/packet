package common

import "github.com/rothskeller/packet/message"

// RenderTable1 returns the initial set of label/value pairs for standard
// fields.
func (sf *StdFields) RenderTable1() []message.LabelValue {
	return []message.LabelValue{
		{Label: "Origin Message Number", Value: sf.OriginMsgID},
		{Label: "Destination Message Number", Value: sf.DestinationMsgID},
		{Label: "Message Date/Time", Value: SmartJoin(sf.MessageDate, sf.MessageTime, " ")},
		{Label: "Handling Order", Value: sf.Handling},
		{Label: "To ICS Position", Value: sf.ToICSPosition},
		{Label: "To Location", Value: sf.ToLocation},
		{Label: "To Name", Value: sf.ToName},
		{Label: "To Contact Info", Value: sf.ToContact},
		{Label: "From ICS Position", Value: sf.FromICSPosition},
		{Label: "From Location", Value: sf.FromLocation},
		{Label: "From Name", Value: sf.FromName},
		{Label: "From Contact Info", Value: sf.FromContact},
	}
}

// RenderTable2 returns the final set of label/value pairs for standard fields.
func (sf *StdFields) RenderTable2() []message.LabelValue {
	return []message.LabelValue{
		{Label: "Handled", Value: SmartJoin(sf.OpDate, sf.OpTime, " ")},
		{Label: "Operator", Value: SmartJoin(sf.OpName, sf.OpCall, " ")},
		{Label: "Relay From", Value: sf.OpRelayRcvd},
		{Label: "Relay Through", Value: sf.OpRelaySent},
	}
}

// SmartJoin joins the two provided strings with the provided separator, but
// only when both are non-empty.  If one is empty, it returns the other.  If
// both are empty, it returns an empty string.
func SmartJoin(a, b, sep string) string {
	if a != "" && b != "" {
		return a + sep + b
	}
	if a == "" {
		return b
	}
	return a
}
