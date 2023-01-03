package analyze

import (
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/english"
	"github.com/rothskeller/packet/xscmsg"
)

func init() {
	Problems[ProbMessageTypeWrong.Code] = ProbMessageTypeWrong
}

// ProbMessageTypeWrong is raised when the message type is not one of the
// expected types for the session, and the message is coming from in-county.
var ProbMessageTypeWrong = &Problem{
	Code: "MessageTypeWrong",
	ifnot: []*Problem{
		// This check does not apply if we have a message that appears
		// to be a form but couldn't be parsed as one.
		ProbFormCorrupt,
	},
	detect: func(a *Analysis) bool {
		// This check does not apply if we have a plain text message
		// coming from outside our BBS network.  (It could be from
		// outside of Santa Clara County where the SCC forms are not
		// available.)
		if a.xsc.Type.Tag == xscmsg.PlainTextTag && config.Get().BBSes[a.msg.FromBBS()] == nil {
			return false
		}
		// The check.
		return !inList(a.session.MessageTypes, a.xsc.Type.Tag)
	},
	Variables: variableMap{
		"AEXPECTTYPE": func(a *Analysis) string {
			var (
				allowed []string
				article string
			)
			for i, code := range a.session.MessageTypes {
				mtype := config.LookupMessageType(code)
				allowed = append(allowed, mtype.Type.Name)
				if i == 0 {
					article = mtype.Type.Article
				}
			}
			return article + " " + english.Conjoin(allowed, "or")
		},
	},
}
