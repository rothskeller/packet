package analyze

import (
	"fmt"
	"strings"

	"steve.rothskeller.net/packet/wppsvr/config"
	"steve.rothskeller.net/packet/xscmsg"
)

// Problem codes
const (
	ProblemFormSubject        = "FormSubject"
	ProblemHandlingOrderCode  = "HandlingOrderCode"
	ProblemSubjectFormat      = "SubjectFormat"
	ProblemSubjectHasSeverity = "SubjectHasSeverity"
)

func init() {
	ProblemLabel[ProblemFormSubject] = "message subject doesn't agree with form contents"
	ProblemLabel[ProblemHandlingOrderCode] = "unknown handling order code"
	ProblemLabel[ProblemSubjectFormat] = "incorrect subject line format"
	ProblemLabel[ProblemSubjectHasSeverity] = "severity on subject line"
}

// checkSubjectLine checks for problems with the Subject line of the message.
// (Specifically, the SCCo-standard parts of it, not the weekly practice parts
// of it.)
func (a *Analysis) checkSubjectLine() {
	// Is the message of a type where the subject line can be derived from
	// the content (i.e., a known form type)?
	if es, ok := a.xsc.(interface{ EncodeSubject() string }); ok {
		// Yes.  But we only want to do so if the form validates.  If it
		// doesn't, the validation errors are enough and subject errors
		// would be redundant.
		if xsc, ok := a.xsc.(interface{ Validate(bool) []string }); ok {
			if problems := xsc.Validate(true); len(problems) != 0 {
				return
			}
		}
		// What should the subject line be, and what is it?
		want := es.EncodeSubject()
		have := a.msg.Header.Get("Subject")
		if have != want {
			a.problems = append(a.problems, &problem{
				code: ProblemFormSubject,
				response: fmt.Sprintf(`
This message has
	Subject: %s
but, based on the contents of the form, it should have
	Subject: %s
PackItForms automatically generates the Subject line from the form contents; it
should not be overridden manually.
`, have, want),
			})
		}
		return
	}
	// Parse the subject line.  If we're not able to parse it, report that
	// problem.
	xscsubj := xscmsg.ParseSubject(a.msg.Header.Get("Subject"))
	if xscsubj == nil {
		a.problems = append(a.problems, problemSubjectFormat(""))
		return
	}
	// The message number is checked elsewhere.
	// Check for a severity code.  There shouldn't be one.
	if xscsubj.SeverityCode != "" {
		a.problems = append(a.problems, &problem{
			code: ProblemSubjectHasSeverity,
			response: fmt.Sprintf(`
The Subject line of this message contains a both a Severity code and a
Handling Order code ("_%s/%s_").  This is an outdated Subject line style.
Current SCCo standards include only the Handling Order code on the Subject
line ("_%s_").
`, xscsubj.SeverityCode, xscsubj.HandlingOrderCode, xscsubj.HandlingOrderCode),
			references: refSubjectLine,
		})
	}
	// Make sure the handling order code is valid.
	if xscsubj.HandlingOrder == 0 {
		a.problems = append(a.problems, &problem{
			code: ProblemHandlingOrderCode,
			response: fmt.Sprintf(`
The Subject line of this message contains an invalid Handling Order code (%s).
The valid codes are "I" for Immediate, "P" for Priority, and "R" for Routine.
`, xscsubj.HandlingOrderCode),
			references: refSubjectLine,
		})
	}
	// If this is a plain text message, there shouldn't be a form name in
	// the subject line.
	if _, ok := a.xsc.(*config.PlainTextMessage); ok {
		if xscsubj.FormTag != "" {
			// Empirically, 90% of the time this happens because the subject
			// erroneously has an underscore after the word "Practice", and
			// so "Practice" got reported as the form name.
			if strings.EqualFold(xscsubj.FormTag, "Practice") {
				a.problems = append(a.problems, problemSubjectFormat(`
Note that there is no underline after the word "Practice".`))
			} else {
				a.problems = append(a.problems, problemSubjectFormat(`
Note that there is no form name between the handling order and the word
"Practice" for plain text messages.  If this is in fact a form message, it
is improperly encoded and was not recognized as a form.`))
			}
		}
	}
}

func problemSubjectFormat(note string) *problem {
	return &problem{
		code: ProblemSubjectFormat,
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
