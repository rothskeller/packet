// Package analyze handles the parsing, analysis, and storage of incoming
// messages, and the generation of responses to them.
package analyze

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
	"github.com/rothskeller/packet/message/jurisstat"
	"github.com/rothskeller/packet/message/plaintext"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/store"
)

// An Analysis contains the analysis of a received message.
type Analysis struct {
	// raw is the raw received message.
	raw string
	// env is the envelope of the received message.
	env *envelope.Envelope
	// subject is the message subject line.
	subject string
	// body is the message body (encoded).
	body string
	// msg is the decoded message contents.
	msg message.Message
	// Hash is the hash of the raw received message, used for deduplication.
	Hash string
	// localID is the local message ID of the received message.
	localID string
	// session is the session for which the message was received.
	session *store.Session
	// toBBS is the name of the BBS to which the message was sent (i.e.,
	// the one from which we retrieved it).
	toBBS string

	// The following are items computed by getAnalysisData.

	// key is the set of key fields from the message.
	key *message.KeyFields
	// severity is the severity code parsed from the subject line.
	severity string
	// formtag is the form tag parsed from the subject line.
	formtag string
	// jurisdiction is the jurisdiction decoded from the practice subject.
	jurisdiction string
	// netDate is the net date decoded from the practice subject.
	netDate time.Time
	// fromBBS is the name of the BBS from which the message was sent, if
	// discernable.
	fromBBS string
	// FromCallSign is the call sign of the message sender, which can come
	// from any of several places.
	FromCallSign string
	// corruptForm indicates that the message, parsed as plain text, looks
	// like a form that couldn't be parsed.
	corruptForm bool

	// The following are the result of the analysis.

	// problems is the set of problems found with the message.
	problems map[string]struct{}
	// reportSubject is the subject of the problem report message.
	reportSubject string
	// reportText is the text of the problem report message.
	reportText strings.Builder
	// references is the bitmask of references that should be cited in the
	// problem report message.
	references reference
	// invalid is a flag indicating that the message should not be counted.
	invalid bool
	// responses is a list of messages that should be sent in response to
	// this received message.  These would include delivery receipts and/or
	// problem reports.
	responses []*response
}

// A response is a single outgoing message that should be sent in response to
// the received message.
type response struct {
	to   []string
	body string
}

// astore is the interface that the store passed into analyze package functions
// must implement.  In production it will be *store.Store, but it can be stubbed
// for testing.
type astore interface {
	HasMessageHash(string) string
	NextMessageID(string) string
	SaveMessage(*store.Message)
}

// Analyze analyzes a single received message, and returns its analysis.
func Analyze(st astore, session *store.Session, bbs, raw string) *Analysis {
	var (
		a   Analysis
		sum [20]byte
		err error
	)
	// Store the basic message information in the analysis.
	a.raw = raw
	a.session = session
	a.toBBS = bbs
	sum = sha1.Sum([]byte(raw))
	a.Hash = hex.EncodeToString(sum[:])
	// Log receipt of the message.
	a.env, a.subject, a.body, err = envelope.ParseRetrieved(raw, bbs, "")
	if err != nil && a.env.ReturnAddr == "" {
		log.Printf("Received at %s@%s: [UNPARSEABLE with hash %s]", session.CallSign, bbs, a.Hash)
	} else {
		log.Printf("Received at %s@%s: from %q subject %q", session.CallSign, bbs, a.env.ReturnAddr, a.subject)
	}
	// If we've already handled the message, stop.
	if a.localID = st.HasMessageHash(a.Hash); a.localID != "" {
		log.Printf("=> already handled as %s", a.localID)
		return nil
	}
	// Assign it a local message ID.
	a.localID = st.NextMessageID(a.session.Prefix)
	// Determine the message type (if no parse error and not a bounce).
	if err == nil && !a.env.Autoresponse {
		a.msg = message.Decode(a.subject, a.body)
	}
	// Find the problems with the message.
	a.problems = make(map[string]struct{})
	if a.humanMessage(err) {
		a.getAnalysisData()
		a.findProblems()
	}
	return &a
}

func (a *Analysis) humanMessage(parseErr error) bool {
	return !a.messageCorrupt(parseErr) && !a.bounceMessage() && !a.deliveryReceipt() && !a.readReceipt()
}

var (
	// practiceCallSignRE is a permissive regexp that looks for a call sign
	// in the "Practice ..." portion of a subject line or field.  It matches
	// anything that looks like a call sign (FCC or tactical) in the word
	// after "Practice".  The regexp returns the matched call sign, which
	// may be in any case.
	practiceCallSignRE = regexp.MustCompile(`(?i)(?:\b|_)Practice[\W_]+(A[A-L][0-9][A-Z]{1,3}|[KNW][A-Z][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3}|[A-Z][A-Z0-9]{5})\b`)
	// fromCallSignRE extracts the fromCallSign from the return address.  It
	// looks for a call sign at the start of the string, followed either by
	// a %, an @, or the end of the string.  It is not case-sensitive.  The
	// substring returned is the call sign.
	fromCallSignRE = regexp.MustCompile(`(?i)^(A[A-L][0-9][A-Z]{1,3}|[KNW][A-Z][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3}|[A-Z][A-Z0-9]{5})(?:@|%|$)`)
	// fccCallSignRE matches a legal FCC call sign.  It is not
	// case-sensitive.
	fccCallSignRE = regexp.MustCompile(`(?i)^(?:A[A-L][0-9][A-Z]{1,3}|[KNW][A-Z][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3})$`)
	// practiceRE matches a correctly formatted practice subject.  The
	// subject must have the word Practice followed by four comma-separated
	// fields (with whitespace also allowed between the fields).  The RE
	// returns the third field (the jurisdiction) and the fourth field (the
	// date) as substrings so that they can be further checked and stored.
	// A comma is allowed after the word "Practice", which doesn't exactly
	// conform to the required syntax, but it is a very common mistake and
	// not worth penalizing.
	practiceRE = regexp.MustCompile(`(?i)^Practice[,\s]+(?:A[A-L][0-9][A-Z]{1,3}|[KNW][A-Z][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3}|[A-Z][A-Z0-9]{5})\s*,[^,]+,([^,]+),\s*((?:0?[1-9]|1[0-2])/(?:0?[1-9]|[12]\d|3[01])/20\d\d)\s*$`)
	// fromBBSRE matches a return address from a BBS, and returns the BBS
	// name.  It is the first word of the address domain, as long as that
	// address looks like a call sign and the rest of the domain is
	// ".ampr.org" or a ".#" BBS network domain.
	fromBBSRE = regexp.MustCompile(`(?i)^[^%@]+[%@](A[A-L][0-9][A-Z]{1,3}|[KNW][A-Z][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3})(?:\.ampr\.org(?:@.*)?|\.#.*)?$`)
)

func (a *Analysis) getAnalysisData() {
	if f, ok := a.msg.(message.IKeyFields); ok {
		a.key = f.KeyFields()
	} else {
		a.key = new(message.KeyFields)
		a.key.OriginMsgID, a.severity, a.key.Handling, a.formtag, a.key.Subject = common.DecodeSubject(a.subject)
		if strings.EqualFold(a.formtag, "Practice") {
			// Problem will be reported in practiceAsFormName, but
			// meanwhile let's fix it so that the practice subject
			// can be parsed.
			a.key.Subject = "Practice " + a.key.Subject
		}
	}
	if match := fromBBSRE.FindStringSubmatch(a.env.ReturnAddr); match != nil {
		a.fromBBS = strings.ToUpper(match[1])
	}
	if match := practiceCallSignRE.FindStringSubmatch(a.subject); match != nil {
		a.FromCallSign = strings.ToUpper(match[1])
	} else if match := practiceCallSignRE.FindStringSubmatch(a.key.Subject); match != nil {
		a.FromCallSign = strings.ToUpper(match[1])
	} else if match := fromCallSignRE.FindStringSubmatch(a.env.ReturnAddr); match != nil && (fccCallSignRE.MatchString(match[1]) || a.fromBBS != "") {
		// A tactical call sign in the return address counts only if
		// the message is coming from within our BBS network.
		a.FromCallSign = strings.ToUpper(match[1])
	} else {
		a.FromCallSign = a.key.OpCall
	}
	if f, ok := a.msg.(*jurisstat.JurisStat); ok && message.OlderVersion(f.FormVersion, "2.2") {
		// If we have an old Municipal Status form, the subject doesn't
		// have the full practice details; it only has the jurisdiction.
		a.jurisdiction = a.key.Subject
	} else if match := practiceRE.FindStringSubmatch(a.key.Subject); match != nil {
		a.jurisdiction = strings.TrimSpace(match[1])
		a.netDate, _ = time.ParseInLocation("1/2/2006", match[2], time.Local)
	}
	if abbr, ok := config.Get().Jurisdictions[strings.ToUpper(a.jurisdiction)]; ok {
		a.jurisdiction = abbr
	}
	if m, ok := a.msg.(*plaintext.PlainText); ok &&
		(strings.Contains(m.Body, "!SCCoPIFO!") || strings.Contains(m.Body, "!PACF!") || strings.Contains(m.Body, "!/ADDON!")) {
		a.corruptForm = true
	}
}

func (a *Analysis) findProblems() {
	if !a.messageFromWinlink() {
		a.messageNotPlainText()
	}
	a.messageNotASCII()
	a.fromBBSDown()
	if !a.toBBSDown() {
		a.toBBSWrong()
	}
	if !a.practiceSubjectFormat() { // also handles formPracticeSubject
		a.unknownJurisdiction()
	}
	if !a.msgNumFormat() {
		a.msgNumPrefix()
	}
	if !a.messageTooLate() &&
		!a.messageTooEarly() {
		a.sessionDate()
	}
	a.messageTypeWrong()
	if _, ok := a.msg.(message.IKeyFields); ok {
		if !a.formNoCallSign() {
			a.callSignConflict()
		}
		a.formSubject()
		a.formInvalid()
		a.pifoVersion()
		a.formVersion()
		a.formDestination() // also handles formToICSPosition and formToLocation
		a.formHandlingOrder()
	} else {
		a.noCallSign()
		if !a.subjectFormat() {
			a.handlingOrderCode()
			a.subjectHasSeverity()
			a.practiceAsFormName()
		}
		if _, ok := a.msg.(*plaintext.PlainText); ok {
			if !a.formCorrupt() {
				a.subjectPlainForm()
			}
		}
	}
}

// Commit commits the analyzed message to the database.
func (a *Analysis) Commit(st astore) {
	var (
		m   store.Message
		tag string
	)
	if a == nil { // message already handled, nothing to commit
		return
	}
	m.LocalID = a.localID
	m.Hash = a.Hash
	m.Message = a.raw
	m.Session = a.session.ID
	m.FromAddress = a.env.ReturnAddr
	m.FromCallSign = a.FromCallSign
	m.ToBBS = a.toBBS
	m.Jurisdiction = a.jurisdiction
	m.MessageType = a.MessageType()
	m.Subject = a.subject
	m.DeliveryTime = a.env.Date
	m.FromBBS = a.fromBBS
	m.Problems, m.Actions = a.ProblemsActions()
	st.SaveMessage(&m)
	if a.msg != nil {
		tag = a.msg.Type().Tag
	} else {
		tag = "-"
	}
	log.Printf("=> %s %s %s", m.LocalID, tag, strings.Join(m.Problems, ","))
}

// MessageType returns the analyzed type of the message.
func (a *Analysis) MessageType() string {
	if a.msg != nil {
		return a.msg.Type().Tag
	}
	return ""
}

// ProblemsActions returns the list of problem codes for the message and the
// resulting action flags.
func (a *Analysis) ProblemsActions() (problems []string, actions config.Action) {
	for p := range a.problems {
		actions |= config.Get().ProblemActionFlags[p]
		problems = append(problems, p)
	}
	sort.Strings(problems)
	return
}
