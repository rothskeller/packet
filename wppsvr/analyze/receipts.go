package analyze

import (
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/delivrcpt"
	"github.com/rothskeller/packet/xscmsg/readrcpt"
)

func init() {
	Problems[ProbDeliveryReceipt.Code] = ProbDeliveryReceipt
	Problems[ProbReadReceipt.Code] = ProbReadReceipt
}

// ProbDeliveryReceipt is raised for any delivery receipt message.  This check
// has the side effect of determining the message type and setting a.xsc.
var ProbDeliveryReceipt = &Problem{
	Code:  "DeliveryReceipt",
	ifnot: []*Problem{ProbBounceMessage},
	detect: func(a *Analysis) (bool, string) {
		// Find out whether the message is a known type.  If it's not,
		// we'll put it in our pseudo-type for plain text or unknown
		// form, whichever fits.
		a.xsc = xscmsg.Recognize(a.msg, true)
		// Is it a delivery receipt?
		return a.xsc.Type.Tag == delivrcpt.Tag, ""
	},
}

// ProbReadReceipt is raised for any READ receipt message.
var ProbReadReceipt = &Problem{
	Code:  "ReadReceipt",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbBounceMessage},
	detect: func(a *Analysis) (bool, string) {
		return a.xsc.Type.Tag == readrcpt.Tag, ""
	},
}
