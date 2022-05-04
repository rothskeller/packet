package analyze

import (
	"fmt"

	"steve.rothskeller.net/packet/wppsvr/config"
	"steve.rothskeller.net/packet/wppsvr/english"
)

// Problem codes
const (
	ProblemMessageTypeWrong = "MessageTypeWrong"
)

func init() {
	ProblemLabel[ProblemMessageTypeWrong] = "incorrect message type"
}

// checkCorrectForm verifies that the received form has the expected type for
// the practice session.
func (a *Analysis) checkCorrectForm() {
	var (
		tag     string
		allowed []string
		article string
	)
	tag = a.xsc.TypeTag()
	for _, mtype := range a.session.MessageTypes {
		if mtype == tag {
			return
		}
	}
	for i, code := range a.session.MessageTypes {
		mtype := config.LookupMessageType(code)
		allowed = append(allowed, mtype.TypeName())
		if i == 0 {
			article = mtype.TypeArticle()
		}
	}
	a.problems = append(a.problems, &problem{
		code: ProblemMessageTypeWrong,
		response: fmt.Sprintf(`
This message is %s %s.  For the %s on %s, %s %s is expected.
`, a.xsc.TypeArticle(), a.xsc.TypeName(), a.session.Name, a.session.End.Format("January 2"), article,
			english.Conjoin(allowed, "or")),
		references: refWeeklyPractice,
	})
}
