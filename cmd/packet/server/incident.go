package server

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json/v2"
	"errors"
	"net/http"
	"strconv"

	"github.com/rothskeller/packet/form/htmlop"
	"github.com/rothskeller/packet/incident"
	"golang.org/x/net/html"
)

//go:embed incident.html
var incidentHTML []byte

// serveGetIncident handles GET /incident requests.  They will have a dir=
// parameter specifying the incident directory.
func (s *Server) serveGetIncident(w http.ResponseWriter, r *http.Request) {
	var (
		dir  string
		doc  *html.Node
		err  error
		vars = map[string]string{}
	)
	if dir = r.FormValue("dir"); dir == "" {
		ErrorPage(w, http.StatusBadRequest, errors.New("incident directory is a required parameter"), nil)
		return
	}
	err = incident.Read(dir, func(i *incident.Incident) error {
		vars["DIR"] = dir
		if i.Config.IncidentName != "" {
			vars["INCNAME"] = i.Config.IncidentName
		} else if i.Config.ActivationNum != "" {
			vars["INCNAME"] = i.Config.ActivationNum
		} else {
			vars["INCNAME"] = dir
		}
		if i.Config.TacCall != "" {
			vars["IDENT"] = i.Config.TacCall + "@" + i.Config.ConnectBBS
		} else {
			vars["IDENT"] = i.Config.OpCall + "@" + i.Config.ConnectBBS
		}
		if !i.Config.NoSendReceipts {
			vars["GENDRS"] = "checked"
		}
		return nil
	})
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError, err, nil)
		return
	}
	vars["VERSION"] = "4.0.0" // TODO: compute this
	if doc, err = html.Parse(bytes.NewReader(incidentHTML)); err != nil {
		ErrorPage(w, http.StatusInternalServerError, err, nil)
		return
	}
	htmlop.Expand(doc, vars)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html.Render(w, doc)
}

// serveGetIncidentLog handles GET /incident/log requests.  They will have a
// dir= parameter specifying the incident directory and a seq= parameter
// specifying the sequence number of the data already held by the client.  This
// is a long polling request; it will wait until there is data available for
// a higher sequence number for the incident and then return the list of log
// entries with a higher sequence number.
func (s *Server) serveGetIncidentLog(w http.ResponseWriter, r *http.Request) {
	var (
		dir  string
		seq  int
		ilog []*incident.LogEntry
		err  error
	)
	if dir = r.FormValue("dir"); dir == "" {
		http.Error(w, "dir is required", http.StatusBadRequest)
		return
	} else if !incident.IsIncident(dir) {
		http.Error(w, "dir is not an incident", http.StatusBadRequest)
		return
	}
	if seq, err = strconv.Atoi(r.FormValue("seq")); err != nil || seq < 0 {
		http.Error(w, "seq is missing or invalid", http.StatusBadRequest)
		return
	}
	if err = incident.Watch(r.Context(), dir, seq); err == context.Canceled {
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = incident.Read(dir, func(i *incident.Incident) error {
		for _, e := range i.Log {
			if e.Seq > seq {
				ilog = append(ilog, e)
			}
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "nostore, private")
	json.MarshalWrite(w, ilog)
}
