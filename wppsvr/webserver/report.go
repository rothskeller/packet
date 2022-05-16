package webserver

import (
	"io"
	"net/http"
	"time"

	"steve.rothskeller.net/packet/wppsvr/report"
	"steve.rothskeller.net/packet/wppsvr/store"
)

func (ws *webserver) serveReport(w http.ResponseWriter, r *http.Request) {
	var session *store.Session

	if checkLoggedIn(w, r) == "" {
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
	if session.Report == "" {
		session.Report = report.Generate(ws.st, session).RenderHTML()
	}
	if session.Report[0] == '<' {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	} else {
		w.Header().Set("Content-Type", "text/plain")
	}
	io.WriteString(w, session.Report)
}
