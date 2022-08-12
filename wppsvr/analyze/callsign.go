package analyze

func init() {
	Problems[ProbCallSignConflict.Code] = ProbCallSignConflict
	Problems[ProbNoCallSign.Code] = ProbNoCallSign
}

type formWithOpCall interface {
	Operator() (string, string)
}

// ProbNoCallSign is raised if we can't find the sender's call sign in the
// message.
var ProbNoCallSign = &Problem{
	Code:  "NoCallSign",
	after: []*Problem{ProbPracticeSubjectFormat, ProbDeliveryReceipt}, // set a.subjectCallSign, a.xsc
	ifnot: []*Problem{ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		var formCS string

		// Extract the call sign from the the OpCall field of the form,
		// if any.
		if op, ok := a.xsc.(formWithOpCall); ok {
			_, formCS = op.Operator()
		}
		// If we didn't find a call sign anywhere, report the problem.
		if a.msg.FromCallSign() == "" && a.subjectCallSign == "" && formCS == "" {
			if _, ok := a.xsc.(formWithOpCall); ok {
				return true, "form"
			}
			return true, "plain"
		}
		// We did find call signs in one or more of those places.  The
		// one we'll count is the one from the subject line, if present;
		// else the form OpCall, if present; else the return address.
		if a.subjectCallSign != "" {
			a.fromCallSign = a.subjectCallSign
		} else if formCS != "" {
			a.fromCallSign = formCS
		} else {
			a.fromCallSign = a.msg.FromCallSign()
		}
		return false, ""
	},
}

// ProbCallSignConflict is raised if the call sign in the OpCall field of a form
// doesn't match the call sign after "Practice" on the subject line.
var ProbCallSignConflict = &Problem{
	Code:  "CallSignConflict",
	after: []*Problem{ProbNoCallSign, ProbDeliveryReceipt}, // set a.fromCallSign, a.xsc
	ifnot: []*Problem{ProbNoCallSign, ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		if op, ok := a.xsc.(formWithOpCall); ok {
			_, formCS := op.Operator()
			if formCS != "" && formCS != a.fromCallSign {
				return true, ""
			}
		}
		return false, ""
	},
	Variables: variableMap{
		"OPCALL": func(a *Analysis) string {
			_, formCS := a.xsc.(formWithOpCall).Operator()
			return formCS
		},
	},
}
