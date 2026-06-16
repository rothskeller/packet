package analyze

// This file contains the problem checks that are run against all human
// messages.  They appear in the order they are run, although some are skipped
// based on the message type or the results of previous checks.

import (
	"fmt"
	"html"
	"regexp"
	"slices"
	"strings"

	"github.com/rothskeller/packet/v4/errors"
	"github.com/rothskeller/packet/v4/form"
	"github.com/rothskeller/packet/v4/message"
	"github.com/rothskeller/packet/v4/message/field"
	"github.com/rothskeller/packet/v4/message/msgifc"
	"github.com/rothskeller/packet/v4/message/receipt"
	"github.com/rothskeller/packet/v4/message/subject"
	"github.com/rothskeller/packet/v4/wppsvr/config"
	"github.com/rothskeller/packet/v4/wppsvr/english"
	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	// fromCallSignRE extracts the fromCallSign from the return address.  It
	// looks for a call sign at the start of the string, followed either by
	// a %, an @, or the end of the string.  It is not case-sensitive.  The
	// substring returned is the call sign.
	fromCallSignRE = regexp.MustCompile(`(?i)^(A[A-L][0-9][A-Z]{1,3}|[KNW][A-Z][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3}|[A-Z][A-Z0-9]{5})(?:@|%|$)`)
	// fccCallSignRE matches a legal FCC call sign.  It is not
	// case-sensitive.
	fccCallSignRE = regexp.MustCompile(`(?i)^(?:A[A-L][0-9][A-Z]{1,3}|[KNW][A-Z][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3})$`)
	// fromBBSRE matches a return address from a BBS, and returns the BBS
	// name.  It is the first word of the address domain, as long as that
	// address looks like a call sign and the rest of the domain is
	// ".ampr.org", ".scc-ares-races.org", or a ".#" BBS network domain.
	fromBBSRE = regexp.MustCompile(`(?i)^[^%@]+[%@](A[A-L][0-9][A-Z]{1,3}|[KNW][A-Z][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3})(?:\.(?:ampr|scc-ares-races)\.org(?:@.*)?|\.#.*)?$`)
	// msgnumRE matches a valid packet message number.
	msgnumRE = regexp.MustCompile(`^(?:[A-Z][A-Z][A-Z]|[A-Z][0-9][A-Z0-9]|[0-9][A-Z][A-Z])-\d\d\d+[AC-HJ-NPR-Y]$`)
)

// messageCounts returns whether the message should be counted as a check-in.
// If not, summary and analysis are set appropriately.
func (a *Analysis) messageCounts(parseErr error) bool {
	var opcallFound bool

	// First, fail immediately if the message is not a human message.
	if parseErr != nil {
		a.score, a.outOf = 0, 1
		a.setSummary("message could not be parsed")
		fmt.Fprintf(a.analysis, "<h2>Message Could Not Be Parsed</h2><p>This message could not be parsed as a valid RFC-4155 or RFC-5322 message.  The parse error is “<tt>%s</tt>”.</p>",
			html.EscapeString(parseErr.Error()))
		return false
	}
	if a.msg.Autoresponse() {
		a.score, a.outOf = 0, 1
		a.setSummary("message has no return address (probably auto-response)")
		a.analysis.WriteString("<h2>Message Has No Return Address</h2><p>This message has no return address, which normally means that it is an auto-response message (e.g., an out-of-office response or a bounce message).  It will not be counted.</p>")
		return false
	}
	if a.msg.Type() == receipt.DeliveryReceipt {
		// No summary or analysis for delivery receipt.
		a.score, a.outOf = 0, 1
		return false
	}
	if a.msg.Type() == receipt.ReadReceipt {
		a.score, a.outOf = 0, 1
		a.setSummary("unexpected READ receipt message")
		a.analysis.WriteString(`<h2>Unexpected READ Receipt Message</h2><p>This message is an Outpost “read receipt,” which should not have been sent.  Most likely, your Outpost installation has the “Auto-Read Receipt” setting turned on.  The SCCo “Standard Outpost Configuration Instructions” (available on the <a href="https://www.scc-ares-races.org/services/data/bbs">“Packet BBS Service” page</a> of the county ARES website) specifies that this setting should be turned off.  You can find it on the Receipts tab of the Message Settings dialog in Outpost.</p>`)
		return false
	}
	// Check that it was sent to a correct BBS.
	if slices.Contains(a.session.DownBBSes, a.sm.ToBBS) {
		a.setSummary("message to incorrect BBS (simulated outage)")
		fmt.Fprintf(a.analysis, "<h2>Message to Incorrect BBS</h2><p>This message was sent to %[1]s at %[2]s, but %[2]s has a simulated outage for %[3]s on %[4]s.  This message will not be counted.  Practice messages for this session must be sent to %[1]s at %[5]s.</p>",
			a.session.CallSign, a.sm.ToBBS, html.EscapeString(a.session.Name), a.session.End.Format("January 2"),
			english.Conjoin(a.session.ToBBSes, "or"))
	} else if !slices.Contains(a.session.ToBBSes, a.sm.ToBBS) {
		a.setSummary("message to incorrect BBS")
		fmt.Fprintf(a.analysis, "<h2>Message to Incorrect BBS</h2><p>This message was sent to %[1]s at %[2]s, but practice messages for %[3]s on %[4]s must be sent to %[1]s at %[5]s.  This message will not be counted.</p>",
			a.session.CallSign, a.sm.ToBBS, html.EscapeString(a.session.Name), a.session.End.Format("January 2"),
			english.Conjoin(a.session.ToBBSes, "or"))
	}
	// Check that it was sent after the start of the session.
	rcvdate := a.msg.BBSRxDate()
	if rcvdate.IsZero() {
		rcvdate = a.msg.Date()
	}
	if rcvdate.Before(a.session.Start) {
		a.setSummary("message sent outside of practice session")
		fmt.Fprintf(a.analysis, "<h2>Message Sent Outside of Practice Session</h2><p>This message arrived at %s on %s.  However, practice messages for %s aren’t accepted until %s.  This message will not be counted.</p>",
			a.sm.ToBBS, rcvdate.Format("2006-01-02 at 15:04"), html.EscapeString(a.session.Name),
			a.session.Start.Format("2006-01-02 at 15:04"))
	}
	// Check that we have a call sign so that we know whom to credit.  To do
	// that, we need to know what BBS the message came from, if any.
	if match := fromBBSRE.FindStringSubmatch(a.msg.ReturnAddr()); match != nil {
		a.sm.FromBBS = strings.ToUpper(match[1])
	}
	if match := fromCallSignRE.FindStringSubmatch(a.msg.ReturnAddr()); match != nil && (fccCallSignRE.MatchString(match[1]) || a.sm.FromBBS != "") {
		// We'll take an FCC call sign in the return address as the
		// call sign to credit.  A tactical call sign in the return
		// address counts only if the message is coming from a BBS;
		// otherwise we can't be sure it's a tactical call sign.
		a.sm.FromCallSign = strings.ToUpper(match[1])
	} else {
		// Look for an OpCall field.
		for f := range a.msg.Fields() {
			if f.Common() == field.COperatorCall {
				opcallFound = true
				a.sm.FromCallSign = f.Value(a.msg)
				break
			}
		}
	}
	if a.sm.FromCallSign == "N6SBC" {
		// Special case request from Timothy Takeuchi, 2023-09.
		a.sm.FromCallSign = "XBEEOC"
		a.sm.Jurisdiction = "XBE"
	}
	if a.sm.FromCallSign == "" {
		a.setSummary("no call sign in message")
		if opcallFound {
			a.analysis.WriteString(`<h2>No Call Sign in Message</h2><p>This message cannot be counted because it’s not clear who sent it.  There is no call sign in the return address or in the Operator Call field of the form.  In order for the message to count, there must be a call sign in at least one of those places.</p>`)
		} else {
			a.analysis.WriteString(`<h2>No Call Sign in Message</h2><p>This message cannot be counted because it’s not clear who sent it.  There is no call sign in the return address.  In order the message to count, it must come from a BBS mailbox or email account whose name is a call sign.`)
		}
	}
	// If any problems have been reported to this point, the message can't
	// be counted.
	if a.sm.Summary != "" {
		a.score, a.outOf = 0, 1
		return false
	}
	a.score, a.outOf = 1, 1
	return true
}

// checkCorrectness verifies that the message is properly encoded and valid.
// These checks are run for all messages, whether or not we have a model message
// to compare against.  Any problems are added to the analysis.
func (a *Analysis) checkCorrectness() {
	// Make sure the message is plain text.
	a.outOf++
	if a.msg.IsMultipart() {
		a.setSummary("not a plain text message")
		a.analysis.WriteString("<h2>Not a Plain Text Message</h2><p>This message is not a plain text message. All SCCo packet messages should be plain text only.  (“Rich text” or HTML-formatted messages, common in email systems, are far larger than plain text messages and put too much strain on the packet infrastructure.)  Please configure your software to send plain text messages when sending to an SCCo BBS.</p>")
	} else {
		a.score++
	}
	// Make sure the message has only ASCII characters.
	a.outOf++
	if strings.IndexFunc(a.msg.Body().EncodedBody(), nonASCII) >= 0 {
		a.setSummary("message has non-ASCII characters")
		a.analysis.WriteString("<h2>Message Has Non-ASCII Characters</h2><p>This message contains characters that are not in the standard ASCII character set (i.e., not on a standard keyboard). Non-standard characters should be avoided in packet messages, because the receiving system may not know how to render them.  Note that some software may introduce undesired non-standard characters (e.g., Microsoft Word’s “smart quotes” feature). If you use message text composed in such software, make sure those features are disabled.</p>")
	} else {
		a.score++
	}
	// Make sure the message came from a BBS that is up.
	a.outOf++
	if slices.Contains(a.session.DownBBSes, a.sm.FromBBS) {
		a.setSummary("message from incorrect BBS (simulated outage)")
		fmt.Fprintf(a.analysis, "<h2>Message from Incorrect BBS</h2><p>This message was sent from %s, which has a simulated outage for %s on %s.  Practice messages should not be sent from BBSes that have a simulated outage.</p>",
			a.sm.FromBBS, html.EscapeString(a.session.Name), a.session.End.Format("January 2"))
	} else {
		a.score++
	}
	// Some checks only apply to form messages (of known form types).
	if fb, ok := a.msg.Body().(*form.FormBody); ok {
		var omi, handling, summary string
		// Get the field values of interest.
		for f := range a.msg.Fields() {
			switch f.Common() {
			case field.COriginMessageID:
				omi = f.Value(a.msg)
			case field.CHandling:
				handling = f.Value(a.msg)
			case field.CMessageSummary:
				summary = f.Value(a.msg)
			}
		}
		// Make sure the message subject matches the form.
		a.outOf++
		subj, _ := form.NewFormSubject(omi, handling, a.msg.Type().Tag(), summary)
		if a.msg.Subject().EncodedSubject() != subj.EncodedSubject() {
			a.setSummary("message subject doesn't agree with form contents")
			fmt.Fprintf(a.analysis, `<h2>Message Subject Doesn’t Agree with Form Contents</h2><p style="margin-bottom:0">This message has</p><div style="margin-left:2rem"><tt>Subject: %s</tt></div><div>but, based on the contents of the form, it should have</div><div style="margin-left:2rem"><tt>Subject: %s</tt></div><p style="margin-top:0">PackItForms automatically generates the Subject line from the form contents; it should not be overridden manually.</p>`,
				html.EscapeString(a.msg.Subject().EncodedSubject()), html.EscapeString(subj.EncodedSubject()))
		} else {
			a.score++
		}
		// Make sure the message is valid according to PackItForms' rules.
		if err := message.ValidateMessage(a.msg, msgifc.VPIFOOnly|msgifc.VPacket); err != nil {
			problems := errors.UnwrapJoined(err)
			a.outOf += len(problems)
			a.setSummary("invalid form contents")
			a.analysis.WriteString(`<h2>Invalid Form Contents</h2><p style="margin-bottom:0">This message contains a form with invalid contents:</p><ul style="margin-top:0;margin-bottom:0">`)
			for _, problem := range problems {
				fmt.Fprintf(a.analysis, "<li>%s</li>", html.EscapeString(problem.Error()))
			}
			a.analysis.WriteString(`</ul><p style="margin-top:0">Please verify the correctness of the form before sending.</p>`)
		}
		// Make sure the PIFO and form versions are up to date.
		a.outOf++
		minPIFO := config.Get().MinPIFOVersion
		if OlderVersion(fb.PIFOVersion(), minPIFO) {
			a.setSummary("PackItForms version out of date")
			fmt.Fprintf(a.analysis, "<h2>PackItForms Version Out of Date</h2><p>This message used version %s of PackItForms to encode the form, but that version is not current.  Please use PackItForms version %s or newer to encode messages containing forms.</p>",
				fb.PIFOVersion(), minPIFO)
		} else {
			a.score++
		}
		if mt := config.Get().MessageTypes[a.msg.Type().Tag()]; mt != nil {
			a.outOf++
			minForm := mt.MinimumVersion
			if OlderVersion(fb.FormVersion(), minForm) {
				a.setSummary("form version out of date")
				fmt.Fprintf(a.analysis, "<h2>Form Version Out of Date</h2><p>This message contains version %s of %s, but that version is not current.  Please use version %s or newer of the form.  (You can get the newer form by updating your PackItForms installation.)",
					fb.FormVersion(), html.EscapeString(a.msg.Type().Name()), minForm)
			} else {
				a.score++
			}
		}
		// Make sure the form didn't have any spurious fields.
		a.outOf++
		haveFields := sets.New(fb.FieldList()...)
		for f := range a.msg.Fields() {
			haveFields.Delete(f.Tag())
		}
		if len(haveFields) != 0 {
			a.setSummary("form has extra fields")
			if len(haveFields) == 1 {
				field, _ := haveFields.PopAny()
				fmt.Fprintf(a.analysis, "<h2>Form Has Extra Fields</h2><p>This message contains an extra field (%s) which is not expected in version %s of %s.",
					field, fb.FormVersion(), html.EscapeString(a.msg.Type().Name()))
			} else {
				fmt.Fprintf(a.analysis, "<h2>Form Has Extra Fields</h2><p>This message contains extra fields (%s) which are not expected in version %s of %s.",
					strings.Join(haveFields.UnsortedList(), ", "), fb.FormVersion(), html.EscapeString(a.msg.Type().Name()))
			}
		} else {
			a.score++
		}
		a.checkMessageNumber(omi)
	} else { // checks for plain text messages (or forms of unknown type)
		// Check the message subject format.
		a.outOf++
		if a.msg.Subject().SubjectMessageID() == "" {
			a.setSummary("incorrect subject line format")
			a.analysis.WriteString(`<h2>Incorrect Subject Line Format</h2><p>This message has an incorrect subject line format.  According to the SCCo “Standard Packet Message Subject Line” (available on the <a href="https://www.scc-ares-races.org/services/data/bbs">“Packet BBS Service” page</a> of the county ARES website), the subject line should look like <tt>AAA-111P_R_Subject</tt>, where <tt>AAA-111P</tt> is the message number, <tt>R</tt> is the handling order code, and <tt>Subject</tt> is the message subject.</p>`)
		} else {
			if _, err := subject.NewPlainSubject(a.msg.Subject().SubjectMessageID(), a.msg.Subject().SubjectHandling(), a.msg.Subject().SubjectSummary()); err != nil {
				a.setSummary("incorrect subject line format")
				fmt.Fprintf(a.analysis, `<h2>Incorrect Subject Line Format</h2><p>This message has an incorrect subject line format.  %s    According to the SCCo “Standard Packet Message Subject Line” (available on the <a href="https://www.scc-ares-races.org/services/data/bbs">“Packet BBS Service” page</a> of the county ARES website), the subject line should look like <tt>AAA-111P_R_Subject</tt>, where <tt>AAA-111P</tt> is the message number, <tt>R</tt> is the handling order code, and <tt>Subject</tt> is the message subject.</p>`, html.EscapeString(err.Error()))
			} else {
				a.score++
			}
			a.checkMessageNumber(a.msg.Subject().SubjectMessageID())
		}
		// If this is actually a plain text message (and not an unknown)
		// form type), there are a couple more things. to check.
		if a.msg.Type() == message.PlainMessage {
			a.outOf++
			body := a.msg.Body().EncodedBody()
			if strings.Contains(body, "!SCCoPIFO!") || strings.Contains(body, "!PACF!") || strings.Contains(body, "!/ADDON!") {
				a.setSummary("incorrectly encoded form")
				a.analysis.WriteString(`<h2>Incorrectly Encoded Form</h2><p>This message appears to contain an encoded form, but the encoding is incorrect.  It appears to have been created or edited by software other than the current PackItForms software.  Please use current PackItForms software to encode messages containing forms.</p>`)
			} else {
				a.score++
			}
		}
	}
}

func nonASCII(r rune) bool {
	return r > 126 || (r < 32 && r != '\t' && r != '\n')
}

// checkMessageNumber checks the validity of the message number passed to it.
// (It comes from different places in forms and non-forms messages.)
func (a *Analysis) checkMessageNumber(msgid string) {
	if msgid != "" {
		a.outOf++
		if !msgnumRE.MatchString(msgid) {
			a.setSummary("incorrect message number format")
			a.analysis.WriteString(`<h2>Incorrect Message Number Format</h2><p style="margin-bottom:0">The message number of this message is not formatted correctly.  According to the SCCo “Standard Packet Message Subject Line” document (available on the <a href="https://www.scc-ares-races.org/services/data/bbs">“Packet BBS Service” page</a> of the county ARES website), it should have a format like "XND-042P", containing:</p><ul style="margin-top:0;margin-bottom:0"><li>a three-character prefix (usually the last three characters of the sender's call sign),</li><li>a dash,</li><li>a number with at least three digits, and</li><li>a “P”, “M”, or “R” suffix.</ul><p style="margin-top:0">All letters should be upper case.  In Outpost, the format of the message number is set in the Message Settings dialog, which should be configured according to the SCCo “Standard Outpost Configuration Instructions” (available on the same page).</p>`)
		} else if fccCallSignRE.MatchString(a.sm.FromCallSign) {
			act := msgid[:3]
			exp := a.sm.FromCallSign[len(a.sm.FromCallSign)-3:]
			if act != exp {
				a.setSummary("incorrect message number prefix")
				fmt.Fprintf(a.analysis, `<h2>Incorrect Message Number Prefix</h2><p>The message number of this message has the prefix “%s”.  According to the SCCo “Standard Packet Message Subject Line” document (available on the <a href="https://www.scc-ares-races.org/services/data/bbs">“Packet BBS Service” page</a> of the county ARES website), the prefix should be the last three characters of your call sign, “%s”.</p>`,
					html.EscapeString(act), exp)
			} else {
				a.score++
			}
		} else {
			a.score++
		}
	}
}

// checkNonModel runs checks against received messages when the session does not
// have a model message to compare against.  Any problems are added to the
// analysis.
func (a *Analysis) checkNonModel() {
	// Check the recommended routing fields.
	if mtc := config.Get().MessageTypes[a.msg.Type().Tag()]; mtc != nil {
		var (
			handling string
			topos    string
			toloc    string
			badpos   bool
			badloc   bool
			exppos   string
			exploc   string
			exphand  string
		)
		for f := range a.msg.Fields() {
			switch f.Common() {
			case field.CHandling:
				handling = f.Value(a.msg)
			case field.CToICSPosition:
				topos = f.Value(a.msg)
			case field.CToLocation:
				toloc = f.Value(a.msg)
			}
		}
		if len(mtc.ToICSPosition) != 0 {
			a.outOf++
			if badpos = !slices.Contains(mtc.ToICSPosition, topos); !badpos {
				a.score++
			}
		}
		if len(mtc.ToLocation) != 0 {
			a.outOf++
			if badloc = !slices.Contains(mtc.ToLocation, toloc); !badloc {
				a.score++
			}
		}
		if badpos {
			var positions []string
			for _, pos := range mtc.ToICSPosition {
				positions = append(positions, "“"+html.EscapeString(pos)+"”")
			}
			exppos = english.Conjoin(positions, "or")
		}
		if badloc {
			var locations []string
			for _, loc := range mtc.ToLocation {
				locations = append(locations, "“"+html.EscapeString(loc)+"”")
			}
			exploc = english.Conjoin(locations, "or")
		}
		if badpos && badloc {
			a.setSummary("incorrect destination for form")
			fmt.Fprintf(a.analysis, `<h2>Incorrect Destination for Form</h2><p>This message form is addressed to ICS Position “%s” at Location “%s”.  According to the “SCCo ARES/RACES Recommended Form Routing” document (available on the <a href="https://www.scc-ares-races.org/operations/forms/go-kit">“Go Kit Forms” page</a> of the county ARES website), %s should be addressed to %s at %s.</p>`,
				html.EscapeString(topos), html.EscapeString(toloc), html.EscapeString(a.msg.Type().Name()), exppos, exploc)
		} else if badpos {
			a.setSummary(`incorrect "To ICS Position" for form`)
			fmt.Fprintf(a.analysis, `<h2>Incorrect “To ICS Position” for Form</h2><p>This message form is addressed to ICS Position “%s”.  According to the “SCCo ARES/RACES Recommended Form Routing” document (available on the <a href="https://www.scc-ares-races.org/operations/forms/go-kit">“Go Kit Forms” page</a> of the county ARES website), %s should be addressed to ICS Position %s.</p>`,
				html.EscapeString(topos), html.EscapeString(a.msg.Type().Name()), exppos)
		} else if badloc {
			a.setSummary(`incorrect "To Location" for form`)
			fmt.Fprintf(a.analysis, `<h2>Incorrect “To Location” for Form</h2><p>This message form is addressed to Location “%s”.  According to the “SCCo ARES/RACES Recommended Form Routing” document (available on the <a href="https://www.scc-ares-races.org/operations/forms/go-kit">“Go Kit Forms” page</a> of the county ARES website), %s should be addressed to Location %s.</p>`,
				html.EscapeString(toloc), html.EscapeString(a.msg.Type().Name()), exploc)
		}
		// Make sure the message has a handling order allowed by the
		// recommended routing cheat sheet.
		exphand = mtc.HandlingOrder
		if exphand == "computed" {
			exphand = config.ComputeRecommendedHandlingOrder(a.msg)
		}
		if exphand != "" {
			a.outOf++
			if exphand != "" && exphand != handling {
				a.setSummary("incorrect handling order for form")
				fmt.Fprintf(a.analysis, `<h2>Incorrect Handling Order for Form</h2><p>This message has handling order “%s”.  According to the “SCCo ARES/RACES Recommended Form Routing” document (available on the <a href="https://www.scc-ares-races.org/operations/forms/go-kit">“Go Kit Forms” page</a> of the county ARES website), it should have handling order “%s”.</p>`,
					html.EscapeString(handling), exphand)
			} else {
				a.score++
			}
		}
	}
	// Make sure the message is of a type allowed for the session.
	a.outOf++
	allowed := a.session.MessageTypes
	if config.Get().BBSes[a.sm.FromBBS] == nil {
		// Plain text messages are always OK when they come from outside
		// the county BBS system.
		allowed = append(allowed, "plain")
	}
	if a.msg.Type() == message.PlainMessage {
		body := a.msg.Body().EncodedBody()
		if strings.Contains(body, "!SCCoPIFO!") || strings.Contains(body, "!PACF!") || strings.Contains(body, "!/ADDON!") {
			// Allow a "plain text" message containing a corrupt form; that
			// problem gets reported elsewhere.
			allowed = append(allowed, "plain")
		}
	}
	if !slices.Contains(allowed, a.msg.Type().Tag()) {
		var allowed []string
		for i, code := range a.session.MessageTypes {
			var name string
			if mt := message.FindCreateTag(code); mt != nil {
				name = mt.Name()
			} else {
				name = "a " + code
			}
			if i != 0 {
				_, name, _ = strings.Cut(name, " ")
			}
			allowed = append(allowed, name)
		}
		a.setSummary("incorrect message type")
		fmt.Fprintf(a.analysis, "<h2>Incorrect Message Type</h2><p>This message is %s.  For the %s on %s, %s is expected.</p>",
			html.EscapeString(a.msg.Type().Name()), html.EscapeString(a.session.Name),
			a.session.End.Format("January 2"), english.Conjoin(allowed, "or"))
	} else {
		a.score++
	}
}

// compareAgainstModel compares the received message against the model message
// for the session.  Any problems are added to the analysis.
func (a *Analysis) compareAgainstModel() {
	// Make sure the received message is the same type as the model.
	if a.msg.Type().Tag() != a.session.ModelMsg.Type().Tag() {
		a.outOf *= 2 // Give a 50% score.
		a.setSummary("incorrect message type")
		_, tname, _ := strings.Cut(a.session.ModelMsg.Type().Name(), " ")
		fmt.Fprintf(a.analysis, "<h2>Incorrect Message Type</h2><p>This message is %s.  For the %s on %s, operators are expected to send a copy of the provided %s.</p>",
			html.EscapeString(a.msg.Type().Name()), html.EscapeString(a.session.Name),
			a.session.End.Format("January 2"), html.EscapeString(tname))
		return
	}
	// Compare the message against the model.
	score, outOf, fields := message.Compare(a.session.ModelMsg, a.msg)
	// The model may have left destination or handling blank, as an exercise
	// for the operator to look them up in the recommended routing cheat
	// sheet.  If so, we need to fix up the results of the comparison for
	// that.
	var recRouteMismatch []string
	if mtc := config.Get().MessageTypes[a.session.ModelMsg.Type().Tag()]; mtc != nil {
		score, recRouteMismatch = a.fixupRecRouteFields(score, fields, mtc)
	}
	a.score += score
	a.outOf += outOf
	if score == outOf {
		return // No need to emit the comparison.
	}
	a.setSummary("message not transcribed correctly")
	a.analysis.WriteString(`<h2>Message Not Transcribed Correctly</h2><p>There are differences between this message and the model message provided for this practice session:</p><div class="comparison"><div class="head"><div class="label">Field Name</div><div class="vmodel">Model Message</div><div class="vrecv">Received Message</div></div>`)
	for _, f := range fields {
		fmt.Fprintf(a.analysis, `<div class="field"><div class="label">%s</div><div class="vmodel">%s</div><div class="vrecv">%s</div></div>`,
			html.EscapeString(f.Label), formatFieldValue(f.Expected, f.ExpectedMask), formatFieldValue(f.Actual, f.ActualMask))
	}
	a.analysis.WriteString(`</div>`)
	if len(recRouteMismatch) != 0 {
		var plural string
		if len(recRouteMismatch) == 1 {
			plural = " was"
		} else {
			plural = "s were"
		}
		fmt.Fprintf(a.analysis, `<p>NOTE: The %s field%s not provided in the model message.  Recommended values for key fields should be filled in based on the “SCCo ARES/RACES Recommended Form Routing” document (available on the <a href="https://www.scc-ares-races.org/operations/forms/go-kit">“Go Kit Forms” page</a> of the county ARES website) when the message author does not provide them.</p>`,
			english.Conjoin(recRouteMismatch, "and"), plural)
	}
}

// fixupRecRouteFields modifies the comparison of the fields covered by the
// recommended routing cheat sheet, to address the possibility that they weren't
// supplied in the model message.
func (a *Analysis) fixupRecRouteFields(score int, fields []*field.ComparedField, mtc *config.MessageTypeConfig) (_ int, mismatches []string) {
	for _, f := range fields {
		switch f.Label {
		case "To ICS Position":
			if f.Expected == "" && len(mtc.ToICSPosition) != 0 {
				f.ExpectedMask = "_"
				if slices.Contains(mtc.ToICSPosition, f.Actual) {
					f.Expected = f.Actual
					score += f.OutOf - f.Score
					f.Score = f.OutOf
					f.ActualMask = " "
				} else {
					mismatches = append(mismatches, "“To ICS Position”")
					f.Expected = english.Conjoin(mtc.ToICSPosition, "or")
					f.Label += " [See NOTE]"
					score -= f.Score
					f.Score = 0
					f.ActualMask = "*"
				}
			}
		case "To Location":
			if f.Expected == "" && len(mtc.ToLocation) != 0 {
				f.ExpectedMask = "_"
				if slices.Contains(mtc.ToLocation, f.Actual) {
					f.Expected = f.Actual
					score += f.OutOf - f.Score
					f.Score = f.OutOf
					f.ActualMask = " "
				} else {
					mismatches = append(mismatches, "“To Location”")
					f.Expected = english.Conjoin(mtc.ToLocation, "or")
					f.Label += " [See NOTE]"
					score -= f.Score
					f.Score = 0
					f.ActualMask = "*"
				}
			}
		case "Handling":
			if f.Expected == "" {
				handling := mtc.HandlingOrder
				if handling == "computed" {
					handling = config.ComputeRecommendedHandlingOrder(a.msg)
				}
				if handling != "" {
					f.Expected = handling
					f.ExpectedMask = "_"
					if handling == f.Actual {
						score += f.OutOf - f.Score
						f.Score = f.OutOf
						f.ActualMask = " "
					} else {
						mismatches = append(mismatches, "“Handling”")
						f.Label += " [See NOTE]"
						score -= f.Score
						f.Score = 0
						f.ActualMask = "*"
					}
				}
			}
		}
	}
	return score, mismatches
}

func formatFieldValue(value, mask string) string {
	var (
		sb    strings.Builder
		style string
	)
	for i, c := range value {
		nstyle := style
		if i < len(mask) {
			var ok bool
			if nstyle, ok = styles[mask[i]]; !ok {
				nstyle = "major"
			}
		}
		if nstyle != style {
			if style != "" {
				sb.WriteString("</span>")
			}
			style = nstyle
			if style != "" {
				fmt.Fprintf(&sb, "<span class=%s>", style)
			}
		}
		switch c {
		case '<':
			sb.WriteString("&lt;")
		case '>':
			sb.WriteString("&gt;")
		case '&':
			sb.WriteString("&amp;")
		default:
			sb.WriteRune(c)
		}
	}
	if style != "" {
		sb.WriteString("</span>")
	}
	return sb.String()
}

var styles = map[byte]string{
	' ': "",
	'_': "recroute",
	'~': "minor",
	'*': "major",
}
