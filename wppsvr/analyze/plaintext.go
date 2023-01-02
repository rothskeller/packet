package analyze

import (
	"strings"

	"github.com/rothskeller/packet/pktmsg"
)

func init() {
	Problems[ProbMessageNotASCII.Code] = ProbMessageNotASCII
	Problems[ProbMessageNotPlainText.Code] = ProbMessageNotPlainText
}

// ProbMessageNotASCII is raised when a message contains non-ASCII characters.
var ProbMessageNotASCII = &Problem{
	Code: "MessageNotASCII",
	detect: func(a *Analysis) (bool, string) {
		return strings.IndexFunc(a.msg.Body, nonASCII) >= 0, ""
	},
}

// ProbMessageNotPlainText is raised when a message uses a content-type other
// than text/plain or a content-transfer-encoding other than binary, 7bit, or
// 8bit.
var ProbMessageNotPlainText = &Problem{
	Code: "MessageNotPlainText",
	detect: func(a *Analysis) (bool, string) {
		if a.msg.Flags&pktmsg.NotPlainText == 0 {
			return false, ""
		}
		if strings.Contains(a.msg.ReturnAddress(), "winlink.org") {
			return true, "winlink"
		}
		return true, "normal"
	},
}

func nonASCII(r rune) bool {
	return r > 126 || (r < 32 && r != '\t' && r != '\n')
}
