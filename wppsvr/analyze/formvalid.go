package analyze

import (
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
)

func init() {
	Problems[ProbFormCorrupt.Code] = ProbFormCorrupt
	Problems[ProbFormInvalid.Code] = ProbFormInvalid
}

// ProbFormCorrupt is raised when the body looks like a form encoding but
// couldn't be successfully decoded.
var ProbFormCorrupt = &Problem{
	Code:  "FormCorrupt",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		if pktmsg.IsForm(a.msg.Body) {
			if a.xsc.Type.Tag == xscmsg.PlainTextTag {
				return true, ""
			}
		}
		return false, ""
	},
}

// ProbFormInvalid is raised when the form has invalid field values.
var ProbFormInvalid = &Problem{
	Code:  "FormInvalid",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbFormCorrupt, ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		if a.xsc.Type.Tag != xscmsg.PlainTextTag {
			if problems := a.xsc.Validate(true); len(problems) != 0 {
				return true, ""
			}
		}
		return false, ""
	},
	Variables: variableMap{
		"PROBLEMS": func(a *Analysis) string {
			return "    " + strings.Join(a.xsc.Validate(true), "\n    ")
		},
	},
}
