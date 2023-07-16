package store

import (
	"database/sql"
	"time"
)

// A Message describes a single received message.
type Message struct {
	LocalID      string    `yaml:"localID"`
	Hash         string    `yaml:"hash"`
	DeliveryTime time.Time `yaml:"deliveryTime"`
	Message      string    `yaml:"message"`
	Session      int       `yaml:"session"`
	FromAddress  string    `yaml:"fromAddress"`
	FromCallSign string    `yaml:"fromCallSign"`
	FromBBS      string    `yaml:"fromBBS"`
	ToBBS        string    `yaml:"toBBS"`
	Jurisdiction string    `yaml:"jurisdiction"`
	MessageType  string    `yaml:"messageType"`
	Score        int       `yaml:"score"`
	Summary      string    `yaml:"summary"`
	Analysis     string    `yaml:"analysis"`
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

// GetMessage returns the message with the specified local ID, or nil if there
// is none.
func (st *Store) GetMessage(localID string) *Message {
	var (
		m   Message
		err error
	)
	m.LocalID = localID
	st.mutex.RLock()
	defer st.mutex.RUnlock()
	err = st.dbh.QueryRow("SELECT session, hash, deliverytime, message, fromaddress, fromcallsign, frombbs, tobbs, jurisdiction, messagetype, score, summary, analysis FROM message WHERE id=?", localID).
		Scan(&m.Session, &m.Hash, &m.DeliveryTime, &m.Message, &m.FromAddress, &m.FromCallSign, &m.FromBBS, &m.ToBBS, &m.Jurisdiction, &m.MessageType, &m.Score, &m.Summary, &m.Analysis)
	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
	return &m
}

// GetMessageByHash returns the message with the specified hash, or nil if there
// is none.
func (st *Store) GetMessageByHash(hash string) *Message {
	var (
		m   Message
		err error
	)
	m.Hash = hash
	st.mutex.RLock()
	defer st.mutex.RUnlock()
	err = st.dbh.QueryRow("SELECT id, session, deliverytime, message, fromaddress, fromcallsign, frombbs, tobbs, jurisdiction, messagetype, score, summary, analysis FROM message WHERE hash=?", hash).
		Scan(&m.LocalID, &m.Session, &m.DeliveryTime, &m.Message, &m.FromAddress, &m.FromCallSign, &m.FromBBS, &m.ToBBS, &m.Jurisdiction, &m.MessageType, &m.Score, &m.Summary, &m.Analysis)
	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
	return &m
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
	rows, err = st.dbh.Query("SELECT id, hash, deliverytime, message, fromaddress, fromcallsign, frombbs, tobbs, jurisdiction, messagetype, score, summary, analysis FROM message WHERE session=? ORDER BY deliverytime", sessionID)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var m Message

		err = rows.Scan(&m.LocalID, &m.Hash, &m.DeliveryTime, &m.Message, &m.FromAddress, &m.FromCallSign, &m.FromBBS, &m.ToBBS, &m.Jurisdiction, &m.MessageType, &m.Score, &m.Summary, &m.Analysis)
		if err != nil {
			panic(err)
		}
		m.Session = sessionID
		messages = append(messages, &m)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}
	return messages
}

// HasMessageHash looks to see whether the database already contains a message
// with the specified hash.  If so, it returns the ID of that message; if not,
// it returns an empty string.
func (st *Store) HasMessageHash(hash string) (id string) {
	var (
		err error
	)
	st.mutex.RLock()
	defer st.mutex.RUnlock()
	err = st.dbh.QueryRow("SELECT id FROM message WHERE hash=?", hash).Scan(&id)
	switch err {
	case nil:
		return id
	case sql.ErrNoRows:
		return ""
	default:
		panic(err)
	}
}

// SaveMessage saves a message to the database.
func (st *Store) SaveMessage(m *Message) {
	var err error

	st.mutex.Lock()
	defer st.mutex.Unlock()
	_, err = st.dbh.Exec("INSERT OR REPLACE INTO message (id, hash, deliverytime, message, session, fromaddress, fromcallsign, frombbs, tobbs, jurisdiction, messagetype, score, summary, analysis) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		m.LocalID, m.Hash, m.DeliveryTime, m.Message, m.Session, m.FromAddress, m.FromCallSign, m.FromBBS, m.ToBBS, m.Jurisdiction, m.MessageType, m.Score, m.Summary, m.Analysis)
	if err != nil {
		panic(err)
	}
}
