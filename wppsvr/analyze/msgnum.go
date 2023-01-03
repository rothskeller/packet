package analyze

import (
	"regexp"

	"github.com/rothskeller/packet/xscmsg"
)

func init() {
	Problems[ProbMsgNumFormat.Code] = ProbMsgNumFormat
	Problems[ProbMsgNumPrefix.Code] = ProbMsgNumPrefix
}

var msgnumRE = regexp.MustCompile(`^(?:[A-Z][A-Z][A-Z]|[A-Z][0-9][A-Z0-9]|[0-9][A-Z][A-Z])-\d\d\d+[PMR]$`)

// ProbMsgNumFormat is raised when the message number does not meet county
// standards.
var ProbMsgNumFormat = &Problem{
	Code: "MsgNumFormat",
	detect: func(a *Analysis) bool {
		// If the message is a form with an Origin Message Number field,
		// check the message number in that field.
		if f := a.xsc.KeyField(xscmsg.FOriginMsgNo); f != nil {
			return !msgnumRE.MatchString(f.Value)
		}
		// If we were able to parse the subject line of the message,
		// check the message number from that.
		if a.subject != nil {
			return !msgnumRE.MatchString(a.subject.MessageNumber)
		}
		// No message number to check.  The problem will be reported
		// elsewhere as a bad subject line.
		return false
	},
}

var fccCallSignRE = regexp.MustCompile(`^[AKNW][A-Z]?[0-9][A-Z]{1,3}$`)

// ProbMsgNumPrefix is raised when the message number prefix is not the last
// three characters of the sender's call sign.
var ProbMsgNumPrefix = &Problem{
	Code: "MsgNumPrefix",
	ifnot: []*Problem{
		// This check does not apply if the message number format is
		// wrong.
		ProbMsgNumFormat,
	},
	detect: func(a *Analysis) bool {
		// This check does not apply if the sender call sign is not an
		// FCC call sign (i.e., missing or tactical call).
		if !fccCallSignRE.MatchString(a.FromCallSign) {
			return false
		}
		// The prefix we want is the last three characters of the sender
		// call sign.
		want := a.FromCallSign[len(a.FromCallSign)-3:]
		// If the message is a form with an Origin Message Number field,
		// check the message number in that field.
		if f := a.xsc.KeyField(xscmsg.FOriginMsgNo); f != nil {
			return f.Value[:3] != want
		}
		// If we were able to parse the subject line of the message,
		// check the message number from that.
		if a.subject != nil {
			return a.subject.MessageNumber[:3] != want
		}
		// No message number to check.  The problem will be reported
		// elsewhere as a bad subject line.
		return false
	},
	Variables: variableMap{
		"ACTUALPFX": func(a *Analysis) string {
			if f := a.xsc.KeyField(xscmsg.FOriginMsgNo); f != nil {
				return f.Value[:3]
			}
			xscsubj := xscmsg.ParseSubject(a.msg.Header.Get("Subject"))
			return xscsubj.MessageNumber[:3]
		},
		"EXPECTPFX": func(a *Analysis) string {
			return a.FromCallSign[len(a.FromCallSign)-3:]
		},
	},
}
