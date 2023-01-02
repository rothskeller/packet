package analyze

func init() {
	Problems[ProbBounceMessage.Code] = ProbBounceMessage
	Problems[ProbMessageCorrupt.Code] = ProbMessageCorrupt
}

// ProbBounceMessage is raised for messages that appear to be auto-responses.
var ProbBounceMessage = &Problem{
	Code: "BounceMessage",
	// detection is hard-coded in Analyze
}

// ProbMessageCorrupt is raised for messages that can't be parsed as valid email
// messages.
var ProbMessageCorrupt = &Problem{
	Code:   "MessageCorrupt",
	detect: func(*Analysis) bool { return false },
	// detection is hard-coded in Analyze
}
