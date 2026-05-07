package store

import (
	"time"

	"github.com/rothskeller/packet/wppsvr/db"
)

const sendTimeFormat = "2006-01-02 15:04:05.999999999-07:00"

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
	db.SQL(st.conn, "SELECT id, sendto, subject, body, sendtime, sendercall, senderbbs FROM response WHERE responseto=? ORDER BY id", func(st *db.St) {
		st.BindText(to)
		for st.Step() {
			var r Response

			r.LocalID = st.ColumnText()
			r.ResponseTo = to
			r.To = st.ColumnText()
			r.Subject = st.ColumnText()
			r.Body = st.ColumnText()
			r.SendTime = st.ColumnTime(sendTimeFormat)
			r.SenderCall = st.ColumnText()
			r.SenderBBS = st.ColumnText()
			responses = append(responses, &r)
		}
	})
	return responses
}

// SaveResponse saves an outgoing response to the database.
func (st *Store) SaveResponse(r *Response) {
	db.Transaction(st.conn, true, func() error {
		db.SQL(st.conn, "INSERT INTO response (id, responseto, sendto, subject, body, sendtime, sendercall, senderbbs) VALUES (?,?,?,?,?,?,?,?)", func(st *db.St) {
			st.BindText(r.LocalID)
			st.BindText(r.ResponseTo)
			st.BindText(r.To)
			st.BindText(r.Subject)
			st.BindText(r.Body)
			st.BindTime(r.SendTime, sendTimeFormat)
			st.BindText(r.SenderCall)
			st.BindText(r.SenderBBS)
			st.Step()
		})
		return nil
	})
}
