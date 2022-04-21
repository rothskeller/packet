package analyze

import (
	"fmt"

	"steve.rothskeller.net/packet/pktmsg"
)

// Problem codes
const (
	ProblemFormDestination = "FormDestination"
)

func init() {
	ProblemLabel[ProblemFormDestination] = "incorrect destination for form"
}

// checkFormDestination determines whether the message has the correct
// destination.
func (a *Analysis) checkFormDestination() {
	// The correct destination depends on the form type, and we only check
	// it for the form types we know.
	switch msg := a.msg.(type) {
	case *pktmsg.RxEOC213RRForm:
		if msg.ToICSPosition != "Planning Section" && msg.ToLocation != "County EOC" {
			a.problems = append(a.problems, &problem{
				code:    ProblemFormDestination,
				subject: "Incorrect destination for form",
				response: fmt.Sprintf(`
This message form is addressed to ICS Position %q at Location %q.  EOC-213RR
messages should be addressed to "Planning Section" at "County EOC".
`, msg.ToICSPosition, msg.ToLocation),
				references: refFormRouting,
			})
		} else if msg.ToICSPosition != "Planning Section" {
			a.problems = append(a.problems, &problem{
				code:    ProblemFormDestination,
				subject: "Incorrect destination for form",
				response: fmt.Sprintf(`
This message form is addressed to ICS position %q.  EOC-213RR messages should
be addressed to "Planning Section".
`, msg.ToICSPosition),
				references: refFormRouting,
			})
		} else if msg.ToLocation != "County EOC" {
			a.problems = append(a.problems, &problem{
				code:    ProblemFormDestination,
				subject: "Incorrect destination for form",
				response: fmt.Sprintf(`
This message form is addressed to ICS Location %q.  EOC-213RR messages should
be addressed to "County EOC".
`, msg.ToLocation),
				references: refFormRouting,
			})
		}
	case *pktmsg.RxMuniStatForm:
		if msg.ToICSPosition != "Situation Analysis Unit" && msg.ToICSPosition != "Planning Section" &&
			msg.ToLocation != "County EOC" {
			a.problems = append(a.problems, &problem{
				code:    ProblemFormDestination,
				subject: "Incorrect destination for form",
				response: fmt.Sprintf(`
This message form is addressed to ICS Position %q at Location %q.  OA
Municipal Status messages should be addressed to either "Situation Analysis
Unit" or "Planning Section" at "County EOC".
`, msg.ToICSPosition, msg.ToLocation),
				references: refFormRouting,
			})
		} else if msg.ToICSPosition != "Situation Analysis Unit" && msg.ToICSPosition != "Planning Section" {
			a.problems = append(a.problems, &problem{
				code:    ProblemFormDestination,
				subject: "Incorrect destination for form",
				response: fmt.Sprintf(`
This message form is addressed to ICS Position %q.  OA Municipal Status
messages should be addressed to either "Situation Analysis Unit" or
"Planning Section".
`, msg.ToICSPosition),
				references: refFormRouting,
			})
		} else if msg.ToLocation != "County EOC" {
			a.problems = append(a.problems, &problem{
				code:    ProblemFormDestination,
				subject: "Incorrect destination for form",
				response: fmt.Sprintf(`
This message form is addressed to ICS Location %q.  OA Municipal Status
messages should be addressed to "County EOC".
`, msg.ToLocation),
				references: refFormRouting,
			})
		}
	case *pktmsg.RxSheltStatForm:
		// Note, we don't check the location; there are too many
		// valid possibilities.
		if msg.ToICSPosition != "Mass Care and Shelter Unit" && msg.ToICSPosition != "Care and Shelter Branch" && msg.ToICSPosition != "Operations Section" {
			a.problems = append(a.problems, &problem{
				code:    ProblemFormDestination,
				subject: "Incorrect destination for form",
				response: fmt.Sprintf(`
This message form is addressed to %q.  OA Shelter Status messages should be
addressed to "Mass Care and Shelter Unit", "Care and Shelter Branch", or
"Operations Section".
`, msg.ToICSPosition),
				references: refFormRouting,
			})
		}
	}
}
