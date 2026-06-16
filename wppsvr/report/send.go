package report

import (
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/rothskeller/packet/v4/jnos"
	"github.com/rothskeller/packet/v4/message/body"
	"github.com/rothskeller/packet/v4/message/payload"
	"github.com/rothskeller/packet/v4/message/subject"
	"github.com/rothskeller/packet/v4/wppsvr/config"
	"github.com/rothskeller/packet/v4/wppsvr/store"
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
	subj, _ := subject.NewPlainSubject(st.NextMessageID(session.Prefix), "ROUTINE", "SCCo Packet Practice Report")
	pb := body.NewPlainBody(session.Report)
	pay := payload.NewOutpostPayload(pb, true, true)
	// To avoid potential problems with JNOS line length limits, we
	// send to each recipient separately.
	for _, addr := range sendTo {
		conn.Send(subj.EncodedSubject(), pay.Encode(), addr)
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
