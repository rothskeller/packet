package server

import (
	"bytes"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/rothskeller/packet/form"
	"github.com/rothskeller/packet/form/formdef"
	"github.com/rothskeller/packet/form/formdefs"
	"github.com/rothskeller/packet/form/htmlop"
	"github.com/rothskeller/packet/form/pifover"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/body"
	"github.com/rothskeller/packet/message/payload"
	"github.com/rothskeller/packet/message/subject"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// serveEditAsset handles a GET /assets/{asset...} request.
func (s *Server) serveAsset(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/" + r.PathValue("asset")
	w.Header().Set("Cache-Control", "no-store")
	http.FileServerFS(formdefs.FormsFS).ServeHTTP(w, r)
}

// serveRenderedPDF handles a GET /pdf/{filename} request.
func (s *Server) serveRenderedPDF(w http.ResponseWriter, r *http.Request) {
	fname := r.PathValue("filename")
	if !strings.HasSuffix(fname, ".pdf") || strings.HasPrefix(fname, ".") || strings.ContainsAny(fname, `/\`) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	r.URL.Path = "/" + fname
	w.Header().Set("Cache-Control", "max-age=3600")
	http.FileServerFS(os.DirFS(os.TempDir())).ServeHTTP(w, r)
}

// editCommon is the common parts of starting to edit a message.
func (s *Server) editCommon(w http.ResponseWriter, fields map[string]string, def *formdef.FormDef, index string) {
	var (
		formFile []byte
		formHTML *html.Node
		bundle   string
		defFile  string
		formBuf  bytes.Buffer
		err      error
	)
	// Read and parse the HTML for the form.
	if formFile, err = fs.ReadFile(formdefs.FormsFS, def.HTMLFile); err != nil {
		slog.Error("fs.ReadFile", "f", def.HTMLFile, "err", err)
		ErrorPage(w, http.StatusInternalServerError, err, nil)
		return
	}
	if formHTML, err = html.Parse(bytes.NewReader(formFile)); err != nil {
		slog.Error("html.Parse", "f", def.HTMLFile, "err", err)
		ErrorPage(w, http.StatusInternalServerError, fmt.Errorf("%s: %s", def.HTMLFile, err), nil)
		return
	}
	bundle, _, _ = strings.Cut(def.HTMLFile, "/")
	// If there is a definitions.html in the same directory, read and parse
	// it too, and prepend it to the form HTML.
	defFile = bundle + "/definitions.html"
	if formFile, err = fs.ReadFile(formdefs.FormsFS, defFile); err == nil {
		var defHTML *html.Node
		if defHTML, err = html.Parse(bytes.NewReader(formFile)); err != nil {
			slog.Error("html.Parse", "f", defFile, "err", err)
			ErrorPage(w, http.StatusInternalServerError, fmt.Errorf("%s: %s", defFile, err), nil)
			return
		}
		formBody := findBody(formHTML)
		defBody := findBody(defHTML)
		for c := defBody.LastChild; c != nil; c = defBody.LastChild {
			defBody.RemoveChild(c)
			formBody.InsertBefore(c, formBody.FirstChild)
		}
	}
	fields["assets"] = "/assets/" + bundle
	fields["submit-label"] = "Submit to Outpost"
	fields["submit-url"] = "/outpost-submit"
	fields["addon-name"] = def.AddonName
	fields["form-html"] = def.HTMLName
	fields["form-version"] = def.Version
	if index != "" {
		fields["outpost-index"] = index
	}
	if def.PDFFile != "" {
		fields["pdf-url"] = "/assets/" + def.PDFFile
	}
	// Expand the templates in the form HTML, using the supplied fields.
	htmlop.Expand(formHTML, fields)
	// Fill in the form using the supplied fields.
	htmlop.FillForm(formHTML, fieldsToValues(fields))
	// Render and minimize the result.
	htmlop.Minify(&formBuf, formHTML)
	// Reply with the form.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Write(formBuf.Bytes())
}

// submitCommon reads a submitted form from a request and builds the
// corresponding message.  If it returns nil, it rendered an ErrorPage.
func submitCommon(w http.ResponseWriter, r *http.Request) message.Message {
	var (
		addonName string
		htmlName  string
		version   string
		mtype     message.MType
		def       *formdef.FormDef
		fbody     *form.FormBody
		msgID     string
		handling  string
		summary   string
		pload     *payload.OutpostPayload
		subj      *subject.SCCoSubject
		msg       *message.DraftMessage
		err       error
	)
	addonName, htmlName, version = r.FormValue("addon-name"), r.FormValue("form-html"), r.FormValue("form-version")
	if mtype = message.FindType(func(mt message.MType) bool {
		if mt, ok := mt.(form.FormType); ok {
			return mt.AddonName == addonName && mt.HTMLName == htmlName && mt.Version == version && len(mt.CreateTags) != 0
		}
		return false
	}); mtype == nil {
		slog.Error("form not found", "addon", addonName, "html", htmlName, "ver", version)
		ErrorPage(w, http.StatusInternalServerError, fmt.Errorf("no such form %s/%s/%s", addonName, htmlName, version), nil)
		return nil
	} else {
		def = mtype.(form.FormType).FormDef
	}
	if def.AddonName == "SCCoPIFO" && def.HTMLName == "form-checkin-out.html" {
		return submitCheckInOut(w, r) // special case
	}
	if fbody, err = form.NewFormBody(addonName, htmlName, pifover.PIFOVersion, version); err != nil {
		slog.Error("form.NewFormBody", "addon", addonName, "html", htmlName, "ver", version)
		ErrorPage(w, http.StatusInternalServerError, fmt.Errorf("%s/%s: %s", addonName, htmlName, err), nil)
		return nil
	}
	for f := range def.AllFields() {
		if f.Tag == "" {
			continue
		}
		value := r.FormValue(f.Tag)
		fbody.SetField(f.Tag, r.FormValue(f.Tag))
		switch f.Common {
		case "handling":
			handling = value
		case "messageSummary":
			summary = value
		case "originMessageID":
			msgID = value
		}
	}
	pload = payload.NewOutpostPayload(fbody)
	if strings.HasPrefix(handling, "I") {
		pload.SetUrgent(true)
	}
	subj, _ = subject.NewSCCoSubject(msgID, handling, def.SubjectTag+"_"+summary)
	msg = message.NewDraftMessage(mtype, subj, pload, false)
	return msg
}

func fieldsToValues(fields map[string]string) (values url.Values) {
	values = make(url.Values)
	for k, v := range fields {
		values.Set(k, v)
	}
	return values
}

func findBody(doc *html.Node) *html.Node {
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.Body {
			return n
		}
	}
	return nil
}

// submitCheckInOut is a special case handler for the check-in/out form, which
// is created like a form but then rendered in plain text.  If we ever have
// more cases like this, it may be worth implementing a mechanism for them to
// be data-driven, but as long as there's only one, it's not worth bothering.
func submitCheckInOut(_ http.ResponseWriter, r *http.Request) message.Message {
	var (
		msgno   = r.FormValue("MsgNo")
		inOut   = r.FormValue("in_out")
		opCall  = r.FormValue("operator_call_sign")
		opName  = r.FormValue("operator_name")
		tacCall = r.FormValue("tactical_call_sign")
		tacName = r.FormValue("tactical_name")
		subj    subject.Subject
		bod     body.Body
		pload   payload.Payload
		mtype   cicoMType
	)
	if tacCall != "" {
		subj, _ = subject.NewSCCoSubject(msgno, "ROUTINE", fmt.Sprintf("Check-%s %s, %s", inOut, tacCall, tacName))
		bod = body.NewPlainBody(fmt.Sprintf("Check-%s %s, %s\n%s, %s\n", inOut, tacCall, tacName, opCall, opName))
	} else {
		subj, _ = subject.NewSCCoSubject(msgno, "ROUTINE", fmt.Sprintf("Check-%s %s, %s", inOut, opCall, opName))
		bod = body.NewPlainBody(fmt.Sprintf("Check-%s %s, %s\n", inOut, opCall, opName))
	}
	pload = payload.NewOutpostPayload(bod)
	mtype.BaseMType = message.NewBaseMType(fmt.Sprintf("a check-%s message", strings.ToLower(inOut)), "")
	return message.NewDraftMessage(mtype, subj, pload, false)
}

type cicoMType struct{ *message.BaseMType }

func (cicoMType) Recognize(message.Message) {}
