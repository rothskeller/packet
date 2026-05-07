package webserver

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/wppsvr/report"
	"github.com/rothskeller/packet/wppsvr/store"
)

func (ws *webserver) serveReport(w http.ResponseWriter, r *http.Request) {
	var (
		callsign string
		sessions []*store.Session
		sb       strings.Builder
	)
	if callsign = ws.checkLoggedIn(w, r); callsign == "" {
		return
	}
	if sid, err := strconv.Atoi(r.FormValue("session")); err == nil {
		if session := ws.st.GetSession(sid); session != nil {
			sessions = []*store.Session{session}
		}
	} else if date, err := time.ParseInLocation("2006-01-02", r.FormValue("date"), time.Local); err == nil {
		sessions = ws.st.GetSessions(date, date.AddDate(0, 0, 1))
	}
	if len(sessions) == 0 {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if sessions[0].Flags&store.Imported != 0 {
		// This is a report imported from the old NCO scripts.  Display
		// its report verbatim.
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, sessions[0].Report)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	report.RenderHTMLProlog(&sb)
	for _, session := range sessions {
		rep := report.Generate(ws.st, session)
		if canViewEveryone(callsign) {
			rep.RenderHTMLBody(&sb, "")
		} else {
			rep.RenderHTMLBody(&sb, callsign)
		}
	}
	io.WriteString(w, sb.String())
}
