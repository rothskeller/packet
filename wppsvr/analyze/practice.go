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
	Problems[ProbFormPracticeSubject.Code] = ProbFormPracticeSubject
	Problems[ProbPracticeSubjectFormat.Code] = ProbPracticeSubjectFormat
	Problems[ProbUnknownJurisdiction.Code] = ProbUnknownJurisdiction
}

// PracticeSubject is the parsed contents of the "Practice ..." part of the
// subject line.
type PracticeSubject struct {
	CallSign     string
	Jurisdiction string
	NetDate      time.Time
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

// parsePracticeSubject parses the practice subject and returns the
// corresponding PracticeSubject structure, or nil if it couldn't be parsed
// successfully.  Jurisdictions are converted to their three-letter code if
// recognized; otherwise they are left as given.
func (a *Analysis) parsePracticeSubject() (ps *PracticeSubject) {
	// If we have an old Municipal Status form, the subject doesn't have the
	// full practice details; it only has the jurisdiction.
	if a.xsc.Type.Tag == jurisstat.Tag21 {
		ps = &PracticeSubject{Jurisdiction: a.subject.Subject}
	} else if match := practiceRE.FindStringSubmatch(a.subject.Subject); match != nil {
		ps = &PracticeSubject{
			CallSign:     strings.ToUpper(match[1]),
			Jurisdiction: strings.TrimSpace(match[2]),
		}
		ps.NetDate, _ = time.ParseInLocation("1/2/2006", match[3], time.Local)
	}
	if ps != nil {
		if abbr, ok := config.Get().Jurisdictions[strings.ToUpper(ps.Jurisdiction)]; ok {
			ps.Jurisdiction = abbr
		}
	}
	return ps
}

// ProbFormPracticeSubject is raised when the practice message details in the
// subject field of a form message have the wrong format.
var ProbFormPracticeSubject = &Problem{
	Code: "FormPracticeSubject",
	ifnot: []*Problem{
		// This check does not apply if the subject line of the message
		// doesn't match the subject field of the form.
		ProbFormSubject,
		// This check doesn't apply if the subject line of the message
		// is ill-formed.
		ProbSubjectFormat,
		// This check doesn't apply if the subject line of the message
		// has the word "Practice" as the form name.
		ProbPracticeAsFormName,
	},
	detect: func(a *Analysis) bool {
		// This check only applies to forms with a subject field.
		if a.xsc.Type.Tag == xscmsg.PlainTextTag || a.xsc.KeyField(xscmsg.FSubject) == nil {
			return false
		}
		// The check.
		return a.Practice == nil
	},
	Variables: variableMap{
		"SUBJECTFIELD": func(a *Analysis) string {
			return a.xsc.KeyField(xscmsg.FSubject).Def.Label
		},
	},
}

// ProbPracticeSubjectFormat is raised when the practice message details in the
// subject line of a non-form message have the wrong format.
var ProbPracticeSubjectFormat = &Problem{
	Code: "PracticeSubjectFormat",
	ifnot: []*Problem{
		// This check does not apply if the message appears to contain
		// a form that couldn't be parsed.
		ProbFormCorrupt,
		// This check does not apply if the subject line of the message
		// is ill-formed.
		ProbSubjectFormat,
		// This check does not apply if the subject line of the message
		// has the word "Practice" as the form name.
		ProbPracticeAsFormName,
	},
	detect: func(a *Analysis) bool {
		// This check only applies to plain text messages.
		if a.xsc.Type.Tag != xscmsg.PlainTextTag {
			return false
		}
		// The check.
		return a.Practice == nil
	},
}

// ProbUnknownJurisdiction is raised when the provided jurisdiction is not one
// of the recognized ones.
var ProbUnknownJurisdiction = &Problem{
	Code: "UnknownJurisdiction",
	ifnot: []*Problem{
		// This check does not apply if the message appears to contain a
		// form that couldn't be parsed.
		ProbFormCorrupt,
		// This check does not apply if the subject line of the message
		// doesn't match the subject field of the form.
		ProbFormSubject,
		// This check does not apply if the subject line of the message
		// is ill-formed.
		ProbSubjectFormat,
		// This check does not apply if the subject line of the message
		// has the word "Practice" as the form name.
		ProbPracticeAsFormName,
		// This check does not apply if the practice details are
		// ill-formed.
		ProbFormPracticeSubject,
		ProbPracticeSubjectFormat,
	},
	detect: func(a *Analysis) bool {
		return config.Get().Jurisdictions[a.Practice.Jurisdiction] == ""
	},
	Variables: variableMap{
		"JURISDICTION": func(a *Analysis) string { return a.Practice.Jurisdiction },
	},
}
