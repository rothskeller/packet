package analyze

import (
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
)

func init() {
	Problems[ProbFormSubject.Code] = ProbFormSubject
	Problems[ProbHandlingOrderCode.Code] = ProbHandlingOrderCode
	Problems[ProbSubjectFormat.Code] = ProbSubjectFormat
	Problems[ProbSubjectHasSeverity.Code] = ProbSubjectHasSeverity
}

// ProbFormSubject is raised when the subject line of the message does not agree
// with what the embedded form would generate.
var ProbFormSubject = &Problem{
	Code:  "FormSubject",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbFormInvalid, ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		// Is the message of a type where the subject line can be
		// derived from the content (i.e., a known form type)?
		if a.xsc.Type.Tag != xscmsg.PlainTextTag {
			if a.xsc.Subject() != a.msg.Header.Get("Subject") {
				println(a.xsc.Subject(), a.msg.Header.Get("Subject"))
			}
			return a.xsc.Subject() != a.msg.Header.Get("Subject"), ""
		}
		return false, ""
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
	Code:  "HandlingOrderCode",
	ifnot: []*Problem{ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		xscsubj := xscmsg.ParseSubject(a.msg.Header.Get("Subject"))
		return xscsubj != nil && xscsubj.HandlingOrder == 0, ""
	},
	Variables: variableMap{
		"HANDLING": func(a *Analysis) string {
			return xscmsg.ParseSubject(a.msg.Header.Get("Subject")).HandlingOrderCode
		},
	},
}

// ProbSubjectFormat is raised when the subject line of the message does not
// have the proper format.
var ProbSubjectFormat = &Problem{
	Code:  "SubjectFormat",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbFormSubject, ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		xscsubj := xscmsg.ParseSubject(a.msg.Header.Get("Subject"))
		if xscsubj == nil {
			return true, "parse"
		}
		if a.xsc.Type.Tag == xscmsg.PlainTextTag && xscsubj.FormTag != "" {
			// Empirically, 90% of the time this happens because the subject
			// erroneously has an underscore after the word "Practice", and
			// so "Practice" got reported as the form name.
			if strings.EqualFold(xscsubj.FormTag, "Practice") {
				return true, "practice"
			}
			if !pktmsg.IsForm(a.msg.Body) {
				return true, "plainform"
			}
			// It's a corrupt form.  That gets reported
			// elsewhere, no need to pile on here.
		}
		return false, ""
	},
}

// ProbSubjectHasSeverity is raised when the subject line of the message
// contains a severity code.
var ProbSubjectHasSeverity = &Problem{
	Code:  "SubjectHasSeverity",
	ifnot: []*Problem{ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		xscsubj := xscmsg.ParseSubject(a.msg.Header.Get("Subject"))
		return xscsubj != nil && xscsubj.SeverityCode != "", ""
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
