package server

import (
	"fmt"
	"html"
	"net/http"
	"strings"
)

// ErrPage emits an error page with the specified HTTP status and error message.
func (s *Server) ErrPage(w http.ResponseWriter, err string, httpStatus int) {
	err = html.EscapeString(err)
	err = strings.ReplaceAll(err, "\n", "<br/>")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(httpStatus)
	fmt.Fprintf(w, `<html><head><title>Problem</title></head><body><h3 id="something-went-wrong"><span style="font-size:24pt;color:red">⚠︎ </span>Something went wrong.</h3>This information might help resolve the problem:<br/><br/>%s</body></html>`, err)
}
