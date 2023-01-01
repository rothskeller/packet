// Package analyze handles the parsing, analysis, and storage of incoming
// messages, and the generation of responses to them.
package analyze

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/store"
	"github.com/rothskeller/packet/xscmsg"
)

// An Analysis contains the analysis of a received message.
type Analysis struct {
	// raw is the raw received message.
	raw string
	// msg is the decoded received message.
	msg *pktmsg.Message
	// xsc is the recognized XSC message of the received message.
	xsc *xscmsg.Message
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
	// problems is the set of problems found with the message.  Each one
	// maps to the appropriate response key for that problem.
	problems map[*Problem]string
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
// map of the Analysis structure.
func (a *Analysis) checkMessage(parseerr error) {
	a.problems = make(map[*Problem]string)
	if parseerr != nil {
		a.problems[ProbMessageCorrupt] = ""
		return
	}
PROBLEMS:
	for _, prob := range orderedProblems() {
		for _, ifnot := range prob.ifnot {
			if _, ok := a.problems[ifnot]; ok {
				continue PROBLEMS
			}
		}
		if found, responseKey := prob.detect(a); found {
			a.problems[prob] = responseKey
		}
	}
}

// Commit commits the analyzed message to the database.
func (a *Analysis) Commit(st astore) {
	var (
		m        store.Message
		tag      string
		problems = config.Get().Problems
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
		m.MessageType = a.xsc.Type.Tag
	}
	m.Subject = a.msg.Header.Get("Subject")
	m.DeliveryTime = a.msg.Date()
	m.FromBBS = a.msg.FromBBS()
	for p := range a.problems {
		m.Actions |= problems[p.Code].ActionFlags
		m.Problems = append(m.Problems, p.Code)
	}
	sort.Strings(m.Problems)
	st.SaveMessage(&m)
	if a.xsc != nil {
		tag = a.xsc.Type.Tag
	} else {
		tag = "-"
	}
	log.Printf("=> %s %s %s", m.LocalID, tag, strings.Join(m.Problems, ","))
}
