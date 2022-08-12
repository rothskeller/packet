package analyze

import (
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/english"
)

func init() {
	Problems[ProbMessageTypeWrong.Code] = ProbMessageTypeWrong
}

// ProbMessageTypeWrong is raised when the message type is not one of the
// expected types for the session, and the message is coming from in-county.
var ProbMessageTypeWrong = &Problem{
	Code:  "MessageTypeWrong",
	after: []*Problem{ProbDeliveryReceipt}, // sets a.xsc
	ifnot: []*Problem{ProbFormCorrupt, ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		var tag = a.xsc.TypeTag()
		for _, mtype := range a.session.MessageTypes {
			if mtype == tag {
				return false, ""
			}
		}
		if _, ok := a.xsc.(*config.PlainTextMessage); ok {
			// It's a plain text message and we're expecting a form.  That's
			// OK if it's coming from somewhere other than one of our BBSes.
			if _, ok := config.Get().BBSes[a.msg.FromBBS()]; !ok {
				return false, ""
			}
		}
		return true, ""
	},
	Variables: variableMap{
		"AEXPECTTYPE": func(a *Analysis) string {
			var (
				allowed []string
				article string
			)
			for i, code := range a.session.MessageTypes {
				mtype := config.LookupMessageType(code)
				allowed = append(allowed, mtype.TypeName())
				if i == 0 {
					article = mtype.TypeArticle()
				}
			}
			return article + " " + english.Conjoin(allowed, "or")
		},
	},
	references: refWeeklyPractice,
}
