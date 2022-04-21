package analyze

import "steve.rothskeller.net/packet/pktmsg"

// Problem codes
const (
	ProblemBounceMessage   = "BounceMessage"
	ProblemDeliveryReceipt = "DeliveryReceipt"
	ProblemMessageCorrupt  = "MessageCorrupt"
	ProblemReadReceipt     = "ReadReceipt"
)

func init() {
	ProblemLabel[ProblemBounceMessage] = "message has no return address (probably auto-response)"
	ProblemLabel[ProblemDeliveryReceipt] = "DELIVERED receipt message"
	ProblemLabel[ProblemReadReceipt] = "unexpected READ receipt message"
	ProblemLabel[ProblemMessageCorrupt] = "message could not be parsed"
}

// checkNonHuman looks for non-human-generated messages and generates
// appropriate problem responses for them.
func (a *Analysis) checkNonHuman() {
	// First, look for messages that couldn't be parsed.  Those will be
	// listed in the weekly report but will not get any response and will
	// not count as checkins.
	if pe := a.msg.Base().ParseError; pe != "" {
		a.problems = append(a.problems, &problem{
			code: ProblemMessageCorrupt,
		})
		return
	}
	// Next, look for messages without a return address.  These are bounce
	// messages, vacation auto-responders, or similar auto-generated
	// messages.  Again, they will be listed in the weekly report but will
	// not get any response and will not count as checkins.
	if a.msg.Base().ReturnAddress == "" {
		a.problems = append(a.problems, &problem{
			code: ProblemBounceMessage,
		})
		return
	}
	// Next, look for delivery receipt messages.  These are not listed in
	// the weekly report, they get no response, and they don't count as
	// check-ins.
	if _, ok := a.msg.(*pktmsg.RxDeliveryReceipt); ok {
		a.problems = append(a.problems, &problem{
			code: ProblemDeliveryReceipt,
		})
	}
	// Finally, look for read receipt messages.  Those get a response back
	// to the sender telling them to turn off read receipts.  They don't
	// count as checkins.
	if _, ok := a.msg.(*pktmsg.RxReadReceipt); ok {
		a.problems = append(a.problems, &problem{
			code: ProblemReadReceipt,
			response: `
This message is an Outpost "read receipt", which should not have been sent.
Most likely, your Outpost installation has the "Auto-Read Receipt" setting
turned on.  The SCCo-standard Outpost configuration specifies that this setting
should be turned off.  You can find it on the Receipts tab of the Message
Settings dialog in Outpost.
`,
			references: refOutpostConfig,
		})
	}
}
