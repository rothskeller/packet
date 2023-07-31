package ics213

import (
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

// Compare compares two messages.  It returns a score indicating how closely
// they match, and the detailed comparisons of each field in the message.  The
// comparison is not symmetric:  the receiver of the call is the "expected"
// message and the argument is the "actual" message.
func (exp *ICS213) Compare(actual message.Message) (score, outOf int, fields []*message.CompareField) {
	var (
		act *ICS213
		ok  bool
	)
	if act, ok = actual.(*ICS213); !ok {
		return 0, 1, []*message.CompareField{{
			Label: "Message Type",
			Score: 0, OutOf: 1,
			Expected:     Type.Name,
			ExpectedMask: strings.Repeat("*", len(Type.Name)),
			Actual:       actual.Type().Name,
			ActualMask:   strings.Repeat("*", len(actual.Type().Name)),
		}}
	}
	fields = []*message.CompareField{
		common.CompareDate("Date", exp.Date, act.Date),
		common.CompareTime("Time", exp.Time, act.Time),
		common.CompareExact("Situation Severity", exp.Severity, act.Severity),
		common.CompareExact("Handling", exp.Handling, act.Handling),
		common.CompareExact("Take Action", exp.TakeAction, act.TakeAction),
		common.CompareExact("Reply", exp.Reply, act.Reply),
		common.CompareText("Reply By", exp.ReplyBy, act.ReplyBy),
		common.CompareCheckbox("For Your Information", exp.FYI, act.FYI),
		common.CompareText("To ICS Position", exp.ToICSPosition, act.ToICSPosition),
		common.CompareText("To Location", exp.ToLocation, act.ToLocation),
		common.CompareText("To Name", exp.ToName, act.ToName),
		common.ComparePhoneNumber("To Telephone", exp.ToTelephone, act.ToTelephone),
		common.CompareText("From ICS Position", exp.FromICSPosition, act.FromICSPosition),
		common.CompareText("From Location", exp.FromLocation, act.FromLocation),
		common.CompareText("From Name", exp.FromName, act.FromName),
		common.ComparePhoneNumber("From Telephone", exp.FromTelephone, act.FromTelephone),
		common.CompareText("Subject", exp.Subject, act.Subject),
		common.CompareExact("Reference", exp.Reference, act.Reference),
		common.CompareText("Message", exp.Message, act.Message),
	}
	return common.ConsolidateCompareFields(fields)
}
