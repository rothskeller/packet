package webserver

import (
	"embed"
	"net/http"
	"net/http/cgi"

	"github.com/rothskeller/packet/v4/wppsvr/config"
	"github.com/rothskeller/packet/v4/wppsvr/store"
)

//go:embed *.css *.js
var static embed.FS

// Run starts the web server running in the background.
func Run(st *store.Store) (err error) {
	var ws webserver

	ws.st = st
	ws.Setup()
	go http.ListenAndServe(config.Get().ListenAddr, nil)
	return nil
}

// HandleCGI handles a single CGI request.
func HandleCGI(st *store.Store) {
	var ws webserver

	ws.st = st
	ws.Setup()
	cgi.Serve(nil)
}

// Setup sets the web server routes.
func (ws *webserver) Setup() {
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
}

type webserver struct {
	st *store.Store
}
