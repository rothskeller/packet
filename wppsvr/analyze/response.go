package analyze

import (
	"fmt"
	"regexp"
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
		count      int
		counts     string
		invalid    bool
		invalids   string
		rbody      strings.Builder
		wrapper    *english.Wrapper
		problems   = config.Get().Problems
		reftext    = config.Get().References
		references = map[string]bool{}
	)
	// Count the number of problems to include in the message, and note
	// whether any of them prevent the message from counting.  We need that
	// information for the header of the message.
	for p := range a.problems {
		pdef := problems[p.Code]
		if pdef.ActionFlags&config.ActionRespond != 0 {
			count++
			subject = pdef.Label
		}
		if pdef.ActionFlags&config.ActionDontCount != 0 {
			invalid = true
		}
		for _, ref := range strings.Fields(pdef.References) {
			references[ref] = true
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
	for p, rk := range a.problems {
		if problems[p.Code].ActionFlags&config.ActionRespond != 0 {
			writeProblemResponse(a, p, problems[p.Code].Responses[rk], wrapper)
		}
	}
	// Add the references.
	wrapper.WriteString("\nFor more information:\n")
	for ref := range references {
		if ref != "packetGroup" { // save it for last
			writeReference(wrapper, reftext[ref])
		}
	}
	writeReference(wrapper, reftext["packetGroup"])
	wrapper.Close()
	return subject, rbody.String()
}
func writeReference(wrapper *english.Wrapper, reftext string) {
	for _, line := range strings.Split(strings.TrimSpace(reftext), "\n") {
		wrapper.WriteString("  " + line + "\n")
	}
}

var variableRE = regexp.MustCompile(`\{([^}]*)\}`)

func writeProblemResponse(a *Analysis, p *Problem, response string, wrapper *english.Wrapper) {
	wrapper.WriteString("\n")
	for {
		if match := variableRE.FindStringIndex(response); match != nil {
			wrapper.WriteString(response[:match[0]])
			if fn, ok := p.Variables[response[match[0]+1:match[1]-1]]; ok && fn != nil {
				wrapper.WriteString(fn(a))
			} else {
				wrapper.WriteString(Variables[response[match[0]+1:match[1]-1]](a))
			}
			response = response[match[1]:]
		} else {
			wrapper.WriteString(response)
			return
		}
	}
}
