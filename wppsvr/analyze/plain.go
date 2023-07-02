package analyze

import (
	"strings"
)

// This file contains the problem checks that are run only against messages that
// are not known PackItForms forms messages.  Usually these are plain text
// messages, but they could also be unknown forms, or corrupted forms.  These
// checks are run after the common checks.  They appear in the order they are
// run, although some are skipped based on the message type or the results of
// previous checks.

func init() {
	ProblemLabels["NoCallSign"] = "no call sign in message"
	ProblemLabels["SubjectFormat"] = "incorrect subject line format"
	ProblemLabels["SubjectHasSeverity"] = "severity on subject line"
	ProblemLabels["HandlingOrderCode"] = "unknown handling order code"
	ProblemLabels["PracticeAsFormName"] = `incorrect subject line format (underline after "Practice")`
	ProblemLabels["FormCorrupt"] = "incorrectly encoded form"
	ProblemLabels["SubjectPlainForm"] = "form name in subject of non-form message"
}

func (a *Analysis) noCallSign() bool {
	if a.FromCallSign == "" {
		return a.reportProblem("NoCallSign", 0, noCallSignResponse)
	}
	return false
}

const noCallSignResponse = `This message cannot be counted because it's not
clear who sent it.  There is no call sign in the return address or after the
word "Practice" on the subject line.  In order for a message to count, there
must be a call sign in at least one of those places.`

func (a *Analysis) subjectFormat() bool {
	if a.key.OriginMsgID == "" {
		return a.reportProblem("SubjectFormat", refSubjectLine|refWeeklyPractice, subjectFormatResponse)
	}
	return false
}

const subjectFormatResponse = `This message has an incorrect subject line
format.  According to SCCo standards, the subject should look like
    AAA-111P_R_Practice A6AAA, Charlie, Nowhereville, 03/16/2019
    (msgnum) |            |    (name)   (city)        (net date)
      (handling order)  (call sign)`

func (a *Analysis) subjectHasSeverity() bool {
	if a.severity != "" {
		var hcode = a.key.Handling
		if hcode != "" {
			hcode = hcode[:1]
		}
		return a.reportProblem("SubjectHasSeverity", refSubjectLine, subjectHasSeverityResponse, a.severity, hcode)
	}
	return false
}

const subjectHasSeverityResponse = `The Subject line of this message contains
both a Severity code and a Handling Order code ("_%s/%s_").  This is an outdated
Subject line style.  Current SCCo standards include only the Handling Order code
on the Subject line ("_%[2]s_").`

func (a *Analysis) handlingOrderCode() bool {
	switch a.key.Handling {
	case "ROUTINE", "PRIORITY", "IMMEDIATE":
		return false
	default:
		return a.reportProblem("HandlingOrderCode", refSubjectLine, handlingOrderCodeResponse, a.key.Handling)
	}
}

const handlingOrderCodeResponse = `The Subject line of this message contains an
invalid Handling Order code (%s). The valid codes are "I" for Immediate, "P" for
Priority, and "R" for Routine.`

func (a *Analysis) practiceAsFormName() bool {
	if strings.EqualFold(a.formtag, "Practice") {
		return a.reportProblem("PracticeAsFormName", refSubjectLine|refWeeklyPractice, practiceAsFormNameResponse)
	}
	return false
}

const practiceAsFormNameResponse = `This message has an incorrect subject line
format.  According to SCCo standards, the subject should look like
    AAA-111P_R_Practice A6AAA, Charlie, Nowhereville, 03/16/2019
    (msgnum) |            |    (name)   (city)        (net date)
      (handling order)  (call sign)
Note that there is no underline after the word "Practice".`

func (a *Analysis) formCorrupt() bool {
	if a.corruptForm {
		return a.reportProblem("FormCorrupt", 0, formCorruptResponse)
	}
	return false
}

const formCorruptResponse = `This message appears to contain an encoded form,
but the encoding is incorrect.  It appears to have been created or edited by
software other than the current PackItForms software.  Please use current
PackItForms software to encode messages containing forms.`

func (a *Analysis) subjectPlainForm() bool {
	if a.formtag != "" && !strings.EqualFold(a.formtag, "Practice") {
		return a.reportProblem("SubjectPlainForm", refSubjectLine, subjectPlainFormResponse)
	}
	return false
}

const subjectPlainFormResponse = `This message has a form name on the subject
line, but does not contain a recognizable form.  If this is a plain text
message, there should be no form name between the handling order and the word
"Practice".  If this is a form message, the form is improperly encoded and could
not be recognized.`
