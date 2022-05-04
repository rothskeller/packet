package analyze

import (
	"fmt"
)

// Problem codes
const (
	ProblemCallSignConflict = "CallSignConflict"
	ProblemNoCallSign       = "NoCallSign"
)

func init() {
	ProblemLabel[ProblemCallSignConflict] = "call sign conflict"
	ProblemLabel[ProblemNoCallSign] = "no call sign in message"
}

// checkCallSign looks for the sender's call sign in various places, and makes
// sure they all agree.
func (a *Analysis) checkCallSign() {
	var (
		formCS string
		isform bool
	)
	// Extract the call sign from the the OpCall field of the form, if any.
	if op, ok := a.xsc.(interface{ Operator() (string, string) }); ok {
		_, formCS = op.Operator()
		isform = true
	}
	// If we didn't find a call sign anywhere, we can't count this message.
	// (The text of the response differs depending on whether it's a form
	// message.)
	if a.msg.FromCallSign() == "" && a.subjectCallSign == "" && formCS == "" {
		if isform {
			a.problems = append(a.problems, &problem{
				code: ProblemNoCallSign,
				response: fmt.Sprintf(`
This message cannot be counted because it's not clear who sent it.  There
is no call sign in the return address, or after the word "Practice" on the
subject line, or in the Operator Call field of the form.  In order for a
message to count, there must be a call sign in at least one of those places.
`),
			})
		} else {
			a.problems = append(a.problems, &problem{
				code: ProblemNoCallSign,
				response: fmt.Sprintf(`
This message cannot be counted because it's not clear who sent it.  There
is no call sign in the return address or after the word "Practice" on the
subject line.  In order for a message to count, there must be a call sign in
at least one of those places.
`),
			})
		}
		return
	}
	// We did find call signs in one or more of those places.  The one we'll
	// count is the one from the subject line, if present; else the form
	// OpCall, if present; else the return address.
	if a.subjectCallSign != "" {
		a.fromCallSign = a.subjectCallSign
	} else if formCS != "" {
		a.fromCallSign = formCS
	} else {
		a.fromCallSign = a.msg.FromCallSign()
	}
	// If the one in the return address doesn't match the one we chose,
	// that's OK.  But if the one in the form doesn't match the one from the
	// subject line, that's a problem to be reported.
	if formCS != "" && formCS != a.fromCallSign {
		a.problems = append(a.problems, &problem{
			code: ProblemCallSignConflict,
			response: fmt.Sprintf(`
This message has conflicting call signs.  The Subject line says the call sign
is %s, but the Operator Call Sign field of the form says %s.  The two should
agree.  (This message will be counted as a practice attempt by %s.)
`, a.subjectCallSign, formCS, a.fromCallSign),
		})
	}
}
