package analyze

import (
	"fmt"
)

// Problem codes
const (
	ProblemFormSubjectConflict = "FormSubjectConflict"
)

func init() {
	ProblemLabel[ProblemFormSubjectConflict] = "message subject doesn't agree with form contents"
}

// checkFormSubject checks to be sure that the Subject line of the message
// matches the contents of the form.
func (a *Analysis) checkFormSubject() {
	var shouldbe string

	// This check is only performed on forms for which we know how to
	// compute the subject line, and then only if the form contents are
	// valid.
	if form, ok := a.msg.(interface{ Valid() bool }); !ok || !form.Valid() {
		return
	}
	if form, ok := a.msg.(interface{ EncodeSubjectLine() string }); ok {
		shouldbe = form.EncodeSubjectLine()
	} else {
		return
	}
	if sl := a.msg.Base().SubjectLine; sl != shouldbe {
		a.problems = append(a.problems, &problem{
			code:    ProblemFormSubjectConflict,
			subject: "Message subject doesn't agree with form contents",
			response: fmt.Sprintf(`
This message has
    Subject: %s
but, based on the contents of the form, it should have
    Subject: %s
PackItForms automatically generates the Subject line from the form contents; it
should not be overridden manually.
`, sl, shouldbe),
		})
	}
}
