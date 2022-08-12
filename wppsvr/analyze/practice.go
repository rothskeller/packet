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

// ProbPracticeSubjectFormat is raised when the practice message details in the
// subject line have the wrong format.
var ProbPracticeSubjectFormat = &Problem{
	Code:  "PracticeSubjectFormat",
	Label: "incorrect practice message details",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbFormSubject, ProbSubjectFormat, ProbFormCorrupt, ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		subject := xscmsg.ParseSubject(a.msg.Header.Get("Subject")).Subject
		// Don't check Jurisdiction Status forms; their subject doesn't
		// have the full practice details.  It just has the jurisdiction, which we save.
		if _, ok := a.xsc.(*jurisstat.Form); ok {
			a.jurisdiction = jurisdictionMap[strings.ToUpper(subject)]
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
