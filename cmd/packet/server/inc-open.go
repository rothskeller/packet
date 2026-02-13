package server

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/osdep"
	"github.com/rothskeller/packet/form/htmlop"
	"github.com/rothskeller/packet/incident"
	"golang.org/x/net/html"
)

//go:embed inc-open.html
var incOpenHTML []byte

// serveIncidentOpen handles GET and POST /incident-open requests.  They should
// have, at minimum, a dir= parameter indicating the current directory to start
// or continue browsing from.  (If not, the user's home directory is used.)
func (s *Server) serveIncidentOpen(w http.ResponseWriter, r *http.Request) {
	var (
		subdirs []string
		doc     *html.Node
		err     error
		vars    = make(map[string]string)
		dir     = r.FormValue("dir")
	)
	if maybeShowREADME(w, r) {
		return
	}
	if dir == "" {
		dir = r.FormValue("return")
	}
	if dir == "" {
		dir = osdep.HomeDir
	}
	if r.FormValue("mkdir") != "" && r.FormValue("mkname") != "" {
		dir = filepath.Join(dir, r.FormValue("mkname"))
		if err = os.Mkdir(dir, 0777); err != nil {
			vars["error"] = err.Error()
		}
	} else if subdir := r.FormValue("subdir"); subdir != "" {
		dir = filepath.Join(dir, subdir)
	} else if r.FormValue("open") != "" {
		if incident.IsIncident(dir) {
			http.Redirect(w, r, "/incident?dir="+url.QueryEscape(dir), http.StatusSeeOther)
			return
		} else if incident.IsUnsafeIncidentDir(dir) {
			vars["error"] = fmt.Sprintf("%s is not a proper directory to create an incident in.  Make a subdirectory with an incident-specific name instead.", dir)
		} else if !isWritable(dir) {
			vars["error"] = dir + " is not writable."
		} else if err = incident.Create(dir, func(i *incident.Incident) error { return nil }); err != nil {
			vars["error"] = fmt.Sprintf("Creating incident in %s: %s", dir, err)
		} else {
			http.Redirect(w, r, "/incident-config?dir="+url.QueryEscape(dir), http.StatusSeeOther)
			return
		}
	}
	vars["DIR"] = dir
	if ents, err := os.ReadDir(dir); err != nil {
		vars["error"] = err.Error()
	} else {
		for _, ent := range ents {
			if ent.IsDir() && !strings.HasPrefix(ent.Name(), ".") {
				subdirs = append(subdirs, ent.Name())
			}
		}
		slices.Sort(subdirs)
		vars["subdirs"] = strings.Join(subdirs, ";")
	}
	if ret := r.FormValue("return"); ret != "" {
		if r.FormValue("cancel") != "" {
			http.Redirect(w, r, "/incident?dir="+url.QueryEscape(ret), http.StatusSeeOther)
			return
		}
		vars["return"] = ret
	}
	if parent := filepath.Dir(dir); parent != dir {
		vars["parent"] = "true"
	}
	if vars["parent"] != "" || vars["subdirs"] != "" {
		vars["browse"] = "true"
	}
	if isWritable(dir) {
		vars["mkdir"] = "true"
	}
	if doc, err = htmlop.Parse(bytes.NewReader(incOpenHTML)); err != nil {
		ErrPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	htmlop.Expand(doc, vars)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	htmlop.Minify(w, doc)
}

func isWritable(dir string) bool {
	if fh, err := os.CreateTemp(dir, ""); err != nil {
		return false
	} else {
		fh.Close()
		os.Remove(fh.Name())
		return true
	}
}
