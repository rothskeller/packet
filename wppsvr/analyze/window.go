package analyze

func init() {
	Problems[ProbMessageTooEarly.Code] = ProbMessageTooEarly
	Problems[ProbMessageTooLate.Code] = ProbMessageTooLate
	Problems[ProbSessionDate.Code] = ProbSessionDate
}

// ProbMessageTooEarly is raised when
var ProbMessageTooEarly = &Problem{
	Code:  "MessageTooEarly",
	ifnot: []*Problem{ProbMessageTooLate},
	detect: func(a *Analysis) (bool, string) {
		return a.msg.Date().Before(a.session.Start), ""
	},
	Variables: variableMap{
		"SESSIONSTART": func(a *Analysis) string {
			return a.session.Start.Format("2006-01-02 at 15:04")
		},
	},
}

// ProbMessageTooLate is raised when
var ProbMessageTooLate = &Problem{
	Code: "MessageTooLate",
	detect: func(a *Analysis) (bool, string) {
		return a.msg.Date().Before(a.session.Start) && a.Practice != nil && !a.Practice.NetDate.IsZero() && a.Practice.NetDate.Before(a.session.Start), ""
	},
	Variables: variableMap{
		"SUBJECTDATE": func(a *Analysis) string {
			return a.Practice.NetDate.Format("January 2")
		},
	},
}

// ProbSessionDate is raised when
var ProbSessionDate = &Problem{
	Code:  "SessionDate",
	ifnot: []*Problem{ProbMessageTooEarly, ProbMessageTooLate},
	detect: func(a *Analysis) (bool, string) {
		if a.Practice == nil || a.Practice.NetDate.IsZero() {
			return false, ""
		}
		nd := a.Practice.NetDate
		return nd.Year() != a.session.End.Year() || nd.Month() != a.session.End.Month() || nd.Day() != a.session.End.Day(), ""
	},
	Variables: variableMap{
		"SUBJECTDATE": func(a *Analysis) string {
			return a.Practice.NetDate.Format("January 2")
		},
	},
}
