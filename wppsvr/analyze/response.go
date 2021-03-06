package analyze

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/english"
	"github.com/rothskeller/packet/wppsvr/store"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/delivrcpt"
	"github.com/rothskeller/packet/xscmsg/readrcpt"
)

// The time.Now function can be overridden by tests.
var now = time.Now

// Responses returns the list of messages that should be sent in response to the
// analyzed message.
func (a *Analysis) Responses(st astore) (list []*store.Response) {
	if a == nil { // message already handled, no responses needed
		return nil
	}
	switch a.xsc.(type) {
	case nil, *delivrcpt.DeliveryReceipt, *readrcpt.ReadReceipt:
		break
	default:
		var dr = xscmsg.Create("DELIVERED").(*delivrcpt.DeliveryReceipt)
		dr.DeliveredSubject = a.msg.Header.Get("Subject")
		dr.DeliveredTime = now().Format("01/02/2006 15:04:05")
		dr.DeliveredTo = fmt.Sprintf("%s@%s.ampr.org", strings.ToLower(a.session.CallSign), strings.ToLower(a.toBBS))
		dr.LocalMessageID = a.localID
		var r store.Response
		r.LocalID = st.NextMessageID(a.session.Prefix)
		r.ResponseTo = a.localID
		r.To = a.msg.ReturnAddress()
		var drmsg = dr.Message(false)
		r.Subject = drmsg.Header.Get("Subject")
		r.Body = drmsg.EncodeBody(false)
		r.SenderCall = a.session.CallSign
		r.SenderBBS = a.toBBS
		list = append(list, &r)
	}
	if rsubject, rbody := a.responseMessage(); rsubject != "" {
		var rm = pktmsg.New()
		rm.Body = rbody
		var r store.Response
		r.LocalID = st.NextMessageID(a.session.Prefix)
		r.ResponseTo = a.localID
		r.To = a.msg.ReturnAddress()
		r.Subject = xscmsg.EncodeSubject(r.LocalID, xscmsg.HandlingRoutine, "", rsubject)
		r.Body = rm.EncodeBody(false)
		r.SenderCall = a.session.CallSign
		r.SenderBBS = a.toBBS
		list = append(list, &r)
	}
	return list
}

// responseMessage generates the problem response message from the problems
// logged with the message.
func (a *Analysis) responseMessage() (subject, body string) {
	var (
		count      int
		counts     string
		invalid    bool
		invalids   string
		rbody      strings.Builder
		wrapper    *english.Wrapper
		actions    = config.Get().ProblemActionFlags
		references = refPacketGroup
	)
	// Count the number of problems to include in the message, and note
	// whether any of them prevent the message from counting.  We need that
	// information for the header of the message.
	for _, p := range a.problems {
		action, ok := actions[p.code]
		if !ok {
			log.Printf("ERROR: config doesn't specify how to handle %s; ignoring", p.code)
		}
		if action&config.ActionRespond != 0 {
			count++
			subject = ProblemLabel[p.code]
			references |= p.references
		}
		if action&config.ActionDontCount != 0 {
			invalid = true
		}
	}
	switch count {
	case 0:
		return "", "" // no problem response message should be sent
	case 1:
		subject = strings.ToUpper(subject[:1]) + subject[1:]
	default:
		counts = "s"
		subject = "Issues with packet practice message"
	}
	if invalid {
		invalids = "  The message will not be counted."
	}
	// Generate the header of the message.
	wrapper = english.NewWrapper(&rbody)
	fmt.Fprintf(wrapper, `The packet practice message
    From: %s
    To: %s@%s
    Subject: %s
    Date: %s
has the following issue%s.%s
`,
		a.msg.ReturnAddress(),
		strings.ToLower(a.session.CallSign), strings.ToLower(a.toBBS),
		a.msg.Header.Get("Subject"),
		a.msg.Header.Get("Date"),
		counts, invalids)
	// Add the paragraphs describing each problem.
	for _, p := range a.problems {
		if actions[p.code]&config.ActionRespond != 0 {
			wrapper.WriteString(p.response)
		}
	}
	// Add the references.
	wrapper.WriteString("\nFor more information:")
	for _, r := range allReferences {
		if references&r != 0 {
			wrapper.WriteString(referenceText[r])
		}
	}
	wrapper.Close()
	return subject, rbody.String()
}
