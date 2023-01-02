// test-history is a test to be run on the wppsvr system after making changes
// to it.  It runs through the entire saved history of past practice sessions,
// and makes sure that the current wppsvr code would produce the same results.
// Any discrepancies are flagged.
package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/rothskeller/packet/wppsvr/analyze"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/store"
	_ "github.com/rothskeller/packet/xscmsg/all" // register message types
)

func main() {
	var (
		st  *store.Store
		err error
	)
	if err = config.Read(analyze.KnownProblems()); err != nil {
		log.Fatal(err)
	}
	if st, err = store.Open(); err != nil {
		log.Fatal(err)
	}
	log.SetOutput(io.Discard)
	for _, session := range st.GetRealizedSessions(time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local), time.Now().AddDate(0, 0, 8)) {
		if session.Imported {
			continue
		}
		testHistorySession(st, session)
	}
}

func testHistorySession(st *store.Store, session *store.Session) {
	for _, message := range st.GetSessionMessages(session.ID) {
		analysis := analyze.Analyze(storeStub{}, session, message.ToBBS, message.Message)
		compareAnalyses(message, analysis)
	}
}

func compareAnalyses(m *store.Message, a *analyze.Analysis) {
	if m.MessageType != a.MessageType() {
		fmt.Printf("%s: MessageType %q => %q\n", m.LocalID, m.MessageType, a.MessageType())
	}
	if m.FromCallSign != a.FromCallSign {
		fmt.Printf("%s: FromCallSign %q => %q\n", m.LocalID, m.FromCallSign, a.FromCallSign)
	}
	var aj string
	if a.Practice != nil {
		aj = a.Practice.Jurisdiction
	}
	if m.Jurisdiction != aj {
		fmt.Printf("%s: Jurisdiction %q => %q\n", m.LocalID, m.Jurisdiction, aj)
	}
	problems, actions := a.ProblemsActions()
	if m.Actions != actions {
		fmt.Printf("%s: Actions %x => %x\n", m.LocalID, m.Actions, actions)
	}
	mplist := strings.Join(m.Problems, ", ")
	aplist := strings.Join(problems, ", ")
	if mplist != aplist {
		fmt.Printf("%s: Problems %s => %s\n", m.LocalID, mplist, aplist)
	}
}

// The storeStub that we pass to Analyze never reports any message as having
// been seen before; never saves anything; and returns "XXX-100P" as the local
// message ID for all messages.
type storeStub struct{}

func (storeStub) HasMessageHash(string) string { return "" }
func (storeStub) NextMessageID(string) string  { return "XXX-100P" }
func (storeStub) SaveMessage(*store.Message)   {}
