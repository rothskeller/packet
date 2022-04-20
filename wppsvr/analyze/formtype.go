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
		code    string
		allowed []string
		article string
	)
	// This check only applies to human messages.
	if a.msg.Message() == nil {
		return
	}
	code = a.msg.TypeCode()
	for _, mtype := range a.session.MessageTypes {
		if mtype == code {
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
		code:    ProblemMessageTypeWrong,
		subject: "Incorrect message type",
		response: fmt.Sprintf(`
This message is %s %s.  For the %s on %s, %s %s is expected.
`, a.msg.TypeArticle(), a.msg.TypeName(), a.session.Name, a.session.End.Format("January 2"), article,
			english.Conjoin(allowed, "or")),
		references: refWeeklyPractice,
	})
}
