package analyze

import (
	"fmt"

	"steve.rothskeller.net/packet/wppsvr/english"
)

// Problem codes
const (
	ProblemFromBBSDown = "FromBBSDown"
	ProblemToBBSDown   = "ToBBSDown"
	ProblemToBBS       = "ToBBS"
)

func init() {
	ProblemLabel[ProblemToBBSDown] = "message to incorrect BBS (simulated down)"
	ProblemLabel[ProblemToBBS] = "message to incorrect BBS"
	ProblemLabel[ProblemFromBBSDown] = "message from incorrect BBS"
}

// checkBBS verifies that the message was sent from and to a valid BBS.
func (a *Analysis) checkBBS() {
	var (
		found bool
		msg   = a.msg.Message()
	)
	// This check doesn't apply to non-human messages.
	if msg == nil {
		return
	}
	for _, to := range a.session.ToBBSes {
		if to == a.toBBS {
			found = true
			break
		}
	}
	if !found {
		for _, down := range a.session.DownBBSes {
			if down == a.toBBS {
				found = true
				break
			}
		}
		if found {
			a.problems = append(a.problems, &problem{
				code: ProblemToBBSDown,
				response: fmt.Sprintf(`
This message was sent to %s at %s, but %s is simulated down for %s on %s.
Practice messages for this session must be sent to %s at %s.
`, a.session.CallSign, a.toBBS, a.toBBS, a.session.Name, a.session.End.Format("January 2"), a.session.CallSign,
					english.Conjoin(a.session.ToBBSes, "or")),
				references: refWeeklyPractice,
			})
		} else {
			a.problems = append(a.problems, &problem{
				code: ProblemToBBS,
				response: fmt.Sprintf(`
This message was sent to %s at %s.  Practice messages for %s on %s must be
sent to %s at %s.
`, a.session.CallSign, a.toBBS, a.session.Name, a.session.End.Format("January 2"), a.session.CallSign,
					english.Conjoin(a.session.ToBBSes, "or")),
				references: refWeeklyPractice,
			})
		}
	}
	for _, down := range a.session.DownBBSes {
		if down == msg.FromBBS {
			a.problems = append(a.problems, &problem{
				code: ProblemFromBBSDown,
				response: fmt.Sprintf(`
This message was sent from %s, which is simulated down for %s on %s.  Practice
messages should not be sent from BBSes that are simulated down.
`, msg.FromBBS, a.session.Name, a.session.End.Format("January 2")),
				references: refWeeklyPractice,
			})
		}
	}
}
