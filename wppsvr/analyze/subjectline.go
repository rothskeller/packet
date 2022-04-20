package analyze

import (
	"fmt"
	"strings"
)

// Problem codes
const (
	ProblemHandlingOrderCode  = "HandlingOrderCode"
	ProblemSubjectFormat      = "SubjectFormat"
	ProblemSubjectHasSeverity = "SubjectHasSeverity"
)

func init() {
	ProblemLabel[ProblemHandlingOrderCode] = "unknown handling order code"
	ProblemLabel[ProblemSubjectFormat] = "incorrect subject line format"
	ProblemLabel[ProblemSubjectHasSeverity] = "severity on subject line"
}

// checkSubjectLine checks for problems with the Subject line of the message.
// (Specifically, the SCCo-standard parts of it, not the weekly practice parts
// of it.)
func (a *Analysis) checkSubjectLine() {
	// These checks only apply to human messages.
	var msg = a.msg.Message()
	if msg == nil {
		return
	}
	// These checks do not apply to forms for which we are able to generate
	// the subject line from the form contents.  Those subject lines are
	// checked separately.
	if _, ok := a.msg.(interface{ EncodeSubjectLine() string }); ok {
		return
	}
	// First, check the overall format of the subject line.  If pktmsg
	// wasn't able to parse it, it will have put the entire SubjectLine into
	// the Subject field.  Check for that.
	if msg.Subject == msg.SubjectLine {
		a.problems = append(a.problems, problemSubjectFormat(""))
		return
	}
	// There's no need to check the message number here.  If this is a plain
	// text message or an unknown form type, it is checked by
	// checkMessageNumber.  If it's a known form type, checkMessageNumber
	// checks the number in the form, and checkFormSubject makes sure the
	// one in the subject line matches.

	// Next, check for a severity code.  There shouldn't be one.
	if msg.SeverityCode != "" {
		a.problems = append(a.problems, &problem{
			code:    ProblemSubjectHasSeverity,
			subject: "Severity on Subject line",
			response: fmt.Sprintf(`
The Subject line of this message contains a both a Severity code and a
Handling Order code ("_%s/%s_").  This is an outdated Subject line style.
Current SCCo standards include only the Handling Order code on the Subject
line ("_%s_").
`, msg.SeverityCode, msg.HandlingOrderCode, msg.HandlingOrderCode),
			references: refSubjectLine,
		})
	}
	// Next, make sure the handling order code is valid.
	if msg.HandlingOrder == 0 {
		a.problems = append(a.problems, &problem{
			code:    ProblemHandlingOrderCode,
			subject: "Unknown handling order code",
			response: fmt.Sprintf(`
The Subject line of this message contains an invalid Handling Order code (%s).
The valid codes are "I" for Immediate, "P" for Priority, and "R" for Routine.
`, msg.HandlingOrderCode),
			references: refSubjectLine,
		})
	}
	// Next, check the form name.  If we have a known form and the form name
	// is wrong, checkFormSubject will catch that.  If we have an unknown
	// form, checkCorrectForm will catch that.  The only case we need to
	// check here is a plain text message that has a form name in the
	// subject.  We'll report it as a subject format problem, with an
	// additional note.
	if a.msg.Form() == nil && msg.FormName != "" {
		// Empirically, 90% of the time this happens because the subject
		// erroneously has an underscore after the word "Practice", and
		// so "Practice" got reported as the form name.
		if strings.EqualFold(msg.FormName, "Practice") {
			a.problems = append(a.problems, problemSubjectFormat(`
Note that there is no underline after the word "Practice".`))
		} else {
			a.problems = append(a.problems, problemSubjectFormat(`
Note that there is no form name between the handling order and the word
"Practice" for plain text messages.`))
		}
	}
}

func problemSubjectFormat(note string) *problem {
	return &problem{
		code:    ProblemSubjectFormat,
		subject: "Incorrect Subject line format",
		response: fmt.Sprintf(`
This message has an incorrect subject line format.  According to SCCo
standards, the subject should look like
    AAA-111P_R_Practice A6AAA, Charlie, Nowhereville, 2019-03-16
    (msgnum) |            |    (name)   (city)        (net date)
      (handling order)  (call sign)%s
`, note),
		references: refWeeklyPractice | refSubjectLine,
	}
}
