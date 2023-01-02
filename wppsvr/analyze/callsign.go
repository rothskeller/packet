package analyze

import (
	"strings"

	"github.com/rothskeller/packet/xscmsg"
)

func init() {
	Problems[ProbCallSignConflict.Code] = ProbCallSignConflict
	Problems[ProbFormNoCallSign.Code] = ProbFormNoCallSign
	Problems[ProbNoCallSign.Code] = ProbNoCallSign
}

func (a *Analysis) fromCallSign() string {
	var fromCS, subjCS, formCS string

	// Extract the call sign from the the OpCall field of the form, if any.
	if f := a.xsc.KeyField(xscmsg.FOpCall); f != nil {
		formCS = strings.ToUpper(f.Value)
	}
	// Extract the call sign from the practice subject, if any.
	if a.Practice != nil {
		subjCS = a.Practice.CallSign
	}
	// Extract the call sign from the From address of the message.  If we
	// find a non-FCC call and the message is coming from outside the BBS
	// network, it doesn't count.  We only accept tactical calls on the From
	// line when within the BBS network.
	fromCS = a.msg.FromCallSign()
	if !fccCallSignRE.MatchString(fromCS) && a.msg.FromBBS() == "" {
		fromCS = ""
	}
	// Choose "the" from call sign.  We give precedence to the subject line,
	// then the From line, then the form OpCall field.
	switch {
	case subjCS != "":
		return subjCS
	case fromCS != "":
		return fromCS
	default:
		return formCS
	}
}

// ProbFormNoCallSign is raised if we can't find the sender's call sign in a
// form message.
var ProbFormNoCallSign = &Problem{
	Code: "FormNoCallSign",
	detect: func(a *Analysis) bool {
		if a.FromCallSign == "" {
			if f := a.xsc.KeyField(xscmsg.FOpCall); f != nil {
				return true
			}
		}
		return false
	},
}

// ProbNoCallSign is raised if we can't find the sender's call sign in a
// non-form message.
var ProbNoCallSign = &Problem{
	Code: "NoCallSign",
	detect: func(a *Analysis) bool {
		if a.FromCallSign == "" {
			if f := a.xsc.KeyField(xscmsg.FOpCall); f == nil {
				return true
			}
		}
		return false
	},
}

// ProbCallSignConflict is raised if the call sign in the OpCall field of a form
// doesn't match the call sign after "Practice" on the subject line.  This check
// only applies if the latter call sign is an FCC call sign; if it's a tactical
// call, a mismatch is OK.
var ProbCallSignConflict = &Problem{
	Code:  "CallSignConflict",
	ifnot: []*Problem{ProbNoCallSign},
	detect: func(a *Analysis) bool {
		if !fccCallSignRE.MatchString(a.FromCallSign) {
			// The from call sign is a tactical call, so the form
			// OpCall is allowed to be different.
			return false
		}
		if f := a.xsc.KeyField(xscmsg.FOpCall); f != nil {
			formCS := strings.ToUpper(f.Value)
			if formCS != "" && formCS != a.FromCallSign {
				return true
			}
		}
		return false
	},
	Variables: variableMap{
		"OPCALL": func(a *Analysis) string {
			return strings.ToUpper(a.xsc.KeyField(xscmsg.FOpCall).Value)
		},
	},
}
