package analyze

import (
	"strings"

	"github.com/rothskeller/packet/pktmsg"
)

func init() {
	Problems[ProbMessageFromWinlink.Code] = ProbMessageFromWinlink
	Problems[ProbMessageNotASCII.Code] = ProbMessageNotASCII
	Problems[ProbMessageNotPlainText.Code] = ProbMessageNotPlainText
}

// ProbMessageNotASCII is raised when a message contains non-ASCII characters.
var ProbMessageNotASCII = &Problem{
	Code: "MessageNotASCII",
	detect: func(a *Analysis) bool {
		return strings.IndexFunc(a.msg.Body, nonASCII) >= 0
	},
}

func nonASCII(r rune) bool {
	return r > 126 || (r < 32 && r != '\t' && r != '\n')
}

// ProbMessageFromWinlink is raised when a message comes from a winlink.org
// address and has a content-transfer-encoding of quoted-printable.  (All
// Winlink messages do, at the time this code was written.)  This is called out
// separately from ProbMessageNotPlainText because the problem response text is
// different.
var ProbMessageFromWinlink = &Problem{
	Code: "MessageFromWinlink",
	detect: func(a *Analysis) bool {
		return strings.Contains(a.msg.ReturnAddress(), "winlink.org") &&
			a.msg.Header.Get("Content-Transfer-Encoding") == "quoted-printable"
	},
}

// ProbMessageNotPlainText is raised when a message uses a content-type other
// than text/plain or a content-transfer-encoding other than binary, 7bit, or
// 8bit (and the return address does not contain "winlink.org").
var ProbMessageNotPlainText = &Problem{
	Code: "MessageNotPlainText",
	ifnot: []*Problem{
		// This check does not apply to messages that came from Winlink;
		// they get their own special problem message.
		ProbMessageFromWinlink,
	},
	detect: func(a *Analysis) bool {
		return a.msg.Flags&pktmsg.NotPlainText != 0
	},
}
