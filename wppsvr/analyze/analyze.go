// Package analyze handles the parsing, analysis, and storage of incoming
// messages, and the generation of responses to them.
package analyze

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"sort"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/store"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/delivrcpt"
	"github.com/rothskeller/packet/xscmsg/readrcpt"
)

// An Analysis contains the analysis of a received message.
type Analysis struct {
	// raw is the raw received message.
	raw string
	// msg is the decoded received message.
	msg *pktmsg.Message
	// xsc is the recognized XSC message of the received message.
	xsc *xscmsg.Message
	// subject is the parsed XSC subject of the received message.
	subject *xscmsg.XSCSubject
	// Practice is the parsed practice message details of the received
	// message.
	Practice *PracticeSubject
	// Hash is the hash of the raw received message, used for deduplication.
	Hash string
	// localID is the local message ID of the received message.
	localID string
	// session is the session for which the message was received.
	session *store.Session
	// toBBS is the name of the BBS to which the message was sent (i.e.,
	// the one from which we retrieved it).
	toBBS string
	// problems is the set of problems found with the message.  Each one
	// maps to the appropriate response key for that problem.
	problems map[*Problem]struct{}
	// responses is a list of messages that should be sent in response to
	// this received message.  These would include delivery receipts and/or
	// problem reports.
	responses []*response
	// FromCallSign is the call sign of the message sender, which can come
	// from any of several places.
	FromCallSign string
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
	a.msg, err = pktmsg.ParseMessage(raw)
	if err != nil && a.msg.ReturnAddress() == "" {
		log.Printf("Received at %s@%s: [UNPARSEABLE with hash %s]", session.CallSign, bbs, a.Hash)
	} else {
		log.Printf("Received at %s@%s: from %q subject %q",
			session.CallSign, bbs, a.msg.ReturnAddress(), a.msg.Header.Get("Subject"))
	}
	// If we've already handled the message, stop.
	if a.localID = st.HasMessageHash(a.Hash); a.localID != "" {
		log.Printf("=> already handled as %s", a.localID)
		return nil
	}
	// Assign it a local message ID.
	a.localID = st.NextMessageID(a.session.Prefix)
	// If we had a parse error, note that and don't analyze further.
	a.problems = make(map[*Problem]struct{})
	if err != nil {
		a.problems[ProbMessageCorrupt] = struct{}{}
		return &a
	}
	// If it was a bounce message, note that and don't analyze further.
	if a.msg.Flags&pktmsg.AutoResponse != 0 {
		a.problems[ProbBounceMessage] = struct{}{}
		return &a
	}
	// Determine what kind of message it was.
	a.xsc = xscmsg.Recognize(a.msg, true)
	// If it was a receipt, note that and don't analyze further.
	switch a.xsc.Type.Tag {
	case delivrcpt.Tag:
		a.problems[ProbDeliveryReceipt] = struct{}{}
		return &a
	case readrcpt.Tag:
		a.problems[ProbReadReceipt] = struct{}{}
		return &a
	}
	// Parse the subject line.
	if a.subject = xscmsg.ParseSubject(a.msg.Header.Get("Subject")); a.subject != nil {
		a.Practice = a.parsePracticeSubject()
	}
	// Determine the sender's call sign, looking at various sources.
	a.FromCallSign = a.fromCallSign()
	// Check the message for problems and log them in the analysis.
PROBLEMS:
	for _, prob := range orderedProblems() {
		for _, ifnot := range prob.ifnot {
			if _, ok := a.problems[ifnot]; ok {
				continue PROBLEMS
			}
		}
		if found := prob.detect(&a); found {
			a.problems[prob] = struct{}{}
		}
	}
	return &a
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
	m.FromAddress = a.msg.ReturnAddress()
	m.FromCallSign = a.FromCallSign
	m.ToBBS = a.toBBS
	if a.Practice != nil {
		m.Jurisdiction = a.Practice.Jurisdiction
	}
	m.MessageType = a.MessageType()
	m.Subject = a.msg.Header.Get("Subject")
	m.DeliveryTime = a.msg.Date()
	m.FromBBS = a.msg.FromBBS()
	m.Problems, m.Actions = a.ProblemsActions()
	st.SaveMessage(&m)
	if a.xsc != nil {
		tag = a.xsc.Type.Tag
	} else {
		tag = "-"
	}
	log.Printf("=> %s %s %s", m.LocalID, tag, strings.Join(m.Problems, ","))
}

// MessageType returns the analyzed type of the message.
func (a *Analysis) MessageType() string {
	if a.xsc != nil {
		return a.xsc.Type.Tag
	}
	return ""
}

// ProblemsActions returns the list of problem codes for the message and the
// resulting action flags.
func (a *Analysis) ProblemsActions() (problems []string, actions config.Action) {
	for p := range a.problems {
		actions |= config.Get().Problems[p.Code].ActionFlags
		problems = append(problems, p.Code)
	}
	sort.Strings(problems)
	return
}
