package webserver

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

func (ws *webserver) serveModelImage(w http.ResponseWriter, r *http.Request) {
	var (
		fh *os.File
	)
	if sid, err := strconv.Atoi(r.FormValue("session")); err == nil {
		if pnum, err := strconv.Atoi(r.FormValue("page")); err == nil {
			fh = ws.st.ModelImage(sid, pnum)
		}
	}
	if fh == nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	defer fh.Close()
	w.Header().Set("Cache-Control", "nostore")
	http.ServeContent(w, r, fh.Name(), time.Time{}, fh)
}
