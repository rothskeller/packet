// Package analyze handles the parsing, analysis, and storage of incoming
// messages, and the generation of responses to them.
package analyze

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"strings"
	"time"

	"steve.rothskeller.net/packet/pktmsg"
	"steve.rothskeller.net/packet/wppsvr/store"
)

// ProblemLabel is a map from problem code to problem label (the short string
// that appears under a message in a report when the message has that problem).
var ProblemLabel = map[string]string{}

// An Analysis contains the analysis of a received message.
type Analysis struct {
	// msg is the received message itself.
	msg pktmsg.ParsedMessage
	// hash is the hash of the raw received message, used for deduplication.
	hash string
	// localID is the local message ID of the received message.
	localID string
	// session is the session for which the message was received.
	session *store.Session
	// toBBS is the name of the BBS to which the message was sent (i.e.,
	// the one from which we retrieved it).
	toBBS string
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
	// subject is the subject line for the response message reporting this
	// problem.
	subject string
	// response is the text of the response message reporting this problem.
	response string
	// references is a bitmask of the references that should be listed at
	// the bottom of the response message.
	references reference
	// invalid is a flag indicating that the message should not be counted
	// as a check-in (even an erroneous one).
	invalid bool
	// warning is a flag indicating that this problem should not be counted
	// as an error.
	warning bool
}

// Analyze analyzes a single received message, and returns its analysis.
func Analyze(st *store.Store, session *store.Session, bbs, raw string) *Analysis {
	var (
		a   Analysis
		sum [20]byte
	)
	a.msg = pktmsg.ParseMessage(raw)
	log.Printf("Received at %s@%s: from %q subject %q",
		session.CallSign, bbs, a.msg.Base().ReturnAddress, a.msg.Base().SubjectLine)
	// Hash the message and find out whether we've already handled it.
	sum = sha1.Sum([]byte(raw))
	a.hash = hex.EncodeToString(sum[:])
	if st.HasMessageHash(a.hash) {
		log.Printf("=> message already handled")
		return nil
	}
	// Store the basic message information in the analysis.
	a.session = session
	a.toBBS = bbs
	a.localID = st.NextMessageID(a.session.Prefix)
	// Run all of the checks against the message.
	a.checkNonHuman()          // was the message a non-human message?
	a.checkPlainText()         // was the message entirely plain ASCII text?
	a.checkValidForm()         // is the form valid? (encoding, field values, required fields)
	a.checkFormVersion()       // is the form using a current encoding and version?
	a.checkFormSubject()       // does the subject agree with the form?
	a.checkSubjectLine()       // does the message have a valid subject line?
	a.checkMessageNumber()     // does the message number have the correct format?
	a.checkFormHandlingOrder() // does the form have the correct handling order?
	a.checkFormDestination()   // does the form have the correct destination?
	a.checkPracticeSubject()   // does the message have the correct "Practice ..." subject?
	a.checkCallSign()          // can we find the sender's call sign?
	a.checkPracticeWindow()    // was the message sent within the practice window for the net?
	a.checkBBS()               // was the message sent from or to the wrong BBS?
	a.checkCorrectForm()       // did the message use the correct form?
	return &a
}

// Commit commits the analyzed message to the database.
func (a *Analysis) Commit(st *store.Store) {
	var m store.Message

	if a == nil { // message already handled, nothing to commit
		return
	}
	m.LocalID = a.localID
	m.Hash = a.hash
	m.Message = a.msg.Base().RawMessage
	m.Session = a.session.ID
	m.FromAddress = a.msg.Base().ReturnAddress
	m.FromCallSign = a.fromCallSign
	m.ToBBS = a.toBBS
	m.Subject = a.msg.Base().SubjectLine
	m.DeliveryTime = a.msg.Base().DeliveryTime
	if am := a.msg.Message(); am != nil {
		m.FromBBS = am.FromBBS
	}
	m.Valid, m.Correct = true, true
	for _, p := range a.problems {
		if p.invalid {
			m.Valid = false
		}
		if !p.warning {
			m.Correct = false
		}
		m.Problems = append(m.Problems, p.code)
	}
	st.SaveMessage(&m)
	log.Printf("=> %s %s %s", m.LocalID, a.msg.TypeCode(), strings.Join(m.Problems, ","))
}
