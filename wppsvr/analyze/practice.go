package analyze

import (
	"regexp"
	"strings"
	"time"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/eoc213rr"
	"github.com/rothskeller/packet/xscmsg/ics213"
)

// Problem codes
const (
	ProblemPracticeSubjectFormat = "PracticeSubjectFormat"
)

func init() {
	ProblemLabel[ProblemPracticeSubjectFormat] = "incorrect Practice subject format"
}

// practiceRE matches a correctly formatted practice subject.  The subject must
// have the word Practice followed by four comma-separated fields (with
// whitespace also allowed between the fields).  The RE returns the first field
// (the call sign) and the fourth field (the date) as substrings so that they
// can be further checked.
var practiceRE = regexp.MustCompile(`(?i)^Practice\s+([AKNW][A-Z]?[0-9][A-Z]{1,3})\s*,[^,]+,([^,]+),\s*([^,]+?)\s*$`)

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

// checkPracticeSubject makes sure that the subject starts with the word
// Practice followed by appropriate data.
func (a *Analysis) checkPracticeSubject() {
	// Parse the subject line to get the true subject.
	xscsubj := xscmsg.ParseSubject(a.msg.Header.Get("Subject"))
	if xscsubj == nil {
		return // unparseable subject line reported elsewhere
	}
	// This check only applies to plain text messages, ICS-213 forms, and
	// EOC-213RR forms.  All other known form types have different subject
	// lines, and unknown form types get an error about that instead.
	switch a.xsc.(type) {
	case *config.PlainTextMessage:
		if pktmsg.IsForm(a.msg.Body) { // corrupt form
			return
		}
	case *ics213.Form, *eoc213rr.Form:
		break
	default:
		// Try to extract the jurisdiction from the entire subject.
		a.jurisdiction = jurisdictionMap[strings.ToUpper(xscsubj.Subject)]
		return
	}
	var match = practiceRE.FindStringSubmatch(xscsubj.Subject)
	if match == nil {
		a.problems = append(a.problems, &problem{
			code: ProblemPracticeSubjectFormat,
			response: `
The Subject of this message does not have the correct format.  After the
message number, handling order, and form type, it should have the word
"Practice" followed by four comma-separated fields:
    Practice CallSign, FirstName, Jurisdiction, Date
`,
			references: refWeeklyPractice,
		})
		return
	}
	a.subjectCallSign = strings.ToUpper(match[1])
	a.jurisdiction = jurisdictionMap[strings.ToUpper(strings.TrimSpace(match[2]))]
	// Look for a date.
	for fmt, re := range dateREs {
		if re.MatchString(match[3]) {
			if date, err := time.ParseInLocation(fmt, match[3], time.Local); err == nil {
				a.subjectDate = date
				break
			}
		}
	}
}

var jurisdictionMap = map[string]string{
	"ALAMEDA":              "XAL",
	"ALAMEDA CO.":          "XAL",
	"ALAMEDA COUNTY":       "XAL",
	"CAMPBELL":             "CBL",
	"CBL":                  "CBL",
	"CONTRA COSTA":         "XCC",
	"CONTRA COSTA CO.":     "XCC",
	"CONTRA COSTA COUNTY":  "XCC",
	"CUP":                  "CUP",
	"CUPERTINO":            "CUP",
	"GIL":                  "GIL",
	"GILROY":               "GIL",
	"LAH":                  "LAH",
	"LG/MS":                "LGT",
	"LGT":                  "LGT",
	"LMP":                  "LMP",
	"LOMA PRIETA":          "LMP",
	"LOS ALTOS HILLS":      "LAH",
	"LOS ALTOS":            "LOS",
	"LOS GATOS":            "LGT",
	"LOS":                  "LOS",
	"MARIN":                "XMR",
	"MARIN CO.":            "XMR",
	"MARIN COUNTY":         "XMR",
	"MILPITAS":             "MLP",
	"MLP":                  "MLP",
	"MONTE SERENO":         "MSO",
	"MONTEREY":             "XMY",
	"MONTEREY CO.":         "XMY",
	"MONTEREY COUNTY":      "XMY",
	"MORGAN HILL":          "MRG",
	"MOUNTAIN VIEW":        "MTV",
	"MRG":                  "MRG",
	"MSO":                  "MSO",
	"MTV":                  "MTV",
	"MV":                   "MTV",
	"NAM":                  "NAM",
	"NASA/AMES":            "NAM",
	"PAF":                  "PAF",
	"PALO ALTO":            "PAF",
	"SAN BENITO":           "XBE",
	"SAN BENITO CO.":       "XBE",
	"SAN BENITO COUNTY":    "XBE",
	"SAN FRANCISCO":        "XSF",
	"SAN FRANCISCO CO.":    "XSF",
	"SAN FRANCISCO COUNTY": "XSF",
	"SAN JOSE":             "SJC",
	"SAN MATEO":            "XSM",
	"SAN MATEO CO.":        "XSM",
	"SAN MATEO COUNTY":     "XSM",
	"SANTA CLARA":          "SNC",
	"SANTA CRUZ":           "XCZ",
	"SANTA CRUZ CO.":       "XCZ",
	"SANTA CRUZ COUNTY":    "XCZ",
	"SAR":                  "SAR",
	"SARATOGA":             "SAR",
	"SJC":                  "SJC",
	"SNC":                  "SNC",
	"SNY":                  "SNY",
	"STANFORD":             "STU",
	"STANFORD UNIVERSITY":  "STU",
	"STU":                  "STU",
	"SUNNYVALE":            "SNY",
	"XAL":                  "XAL",
	"XBE":                  "XBE",
	"XCC":                  "XCC",
	"XCZ":                  "XCZ",
	"XMR":                  "XMR",
	"XMY":                  "XMY",
	"XSC":                  "XSC",
	"XSF":                  "XSF",
	"XSM":                  "XSM",
}
