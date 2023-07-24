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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/wppsvr/interval"
)

// A Session defines the parameters of a single session instance.
type Session struct {
	ID           int          `yaml:"id"`
	CallSign     string       `yaml:"callSign"`
	Name         string       `yaml:"name"`
	Prefix       string       `yaml:"prefix"`
	Start        time.Time    `yaml:"start"`
	End          time.Time    `yaml:"end"`
	ReportToText []string     `yaml:"-"`
	ReportToHTML []string     `yaml:"-"`
	ToBBSes      []string     `yaml:"toBBSes"`
	DownBBSes    []string     `yaml:"downBBSes"`
	Retrieve     []*Retrieval `yaml:"retrieve"`
	MessageTypes []string     `yaml:"messageTypes"`
	ModelMessage string       `yaml:"modelMessage"`
	Instructions string       `yaml:"instructions"`
	RetrieveAt   string       `yaml:"retrieveAt"`
	Report       string       `yaml:"-"`
	Flags        SessionFlags `yaml:"flags"`

	ModelMsg         message.ICompare  `yaml:"-"`
	RetrieveInterval interval.Interval `yaml:"-"`
}

// A Retrieval describes a single scheduled retrieval for a session.
type Retrieval struct {
	BBS     string    `yaml:"bbs"`
	LastRun time.Time `yaml:"lastRun"`
}

// SessionFlags is a collection of flags describing a session.
type SessionFlags uint8

// Values for SessionFlags
const (
	Running SessionFlags = (1 << iota)
	ExcludeFromWeek
	DontKillMessages
	DontSendResponses
	Imported
	Modified
	ReportToSenders
)

// GetRunningSessions returns the (unordered) list of all running sessions.
func (s *Store) GetRunningSessions() (list []*Session) {
	// Running sessions are always realized in the database, because the act
	// of setting their running flag causes them to be realized.  So this is
	// just a database query.
	return s.getSessionsWhere("flags&1")
}

// ExistSessions returns whether any sessions exist in the specified time range
// (inclusive start, exclusive end).
func (s *Store) ExistSessions(start, end time.Time) bool {
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

// GetSession returns the session with the specified ID, or nil if there is
// none.
func (s *Store) GetSession(id int) *Session {
	if list := s.getSessionsWhere("id=?", id); len(list) != 0 {
		return list[0]
	}
	return nil
}

// GetSessions returns the set of sessions that end during the specified
// time range (inclusive start, exclusive end).  The sessions are sorted by end
// time, then by call sign.
func (s *Store) GetSessions(start, end time.Time) (list []*Session) {
	return s.getSessionsWhere("end>=? AND end<? ORDER BY end, callsign", start, end)
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
	rows, err = s.dbh.Query("SELECT id, callsign, name, prefix, start, end, reporttotext, reporttohtml, tobbses, downbbses, messagetypes, modelmessage, instructions, retrieveat, report, flags FROM session WHERE "+where, args...)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var (
			session      Session
			reporttotext string
			reporttohtml string
			tobbses      string
			downbbses    string
			messagetypes string
			rows2        *sql.Rows
		)
		err = rows.Scan(&session.ID, &session.CallSign, &session.Name, &session.Prefix, &session.Start, &session.End, &reporttotext, &reporttohtml, &tobbses, &downbbses, &messagetypes, &session.ModelMessage, &session.Instructions, &session.RetrieveAt, &session.Report, &session.Flags)
		if err != nil {
			panic(err)
		}
		session.ReportToText = split(reporttotext)
		session.ReportToHTML = split(reporttohtml)
		session.ToBBSes = split(tobbses)
		session.DownBBSes = split(downbbses)
		session.MessageTypes = split(messagetypes)
		session.RetrieveInterval = interval.Parse(session.RetrieveAt)
		if session.ModelMessage != "" {
			if env, body, err := envelope.ParseSaved(session.ModelMessage); err == nil {
				session.ModelMsg = message.Decode(env.SubjectLine, body).(message.ICompare)
			} else {
				panic(err)
			}
		}
		rows2, err = s.dbh.Query(`SELECT bbs, lastrun FROM retrieval WHERE session=?`, session.ID)
		if err != nil {
			panic(err)
		}
		for rows2.Next() {
			var r Retrieval

			if err = rows2.Scan(&r.BBS, &r.LastRun); err != nil {
				panic(err)
			}
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
	result, err = tx.Exec("INSERT INTO session (callsign, name, prefix, start, end, reporttotext, reporttohtml, tobbses, downbbses, messagetypes, modelmessage, instructions, retrieveat, report, flags) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		session.CallSign, session.Name, session.Prefix, session.Start, session.End,
		strings.Join(session.ReportToText, ";"), strings.Join(session.ReportToHTML, ";"), strings.Join(session.ToBBSes, ";"),
		strings.Join(session.DownBBSes, ";"), strings.Join(session.MessageTypes, ";"), session.ModelMessage,
		session.Instructions, session.RetrieveAt, session.Report, session.Flags)
	if err != nil {
		panic(err)
	}
	id, err = result.LastInsertId()
	if err != nil {
		panic(err)
	}
	session.ID = int(id)
	for _, r := range session.Retrieve {
		_, err = tx.Exec("INSERT INTO retrieval (session, bbs, lastrun) VALUES (?,?,?)",
			session.ID, r.BBS, r.LastRun)
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
	_, err = tx.Exec("UPDATE session SET (callsign, name, prefix, start, end, reporttotext, reporttohtml, tobbses, downbbses, messagetypes, modelmessage, instructions, retrieveat, report, flags) = (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) WHERE id=?",
		session.CallSign, session.Name, session.Prefix, session.Start, session.End,
		strings.Join(session.ReportToText, ";"), strings.Join(session.ReportToHTML, ";"), strings.Join(session.ToBBSes, ";"),
		strings.Join(session.DownBBSes, ";"), strings.Join(session.MessageTypes, ";"), session.ModelMessage,
		session.Instructions, session.RetrieveAt, session.Report, session.Flags, session.ID)
	if err != nil {
		panic(err)
	}
	if _, err = tx.Exec(`DELETE FROM retrieval WHERE session=?`, session.ID); err != nil {
		panic(err)
	}
	for _, r := range session.Retrieve {
		_, err = tx.Exec("INSERT INTO retrieval (session, bbs, lastrun) VALUES (?,?,?)",
			session.ID, r.BBS, r.LastRun)
		if err != nil {
			panic(err)
		}
	}
	if err = tx.Commit(); err != nil {
		panic(err)
	}
}

// DeleteSession deletes a session.
func (s *Store) DeleteSession(session *Session) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, err := s.dbh.Exec(`DELETE FROM session WHERE id=?`, session.ID); err != nil {
		panic(err)
	}
}

// ModelImageCount returns the number of model images associated with the
// session.  The images use 1-based numbering, so they are numbered 1 through
// the return value of this function, inclusive.
func (s *Store) ModelImageCount(sid int) (count int) {
	prefix := fmt.Sprintf("s%d", sid)
	matches, _ := filepath.Glob(prefix + "p*.*")
	for _, match := range matches {
		pstr := match[len(prefix)+1 : len(match)-len(filepath.Ext(match))]
		if pnum, err := strconv.Atoi(pstr); err == nil && pnum > count {
			count = pnum
		}
	}
	return count
}

// ModelImage returns an open file handle to the specified model image page
// number, or nil if there is no such image.  Model image page numbers start at
// 1.  It is the caller's responsibility to close the handle.
func (s *Store) ModelImage(sid int, pnum int) (fh *os.File) {
	matches, _ := filepath.Glob(fmt.Sprintf("s%dp%d.*", sid, pnum))
	if len(matches) == 1 {
		fh, _ = os.Open(matches[0])
	}
	return fh
}

// DeleteModelImages removes all model images for the specified session.
func (s *Store) DeleteModelImages(sid int) {
	prefix := fmt.Sprintf("s%d", sid)
	matches, _ := filepath.Glob(prefix + "p*.*")
	for _, match := range matches {
		os.Remove(match)
	}
}

// SaveModelImage saves the specified model image for the specified session.
func (s *Store) SaveModelImage(sid int, pnum int, name string, body io.Reader) {
	fname := fmt.Sprintf("s%dp%d%s", sid, pnum, filepath.Ext(name))
	if fh, err := os.Create(fname); err == nil {
		io.Copy(fh, body)
		fh.Close()
	}
}
