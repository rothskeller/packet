package report

import (
	"log"

	"github.com/rothskeller/packet/jnos"
	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/wppsvr/store"
	"github.com/rothskeller/packet/xscmsg"
)

// Send generates the report for the session and sends it to all designated
// recipients, through the supplied open BBS connection.
func Send(st Store, conn *jnos.Conn, session *store.Session) {
	report := Generate(st, session)
	sendTo := session.ReportTo
	if len(sendTo) != 0 && sendTo[0] == "MESSAGE-SENDERS" {
		sendTo = append(report.Participants, sendTo[1:]...)
	}
	session.Report = report.RenderPlainText()
	st.UpdateSession(session)
	var rm = pktmsg.New()
	rm.Body = session.Report
	subject := xscmsg.EncodeSubject(st.NextMessageID(session.Prefix), xscmsg.HandlingRoutine, "", "SCCo Packet Practice Report")
	body := rm.EncodeBody(false)
	// To avoid potential problems with JNOS line length limits, we
	// send to each recipient separately.
	// conn.Send(subject, body, sendTo...)
	for _, addr := range sendTo {
		conn.Send(subject, body, addr)
	}
	log.Printf("Sent report for %s on %s.", session.Name, session.End.Format("2006-01-02"))
}
