package server

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"
)

// ErrorPage emits an error page with the specified HTTP status and error
// message.  If state is non-nil, it is JSON-encoded and shown on the error
// page.
func ErrorPage(w http.ResponseWriter, httpStatus int, err error, state any) {
	var (
		errorString string
		stateString string
	)
	errorString = html.EscapeString(err.Error())
	errorString = strings.ReplaceAll(errorString, "\n", "<br/>")
	if state != nil {
		if by, err := json.Marshal(state); err == nil {
			stateString = string(by)
			if strings.HasPrefix(stateString, "{") {
				// Enable line wrapping.
				stateString = strings.ReplaceAll(stateString, `","`, `", "`)
			}
			stateString = "<br/>" + html.EscapeString(stateString)
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(httpStatus)
	fmt.Fprintf(w, `<html><head><title>Problem</title></head><body><h3 id="something-went-wrong"><span style="font-size:24pt;color:red">⚠︎ </span>Something went wrong.</h3>This information might help resolve the problem:<br/><br/>%s%s</body></html>`,
		errorString, stateString)
}
