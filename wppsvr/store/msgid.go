package store

import (
	"fmt"

	"github.com/rothskeller/packet/wppsvr/db"
)

// NextMessageID returns the next message ID in the sequence with the specified
// prefix.
func (st *Store) NextMessageID(prefix string) string {
	var num int

	db.Transaction(st.conn, true, func() error {
		db.SQL(st.conn, "SELECT num FROM msgnum WHERE prefix=?", func(st *db.St) {
			st.BindText(prefix)
			if st.Step() {
				num = st.ColumnInt() + 1
			} else {
				num = 1
			}
		})
		db.SQL(st.conn, "INSERT OR REPLACE INTO msgnum (prefix, num) VALUES (?,?)", func(st *db.St) {
			st.BindText(prefix)
			st.BindInt(num)
			st.Step()
		})
		return nil
	})
	return fmt.Sprintf("%s-%03dP", prefix, num)
}
