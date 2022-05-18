package analyze

import (
	"github.com/rothskeller/packet/xscmsg/delivrcpt"
	"github.com/rothskeller/packet/xscmsg/readrcpt"
)

// Problem codes
const (
	ProblemDeliveryReceipt = "DeliveryReceipt"
	ProblemReadReceipt     = "ReadReceipt"
)

func init() {
	ProblemLabel[ProblemDeliveryReceipt] = "DELIVERED receipt message"
	ProblemLabel[ProblemReadReceipt] = "unexpected READ receipt message"
}

// checkReceipts looks for delivery and read receipt messages and generates
// appropriate problem responses for them.
func (a *Analysis) checkReceipts() bool {
	if _, ok := a.xsc.(*delivrcpt.DeliveryReceipt); ok {
		a.problems = append(a.problems, &problem{code: ProblemDeliveryReceipt})
		return true
	}
	if _, ok := a.xsc.(*readrcpt.ReadReceipt); ok {
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
		return true
	}
	return false
}
