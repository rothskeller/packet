package webserver

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/plaintext"
	"github.com/rothskeller/packet/wppsvr/english"
	"github.com/rothskeller/packet/wppsvr/htmlb"
	"github.com/rothskeller/packet/wppsvr/store"
)

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
	html.E("link rel=stylesheet href=/static/instructions.css")
	html.E("div id=org>Santa Clara County ARES<sup>®</sup>/RACES")
	html.E("div id=title>Weekly Packet Practice")
	html.E("div id=subtitle>%s — %s", session.Name, session.End.Format("January 2, 2006"))
	// Write the main instructions.
	main := html.E("div id=main")
	para := main.E("p")
	para.R("For this practice session, please send ")
	switch msg := session.ModelMsg.(type) {
	case nil:
		var article string
		var names []string
		for i, tag := range session.MessageTypes {
			if tag == "plain" {
				plainText = true
			} else {
				needHandling, needDestination = true, true
			}
			if mt := message.RegisteredTypes[tag]; mt != nil {
				if i == 0 {
					article = mt.Article
				}
				names = append(names, mt.Name)
			} else {
				if i == 0 {
					article = "a"
				}
				names = append(names, tag)
			}
		}
		para.TF("%s %s", article, english.Conjoin(names, "or"))
	case *plaintext.PlainText:
		plainText = true
		switch msg.Handling {
		case "IMMEDIATE":
			para.R("an immediate")
		case "":
			para.R("a")
		default:
			para.TF("a %s", strings.ToLower(msg.Handling))
		}
		para.R(" plain text message with the text shown below")
	default:
		para.TF("the %s shown below", msg.Type().Name)
		if f, ok := msg.(message.IKeyFields); ok {
			kf := f.KeyFields()
			if kf.Handling == "" {
				needHandling = true
			}
			if kf.ToICSPosition == "" || kf.ToLocation == "" {
				needDestination = true
			}
		}
	}
	para.TF(" to %s at the appropriate BBS.  The message must be received there between %s and %s.",
		session.CallSign, session.Start.Format("15:04 on Monday"), session.End.Format("15:04 on Monday"))
	switch len(session.DownBBSes) {
	case 0:
		break
	case 1:
		para.TF(" Do not use or send to %s during this session; it has a simulated outage.", session.DownBBSes[0])
	default:
		para.TF(" Do not use or send to %s during this session; they have simulated outages.", english.Conjoin(session.DownBBSes, "or"))
	}
	if session.Instructions != "" {
		main.R(session.Instructions)
	}
	if msg, ok := session.ModelMsg.(*plaintext.PlainText); ok {
		grid := main.E("div id=plainmodel")
		grid.E("div>Subject:")
		grid.E("div>%s", msg.Subject)
		grid.E("div>Message:")
		grid.E("div>%s", msg.Body)
	}
	main.E("p style=margin-bottom:0>The following references may be helpful to you:")
	list := main.E("ul style=margin-top:0")
	list.E("li>The ").E("a href=https://www.scc-ares-races.org/freqs/packet-freqs.html target=_blank>Packet Frequency and BBS Listings").
		P().TF(" will tell you which BBS to use to reach %s.", session.CallSign)
	list.E("li>The “Standard Outpost Configuration Instructions”, available on the ").
		E("a href=https://www.scc-ares-races.org/data/packet/index.html target=_blank>Packet BBS Service").
		P().R(" page, will tell you how to configure Outpost to send messages following county standards.")
	if plainText {
		list.E("li>The “Standard Packet Message Subject Line”, available on the ").
			E("a href=https://www.scc-ares-races.org/data/packet/index.html target=_blank>Packet BBS Service").
			P().R(" page, will tell you how to construct the subject line of your message to conform to county standards.")
	}
	if needHandling || needDestination {
		li := list.E("li>The “RACES Recommended Form Routing Cheat Sheet”, available on the ")
		li.E("a href=https://www.scc-ares-races.org/operations/go-kit-forms.html target=_blank>Go Kit Forms")
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
		E("a href=https://www.scc-ares-races.org/discuss-groups.html target=_blank>Discussion Groups").
		P().R(" page.")
	count := ws.st.ModelImageCount(session.ID)
	for pnum := 1; pnum <= count; pnum++ {
		html.E("img class=modelimage src=/session/image?session=%d&page=%d", session.ID, pnum)
	}
}
