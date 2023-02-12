package webserver

import (
	"bytes"
	_ "embed" // -
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/report"
	"github.com/rothskeller/packet/wppsvr/store"
)

// contentMarker is the Marker in calendar.html for the place where the variable
// content should be inserted.
var contentMarker = []byte("@@CONTENT@@")

//go:embed "calendar.html"
var calendarHTML []byte

var callsignRE = regexp.MustCompile(`(?i)^(?:A[A-L]|[KNW][A-Z]?)[0-9][A-Z]{1,3}$`)

// serveCalendar handles GET /calendar requests.
func (ws *webserver) serveCalendar(w http.ResponseWriter, r *http.Request) {
	var (
		content  int
		view     string
		callsign string
		year     = time.Now().Year()
	)
	if callsign = checkLoggedIn(w, r); callsign == "" {
		return
	}
	// What are we trying to view?  The check-in counts, or the results for
	// a particular call sign?  And for which year?
	view = r.FormValue("view")
	if view != "counts" {
		if callsignRE.MatchString(view) && isAllowedToView(callsign, view) {
			view = strings.ToUpper(view)
		} else {
			view = callsign
		}
	}
	if y, err := strconv.Atoi(r.FormValue("year")); err == nil && y > 2000 && y < 3000 {
		year = y
	}
	// Write the preamble.
	w.Header().Set("Cache-Control", "nostore")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	content = bytes.Index(calendarHTML, contentMarker)
	w.Write(calendarHTML[:content])
	// Write the view options.
	if view == "counts" {
		fmt.Fprintf(w, `<div id="view">Viewing Check-In Counts&nbsp&nbsp;|&nbsp; <a href="?year=%d&view=%s">View %s Results</a></div>`,
			year, callsign, callsign)
	} else {
		fmt.Fprintf(w, `<div id="view">Viewing %s Results&nbsp&nbsp;|&nbsp; <a href="?year=%d&view=counts">View Check-In Counts</a></div>`,
			view, year)
	}
	// Write the year selector.
	if ws.yearHasSessions(year-1) || canEditSessions(callsign) {
		fmt.Fprintf(w, `<div id="year"><div id="lastyear"><a href="?year=%d&view=%s">&lt; %d</a></div>`, year-1, view, year-1)
	} else {
		fmt.Fprintf(w, `<div id="year"><div id="lastyear">&lt; %d</div>`, year-1)
		// The year doesn't actually appear (color: transparent), but
		// this makes the spacing correct.
	}
	fmt.Fprintf(w, `<div id="thisyear">%d</div>`, year)
	if ws.yearHasSessions(year+1) || canEditSessions(callsign) {
		fmt.Fprintf(w, `<div id="nextyear"><a href="?year=%d&view=%s">%d &gt;</a></div></div>`, year+1, view, year+1)
	} else {
		fmt.Fprintf(w, `<div id="nextyear">%d &gt;</div></div>`, year+1)
	}
	// Write the months.
	io.WriteString(w, `<div id="calendar">`)
	for month := time.January; month <= time.December; month++ {
		ws.serveCalendarMonth(w, r, year, month, view)
	}
	io.WriteString(w, `</div>`)
	// Close out the HTML.
	w.Write(calendarHTML[content+len(contentMarker):])
}

// isAllowedToView returns whether the viewer (identified by callsign) is
// allowed to view the specified view (also a call sign).  It returns true if
// the two are the same call sign or if the viewer is in the privilege list.
func isAllowedToView(callsign, view string) bool {
	if strings.EqualFold(callsign, view) {
		return true
	}
	return canViewEveryone(callsign)
}

func (ws *webserver) yearHasSessions(year int) bool {
	return ws.st.ExistRealizedSessions(
		time.Date(year, 1, 1, 0, 0, 0, 0, time.Local),
		time.Date(year+1, 1, 1, 0, 0, 0, 0, time.Local),
	)
}

func (ws *webserver) serveCalendarMonth(w http.ResponseWriter, r *http.Request, year int, month time.Month, view string) {
	var (
		date     time.Time
		sessions []*store.Session
	)
	sessions = ws.st.GetSessions(
		time.Date(year, month, 1, 0, 0, 0, 0, time.Local),
		time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local),
	)
	fmt.Fprintf(w, `<div class="month"><div class="monthname">%s</div><div class="weekday">S</div><div class="weekday">M</div><div class="weekday">T</div><div class="weekday">W</div><div class="weekday">T</div><div class="weekday">F</div><div class="weekday">S</div>`, month.String())
	date = time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	for i := 0; i < int(date.Weekday()); i++ {
		io.WriteString(w, `<div></div>`)
	}
	for ; date.Month() == month; date = date.AddDate(0, 0, 1) {
		var (
			session *store.Session
			ciClass string
			ciValue string
		)
		if len(sessions) != 0 && sameDay(sessions[0].End, date) {
			session = sessions[0]
			for len(sessions) != 0 && sameDay(sessions[0].End, date) {
				sessions = sessions[1:]
			}
		}
		if session == nil {
			fmt.Fprintf(w, `<div class="day"><div class="date">%d</div></div>`, date.Day())
			continue
		}
		if view == "counts" {
			ciClass = "count"
			ciValue = strconv.Itoa(ws.sessionCheckInCount(session))
		} else {
			ciClass, ciValue = ws.calendarCell(session, view)
		}
		fmt.Fprintf(w, `<div class="day net"><div class="date"><a href="report?date=%s">%d</a></div><div class="%s">%s</div></div>`,
			date.Format("2006-01-02"), date.Day(), ciClass, ciValue)
	}
	for i := int(date.Weekday()); i%7 != 0; i++ {
		io.WriteString(w, `<div></div>`)
	}
	io.WriteString(w, `</div>`)
}

// sessionCheckInCount determines the check-in count for the session, for
// display in the calendar.
func (ws *webserver) sessionCheckInCount(session *store.Session) int {
	if session.ID == 0 {
		return 0
	}
	rpt := report.Generate(ws.st, session)
	return rpt.UniqueCallSigns
}

// calendarCell returns the classname and value for a calendar cell, based on
// whether the specified call sign checked into the specified session, and
// whether they did so with or without error.  Only the last check-in from each
// distinct from address counts in determining whether there was error.
func (ws *webserver) calendarCell(session *store.Session, callsign string) (class, value string) {
	class, value = "noci", "—"
	if session.ID == 0 {
		return
	}
	var fromAddrs = make(map[string]bool)
	for _, message := range ws.st.GetSessionMessages(session.ID) {
		if message.FromCallSign != callsign {
			continue
		}
		if message.Actions&(config.ActionDontCount|config.ActionDropMsg) != 0 {
			continue
		}
		class, value = "ok", "✓"
		fromAddrs[message.FromAddress] = message.Actions&config.ActionError != 0
	}
	for _, haserr := range fromAddrs {
		if haserr {
			class, value = "error", "✕"
		}
	}
	return
}

func sameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}
