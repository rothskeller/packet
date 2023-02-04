package analyze

import (
	"fmt"
	"strings"
	"time"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/wppsvr/english"
	"github.com/rothskeller/packet/wppsvr/store"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/delivrcpt"
	"github.com/rothskeller/packet/xscmsg/readrcpt"
)

// The time.Now function can be overridden by tests.
var now = time.Now

type reference uint

const (
	refBBSList reference = 1 << iota
	refOutpostConfig
	refRouting
	refSubjectLine
	refWeeklyPractice
)

// Responses returns the list of messages that should be sent in response to the
// analyzed message.
func (a *Analysis) Responses(st astore) (list []*store.Response) {
	if a == nil || a.xsc == nil { // message already handled, no responses needed
		return nil
	}
	switch a.xsc.Type.Tag {
	case delivrcpt.Tag, readrcpt.Tag:
		break
	default:
		var dr = xscmsg.Create(delivrcpt.Tag)
		dr.Field("DeliveredSubject").Value = a.msg.Header.Get("Subject")
		dr.Field("DeliveredTime").Value = now().Format("01/02/2006 15:04:05")
		dr.Field("DeliveredTo").Value = fmt.Sprintf("%s@%s.ampr.org", strings.ToLower(a.session.CallSign), strings.ToLower(a.toBBS))
		dr.Field("LocalMessageID").Value = a.localID
		var r store.Response
		r.LocalID = st.NextMessageID(a.session.Prefix)
		r.ResponseTo = a.localID
		r.To = a.msg.ReturnAddress()
		r.Subject = dr.Subject()
		r.Body = dr.Body(false)
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
		counts   string
		invalids string
		rbody    strings.Builder
		wrapper  *english.Wrapper
	)
	if a.reportSubject == "" {
		return "", "" // no problem response message should be sent
	}
	if a.reportSubject == multipleProblemsSubject {
		counts = "s"
	}
	if a.invalid {
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
	wrapper.WriteString(a.reportText.String())
	// Add the references.
	wrapper.WriteString("For more information:")
	if a.references&refBBSList != 0 {
		wrapper.WriteString(bbsListRefText)
	}
	if a.references&refOutpostConfig != 0 {
		wrapper.WriteString(outpostConfigRefText)
	}
	if a.references&refRouting != 0 {
		wrapper.WriteString(routingRefText)
	}
	if a.references&refSubjectLine != 0 {
		wrapper.WriteString(subjectLineRefText)
	}
	if a.references&refWeeklyPractice != 0 {
		wrapper.WriteString(weeklyPracticeRefText)
	}
	wrapper.WriteString(packetGroupRefText)
	wrapper.Close()
	return a.reportSubject, rbody.String()
}

const bbsListRefText = `
  * The "County Packet Frequency List and BBS Info" page on the county
    ARES/RACES website gives a list of the known jurisdiction names and their
    abbreviations.  It is available at
    https://www.scc-ares-races.org/freqs/packet/freqs.html#assignments`
const outpostConfigRefText = `
  * The "Standard Outpost Configuration Instructions" document describes how
    to configure the Outpost messaging software to send messages following
    county standards.  It is available from the "Packet BBS Service" page at
    https://www.scc-ares-races.org/data/packet/index.html`
const routingRefText = `
  * The "SCCo ARES/RACES Recommended Form Routing" document gives
    recommendations for, among other things, what handling orders should be
    used for different types of forms, and what positions and locations they
    should be sent to.  It is available from the "Go Kit Forms" page at
    https://www.scc-ares-races.org/operations/go-kit-forms.html`
const subjectLineRefText = `
  * The "Standard Packet Message Subject Line" document describes how to
    compose the subject line of a packet message following county standards.
    It is available from the "Packet BBS Service" page at
    https://www.scc-ares-races.org/data/packet/index.html`
const weeklyPracticeRefText = `
  * The "Weekly SPECS/SVECS Packet Practice" page on the county ARES/RACES
    website gives details of the packet practice exercise, including the net
    practice schedules, the schedule of what type of message to send, the
    schedule of simulated outages of BBS systems, and the format of the subject
    line for practice messages.  It is available at
    https://www.scc-ares-races.org/data/packet/weekly-packet-practice.html`
const packetGroupRefText = `
  * If you need assistance, you can request it in the packet discussion group.
    To sign up for this group, see the Discussion Groups page at
    https://www.scc-ares-races.org/discuss-groups.html
` // This one is always last and has a newline at the end on purpose.
