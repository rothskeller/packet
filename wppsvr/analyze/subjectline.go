package analyze

import (
	"strings"

	"github.com/rothskeller/packet/pktmsg"
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
	Code:  "FormSubject",
	ifnot: []*Problem{ProbFormInvalid},
	detect: func(a *Analysis) bool {
		// Is the message of a type where the subject line can be
		// derived from the content (i.e., a known form type)?
		if a.xsc.Type.Tag != xscmsg.PlainTextTag {
			return a.xsc.Subject() != a.msg.Header.Get("Subject")
		}
		return false
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
		xscsubj := xscmsg.ParseSubject(a.msg.Header.Get("Subject"))
		return xscsubj != nil && xscsubj.HandlingOrder == 0
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
	Code:  "PracticeAsFormName",
	ifnot: []*Problem{ProbFormSubject},
	detect: func(a *Analysis) bool {
		return a.xsc.Type.Tag == xscmsg.PlainTextTag && a.subject != nil && strings.EqualFold(a.subject.FormTag, "Practice")
	},
}

// ProbSubjectPlainForm is raised when the subject line of a plain text message
// has a form name in it (other than "Practice").
var ProbSubjectPlainForm = &Problem{
	Code:  "SubjectPlainForm",
	ifnot: []*Problem{ProbPracticeAsFormName, ProbFormSubject},
	detect: func(a *Analysis) bool {
		return a.xsc.Type.Tag == xscmsg.PlainTextTag && a.subject != nil && a.subject.FormTag != "" && !pktmsg.IsForm(a.msg.Body)
	},
}

// ProbSubjectHasSeverity is raised when the subject line of the message
// contains a severity code.
var ProbSubjectHasSeverity = &Problem{
	Code: "SubjectHasSeverity",
	detect: func(a *Analysis) bool {
		xscsubj := xscmsg.ParseSubject(a.msg.Header.Get("Subject"))
		return xscsubj != nil && xscsubj.SeverityCode != ""
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
