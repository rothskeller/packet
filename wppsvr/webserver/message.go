package webserver

import (
	_ "embed" // -
	"net/http"
	"strings"

	"github.com/rothskeller/packet/wppsvr/htmlb"
	"github.com/rothskeller/packet/wppsvr/store"
)

func (ws *webserver) serveMessage(w http.ResponseWriter, r *http.Request) {
	var (
		callsign string
		msg      *store.Message
	)
	if hash := r.FormValue("hash"); hash != "" {
		if msg = ws.st.GetMessageByHash(hash); msg == nil {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
	} else {
		if callsign = checkLoggedIn(w, r); callsign == "" {
			return
		}
		if msg = ws.st.GetMessage(r.FormValue("id")); msg == nil {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
		if msg.FromCallSign != callsign && !canViewEveryone(callsign) {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return
		}
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
	html.E("link rel=stylesheet href=/static/message.css")
	html.E("div id=org>Santa Clara County ARES<sup>®</sup>/RACES")
	html.E("div id=title>Weekly Packet Practice")
	html.E("div id=subtitle>Message Evaluation")
	body := html.E("div id=results")
	// Display the summary.
	lr := body.E("div id=lr")
	var fromcall = msg.FromCallSign
	if strings.HasPrefix(strings.ToUpper(msg.FromAddress), fromcall) {
		fromcall = ""
	}
	if fromcall != "" {
		lr.E("div id=msgid>%s from %s (%s)", msg.LocalID, msg.FromAddress, fromcall)
	} else {
		lr.E("div id=msgid>%s from %s", msg.LocalID, msg.FromAddress)
	}
	lr.E("div id=score>Score: %d%%", msg.Score)
	body.E("div id=rawmsg>%s", msg.Message)
	if msg.Analysis != "" {
		body.R(msg.Analysis)
		body.E("h2>For Assistance")
		body.E(`p>If you need assistance, the best place to request it is the <kbd>packet@scc-ares-races.groups.io</kbd> discussion group.  To sign up for this group, see the <a href="https://www.scc-ares-races.org/discuss-groups.html">Discussion Groups</a> page on the county ARES website.</p>`)
	} else {
		body.E("h2>No Issues Found")
		body.E("p>Thank you for checking in to the net successfully.")
	}
}
