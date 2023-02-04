package analyze

import "github.com/rothskeller/packet/wppsvr/english"

func init() {
	ProblemLabels["ToBBS"] = "message to incorrect BBS"
	ProblemLabels["ToBBSDown"] = "message to incorrect BBS (simulated outage)"
}

// checkReceiver checks for problems with who and where the message was sent to.
func (a *Analysis) checkReceiver() {
	// Check whether the message was sent to a simulated down BBS.
	if inList(a.session.DownBBSes, a.toBBS) {
		a.reportProblem("ToBBSDown", refWeeklyPractice, toBBSDownResponse,
			a.session.CallSign, a.toBBS, a.session.Name, a.session.End.Format("January 2"),
			english.Conjoin(a.session.ToBBSes, "or"))
		return
	}
	// Check whether the message was sent to a valid destination BBS.
	if !inList(a.session.ToBBSes, a.toBBS) {
		a.reportProblem("ToBBS", refWeeklyPractice, toBBSResponse,
			a.session.CallSign, a.toBBS, a.session.Name, a.session.End.Format("January 2"),
			english.Conjoin(a.session.ToBBSes, "or"))
	}
}

const toBBSResponse = `This message was sent to %[1]s at %[2]s.  Practice
messages for %[3]s on %[4]s must be sent to %[1]s at %[5]s.`
const toBBSDownResponse = `This message was sent to %[1]s at %[2]s, but %[2]s
has a simulated outage for %[3]s on %[4]s.  Practice messages for this session
must be sent to %[1]s at %[5]s.`
