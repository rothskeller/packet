package analyze

func init() {
	ProblemLabels["MessageTooEarly"] = "message before start of practice session"
	ProblemLabels["MessageTooLate"] = "message after end of practice session"
	ProblemLabels["SessionDate"] = "incorrect net date in subject"
}

// checkDate checks for problems with when the message was sent and what session
// it was intended for.
func (a *Analysis) checkDate() {
	// Check whether the message date is before the start of the current
	// session.
	if a.msg.Date().Before(a.session.Start) {
		// If we know (from the "Practice ..." subject) what session
		// this message was intended for, and it was before the current
		// session, this message is too late.  Otherwise we call it too
		// early.
		if a.Practice != nil && !a.Practice.NetDate.IsZero() && a.Practice.NetDate.Before(a.session.Start) {
			a.reportProblem("MessageTooLate", refWeeklyPractice, messageTooLateResponse,
				a.toBBS, a.msg.Date().Format("2006-01-02 at 15:04"), a.session.Name,
				a.Practice.NetDate.Format("January 2"))
		} else {
			a.reportProblem("MessageTooEarly", refWeeklyPractice, messageTooEarlyResponse,
				a.toBBS, a.msg.Date().Format("2006-01-02 at 15:04"), a.session.Name,
				a.session.Start.Format("2006-01-02 at 15:04"))
		}
		return
	}
	// The message isn't too early for the current session.  But if we know
	// what session it was intended for, let's make sure it was the current
	// one.
	if a.Practice != nil {
		nd := a.Practice.NetDate
		if !nd.IsZero() && (nd.Year() != a.session.End.Year() || nd.Month() != a.session.End.Month() || nd.Day() != a.session.End.Day()) {
			a.reportProblem("SessionDate", refWeeklyPractice, sessionDateResponse,
				a.session.Name, a.session.End.Format("January 2"), a.Practice.NetDate.Format("January 2"))
		}
	}
}

const messageTooEarlyResponse = `This message arrived at %s on %s.  However,
practice messages for %s aren't accepted until %s.`
const messageTooLateResponse = `This message arrived at %s on %s.  That was too
late to be counted for the %s on %s.`
const sessionDateResponse = `This message is being counted for %s on %s, but the
subject line says it's intended for a net on %s.  This may indicate that the
message was sent to the wrong net.`
