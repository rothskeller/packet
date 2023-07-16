package analyze

// This file contains the problem checks that are run against all human
// messages.  They appear in the order they are run, although some are skipped
// based on the message type or the results of previous checks.

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
	"github.com/rothskeller/packet/message/delivrcpt"
	"github.com/rothskeller/packet/message/plaintext"
	"github.com/rothskeller/packet/message/readrcpt"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/english"
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
	// ".ampr.org" or a ".#" BBS network domain.
	fromBBSRE = regexp.MustCompile(`(?i)^[^%@]+[%@](A[A-L][0-9][A-Z]{1,3}|[KNW][A-Z][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3})(?:\.ampr\.org(?:@.*)?|\.#.*)?$`)
	// msgnumRE matches a valid packet message number.
	msgnumRE = regexp.MustCompile(`^(?:[A-Z][A-Z][A-Z]|[A-Z][0-9][A-Z0-9]|[0-9][A-Z][A-Z])-\d\d\d+[PMR]$`)
)

// messageCounts returns whether the message should be counted as a check-in.
// If not, summary and analysis are set appropriately.
func (a *Analysis) messageCounts(parseErr error) bool {
	// First, fail immediately if the message is not a human message.
	if parseErr != nil {
		a.score, a.outOf = 0, 1
		a.setSummary("message could not be parsed")
		fmt.Fprintf(a.analysis, "<h2>Message Could Not Be Parsed</h2><p>This message could not be parsed as a valid RFC-4155 or RFC-5322 message.  The parse error is “<tt>%s</tt>”.</p>",
			html.EscapeString(parseErr.Error()))
		return false
	}
	if a.env.Autoresponse {
		a.score, a.outOf = 0, 1
		a.setSummary("message has no return address (probably auto-response)")
		a.analysis.WriteString("<h2>Message Has No Return Address</h2><p>This message has no return address, which normally means that it is an auto-response message (e.g., an out-of-office response or a bounce message).  It will not be counted.</p>")
		return false
	}
	if _, ok := a.msg.(*delivrcpt.DeliveryReceipt); ok {
		// No summary or analysis for delivery receipt.
		a.score, a.outOf = 0, 1
		return false
	}
	if _, ok := a.msg.(*readrcpt.ReadReceipt); ok {
		a.score, a.outOf = 0, 1
		a.setSummary("unexpected READ receipt message")
		a.analysis.WriteString(`<h2>Unexpected READ Receipt Message</h2><p>This message is an Outpost “read receipt,” which should not have been sent.  Most likely, your Outpost installation has the “Auto-Read Receipt” setting turned on.  The SCCo “Standard Outpost Configuration Instructions” (available on the <a href="https://www.scc-ares-races.org/data/packet/index.html">“Packet BBS Service” page</a> of the county ARES website) specifies that this setting should be turned off.  You can find it on the Receipts tab of the Message Settings dialog in Outpost.</p>`)
		return false
	}
	// Check that it was sent to a correct BBS.
	if inList(a.session.DownBBSes, a.sm.ToBBS) {
		a.setSummary("message to incorrect BBS (simulated outage)")
		fmt.Fprintf(a.analysis, "<h2>Message to Incorrect BBS</h2><p>This message was sent to %[1]s at %[2]s, but %[2]s has a simulated outage for %[3]s on %[4]s.  This message will not be counted.  Practice messages for this session must be sent to %[1]s at %[5]s.</p>",
			a.session.CallSign, a.sm.ToBBS, html.EscapeString(a.session.Name), a.session.End.Format("January 2"),
			english.Conjoin(a.session.ToBBSes, "or"))
	} else if !inList(a.session.ToBBSes, a.sm.ToBBS) {
		a.setSummary("message to incorrect BBS")
		fmt.Fprintf(a.analysis, "<h2>Message to Incorrect BBS</h2><p>This message was sent to %[1]s at %[2]s, but practice messages for %[3]s on %[4]s must be sent to %[1]s at %[5]s.  This message will not be counted.</p>",
			a.session.CallSign, a.sm.ToBBS, html.EscapeString(a.session.Name), a.session.End.Format("January 2"),
			english.Conjoin(a.session.ToBBSes, "or"))
	}
	// Check that it was sent after the start of the session.
	var rcvdate = a.env.BBSReceivedDate
	if rcvdate.IsZero() {
		rcvdate = a.env.Date
	}
	if rcvdate.Before(a.session.Start) {
		a.setSummary("message sent outside of practice session")
		fmt.Fprintf(a.analysis, "<h2>Message Sent Outside of Practice Session</h2><p>This message arrived at %s on %s.  However, practice messages for %s aren’t accepted until %s.  This message will not be counted.</p>",
			a.sm.ToBBS, rcvdate.Format("2006-01-02 at 15:04"), html.EscapeString(a.session.Name),
			a.session.Start.Format("2006-01-02 at 15:04"))
	}
	// Check that we have a call sign so that we know whom to credit.  To do
	// that, we need to know what BBS the message came from, if any.
	if match := fromBBSRE.FindStringSubmatch(a.env.ReturnAddr); match != nil {
		a.sm.FromBBS = strings.ToUpper(match[1])
	}
	if f, ok := a.msg.(message.IKeyFields); ok {
		a.key = f.KeyFields()
	}
	if match := fromCallSignRE.FindStringSubmatch(a.env.ReturnAddr); match != nil && (fccCallSignRE.MatchString(match[1]) || a.sm.FromBBS != "") {
		// We'll take an FCC call sign in the return address as the
		// call sign to credit.  A tactical call sign in the return
		// address counts only if the message is coming from a BBS;
		// otherwise we can't be sure it's a tactical call sign.
		a.sm.FromCallSign = strings.ToUpper(match[1])
	} else if a.key != nil {
		// No call sign in the return address.  But it's a form with an
		// OpCall field; hopefully there's one there.
		a.sm.FromCallSign = a.key.OpCall
	}
	if a.sm.FromCallSign == "" {
		a.setSummary("no call sign in message")
		if _, ok := a.msg.(message.IKeyFields); ok {
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
	if a.env.NotPlainText {
		if strings.Contains(a.env.ReturnAddr, "winlink.org") {
			a.setSummary("message sent from Winlink")
			a.analysis.WriteString("<h2>Message Sent from Winlink</h2><p>This message was sent from Winlink.  Winlink should not be used for emergency communications, unless no alternatives are available, because it uses a message encoding system (“quoted-printable”) that Outpost cannot decode.  As a result, some messages (particularly those with long lines and those containing equals signs) may be garbled in transmission.</p>")
		} else {
			a.setSummary("not a plain text message")
			a.analysis.WriteString("<h2>Not a Plain Text Message</h2><p>This message is not a plain text message. All SCCo packet messages should be plain text only.  (“Rich text” or HTML-formatted messages, common in email systems, are far larger than plain text messages and put too much strain on the packet infrastructure.)  Please configure your software to send plain text messages when sending to an SCCo BBS.</p>")
		}
	} else {
		a.score++
	}
	// Make sure the message has only ASCII characters.
	a.outOf++
	if strings.IndexFunc(a.body, nonASCII) >= 0 {
		a.setSummary("message has non-ASCII characters")
		a.analysis.WriteString("<h2>Message Has Non-ASCII Characters</h2><p>This message contains characters that are not in the standard ASCII character set (i.e., not on a standard keyboard). Non-standard characters should be avoided in packet messages, because the receiving system may not know how to render them.  Note that some software may introduce undesired non-standard characters (e.g., Microsoft Word’s “smart quotes” feature). If you use message text composed in such software, make sure those features are disabled.</p>")
	} else {
		a.score++
	}
	// Make sure the message came from a BBS that is up.
	a.outOf++
	if inList(a.session.DownBBSes, a.sm.FromBBS) {
		a.setSummary("message from incorrect BBS (simulated outage)")
		fmt.Fprintf(a.analysis, "<h2>Message from Incorrect BBS</h2><p>This message was sent from %s, which has a simulated outage for %s on %s.  Practice messages should not be sent from BBSes that have a simulated outage.</p>",
			a.sm.FromBBS, html.EscapeString(a.session.Name), a.session.End.Format("January 2"))
	} else {
		a.score++
	}
	// Some checks only apply to form messages (of known form types).
	if a.key != nil {
		// Make sure the message subject matches the form.
		if f, ok := a.msg.(message.IEncode); ok {
			a.outOf++
			subject := f.EncodeSubject()
			if a.subject != subject && a.subject != strings.TrimRight(subject, " ") {
				a.setSummary("message subject doesn't agree with form contents")
				fmt.Fprintf(a.analysis, `<h2>Message Subject Doesn’t Agree with Form Contents</h2><p style="margin-bottom:0">This message has</p><div style="margin-left:2rem"><tt>Subject: %s</tt></div><div>but, based on the contents of the form, it should have</div><div style="margin-left:2rem"><tt>Subject: %s</tt></div><p style="margin-top:0">PackItForms automatically generates the Subject line from the form contents; it should not be overridden manually.</p>`,
					html.EscapeString(a.subject), html.EscapeString(subject))
			} else {
				a.score++
			}
		}
		// Make sure the message is valid according to PackItForms' rules.
		if f, ok := a.msg.(message.IValidate); ok {
			if problems := f.Validate(); len(problems) != 0 {
				a.outOf += len(problems)
				a.setSummary("invalid form contents")
				a.analysis.WriteString(`<h2>Invalid Form Contents</h2><p style="margin-bottom:0">This message contains a form with invalid contents:</p><ul style="margin-top:0;margin-bottom:0">`)
				for _, problem := range problems {
					fmt.Fprintf(a.analysis, "<li>%s</li>", html.EscapeString(problem))
				}
				a.analysis.WriteString(`</ul><p style="margin-top:0">Please verify the correctness of the form before sending.</p>`)
			}
		}
		// Make sure the PIFO and form versions are up to date.
		a.outOf += 2
		var minPIFO = config.Get().MinPIFOVersion
		var minForm = config.Get().MessageTypes[a.msg.Type().Tag].MinimumVersion
		if message.OlderVersion(a.key.PIFOVersion, minPIFO) {
			a.setSummary("PackItForms version out of date")
			fmt.Fprintf(a.analysis, "<h2>PackItForms Version Out of Date</h2><p>This message used version %s of PackItForms to encode the form, but that version is not current.  Please use PackItForms version %s or newer to encode messages containing forms.</p>",
				a.key.PIFOVersion, minPIFO)
		} else {
			a.score++
		}
		if message.OlderVersion(a.key.FormVersion, minForm) {
			a.setSummary("form version out of date")
			fmt.Fprintf(a.analysis, "<h2>Form Version Out of Date</h2><p>This message contains version %s of the %s, but that version is not current.  Please use version %s or newer of the form.  (You can get the newer form by updating your PackItForms installation.)",
				a.key.FormVersion, html.EscapeString(a.msg.Type().Name), minForm)
		} else {
			a.score++
		}
		a.checkMessageNumber(a.key.OriginMsgID)
	} else { // checks for plain text messages (or forms of unknown type)
		// Check the message subject format.
		a.outOf += 3
		msgid, severity, handling, formtag, _ := common.DecodeSubject(a.subject)
		if msgid == "" {
			a.setSummary("incorrect subject line format")
			a.analysis.WriteString(`<h2>Incorrect Subject Line Format</h2><p>This message has an incorrect subject line format.  According to the SCCo “Standard Packet Message Subject Line” (available on the <a href="https://www.scc-ares-races.org/data/packet/index.html">“Packet BBS Service” page</a> of the county ARES website), the subject line should look like <tt>AAA-111P_R_Subject</tt>, where <tt>AAA-111P</tt> is the message number, <tt>R</tt> is the handling order code, and <tt>Subject</tt> is the message subject.</p>`)
		} else {
			a.checkMessageNumber(msgid)
			a.score++
			if severity != "" {
				a.setSummary("severity on subject line")
				fmt.Fprintf(a.analysis, `<h2>Severity on Subject Line</h2><p>The subject line of this message contains both a Severity code and a Handling Order code (“_%s/%s_”).  This is an outdated subject line style.  The current SCCo “Standard Packet Message Subject Line” (available on the <a href="https://www.scc-ares-races.org/data/packet/index.html">“Packet BBS Service” page</a> of the county ARES website) includes only the Handling Order code on the Subject line (“_%[2]s_”).</p>`,
					severity, handling)
			} else {
				a.score++
			}
			switch handling {
			case "R", "P", "I":
				a.score++
			case "":
				a.setSummary("missing handling order code")
				a.analysis.WriteString(`<h2>Missing Handling Order Code on Subject Line</h2><p>The Subject line of this message does not contain a Handling Order code. As documented in the SCCo “Standard Packet Message Subject Line” (available on the <a href="https://www.scc-ares-races.org/data/packet/index.html">“Packet BBS Service” page</a> of the county ARES website), it must contain an “I” for Immediate, “P” for Priority, or “R” for Routine.</p>`)
			default:
				a.setSummary("unknown handling order code")
				fmt.Fprintf(a.analysis, `<h2>Unknown Handling Order Code on Subject Line</h2><p>The Subject line of this message contains an invalid Handling Order code (“%s”). As documented in the SCCo “Standard Packet Message Subject Line” (available on the <a href="https://www.scc-ares-races.org/data/packet/index.html">“Packet BBS Service” page</a> of the county ARES website), the valid codes are “I” for Immediate, “P” for Priority, and “R” for Routine.</p>`,
					html.EscapeString(handling))
			}
		}
		// If this is actually a plain text message (and not an unknown)
		// form type), there are a couple more things. to check.
		if m, ok := a.msg.(*plaintext.PlainText); ok {
			a.outOf++
			if strings.Contains(m.Body, "!SCCoPIFO!") || strings.Contains(m.Body, "!PACF!") || strings.Contains(m.Body, "!/ADDON!") {
				a.setSummary("incorrectly encoded form")
				a.analysis.WriteString(`<h2>Incorrectly Encoded Form</h2><p>This message appears to contain an encoded form, but the encoding is incorrect.  It appears to have been created or edited by software other than the current PackItForms software.  Please use current PackItForms software to encode messages containing forms.</p>`)
			} else if formtag != "" {
				a.setSummary("form name in subject of non-form message")
				fmt.Fprintf(a.analysis, "<h2>Form Name in Subject Line of Non-Form Message</h2><p>This message has a form name (“%s”) on the subject line, but does not contain a recognizable form.  If this is a plain text message, there should be no form name between the handling order code and the subject.  If this is a form message, the form is improperly encoded and could not be recognized.</p>",
					html.EscapeString(formtag))
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
			a.analysis.WriteString(`<h2>Incorrect Message Number Format</h2><p style="margin-bottom:0">The message number of this message is not formatted correctly.  According to the SCCo “Standard Packet Message Subject Line” document (available on the <a href="https://www.scc-ares-races.org/data/packet/index.html">“Packet BBS Service” page</a> of the county ARES website), it should have a format like "XND-042P", containing:</p><ul style="margin-top:0;margin-bottom:0"><li>a three-character prefix (usually the sender's call sign suffix),</li><li>a dash,</li><li>a number with at least three digits, and</li><li>a “P”, “M”, or “R” suffix.</ul><p style="margin-top:0">All letters should be upper case.  In Outpost, the format of the message number is set in the Message Settings dialog, which should be configured according to the SCCo “Standard Outpost Configuration Instructions” (available on the same page).</p>`)
		} else if fccCallSignRE.MatchString(a.sm.FromCallSign) {
			act := msgid[:3]
			exp := a.sm.FromCallSign[len(a.sm.FromCallSign)-3:]
			if act != exp {
				a.setSummary("incorrect message number prefix")
				fmt.Fprintf(a.analysis, `<h2>Incorrect Message Number Prefix</h2><p>The message number of this message has the prefix “%s”.  According to the SCCo “Standard Packet Message Subject Line” document (available on the <a href="https://www.scc-ares-races.org/data/packet/index.html">“Packet BBS Service” page</a> of the county ARES website), the prefix should be the last three characters of your call sign, “%s”.</p>`,
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
	if a.key != nil {
		// Make sure the message has a destination allowed by the
		// recommended routing cheat sheet.
		var (
			mtc      *config.MessageTypeConfig
			badpos   bool
			badloc   bool
			exppos   string
			exploc   string
			handling string
		)
		mtc = config.Get().MessageTypes[a.msg.Type().Tag]
		if len(mtc.ToICSPosition) != 0 {
			a.outOf++
			if badpos = !inList(mtc.ToICSPosition, a.key.ToICSPosition); !badpos {
				a.score++
			}
		}
		if len(mtc.ToLocation) != 0 {
			a.outOf++
			if badloc = !inList(mtc.ToLocation, a.key.ToLocation); !badloc {
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
			fmt.Fprintf(a.analysis, `<h2>Incorrect Destination for Form</h2><p>This message form is addressed to ICS Position “%s” at Location “%s”.  According to the “SCCo ARES/RACES Recommended Form Routing” document (available on the <a href="https://www.scc-ares-races.org/operations/go-kit-forms.">“Go Kit Forms” page</a> of the county ARES website), %ss should be addressed to %s at %s.</p>`,
				html.EscapeString(a.key.ToICSPosition), html.EscapeString(a.key.ToLocation), html.EscapeString(a.msg.Type().Name), exppos, exploc)
		} else if badpos {
			a.setSummary(`incorrect "To ICS Position" for form`)
			fmt.Fprintf(a.analysis, `<h2>Incorrect “To ICS Position” for Form</h2><p>This message form is addressed to ICS Position “%s”.  According to the “SCCo ARES/RACES Recommended Form Routing” document (available on the <a href="https://www.scc-ares-races.org/operations/go-kit-forms.html">“Go Kit Forms” page</a> of the county ARES website), %ss should be addressed to ICS Position %s.</p>`,
				html.EscapeString(a.key.ToICSPosition), html.EscapeString(a.msg.Type().Name), exppos)
		} else if badloc {
			a.setSummary(`incorrect "To Location" for form`)
			fmt.Fprintf(a.analysis, `<h2>Incorrect “To Location” for Form</h2><p>This message form is addressed to Location “%s”.  According to the “SCCo ARES/RACES Recommended Form Routing” document (available on the <a href="https://www.scc-ares-races.org/operations/go-kit-forms.html">“Go Kit Forms” page</a> of the county ARES website), %ss should be addressed to Location %s.</p>`,
				html.EscapeString(a.key.ToLocation), html.EscapeString(a.msg.Type().Name), exploc)
		}
		// Make sure the message has a handling order allowed by the
		// recommended routing cheat sheet.
		handling = config.Get().MessageTypes[a.msg.Type().Tag].HandlingOrder
		if handling == "computed" {
			handling = config.ComputeRecommendedHandlingOrder(a.msg)
		}
		if handling != "" {
			a.outOf++
			if handling != "" && handling != a.key.Handling {
				a.setSummary("incorrect handling order for form")
				fmt.Fprintf(a.analysis, `<h2>Incorrect Handling Order for Form</h2><p>This message has handling order %s.  According to the “SCCo ARES/RACES Recommended Form Routing” document (available on the <a href="https://www.scc-ares-races.org/operations/go-kit-forms.html">“Go Kit Forms” page</a> of the county ARES website), it should have handling order %s.</p>`,
					html.EscapeString(a.key.Handling), handling)
			} else {
				a.score++
			}
		}
	}
	// Make sure the message is of a type allowed for the session.
	a.outOf++
	var allowed = a.session.MessageTypes
	if config.Get().BBSes[a.sm.FromBBS] == nil {
		// Plain text messages are always OK when they come from outside
		// the county BBS system.
		allowed = append(allowed, plaintext.Type.Tag)
	}
	if m, ok := a.msg.(*plaintext.PlainText); ok &&
		(strings.Contains(m.Body, "!SCCoPIFO!") || strings.Contains(m.Body, "!PACF!") || strings.Contains(m.Body, "!/ADDON!")) {
		// Allow a "plain text" message containing a corrupt form; that
		// problem gets reported elsewhere.
		allowed = append(allowed, plaintext.Type.Tag)
	}
	if !inList(allowed, a.msg.Type().Tag) {
		var (
			allowed []string
			article string
		)
		for i, code := range a.session.MessageTypes {
			mtype := message.RegisteredTypes[code]
			allowed = append(allowed, html.EscapeString(mtype.Name))
			if i == 0 {
				article = mtype.Article
			}
		}
		a.setSummary("incorrect message type")
		fmt.Fprintf(a.analysis, "<h2>Incorrect Message Type</h2><p>This message is %s %s.  For the %s on %s, %s %s is expected.</p>",
			a.msg.Type().Article, html.EscapeString(a.msg.Type().Name), html.EscapeString(a.session.Name),
			a.session.End.Format("January 2"), article, english.Conjoin(allowed, "or"))
	} else {
		a.score++
	}
}

// compareAgainstModel compares the received message against the model message
// for the session.  Any problems are added to the analysis.
func (a *Analysis) compareAgainstModel() {
	// Make sure the received message is the same type as the model.
	if a.msg.Type() != a.session.ModelMsg.Type() {
		a.outOf *= 2 // Give a 50% score.
		a.setSummary("incorrect message type")
		fmt.Fprintf(a.analysis, "<h2>Incorrect Message Type</h2><p>This message is %s %s.  For the %s on %s, operators are expected to send a copy of the provided %s.</p>",
			a.msg.Type().Article, html.EscapeString(a.msg.Type().Name), html.EscapeString(a.session.Name),
			a.session.End.Format("January 2"), html.EscapeString(a.session.ModelMsg.Type().Name))
		return
	}
	// Compare the message against the model.
	var score, outOf, fields = a.session.ModelMsg.Compare(a.msg)
	// The model may have left destination or handling blank, as an exercise
	// for the operator to look them up in the recommended routing cheat
	// sheet.  If so, we need to fix up the results of the comparison for
	// that.
	var recRouteMismatch []string
	if mtc := config.Get().MessageTypes[a.session.ModelMsg.Type().Tag]; mtc != nil {
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
		fmt.Fprintf(a.analysis, `<p>NOTE: The %s field%s not provided in the model message.  Recommended values for key fields should be filled in based on the “SCCo ARES/RACES Recommended Form Routing” document (available on the <a href="https://www.scc-ares-races.org/operations/go-kit-forms.">“Go Kit Forms” page</a> of the county ARES website) when the message author does not provide them.</p>`,
			english.Conjoin(recRouteMismatch, "and"), plural)
	}
}

// fixupRecRouteFields modifies the comparison of the fields covered by the
// recommended routing cheat sheet, to address the possibility that they weren't
// supplied in the model message.
func (a *Analysis) fixupRecRouteFields(score int, fields []*message.CompareField, mtc *config.MessageTypeConfig) (_ int, mismatches []string) {
	for _, f := range fields {
		switch f.Label {
		case "To ICS Position":
			if f.Expected == "" && len(mtc.ToICSPosition) != 0 {
				f.ExpectedMask = "_"
				if inList(mtc.ToICSPosition, f.Actual) {
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
				if inList(mtc.ToLocation, f.Actual) {
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
		var nstyle = style
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
