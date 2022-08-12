package analyze

import (
	"regexp"

	"github.com/rothskeller/packet/xscmsg"
)

func init() {
	Problems[ProbMsgNumFormat.Code] = ProbMsgNumFormat
}

type formWithOriginNumber interface {
	OriginNumber() string
}

var msgnumRE = regexp.MustCompile(`^(?:[A-Z][A-Z][A-Z]|[A-Z][0-9][A-Z0-9]|[0-9][A-Z][A-Z])-\d\d\d+[PM]$`)

// ProbMsgNumFormat is raised when the message number does not meet county
// standards.
var ProbMsgNumFormat = &Problem{
	Code:  "MsgNumFormat",
	Label: "incorrect message number format",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		if xsc, ok := a.xsc.(formWithOriginNumber); ok {
			// It's a form, so check the number in the form.
			return !msgnumRE.MatchString(xsc.OriginNumber()), ""
		}
		if xscsubj := xscmsg.ParseSubject(a.msg.Header.Get("Subject")); xscsubj != nil {
			// It's not a form, but we were able to parse a message
			// number out of the subject line, so check that.
			return !msgnumRE.MatchString(xscsubj.MessageNumber), ""
		}
		// No message number to check.  The problem will be
		// reported elsewhere as a bad subject line.
		return false, ""
	},
	references: refOutpostConfig,
}
