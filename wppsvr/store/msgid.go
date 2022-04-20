package store

import (
	"database/sql"
	"fmt"
)

// NextMessageID returns the next message ID in the sequence with the specified
// prefix.
func (st *Store) NextMessageID(prefix string) string {
	var (
		num int
		err error
	)
	st.mutex.Lock()
	defer st.mutex.Unlock()
	err = st.dbh.QueryRow("SELECT num FROM msgnum WHERE prefix=?", prefix).Scan(&num)
	switch err {
	case nil:
		num++
	case sql.ErrNoRows:
		num = 1
	default:
		panic(err)
	}
	_, err = st.dbh.Exec("UPDATE msgnum SET num=? WHERE prefix=?", num, prefix)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s-%03d", prefix, num)
}
