package main

import (
	"log"
	"time"

	"steve.rothskeller.net/packet/wppsvr/config"
	"steve.rothskeller.net/packet/wppsvr/report"
	"steve.rothskeller.net/packet/wppsvr/retrieve"
	"steve.rothskeller.net/packet/wppsvr/store"
)

// closeSessions closes any sessions that are past their end time and sends
// reports for them.
func closeSessions(st *store.Store) {
	var now = time.Now()

	for _, session := range st.GetRunningSessions() {
		if session.End.Before(now) {
			if len(session.ReportTo) != 0 || st.SessionHasMessages(session.ID) {
				session.Running = false
				st.UpdateSession(session)
				log.Printf("Ended session for %s ending %s.", session.Name, session.End.Format("2006-01-02 15:04"))
				var conn = retrieve.ConnectToBBS(session.ToBBSes[0], session.CallSign)
				report.Send(st, conn, session)
				conn.Close()
			} else {
				st.DeleteSession(session.ID)
				log.Printf("Ended and deleted empty session for %s ending %s.",
					session.Name, session.End.Format("2006-01-02 15:04"))
			}
		}
	}
}

// createSessions starts new session instances for any defined sessions that
// don't have a running instance.  This includes both the sessions that were
// closed moments ago by closeSessions() as well as any new sessions added to
// the config.
func createSessions(st *store.Store) {
	var running = make(map[string]bool)

	for _, session := range st.GetRunningSessions() {
		running[session.CallSign] = true
	}
	for callSign, sessionParams := range config.Get().Sessions {
		if !running[callSign] {
			createSession(st, callSign, sessionParams)
		}
	}
}

// createSession starts a new session with the specified parameters.
func createSession(st *store.Store, callSign string, params *config.SessionConfig) {
	var (
		session store.Session
		now     = time.Now()
	)
	session.CallSign = callSign
	session.Name = params.Name
	if params.StartInterval.Match(now) {
		session.Start = now
	} else {
		session.Start = params.StartInterval.Prev(now)
	}
	session.End = params.EndInterval.Next(session.Start)
	session.ReportTo = params.ReportTo
	session.GenerateWeekSummary = params.GenerateWeekSummary
	session.ExcludeFromWeekSummary = params.ExcludeFromWeekSummary
	session.ToBBSes = params.ToBBSes.AllFor(session.End)
	session.DownBBSes = params.DownBBSes.AllFor(session.End)
	session.RetrieveFromBBSes = params.RetrieveFromBBSes
	session.RetrieveAt, session.RetrieveAtInterval = params.RetrieveAt, params.RetrieveAtInterval
	session.MessageTypes = params.MessageTypes.AllFor(session.End)
	session.Running = true
	st.CreateSession(&session)
	log.Printf("Started session for %s ending %s.", session.Name, session.End.Format("2006-01-02 15:04"))

	// Log any problems with the scheduled parts of the session config.
	if len(session.ToBBSes) == 0 {
		log.Printf("ERROR: %s has no valid destination BBSes", callSign)
	}
	for _, down := range session.DownBBSes {
		for _, to := range session.ToBBSes {
			if down == to {
				log.Printf("ERROR: %s lists %s as both down and valid", callSign, down)
			}
		}
	}
	if len(session.MessageTypes) == 0 {
		log.Printf("ERROR: %s has no valid message types", callSign)
	}
}
