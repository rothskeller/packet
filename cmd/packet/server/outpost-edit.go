package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/form"
	"github.com/rothskeller/packet/form/formdef"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/payload"
	"github.com/rothskeller/packet/message/subject"
)

const (
	opdirectTimeout = 30 * time.Second
	opdirectURL     = "http://127.0.0.1:9334/TBD" // yes, really
)

// outpostNewRequest handles a GET /outpost-new request, which is sent by the
// "packet outpost new" command.  The form parameters are:
//
//   - addon: name of the Outpost addon responsible for the form.  Required.
//   - msgtype: message type, also the name of the form HTML file.  Required.
//   - msgID: default origin message ID for new message.  Required.
//   - opName: operator's name.  Required.
//   - opCall: operator's FCC call sign.  Required.
//   - tacName: tactical station name.  Optional.
//   - tacCall: tactical station call sign.  Optional.
func (s *Server) outpostNewRequest(w http.ResponseWriter, r *http.Request) {
	var (
		def     *formdef.FormDef
		addon   string
		msgtype string
		fields  = make(map[string]string)
	)
	if maybeShowREADME(w, r) {
		return
	}
	// Find the definition of the form.
	addon, msgtype = r.FormValue("addon"), r.FormValue("msgtype")
	if mtype := message.FindType(func(mt message.MType) bool {
		if mt, ok := mt.(form.FormType); ok {
			return mt.AddonName == addon && mt.HTMLName == msgtype && len(mt.CreateTags) != 0
		}
		return false
	}); mtype == nil {
		slog.Error("no form definition", "addon", addon, "html", msgtype)
		ErrorPage(w, http.StatusInternalServerError, fmt.Errorf("no editable form definition found for %s/%s", addon, msgtype), nil)
		return
	} else {
		def = mtype.(form.FormType).FormDef
	}
	// Walk through the fields of the form, setting fields.
	for f := range def.AllFields() {
		// Default values for fields.
		if f.Tag != "" && f.Value != "" {
			fields[f.Tag] = f.Value // default value
		}
		// Well-known fields with supplied values.
		switch f.Common {
		case "messageDate", "formDate":
			fields[f.Tag] = time.Now().Format("01/02/2006")
		case "operatorCall":
			fields[f.Tag] = r.FormValue("opCall")
		case "operatorName":
			fields[f.Tag] = r.FormValue("opName")
		case "originMessageID":
			fields[f.Tag] = r.FormValue("msgID")
		case "tacticalCall":
			if v := r.FormValue("tacCall"); v != "" {
				fields[f.Tag] = v
			}
		case "tacticalName":
			if v := r.FormValue("tacName"); v != "" {
				fields[f.Tag] = v
			}
		}
	}
	s.editCommon(w, fields, def, "")
	slog.Info("rendered editor for new form", "addon", addon, "html", msgtype, "msgID", r.FormValue("msgID"))
}

// outpostEditRequest handles a GET /outpost-edit request, which is sent by the
// "packet outpost-edit" command.  The form parameters are:
//
//   - msgfile: name of file containing message.  Required.
//   - index:  Outpost index of the message.  Required.
func (s *Server) outpostEditRequest(w http.ResponseWriter, r *http.Request) {
	var (
		index  string
		msg    message.Message
		def    *formdef.FormDef
		body   *form.FormBody
		msgID  string
		err    error
		fields = make(map[string]string)
	)
	if maybeShowREADME(w, r) {
		return
	}
	// Get the Outpost message index.
	if _, err = strconv.Atoi(r.FormValue("index")); err != nil {
		slog.Error("missing/invalid index")
		ErrorPage(w, http.StatusBadRequest, fmt.Errorf("missing/invalid index parameter"), nil)
	} else {
		index = r.FormValue("index")
		delete(r.Form, "index")
	}
	// Open, read, and parse the message.
	if msg, err = message.ReadNoHeader(r.FormValue("msgfile")); msg == nil {
		slog.Error("read message from Outpost", "f", r.FormValue("msgfile"), "err", err)
		ErrorPage(w, http.StatusBadRequest, fmt.Errorf("reading message: %s", err), nil)
		return
	} else if mt, ok := msg.Type().(form.FormType); !ok {
		slog.Error("message from Outpost is not a form", "f", r.FormValue("msgfile"), "type", fmt.Sprintf("%T", msg.Type()))
		ErrorPage(w, http.StatusBadRequest, errors.New("message is not a recognized form"), nil)
		return
	} else {
		def = mt.FormDef
	}
	body = msg.Body().(*form.FormBody)
	for fd := range def.AllFields() {
		if fd.Tag != "" {
			fields[fd.Tag] = body.Field(fd.Tag)
		}
		if fd.Common == "originMessageID" {
			msgID = body.Field(fd.Tag)
		}
	}
	s.editCommon(w, fields, def, index)
	slog.Info("rendered editor for existing form", "addon", def.AddonName, "html", def.HTMLName, "msgID", msgID)
}

var eofRE = regexp.MustCompile(`(?i)%23EOF`)

// outpostSubmit handles POST /outpost-submit requests, which are submissions
// of forms edited through an outpost-new or outpost-edit request.
func (s *Server) outpostSubmit(w http.ResponseWriter, r *http.Request) {
	var (
		msg    message.Message
		addon  string
		msgID  string
		bbuild strings.Builder
		body   string
	)
	if msg = submitCommon(w, r); msg == nil {
		return // ErrorPage emitted
	}
	if ft, ok := msg.Type().(form.FormType); ok {
		addon = ft.AddonName
		// Fill in the OpDate and OpTime fields.
		for fd := range ft.AllFields() {
			switch fd.Common {
			case "operatorDate":
				msg.Body().(*form.FormBody).SetField(fd.Tag, time.Now().Format("01/02/2006"))
			case "operatorTime":
				msg.Body().(*form.FormBody).SetField(fd.Tag, time.Now().Format("15:04"))
			}
		}
	} else {
		// Must be the special case Check-In/Out.
		addon = "SCCoPIFO"
	}
	// The message body posted to Opdirect must be specially constructed
	// because Opdirect is sensitive to the order of elements in it.  It
	// also requires an & before the first parameter.
	fmt.Fprintf(&bbuild, "adn=%s", url.QueryEscape(addon))
	if index := r.FormValue("outpost-index"); index != "" {
		fmt.Fprintf(&bbuild, "&upd=%s", url.QueryEscape(index))
	}
	fmt.Fprintf(&bbuild, "&sub=%s", url.QueryEscape(msg.Subject().EncodedSubject()))
	if msg.Payload().(*payload.OutpostPayload).Urgent() {
		io.WriteString(&bbuild, "&urg=TRUE")
	} else {
		io.WriteString(&bbuild, "&urg=FALSE")
	}
	body = msg.Body().EncodedBody()
	body = strings.ReplaceAll(body, "\n", "\r\n")
	fmt.Fprintf(&bbuild, "&msg=%s", url.QueryEscape(body))
	body = bbuild.String()
	// Outpost treats #EOF as a PacFORM end of file marker, so we need to
	// make sure it doesn't appear anywhere in the data.
	body = eofRE.ReplaceAllStringFunc(body, func(s string) string {
		if strings.HasPrefix(s, "%23E") {
			return "%23%45" + s[4:]
		} else {
			return "%23%65" + s[4:]
		}
	})
	// Then we add the desired EOF marker.  Yes, all of this.
	body += "&4VAO=%0D%0A%23EOF"
	msgID = msg.Subject().(*subject.SCCoSubject).SubjectMessageID()
	if sendToOpdirect(w, r, body, msgID) {
		serveMessagePDF(w, r, msg)
	}
}

// sendToOpdirect sends the submitted message to Opdirect, and handles its
// various possible responses.  It returns true if successful (and nothing
// emitted); false if an error occurs (and an ErrorPage has been emitted).
func sendToOpdirect(w http.ResponseWriter, r *http.Request, body, msgID string) bool {
	var (
		ctx      context.Context
		cancel   func()
		req      *http.Request
		resp     *http.Response
		response string
		err      error
	)
	if runtime.GOOS != "windows" {
		// If we're running "packet outpost ..." on a non-Windows
		// system, it's presumably for testing/debugging.  Since there
		// can't be any Outpost to talk to, just pretend we did.
		return true
	}
	// Build and issue the request to Opdirect.
	ctx, cancel = context.WithTimeout(r.Context(), opdirectTimeout)
	defer cancel()
	req, _ = http.NewRequestWithContext(ctx, http.MethodPost, opdirectURL, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if resp, err = http.DefaultClient.Do(req); err != nil {
		slog.Error("post to opdirect", "url", opdirectURL, "err", err)
		if strings.Contains(err.Error(), "connection refused") {
			err = fmt.Errorf(" It appears that Opdirect is not running.\n%w", err)
		}
		ErrorPage(w, http.StatusInternalServerError, err, nil)
		return false
	}
	// Check the result from Opdirect.
	if by, _ := io.ReadAll(resp.Body); true {
		response = string(by)
	}
	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.Error("post to opdirect", "url", opdirectURL, "code", resp.StatusCode, "status", resp.Status)
		ErrorPage(w, http.StatusInternalServerError, fmt.Errorf(" Opdirect returned error %d %s",
			resp.StatusCode, resp.Status), response)
		return false
	} else if err = opdirectReturnCode(response); err != nil {
		slog.Error("post to opdirect", "url", opdirectURL, "err", err)
		ErrorPage(w, http.StatusInternalServerError, err, response)
		return false
	} else if strings.Contains(response, "Your PacFORMS submission was successful!") {
		slog.Error("PacFORMS response from opdirect")
		ErrorPage(w, http.StatusInternalServerError, errors.New(" It appears you are running an obsolete version of Outpost. "), response)
		return false
	}
	slog.Info("message accepted by opdirect", "msgID", msgID)
	return true
}

var odrcRE = regexp.MustCompile(`(?i)<\s*meta\s+[^>]*\bname\s*=\s*"OpDirectReturnCode"[^>]*`)
var odrcContentRE = regexp.MustCompile(`\s+content\s*=\s*"\s*([0-9]+)\s*"`)

// opdirectReturnCode looks for an Opdirect return code embedded in the response
// and returns an appropriate error.
func opdirectReturnCode(response string) (err error) {
	if match := odrcRE.FindString(response); match != "" {
		if match2 := odrcContentRE.FindStringSubmatch(match); match2 != nil {
			status, _ := strconv.Atoi(match2[1])
			if status < 200 || status >= 300 {
				return fmt.Errorf(" Opdirect returned error %d", status)
			}
			return nil
		}
	}
	return errors.New(" The response from Opdirect did not contain a status code. ")
}

// serveMessagePDF generates a PDF with the message contents and serves it as
// the response to the current request.  If something goes wrong it may serve an
// error page instead.
func serveMessagePDF(w http.ResponseWriter, r *http.Request, msg message.Message) {
	var (
		fname string
		err   error
	)
	// Create a temp file for the PDF.
	if fname, err = CreateTempPDF(msg); err != nil {
		ErrorPage(w, http.StatusInternalServerError, fmt.Errorf("unable to create PDF: %w", err), nil)
		return
	}
	// Send a redirect to fetch that file.  We can't serve the PDF directly
	// because the browser will want to re-fetch it to handle save, print,
	// or reload, so it has to have its own URL.
	http.Redirect(w, r, "/pdf/"+filepath.Base(fname), http.StatusSeeOther)
	slog.Debug("redirecting to PDF", "url", "/pdf"+filepath.Base(fname))
}

// CreateTempPDF creates a PDF rendering of the supplied message in a temporary
// file, and returns the filename.
func CreateTempPDF(msg message.Message) (fname string, err error) {
	var (
		def    *formdef.FormDef
		msgid  string
		tmpdir = os.TempDir()
	)
	// Get the origin message number.
	if ft, ok := msg.Type().(form.FormType); ok {
		def = ft.FormDef
		for fd := range def.AllFields() {
			if fd.Common == "originMessageID" {
				msgid = msg.Body().(*form.FormBody).Field(fd.Tag)
				break
			}
		}
	}
	if msgid == "" {
		// No message number found, so create one with an UNK- prefix.
		for seq := 0; ; seq++ {
			fname = filepath.Join(tmpdir, fmt.Sprintf("UNK-%03dP.pdf", seq))
			if _, err := os.Stat(fname); os.IsNotExist(err) {
				break
			} else if err != nil {
				return "", err
			}
		}
	} else {
		// Use the message number as the filename, adding a sequence
		// number if needed.
		fname = filepath.Join(tmpdir, msgid+".pdf")
		for seq := 2; ; seq++ {
			if _, err := os.Stat(fname); os.IsNotExist(err) {
				break
			} else if err != nil {
				return "", err
			}
			fname = filepath.Join(tmpdir, fmt.Sprintf("%s (%d).pdf", msgid, seq))
		}
	}
	err = msg.Type().RenderPDF(msg, fname, "")
	if _, ok := err.(message.Warning); ok {
		slog.Warn("RenderPDF", "w", err)
		err = nil // ignore warnings
	} else if err != nil {
		slog.Error("RenderPDF", "f", fname, "err", err)
		return "", fmt.Errorf("unable to create PDF: %s", err)
	}
	slog.Debug("created temp PDF", "f", fname)
	return fname, nil
}
