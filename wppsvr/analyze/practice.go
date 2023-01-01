package analyze

import (
	"regexp"
	"strings"
	"time"

	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/jurisstat"
)

func init() {
	Problems[ProbPracticeSubjectFormat.Code] = ProbPracticeSubjectFormat
	Problems[ProbUnknownJurisdiction.Code] = ProbUnknownJurisdiction
}

// practiceRE matches a correctly formatted practice subject.  The subject must
// have the word Practice followed by four comma-separated fields (with
// whitespace also allowed between the fields).  The RE returns the first field
// (the call sign), the third field (the jurisdiction), and the fourth field
// (the date) as substrings so that they can be further checked and stored.
// A comma is allowed after the word "Practice", which doesn't exactly conform
// to the required syntax, but it is a very common mistake and not worth
// penalizing.
var practiceRE = regexp.MustCompile(`(?i)^Practice[,\s]+([AKNW][A-Z]?[0-9][A-Z]{1,3}|[A-Z][A-Z0-9]{5})\s*,[^,]+,([^,]+),\s*((?:0?[1-9]|1[0-2])/(?:0?[1-9]|[12]\d|3[01])/20\d\d)\s*$`)

// ProbPracticeSubjectFormat is raised when the practice message details in the
// subject line have the wrong format.
var ProbPracticeSubjectFormat = &Problem{
	Code:  "PracticeSubjectFormat",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbFormSubject, ProbSubjectFormat, ProbFormCorrupt, ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		jurisdictionMap := config.Get().Jurisdictions
		subject := xscmsg.ParseSubject(a.msg.Header.Get("Subject")).Subject
		// Don't check Municipal Status forms; their subject doesn't
		// have the full practice details.  It just has the
		// jurisdiction, which we save.
		if a.xsc.Type.Tag == jurisstat.Tag21 {
			a.jurisdiction = subject
			if abbr, ok := jurisdictionMap[strings.ToUpper(a.jurisdiction)]; ok {
				a.jurisdiction = abbr
			}
			return false, ""
		}
		// Parse the subject line to get the true subject.
		xscsubj := xscmsg.ParseSubject(a.msg.Header.Get("Subject"))
		if xscsubj == nil {
			return false, "" // unparseable subject line reported elsewhere
		}
		// Do we have a valid practice message subject?
		var match = practiceRE.FindStringSubmatch(subject)
		if match == nil {
			if a.xsc.Type.Tag == xscmsg.PlainTextTag {
				return true, "plain"
			}
			if a.xsc.Field(xscmsg.FSubject) != nil {
				return true, "form"
			}
			// It's a form that can't be used for weekly practice
			// (e.g. check-in form).  We'll ignore the bad subject
			// and just give them the "wrong form" message.
			return false, ""
		}
		// Yes, we do, so save the information from it for other
		// analysis steps.
		a.subjectCallSign = strings.ToUpper(match[1])
		a.jurisdiction = strings.TrimSpace(match[2])
		if abbr, ok := jurisdictionMap[strings.ToUpper(a.jurisdiction)]; ok {
			a.jurisdiction = abbr
		}
		a.subjectDate, _ = time.ParseInLocation("1/2/2006", match[3], time.Local)
		return false, ""
	},
	Variables: variableMap{
		"SUBJECTFIELD": func(a *Analysis) string {
			return a.xsc.Field(xscmsg.FSubject).Def.Tag
		},
	},
}

// ProbUnknownJurisdiction is raised when the provided jurisdiction is not one
// of the recognized ones.
var ProbUnknownJurisdiction = &Problem{
	Code:  "UnknownJurisdiction",
	after: []*Problem{ProbPracticeSubjectFormat}, // sets a.jurisdiction
	ifnot: []*Problem{ProbPracticeSubjectFormat, ProbFormSubject, ProbSubjectFormat, ProbFormCorrupt, ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		_, ok := config.Get().Jurisdictions[a.jurisdiction]
		return !ok, ""
	},
	Variables: variableMap{
		"JURISDICTION": func(a *Analysis) string { return a.jurisdiction },
	},
}
