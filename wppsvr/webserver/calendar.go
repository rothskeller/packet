package webserver

import (
	_ "embed" // -
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/wppsvr/htmlb"
	"github.com/rothskeller/packet/wppsvr/report"
	"github.com/rothskeller/packet/wppsvr/store"
)

var callsignRE = regexp.MustCompile(`(?i)^(?:A[A-L][0-9][A-Z]{1,3}|[KNW][A-Z][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3})$`)

// serveCalendar handles GET /calendar requests.
func (ws *webserver) serveCalendar(w http.ResponseWriter, r *http.Request) {
	var (
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
	// Start the HTML page.
	w.Header().Set("Cache-Control", "nostore")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := htmlb.HTML(w)
	defer html.Close()
	html.E("meta charset=utf-8")
	html.E("title>Weekly Packet Practice - Santa Clara County ARES/RACES")
	html.E("meta name=viewport content='width=device-width, initial-scale=1'")
	html.E("link rel=stylesheet href=/static/common.css")
	html.E("link rel=stylesheet href=/static/calendar.css")
	html.E("script src=/static/calendar.js")
	html.E("div id=org>Santa Clara County ARES<sup>®</sup>/RACES")
	html.E("div id=title>Weekly Packet Practice")
	html.E("div id=key>Click on any net date to see its report.")
	// Write the view options.
	if view == "counts" {
		html.E("div id=view>Viewing Check-In Counts  |  ").
			E("a href=?year=%d&view=%s>View %s Results", year, callsign, callsign)
	} else {
		html.E("div id=view>Viewing %s Results  |  ", view).
			E("a href=?year=%d&view=counts>View Check-In Counts</a></div>", year)
	}
	// Write the year selector.  Last year and next year are always written
	// so the spacing is correct, but if those years don't have sessions,
	// they appear transparent.
	yearsel := html.E("div id=year")
	lastyear := yearsel.E("div id=lastyear")
	if ws.yearHasSessions(year - 1) {
		lastyear = lastyear.E("a href=?year=%d&view=%s", year-1, view)
	}
	lastyear.TF("< %d", year-1)
	yearsel.E("div id=thisyear>%d", year)
	nextyear := yearsel.E("div id=nextyear")
	if ws.yearHasSessions(year + 1) {
		nextyear = nextyear.E("a href=?year=%d&view=%s", year+1, view)
	}
	nextyear.TF("%d >", year+1)
	// Write the months.
	calendar := html.E("div id=calendar")
	for month := time.January; month <= time.December; month++ {
		ws.serveCalendarMonth(calendar, year, month, view)
	}
	// Give a link to the session editor, for those who can use it.
	if canEditSessions(callsign) {
		html.E("a id=edit href=/sessions>Edit Practice Session Definitions")
	}
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
	return ws.st.ExistSessions(
		time.Date(year, 1, 1, 0, 0, 0, 0, time.Local),
		time.Date(year+1, 1, 1, 0, 0, 0, 0, time.Local),
	)
}

func (ws *webserver) serveCalendarMonth(calendar *htmlb.Element, year int, month time.Month, view string) {
	var (
		date     time.Time
		sessions []*store.Session
	)
	sessions = ws.st.GetSessions(
		time.Date(year, month, 1, 0, 0, 0, 0, time.Local),
		time.Date(year, month+1, 1, 0, 0, 0, 0, time.Local),
	)
	mdiv := calendar.E("div class=month")
	mdiv.E("div class=monthname>%s", month.String())
	mdiv.E("div class=weekday>S")
	mdiv.E("div class=weekday>M")
	mdiv.E("div class=weekday>T")
	mdiv.E("div class=weekday>W")
	mdiv.E("div class=weekday>T")
	mdiv.E("div class=weekday>F")
	mdiv.E("div class=weekday>S")
	date = time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	for i := 0; i < int(date.Weekday()); i++ {
		mdiv.E("div")
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
			mdiv.E("div class=day").E("div class=date>%d", date.Day())
			continue
		}
		if view == "counts" {
			ciClass = "count"
			ciValue = strconv.Itoa(ws.sessionCheckInCount(session))
		} else {
			ciClass, ciValue = ws.calendarCell(session, view)
		}
		ddiv := mdiv.E("div class='day net'")
		ddiv.E("div class=date").E("a href=report?session=%d>%d", session.ID, date.Day())
		ddiv.E("div class=%s>%s", ciClass, ciValue)
	}
	for i := int(date.Weekday()); i%7 != 0; i++ {
		mdiv.E("div")
	}
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
		return "noci", "—"
	}
	var fromAddrs = make(map[string]int)
	var minscore = 0
	for _, message := range ws.st.GetSessionMessages(session.ID) {
		if message.FromCallSign != callsign {
			continue
		}
		if message.Score == 0 {
			continue
		}
		minscore = 100
		fromAddrs[message.FromAddress] = message.Score
	}
	for _, score := range fromAddrs {
		if score < minscore {
			minscore = score
		}
	}
	switch {
	case minscore == 0:
		return "noci", "—"
	case minscore == 100:
		return "ok", "✓"
	case minscore >= 90:
		return "warn", "⚠︎"
	default:
		return "error", "✕"
	}
}

func sameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}
