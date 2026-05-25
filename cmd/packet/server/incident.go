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
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/cmd/packet/osdep"
	"github.com/rothskeller/packet/form/htmlop"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/packetver"
	"golang.org/x/net/html"
)

//go:embed incident.html
var incidentHTML []byte

type incidentData struct {
	BBS         string        `json:"bbs"`
	Dir         string        `json:"dir"`
	GenDRs      bool          `json:"genDRs"`
	Ident       string        `json:"ident"`
	Manual      bool          `json:"manual"`
	MsgTypes    []messageType `json:"msgtypes"`
	Name        string        `json:"name"`
	ServerPrint bool          `json:"serverPrint"`
	Version     string        `json:"version"`
	ViewFlags   viewFlags     `json:"viewFlags"`
}
type messageType struct {
	Tag  string `json:"tag"`
	Key  string `json:"key"`
	Name string `json:"name"`
}
type viewFlags struct {
	Full     bool `json:"full"`
	Large    bool `json:"large"`
	Receipts bool `json:"receipts"`
}

// serveGetIncident handles GET /incident requests.  They will have a dir=
// parameter specifying the incident directory.
func (s *Server) serveGetIncident(w http.ResponseWriter, r *http.Request) {
	var (
		dir  string
		doc  *html.Node
		err  error
		data incidentData
		vars = map[string]string{}
	)
	if maybeShowREADME(w, r) {
		return
	}
	if dir = r.FormValue("dir"); dir == "" {
		ErrPage(w, "The GET /incident request is missing the required dir= parameter.", http.StatusBadRequest)
		return
	}
	err = incident.Write(dir, func(i *incident.Incident) error {
		i.UpdateIncDefaults() // marks incident as recently used
		data.Dir = dir
		data.Ident = i.Config.ActiveCall()
		data.BBS = i.Config.ConnectBBS
		if i.Config.ConnectType == incident.ConnectNone {
			data.Manual = true
			vars["manual"] = "true"
		}
		if i.Config.AllowVoice {
			vars["allowvoice"] = "true"
		}
		if i.Config.IncidentName != "" {
			data.Name = i.Config.IncidentName
		} else if i.Config.ActivationNum != "" {
			data.Name = i.Config.ActivationNum
		} else {
			data.Name = dir
		}
		data.GenDRs = !i.Config.NoSendReceipts
		data.ViewFlags.Full = i.Config.ViewFlags&incident.ViewFull != 0
		data.ViewFlags.Large = i.Config.ViewFlags&incident.ViewLarge != 0
		data.ViewFlags.Receipts = i.Config.ViewFlags&incident.ViewReceipts != 0
		return nil
	})
	if err != nil {
		ErrPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Version = packetver.Version
	data.MsgTypes = newMessageTypeList()
	data.ServerPrint = osdep.PrintPDFCommand("x") != nil
	if by, _ := json.Marshal(data); true {
		vars["INCIDENTDATA"] = string(by)
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
// entries with a higher sequence number.  However, it will wait for at most
// five minutes, so that the client sends a new request at least that often, so
// that the idle timeout doesn't trigger.
func (s *Server) serveGetIncidentLog(w http.ResponseWriter, r *http.Request) {
	var (
		dir    string
		seq    int
		ctx    context.Context
		cancel func()
		upd    ilogUpdate
		err    error
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
	ctx, cancel = context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()
	if err = incident.Watch(ctx, dir, seq, s.stop); err == context.Canceled {
		return
	} else if err != nil && err != context.DeadlineExceeded {
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
func newMessageTypeList() (types []messageType) {
	var list []message.MType
	for mt := range message.AllTypes() {
		if _, ok := mt.(message.EditableMType); ok {
			list = append(list, mt)
		}
	}
	slices.SortFunc(list, message.CompareTypes)
	for _, mt := range list {
		emt := mt.(message.EditableMType)
		_, name, _ := strings.Cut(mt.Name(), " ") // drop "a" or "an"
		name = titleCaseRE.ReplaceAllStringFunc(name, strings.ToUpper)
		types = append(types, messageType{Tag: emt.CreateTag(), Key: emt.CreateKey(), Name: name})
	}
	return types
}

func (s *Server) servePostViewICS309(w http.ResponseWriter, r *http.Request) {
	var (
		dir   string
		fname string
	)
	signature := r.FormValue("signature")
	serveIncident(w, r, false, func(i *incident.Incident) (err error) {
		dir = i.Dir
		fname = "ICS-309"
		if i.Config.ActivationNum != "" {
			fname += " " + i.Config.ActivationNum
		}
		if i.Config.TacCall != "" {
			fname += " " + i.Config.TacCall
		} else if i.Config.OpCall != "" {
			fname += " " + i.Config.OpCall
		}
		if !i.Config.OpStart.IsZero() {
			fname += i.Config.OpStart.Format(" 2006-01-02")
		}
		fname += ".pdf"
		return i.GenerateICS309(signature)
	}, func() error {
		if fh, err := os.Open(filepath.Join(dir, "ics309.pdf")); err != nil {
			return err
		} else {
			w.Header().Set("Content-Type", "application/pdf")
			w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", fname))
			http.ServeContent(w, r, "ics309.pdf", time.Time{}, fh)
			fh.Close()
			return nil
		}
	})
}
