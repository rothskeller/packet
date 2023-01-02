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
	Code: "FormCorrupt",
	detect: func(a *Analysis) bool {
		if pktmsg.IsForm(a.msg.Body) {
			if a.xsc.Type.Tag == xscmsg.PlainTextTag {
				return true
			}
		}
		return false
	},
}

// ProbFormInvalid is raised when the form has invalid field values.
var ProbFormInvalid = &Problem{
	Code:  "FormInvalid",
	ifnot: []*Problem{ProbFormCorrupt},
	detect: func(a *Analysis) bool {
		if a.xsc.Type.Tag != xscmsg.PlainTextTag {
			if problems := a.xsc.Validate(true); len(problems) != 0 {
				return true
			}
		}
		return false
	},
	Variables: variableMap{
		"PROBLEMS": func(a *Analysis) string {
			return "    " + strings.Join(a.xsc.Validate(true), "\n    ")
		},
	},
}
