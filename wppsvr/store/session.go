package store

import (
	"database/sql"
	"strings"
	"time"

	"steve.rothskeller.net/packet/wppsvr/config"
)

// A Session defines the parameters of a single session instance.
type Session struct {
	ID                     int
	CallSign               string
	Name                   string
	Prefix                 string
	Start                  time.Time
	End                    time.Time
	GenerateWeekSummary    bool
	ExcludeFromWeekSummary bool
	ReportTo               []string
	ToBBSes                []string
	DownBBSes              []string
	RetrieveFromBBSes      []string
	RetrieveAt             string
	RetrieveAtInterval     *config.Interval
	MessageTypes           []string
	Modified               bool
	Running                bool
	Report                 string
}

// GetRunningSessions returns the (unordered) list of all running sessions.
func (s *Store) GetRunningSessions() (list []*Session) {
	return s.getSessionsWhere("running")
}

// GetSessionsEnding returns the (unordered) list of sessions that end during
// the specified time range (inclusive start, exclusive end).
func (s *Store) GetSessionsEnding(start, end time.Time) (list []*Session) {
	return s.getSessionsWhere("end>=? AND end<?", start, end)
}

// getSessionsWhere returns the (unordered) list of sessions matching the
// specified criteria.
func (s *Store) getSessionsWhere(where string, args ...interface{}) (list []*Session) {
	var (
		rows *sql.Rows
		err  error
	)
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	rows, err = s.dbh.Query("SELECT id, callsign, name, prefix, start, end, generateweeksummary, excludefromweeksummary, reportto, tobbses, downbbses, retrievefrombbses, retrieveat, messagetypes, modified, running, report FROM session WHERE "+where, args...)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var (
			session       Session
			reportto      string
			tobbses       string
			downbbses     string
			retrievebbses string
			messagetypes  string
		)
		err = rows.Scan(&session.ID, &session.CallSign, &session.Name, &session.Prefix, &session.Start, &session.End, &session.GenerateWeekSummary, &session.ExcludeFromWeekSummary, &reportto, &tobbses, &downbbses, &retrievebbses, &session.RetrieveAt, &messagetypes, &session.Modified, &session.Running, &session.Report)
		if err != nil {
			panic(err)
		}
		session.ReportTo = split(reportto)
		session.ToBBSes = split(tobbses)
		session.DownBBSes = split(downbbses)
		session.RetrieveFromBBSes = split(retrievebbses)
		session.RetrieveAtInterval, _ = config.ParseInterval(session.RetrieveAt)
		session.MessageTypes = split(messagetypes)
		list = append(list, &session)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}
	return list
}

// CreateSession creates a new session.
func (s *Store) CreateSession(session *Session) {
	var (
		result sql.Result
		id     int64
		err    error
	)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	result, err = s.dbh.Exec("INSERT INTO session (callsign, name, prefix, start, end, generateweeksummary, excludefromweeksummary, reportto, tobbses, downbbses, retrievefrombbses, retrieveat, messagetypes, modified, running, report) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		session.CallSign, session.Name, session.Prefix, session.Start, session.End,
		session.GenerateWeekSummary, session.ExcludeFromWeekSummary,
		strings.Join(session.ReportTo, ","), strings.Join(session.ToBBSes, ","),
		strings.Join(session.DownBBSes, ","), strings.Join(session.RetrieveFromBBSes, ","),
		session.RetrieveAt, strings.Join(session.MessageTypes, ","),
		session.Modified, session.Running, session.Report)
	if err != nil {
		panic(err)
	}
	id, err = result.LastInsertId()
	if err != nil {
		panic(err)
	}
	session.ID = int(id)
}

// UpdateSession updates an existing session.
func (s *Store) UpdateSession(session *Session) {
	var err error

	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, err = s.dbh.Exec("UPDATE session SET (callsign, name, prefix, start, end, generateweeksummary, excludefromweeksummary, reportto, tobbses, downbbses, retrievefrombbses, retrieveat, messagetypes, modified, running, report) = (?,?,?,?,?,?,?,?,?,?,?,?,?,?) WHERE id=?",
		session.CallSign, session.Name, session.Prefix, session.Start, session.End,
		session.GenerateWeekSummary, session.ExcludeFromWeekSummary,
		strings.Join(session.ReportTo, ","), strings.Join(session.ToBBSes, ","),
		strings.Join(session.DownBBSes, ","), strings.Join(session.RetrieveFromBBSes, ","),
		session.RetrieveAt, strings.Join(session.MessageTypes, ","),
		session.Modified, session.Running, session.Report, session.ID)
	if err != nil {
		panic(err)
	}
}

// DeleteSession deletes a session.
func (s *Store) DeleteSession(sessionID int) {
	var err error

	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, err = s.dbh.Exec("DELETE FROM session WHERE id=?", sessionID)
	if err != nil {
		panic(err)
	}
}
