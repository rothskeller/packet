package analyze

import (
	"regexp"

	"steve.rothskeller.net/packet/pktmsg"
)

// Problem codes
const (
	ProblemMsgNumFormat = "MsgNumFormat"
)

func init() {
	ProblemLabel[ProblemMsgNumFormat] = "incorrect message number format"
}

var msgnumRE = regexp.MustCompile(`^(?:[A-Z][A-Z][A-Z]|[A-Z][0-9][A-Z0-9]|[0-9][A-Z][A-Z])-\d\d\d+[A-Z]$`)

// checkMessageNumber checks for problems with the message number.
func (a *Analysis) checkMessageNumber() {
	var msgnum string

	if a.msg.Message() == nil {
		return // disable this check for non-human messages
	} else if msg := a.msg.SCCoForm(); msg != nil && msg.OriginMessageNumber != "" {
		msgnum = msg.OriginMessageNumber
	} else if msg, ok := a.msg.(*pktmsg.RxICS213Form); ok && msg.MessageNumber != "" {
		msgnum = msg.MessageNumber
	} else if msg := a.msg.Form(); msg != nil && msg.Fields["MsgNo"] != "" {
		msgnum = msg.Fields["MsgNo"]
	} else if msg := a.msg.Message(); msg.MessageNumber != "" {
		msgnum = msg.MessageNumber
	} else {
		// We couldn't find a message number to check.  The problem will
		// be reported elsewhere (either as an invalid form or an
		// invalid subject line).
		return
	}
	if !msgnumRE.MatchString(msgnum) {
		a.problems = append(a.problems, &problem{
			code:    ProblemMsgNumFormat,
			subject: "Incorrect message number format",
			response: `
The message number of this message is not formatted correctly.  It should have
a format like "XND-042P", containing:
  - a three-character prefix (usually the sender's call sign suffix),
  - a dash,
  - a number with at least three digits, and
  - a "P" suffix.
All letters should be upper case.  In Outpost, the format of the message
number is set in the Message Settings dialog, which should be configured
according to county standards.
`,
			references: refOutpostConfig,
		})
	}
}
