package store

import (
	"time"

	"github.com/rothskeller/packet/wppsvr/db"
)

const deliveryTimeFormat = "2006-01-02 15:04:05-07:00"

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
func (st *Store) SessionHasMessages(sessionID int) (found bool) {
	db.SQL(st.conn, "SELECT 1 FROM message WHERE session=? LIMIT 1", func(st *db.St) {
		st.BindInt(sessionID)
		if st.Step() {
			found = true
		}
	})
	return found
}

// GetMessage returns the message with the specified local ID, or nil if there
// is none.
func (st *Store) GetMessage(localID string) (m *Message) {
	db.SQL(st.conn, "SELECT session, hash, deliverytime, message, fromaddress, fromcallsign, frombbs, tobbs, jurisdiction, messagetype, score, summary, analysis FROM message WHERE id=?", func(st *db.St) {
		st.BindText(localID)
		if st.Step() {
			m = new(Message)
			m.LocalID = localID
			m.Session = st.ColumnInt()
			m.Hash = st.ColumnText()
			m.DeliveryTime = st.ColumnTime(deliveryTimeFormat)
			m.Message = st.ColumnText()
			m.FromAddress = st.ColumnText()
			m.FromCallSign = st.ColumnText()
			m.FromBBS = st.ColumnText()
			m.ToBBS = st.ColumnText()
			m.Jurisdiction = st.ColumnText()
			m.MessageType = st.ColumnText()
			m.Score = st.ColumnInt()
			m.Summary = st.ColumnText()
			m.Analysis = st.ColumnText()
		}
	})
	return m
}

// GetMessageByHash returns the message with the specified hash, or nil if there
// is none.
func (st *Store) GetMessageByHash(hash string) (m *Message) {
	db.SQL(st.conn, "SELECT id, session, deliverytime, message, fromaddress, fromcallsign, frombbs, tobbs, jurisdiction, messagetype, score, summary, analysis FROM message WHERE hash=?", func(st *db.St) {
		st.BindText(hash)
		if st.Step() {
			m = new(Message)
			m.Hash = hash
			m.LocalID = st.ColumnText()
			m.Session = st.ColumnInt()
			m.DeliveryTime = st.ColumnTime(deliveryTimeFormat)
			m.Message = st.ColumnText()
			m.FromAddress = st.ColumnText()
			m.FromCallSign = st.ColumnText()
			m.FromBBS = st.ColumnText()
			m.ToBBS = st.ColumnText()
			m.Jurisdiction = st.ColumnText()
			m.MessageType = st.ColumnText()
			m.Score = st.ColumnInt()
			m.Summary = st.ColumnText()
			m.Analysis = st.ColumnText()
		}
	})
	return m
}

// GetSessionMessages returns the set of messages received for the session, in
// the order they were delivered to the BBS at which they were received.
func (st *Store) GetSessionMessages(sessionID int) (messages []*Message) {
	db.SQL(st.conn, "SELECT id, hash, deliverytime, message, fromaddress, fromcallsign, frombbs, tobbs, jurisdiction, messagetype, score, summary, analysis FROM message WHERE session=? ORDER BY deliverytime", func(st *db.St) {
		st.BindInt(sessionID)
		for st.Step() {
			var m Message

			m.Session = sessionID
			m.LocalID = st.ColumnText()
			m.Hash = st.ColumnText()
			m.DeliveryTime = st.ColumnTime(deliveryTimeFormat)
			m.Message = st.ColumnText()
			m.FromAddress = st.ColumnText()
			m.FromCallSign = st.ColumnText()
			m.FromBBS = st.ColumnText()
			m.ToBBS = st.ColumnText()
			m.Jurisdiction = st.ColumnText()
			m.MessageType = st.ColumnText()
			m.Score = st.ColumnInt()
			m.Summary = st.ColumnText()
			m.Analysis = st.ColumnText()
			messages = append(messages, &m)
		}
	})
	return messages
}

// HasMessageHash looks to see whether the database already contains a message
// with the specified hash.  If so, it returns the ID of that message; if not,
// it returns an empty string.
func (st *Store) HasMessageHash(hash string) (id string) {
	db.SQL(st.conn, "SELECT id FROM message WHERE hash=?", func(st *db.St) {
		st.BindText(hash)
		if st.Step() {
			id = st.ColumnText()
		}
	})
	return id
}

// SaveMessage saves a message to the database.
func (st *Store) SaveMessage(m *Message) {
	db.Transaction(st.conn, true, func() error {
		db.SQL(st.conn, "INSERT OR REPLACE INTO message (id, hash, deliverytime, message, session, fromaddress, fromcallsign, frombbs, tobbs, jurisdiction, messagetype, score, summary, analysis) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)", func(st *db.St) {
			st.BindText(m.LocalID)
			st.BindText(m.Hash)
			st.BindTime(m.DeliveryTime, deliveryTimeFormat)
			st.BindText(m.Message)
			st.BindInt(m.Session)
			st.BindText(m.FromAddress)
			st.BindText(m.FromCallSign)
			st.BindText(m.FromBBS)
			st.BindText(m.ToBBS)
			st.BindText(m.Jurisdiction)
			st.BindText(m.MessageType)
			st.BindInt(m.Score)
			st.BindText(m.Summary)
			st.BindText(m.Analysis)
			st.Step()
		})
		return nil
	})
}
