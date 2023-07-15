package report

import (
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/jnos"
	"github.com/rothskeller/packet/message/common"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/store"
)

// Send generates the report for the session and sends it to all designated
// recipients, through the supplied open BBS connection and/or via SMTP.  It
// also stores the report in the session.
func Send(st Store, conn *jnos.Conn, session *store.Session) {
	report := Generate(st, session)
	sendTo := session.ReportToText
	if session.Flags&store.ReportToSenders != 0 {
		sendTo = append(sendTo, report.Participants...)
	}
	session.Report = report.RenderPlainText()
	st.UpdateSession(session)
	subject := common.EncodeSubject(st.NextMessageID(session.Prefix), "ROUTINE", "", "SCCo Packet Practice Report")
	body := new(envelope.Envelope).RenderBody(session.Report)
	// To avoid potential problems with JNOS line length limits, we
	// send to each recipient separately.
	// conn.Send(subject, body, sendTo...)
	for _, addr := range sendTo {
		conn.Send(subject, body, addr)
	}
	if len(session.ReportToHTML) != 0 {
		if err := report.SendHTML(session.ReportToHTML); err != nil {
			log.Printf("ERROR: %s", err)
		}
	}
	log.Printf("Sent report for %s on %s.", session.Name, session.End.Format("2006-01-02"))
}

// SendHTML sends the report in HTML format to the specified address(es).
func (r *Report) SendHTML(to []string) error {
	conf := config.Get().SMTP
	var addrs []string
	for i, t := range to {
		if addr, err := mail.ParseAddress(t); err == nil {
			addrs = append(addrs, addr.Address)
			to[i] = addr.String()
		} else {
			return fmt.Errorf("address %q: %s", t, err)
		}
	}
	msg := r.RenderEmail(strings.Join(to, ", "))
	from, _ := mail.ParseAddress(conf.From)
	host, _, _ := net.SplitHostPort(conf.Server)
	auth := smtp.PlainAuth("", conf.Username, conf.Password, host)
	if err := smtp.SendMail(conf.Server, auth, from.Address, addrs, []byte(msg)); err != nil {
		return fmt.Errorf("smtp.SendMail: %s", err)
	}
	return nil
}
