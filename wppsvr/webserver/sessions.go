package webserver

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/wppsvr/htmlb"
)

// serveSessionList displays the list of defined sessions and allows them to be
// edited.
func (ws *webserver) serveSessionList(w http.ResponseWriter, r *http.Request) {
	var (
		callsign string
		year     = time.Now().Year()
	)
	if callsign = checkLoggedIn(w, r); callsign == "" {
		return
	}
	if !canEditSessions(callsign) {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}
	// Which year are we trying to view?
	if y, err := strconv.Atoi(r.FormValue("year")); err == nil && y > 2000 && y < 3000 {
		year = y
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
	html.E("link rel=stylesheet href=/static/sessions.css")
	html.E("div id=org>Santa Clara County ARES<sup>®</sup>/RACES")
	html.E("div id=title>Weekly Packet Practice")
	html.E("div id=subtitle>Practice Sessions Editor")
	// Write the year selector.
	yearsel := html.E("div id=year")
	yearsel.E("div id=lastyear").E("a href=?year=%d>< %d", year-1, year-1)
	yearsel.E("div id=thisyear>%d", year)
	yearsel.E("div id=nextyear").E("a href=?year=%d>%d >", year+1, year+1)
	// Start the table.
	table := html.E("table id=sessions")
	tr := table.E("tr")
	tr.E("th>Date")
	tr.E("th>Name")
	tr.E("th>Down")
	tr.E("th>Message Type")
	// Show the session data.
	for _, session := range ws.st.GetSessions(
		time.Date(year, 1, 1, 0, 0, 0, 0, time.Local),
		time.Date(year+1, 1, 1, 0, 0, 0, 0, time.Local),
	) {
		tr = table.E("tr")
		tr.E("td").E("a href=/session?id=%d>%s", session.ID, session.End.Format("Jan 02"))
		tr.E("td>%s", session.Name)
		tr.E("td>%s", strings.Join(session.DownBBSes, ", "))
		if session.ModelMsg != nil {
			tr.E("td>Model %s", session.ModelMsg.Type().Tag)
		} else {
			tr.E("td>Any %s", strings.Join(session.MessageTypes, ", "))
		}
	}
	// Creation hint.
	html.E("div id=newhint>To create a new session, click on an existing similar session,<br>modify it, and click “Save Copy”.")
}
