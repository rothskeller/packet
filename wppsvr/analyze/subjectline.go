package analyze

import (
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/xscmsg"
)

func init() {
	Problems[ProbFormSubject.Code] = ProbFormSubject
	Problems[ProbHandlingOrderCode.Code] = ProbHandlingOrderCode
	Problems[ProbSubjectFormat.Code] = ProbSubjectFormat
	Problems[ProbSubjectHasSeverity.Code] = ProbSubjectHasSeverity
}

type formWithEncodableSubject interface {
	EncodeSubject() string
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
		if es, ok := a.xsc.(formWithEncodableSubject); ok {
			return es.EncodeSubject() != a.msg.Header.Get("Subject"), ""
		}
		return false, ""
	},
	Variables: variableMap{
		"ACTUALSUBJ": func(a *Analysis) string {
			return a.msg.Header.Get("Subject")
		},
		"EXPECTSUBJ": func(a *Analysis) string {
			return a.xsc.(formWithEncodableSubject).EncodeSubject()
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
	references: refSubjectLine,
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
		if _, ok := a.xsc.(*config.PlainTextMessage); ok && xscsubj.FormTag != "" {
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
	references: refWeeklyPractice | refSubjectLine,
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
	references: refSubjectLine,
}
