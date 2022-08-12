package analyze

func init() {
	Problems[ProbMessageTooEarly.Code] = ProbMessageTooEarly
	Problems[ProbMessageTooLate.Code] = ProbMessageTooLate
	Problems[ProbSessionDate.Code] = ProbSessionDate
}

// ProbMessageTooEarly is raised when
var ProbMessageTooEarly = &Problem{
	Code:  "MessageTooEarly",
	after: []*Problem{ProbPracticeSubjectFormat},
	ifnot: []*Problem{ProbMessageTooLate, ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		return a.msg.Date().Before(a.session.Start), ""
	},
	Variables: variableMap{
		"SESSIONSTART": func(a *Analysis) string {
			return a.session.Start.Format("2006-01-02 at 15:04")
		},
	},
	references: refWeeklyPractice,
}

// ProbMessageTooLate is raised when
var ProbMessageTooLate = &Problem{
	Code:  "MessageTooLate",
	after: []*Problem{ProbPracticeSubjectFormat},
	ifnot: []*Problem{ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		return a.msg.Date().Before(a.session.Start) && !a.subjectDate.IsZero() && a.subjectDate.Before(a.session.Start), ""
	},
	Variables: variableMap{
		"SUBJECTDATE": func(a *Analysis) string {
			return a.subjectDate.Format("January 2")
		},
	},
	references: refWeeklyPractice,
}

// ProbSessionDate is raised when
var ProbSessionDate = &Problem{
	Code:  "SessionDate",
	after: []*Problem{ProbPracticeSubjectFormat},
	ifnot: []*Problem{ProbMessageTooEarly, ProbMessageTooLate, ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		return !a.subjectDate.IsZero() &&
				(a.subjectDate.Year() != a.session.End.Year() ||
					a.subjectDate.Month() != a.session.End.Month() ||
					a.subjectDate.Day() != a.session.End.Day()),
			""
	},
	Variables: variableMap{
		"SUBJECTDATE": func(a *Analysis) string {
			return a.subjectDate.Format("January 2")
		},
	},
	references: refWeeklyPractice,
}
