package store

import (
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
