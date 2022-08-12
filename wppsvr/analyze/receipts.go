package analyze

import (
	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/wppsvr/config"
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
		if a.xsc = xscmsg.Recognize(a.msg, true); a.xsc == nil {
			if form := pktmsg.ParseForm(a.msg.Body, true); form != nil {
				a.xsc = &config.UnknownForm{M: a.msg, F: form}
			} else {
				a.xsc = &config.PlainTextMessage{M: a.msg}
			}
		}
		// Is it a delivery receipt?
		_, ok := a.xsc.(*delivrcpt.DeliveryReceipt)
		return ok, ""
	},
}

// ProbReadReceipt is raised for any READ receipt message.
var ProbReadReceipt = &Problem{
	Code:  "ReadReceipt",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbBounceMessage},
	detect: func(a *Analysis) (bool, string) {
		_, ok := a.xsc.(*readrcpt.ReadReceipt)
		return ok, ""
	},
}
