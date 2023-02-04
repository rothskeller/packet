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
	// Determine the message type (if no parse error and not a bounce).
	if err == nil && a.msg.Flags&pktmsg.AutoResponse == 0 {
		a.xsc = xscmsg.Recognize(a.msg, true)
	}
	// Find the problems with the message.
	a.findProblems(err)
	return &a
}

// findProblems finds all of the problems with the message.
func (a *Analysis) findProblems(parseErr error) {
	a.problems = make(map[string]struct{})
	if a.notPracticeMessage(parseErr) {
		return
	}
	a.checkSender() // sets a.FromCallSign
	a.checkReceiver()
	a.checkSubject() // sets a.subject and a.Practice
	a.checkDate()
	a.checkBody()
}

// notPracticeMessage checks for problems that indicate the message is not a
// practice message at all, and shouldn't get the full analysis.
func (a *Analysis) notPracticeMessage(parseErr error) bool {
	switch {
	case parseErr != nil:
		// If we had a parse error, note that and don't analyze further.
		a.problems["MessageCorrupt"] = struct{}{}
		return true
	case a.msg.Flags&pktmsg.AutoResponse != 0:
		// If it was a bounce message, note that and don't analyze further.
		a.problems["BounceMessage"] = struct{}{}
		return true
	case a.xsc.Type.Tag == delivrcpt.Tag:
		a.problems["DeliveryReceipt"] = struct{}{}
		return true
	case a.xsc.Type.Tag == readrcpt.Tag:
		a.reportProblem("ReadReceipt", refOutpostConfig, readReceiptResponse)
		return true
	}
	return false
}

// MultipleMessagesFromAddress is a pseudo-problem used during reporting, to
// mark messages that have been superseded by other messages.
const MultipleMessagesFromAddress = "MultipleMessagesFromAddress"

func init() {
	ProblemLabels["BounceMessage"] = "message has no return address (probably auto-response)"
	ProblemLabels["DeliveryReceipt"] = "DELIVERED receipt message"
	ProblemLabels["MessageCorrupt"] = "message could not be parsed"
	ProblemLabels[MultipleMessagesFromAddress] = "multiple messages from this address"
	ProblemLabels["ReadReceipt"] = "unexpected READ receipt message"
}

const readReceiptResponse = `This message is an Outpost "read receipt", which
should not have been sent.  Most likely, your Outpost installation has the
"Auto-Read Receipt" setting turned on.  The SCCo-standard Outpost configuration
specifies that this setting should be turned off.  You can find it on the
Receipts tab of the Message Settings dialog in Outpost.`

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
		actions |= config.Get().ProblemActionFlags[p]
		problems = append(problems, p)
	}
	sort.Strings(problems)
	return
}
