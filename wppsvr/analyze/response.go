package analyze

import (
	"fmt"
	"strings"
	"time"

	"github.com/rothskeller/packet/message/delivrcpt"
	"github.com/rothskeller/packet/message/readrcpt"
	"github.com/rothskeller/packet/wppsvr/store"
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
	switch a.msg.(type) {
	case *delivrcpt.DeliveryReceipt, *readrcpt.ReadReceipt:
		break
	default:
		var dr delivrcpt.DeliveryReceipt
		dr.MessageSubject = a.subject
		dr.MessageTo = fmt.Sprintf("%s@%s.ampr.org", strings.ToLower(a.session.CallSign), strings.ToLower(a.sm.ToBBS))
		dr.DeliveredTime = now().Format("01/02/2006 15:04:05")
		dr.LocalMessageID = a.sm.LocalID
		switch a.sm.Score {
		case 0:
			dr.ExtraText = fmt.Sprintf("MESSAGE WAS NOT COUNTED as a check-in to the %s on %s.\nReason: %s\nFor more information, visit https://scc-ares-races.org/pacpractice",
				a.session.Name, a.session.End.Format("January 2"), a.sm.Summary)
		case 100:
			dr.ExtraText = fmt.Sprintf("100%% correct check-in to the %s on %s.",
				a.session.Name, a.session.End.Format("January 2"))
		default:
			dr.ExtraText = fmt.Sprintf("%d%% score for check-in to the %s on %s.\nReason: %s\nFor more information, visit https://scc-ares-races.org/pacpractice",
				a.sm.Score, a.session.Name, a.session.End.Format("January 2"), a.sm.Summary)
		}
		var r store.Response
		r.LocalID = st.NextMessageID(a.session.Prefix)
		r.ResponseTo = a.sm.LocalID
		r.To = a.env.ReturnAddr
		r.Subject = dr.EncodeSubject()
		r.Body = dr.EncodeBody()
		r.SenderCall = a.session.CallSign
		r.SenderBBS = a.sm.ToBBS
		list = append(list, &r)
	}
	return list
}
