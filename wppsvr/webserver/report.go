package webserver

import (
	"io"
	"net/http"
	"time"

	"github.com/rothskeller/packet/wppsvr/report"
	"github.com/rothskeller/packet/wppsvr/store"
)

func (ws *webserver) serveReport(w http.ResponseWriter, r *http.Request) {
	var (
		callsign string
		session  *store.Session
	)
	if callsign = checkLoggedIn(w, r); callsign == "" {
		return
	}
	if date, err := time.ParseInLocation("2006-01-02", r.FormValue("date"), time.Local); err == nil {
		if sessions := ws.st.GetSessions(date, date.AddDate(0, 0, 1)); len(sessions) != 0 {
			session = sessions[0]
		}
	}
	if session == nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if session.Imported {
		// This is a report imported from the old NCO scripts.  Display
		// its report verbatim.
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, session.Report)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if canViewEveryone(callsign) {
		io.WriteString(w, report.Generate(ws.st, session).RenderHTML(""))
	} else {
		io.WriteString(w, report.Generate(ws.st, session).RenderHTML(callsign))
	}
}
