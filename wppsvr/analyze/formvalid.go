package analyze

import (
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/wppsvr/config"
)

func init() {
	Problems[ProbFormCorrupt.Code] = ProbFormCorrupt
	Problems[ProbFormInvalid.Code] = ProbFormInvalid
}

type validatableForm interface {
	Validate(bool) []string
}

// ProbFormCorrupt is raised when the body looks like a form encoding but
// couldn't be successfully decoded.
var ProbFormCorrupt = &Problem{
	Code:  "FormCorrupt",
	Label: "incorrectly encoded form",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		if pktmsg.IsForm(a.msg.Body) {
			if _, ok := a.xsc.(*config.PlainTextMessage); ok {
				return true, ""
			}
		}
		return false, ""
	},
}

// ProbFormInvalid is raised when the form has invalid field values.
var ProbFormInvalid = &Problem{
	Code:  "FormInvalid",
	Label: "invalid form contents",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbFormCorrupt, ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		if xsc, ok := a.xsc.(validatableForm); ok {
			if problems := xsc.Validate(true); len(problems) != 0 {
				return true, ""
			}
		}
		return false, ""
	},
	Variables: variableMap{
		"PROBLEMS": func(a *Analysis) string {
			return "    " + strings.Join(a.xsc.(validatableForm).Validate(true), "\n    ")
		},
	},
}
