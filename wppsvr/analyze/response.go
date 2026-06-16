package analyze

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rothskeller/packet/v4/message/receipt"
	"github.com/rothskeller/packet/v4/wppsvr/store"
)

// The time.Now function can be overridden by tests.
var now = time.Now

// Responses returns the list of messages that should be sent in response to the
// analyzed message.  In current implementation, it only returns delivery
// receipts.
func (a *Analysis) Responses(st astore) (list []*store.Response) {
	if a == nil || a.msg == nil { // message already handled, no responses needed
		return nil
	}
	switch a.msg.Type() {
	case receipt.DeliveryReceipt, receipt.ReadReceipt:
		break
	default:
		var extra string
		switch a.sm.Score {
		case 0:
			extra = fmt.Sprintf("MESSAGE WAS NOT COUNTED as a check-in to the %s on %s.\nReason: %s\nFor more information, visit https://wpp.scc-ares-races.org",
				a.session.Name, a.session.End.Format("January 2"), a.sm.Summary)
		case 100:
			extra = fmt.Sprintf("100%% correct check-in to the %s on %s.",
				a.session.Name, a.session.End.Format("January 2"))
		default:
			extra = fmt.Sprintf("%d%% score for check-in to the %s on %s.\nReason: %s\nFor more information, visit https://wpp.scc-ares-races.org",
				a.sm.Score, a.session.Name, a.session.End.Format("January 2"), a.sm.Summary)
		}
		dr, err := receipt.NewDeliveryReceipt(
			a.msg.ReturnAddr(),
			fmt.Sprintf("%s@%s.ampr.org", strings.ToLower(a.session.CallSign), strings.ToLower(a.sm.ToBBS)),
			a.msg.Subject().EncodedSubject(),
			a.sm.LocalID,
			now(),
			extra)
		if err != nil {
			log.Printf("ERROR: can't generate delivery receipt: %s", err)
			break
		}
		var r store.Response
		r.LocalID = st.NextMessageID(a.session.Prefix)
		r.ResponseTo = a.sm.LocalID
		r.To = a.msg.ReturnAddr()
		r.Subject = dr.Subject().EncodedSubject()
		r.Body = dr.Payload().Encode()
		r.SenderCall = a.session.CallSign
		r.SenderBBS = a.sm.ToBBS
		list = append(list, &r)
	}
	return list
}
