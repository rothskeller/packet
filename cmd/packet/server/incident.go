package server

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json/v2"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rothskeller/packet/form/htmlop"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/packetver"
	"golang.org/x/net/html"
)

//go:embed incident.html
var incidentHTML []byte

var serverPrintCmd string
var serverPrintOnce sync.Once

// serveGetIncident handles GET /incident requests.  They will have a dir=
// parameter specifying the incident directory.
func (s *Server) serveGetIncident(w http.ResponseWriter, r *http.Request) {
	var (
		dir  string
		doc  *html.Node
		err  error
		vars = map[string]string{}
	)
	if maybeShowREADME(w, r) {
		return
	}
	if dir = r.FormValue("dir"); dir == "" {
		ErrPage(w, "The GET /incident request is missing the required dir= parameter.", http.StatusBadRequest)
		return
	}
	err = incident.Read(dir, func(i *incident.Incident) error {
		vars["DIR"] = dir
		vars["IDENT"] = i.Config.ActiveCall()
		vars["BBS"] = i.Config.ConnectBBS
		if i.Config.ConnectType == incident.ConnectNone {
			vars["manual"] = "true"
		}
		if i.Config.IncidentName != "" {
			vars["INCNAME"] = i.Config.IncidentName
		} else if i.Config.ActivationNum != "" {
			vars["INCNAME"] = i.Config.ActivationNum
		} else {
			vars["INCNAME"] = dir
		}
		if !i.Config.NoSendReceipts {
			vars["GENDRS"] = "checked"
		}
		vars["VIEWFLAGS"] = i.Config.ViewFlags.String()
		return nil
	})
	if err != nil {
		ErrPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vars["VERSION"] = packetver.Version
	vars["MTYPES"] = newMessageTypeList()
	serverPrintOnce.Do(setServerCanPrint)
	if serverPrintCmd != "" {
		vars["SERVERPRINT"] = "true"
	}
	if doc, err = html.Parse(bytes.NewReader(incidentHTML)); err != nil {
		slog.Error("parse incident HTML", "err", err)
		ErrPage(w, "The incident.html page could not be parsed.  Please report this error to the author.", http.StatusInternalServerError)
		return
	}
	htmlop.Expand(doc, vars)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html.Render(w, doc)
}

type ilogUpdate struct {
	Seq          int                  `json:"seq"`
	Log          []*incident.LogEntry `json:"log"`
	ConnProgress string               `json:"connProgress"`
	ConnError    string               `json:"connError"`
}

// serveGetIncidentLog handles GET /incident/log requests.  They will have a
// dir= parameter specifying the incident directory and a seq= parameter
// specifying the sequence number of the data already held by the client.  This
// is a long polling request; it will wait until there is data available for
// a higher sequence number for the incident and then return the list of log
// entries with a higher sequence number.
func (s *Server) serveGetIncidentLog(w http.ResponseWriter, r *http.Request) {
	var (
		dir string
		seq int
		upd ilogUpdate
		err error
	)
	if dir = r.FormValue("dir"); dir == "" {
		slog.Error("no incident dir")
		http.Error(w, "dir is required", http.StatusBadRequest)
		return
	} else if !incident.IsIncident(dir) {
		slog.Error("no such incident dir", "dir", dir)
		http.Error(w, "dir is not an incident", http.StatusBadRequest)
		return
	}
	if seq, err = strconv.Atoi(r.FormValue("seq")); err != nil || seq < 0 {
		slog.Error("bad seq number", "seq", r.FormValue("seq"))
		http.Error(w, "seq is missing or invalid", http.StatusBadRequest)
		return
	}
	if err = incident.Watch(r.Context(), dir, seq, s.stop); err == context.Canceled {
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = incident.Read(dir, func(i *incident.Incident) error {
		upd.Seq = i.Seq
		for _, e := range i.Log {
			if e.Seq > seq {
				upd.Log = append(upd.Log, e)
			}
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, private")
	json.MarshalWrite(w, upd)
}

var titleCaseRE = regexp.MustCompile(`(?:^|[- ])[a-z]`)

// newMessageTypeList returns a semicolon-separated string of creatable message
// types.  Each element is a colon-separated string of tag, key, and name in
// title case.
func newMessageTypeList() string {
	var list []message.MType
	for mt := range message.AllTypes() {
		if _, ok := mt.(message.EditableMType); ok {
			list = append(list, mt)
		}
	}
	slices.SortFunc(list, message.CompareTypes)
	var data []string
	for _, mt := range list {
		emt := mt.(message.EditableMType)
		_, name, _ := strings.Cut(mt.Name(), " ") // drop "a" or "an"
		name = titleCaseRE.ReplaceAllStringFunc(name, strings.ToUpper)
		data = append(data, fmt.Sprintf("%s:%s:%s", emt.CreateTag(), emt.CreateKey(), name))
	}
	return strings.Join(data, ";")
}

func setServerCanPrint() {
	var err error

	if serverPrintCmd, err = exec.LookPath("lpr"); err != nil || serverPrintCmd == "" {
		if serverPrintCmd, err = exec.LookPath("lp"); err != nil {
			serverPrintCmd = ""
		}
	}
}

func (s *Server) servePostViewICS309(w http.ResponseWriter, r *http.Request) {
	var dir string

	signature := r.FormValue("signature")
	serveIncident(w, r, false, func(i *incident.Incident) (err error) {
		dir = i.Dir
		return i.GenerateICS309(signature)
	}, func() error {
		if fh, err := os.Open(filepath.Join(dir, "ICS-309.pdf")); err != nil {
			return err
		} else {
			w.Header().Set("Content-Type", "application/pdf")
			http.ServeContent(w, r, "ICS-309.pdf", time.Time{}, fh)
			fh.Close()
			return nil
		}
	})
}
