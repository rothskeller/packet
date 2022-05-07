package analyze

import (
	"regexp"

	"steve.rothskeller.net/packet/xscmsg"
)

// Problem codes
const (
	ProblemMsgNumFormat = "MsgNumFormat"
)

func init() {
	ProblemLabel[ProblemMsgNumFormat] = "incorrect message number format"
}

var msgnumRE = regexp.MustCompile(`^(?:[A-Z][A-Z][A-Z]|[A-Z][0-9][A-Z0-9]|[0-9][A-Z][A-Z])-\d\d\d+[PM]$`)

// checkMessageNumber checks for problems with the message number.
func (a *Analysis) checkMessageNumber() {
	var msgnum string

	// Which message number should we check?
	if xsc, ok := a.xsc.(interface{ OriginNumber() string }); ok {
		// It's a form, so check the number in the form.
		msgnum = xsc.OriginNumber()
		// But don't check it if the form failed to validate; the
		// validation errors probably already reported it.
		if xsc, ok := a.xsc.(interface{ Validate(bool) []string }); ok && len(xsc.Validate(true)) != 0 {
			return
		}
	} else if xscsubj := xscmsg.ParseSubject(a.msg.Header.Get("Subject")); xscsubj != nil {
		// It's not a form, but we were able to parse a message number
		// out of the subject line, so check that.
		msgnum = xscsubj.MessageNumber
	} else {
		// No message number to check.  The problem will be reported
		// elsewhere as a bad subject line.
		return
	}
	if !msgnumRE.MatchString(msgnum) {
		a.problems = append(a.problems, &problem{
			code: ProblemMsgNumFormat,
			response: `
The message number of this message is not formatted correctly.  It should have
a format like "XND-042P", containing:
  - a three-character prefix (usually the sender's call sign suffix),
  - a dash,
  - a number with at least three digits, and
  - a "P" or "M" suffix.
All letters should be upper case.  In Outpost, the format of the message
number is set in the Message Settings dialog, which should be configured
according to county standards.
`,
			references: refOutpostConfig,
		})
	}
}
