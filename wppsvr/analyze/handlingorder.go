package analyze

import (
	"fmt"

	"steve.rothskeller.net/packet/pktmsg"
)

// Problem codes
const (
	ProblemFormHandlingOrder = "FormHandlingOrder"
)

func init() {
	ProblemLabel[ProblemFormHandlingOrder] = "incorrect handling order for form"
}

// checkFormHandlingOrder determines whether the message has the correct
// handling order based on the form contents.
func (a *Analysis) checkFormHandlingOrder() {
	// The correct handling order depends on the form type, and we only
	// check it for the form types we know.
	switch msg := a.msg.(type) {
	case *pktmsg.RxICS213Form:
		var shouldbe pktmsg.HandlingOrder
		// According to the recommended form routing, the handling order
		// depends on the severity.
		switch msg.Severity {
		case pktmsg.SeverityEmergency:
			shouldbe = pktmsg.HandlingImmediate
		case pktmsg.SeverityUrgent:
			shouldbe = pktmsg.HandlingPriority
		case pktmsg.SeverityOther:
			shouldbe = pktmsg.HandlingRoutine
		}
		if shouldbe != 0 && msg.HandlingOrder != shouldbe {
			a.problems = append(a.problems, &problem{
				code:    ProblemFormHandlingOrder,
				subject: "Incorrect handling order for form",
				response: fmt.Sprintf(`
This message has severity %s and handling order %s.  ICS-213 messages with
severity %s should have handling order %s.
`, msg.Severity, msg.HandlingOrder, msg.Severity, shouldbe),
				references: refFormRouting,
			})
		}
	case *pktmsg.RxEOC213RRForm:
		var shouldbe pktmsg.HandlingOrder
		// According to the recommended form routing, the handling order
		// depends on the priority.
		switch msg.Priority {
		case "Now", "High":
			shouldbe = pktmsg.HandlingImmediate
		case "Medium":
			shouldbe = pktmsg.HandlingPriority
		case "Low":
			shouldbe = pktmsg.HandlingRoutine
		}
		if shouldbe != 0 && msg.HandlingOrder != shouldbe {
			a.problems = append(a.problems, &problem{
				code:    ProblemFormHandlingOrder,
				subject: "Incorrect handling order for form",
				response: fmt.Sprintf(`
This message has priority %s and handling order %s.  EOC-213RR messages with
priority %s should have handling order %s.
`, msg.Priority, msg.HandlingOrder, msg.Priority, shouldbe),
				references: refFormRouting,
			})
		}
	case *pktmsg.RxMuniStatForm:
		if msg.HandlingOrder != pktmsg.HandlingImmediate {
			a.problems = append(a.problems, &problem{
				code:    ProblemFormHandlingOrder,
				subject: "Incorrect handling order for form",
				response: fmt.Sprintf(`
This message handling order %s.  OA Municipal Status messages should have
handling order IMMEDIATE.
`, msg.HandlingOrder),
				references: refFormRouting,
			})
		}
	}
}
