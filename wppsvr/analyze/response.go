package analyze

import (
	"fmt"
	"strings"
	"time"

	"steve.rothskeller.net/packet/pktmsg"
	"steve.rothskeller.net/packet/wppsvr/english"
	"steve.rothskeller.net/packet/wppsvr/store"
)

// Responses returns the list of messages that should be sent in response to the
// analyzed message.
func (a *Analysis) Responses(st *store.Store) (list []*store.Response) {
	if a == nil { // message already handled, no responses needed
		return nil
	}
	if msg := a.msg.Message(); msg != nil {
		var dr pktmsg.TxDeliveryReceipt
		dr.DeliveredSubject = msg.SubjectLine
		dr.DeliveredTime = time.Now()
		dr.DeliveredTo = fmt.Sprintf("%s@%s.ampr.org", strings.ToLower(a.session.CallSign), strings.ToLower(a.toBBS))
		dr.LocalMessageID = a.localID
		var r store.Response
		r.LocalID = st.NextMessageID(a.session.Prefix)
		r.ResponseTo = a.localID
		r.To = a.msg.Base().ReturnAddress
		r.Subject, r.Body, _ = dr.Encode()
		r.SenderCall = a.session.CallSign
		r.SenderBBS = a.toBBS
		list = append(list, &r)
	}
	if rsubject, rbody := a.responseMessage(); rsubject != "" {
		var rm pktmsg.TxMessage
		rm.Body = rbody
		rm.HandlingOrder = pktmsg.HandlingRoutine
		rm.MessageNumber = st.NextMessageID(a.session.Prefix)
		rm.Subject = rsubject
		var r store.Response
		r.LocalID = rm.MessageNumber
		r.ResponseTo = a.localID
		r.To = a.msg.Base().ReturnAddress
		r.Subject, r.Body, _ = rm.Encode()
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
		references reference = refPacketGroup
	)
	wrapper = english.NewWrapper(&rbody)
	// Count the number of problems to include in the message, and note
	// whether any of them prevent the message from counting.  We need that
	// information for the header of the message.
	for _, p := range a.problems {
		if p.subject == "" {
			continue
		}
		count++
		subject = p.subject
		if p.invalid {
			invalid = true
		}
		references |= p.references
	}
	switch count {
	case 0:
		return "", "" // no problem response message should be sent
	case 1:
		break
	default:
		counts = "s"
		subject = "Issues with packet practice message"
	}
	if invalid {
		invalids = "  The message will not be counted."
	}
	// Generate the header of the message.
	fmt.Fprintf(wrapper, `The packet practice message
    From: %s
    To: %s@%s
    Subject: %s
    Date: %s
has the following issue%s.%s
`,
		a.msg.Base().ReturnAddress,
		strings.ToLower(a.session.CallSign), strings.ToLower(a.toBBS),
		a.msg.Base().SubjectLine,
		a.msg.Base().DateLine,
		counts, invalids)
	// Add the paragraphs describing each problem.
	for _, p := range a.problems {
		if p.subject == "" {
			continue
		}
		wrapper.WriteString(p.response)
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
