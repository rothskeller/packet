package server

import (
	"fmt"
	"html"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rothskeller/packet/form/formdefs"
)

// maybeShowREADME looks to see if there is a README.txt in the root of the
// forms directory.  If so, it serves it, deletes it, and returns true.  If not,
// it does nothing and returns false.
func maybeShowREADME(w http.ResponseWriter, r *http.Request) bool {
	var (
		rmfile string
		data   []byte
		err    error
	)
	rmfile = filepath.Join(formdefs.FormsDir(), "README.txt")
	if data, err = os.ReadFile(rmfile); os.IsNotExist(err) {
		return false
	} else if err != nil {
		slog.Error("os.ReadFile", "f", rmfile, "err", err)
		return false
	}
	os.Remove(rmfile)
	w.Header().Set("Cache-Control", "no-store, private")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html><title>Updated Forms</title><h1>Updated Forms</h1><pre>%s</pre><div style="margin-top:2rem"><button onclick="location.href='%s'">continue</button></div>`,
		html.EscapeString(string(data)), html.EscapeString(r.URL.String()))
	return true
}
