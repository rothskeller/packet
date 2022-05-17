package report

import (
	"log"

	"steve.rothskeller.net/packet/jnos"
	"steve.rothskeller.net/packet/pktmsg"
	"steve.rothskeller.net/packet/wppsvr/store"
	"steve.rothskeller.net/packet/xscmsg"
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
	conn.Send(subject, rm.EncodeBody(false), sendTo...)
	log.Printf("Sent report for %s on %s.", session.Name, session.End.Format("2006-01-02"))
}
