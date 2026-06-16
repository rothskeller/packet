package server

import (
	"fmt"
	"html"
	"net/http"

	"github.com/rothskeller/packet/v4/form/formdefs"
)

// maybeShowREADME looks to see if there is a README.txt in the root of the
// forms directory.  If so, it serves it, deletes it, and returns true.  If not,
// it does nothing and returns false.
func maybeShowREADME(w http.ResponseWriter, r *http.Request) bool {
	readme := formdefs.GetFlushREADME()
	if readme != "" {
		w.Header().Set("Cache-Control", "no-store, private")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!DOCTYPE html><title>Updated Forms</title><h1>Updated Forms</h1><pre>%s</pre><div style="margin-top:2rem"><button onclick="location.href='%s'">Continue</button></div>`,
			html.EscapeString(readme), html.EscapeString(r.URL.String()))
	}
	return readme != ""
}
