// Package analyze handles the parsing, analysis, and storage of incoming
// messages, and the generation of responses to them.
package analyze

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"strings"
	"time"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/store"
	"github.com/rothskeller/packet/xscmsg"
)

// Problem codes
const (
	ProblemBounceMessage  = "BounceMessage"
	ProblemMessageCorrupt = "MessageCorrupt"
)

// ProblemLabel is a map from problem code to problem label (the short string
// that appears under a message in a report when the message has that problem).
var ProblemLabel = map[string]string{
	ProblemBounceMessage:  "message has no return address (probably auto-response)",
	ProblemMessageCorrupt: "message could not be parsed",
	// others are added by init funcs in the files that detect those problems
}

// An Analysis contains the analysis of a received message.
type Analysis struct {
	// raw is the raw received message.
	raw string
	// msg is the decoded received message.
	msg *pktmsg.Message
	// xsc is the recognized XSC message of the received message.
	xsc xscmsg.XSCMessage
	// hash is the hash of the raw received message, used for deduplication.
	hash string
	// localID is the local message ID of the received message.
	localID string
	// session is the session for which the message was received.
	session *store.Session
	// toBBS is the name of the BBS to which the message was sent (i.e.,
	// the one from which we retrieved it).
	toBBS string
	// jurisdiction is the jurisdiction of the sender, if known.
	jurisdiction string
	// problems is the list of problems found with the message.
	problems []*problem
	// responses is a list of messages that should be sent in response to
	// this received message.  These would include delivery receipts and/or
	// problem reports.
	responses []*response
	// fromCallSign is the call sign of the message sender, which can come
	// from any of several places.
	fromCallSign string
	// subjectCallSign is the call sign found in the Practice part of the
	// subject line, if any.
	subjectCallSign string
	// subjectDate is the net date found in the Practice part of the subject
	// line, if any.
	subjectDate time.Time
}

// A response is a single outgoing message that should be sent in response to
// the received message.
type response struct {
	to   []string
	body string
}

// A problem describes a problem with a message found during analysis.
type problem struct {
	// code is a machine-readable code that identifies the type of problem.
	code string
	// response is the text of the response message reporting this problem.
	response string
	// references is a bitmask of the references that should be listed at
	// the bottom of the response message.
	references reference
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
	a.hash = hex.EncodeToString(sum[:])
	// Log receipt of the message.
	a.msg, err = pktmsg.ParseMessage(raw)
	if err != nil && a.msg.ReturnAddress() == "" {
		log.Printf("Received at %s@%s: [UNPARSEABLE with hash %s]", session.CallSign, bbs, a.hash)
	} else {
		log.Printf("Received at %s@%s: from %q subject %q",
			session.CallSign, bbs, a.msg.ReturnAddress(), a.msg.Header.Get("Subject"))
	}
	// If we've already handled the message, stop.
	if a.localID = st.HasMessageHash(a.hash); a.localID != "" {
		log.Printf("=> already handled as %s", a.localID)
		return nil
	}
	// Assign it a local message ID.
	a.localID = st.NextMessageID(a.session.Prefix)
	// Check the message for problems and log them in the analysis.
	a.checkMessage(err)
	return &a
}

// checkMessage checks the message for problems, logging them in the problems
// list of the Analysis structure.
func (a *Analysis) checkMessage(parseerr error) {
	if parseerr != nil {
		a.problems = append(a.problems, &problem{code: ProblemMessageCorrupt})
		return
	}
	if a.msg.Flags&pktmsg.AutoResponse != 0 {
		a.problems = append(a.problems, &problem{code: ProblemBounceMessage})
		return
	}
	// Find out whether the message is a known type.  If it's not, we'll put
	// it in our pseudo-type for plain text or unknown form, whichever fits.
	if a.xsc = xscmsg.Recognize(a.msg, true); a.xsc == nil {
		if form := pktmsg.ParseForm(a.msg.Body, true); form != nil {
			a.xsc = &config.UnknownForm{M: a.msg, F: form}
		} else {
			a.xsc = &config.PlainTextMessage{M: a.msg}
		}
	}
	// If the message is a delivery or read receipt, we don't want to run
	// most of the checks because it's not actually a practice attempt.
	if a.checkReceipts() {
		return
	}
	// Run all of the various checks against the message.
	a.checkPlainText()         // was the message entirely plain ASCII text?
	a.checkValidForm()         // is the form valid? (encoding, field values, required fields)
	a.checkSubjectLine()       // does the message have a valid/correct subject line?
	a.checkMessageNumber()     // does the message number have the correct format?
	a.checkPracticeSubject()   // does the message have the correct "Practice ..." subject?
	a.checkPracticeWindow()    // was the message sent within the practice window for the net?
	a.checkFormVersion()       // is the form using a current encoding and version?
	a.checkCallSign()          // can we find the sender's call sign?
	a.checkBBS()               // was the message sent from or to the wrong BBS?
	a.checkCorrectForm()       // did the message use the correct form?
	a.checkFormHandlingOrder() // does the form have the correct handling order?
	a.checkFormDestination()   // does the form have the correct destination?
	// NOTE: these checks are mostly unordered.  However,
	// checkPracticeWindow and checkCallSign both have to come after
	// checkPracticeSubject because they rely on data extracted from the
	// practice subject.
}

// Commit commits the analyzed message to the database.
func (a *Analysis) Commit(st astore) {
	var (
		m       store.Message
		tag     string
		actions = config.Get().ProblemActionFlags
	)
	if a == nil { // message already handled, nothing to commit
		return
	}
	m.LocalID = a.localID
	m.Hash = a.hash
	m.Message = a.raw
	m.Session = a.session.ID
	m.FromAddress = a.msg.ReturnAddress()
	m.FromCallSign = a.fromCallSign
	m.ToBBS = a.toBBS
	m.Jurisdiction = a.jurisdiction
	if a.xsc != nil {
		m.MessageType = a.xsc.TypeTag()
	}
	m.Subject = a.msg.Header.Get("Subject")
	m.DeliveryTime = a.msg.Date()
	m.FromBBS = a.msg.FromBBS()
	for _, p := range a.problems {
		m.Actions |= actions[p.code]
		m.Problems = append(m.Problems, p.code)
	}
	st.SaveMessage(&m)
	if a.xsc != nil {
		tag = a.xsc.TypeTag()
	} else {
		tag = "-"
	}
	log.Printf("=> %s %s %s", m.LocalID, tag, strings.Join(m.Problems, ","))
}
