// reanalyze removes all of the analyzed messages from a session, and then
// re-analyzes and re-adds them.  It then re-generates the report for the
// session.  No response or problem messages are actually sent, and the report
// is not sent either.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rothskeller/packet/message/allmsg"
	"github.com/rothskeller/packet/wppsvr/analyze"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/report"
	"github.com/rothskeller/packet/wppsvr/store"
)

func main() {
	var (
		st       *store.Store
		session  *store.Session
		messages []*store.Message
		err      error
	)
	allmsg.Register()
	if err = config.Read(); err != nil {
		log.Fatal(err)
	}
	if st, err = store.Open(); err != nil {
		log.Fatal(err)
	}
	if len(os.Args) == 2 {
		date, _ := time.Parse("2006-01-02", os.Args[1])
		if !date.IsZero() {
			sessions := st.GetSessions(
				time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local),
				time.Date(date.Year(), date.Month(), date.Day()+1, 0, 0, 0, 0, time.Local),
			)
			if len(sessions) == 1 {
				session = sessions[0]
			}
		}
	}
	if session == nil {
		fmt.Fprintf(os.Stderr, "usage: reanalyze session-end-date\n")
		os.Exit(2)
	}
	messages = st.GetSessionMessages(session.ID)
	for _, message := range messages {
		var fs = &filteredStore{st: st, id: message.LocalID}
		analysis := analyze.Analyze(fs, session, message.ToBBS, message.Message)
		analysis.Commit(fs)
	}
	if session.Flags&store.Running == 0 && session.Report != "" {
		rpt := report.Generate(st, session)
		session.Report = rpt.RenderPlainText()
		st.UpdateSession(session)
	}
}

type filteredStore struct {
	st *store.Store
	id string
}

func (fs *filteredStore) HasMessageHash(string) string {
	return "" // we never already have the message.
}
func (fs *filteredStore) NextMessageID(string) string {
	return fs.id // return the ID that the message was given originally
}
func (fs *filteredStore) SaveMessage(m *store.Message) {
	fs.SaveMessage(m)
}
