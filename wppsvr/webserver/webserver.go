package webserver

import (
	"net/http"

	"steve.rothskeller.net/packet/wppsvr/store"
)

// Run starts the web server running in the background.
func Run(st *store.Store) (err error) {
	var ws webserver

	ws.st = st
	http.Handle("/", http.HandlerFunc(ws.serveFrontPage))
	http.Handle("/calendar", http.HandlerFunc(ws.serveCalendar))
	http.Handle("/login", http.HandlerFunc(ws.serveLogin))
	go http.ListenAndServe("localhost:8000", nil)
	return nil
}

type webserver struct {
	st *store.Store
}
