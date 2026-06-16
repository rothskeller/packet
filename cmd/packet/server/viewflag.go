package server

import (
	"net/http"

	"github.com/rothskeller/packet/v4/incident"
)

// servePostSetViewFlag handles POST /set-view-flag requests, which carry a
// dir= parameter identifying the incident and a compact=, receipts=, or large=
// parameter set to a boolean value.  It stores the updated flag value and
// returns 204 on success, or an error with text/plain error message on failure.
func (s *Server) servePostSetViewFlag(w http.ResponseWriter, r *http.Request) {
	s.outpost = false
	dir := r.FormValue("dir")
	err := incident.Write(dir, func(i *incident.Incident) error {
		switch r.FormValue("compact") {
		case "true":
			i.Config.ViewFlags &^= incident.ViewFull
		case "false":
			i.Config.ViewFlags |= incident.ViewFull
		}
		switch r.FormValue("large") {
		case "false":
			i.Config.ViewFlags &^= incident.ViewLarge
		case "true":
			i.Config.ViewFlags |= incident.ViewLarge
		}
		switch r.FormValue("receipts") {
		case "false":
			i.Config.ViewFlags &^= incident.ViewReceipts
		case "true":
			i.Config.ViewFlags |= incident.ViewReceipts
		}
		i.UpdateIncDefaults()
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
