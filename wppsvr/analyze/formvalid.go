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
		return a.xsc.Type.Tag == xscmsg.PlainTextTag && pktmsg.IsForm(a.msg.Body)
	},
}

// ProbFormInvalid is raised when the form has invalid field values.
var ProbFormInvalid = &Problem{
	Code: "FormInvalid",
	detect: func(a *Analysis) bool {
		// This message only applies to forms.
		if a.xsc.Type.Tag == xscmsg.PlainTextTag {
			return false
		}
		// The check.
		return len(a.xsc.Validate(true)) != 0
	},
	Variables: variableMap{
		"PROBLEMS": func(a *Analysis) string {
			return "    " + strings.Join(a.xsc.Validate(true), "\n    ")
		},
	},
}
