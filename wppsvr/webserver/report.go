package webserver

import (
	"io"
	"net/http"
	"strconv"

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
	if sid, err := strconv.Atoi(r.FormValue("session")); err == nil {
		session = ws.st.GetSession(sid)
	}
	if session == nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if session.Flags&store.Imported != 0 {
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
