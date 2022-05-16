package webserver

import (
	"bytes"
	_ "embed" // -
	"fmt"
	"net/http"
	"strings"
	"time"

	"steve.rothskeller.net/packet/wppsvr/store"
)

//go:embed "message.html"
var messageHTML []byte

func (ws *webserver) serveMessage(w http.ResponseWriter, r *http.Request) {
	var (
		callsign  string
		message   *store.Message
		responses []*store.Response
		content   int
	)
	if callsign = checkLoggedIn(w, r); callsign == "" {
		return
	}
	if message = ws.st.GetMessage(r.FormValue("id")); message == nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if message.FromCallSign != callsign && !canViewEveryone(callsign) {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}
	responses = ws.st.GetResponses(message.LocalID)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	content = bytes.Index(messageHTML, contentMarker)
	w.Write(messageHTML[:content])
	fmt.Fprintf(w, `<div class="heading">Received Message %s</div><pre class="message">%s</pre>`,
		message.LocalID, message.Message)
	for i, resp := range responses {
		fmt.Fprintf(w, `<div class="heading">Response #%d</div><pre class="message">From: %s@%s.ampr.org
To: %s
Subject: %s
Date: %s

%s</pre>`, i+1, strings.ToLower(resp.SenderCall), strings.ToLower(resp.SenderBBS),
			resp.To, resp.Subject, resp.SendTime.Format(time.RFC822), resp.Body)
	}
	w.Write(messageHTML[content+len(contentMarker):])
}
