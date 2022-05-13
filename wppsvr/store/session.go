package store

/*
Sessions are unusual because they aren't all stored as rows in the "session"
database table.  The set of sessions exposed by the "store" package is the union
of those stored in the database table and the set of future sessions defined by
the configuration.  The fact that some of them are not realized in the database
is an internal implementation detail.
*/

import (
	"database/sql"
	"sort"
	"strings"
	"time"

	"steve.rothskeller.net/packet/wppsvr/config"
)

// A Session defines the parameters of a single session instance.
type Session struct {
	ID                 int                `yaml:"id"`
	CallSign           string             `yaml:"callSign"`
	Name               string             `yaml:"name"`
	Prefix             string             `yaml:"prefix"`
	Start              time.Time          `yaml:"start"`
	End                time.Time          `yaml:"end"`
	ExcludeFromWeek    bool               `yaml:"-"`
	ReportTo           []string           `yaml:"-"`
	ToBBSes            []string           `yaml:"toBBSes"`
	DownBBSes          []string           `yaml:"downBBSes"`
	RetrieveFromBBSes  []string           `yaml:"-"`
	RetrieveAt         []string           `yaml:"-"`
	RetrieveAtInterval []*config.Interval `yaml:"-"`
	MessageTypes       []string           `yaml:"messageTypes"`
	Modified           bool               `yaml:"-"`
	Running            bool               `yaml:"-"`
	Report             string             `yaml:"-"`
}

// GetRunningSessions returns the (unordered) list of all running sessions.
func (s *Store) GetRunningSessions() (list []*Session) {
	// Running sessions are always realized in the database, because the act
	// of setting their running flag causes them to be realized.  So this is
	// just a database query.
	return s.getSessionsWhere("running")
}

// GetSessions returns the set of sessions that end during the specified time
// range (inclusive start, exclusive end).
func (s *Store) GetSessions(start, end time.Time) (list []*Session) {
	var sched []*Session

	// Start by retrieving all realized sessions in the range from the
	// database.
	list = s.getSessionsWhere("end>=? AND end<? ORDER BY end, callsign", start, end)
	// If the end of the range is in the past, we're done.
	if !time.Now().Before(end) {
		return list
	}
	// Get the list of future sessions defined by the configuration for the
	// specified range.
	if start.Before(time.Now()) {
		start = time.Now()
	}
	sched = getConfiguredSessions(start, end)
	// Add those scheduled sessions to the list, but only where they don't
	// overlap sessions already in the list.
	for _, session := range sched {
		var overlap bool
		for _, listed := range list {
			if sessionOverlap(session, listed) {
				overlap = true
				break
			}
		}
		if !overlap {
			list = append(list, session)
		}
	}
	// Sort the resulting list.
	sort.Slice(list, func(i, j int) bool {
		if !list[i].End.Equal(list[j].End) {
			return list[i].End.Before(list[j].End)
		}
		return list[i].CallSign < list[j].CallSign
	})
	return list
}
func sessionOverlap(a, b *Session) bool {
	if a.CallSign != b.CallSign {
		return false
	}
	return a.Start.Before(b.End) && b.Start.Before(a.End)
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
	rows, err = s.dbh.Query("SELECT id, callsign, name, prefix, start, end, excludefromweek, reportto, tobbses, downbbses, retrievefrombbses, retrieveat, messagetypes, modified, running, report FROM session WHERE "+where, args...)
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
		err = rows.Scan(&session.ID, &session.CallSign, &session.Name, &session.Prefix, &session.Start, &session.End, &session.ExcludeFromWeek, &reportto, &tobbses, &downbbses, &retrievebbses, &retrieveats, &messagetypes, &session.Modified, &session.Running, &session.Report)
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
	result, err = s.dbh.Exec("INSERT INTO session (callsign, name, prefix, start, end, excludefromweek, reportto, tobbses, downbbses, retrievefrombbses, retrieveat, messagetypes, modified, running, report) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		session.CallSign, session.Name, session.Prefix, session.Start, session.End, session.ExcludeFromWeek,
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

	if session.ID == 0 {
		// This is an unrealized session; we actually need to create it.
		s.CreateSession(session)
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, err = s.dbh.Exec("UPDATE session SET (callsign, name, prefix, start, end, excludefromweek, reportto, tobbses, downbbses, retrievefrombbses, retrieveat, messagetypes, modified, running, report) = (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) WHERE id=?",
		session.CallSign, session.Name, session.Prefix, session.Start, session.End, session.ExcludeFromWeek,
		strings.Join(session.ReportTo, ";"), strings.Join(session.ToBBSes, ";"),
		strings.Join(session.DownBBSes, ";"), strings.Join(session.RetrieveFromBBSes, ";"),
		strings.Join(session.RetrieveAt, ";"), strings.Join(session.MessageTypes, ";"),
		session.Modified, session.Running, session.Report, session.ID)
	if err != nil {
		panic(err)
	}
}

// getConfiguredSessions returns the sessions defined by the configuration that
// end during the specified time range (inclusive start, exclusive end).  They
// may or may not be realized in the database; the ID fields are not filled in.
func getConfiguredSessions(start, end time.Time) (list []*Session) {
	for callSign, params := range config.Get().Sessions {
		sessend := params.EndInterval.Next(start.Add(-time.Second))
		for sessend.Before(end) {
			var (
				session Session
			)
			session.CallSign = callSign
			session.Name = params.Name
			session.Start = params.StartInterval.Prev(sessend)
			session.End = sessend
			session.ReportTo = params.ReportTo
			session.ExcludeFromWeek = params.ExcludeFromWeek
			session.ToBBSes = params.ToBBSes.AllFor(session.End)
			session.DownBBSes = params.DownBBSes.AllFor(session.End)
			session.RetrieveFromBBSes = params.RetrieveFromBBSes
			session.RetrieveAt, session.RetrieveAtInterval = params.RetrieveAt, params.RetrieveAtInterval
			session.MessageTypes = params.MessageTypes.AllFor(session.End)
			list = append(list, &session)
			sessend = params.EndInterval.Next(sessend)
		}
	}
	return list
}
