package analyze

import "github.com/rothskeller/packet/pktmsg"

// ProbBounceMessage is raised for messages that appear to be auto-responses.
var ProbBounceMessage = &Problem{
	Code: "BounceMessage",
	detect: func(a *Analysis) (bool, string) {
		return a.msg.Flags&pktmsg.AutoResponse != 0, ""
	},
}

// ProbMessageCorrupt is raised for messages that can't be parsed as valid email
// messages.
var ProbMessageCorrupt = &Problem{
	Code:   "MessageCorrupt",
	detect: func(*Analysis) (bool, string) { return false, "" },
	// detection is hard-coded in checkMessage
}

func init() {
	Problems[ProbBounceMessage.Code] = ProbBounceMessage
	Problems[ProbMessageCorrupt.Code] = ProbMessageCorrupt
}
