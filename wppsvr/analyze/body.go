package analyze

import (
	"strconv"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/english"
	"github.com/rothskeller/packet/xscmsg"
)

func init() {
	ProblemLabels["FormCorrupt"] = "incorrectly encoded form"
	ProblemLabels["FormDestination"] = "incorrect destination for form"
	ProblemLabels["FormHandlingOrder"] = "incorrect handling order for form"
	ProblemLabels["FormInvalid"] = "invalid form contents"
	ProblemLabels["FormToICSPosition"] = `incorrect "To ICS Position" for form`
	ProblemLabels["FormToLocation"] = `incorrect "To Location" for form`
	ProblemLabels["FormVersion"] = "form version out of date"
	ProblemLabels["MessageFromWinlink"] = "message sent from Winlink"
	ProblemLabels["MessageNotASCII"] = "message has non-ASCII characters"
	ProblemLabels["MessageNotPlainText"] = "not a plain text message"
	ProblemLabels["MessageTypeWrong"] = "incorrect message type"
	ProblemLabels["PIFOVersion"] = "PackItForms version out of date"
}

// checkBody checks for problems with the message body.
func (a *Analysis) checkBody() {
	a.checkBodyFormat()
	a.checkMessageType()
	a.checkForm()
}

// checkBodyFormat checks for problems with the formatting of the message body.
func (a *Analysis) checkBodyFormat() {
	// Check for a message body that is not in text/plain format or has a
	// non-identity transfer encoding.
	if strings.Contains(a.msg.ReturnAddress(), "winlink.org") &&
		a.msg.Header.Get("Content-Transfer-Encoding") == "quoted-printable" {
		// Messages from winlink.org with quoted-printable encoding get
		// a special error message.
		a.reportProblem("MessageFromWinlink", 0, messageFromWinlinkResponse)
	} else if a.msg.Flags&pktmsg.NotPlainText != 0 {
		a.reportProblem("MessageNotPlainText", 0, messageNotPlainTextResponse)
	}
	// Check for non-ASCII characters in the message body.
	if strings.IndexFunc(a.msg.Body, nonASCII) >= 0 {
		a.reportProblem("MessageNotASCII", 0, messageNotASCIIResponse)
	}
}
func nonASCII(r rune) bool {
	return r > 126 || (r < 32 && r != '\t' && r != '\n')
}

// checkMessageType checks to be sure the message type matches expectations.
func (a *Analysis) checkMessageType() {
	if a.xsc.Type.Tag == xscmsg.PlainTextTag {
		// If the body looks like it contains a form, but the message
		// was parsed as plain text, it means the form encoding is
		// corrupt.  Report that.
		if pktmsg.IsForm(a.msg.Body) {
			a.reportProblem("FormCorrupt", 0, formCorruptResponse)
			return
		}
		// If the message came from outside the county BBS system, plain
		// text messages are always OK and we don't go on to check the
		// message type.
		if config.Get().BBSes[a.msg.FromBBS()] == nil {
			return
		}
	}
	// Make sure the message has a type expected for the current session.
	if !inList(a.session.MessageTypes, a.xsc.Type.Tag) {
		var (
			allowed []string
			article string
		)
		for i, code := range a.session.MessageTypes {
			mtype := xscmsg.RegisteredTypes[code]
			allowed = append(allowed, mtype.Name)
			if i == 0 {
				article = mtype.Article
			}
		}
		a.reportProblem("MessageTypeWrong", refWeeklyPractice, messageTypeWrongResponse,
			a.xsc.Type.Article, a.xsc.Type.Name, a.session.Name, a.session.End.Format("January 2"),
			article, english.Conjoin(allowed, "or"))
	}
}

// checkForm checks for problems with the form embedded in the message body.
func (a *Analysis) checkForm() {
	if a.xsc.Type.Tag == xscmsg.PlainTextTag {
		return // this function only applies to forms
	}
	// Check the PIFO version.
	minPIFO := config.Get().MinPIFOVersion
	if xscmsg.OlderVersion(a.xsc.RawForm.PIFOVersion, minPIFO) {
		a.reportProblem("PIFOVersion", 0, pifoVersionResponse, a.xsc.RawForm.PIFOVersion, minPIFO)
	}
	// Check the validity of the form contents.
	if problems := a.xsc.Validate(true); len(problems) != 0 {
		a.reportProblem("FormInvalid", 0, formInvalidResponse1+strings.Join(problems, "\n    ")+formInvalidResponse2)
	}
	// Get the message type data.
	mtc := config.Get().MessageTypes[a.xsc.Type.Tag]
	if mtc == nil {
		return // no message type data, must be unknown form type
	}
	// Check the form version.
	if mtc.MinimumVersion != "" {
		if xscmsg.OlderVersion(a.xsc.RawForm.FormVersion, mtc.MinimumVersion) {
			a.reportProblem("FormVersion", 0, formVersionResponse, a.xsc.RawForm.FormVersion, a.xsc.Type.Name, mtc.MinimumVersion)
		}
	}
	badpos := len(mtc.ToICSPosition) != 0 && !inList(mtc.ToICSPosition, a.xsc.KeyField(xscmsg.FToICSPosition).Value)
	badloc := len(mtc.ToLocation) != 0 && !inList(mtc.ToLocation, a.xsc.KeyField(xscmsg.FToLocation).Value)
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
		a.reportProblem("FormDestination", refRouting, formDestinationResponse,
			a.xsc.KeyField(xscmsg.FToICSPosition).Value, a.xsc.KeyField(xscmsg.FToLocation).Value,
			a.xsc.Type.Name, exppos, exploc)
	} else if badpos {
		a.reportProblem("FormToICSPosition", refRouting, formToICSPositionResponse,
			a.xsc.KeyField(xscmsg.FToICSPosition).Value, a.xsc.Type.Name, exppos)
	} else if badloc {
		a.reportProblem("FormToLocation", refRouting, formToLocationResponse,
			a.xsc.KeyField(xscmsg.FToLocation).Value, a.xsc.Type.Name, exploc)
	}
	// Check the recommended handling order.
	if handling := mtc.HandlingOrder; handling != "" {
		var want xscmsg.HandlingOrder
		if mtc.HandlingOrder == "computed" {
			want = config.ComputedRecommendedHandlingOrder[a.xsc.Type.Tag](a.xsc)
		} else {
			want, _ = xscmsg.ParseHandlingOrder(mtc.HandlingOrder)
		}
		have, _ := xscmsg.ParseHandlingOrder(a.xsc.KeyField(xscmsg.FHandling).Value)
		if want != 0 && want != have {
			a.reportProblem("FormHandlingOrder", refRouting, formHandlingOrderResponse,
				a.xsc.KeyField(xscmsg.FHandling).Value, want)
		}
	}
}

const formCorruptResponse = `This message appears to contain an encoded form,
but the encoding is incorrect.  It appears to have been created or edited by
software other than the current PackItForms software.  Please use current
PackItForms software to encode messages containing forms.`
const formDestinationResponse = `This message form is addressed to ICS Position
"%s" at Location "%s".  %ss should be addressed to %s at %s.`
const formHandlingOrderResponse = `This message has handling order %s.  It
should have handling order %s.`
const formInvalidResponse1 = "This message contains a form with invalid contents:\n    "
const formInvalidResponse2 = "\nPlease verify the correctness of the form before sending."
const formToICSPositionResponse = `This message form is addressed to ICS
Position "%s".  %ss should be addressed to ICS Position %s.`
const formToLocationResponse = `This message form is addressed to Location
"%s".  %ss should be addressed to Location %s.`
const formVersionResponse = `This message contains version %s of the %s, but
that version is not current.  Please use version %s or newer of the form.  (You
can get the newer form by updating your PackItForms installation.)`
const messageFromWinlinkResponse = `This message was sent from Winlink.  Winlink
should not be used for emergency communications, unless no alternatives are
available, because it uses a message encoding system ("quoted-printable") that
Outpost cannot decode.  As a result, some messages (particularly those with long
lines and those containing equals signs) may be garbled in transmission.`
const messageNotASCIIResponse = `This message contains characters that are not
in the standard ASCII character set (i.e., not on a standard keyboard).
Non-standard characters should be avoided in packet messages, because the
receiving system may not know how to render them.  Note that some software may
introduce undesired non-standard characters (e.g., Microsoft Word's "smart
quotes" feature). If you use message text composed in such software, make sure
those features are disabled.`
const messageNotPlainTextResponse = `This message is not a plain text message.
All SCCo packet messages should be plain text only.  ("Rich text" or
HTML-formatted messages, common in email systems, are far larger than plain text
messages and put too much strain on the packet infrastructure.)  Please
configure your software to send plain text messages when sending to an SCCo
BBS.`
const messageTypeWrongResponse = `This message is %s %s.  For the %s on %s, %s
%s is expected.`
const pifoVersionResponse = `This message used version %s of PackItForms to
encode the form, but that version is not current.  Please use PackItForms
version %s or newer to encode messages containing forms.`
