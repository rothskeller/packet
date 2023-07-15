package webserver

import (
	"embed"
	"net/http"

	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/store"
)

//go:embed *.css *.js
var static embed.FS

// Run starts the web server running in the background.
func Run(st *store.Store) (err error) {
	var ws webserver

	ws.st = st
	http.Handle("/", http.HandlerFunc(ws.serveFrontPage))
	http.Handle("/calendar", http.HandlerFunc(ws.serveCalendar))
	http.Handle("/instructions", http.HandlerFunc(ws.serveInstructions))
	http.Handle("/login", http.HandlerFunc(ws.serveLogin))
	http.Handle("/message", http.HandlerFunc(ws.serveMessage))
	http.Handle("/report", http.HandlerFunc(ws.serveReport))
	http.Handle("/session", http.HandlerFunc(ws.serveSessionEdit))
	http.Handle("/session/image", http.HandlerFunc(ws.serveModelImage))
	http.Handle("/sessions", http.HandlerFunc(ws.serveSessionList))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))
	go http.ListenAndServe(config.Get().ListenAddr, nil)
	return nil
}

type webserver struct {
	st *store.Store
}
