package store

import (
	"time"
)

// A Response is an outgoing message that responds to a received message.
type Response struct {
	LocalID    string
	ResponseTo string
	To         string
	Subject    string
	Body       string
	SendTime   time.Time
	SenderCall string
	SenderBBS  string
}

// SaveResponse saves an outgoing response to the database.
func (st *Store) SaveResponse(r *Response) {
	var err error

	st.mutex.Lock()
	defer st.mutex.Unlock()
	_, err = st.dbh.Exec("INSERT INTO response (id, responseto, to, subject, body, sendtime, sendercall, senderbbs) VALUES (?,?,?,?,?,?,?,?)",
		r.LocalID, r.ResponseTo, r.To, r.Subject, r.Body, r.SendTime, r.SenderCall, r.SenderBBS)
	if err != nil {
		panic(err)
	}
}
