package analyze

func init() {
	Problems[ProbMessageTooEarly.Code] = ProbMessageTooEarly
	Problems[ProbMessageTooLate.Code] = ProbMessageTooLate
	Problems[ProbSessionDate.Code] = ProbSessionDate
}

// ProbMessageTooEarly is raised when a message is received between sessions
// (but not if it was intended for a previous session).
var ProbMessageTooEarly = &Problem{
	Code: "MessageTooEarly",
	ifnot: []*Problem{
		// This check does not apply if the message was too late for a
		// preceding session.
		ProbMessageTooLate,
	},
	detect: func(a *Analysis) bool {
		return a.msg.Date().Before(a.session.Start)
	},
	Variables: variableMap{
		"SESSIONSTART": func(a *Analysis) string {
			return a.session.Start.Format("2006-01-02 at 15:04")
		},
	},
}

// ProbMessageTooLate is raised when a message is received between sessions,
// that is intended for a previous session.
var ProbMessageTooLate = &Problem{
	Code: "MessageTooLate",
	detect: func(a *Analysis) bool {
		// This check does not apply if we can't tell what session the
		// message was intended for.
		if a.Practice == nil || a.Practice.NetDate.IsZero() {
			return false
		}
		// The check.
		return a.msg.Date().Before(a.session.Start) && a.Practice.NetDate.Before(a.session.Start)
	},
	Variables: variableMap{
		"SUBJECTDATE": func(a *Analysis) string {
			return a.Practice.NetDate.Format("January 2")
		},
	},
}

// ProbSessionDate is raised when a message is received during a session that is
// intended for a different session.
var ProbSessionDate = &Problem{
	Code:  "SessionDate",
	ifnot: []*Problem{ProbMessageTooEarly, ProbMessageTooLate},
	detect: func(a *Analysis) bool {
		// This check does not apply if we can't tell what session the
		// message was intended for.
		if a.Practice == nil || a.Practice.NetDate.IsZero() {
			return false
		}
		// The check.
		nd := a.Practice.NetDate
		return nd.Year() != a.session.End.Year() || nd.Month() != a.session.End.Month() || nd.Day() != a.session.End.Day()
	},
	Variables: variableMap{
		"SUBJECTDATE": func(a *Analysis) string {
			return a.Practice.NetDate.Format("January 2")
		},
	},
}
