// Package retrieve handles connecting to BBSes, retrieving messages from them,
// and sending the responses to them.
package retrieve

import (
	"log"
	"sync"
	"time"

	"steve.rothskeller.net/packet/jnos"
	"steve.rothskeller.net/packet/jnos/kpc3plus"
	"steve.rothskeller.net/packet/jnos/telnet"
	"steve.rothskeller.net/packet/wppsvr/analyze"
	"steve.rothskeller.net/packet/wppsvr/config"
	"steve.rothskeller.net/packet/wppsvr/store"
)

// ForRunningSessions retrieves and responds to new messages in all running
// practice sessions.
func ForRunningSessions(st *store.Store) {
	var (
		wg  sync.WaitGroup
		now = time.Now()
	)
	for _, session := range st.GetRunningSessions() {
		var lastRetrievalPoint time.Time

		for _, rai := range session.RetrieveAtInterval {
			point := rai.Prev(now)
			if point.After(lastRetrievalPoint) {
				lastRetrievalPoint = point
			}
		}
		for _, bbs := range session.RetrieveFromBBSes {
			if st.GetLastRetrieval(session.CallSign, bbs).Before(lastRetrievalPoint) {
				wg.Add(1)
				go checkBBS(st, &wg, session, bbs)
			}
		}
	}
	wg.Wait()
}

// ForSession retrieves and responds to new messages in the specified practice
// session.
func ForSession(st *store.Store, session *store.Session) {
	var (
		wg                 sync.WaitGroup
		lastRetrievalPoint time.Time
		now                = time.Now()
	)
	for _, rai := range session.RetrieveAtInterval {
		point := rai.Prev(now)
		if point.After(lastRetrievalPoint) {
			lastRetrievalPoint = point
		}
	}
	for _, bbs := range session.RetrieveFromBBSes {
		if st.GetLastRetrieval(session.CallSign, bbs).Before(lastRetrievalPoint) {
			wg.Add(1)
			go checkBBS(st, &wg, session, bbs)
		}
	}
	wg.Wait()
}

// checkBBS retrieves and responds to new check-in messages on a specific BBS.
func checkBBS(st *store.Store, wg *sync.WaitGroup, session *store.Session, bbsname string) {
	var (
		conn   *jnos.Conn
		err    error
		msgnum = 1
		start  = time.Now()
	)
	defer wg.Done()
	if conn = ConnectToBBS(bbsname, session.CallSign); conn == nil {
		return
	}
	defer func() {
		if err = conn.Close(); err != nil {
			log.Printf("ERROR: closing connection to %s@%s: %s", session.CallSign, bbsname, err)
		}
	}()
	for {
		var message string

		if message, err = conn.Read(msgnum); err != nil {
			log.Printf("ERROR: reading messages to %s@%s: %s", session.CallSign, bbsname, err)
			return
		} else if message == "" { // no more messages
			break
		}
		handleMessage(st, conn, session, bbsname, message, msgnum)
		msgnum++
	}
	st.SetLastRetrieval(session.CallSign, bbsname, start)
}

// ConnectToBBS connects to the specified mailbox on the specified BBS, in the
// manner dictated by the BBS configuration.
func ConnectToBBS(bbsname, mailbox string) (conn *jnos.Conn) {
	// This function is exported because it is also used by
	// wppsvr/sessions.go to connect to the BBS to send end-of-session
	// reports.
	var (
		bbs *config.BBSConfig
		err error
	)
	bbs = config.Get().BBSes[bbsname]
	switch bbs.Transport {
	case "kpc3plus":
		conn, err = kpc3plus.Connect("/dev/tty.usbserial-1410", bbs.AX25, mailbox, "KC6RSC")
	case "telnet":
		conn, err = telnet.Connect(bbs.TCP, mailbox, bbs.Passwords[mailbox])
	}
	if err != nil {
		log.Printf("ERROR: can't connect to %s@%s via %s: %s", mailbox, bbsname, bbs.Transport, err)
		return nil
	}
	return conn
}

// handleMessage handles a single incoming message.
func handleMessage(st *store.Store, conn *jnos.Conn, session *store.Session, bbsname, message string, msgnum int) {
	var (
		analysis  *analyze.Analysis
		responses []*store.Response
		err       error
	)
	analysis = analyze.Analyze(st, session, bbsname, message)
	responses = analysis.Responses(st)
	for _, response := range responses {
		if err = conn.Send(response.Subject, response.Body, response.To); err != nil {
			log.Printf("ERROR: sending message from %s@%s: %s", session.CallSign, bbsname, err)
			return
		}
		response.SendTime = time.Now()
	}
	if err = conn.Kill(msgnum); err != nil {
		log.Printf("ERROR: killing message %d at %s@%s: %s", msgnum, session.CallSign, bbsname, err)
		return
	}
	analysis.Commit(st)
	for _, response := range responses {
		st.SaveResponse(response)
	}
}
