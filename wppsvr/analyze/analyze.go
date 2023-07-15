// Package analyze handles the parsing, analysis, and storage of incoming
// messages, and the generation of responses to them.
package analyze

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"strings"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/wppsvr/store"
)

// An Analysis contains the analysis of a received message.
type Analysis struct {
	// sm is the message record that will be stored in the database.
	sm store.Message
	// env is the envelope of the received message.
	env *envelope.Envelope
	// subject is the message subject line.
	subject string
	// body is the message body (encoded).
	body string
	// msg is the decoded message contents.
	msg message.Message
	// session is the session for which the message was received.
	session *store.Session
	// key is the set of key fields from the message.
	key *message.KeyFields
	// score is the number of score points, 0 <= score <= outOf.
	score int
	// outOf is the maximum number of score points.
	outOf int
	// analysis is a strings.Builder for building sm.Analysis.
	analysis *strings.Builder
}

// astore is the interface that the store passed into analyze package functions
// must implement.  In production it will be *store.Store, but it can be stubbed
// for testing.
type astore interface {
	HasMessageHash(string) string
	NextMessageID(string) string
	SaveMessage(*store.Message)
}

// Analyze analyzes a single received message, and returns its analysis.  The
// analysis is not persisted in the database until its Commit method is called.
func Analyze(st astore, session *store.Session, bbs, raw string) *Analysis {
	var (
		a   Analysis
		sum [20]byte
		err error
	)
	// Store the basic message information in the analysis.
	a.sm.Message = raw
	a.sm.Session = session.ID
	a.session = session
	a.sm.ToBBS = bbs
	sum = sha1.Sum([]byte(raw))
	a.sm.Hash = hex.EncodeToString(sum[:])
	// Log receipt of the message.
	a.env, a.subject, a.body, err = envelope.ParseRetrieved(raw, bbs, "")
	a.sm.DeliveryTime = a.env.Date
	a.sm.FromAddress = a.env.ReturnAddr
	if err != nil && a.env.ReturnAddr == "" {
		log.Printf("Received at %s@%s: [UNPARSEABLE with hash %s]", session.CallSign, bbs, a.sm.Hash)
	} else {
		log.Printf("Received at %s@%s: from %q subject %q", session.CallSign, bbs, a.env.ReturnAddr, a.subject)
	}
	// If we've already handled the message, stop.
	if a.sm.LocalID = st.HasMessageHash(a.sm.Hash); a.sm.LocalID != "" {
		log.Printf("=> already handled as %s", a.sm.LocalID)
		return nil
	}
	// Assign it a local message ID.
	a.sm.LocalID = st.NextMessageID(a.session.Prefix)
	// Determine the message type (if no parse error and not a bounce).
	if err == nil && !a.env.Autoresponse {
		a.msg = message.Decode(a.subject, a.body)
		a.sm.MessageType = a.msg.Type().Tag
	}
	// Find the problems with the message.
	a.analysis = new(strings.Builder)
	if a.messageCounts(err) {
		a.checkCorrectness()
		if a.session.ModelMsg == nil {
			a.checkNonModel()
		} else {
			a.compareAgainstModel()
		}
	}
	if a.sm.Summary == "" && a.score == a.outOf {
		a.sm.Summary = "OK"
	}
	a.sm.Analysis = a.analysis.String()
	a.sm.Score = a.score * 100 / a.outOf
	return &a
}

// setSummary sets the summary line of the analysis, handling the possibility
// of multiple issues.
func (a *Analysis) setSummary(s string) {
	if a.sm.Summary != "" {
		a.sm.Summary = "multiple issues"
	} else {
		a.sm.Summary = s
	}
}

// Commit commits the analyzed message to the database.
func (a *Analysis) Commit(st astore) {
	var tag string

	if a == nil { // message already handled, nothing to commit
		return
	}
	a.fetchJurisdiction()
	st.SaveMessage(&a.sm)
	if a.msg != nil {
		tag = a.msg.Type().Tag
	} else {
		tag = "-"
	}
	log.Printf("=> %s %s %d", a.sm.LocalID, tag, a.sm.Score)
}

func inList[T comparable](list []T, item T) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}
