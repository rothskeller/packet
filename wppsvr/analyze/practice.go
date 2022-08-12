package analyze

import (
	"regexp"
	"strings"
	"time"

	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/ahfacstat"
	"github.com/rothskeller/packet/xscmsg/eoc213rr"
	"github.com/rothskeller/packet/xscmsg/ics213"
	"github.com/rothskeller/packet/xscmsg/jurisstat"
	"github.com/rothskeller/packet/xscmsg/racesmar"
	"github.com/rothskeller/packet/xscmsg/sheltstat"
)

func init() {
	Problems[ProbPracticeSubjectFormat.Code] = ProbPracticeSubjectFormat
	Problems[ProbUnknownJurisdiction.Code] = ProbUnknownJurisdiction
}

// practiceRE matches a correctly formatted practice subject.  The subject must
// have the word Practice followed by four comma-separated fields (with
// whitespace also allowed between the fields).  The RE returns the first field
// (the call sign) and the fourth field (the date) as substrings so that they
// can be further checked.
// A comma is allowed after the word "Practice", which doesn't exactly conform
// to the required syntax, but it is a very common mistake and not worth
// penalizing.
var practiceRE = regexp.MustCompile(`(?i)^Practice[,\s]+([AKNW][A-Z]?[0-9][A-Z]{1,3})\s*,[^,]+,([^,]+),\s*([^,]+?)\s*$`)

// dateREs is a map from date format string to regular expression.  There is no
// standard for the formatting of the date on the subject line, but we try all
// of the common ones.
var dateREs = map[string]*regexp.Regexp{
	// These are all patterns seen in actual check-in data.  The point of
	// parsing this is to verify that the date correctly identifies the
	// check-in period, not to verify the syntax of the date, so we accept
	// as many patterns as we can reasonably parse.
	"1-2-2006":        regexp.MustCompile(`^\d?\d-\d?\d-20\d\d$`),
	"1/2/2006":        regexp.MustCompile(`^\d?\d/\d?\d/20\d\d$`),
	"2 Jan 2006":      regexp.MustCompile(`(?i)^\d?\d (?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec) 20\d\d$`),
	"2-Jan-2006":      regexp.MustCompile(`(?i)^\d?\d-(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)-20\d\d$`),
	"2006-1-2":        regexp.MustCompile(`^20\d\d-\d?\d-\d?\d$`),
	"2006-Jan-2":      regexp.MustCompile(`(?i)^20\d\d-(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)-\d?\d$`),
	"2006/1/2":        regexp.MustCompile(`^20\d\d/\d?\d/\d?\d$`),
	"20060102":        regexp.MustCompile(`^20\d\d\d\d\d\d$`),
	"2Jan2006":        regexp.MustCompile(`(?i)^\d?\d(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)20\d\d$`),
	"Jan 2, 2006":     regexp.MustCompile(`(?i)^(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec) \d?\d, 20\d\d$`),
	"Jan-2-2006":      regexp.MustCompile(`(?i)^(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)-\d?\d-20\d\d$`),
	"Jan. 2, 2006":    regexp.MustCompile(`(?i)^(?:jan|feb|mar|apr|jun|jul|aug|sep|oct|nov|dec)\. \d?\d, 20\d\d$`),
	"January 2, 2006": regexp.MustCompile(`(?i)^(?:january|february|march|april|may|june|july|august|september|october|november|december) \d?\d, 20\d\d$`),
	// Repeat all of the above with two-digit years, except where that
	// doesn't make sense.
	"1-2-06":   regexp.MustCompile(`^\d?\d-\d?\d-\d\d$`),
	"1/2/06":   regexp.MustCompile(`^\d?\d/\d?\d/\d\d$`),
	"2 Jan 06": regexp.MustCompile(`(?i)^\d?\d (?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec) \d\d$`),
	"2-Jan-06": regexp.MustCompile(`(?i)^\d?\d-(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)-\d\d$`),
	"060102":   regexp.MustCompile(`^\d\d\d\d\d\d$`),
	"2Jan06":   regexp.MustCompile(`(?i)^\d?\d(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)\d\d$`),
	"Jan-2-06": regexp.MustCompile(`(?i)^(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)-\d?\d-\d\d$`),
}

// ProbPracticeSubjectFormat is raised when the practice message details in the
// subject line have the wrong format.
var ProbPracticeSubjectFormat = &Problem{
	Code:  "PracticeSubjectFormat",
	Label: "incorrect practice message details",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbFormSubject, ProbSubjectFormat, ProbFormCorrupt, ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		jurisdictionMap := config.Get().Jurisdictions
		subject := xscmsg.ParseSubject(a.msg.Header.Get("Subject")).Subject
		// Don't check Jurisdiction Status forms; their subject doesn't
		// have the full practice details.  It just has the
		// jurisdiction, which we save.
		if _, ok := a.xsc.(*jurisstat.Form); ok {
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
			if _, ok := a.xsc.(*config.PlainTextMessage); ok {
				return true, "plain"
			}
			return true, "form"
		}
		// Yes, we do, so save the information from it for other
		// analysis steps.
		a.subjectCallSign = strings.ToUpper(match[1])
		a.jurisdiction = strings.TrimSpace(match[2])
		if abbr, ok := jurisdictionMap[strings.ToUpper(a.jurisdiction)]; ok {
			a.jurisdiction = abbr
		}
		// Look for a date.
		for fmt, re := range dateREs {
			if re.MatchString(match[3]) {
				if date, err := time.ParseInLocation(fmt, match[3], time.Local); err == nil {
					a.subjectDate = date
					break
				}
			}
		}
		return false, ""
	},
	Variables: variableMap{
		"SUBJECTFIELD": func(a *Analysis) string {
			switch a.xsc.(type) {
			case *ics213.Form:
				return "Subject"
			case *eoc213rr.Form:
				return "Incident Name"
			case *sheltstat.Form:
				return "Shelter Name"
			case *ahfacstat.Form:
				return "Facility Name"
			case *racesmar.Form:
				return "Agency Name"
			default:
				panic("unknown field name for subject")
			}
		},
	},
	references: refWeeklyPractice,
}

// ProbUnknownJurisdiction is raised when the provided jurisdiction is not one
// of the recognized ones.
var ProbUnknownJurisdiction = &Problem{
	Code:  "UnknownJurisdiction",
	Label: "unknown jurisdiction",
	after: []*Problem{ProbPracticeSubjectFormat}, // sets a.jurisdiction
	ifnot: []*Problem{ProbPracticeSubjectFormat, ProbFormSubject, ProbSubjectFormat, ProbFormCorrupt, ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		_, ok := config.Get().Jurisdictions[a.jurisdiction]
		return !ok, ""
	},
	Variables: variableMap{
		"JURISDICTION": func(a *Analysis) string { return a.jurisdiction },
	},
	references: refWeeklyPractice,
}
