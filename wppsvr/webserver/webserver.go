package webserver

import (
	"net/http"

	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/store"
)

// Run starts the web server running in the background.
func Run(st *store.Store) (err error) {
	var ws webserver

	ws.st = st
	http.Handle("/", http.HandlerFunc(ws.serveFrontPage))
	http.Handle("/calendar", http.HandlerFunc(ws.serveCalendar))
	http.Handle("/login", http.HandlerFunc(ws.serveLogin))
	http.Handle("/message", http.HandlerFunc(ws.serveMessage))
	http.Handle("/report", http.HandlerFunc(ws.serveReport))
	go http.ListenAndServe(config.Get().ListenAddr, nil)
	return nil
}

type webserver struct {
	st *store.Store
}
