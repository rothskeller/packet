package store

/*
Sessions are unusual because they aren't all stored as rows in the "session"
database table.  The set of sessions exposed by the "store" package is the union
of those stored in the database table and the set of future sessions defined by
the configuration.  The fact that some of them are not realized in the database
is an internal implementation detail.
*/

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/wppsvr/db"
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

	ModelMsg         message.Message   `yaml:"-"`
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

const (
	startEndFormat = "2006-01-02 15:04:05-07:00"
	lastRunFormat  = "2006-01-02 15:04:05.999999999-07:00"
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
func (s *Store) ExistSessions(start, end time.Time) (found bool) {
	db.SQL(s.conn, "SELECT 1 FROM session WHERE end>=? AND end <? LIMIT 1", func(st *db.St) {
		st.BindTime(start, startEndFormat)
		st.BindTime(end, startEndFormat)
		found = st.Step()
	})
	return found
}

// OverlappingSession returns whether any sessions exist with an overlapping
// time range and the same call sign.  If allow is non-zero, a session with that
// ID doesn't count (generally the one we're trying to test overlap against).
func (s *Store) OverlappingSession(start, end time.Time, callsign string, allow int) (found bool) {
	db.SQL(s.conn, "SELECT 1 FROM session WHERE callsign=? AND start<=? AND ?<=end AND id!=?", func(st *db.St) {
		st.BindText(callsign)
		st.BindTime(end, startEndFormat)
		st.BindTime(start, startEndFormat)
		st.BindInt(allow)
		found = st.Step()
	})
	return found
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
	db.Transaction(s.conn, false, func() error {
		db.SQL(s.conn, "SELECT id, callsign, name, prefix, start, end, reporttotext, reporttohtml, tobbses, downbbses, messagetypes, modelmessage, instructions, retrieveat, report, flags FROM session WHERE "+where, func(st *db.St) {
			var err error

			for _, arg := range args {
				switch arg := arg.(type) {
				case int:
					st.BindInt(arg)
				case time.Time:
					st.BindTime(arg, startEndFormat)
				}
			}
			for st.Step() {
				var session Session

				session.ID = st.ColumnInt()
				session.CallSign = st.ColumnText()
				session.Name = st.ColumnText()
				session.Prefix = st.ColumnText()
				session.Start = st.ColumnTime(startEndFormat)
				session.End = st.ColumnTime(startEndFormat)
				session.ReportToText = split(st.ColumnText())
				session.ReportToHTML = split(st.ColumnText())
				session.ToBBSes = split(st.ColumnText())
				session.DownBBSes = split(st.ColumnText())
				session.MessageTypes = split(st.ColumnText())
				session.ModelMessage = st.ColumnText()
				session.Instructions = st.ColumnText()
				session.RetrieveAt = st.ColumnText()
				session.Report = st.ColumnText()
				session.Flags = SessionFlags(st.ColumnInt())
				session.RetrieveInterval = interval.Parse(session.RetrieveAt)
				if session.ModelMessage != "" {
					if session.ModelMsg, err = message.Parse(session.ModelMessage, "(database)"); err != nil {
						panic(err)
					}
				}
				db.SQL(s.conn, "SELECT bbs, lastrun FROM retrieval WHERE session=?", func(s2 *db.St) {
					s2.BindInt(session.ID)
					for s2.Step() {
						var r Retrieval

						r.BBS = s2.ColumnText()
						r.LastRun = s2.ColumnTime(lastRunFormat)
						session.Retrieve = append(session.Retrieve, &r)
					}
				})
				list = append(list, &session)
			}
		})
		return nil
	})
	return list
}

// CreateSession creates a new session.
func (s *Store) CreateSession(session *Session) {
	db.Transaction(s.conn, true, func() error {
		db.SQL(s.conn, "INSERT INTO session (callsign, name, prefix, start, end, reporttotext, reporttohtml, tobbses, downbbses, messagetypes, modelmessage, instructions, retrieveat, report, flags) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", func(st *db.St) {
			st.BindText(session.CallSign)
			st.BindText(session.Name)
			st.BindText(session.Prefix)
			st.BindTime(session.Start, startEndFormat)
			st.BindTime(session.End, startEndFormat)
			st.BindText(strings.Join(session.ReportToText, ";"))
			st.BindText(strings.Join(session.ReportToHTML, ";"))
			st.BindText(strings.Join(session.ToBBSes, ";"))
			st.BindText(strings.Join(session.DownBBSes, ";"))
			st.BindText(strings.Join(session.MessageTypes, ";"))
			st.BindText(session.ModelMessage)
			st.BindText(session.Instructions)
			st.BindText(session.RetrieveAt)
			st.BindText(session.Report)
			st.BindInt(int(session.Flags))
			st.Step()
		})
		session.ID = int(s.conn.LastInsertRowID())
		db.SQL(s.conn, "INSERT INTO retrieval (session, bbs, lastrun) VALUES (?,?,?)", func(st *db.St) {
			for _, r := range session.Retrieve {
				st.BindInt(session.ID)
				st.BindText(r.BBS)
				st.BindTime(r.LastRun, lastRunFormat)
				st.Step()
				st.Reset()
			}
		})
		return nil
	})
}

// UpdateSession updates an existing session.
func (s *Store) UpdateSession(session *Session) {
	if session.ID == 0 {
		// This is an unrealized session; we actually need to create it.
		s.CreateSession(session)
		return
	}
	db.Transaction(s.conn, true, func() error {
		db.SQL(s.conn, "UPDATE session SET (callsign, name, prefix, start, end, reporttotext, reporttohtml, tobbses, downbbses, messagetypes, modelmessage, instructions, retrieveat, report, flags) = (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) WHERE id=?", func(st *db.St) {
			st.BindText(session.CallSign)
			st.BindText(session.Name)
			st.BindText(session.Prefix)
			st.BindTime(session.Start, startEndFormat)
			st.BindTime(session.End, startEndFormat)
			st.BindText(strings.Join(session.ReportToText, ";"))
			st.BindText(strings.Join(session.ReportToHTML, ";"))
			st.BindText(strings.Join(session.ToBBSes, ";"))
			st.BindText(strings.Join(session.DownBBSes, ";"))
			st.BindText(strings.Join(session.MessageTypes, ";"))
			st.BindText(session.ModelMessage)
			st.BindText(session.Instructions)
			st.BindText(session.RetrieveAt)
			st.BindText(session.Report)
			st.BindInt(int(session.Flags))
			st.BindInt(session.ID)
			st.Step()
		})
		db.SQL(s.conn, "DELETE FROM retrieval WHERE session=?", func(st *db.St) {
			st.BindInt(session.ID)
			st.Step()
		})
		db.SQL(s.conn, "INSERT INTO retrieval (session, bbs, lastrun) VALUES (?,?,?)", func(st *db.St) {
			for _, r := range session.Retrieve {
				st.BindInt(session.ID)
				st.BindText(r.BBS)
				st.BindTime(r.LastRun, lastRunFormat)
				st.Step()
				st.Reset()
			}
		})
		return nil
	})
}

// DeleteSession deletes a session.
func (s *Store) DeleteSession(session *Session) {
	db.Transaction(s.conn, true, func() error {
		db.SQL(s.conn, "DELETE FROM session WHERE id=?", func(st *db.St) {
			st.BindInt(session.ID)
			st.Step()
		})
		return nil
	})
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
