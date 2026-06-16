package webserver

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/rothskeller/packet/v4/message"
	"github.com/rothskeller/packet/v4/message/field"
	"github.com/rothskeller/packet/v4/wppsvr/english"
	"github.com/rothskeller/packet/v4/wppsvr/htmlb"
	"github.com/rothskeller/packet/v4/wppsvr/store"
)

var plainModelReplacer = strings.NewReplacer("¡", "", "-", "‑")

// serveInstructions displays the instructions for a session.
func (ws *webserver) serveInstructions(w http.ResponseWriter, r *http.Request) {
	var (
		session         *store.Session
		plainText       bool
		needHandling    bool
		needDestination bool
	)
	if sid, err := strconv.Atoi(r.FormValue("session")); err == nil {
		session = ws.st.GetSession(sid)
	}
	if session == nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	// Start the HTML page.
	w.Header().Set("Cache-Control", "nostore")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := htmlb.HTML(w)
	defer html.Close()
	html.E("meta charset=utf-8")
	html.E("title>Weekly Packet Practice - Santa Clara County ARES/RACES")
	html.E("meta name=viewport content='width=device-width, initial-scale=1'")
	html.E("link rel=stylesheet href=/static/common.css")
	html.E("link rel=stylesheet href=/static/instructions2.css")
	html.E("div id=org>Santa Clara County ARES<sup>®</sup>/RACES")
	html.E("div id=title").E("a href=/>Weekly Packet Practice")
	html.E("div id=subtitle>%s — %s", session.Name, session.End.Format("January 2, 2006"))
	// Write the main instructions.
	main := html.E("div id=main")
	para := main.E("p")
	para.R("For this practice session, please send ")
	if session.ModelMsg == nil {
		var article string
		var names []string
		for i, tag := range session.MessageTypes {
			if tag == "plain" {
				plainText = true
			} else {
				needHandling, needDestination = true, true
			}
			if mt := message.FindCreateTag(tag); mt != nil {
				words := strings.Fields(mt.Name())
				if i == 0 {
					article = words[0]
				}
				names = append(names, strings.Join(words[1:], " "))
			} else {
				if i == 0 {
					article = "a"
				}
				names = append(names, tag)
			}
		}
		para.TF("%s %s", article, english.Conjoin(names, "or"))
	} else if session.ModelMsg.Type() == message.PlainMessage {
		plainText = true
		var handling string
		for f := range session.ModelMsg.Fields() {
			if f.Common() == field.CHandling {
				handling = f.Value(session.ModelMsg)
				break
			}
		}
		switch handling {
		case "IMMEDIATE":
			para.R("an immediate")
		case "":
			para.R("a")
		default:
			para.TF("a %s", strings.ToLower(handling))
		}
		para.R(" plain text message with the text shown below")
	} else {
		words := strings.Fields(session.ModelMsg.Type().Name())
		para.TF("the %s shown below", strings.Join(words[1:], " "))
		for f := range session.ModelMsg.Fields() {
			switch f.Common() {
			case field.CHandling:
				if f.Value(session.ModelMsg) == "" {
					needHandling = true
				}
			case field.CToICSPosition, field.CToLocation:
				if f.Value(session.ModelMsg) == "" {
					needDestination = true
				}
			}
		}
	}
	para.TF(" to %s at the appropriate BBS.  The message must be received there between %s and %s.",
		session.CallSign, session.Start.Format("15:04 on Monday"), session.End.Format("15:04 on Monday"))
	switch len(session.DownBBSes) {
	case 0:
		para.R(" There are no simulated BBS outages for this session.")
	case 1:
		para.TF(" Do not use or send to %s during this session; it has a simulated outage.", session.DownBBSes[0])
	default:
		para.TF(" Do not use or send to %s during this session; they have simulated outages.", english.Conjoin(session.DownBBSes, "or"))
	}
	if session.Instructions != "" {
		main.R(session.Instructions)
	}
	if session.ModelMsg != nil && session.ModelMsg.Type() == message.PlainMessage {
		grid := main.E("div id=plainmodel")
		grid.E("div>Subject:")
		grid.E("div>%s", plainModelReplacer.Replace(session.ModelMsg.Subject().EncodedSubject()))
		grid.E("div>Message:")
		grid.E("div>%s", plainModelReplacer.Replace(session.ModelMsg.Body().EncodedBody()))
	}
	main.E("p style=margin-bottom:0>The following references may be helpful to you:")
	list := main.E("ul style=margin-top:0")
	list.E("li>The ").E("a href=https://www.scc-ares-races.org/operations/packet/bbs target=_blank>SCCo Packet BBSes").
		P().TF(" page will tell you which BBS to use to reach %s.", session.CallSign)
	list.E("li>The “Standard Outpost Configuration Instructions”, available on the ").
		E("a href=https://www.scc-ares-races.org/services/data/bbs target=_blank>Packet BBS Service").
		P().R(" page, will tell you how to configure Outpost to send messages following county standards.")
	if plainText {
		list.E("li>The “Standard Packet Message Subject Line”, available on the ").
			E("a href=https://www.scc-ares-races.org/services/data/bbs target=_blank>Packet BBS Service").
			P().R(" page, will tell you how to construct the subject line of your message to conform to county standards.")
	}
	if needHandling || needDestination {
		li := list.E("li>The “RACES Recommended Form Routing Cheat Sheet”, available on the ")
		li.E("a href=https://www.scc-ares-races.org/operations/forms/go-kit target=_blank>Go Kit Forms")
		li.R(" page, will tell you how to fill in the ")
		if needHandling && needDestination {
			li.R("handling and destination fields")
		} else if needHandling {
			li.R("handling field")
		} else {
			li.R("destination fields")
		}
		li.R(" for this message.")
	}
	list.E("li>The <kbd>packet@scc-ares-races.groups.io</kbd> mailing list is the best place to ask for help or report problems.  To join that list, see the instructions on the ").
		E("a href=https://www.scc-ares-races.org/about/email-lists target=_blank>Email Discussion Groups").
		P().R(" page.")
	count := ws.st.ModelImageCount(session.ID)
	for pnum := 1; pnum <= count; pnum++ {
		html.E("img class=modelimage src=/session/image?session=%d&page=%d", session.ID, pnum)
	}
}
