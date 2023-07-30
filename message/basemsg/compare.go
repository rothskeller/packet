package basemsg

import "github.com/rothskeller/packet/message"

// TODO defining this functions on BaseMessage means that only messages for
// which it's appropriate can leverage BaseMessage.

// Compare compares two messages.  It returns a score indicating how closely
// they match, and the detailed comparisons of each field in the message.  The
// comparison is not symmetric:  the receiver of the call is the "expected"
// message and the argument is the "actual" message.
func (bm *BaseMessage) Compare(actual message.Message) (score int, outOf int, cfields []*message.CompareField) {
	var act *BaseMessage

	if actual.Type() != bm.MessageType {
		return 0, 1, []*message.CompareField{{
			Label: "Message Type",
			Score: 0, OutOf: 1,
			Expected:     bm.MessageType.Name,
			ExpectedMask: "*",
			Actual:       actual.Type().Name,
			ActualMask:   "*",
		}}
	}
	act = actual.(interface{ Base() *BaseMessage }).Base() // TODO: remove cast when Base is part of message.Message
	for i, expf := range bm.Fields {
		actf := act.Fields[i]
		if expf.Compare != nil && (*expf.Value != "" || *actf.Value != "") {
			comp := expf.Compare(expf.Label, *expf.Value, *actf.Value)
			if comp.OutOf == 0 {
				comp.OutOf = 1
			}
			score, outOf = score+comp.Score, outOf+comp.OutOf
			cfields = append(cfields, comp)
		}
	}
	return score, outOf, cfields
}
