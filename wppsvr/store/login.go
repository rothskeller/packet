package store

import (
	"time"

	"github.com/rothskeller/packet/v4/wppsvr/db"
)

const expiresFormat = "2006-01-02 15:04:05-07:00"

// GetLogin looks up an authorization token to determine whether it is valid.
// If so, it returns the corresponding call sign.  If not, it returns an empty
// string.
func (s *Store) GetLogin(token string) (callsign string) {
	db.Transaction(s.conn, true, func() error {
		db.SQL(s.conn, "DELETE FROM login WHERE expires<?", func(st *db.St) {
			st.BindTime(time.Now(), expiresFormat)
			st.Step()
		})
		return nil
	})
	db.SQL(s.conn, "SELECT callsign FROM login WHERE token=?", func(st *db.St) {
		st.BindText(token)
		if st.Step() {
			callsign = st.ColumnText()
		}
	})
	return callsign
}

// AddLogin adds a login to the database.
func (s *Store) AddLogin(token, callsign string, expires time.Time) {
	db.Transaction(s.conn, true, func() error {
		db.SQL(s.conn, "INSERT INTO login (token, callsign, expires) VALUES (?,?,?)", func(st *db.St) {
			st.BindText(token)
			st.BindText(callsign)
			st.BindTime(expires, expiresFormat)
			st.Step()
		})
		return nil
	})
}
