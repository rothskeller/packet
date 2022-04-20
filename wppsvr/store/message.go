package store

import (
	"database/sql"
	"strings"
	"time"
)

// A Message describes a single received message.
type Message struct {
	LocalID      string
	Hash         string
	DeliveryTime time.Time
	Message      string
	Session      int
	FromAddress  string
	FromCallSign string
	FromBBS      string
	ToBBS        string
	Subject      string
	Problems     []string
	Valid        bool
	Correct      bool
}

// SessionHasMessages returns whether there are any messages stored for the
// specified session.
func (st *Store) SessionHasMessages(sessionID int) bool {
	var (
		dummy int
		err   error
	)
	st.mutex.RLock()
	defer st.mutex.RUnlock()
	err = st.dbh.QueryRow("SELECT 1 FROM message WHERE session=? LIMIT 1", sessionID).Scan(&dummy)
	switch err {
	case nil:
		return true
	case sql.ErrNoRows:
		return false
	default:
		panic(err)
	}
}

// GetSessionMessages returns the set of messages received for the session, in
// the order they were delivered to the BBS at which they were received.
func (st *Store) GetSessionMessages(sessionID int) (messages []*Message) {
	var (
		rows *sql.Rows
		err  error
	)
	st.mutex.RLock()
	defer st.mutex.RUnlock()
	rows, err = st.dbh.Query("SELECT id, hash, deliverytime, message, fromaddress, fromcallsign, frombbs, tobbs, subject, problems, valid, correct FROM message WHERE session=? ORDER BY deliverytime", sessionID)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var (
			m        Message
			problems string
		)
		err = rows.Scan(&m.LocalID, &m.Hash, &m.Message, &m.FromAddress, &m.FromCallSign, &m.FromBBS, &m.ToBBS, &m.Subject, &problems, &m.Valid, &m.Correct)
		if err != nil {
			panic(err)
		}
		m.Problems = split(problems)
		m.Session = sessionID
		messages = append(messages, &m)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}
	return messages
}

// HasMessageHash returns whether the database already contains a message with
// the specified hash.
func (st *Store) HasMessageHash(hash string) bool {
	var (
		dummy int
		err   error
	)
	st.mutex.RLock()
	defer st.mutex.RUnlock()
	err = st.dbh.QueryRow("SELECT 1 FROM message WHERE hash=?", hash).Scan(&dummy)
	switch err {
	case nil:
		return true
	case sql.ErrNoRows:
		return false
	default:
		panic(err)
	}
}

// SaveMessage saves a message to the database.
func (st *Store) SaveMessage(m *Message) {
	var err error

	st.mutex.Lock()
	defer st.mutex.Unlock()
	_, err = st.dbh.Exec("INSERT INTO message (id, hash, deliverytime, message, session, fromaddress, fromcallsign, frombbs, tobbs, subject, problems, valid, correct) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)",
		m.LocalID, m.Hash, m.DeliveryTime, m.Message, m.Session, m.FromAddress, m.FromBBS, m.ToBBS, m.Subject, strings.Join(m.Problems, ","), m.Valid, m.Correct)
	if err != nil {
		panic(err)
	}
}
