package server

import (
	"bytes"
	"cmp"
	_ "embed"
	"errors"
	"fmt"
	"maps"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/v4/cmd/packet/pseudomsg"
	"github.com/rothskeller/packet/v4/form/htmlop"
	"github.com/rothskeller/packet/v4/incident"
	"github.com/rothskeller/packet/v4/jnos/tnc"
	"golang.org/x/net/html"
)

//go:embed config.html
var configForm []byte

func (s *Server) serveGetIncidentConfig(w http.ResponseWriter, r *http.Request) {
	var (
		doc       *html.Node
		err       error
		defs      *incident.IncDefaults
		variables = make(map[string]string)
		form      = make(url.Values)
	)
	s.outpost = false
	if maybeShowREADME(w, r) {
		return
	}
	defs = incident.GetIncDefaults()
	err = incident.Read(r.FormValue("dir"), func(inc *incident.Incident) error {
		form.Set("dir", inc.Dir)
		form.Set("incname", inc.Config.IncidentName)
		form.Set("actnum", inc.Config.ActivationNum)
		if start := inc.Config.OpStart; !start.IsZero() {
			form.Set("opstartdate", start.Format("01/02/2006"))
			form.Set("opstarttime", start.Format("15:04"))
		}
		if end := inc.Config.OpEnd; !end.IsZero() {
			form.Set("openddate", end.Format("01/02/2006"))
			form.Set("opendtime", end.Format("15:04"))
		}
		form.Set("opcall", inc.Config.OpCall)
		form.Set("opname", inc.Config.OpName)
		form.Set("taccall", inc.Config.TacCall)
		form.Set("tacname", inc.Config.TacName)
		form.Set("rxmsgid", inc.Config.RxMessageID)
		form.Set("txmsgid", inc.Config.TxMessageID)
		form.Set("toaddr", inc.Config.DefaultTo)
		form.Set("topos", inc.Config.DefaultToPos)
		form.Set("toloc", inc.Config.DefaultToLoc)
		form.Set("frompos", inc.Config.DefaultFromPos)
		form.Set("fromloc", inc.Config.DefaultFromLoc)
		form.Set("defbody", inc.Config.DefaultBody)
		switch inc.Config.ConnectType {
		case incident.ConnectNone:
			form.Set("conntype", "manual")
			form.Set("bbscall", inc.Config.ConnectBBS)
			variables["nosaved"] = "true"
		case incident.ConnectSerialTNC:
			form.Set("conntype", "serial+tnc")
			form.Set("bbsaddr", inc.Config.ConnectAddress)
			variables["nosaved"] = "true"
		case incident.ConnectTelnet:
			form.Set("conntype", "telnet")
			form.Set("bbscall", inc.Config.ConnectBBS)
			if h, p, err := net.SplitHostPort(inc.Config.ConnectAddress); err == nil {
				form.Set("tcpaddr", h)
				form.Set("tcpport", p)
			}
			if inc.Config.TelnetPassword != "" {
				form.Set("usesaved", "checked")
			}
		default:
			return errors.New("invalid connection type in configuration")
		}
		if !inc.Config.NoSendReceipts {
			form.Set("sendrcpts", "checked")
		}
		form.Set("tnctype", inc.Config.TNCType)
		variables["tnctypes"] = strings.Join(tnc.AllTNCs(), ";")
		form.Set("serport", inc.Config.SerialPort)
		if ports := pseudomsg.GuessSerialPorts(); len(ports) != 0 {
			variables["serports"] = strings.Join(ports, ";")
		}
		if len(defs.TCPAddresses) != 0 {
			var triplets []string
			for bbs, addr := range defs.TCPAddresses {
				h, p, _ := net.SplitHostPort(addr)
				triplets = append(triplets, fmt.Sprintf("%s,%s,%s", bbs, h, p))
			}
			variables["tcpaddrs"] = strings.Join(triplets, ";")
		}
		if len(defs.TelnetPasswords) != 0 {
			variables["havesaved"] = strings.Join(slices.Collect(maps.Keys(defs.TelnetPasswords)), ",")
		}
		if len(inc.Config.BulletinChecks) != 0 {
			areas := slices.Collect(maps.Keys(inc.Config.BulletinChecks))
			slices.SortFunc(areas, func(a, b string) int {
				if strings.HasPrefix(a, "XSC") && strings.HasPrefix(b, "XSC") {
					return cmp.Compare(a, b)
				} else if strings.HasPrefix(a, "XSC") {
					return -1
				} else if strings.HasPrefix(b, "XSC") {
					return +1
				} else {
					return cmp.Compare(a, b)
				}
			})
			var triplets []string
			for _, area := range areas {
				dur := int(inc.Config.BulletinChecks[area].Duration / time.Minute)
				hr, min := dur/60, dur%60
				triplets = append(triplets, fmt.Sprintf("%s,%d,%d", area, hr, min))
			}
			variables["areas"] = strings.Join(triplets, ";")
		}
		form.Set("viewflags", inc.Config.ViewFlags.String())
		return nil
	})
	if err != nil {
		goto ERROR
	}
	if doc, err = htmlop.Parse(bytes.NewReader(configForm)); err != nil {
		goto ERROR
	}
	htmlop.Expand(doc, variables)
	htmlop.FillForm(doc, form)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html.Render(w, doc)
	return
ERROR:
	s.ErrPage(w, err.Error(), http.StatusInternalServerError)
}

func (s *Server) servePostIncidentConfig(w http.ResponseWriter, r *http.Request) {
	var err error

	s.outpost = false
	err = incident.Write(r.FormValue("dir"), func(inc *incident.Incident) error {
		var c = inc.Config.Clone()

		c.IncidentName = r.FormValue("incname")
		c.ActivationNum = r.FormValue("actnum")
		if r.FormValue("opstartdate") == "" {
			c.OpStart, c.OpEnd = time.Time{}, time.Time{}
		} else {
			c.OpStart, _ = time.ParseInLocation("01/02/2006 15:04", r.FormValue("opstartdate")+" "+r.FormValue("opstarttime"), time.Local)
			c.OpEnd, _ = time.ParseInLocation("01/02/2006 15:04", r.FormValue("openddate")+" "+r.FormValue("opendtime"), time.Local)
		}
		c.OpCall = r.FormValue("opcall")
		c.OpName = r.FormValue("opname")
		c.TacCall = r.FormValue("taccall")
		c.TacName = r.FormValue("tacname")
		c.RxMessageID = r.FormValue("rxmsgid")
		c.TxMessageID = r.FormValue("txmsgid")
		c.DefaultTo = r.FormValue("toaddr")
		c.DefaultToPos = r.FormValue("topos")
		c.DefaultToLoc = r.FormValue("toloc")
		c.DefaultFromPos = r.FormValue("frompos")
		c.DefaultFromLoc = r.FormValue("fromloc")
		c.DefaultBody = r.FormValue("defbody")
		switch r.FormValue("conntype") {
		case "manual":
			c.ConnectType = incident.ConnectNone
			c.ConnectBBS = r.FormValue("bbscall")
			c.NoSendReceipts = r.FormValue("sendrcpts") == ""
		case "serial+tnc":
			c.ConnectType = incident.ConnectSerialTNC
			c.NoSendReceipts = false
			c.TNCType = r.FormValue("tnctype")
			c.SerialPort = r.FormValue("serport")
			c.ConnectAddress = r.FormValue("bbsaddr")
			c.ConnectBBS, _, _ = strings.Cut(c.ConnectAddress, "-")
		case "telnet":
			c.ConnectType = incident.ConnectTelnet
			c.NoSendReceipts = false
			c.ConnectBBS = r.FormValue("bbscall")
			c.ConnectAddress = net.JoinHostPort(r.FormValue("tcpaddr"), r.FormValue("tcpport"))
			c.TelnetUser = c.ActiveCall()
			if r.FormValue("usesaved") != "" {
				c.TelnetPassword = incident.GetIncDefaults().TelnetPasswords[c.TelnetUser]
			} else {
				c.TelnetPassword = r.FormValue("bbspwd")
			}
		}
		c.BulletinChecks = make(map[string]incident.CheckFrequency)
		for key := range r.Form {
			if suffix, ok := strings.CutPrefix(key, "area"); ok {
				if area := r.FormValue(key); area != "" {
					hr, _ := strconv.Atoi(r.FormValue("hr" + suffix))
					min, _ := strconv.Atoi(r.FormValue("min" + suffix))
					now := r.FormValue("now"+suffix) != ""
					if hr != 0 || min != 0 || now {
						c.BulletinChecks[area] = incident.CheckFrequency{Duration: time.Duration(hr)*time.Hour + time.Duration(min)*time.Minute}
					}
					if now {
						inc.BulletinAreaChecked(area, time.Time{})
					}
				}
			}
		}
		c.ViewFlags, _ = incident.ParseViewFlags(r.FormValue("viewflags"))
		inc.UpdateConfig(c)
		inc.UpdateIncDefaults()
		return nil
	})
	if err != nil {
		goto ERROR
	}
	http.Redirect(w, r, strings.Replace(r.URL.String(), "-config", "", 1), http.StatusSeeOther)
	return
ERROR:
	s.ErrPage(w, err.Error(), http.StatusInternalServerError)
}
