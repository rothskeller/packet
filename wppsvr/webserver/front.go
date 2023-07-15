package webserver

import (
	"net/http"
	"time"

	"github.com/rothskeller/packet/wppsvr/htmlb"
	"github.com/rothskeller/packet/wppsvr/report"
	"github.com/rothskeller/packet/wppsvr/store"
)

// serveFrontPage handles GET / requests.  These could be either page requests,
// asking for HTML, or fetch requests, asking for JSON.
func (ws *webserver) serveFrontPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	// We will include on the front page any sessions that end between the
	// beginning of today and the end of the day six days forward.
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 0, 7)
	sessions := ws.st.GetSessions(start, end)
	// If any of those sessions end within 4 minutes, our refresh time
	// should be one minute after they end.  (The one minute delay allows
	// time for the final retrieval.)  Otherwise, our refresh time should
	// be 5 minutes.
	refresh := 5 * 60
	for _, session := range sessions {
		delay := int(time.Until(session.End)/time.Second) + 60
		if delay < refresh {
			refresh = delay
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
	html.E("meta http-equiv=refresh content=%d", refresh)
	html.E("link rel=stylesheet href=/static/common.css")
	html.E("link rel=stylesheet href=/static/front.css")
	html.E("script src=/static/front.js")
	html.E("div id=org>Santa Clara County ARES<sup>®</sup>/RACES")
	html.E("div id=title>Weekly Packet Practice")
	// Render the rest of the page.
	ws.renderSessionData(html, sessions)
	renderLoginForm(html)
}

// renderSessionData renders the sessions on the page.
func (ws *webserver) renderSessionData(html *htmlb.Element, sessions []*store.Session) {
	flow := html.E("div id=sessions")
	for _, session := range sessions {
		bubble := flow.E("div class=bubble")
		bubble.E("div class=label>%s", session.Name)
		bubble.E("div class=date>%s", session.End.Format("Monday, January 2"))
		bubble.E("a class=instructions href=/instructions?session=%d", session.ID).R("Instructions")
		rep := report.Generate(ws.st, session)
		bubble.E("div class=count>%d", rep.UniqueCallSigns)
		if rep.UniqueCallSigns != 0 {
			bubble.E("div class=score>%d", rep.AverageValidScore)
		}
		if session.Flags&store.Running != 0 {
			bubble.E("div class=preliminary>preliminary")
		}
	}
}

// renderLoginForm renders the login form on the page.
func renderLoginForm(html *htmlb.Element) {
	html.E("div id=login>For more detail, please log in.")
	form := html.E("form id=form")
	form.E("label for=callsign>Call Sign")
	form.E("input type=text id=callsign name=callsign")
	form.E("label for=password>Password")
	form.E("input type=password id=password name=password")
	form.E("div id=pwdhint>Use your password from scc-ares-races.org.")
	submit := form.E("div id=submitline")
	submit.E("input type=submit value='Log In'")
	submit.E("span id=login-incorrect style=display:none>Login incorrect")
}
