package analyze

// This file contains the problem checks that are run only against known
// PackItForms forms messages.  They are run after the common checks.  They
// appear in the order they are run, although some are skipped based on the
// message type or the results of previous checks.

import (
	"strconv"
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/english"
)

func init() {
	ProblemLabels["FormNoCallSign"] = "no call sign in form message"
	ProblemLabels["CallSignConflict"] = "call sign conflict"
	ProblemLabels["FormSubject"] = "message subject doesn't agree with form contents"
	ProblemLabels["FormInvalid"] = "invalid form contents"
	ProblemLabels["PIFOVersion"] = "PackItForms version out of date"
	ProblemLabels["FormVersion"] = "form version out of date"
	ProblemLabels["FormDestination"] = "incorrect destination for form"
	ProblemLabels["FormToICSPosition"] = `incorrect "To ICS Position" for form`
	ProblemLabels["FormToLocation"] = `incorrect "To Location" for form`
	ProblemLabels["FormHandlingOrder"] = "incorrect handling order for form"
}

func (a *Analysis) formNoCallSign() bool {
	if a.FromCallSign == "" {
		return a.reportProblem("FormNoCallSign", 0, formNoCallSignResponse)
	}
	return false
}

const formNoCallSignResponse = `This message cannot be counted because it's not
clear who sent it.  There is no call sign in the return address, or after the
word "Practice" on the subject line, or in the Operator Call field of the form.
In order for a message to count, there must be a call sign in at least one of
those places.`

func (a *Analysis) callSignConflict() bool {
	if fccCallSignRE.MatchString(a.FromCallSign) {
		if a.key.OpCall != "" && a.key.OpCall != a.FromCallSign {
			return a.reportProblem("CallSignConflict", 0, callSignConflictResponse, a.FromCallSign, a.key.OpCall)
		}
	}
	return false
}

const callSignConflictResponse = `This message has conflicting call signs.  The
Subject line says the call sign is %s, but the Operator Call Sign field of the
form says %s.  The two should agree.  (This message will be counted as a
practice attempt by %[1]s.)`

func (a *Analysis) formSubject() bool {
	subject := a.msg.(message.IEncode).EncodeSubject()
	if a.subject != subject {
		return a.reportProblem("FormSubject", 0, formSubjectResponse, a.subject, subject)
	}
	return false
}

const formSubjectResponse = `This message has
    Subject: %s
but, based on the contents of the form, it should have
    Subject: %s
PackItForms automatically generates the Subject line from the form contents; it
should not be overridden manually.`

func (a *Analysis) formInvalid() bool {
	if f, ok := a.msg.(message.IValidate); ok {
		if problems := f.Validate(); len(problems) != 0 {
			return a.reportProblem("FormInvalid", 0, formInvalidResponse1+strings.Join(problems, "\n    "+formInvalidResponse2))
		}
	}
	return false
}

const formInvalidResponse1 = "This message contains a form with invalid contents:\n    "
const formInvalidResponse2 = "\nPlease verify the correctness of the form before sending."

func (a *Analysis) pifoVersion() bool {
	min := config.Get().MinPIFOVersion
	have := a.key.PIFOVersion
	if message.OlderVersion(have, min) {
		return a.reportProblem("PIFOVersion", 0, pifoVersionResponse, have, min)
	}
	return false
}

const pifoVersionResponse = `This message used version %s of PackItForms to
encode the form, but that version is not current.  Please use PackItForms
version %s or newer to encode messages containing forms.`

func (a *Analysis) formVersion() bool {
	min := config.Get().MessageTypes[a.msg.Type().Tag].MinimumVersion
	have := a.key.FormVersion
	if message.OlderVersion(have, min) {
		return a.reportProblem("FormVersion", 0, formVersionResponse, have, a.msg.Type().Name, min)
	}
	return false
}

const formVersionResponse = `This message contains version %s of the %s, but
that version is not current.  Please use version %s or newer of the form.  (You
can get the newer form by updating your PackItForms installation.)`

func (a *Analysis) formDestination() bool {
	mtc := config.Get().MessageTypes[a.msg.Type().Tag]
	badpos := len(mtc.ToICSPosition) != 0 && !inList(mtc.ToICSPosition, a.key.ToICSPosition)
	badloc := len(mtc.ToLocation) != 0 && !inList(mtc.ToLocation, a.key.ToLocation)
	var exppos, exploc string
	if badpos {
		var positions []string
		for _, pos := range mtc.ToICSPosition {
			positions = append(positions, strconv.Quote(pos))
		}
		exppos = english.Conjoin(positions, "or")
	}
	if badloc {
		var locations []string
		for _, loc := range mtc.ToLocation {
			locations = append(locations, strconv.Quote(loc))
		}
		exploc = english.Conjoin(locations, "or")
	}
	if badpos && badloc {
		return a.reportProblem("FormDestination", refRouting, formDestinationResponse, a.key.ToICSPosition, a.key.ToLocation,
			a.msg.Type().Name, exppos, exploc)
	}
	if badpos {
		return a.reportProblem("FormToICSPosition", refRouting, formToICSPositionResponse,
			a.key.ToICSPosition, a.msg.Type().Name, exppos)
	}
	if badloc {
		return a.reportProblem("FormToLocation", refRouting, formToLocationResponse,
			a.key.ToLocation, a.msg.Type().Name, exploc)
	}
	return false
}

const formDestinationResponse = `This message form is addressed to ICS Position
"%s" at Location "%s".  %ss should be addressed to %s at %s.`
const formToICSPositionResponse = `This message form is addressed to ICS
Position "%s".  %ss should be addressed to ICS Position %s.`
const formToLocationResponse = `This message form is addressed to Location
"%s".  %ss should be addressed to Location %s.`

func (a *Analysis) formHandlingOrder() bool {
	want := config.Get().MessageTypes[a.msg.Type().Tag].HandlingOrder
	if want == "computed" {
		want = config.ComputeRecommendedHandlingOrder(a.msg)
	}
	have := a.key.Handling
	if want != "" && want != have {
		return a.reportProblem("FormHandlingOrder", refRouting, formHandlingOrderResponse, have, want)
	}
	return false
}

const formHandlingOrderResponse = `This message has handling order %s.  It
should have handling order %s.`
