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

	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/interval"
)

// A Session defines the parameters of a single session instance.
type Session struct {
	ID              int          `yaml:"id"`
	CallSign        string       `yaml:"callSign"`
	Name            string       `yaml:"name"`
	Prefix          string       `yaml:"prefix"`
	Start           time.Time    `yaml:"start"`
	End             time.Time    `yaml:"end"`
	ExcludeFromWeek bool         `yaml:"-"`
	ReportTo        []string     `yaml:"-"`
	ToBBSes         []string     `yaml:"toBBSes"`
	DownBBSes       []string     `yaml:"downBBSes"`
	Retrieve        []*Retrieval `yaml:"retrieve"`
	MessageTypes    []string     `yaml:"messageTypes"`
	Modified        bool         `yaml:"-"`
	Running         bool         `yaml:"-"`
	Imported        bool         `yaml:"-"`
	Report          string       `yaml:"-"`
}

// A Retrieval describes a single scheduled retrieval for a session.
type Retrieval struct {
	When              string            `yaml:"interval"`
	BBS               string            `yaml:"bbs"`
	Mailbox           string            `yaml:"mailbox"`
	DontKillMessages  bool              `yaml:"dontKillMessages"`
	DontSendResponses bool              `yaml:"dontSendResponses"`
	LastRun           time.Time         `yaml:"lastRun"`
	Interval          interval.Interval `yaml:"-"`
}

// GetRunningSessions returns the (unordered) list of all running sessions.
func (s *Store) GetRunningSessions() (list []*Session) {
	// Running sessions are always realized in the database, because the act
	// of setting their running flag causes them to be realized.  So this is
	// just a database query.
	return s.getSessionsWhere("running")
}

// ExistRealizedSessions returns whether any realized sessions exist in the
// specified time range (inclusive start, exclusive end).
func (s *Store) ExistRealizedSessions(start, end time.Time) bool {
	var (
		dummy int
	)
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	switch err := s.dbh.QueryRow(`SELECT 1 FROM session WHERE end>=? AND end <? LIMIT 1`, start, end).Scan(&dummy); err {
	case nil:
		return true
	case sql.ErrNoRows:
		return false
	default:
		panic(err)
	}
}

// GetSessions returns the set of sessions that end during the specified time
// range (inclusive start, exclusive end).  The sessions are sorted by end time,
// then by call sign.
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
	rows, err = s.dbh.Query("SELECT id, callsign, name, prefix, start, end, excludefromweek, reportto, tobbses, downbbses, messagetypes, modified, running, imported, report FROM session WHERE "+where, args...)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var (
			session      Session
			reportto     string
			tobbses      string
			downbbses    string
			messagetypes string
			rows2        *sql.Rows
		)
		err = rows.Scan(&session.ID, &session.CallSign, &session.Name, &session.Prefix, &session.Start, &session.End, &session.ExcludeFromWeek, &reportto, &tobbses, &downbbses, &messagetypes, &session.Modified, &session.Running, &session.Imported, &session.Report)
		if err != nil {
			panic(err)
		}
		session.ReportTo = split(reportto)
		session.ToBBSes = split(tobbses)
		session.DownBBSes = split(downbbses)
		session.MessageTypes = split(messagetypes)
		rows2, err = s.dbh.Query(`SELECT interval, bbs, mailbox, dontkillmessages, dontsendresponses, lastrun FROM retrieval WHERE session=?`, session.ID)
		if err != nil {
			panic(err)
		}
		for rows2.Next() {
			var r Retrieval

			if err = rows2.Scan(&r.When, &r.BBS, &r.Mailbox, &r.DontKillMessages, &r.DontSendResponses, &r.LastRun); err != nil {
				panic(err)
			}
			r.Interval = interval.Parse(r.When)
			session.Retrieve = append(session.Retrieve, &r)
		}
		if err = rows2.Err(); err != nil {
			panic(err)
		}
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
		tx     *sql.Tx
		result sql.Result
		id     int64
		err    error
	)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	tx, err = s.dbh.Begin()
	result, err = tx.Exec("INSERT INTO session (callsign, name, prefix, start, end, excludefromweek, reportto, tobbses, downbbses, messagetypes, modified, running, imported, report) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		session.CallSign, session.Name, session.Prefix, session.Start, session.End, session.ExcludeFromWeek,
		strings.Join(session.ReportTo, ";"), strings.Join(session.ToBBSes, ";"),
		strings.Join(session.DownBBSes, ";"), strings.Join(session.MessageTypes, ";"),
		session.Modified, session.Running, session.Imported, session.Report)
	if err != nil {
		panic(err)
	}
	id, err = result.LastInsertId()
	if err != nil {
		panic(err)
	}
	session.ID = int(id)
	for _, r := range session.Retrieve {
		_, err = tx.Exec("INSERT INTO retrieval (session, interval, bbs, mailbox, dontkillmessages, dontsendresponses, lastrun) VALUES (?,?,?,?,?,?,?)",
			session.ID, r.When, r.BBS, r.Mailbox, r.DontKillMessages, r.DontSendResponses, r.LastRun)
		if err != nil {
			panic(err)
		}
	}
	if err = tx.Commit(); err != nil {
		panic(err)
	}
}

// UpdateSession updates an existing session.
func (s *Store) UpdateSession(session *Session) {
	var (
		tx  *sql.Tx
		err error
	)
	if session.ID == 0 {
		// This is an unrealized session; we actually need to create it.
		s.CreateSession(session)
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if tx, err = s.dbh.Begin(); err != nil {
		panic(err)
	}
	_, err = tx.Exec("UPDATE session SET (callsign, name, prefix, start, end, excludefromweek, reportto, tobbses, downbbses, messagetypes, modified, running, imported, report) = (?,?,?,?,?,?,?,?,?,?,?,?,?,?) WHERE id=?",
		session.CallSign, session.Name, session.Prefix, session.Start, session.End, session.ExcludeFromWeek,
		strings.Join(session.ReportTo, ";"), strings.Join(session.ToBBSes, ";"),
		strings.Join(session.DownBBSes, ";"), strings.Join(session.MessageTypes, ";"),
		session.Modified, session.Running, session.Imported, session.Report, session.ID)
	if err != nil {
		panic(err)
	}
	if _, err = tx.Exec(`DELETE FROM retrieval WHERE session=?`, session.ID); err != nil {
		panic(err)
	}
	for _, r := range session.Retrieve {
		_, err = tx.Exec("INSERT INTO retrieval (session, interval, bbs, mailbox, dontkillmessages, dontsendresponses, lastrun) VALUES (?,?,?,?,?,?,?)",
			session.ID, r.When, r.BBS, r.Mailbox, r.DontKillMessages, r.DontSendResponses, r.LastRun)
		if err != nil {
			panic(err)
		}
	}
	if err = tx.Commit(); err != nil {
		panic(err)
	}
}

// getConfiguredSessions returns the sessions defined by the configuration that
// end during the specified time range (inclusive start, exclusive end).  They
// may or may not be realized in the database; the ID fields are not filled in.
func getConfiguredSessions(start, end time.Time) (list []*Session) {
	for callSign, params := range config.Get().Sessions {
		var sessend time.Time

		if sessend = start.Truncate(time.Minute); !start.Equal(sessend) {
			sessend = sessend.Add(time.Minute)
		}
		for ; sessend.Before(end); sessend = sessend.Add(time.Minute) {
			var session Session

			if !params.EndInterval.Match(sessend) {
				continue
			}
			session.CallSign = callSign
			session.Name = params.Name
			session.Prefix = params.Prefix
			session.End = sessend
			for session.Start = sessend.Add(-time.Minute); !params.StartInterval.Match(session.Start); session.Start = session.Start.Add(-time.Minute) {
			}
			session.ReportTo = params.ReportTo
			session.ExcludeFromWeek = params.ExcludeFromWeek
			session.ToBBSes = params.ToBBSes.AllFor(session.End)
			session.DownBBSes = params.DownBBSes.AllFor(session.End)
			session.MessageTypes = params.MessageTypes.AllFor(session.End)
			session.Retrieve = make([]*Retrieval, len(params.Retrieve))
			for i, pr := range params.Retrieve {
				session.Retrieve[i] = &Retrieval{
					When:              pr.When,
					BBS:               pr.BBS,
					Mailbox:           pr.Mailbox,
					DontKillMessages:  pr.DontKillMessages,
					DontSendResponses: pr.DontSendResponses,
					Interval:          pr.Interval,
				}
			}
			list = append(list, &session)
		}
	}
	return list
}

func (s *Store) DeleteSession(session *Session) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, err := s.dbh.Exec(`DELETE FROM session WHERE id=?`, session.ID); err != nil {
		panic(err)
	}

}
