package analyze

import (
	"strings"

	"github.com/rothskeller/packet/xscmsg"
)

func init() {
	Problems[ProbCallSignConflict.Code] = ProbCallSignConflict
	Problems[ProbNoCallSign.Code] = ProbNoCallSign
}

// ProbNoCallSign is raised if we can't find the sender's call sign in the
// message.
var ProbNoCallSign = &Problem{
	Code:  "NoCallSign",
	after: []*Problem{ProbPracticeSubjectFormat, ProbDeliveryReceipt}, // set a.subjectCallSign, a.xsc
	ifnot: []*Problem{ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		var fromCS, formCS string

		// Extract the call sign from the the OpCall field of the form,
		// if any.
		if f := a.xsc.Field(xscmsg.FOpCall); f != nil {
			formCS = strings.ToUpper(f.Value)
		}
		// Extract the call sign from the From address of the message.
		// If we find a non-FCC call and the message is coming from
		// outside the BBS network, it doesn't count.  We only accept
		// tactical calls on the From line when within the BBS network.
		fromCS = a.msg.FromCallSign()
		if !fccCallSignRE.MatchString(fromCS) && a.msg.FromBBS() == "" {
			fromCS = ""
		}
		// If we didn't find a call sign anywhere, report the problem.
		if fromCS == "" && a.subjectCallSign == "" && formCS == "" {
			if f := a.xsc.Field(xscmsg.FOpCall); f != nil {
				return true, "form"
			}
			return true, "plain"
		}
		// We did find call signs in one or more of those places.  Now
		// we need to figure out which one to count as "the" from call
		// sign.  We give precedence to the subject line, then the From
		// line, then the form OpCall field.
		switch {
		case a.subjectCallSign != "":
			a.fromCallSign = a.subjectCallSign
		case fromCS != "":
			a.fromCallSign = fromCS
		default:
			a.fromCallSign = formCS
		}
		return false, ""
	},
}

// ProbCallSignConflict is raised if the call sign in the OpCall field of a form
// doesn't match the call sign after "Practice" on the subject line.  This check
// only applies if the latter call sign is an FCC call sign; if it's a tactical
// call, a mismatch is OK.
var ProbCallSignConflict = &Problem{
	Code:  "CallSignConflict",
	after: []*Problem{ProbNoCallSign, ProbDeliveryReceipt}, // set a.fromCallSign, a.xsc
	ifnot: []*Problem{ProbNoCallSign, ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		if !fccCallSignRE.MatchString(a.fromCallSign) {
			// The from call sign is a tactical call, so the form
			// OpCall is allowed to be different.
			return false, ""
		}
		if f := a.xsc.Field(xscmsg.FOpCall); f != nil {
			formCS := strings.ToUpper(f.Value)
			if formCS != "" && formCS != a.fromCallSign {
				return true, ""
			}
		}
		return false, ""
	},
	Variables: variableMap{
		"OPCALL": func(a *Analysis) string {
			return strings.ToUpper(a.xsc.Field(xscmsg.FOpCall).Value)
		},
	},
}
