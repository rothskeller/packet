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
	report, sendTo := generate(st, session)
	sendTo = append(sendTo, session.ReportTo...)
	session.Report = report
	st.UpdateSession(session)
	var rm = pktmsg.New()
	rm.Body = report
	subject := xscmsg.EncodeSubject(st.NextMessageID(session.Prefix), xscmsg.HandlingRoutine, "", "SCCo Packet Practice Report")
	conn.Send(subject, rm.EncodeBody(false), sendTo...)
	log.Printf("Sent report for %s on %s.", session.Name, session.End.Format("2006-01-02"))
}
