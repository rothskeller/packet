package store

import (
	"database/sql"
	"time"
)

// A Response is an outgoing message that responds to a received message.
type Response struct {
	LocalID    string    `yaml:"localID"`
	ResponseTo string    `yaml:"responseTo"`
	To         string    `yaml:"to"`
	Subject    string    `yaml:"subject"`
	Body       string    `yaml:"body"`
	SendTime   time.Time `yaml:"sendTime"`
	SenderCall string    `yaml:"senderCall"`
	SenderBBS  string    `yaml:"senderBBS"`
}

// GetResponses retrieves the responses for the specified message.
func (st *Store) GetResponses(to string) (responses []*Response) {
	var (
		rows *sql.Rows
		err  error
	)
	st.mutex.RLock()
	defer st.mutex.RUnlock()
	rows, err = st.dbh.Query(`SELECT id, sendto, subject, body, sendtime, sendercall, senderbbs FROM response WHERE responseto=? ORDER BY id`, to)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var r Response
		err = rows.Scan(&r.LocalID, &r.To, &r.Subject, &r.Body, &r.SendTime, &r.SenderCall, &r.SenderBBS)
		if err != nil {
			panic(err)
		}
		r.ResponseTo = to
		responses = append(responses, &r)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}
	return responses
}

// SaveResponse saves an outgoing response to the database.
func (st *Store) SaveResponse(r *Response) {
	var err error

	st.mutex.Lock()
	defer st.mutex.Unlock()
	_, err = st.dbh.Exec("INSERT INTO response (id, responseto, sendto, subject, body, sendtime, sendercall, senderbbs) VALUES (?,?,?,?,?,?,?,?)",
		r.LocalID, r.ResponseTo, r.To, r.Subject, r.Body, r.SendTime, r.SenderCall, r.SenderBBS)
	if err != nil {
		panic(err)
	}
}
