package store

import (
	"database/sql"
	"strings"
	"time"

	"steve.rothskeller.net/packet/wppsvr/config"
)

// A Session defines the parameters of a single session instance.
type Session struct {
	ID                     int                `yaml:"id"`
	CallSign               string             `yaml:"callSign"`
	Name                   string             `yaml:"name"`
	Prefix                 string             `yaml:"prefix"`
	Start                  time.Time          `yaml:"start"`
	End                    time.Time          `yaml:"end"`
	GenerateWeekSummary    bool               `yaml:"-"`
	ExcludeFromWeekSummary bool               `yaml:"-"`
	ReportTo               []string           `yaml:"-"`
	ToBBSes                []string           `yaml:"toBBSes"`
	DownBBSes              []string           `yaml:"downBBSes"`
	RetrieveFromBBSes      []string           `yaml:"-"`
	RetrieveAt             []string           `yaml:"-"`
	RetrieveAtInterval     []*config.Interval `yaml:"-"`
	MessageTypes           []string           `yaml:"messageTypes"`
	Modified               bool               `yaml:"-"`
	Running                bool               `yaml:"-"`
	Report                 string             `yaml:"-"`
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

// PreviousSession returns the session immediately preceding the supplied
// session (and with the same call sign).  It returns nil if there is none.
func (s *Store) PreviousSession(from *Session) *Session {
	list := s.getSessionsWhere("callsign=? and end<? ORDER BY end DESC LIMIT 1", from.CallSign, from.End)
	if len(list) != 0 {
		return list[0]
	}
	return nil
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
			retrieveats   string
			messagetypes  string
		)
		err = rows.Scan(&session.ID, &session.CallSign, &session.Name, &session.Prefix, &session.Start, &session.End, &session.GenerateWeekSummary, &session.ExcludeFromWeekSummary, &reportto, &tobbses, &downbbses, &retrievebbses, &retrieveats, &messagetypes, &session.Modified, &session.Running, &session.Report)
		if err != nil {
			panic(err)
		}
		session.ReportTo = split(reportto)
		session.ToBBSes = split(tobbses)
		session.DownBBSes = split(downbbses)
		session.RetrieveFromBBSes = split(retrievebbses)
		session.RetrieveAt = split(retrieveats)
		session.RetrieveAtInterval = make([]*config.Interval, len(session.RetrieveAt))
		for i, ra := range session.RetrieveAt {
			session.RetrieveAtInterval[i], _ = config.ParseInterval(ra)
		}
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
	result, err = s.dbh.Exec("INSERT INTO session (callsign, name, prefix, start, end, generateweeksummary, excludefromweeksummary, reportto, tobbses, downbbses, retrievefrombbses, retrieveat, messagetypes, modified, running, report) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		session.CallSign, session.Name, session.Prefix, session.Start, session.End,
		session.GenerateWeekSummary, session.ExcludeFromWeekSummary,
		strings.Join(session.ReportTo, ";"), strings.Join(session.ToBBSes, ";"),
		strings.Join(session.DownBBSes, ";"), strings.Join(session.RetrieveFromBBSes, ";"),
		strings.Join(session.RetrieveAt, ";"), strings.Join(session.MessageTypes, ";"),
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
	_, err = s.dbh.Exec("UPDATE session SET (callsign, name, prefix, start, end, generateweeksummary, excludefromweeksummary, reportto, tobbses, downbbses, retrievefrombbses, retrieveat, messagetypes, modified, running, report) = (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) WHERE id=?",
		session.CallSign, session.Name, session.Prefix, session.Start, session.End,
		session.GenerateWeekSummary, session.ExcludeFromWeekSummary,
		strings.Join(session.ReportTo, ";"), strings.Join(session.ToBBSes, ";"),
		strings.Join(session.DownBBSes, ";"), strings.Join(session.RetrieveFromBBSes, ";"),
		strings.Join(session.RetrieveAt, ";"), strings.Join(session.MessageTypes, ";"),
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
