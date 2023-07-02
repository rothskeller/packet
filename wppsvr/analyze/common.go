package analyze

// This file contains the problem checks that are run against all human
// messages.  They appear in the order they are run, although some are skipped
// based on the message type or the results of previous checks.

import (
	"regexp"
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/delivrcpt"
	"github.com/rothskeller/packet/message/plaintext"
	"github.com/rothskeller/packet/message/readrcpt"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/english"
)

// MultipleMessagesFromAddress is a pseudo-problem used during reporting, to
// mark messages that have been superseded by other messages.
const MultipleMessagesFromAddress = "MultipleMessagesFromAddress"

func init() {
	ProblemLabels[MultipleMessagesFromAddress] = "multiple messages from this address"
	ProblemLabels["MessageCorrupt"] = "message could not be parsed"
	ProblemLabels["BounceMessage"] = "message has no return address (probably auto-response)"
	ProblemLabels["DeliveryReceipt"] = "DELIVERED receipt message"
	ProblemLabels["ReadReceipt"] = "unexpected READ receipt message"
	ProblemLabels["MessageFromWinlink"] = "message sent from Winlink"
	ProblemLabels["MessageNotPlainText"] = "not a plain text message"
	ProblemLabels["MessageNotASCII"] = "message has non-ASCII characters"
	ProblemLabels["FromBBSDown"] = "message from incorrect BBS (simulated outage)"
	ProblemLabels["ToBBSDown"] = "message to incorrect BBS (simulated outage)"
	ProblemLabels["ToBBS"] = "message to incorrect BBS"
	ProblemLabels["PracticeSubjectFormat"] = "incorrect practice message details"
	ProblemLabels["FormPracticeSubject"] = "incorrect practice message details in form"
	ProblemLabels["UnknownJurisdiction"] = "unknown jurisdiction"
	ProblemLabels["MsgNumFormat"] = "incorrect message number format"
	ProblemLabels["MsgNumPrefix"] = "incorrect message number prefix"
	ProblemLabels["MessageTooLate"] = "message after end of practice session"
	ProblemLabels["MessageTooEarly"] = "message before start of practice session"
	ProblemLabels["SessionDate"] = "incorrect net date in subject"
	ProblemLabels["MessageTypeWrong"] = "incorrect message type"
}

func (a *Analysis) messageCorrupt(parseErr error) bool {
	if parseErr != nil {
		a.problems["MessageCorrupt"] = struct{}{}
		return true
	}
	return false
}

func (a *Analysis) bounceMessage() bool {
	if a.env.Autoresponse {
		a.problems["BounceMessage"] = struct{}{}
		return true
	}
	return false
}

func (a *Analysis) deliveryReceipt() bool {
	if _, ok := a.msg.(*delivrcpt.DeliveryReceipt); ok {
		a.problems["DeliveryReceipt"] = struct{}{}
		return true
	}
	return false
}

func (a *Analysis) readReceipt() bool {
	if _, ok := a.msg.(*readrcpt.ReadReceipt); ok {
		a.reportProblem("ReadReceipt", refOutpostConfig, readReceiptResponse)
		return true
	}
	return false
}

const readReceiptResponse = `This message is an Outpost "read receipt", which
should not have been sent.  Most likely, your Outpost installation has the
"Auto-Read Receipt" setting turned on.  The SCCo-standard Outpost configuration
specifies that this setting should be turned off.  You can find it on the
Receipts tab of the Message Settings dialog in Outpost.`

func (a *Analysis) messageFromWinlink() bool {
	if a.env.NotPlainText && strings.Contains(a.env.ReturnAddr, "winlink.org") {
		return a.reportProblem("MessageFromWinlink", 0, messageFromWinlinkResponse)
	}
	return false
}

const messageFromWinlinkResponse = `This message was sent from Winlink.  Winlink
should not be used for emergency communications, unless no alternatives are
available, because it uses a message encoding system ("quoted-printable") that
Outpost cannot decode.  As a result, some messages (particularly those with long
lines and those containing equals signs) may be garbled in transmission.`

func (a *Analysis) messageNotPlainText() bool {
	if a.env.NotPlainText {
		return a.reportProblem("MessageNotPlainText", 0, messageNotPlainTextResponse)
	}
	return false
}

const messageNotPlainTextResponse = `This message is not a plain text message.
All SCCo packet messages should be plain text only.  ("Rich text" or
HTML-formatted messages, common in email systems, are far larger than plain text
messages and put too much strain on the packet infrastructure.)  Please
configure your software to send plain text messages when sending to an SCCo
BBS.`

func (a *Analysis) messageNotASCII() bool {
	if strings.IndexFunc(a.body, nonASCII) >= 0 {
		return a.reportProblem("MessageNotASCII", 0, messageNotASCIIResponse)
	}
	return false
}
func nonASCII(r rune) bool {
	return r > 126 || (r < 32 && r != '\t' && r != '\n')
}

const messageNotASCIIResponse = `This message contains characters that are not
in the standard ASCII character set (i.e., not on a standard keyboard).
Non-standard characters should be avoided in packet messages, because the
receiving system may not know how to render them.  Note that some software may
introduce undesired non-standard characters (e.g., Microsoft Word's "smart
quotes" feature). If you use message text composed in such software, make sure
those features are disabled.`

func (a *Analysis) fromBBSDown() bool {
	if inList(a.session.DownBBSes, a.fromBBS) {
		return a.reportProblem("FromBBSDown", refWeeklyPractice, fromBBSDownResponse,
			a.fromBBS, a.session.Name, a.session.End.Format("January 2"))
	}
	return false
}

const fromBBSDownResponse = `This message was sent from %s, which has a
simulated outage for %s on %s.  Practice messages should not be sent from BBSes
that have a simulated outage.`

func (a *Analysis) toBBSDown() bool {
	if inList(a.session.DownBBSes, a.toBBS) {
		return a.reportProblem("ToBBSDown", refWeeklyPractice, toBBSDownResponse,
			a.session.CallSign, a.toBBS, a.session.Name, a.session.End.Format("January 2"),
			english.Conjoin(a.session.ToBBSes, "or"))
	}
	return false
}

const toBBSDownResponse = `This message was sent to %[1]s at %[2]s, but %[2]s
has a simulated outage for %[3]s on %[4]s.  Practice messages for this session
must be sent to %[1]s at %[5]s.`

func (a *Analysis) toBBSWrong() bool {
	if !inList(a.session.ToBBSes, a.toBBS) {
		return a.reportProblem("ToBBS", refWeeklyPractice, toBBSResponse,
			a.session.CallSign, a.toBBS, a.session.Name, a.session.End.Format("January 2"),
			english.Conjoin(a.session.ToBBSes, "or"))
	}
	return false
}

const toBBSResponse = `This message was sent to %[1]s at %[2]s.  Practice
messages for %[3]s on %[4]s must be sent to %[1]s at %[5]s.`

func (a *Analysis) practiceSubjectFormat() bool {
	if mtc := config.Get().MessageTypes[a.msg.Type().Tag]; mtc == nil || mtc.NoPracticeInfo {
		// This is an unknown form type or one that doesn't support
		// Practice... info on the Subject line.  Return without raising
		// an error.
		return false
	}
	if a.Jurisdiction == "" && a.netDate.IsZero() {
		if a.key.SubjectLabel != "" {
			return a.reportProblem("FormPracticeSubject", refWeeklyPractice, formPracticeSubjectResponse, a.key.SubjectLabel)
		}
		return a.reportProblem("PracticeSubjectFormat", refWeeklyPractice, practiceSubjectFormatResponse)
	}
	return false
}

const practiceSubjectFormatResponse = `The Subject of this message does not have
the correct format.  After the message number and handling order, it should have
the word "Practice" followed by four comma-separated fields:
    Practice CallSign, FirstName, Jurisdiction, NetDate
NetDate should be in the form MM/DD/YYYY.`
const formPracticeSubjectResponse = `The %s field of this form does not have the
correct format. It should have the word "Practice" followed by four
comma-separated fields:
    Practice CallSign, FirstName, Jurisdiction, NetDate
NetDate should be in the form MM/DD/YYYY.`

func (a *Analysis) unknownJurisdiction() bool {
	if config.Get().Jurisdictions[a.Jurisdiction] == "" {
		return a.reportProblem("UnknownJurisdiction", refBBSList|refWeeklyPractice, unknownJurisdictionResponse,
			a.Jurisdiction)
	}
	return false
}

const unknownJurisdictionResponse = `The jurisdiction "%s" is not recognized.
Please use one of the recognized jurisdiction names or abbreviations.`

var msgnumRE = regexp.MustCompile(`^(?:[A-Z][A-Z][A-Z]|[A-Z][0-9][A-Z0-9]|[0-9][A-Z][A-Z])-\d\d\d+[PMR]$`)

func (a *Analysis) msgNumFormat() bool {
	if a.key.OriginMsgID != "" && !msgnumRE.MatchString(a.key.OriginMsgID) {
		return a.reportProblem("MsgNumFormat", refOutpostConfig|refSubjectLine, msgNumFormatResponse)
	}
	return false
}

const msgNumFormatResponse = `The message number of this message is not
formatted correctly.  It should have a format like "XND-042P", containing:
  - a three-character prefix (usually the sender's call sign suffix),
  - a dash,
  - a number with at least three digits, and
  - a "P", "M", or "R" suffix.
All letters should be upper case.  In Outpost, the format of the message number
is set in the Message Settings dialog, which should be configured according to
county standards.`

func (a *Analysis) msgNumPrefix() bool {
	if a.key.OriginMsgID != "" {
		if fccCallSignRE.MatchString(a.FromCallSign) {
			act := a.key.OriginMsgID[:3]
			exp := a.FromCallSign[len(a.FromCallSign)-3:]
			if act != exp {
				return a.reportProblem("MsgNumPrefix", refSubjectLine, msgNumPrefixResponse, act, exp)
			}
		}
	}
	return false
}

const msgNumPrefixResponse = `The message number of this message has the prefix
"%s".  The prefix should be the last three characters of your call sign, "%s".`

func (a *Analysis) messageTooLate() bool {
	var date = a.env.BBSReceivedDate
	if date.IsZero() {
		date = a.env.Date
	}
	if date.Before(a.session.Start) && !a.netDate.IsZero() && a.netDate.Before(a.session.Start) {
		return a.reportProblem("MessageTooLate", refWeeklyPractice, messageTooLateResponse,
			a.toBBS, date.Format("2006-01-02 at 15:04"), a.session.Name,
			a.netDate.Format("January 2"))
	}
	return false
}

const messageTooLateResponse = `This message arrived at %s on %s.  That was too
late to be counted for the %s on %s.`

func (a *Analysis) messageTooEarly() bool {
	var date = a.env.BBSReceivedDate
	if date.IsZero() {
		date = a.env.Date
	}
	if date.Before(a.session.Start) {
		return a.reportProblem("MessageTooEarly", refWeeklyPractice, messageTooEarlyResponse,
			a.toBBS, date.Format("2006-01-02 at 15:04"), a.session.Name,
			a.session.Start.Format("2006-01-02 at 15:04"))
	}
	return false
}

const messageTooEarlyResponse = `This message arrived at %s on %s.  However,
practice messages for %s aren't accepted until %s.`

func (a *Analysis) sessionDate() bool {
	if !a.netDate.IsZero() && (a.netDate.Year() != a.session.End.Year() || a.netDate.Month() != a.session.End.Month() || a.netDate.Day() != a.session.End.Day()) {
		return a.reportProblem("SessionDate", refWeeklyPractice, sessionDateResponse,
			a.session.Name, a.session.End.Format("January 2"), a.netDate.Format("January 2"))
	}
	return false
}

const sessionDateResponse = `This message is being counted for %s on %s, but the
subject line says it's intended for a net on %s.  This may indicate that the
message was sent to the wrong net.`

func (a *Analysis) messageTypeWrong() bool {
	if _, ok := a.msg.(*plaintext.PlainText); ok && (config.Get().BBSes[a.fromBBS] == nil || a.corruptForm) {
		// Plain text message is always OK when it comes from outside
		// the county BBS system.  We also avoid raising this message
		// when we have a corrupt form.
		return false
	}
	if inList(a.session.MessageTypes, a.msg.Type().Tag) {
		return false
	}
	var (
		allowed []string
		article string
	)
	for i, code := range a.session.MessageTypes {
		mtype := message.RegisteredTypes[code]
		allowed = append(allowed, mtype.Name)
		if i == 0 {
			article = mtype.Article
		}
	}
	return a.reportProblem("MessageTypeWrong", refWeeklyPractice, messageTypeWrongResponse,
		a.msg.Type().Article, a.msg.Type().Name, a.session.Name, a.session.End.Format("January 2"),
		article, english.Conjoin(allowed, "or"))
}

const messageTypeWrongResponse = `This message is %s %s.  For the %s on %s, %s
%s is expected.`
