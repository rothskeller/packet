package analyze

import (
	"regexp"

	"github.com/rothskeller/packet/xscmsg"
)

func init() {
	Problems[ProbMsgNumFormat.Code] = ProbMsgNumFormat
	Problems[ProbMsgNumPrefix.Code] = ProbMsgNumPrefix
}

type formWithOriginNumber interface {
	OriginNumber() string
}

var msgnumRE = regexp.MustCompile(`^(?:[A-Z][A-Z][A-Z]|[A-Z][0-9][A-Z0-9]|[0-9][A-Z][A-Z])-\d\d\d+[PMR]$`)

// ProbMsgNumFormat is raised when the message number does not meet county
// standards.
var ProbMsgNumFormat = &Problem{
	Code:  "MsgNumFormat",
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
}

var fccCallSignRE = regexp.MustCompile(`^[AKNW][A-Z]?[0-9][A-Z]{1,3}$`)

// ProbMsgNumPrefix is raised when the message number prefix is not the last
// three characters of the sender's call sign.
var ProbMsgNumPrefix = &Problem{
	Code:  "MsgNumPrefix",
	after: []*Problem{ProbDeliveryReceipt, ProbNoCallSign}, // set a.xsc, a.fromCallSign
	ifnot: []*Problem{ProbMsgNumFormat, ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		var msgnum string

		if !fccCallSignRE.MatchString(a.fromCallSign) {
			return false, "" // prefix not checked for tactical calls
		}
		if xsc, ok := a.xsc.(formWithOriginNumber); ok {
			msgnum = xsc.OriginNumber()
		} else if xscsubj := xscmsg.ParseSubject(a.msg.Header.Get("Subject")); xscsubj != nil {
			msgnum = xscsubj.MessageNumber
		} else {
			return false, ""
		}
		return msgnum[:3] != a.fromCallSign[len(a.fromCallSign)-3:], ""
	},
	Variables: variableMap{
		"ACTUALPFX": func(a *Analysis) string {
			if xsc, ok := a.xsc.(formWithOriginNumber); ok {
				return xsc.OriginNumber()[:3]
			}
			xscsubj := xscmsg.ParseSubject(a.msg.Header.Get("Subject"))
			return xscsubj.MessageNumber[:3]
		},
		"EXPECTPFX": func(a *Analysis) string {
			return a.fromCallSign[len(a.fromCallSign)-3:]
		},
	},
}
