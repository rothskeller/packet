package analyze

import (
	"strings"

	"github.com/rothskeller/packet/xscmsg"
)

func init() {
	Problems[ProbFormSubject.Code] = ProbFormSubject
	Problems[ProbHandlingOrderCode.Code] = ProbHandlingOrderCode
	Problems[ProbPracticeAsFormName.Code] = ProbPracticeAsFormName
	Problems[ProbSubjectFormat.Code] = ProbSubjectFormat
	Problems[ProbSubjectHasSeverity.Code] = ProbSubjectHasSeverity
	Problems[ProbSubjectPlainForm.Code] = ProbSubjectPlainForm
}

// ProbFormSubject is raised when the subject line of the message does not agree
// with what the embedded form would generate.
var ProbFormSubject = &Problem{
	Code: "FormSubject",
	ifnot: []*Problem{
		// This check does not apply if the form contents don't pass
		// validation checks (such as, for example, not having anything
		// in the subject field at all).
		ProbFormInvalid,
	},
	detect: func(a *Analysis) bool {
		// This check only applies to forms.
		if a.xsc.Type.Tag == xscmsg.PlainTextTag {
			return false
		}
		// The check.
		return a.xsc.Subject() != a.msg.Header.Get("Subject")
	},
	Variables: variableMap{
		"ACTUALSUBJ": func(a *Analysis) string {
			return a.msg.Header.Get("Subject")
		},
		"EXPECTSUBJ": func(a *Analysis) string {
			return a.xsc.Subject()
		},
	},
}

// ProbHandlingOrderCode is raised when the subject line of the message contains
// an unknown handling order code.
var ProbHandlingOrderCode = &Problem{
	Code: "HandlingOrderCode",
	detect: func(a *Analysis) bool {
		// This check does not apply if the subject line could not be
		// parsed.
		if a.subject == nil {
			return false
		}
		// The check.
		return a.subject.HandlingOrder == 0
	},
	Variables: variableMap{
		"HANDLING": func(a *Analysis) string {
			return xscmsg.ParseSubject(a.msg.Header.Get("Subject")).HandlingOrderCode
		},
	},
}

// ProbSubjectFormat is raised when the subject line of the message could not be
// parsed.
var ProbSubjectFormat = &Problem{
	Code:  "SubjectFormat",
	ifnot: []*Problem{ProbFormSubject},
	detect: func(a *Analysis) bool {
		return a.subject == nil
	},
}

// ProbPracticeAsFormName is raised when the subject line of a plain text
// message has a form name of "Practice".  This means the sender put an
// underline after the word "Practice" rather than a blank.
var ProbPracticeAsFormName = &Problem{
	Code: "PracticeAsFormName",
	detect: func(a *Analysis) bool {
		// This check only applies to plain text messages.
		if a.xsc.Type.Tag != xscmsg.PlainTextTag {
			return false
		}
		// This check does not apply if the subject line could not be
		// parsed.
		if a.subject == nil {
			return false
		}
		// The check.
		return strings.EqualFold(a.subject.FormTag, "Practice")
	},
}

// ProbSubjectPlainForm is raised when the subject line of a plain text message
// has a form name in it (other than "Practice").
var ProbSubjectPlainForm = &Problem{
	Code: "SubjectPlainForm",
	ifnot: []*Problem{
		// This check does not apply if the form name is "Practice".
		ProbPracticeAsFormName,
		// This check does not apply if the message appears to contain
		// a form that couldn't be parsed.
		ProbFormCorrupt,
	},
	detect: func(a *Analysis) bool {
		// This check only applies to plain text messages.
		if a.xsc.Type.Tag != xscmsg.PlainTextTag {
			return false
		}
		// This check does not apply if the subject line could not be
		// parsed.
		if a.subject == nil {
			return false
		}
		// The check.
		return a.subject.FormTag != ""
	},
}

// ProbSubjectHasSeverity is raised when the subject line of the message
// contains a severity code.
var ProbSubjectHasSeverity = &Problem{
	Code: "SubjectHasSeverity",
	detect: func(a *Analysis) bool {
		// This check does not apply if the subject line could not be
		// parsed.
		if a.subject == nil {
			return false
		}
		// The check.
		return a.subject.SeverityCode != ""
	},
	Variables: variableMap{
		"HANDLING": func(a *Analysis) string {
			return xscmsg.ParseSubject(a.msg.Header.Get("Subject")).HandlingOrderCode
		},
		"SEVERITY": func(a *Analysis) string {
			return xscmsg.ParseSubject(a.msg.Header.Get("Subject")).SeverityCode
		},
	},
}
