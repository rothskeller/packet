package analyze

import (
	"regexp"
	"strings"

	"github.com/rothskeller/packet/xscmsg"
)

func init() {
	ProblemLabels["CallSignConflict"] = "call sign conflict"
	ProblemLabels["FormNoCallSign"] = "no call sign in form message"
	ProblemLabels["FromBBSDown"] = "message from incorrect BBS (simulated outage)"
	ProblemLabels["NoCallSign"] = "no call sign in message"
}

// checkSender checks for problems with who and where the message was sent from.
// This has the side effect of setting a.FromCallSign.
func (a *Analysis) checkSender() {
	// First, check whether the message is from a simulated down BBS.
	if inList(a.session.DownBBSes, a.msg.FromBBS()) {
		a.reportProblem("FromBBSDown", refWeeklyPractice, fromBBSDownResponse,
			a.msg.FromBBS(), a.session.Name, a.session.End.Format("January 2"))
	}
	// Next, check whether we have a call sign.
	if a.FromCallSign = a.fromCallSign(); a.FromCallSign == "" {
		// We don't.  Which message we raise depends on whether the
		// message is a form with an OpCall field.
		if a.xsc.KeyField(xscmsg.FOpCall) != nil {
			a.reportProblem("FormNoCallSign", 0, formNoCallSignResponse)
		} else {
			a.reportProblem("NoCallSign", 0, noCallSignResponse)
		}
		return
	}
	// We have a call sign.  If it is an FCC call sign, and this is a form
	// with an OpCall field, and it doesn't match the call sign in that
	// field, raise a conflict error.
	if fccCallSignRE.MatchString(a.FromCallSign) {
		if f := a.xsc.KeyField(xscmsg.FOpCall); f != nil {
			if strings.ToUpper(f.Value) != a.FromCallSign {
				a.reportProblem("CallSignConflict", 0, callSignConflictResponse,
					a.FromCallSign, strings.ToUpper(f.Value))
			}
		}
	}
}

// practiceCallSignRE is a permissive regexp that looks for a call sign in the
// "Practice ..." portion of a subject line or field.  It matches anything that
// looks like a call sign (FCC or tactical) in the word after "Practice".  The
// regexp returns the matched call sign, which may be in any case.
var practiceCallSignRE = regexp.MustCompile(`(?i)(?:\b|_)Practice[\W_]+([AKNW][A-Z]?[0-9][A-Z]{1,3}|[A-Z][A-Z0-9]{5})\b`)

// fromCallSign extracts the sender's call sign from the several different
// places it might appear.
func (a *Analysis) fromCallSign() string {
	// If there's a call sign on the subject line, use that.
	if match := practiceCallSignRE.FindStringSubmatch(a.msg.Header.Get("Subject")); match != nil {
		return strings.ToUpper(match[1])
	}
	// If this is a form with a Subject field and there's a call sign in it,
	// use that.
	if sf := a.xsc.KeyField(xscmsg.FSubject); sf != nil {
		if match := practiceCallSignRE.FindStringSubmatch(sf.Value); match != nil {
			return strings.ToUpper(match[1])
		}
	}
	// If there's a call sign in the From address of the message, use that.
	if cs := a.msg.FromCallSign(); cs != "" {
		// If it's a valid FCC call sign, use it.
		if fccCallSignRE.MatchString(cs) {
			return cs
		}
		// It's not a valid FCC call sign.  If it's from within our BBS
		// network, assume it's a tactical call sign and use it.
		// If it's from within our BBS network, use it.
		if a.msg.FromBBS() != "" {
			return cs
		}
		// It's not a valid FCC call sign and it's not from within our
		// BBS network, so it doesn't count.
	}
	// If the message is a form and has a call sign in an OpCall field, use
	// that.
	if f := a.xsc.KeyField(xscmsg.FOpCall); f != nil && f.Value != "" {
		return strings.ToUpper(f.Value)
	}
	// No call sign found.
	return ""
}

const callSignConflictResponse = `This message has conflicting call signs.  The
Subject line says the call sign is %s, but the Operator Call Sign field of the
form says %s.  The two should agree.  (This message will be counted as a
practice attempt by %[1]s.)`
const formNoCallSignResponse = `This message cannot be counted because it's not
clear who sent it.  There is no call sign in the return address, or after the
word "Practice" on the subject line, or in the Operator Call field of the form.
In order for a message to count, there must be a call sign in at least one of
those places.`
const fromBBSDownResponse = `This message was sent from %s, which has a
simulated outage for %s on %s.  Practice messages should not be sent from BBSes
that have a simulated outage.`
const noCallSignResponse = `This message cannot be counted because it's not
clear who sent it.  There is no call sign in the return address or after the
word "Practice" on the subject line.  In order for a message to count, there
must be a call sign in at least one of those places.`
