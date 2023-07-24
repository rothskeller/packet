package ics213

import (
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

// RenderTable renders the message as a set of field label / field value pairs,
// intended for read-only display to a human.
func (f *ICS213) RenderTable() []message.LabelValue {
	var oplabel, action, reply string

	reply = f.Reply
	if f.ReplyBy != "" {
		reply = common.SmartJoin(reply, "by "+f.ReplyBy, " ")
	}
	switch f.ReceivedSent {
	case "sender":
		oplabel = "Sent"
	case "receiver":
		oplabel = "Received"
	default:
		oplabel = "Action"
	}
	action = common.SmartJoin(f.OpDate, f.OpTime, " ")
	if f.TxMethod == "Other" && f.OtherMethod != "" {
		action = common.SmartJoin(action, "via "+f.OtherMethod, " ")
	} else if f.TxMethod != "" {
		action = common.SmartJoin(action, "via "+f.TxMethod, " ")
	}
	return []message.LabelValue{
		{Label: "Origin Message Number", Value: f.OriginMsgID},
		{Label: "Destination Message Number", Value: f.DestinationMsgID},
		{Label: "Date and Time", Value: common.SmartJoin(f.Date, f.Time, " ")},
		{Label: "Handling Order", Value: f.Handling},
		{Label: "Take Action", Value: f.TakeAction},
		{Label: "Reply", Value: reply},
		{Label: "To ICS Position", Value: f.ToICSPosition},
		{Label: "To Location", Value: f.ToLocation},
		{Label: "To Name", Value: f.ToName},
		{Label: "To Telephone #", Value: f.ToTelephone},
		{Label: "From ICS Position", Value: f.FromICSPosition},
		{Label: "From Location", Value: f.FromLocation},
		{Label: "From Name", Value: f.FromName},
		{Label: "From Telephone #", Value: f.FromTelephone},
		{Label: "Subject", Value: f.Subject},
		{Label: "Reference", Value: f.Reference},
		{Label: "Message", Value: f.Message},
		{Label: oplabel, Value: action},
		{Label: "Operator", Value: common.SmartJoin(f.OpName, f.OpCall, " ")},
		{Label: "Relay From", Value: f.OpRelayRcvd},
		{Label: "Relay Through", Value: f.OpRelaySent},
	}
}
