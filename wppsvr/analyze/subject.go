package analyze

import (
	"regexp"
	"strings"
	"time"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/jurisstat"
)

func init() {
	ProblemLabels["FormPracticeSubject"] = "incorrect practice message details in form"
	ProblemLabels["FormSubject"] = "message subject doesn't agree with form contents"
	ProblemLabels["HandlingOrderCode"] = "unknown handling order code"
	ProblemLabels["MsgNumFormat"] = "incorrect message number format"
	ProblemLabels["MsgNumPrefix"] = "incorrect message number prefix"
	ProblemLabels["PracticeAsFormName"] = `incorrect subject line format (underline after "Practice")`
	ProblemLabels["PracticeSubjectFormat"] = "incorrect practice message details"
	ProblemLabels["SubjectFormat"] = "incorrect subject line format"
	ProblemLabels["SubjectHasSeverity"] = "severity on subject line"
	ProblemLabels["SubjectPlainForm"] = "form name in subject of non-form message"
	ProblemLabels["UnknownJurisdiction"] = "unknown jurisdiction"
}

// PracticeSubject is the parsed contents of the "Practice ..." part of the
// subject line.
type PracticeSubject struct {
	CallSign     string
	Jurisdiction string
	NetDate      time.Time
}

// checkSubject checks for problems with the subject line.  As a side effect, it
// sets a.subject and a.Practice.
func (a *Analysis) checkSubject() {
	a.checkFormSubject()
	a.checkSubjectFormat() // sets a.subject
	a.checkPracticeInfo()  // sets a.Practice
}

// checkFormSubject checks that the computed Subject line based on the form
// contents (if any) matches the actual Subject line of the message.
func (a *Analysis) checkFormSubject() {
	// Note that, for messages that can't compute a Subject line,
	// a.xsc.Subject() will simply return the actual Subject line, so this
	// check will always pass for them.
	act := a.msg.Header.Get("Subject")
	exp := a.xsc.Subject()
	if exp != act {
		// JNOS is known to trim spaces from the end of subject lines.
		// That's not a human error, so it should be accepted.
		exp = strings.TrimRight(exp, " ")
		if exp != act {
			a.reportProblem("FormSubject", 0, formSubjectResponse, act, exp)
		}
	}
}

var msgnumRE = regexp.MustCompile(`^(?:[A-Z][A-Z][A-Z]|[A-Z][0-9][A-Z0-9]|[0-9][A-Z][A-Z])-\d\d\d+[PMR]$`)
var fccCallSignRE = regexp.MustCompile(`^(?:A[A-L]|[KNW][A-Z]?)[0-9][A-Z]{1,3}$`)

// checkSubjectFormat makes sure the Subject line of the message follows the
// expected standard for SCCo packet messages.  As a side effect, it sets
// a.subject.  It does not check the "Practice ..." portion of the Subject line.
func (a *Analysis) checkSubjectFormat() {
	// Parse the subject line.
	if a.subject = xscmsg.ParseSubject(a.msg.Header.Get("Subject")); a.subject == nil {
		a.reportProblem("SubjectFormat", refSubjectLine|refWeeklyPractice, subjectFormatResponse)
		// Couldn't parse it, so none of the other checks in this
		// routine apply.
		return
	}
	// Check the message number format.
	if !msgnumRE.MatchString(a.subject.MessageNumber) {
		a.reportProblem("MsgNumFormat", refOutpostConfig|refSubjectLine, msgNumFormatResponse)
	} else {
		// It's a valid message number.  If the message is coming from
		// an FCC call sign, the prefix of the message number should be
		// the last three characters of the call sign.
		if fccCallSignRE.MatchString(a.FromCallSign) {
			act := a.subject.MessageNumber[:3]
			exp := a.FromCallSign[len(a.FromCallSign)-3:]
			if act != exp {
				a.reportProblem("MsgNumPrefix", refSubjectLine, msgNumPrefixResponse, act, exp)
			}
		}
	}
	// Check the handling order code.
	if a.subject.HandlingOrder == 0 {
		a.reportProblem("HandlingOrderCode", refSubjectLine, handlingOrderCodeResponse, a.subject.HandlingOrderCode)
	}
	// Check for an obsolete severity code.
	if a.subject.SeverityCode != "" {
		a.reportProblem("SubjectHasSeverity", refSubjectLine, subjectHasSeverityResponse,
			a.subject.SeverityCode, a.subject.HandlingOrderCode)
	}
	// Check for a form name where there shouldn't be one.
	if a.xsc.Type.Tag == xscmsg.PlainTextTag && a.subject.FormTag != "" {
		if strings.EqualFold(a.subject.FormTag, "Practice") {
			a.reportProblem("PracticeAsFormName", refSubjectLine|refWeeklyPractice, practiceAsFormNameResponse)
		} else if !pktmsg.IsForm(a.msg.Body) {
			a.reportProblem("SubjectPlainForm", refSubjectLine, subjectPlainFormResponse)
		} // else handled by FormCorrupt in body.go
	}
}

// checkPracticeInfo checks the "Practice ..." portion of the Subject line.
func (a *Analysis) checkPracticeInfo() {
	if mtc := config.Get().MessageTypes[a.xsc.Type.Tag]; mtc == nil || mtc.NoPracticeInfo {
		// This is an unknown form type or one that doesn't support
		// Practice... info on the Subject line.  Return without raising
		// an error.
		return
	}
	if f := a.xsc.KeyField(xscmsg.FSubject); f != nil && a.xsc.Type.Tag != xscmsg.PlainTextTag {
		// This is a form with a Subject field.  Check the practice info
		// from there.
		if a.Practice = a.parsePracticeSubject(f.Value); a.Practice == nil {
			a.reportProblem("FormPracticeSubject", refWeeklyPractice, formPracticeSubjectResponse, f.Def.Label)
			return
		}
	} else if a.subject != nil {
		// This is not a form with a Subject field, but we were able to
		// parse a subject out of the Subject line of the message.
		// Check the practice info from there.
		subject := a.subject.Subject
		if _, ok := a.problems["PracticeAsFormName"]; ok {
			subject = "Practice " + subject
		}
		if a.Practice = a.parsePracticeSubject(subject); a.Practice == nil {
			a.reportProblem("PracticeSubjectFormat", refWeeklyPractice, practiceSubjectFormatResponse)
			return
		}
	} else {
		return // no subject to parse
	}
	// We were able to parse the practice subject.  Do we have a valid
	// jurisdiction?
	if config.Get().Jurisdictions[a.Practice.Jurisdiction] == "" {
		a.reportProblem("UnknownJurisdiction", refBBSList|refWeeklyPractice, unknownJurisdictionResponse,
			a.Practice.Jurisdiction)
	}
}

// practiceRE matches a correctly formatted practice subject.  The subject must
// have the word Practice followed by four comma-separated fields (with
// whitespace also allowed between the fields).  The RE returns the first field
// (the call sign), the third field (the jurisdiction), and the fourth field
// (the date) as substrings so that they can be further checked and stored.
// A comma is allowed after the word "Practice", which doesn't exactly conform
// to the required syntax, but it is a very common mistake and not worth
// penalizing.
var practiceRE = regexp.MustCompile(`(?i)^Practice[,\s]+((?:A[A-L]|[KNW][A-Z]?)[0-9][A-Z]{1,3}|[A-Z][A-Z0-9]{5})\s*,[^,]+,([^,]+),\s*((?:0?[1-9]|1[0-2])/(?:0?[1-9]|[12]\d|3[01])/20\d\d)\s*$`)

// parsePracticeSubject parses the practice subject and returns the
// corresponding PracticeSubject structure, or nil if it couldn't be parsed
// successfully.  Jurisdictions are converted to their three-letter code if
// recognized; otherwise they are left as given.
func (a *Analysis) parsePracticeSubject(subject string) (ps *PracticeSubject) {
	// If we have an old Municipal Status form, the subject doesn't have the
	// full practice details; it only has the jurisdiction.
	if a.xsc.Type.Tag == jurisstat.Tag21 {
		ps = &PracticeSubject{Jurisdiction: subject}
	} else if match := practiceRE.FindStringSubmatch(subject); match != nil {
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

const formPracticeSubjectResponse = `The %s field of this form does not have the
correct format. It should have the word "Practice" followed by four
comma-separated fields:
    Practice CallSign, FirstName, Jurisdiction, NetDate
NetDate should be in the form MM/DD/YYYY.`
const formSubjectResponse = `This message has
    Subject: %s
but, based on the contents of the form, it should have
    Subject: %s
PackItForms automatically generates the Subject line from the form contents; it
should not be overridden manually.`
const handlingOrderCodeResponse = `The Subject line of this message contains an
invalid Handling Order code (%s). The valid codes are "I" for Immediate, "P" for
Priority, and "R" for Routine.`
const msgNumFormatResponse = `The message number of this message is not
formatted correctly.  It should have a format like "XND-042P", containing:
  - a three-character prefix (usually the sender's call sign suffix),
  - a dash,
  - a number with at least three digits, and
  - a "P", "M", or "R" suffix.
All letters should be upper case.  In Outpost, the format of the message number
is set in the Message Settings dialog, which should be configured according to
county standards.`
const msgNumPrefixResponse = `The message number of this message has the prefix
"%s".  The prefix should be the last three characters of your call sign, "%s".`
const practiceAsFormNameResponse = `This message has an incorrect subject line
format.  According to SCCo standards, the subject should look like
    AAA-111P_R_Practice A6AAA, Charlie, Nowhereville, 03/16/2019
    (msgnum) |            |    (name)   (city)        (net date)
      (handling order)  (call sign)
Note that there is no underline after the word "Practice".`
const practiceSubjectFormatResponse = `The Subject of this message does not have
the correct format.  After the message number and handling order, it should have
the word "Practice" followed by four comma-separated fields:
    Practice CallSign, FirstName, Jurisdiction, NetDate
NetDate should be in the form MM/DD/YYYY.`
const subjectFormatResponse = `This message has an incorrect subject line
format.  According to SCCo standards, the subject should look like
    AAA-111P_R_Practice A6AAA, Charlie, Nowhereville, 03/16/2019
    (msgnum) |            |    (name)   (city)        (net date)
      (handling order)  (call sign)`
const subjectHasSeverityResponse = `The Subject line of this message contains
both a Severity code and a Handling Order code ("_%s/%s_").  This is an outdated
Subject line style.  Current SCCo standards include only the Handling Order code
on the Subject line ("_%[2]s_").`
const subjectPlainFormResponse = `This message has a form name on the subject
line, but does not contain a recognizable form.  If this is a plain text
message, there should be no form name between the handling order and the word
"Practice".  If this is a form message, the form is improperly encoded and could
not be recognized.`
const unknownJurisdictionResponse = `The jurisdiction "%s" is not recognized.
Please use one of the recognized jurisdiction names or abbreviations.`
