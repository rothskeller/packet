package store

import (
	"database/sql"
	"time"
)

// GetLastRetrieval returns the time of the last successful retrieval for the
// specified toCallSign from the specified BBS.
func (st *Store) GetLastRetrieval(toCallSign, bbs string) (retrieval time.Time) {
	var err error

	st.mutex.RLock()
	defer st.mutex.RUnlock()
	err = st.dbh.QueryRow("SELECT time FROM retrieval WHERE callsign=? AND bbs=?", toCallSign, bbs).Scan(&retrieval)
	switch err {
	case nil:
		return retrieval
	case sql.ErrNoRows:
		return time.Time{}
	default:
		panic(err)
	}
}

// SetLastRetrieval updates the time of the last successful retrieval for the
// specified toCallSign and the specified BBS.
func (st *Store) SetLastRetrieval(toCallSign, bbs string, retrieval time.Time) {
	var err error

	st.mutex.Lock()
	defer st.mutex.Unlock()
	_, err = st.dbh.Exec("INSERT OR REPLACE INTO retrieval (callsign, bbs, time) VALUES (?,?,?)",
		toCallSign, bbs, retrieval)
	if err != nil {
		panic(err)
	}
}
