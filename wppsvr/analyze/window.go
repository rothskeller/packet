package analyze

import (
	"fmt"
)

// Problem codes
const (
	ProblemMessageTooEarly  = "MessageTooEarly"
	ProblemMessageTooLate   = "MessageTooLate"
	ProblemSessionDateWrong = "SessionDateWrong"
)

func init() {
	ProblemLabel[ProblemMessageTooEarly] = "message before start of practice session"
	ProblemLabel[ProblemMessageTooLate] = "message after end of practice session"
	ProblemLabel[ProblemSessionDateWrong] = "incorrect net date in subject"
}

// checkPracticeWindow verifies that the message was sent within the proper
// window of time.  It also verifies that the report end date matches the net
// date given on the subject line.
func (a *Analysis) checkPracticeWindow() {
	// This check doesn't apply to non-human messages.
	if a.msg.Message() == nil {
		return
	}
	if a.msg.Base().DeliveryTime.Before(a.session.Start) {
		// This message occurred before the start of the current
		// session, so we're not going to count it.  But, which session
		// was it intended for?  Was it too early for the current
		// session, or too late for the previous session?
		if !a.subjectDate.IsZero() && a.subjectDate.Before(a.session.Start) {
			// It was too late for the previous session.
			a.problems = append(a.problems, &problem{
				code:    ProblemMessageTooLate,
				subject: "Message after end of practice session",
				response: fmt.Sprintf(`
This message arrived at %s on %s.  That was too late to be counted for the %s
on %s.
`, a.toBBS, a.msg.Base().DeliveryTime.Format("2006-01-02 at 15:04"), a.session.Name, a.subjectDate.Format("January 2")),
				invalid:    true,
				references: refWeeklyPractice,
			})
			return
		}
		a.problems = append(a.problems, &problem{
			code:    ProblemMessageTooEarly,
			subject: "Message before start of practice session",
			response: fmt.Sprintf(`
This message arrived at %s on %s.  However, practice messages for %s aren't
accepted until %s.
`, a.toBBS, a.msg.Base().DeliveryTime.Format("2006-01-02 at 15:04"), a.session.Name, a.session.Start.Format("2006-01-02 at 15:04")),
			invalid:    true,
			references: refWeeklyPractice,
		})
		return
	}
	// If the subject has a target net date that's wrong, note that.
	if !a.subjectDate.IsZero() &&
		(a.subjectDate.Year() != a.session.End.Year() ||
			a.subjectDate.Month() != a.session.End.Month() ||
			a.subjectDate.Day() != a.session.End.Day()) {
		a.problems = append(a.problems, &problem{
			code:    ProblemSessionDateWrong,
			subject: "Incorrect net date in subject",
			response: fmt.Sprintf(`
This message is being counted for %s on %s, but the subject line says it's
intended for a net on %s.  This may indicate that the message was sent to the
wrong net.
`, a.session.Name, a.session.Start.Format("January 2"), a.subjectDate.Format("January 2")),
			warning:    true, // not an error
			references: refWeeklyPractice,
		})
	}
}
