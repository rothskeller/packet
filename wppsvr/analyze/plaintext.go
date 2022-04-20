package analyze

import (
	"strings"
)

// Problem codes
const (
	ProblemMessageNotASCII     = "MessageNotASCII"
	ProblemMessageNotPlainText = "MessageNotPlainText"
)

func init() {
	ProblemLabel[ProblemMessageNotASCII] = "contains non-ASCII characters"
	ProblemLabel[ProblemMessageNotPlainText] = "not a plain text message"
}

// checkMessageNotPlainText checks for messages that aren't entirely plain text
// with ASCII characters.
func (a *Analysis) checkPlainText() {
	// These checks only apply to human messages.
	if a.msg.Message() == nil {
		return
	}
	// Check for a Content-Type other than text/plain, or a
	// Content-Transfer-Encoding other than 7bit or 8bit.  These conditions
	// are reflected in the NotPlainText flag set by the pktmsg package
	// while the message was being parsed.
	if a.msg.Base().NotPlainText {
		a.problems = append(a.problems, &problem{
			code:    ProblemMessageNotPlainText,
			subject: "Not a plain text message",
			response: `
This message is not a plain text message.  All SCCo packet messages should be
plain text only.  ("Rich text" or HTML-formatted messages, common in email
systems, are far larger than plain text messages and put too much strain on
the packet infrastructure.)  Please configure your software to send plain text
messages when sending to an SCCo BBS.
`,
		})
		return
	}
	// Check for the body containing non-ASCII characters.
	if strings.IndexFunc(a.msg.Base().Body, nonASCII) >= 0 {
		a.problems = append(a.problems, &problem{
			code:    ProblemMessageNotASCII,
			subject: "Non-ASCII characters",
			response: `
This message contains characters that are not in the standard ASCII character
set (i.e., not on a standard keyboard).  Non-standard characters should be
avoided in packet messages, because the receiving system may not know how to
render them.  Note that some software may introduce undesired non-standard
characters (e.g., Microsoft Word's "smart quotes" feature).  If you use
message text composed in such software, make sure those features are disabled.
`,
		})
	}
}
func nonASCII(r rune) bool {
	return r > 126 || (r < 32 && r != '\t' && r != '\n')
}
