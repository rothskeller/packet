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

// fromCallSign extracts the sender's call sign from the several different
// places it might appear.
func (a *Analysis) fromCallSign() string {
	// If there's a call sign on the subject line, use that.
	if a.Practice != nil && a.Practice.CallSign != "" {
		return a.Practice.CallSign
	}
	// If there's a call sign in the From address of the message, use that.
	if cs := a.msg.FromCallSign(); cs != "" {
		// If it's a valid FCC call sign, use it.
		if fccCallSignRE.MatchString(cs) {
			return cs
		}
		// It's not a valid FCC call sign.  If it's from within our BBS
		// network, assume it's a tactical call sign and use it.
		// If it's from within our BBS network, use it.
		if a.msg.FromBBS() != "" {
			return cs
		}
		// It's not a valid FCC call sign and it's not from within our
		// BBS network, so it doesn't count.
	}
	// If the message is a form and has a call sign in an OpCall field, use
	// that.
	if f := a.xsc.KeyField(xscmsg.FOpCall); f != nil && f.Value != "" {
		return strings.ToUpper(f.Value)
	}
	// No call sign found.  (Problem will be reported later.)
	return ""
}

// ProbFormNoCallSign is raised if we can't find the sender's call sign in a
// form message.
var ProbFormNoCallSign = &Problem{
	Code: "FormNoCallSign",
	detect: func(a *Analysis) bool {
		// This check only applies to forms with an OpCall field.
		if a.xsc.KeyField(xscmsg.FOpCall) == nil {
			return false
		}
		// The check.
		return a.FromCallSign == ""
	},
}

// ProbNoCallSign is raised if we can't find the sender's call sign in a
// non-form message.
var ProbNoCallSign = &Problem{
	Code: "NoCallSign",
	detect: func(a *Analysis) bool {
		// This check only applies to non-forms or to forms that don't
		// have an OpCall field.
		if a.xsc.KeyField(xscmsg.FOpCall) != nil {
			return false
		}
		// The check.
		return a.FromCallSign == ""
	},
}

// ProbCallSignConflict is raised if the call sign in the OpCall field of a form
// doesn't match the call sign after "Practice" on the subject line.  This check
// only applies if the latter call sign is an FCC call sign; if it's a tactical
// call, a mismatch is OK.
var ProbCallSignConflict = &Problem{
	Code:  "CallSignConflict",
	ifnot: []*Problem{ProbFormNoCallSign},
	detect: func(a *Analysis) bool {
		// This check only applies to forms with a call sign in an
		// OpCall field.
		var formCS string
		if f := a.xsc.KeyField(xscmsg.FOpCall); f != nil && f.Value != "" {
			formCS = strings.ToUpper(f.Value)
		} else {
			return false
		}
		// This check only applies if the detected FromCallSign is a
		// valid FCC call sign.
		if !fccCallSignRE.MatchString(a.FromCallSign) {
			return false
		}
		// The check.
		return formCS != a.FromCallSign
	},
	Variables: variableMap{
		"OPCALL": func(a *Analysis) string {
			return strings.ToUpper(a.xsc.KeyField(xscmsg.FOpCall).Value)
		},
	},
}
