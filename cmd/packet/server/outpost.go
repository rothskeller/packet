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
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/payload"
)

const (
	opdirectTimeout = 30 * time.Second
	opdirectURL     = "http://127.0.0.1:9334/TBD" // yes, really
)

// outpostNewRequest handles a GET /outpost-new request, which is sent by the
// "packet outpost new" command.  The form parameters are:
//
//   - formtag: create tag of the form to create.  Required.
//   - msgID: default origin message ID for new message.  Required.
//   - opName: operator's name.  Required.
//   - opCall: operator's FCC call sign.  Required.
//   - tacName: tactical station name.  Optional.
//   - tacCall: tactical station call sign.  Optional.
func (s *Server) outpostNewRequest(w http.ResponseWriter, r *http.Request) {
	var (
		formtag string
		mtype   message.EditableMType
		msg     message.Message
	)
	if maybeShowREADME(w, r) {
		return
	}
	// Find the definition of the form.
	formtag = r.FormValue("formtag")
	if mtype = message.FindCreateTag(formtag); mtype == nil {
		slog.Error("no form definition", "formtag", formtag)
		ErrPage(w, fmt.Sprintf("No editable form definition was found for %q.  Please report this to the author.", formtag), http.StatusInternalServerError)
		return
	}
	msg = mtype.NewDraft()
	// Walk through the fields of the form, setting fields.
	for f := range msg.Fields() {
		// Well-known fields with supplied values.
		switch f.Common() {
		case field.CMessageDate, field.CFormDate:
			f.SetValue(msg, time.Now().Format("01/02/2006"))
		case field.COperatorCall:
			f.SetValue(msg, r.FormValue("opCall"))
		case field.COperatorName:
			f.SetValue(msg, r.FormValue("opName"))
		case field.COriginMessageID:
			f.SetValue(msg, r.FormValue("msgID"))
		case field.CTacticalCall:
			f.SetValue(msg, r.FormValue("tacCall"))
		case field.CTacticalName:
			f.SetValue(msg, r.FormValue("tacName"))
		case field.CUseTactical:
			if r.FormValue("tacCall") != "" {
				f.SetValue(msg, "checked")
			}
		case field.CReceiverSender:
			f.SetValue(msg, "sender")
		case field.COperatorMethod:
			f.SetValue(msg, "Other")
		case field.COperatorMethodOther:
			f.SetValue(msg, "Packet")
		}
	}
	s.outpostEditCommon(w, msg, "")
	slog.Info("rendered editor for new form", "formtag", formtag, "msgID", r.FormValue("msgID"))
}

// outpostEditRequest handles a GET /outpost-edit request, which is sent by the
// "packet outpost-edit" command.  The form parameters are:
//
//   - msgfile: name of file containing message.  Required.
//   - index:  Outpost index of the message.  Required.
func (s *Server) outpostEditRequest(w http.ResponseWriter, r *http.Request) {
	var (
		index string
		msg   message.Message
		msgID string
		err   error
	)
	if maybeShowREADME(w, r) {
		return
	}
	// Get the Outpost message index.
	if _, err = strconv.Atoi(r.FormValue("index")); err != nil {
		slog.Error("missing/invalid index")
		ErrPage(w, "The GET /outpost-edit request is missing the required index parameter.  Please report this error to the author.", http.StatusInternalServerError)
	} else {
		index = r.FormValue("index")
		delete(r.Form, "index")
	}
	// Open, read, and parse the message.
	if msg, err = message.ReadNoHeader(r.FormValue("msgfile")); msg == nil {
		slog.Error("read message from Outpost", "f", r.FormValue("msgfile"), "err", err)
		ErrPage(w, "The message provided by Outpost was not in a valid format.  Please report this error to the author.", http.StatusInternalServerError)
		return
	} else if _, ok := msg.Type().(form.EditableFormType); !ok {
		slog.Error("message from Outpost is not a form", "f", r.FormValue("msgfile"), "type", fmt.Sprintf("%T", msg.Type()))
		ErrPage(w, "The message provided by Outpost was not in a valid format.  Please report this error to the author.", http.StatusInternalServerError)
		return
	}
	for f := range msg.Fields() {
		if f.Common() == field.COriginMessageID {
			msgID = f.Value(msg)
		}
	}
	s.outpostEditCommon(w, msg, index)
	slog.Info("rendered editor for existing form", "msgID", msgID)
}

// outpostEditCommon is the common parts of Outpost new and edit message.
func (s *Server) outpostEditCommon(w http.ResponseWriter, msg message.Message, index string) {
	var (
		tag    string
		vars   message.EditHTMLVars
		out    []byte
		err    error
		params = make(url.Values)
	)
	tag = msg.Type().(message.EditableMType).CreateTag()
	params.Set("formtag", tag)
	if index != "" {
		params.Set("outpost-index", index)
	}
	vars.AssetBase = "/assets/" + url.PathEscape(tag)
	vars.SubmitLabel = "Submit to Outpost"
	vars.SubmitURL = "/outpost-submit?" + params.Encode()
	if out, err = msg.Type().(message.EditableMType).EditHTML(msg.(*message.DraftMessage), vars); err != nil {
		ErrPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(out)
}

var eofRE = regexp.MustCompile(`(?i)%23EOF`)

// outpostSubmit handles POST /outpost-submit requests, which are submissions
// of forms edited through an outpost-new or outpost-edit request.
func (s *Server) outpostSubmit(w http.ResponseWriter, r *http.Request) {
	var (
		msg     message.Message
		addon   string
		msgID   string
		bbuild  strings.Builder
		body    string
		formtag string
		mtype   message.EditableMType
		err     error
	)
	formtag = r.FormValue("formtag")
	if mtype = message.FindCreateTag(formtag); mtype == nil {
		slog.Error("form not found", "formtag", formtag)
		ErrPage(w, fmt.Sprintf("The form with tag=%q was not found.  Please report this error to the author.", formtag), http.StatusInternalServerError)
		return
	}
	if msg, err = mtype.FromPOST(r); err != nil {
		ErrPage(w, fmt.Sprintf("The form could not be read (%s).  Please report this error to the author.", err), http.StatusInternalServerError)
		return
	}
	addon = mtype.(form.EditableFormType).AddonName
	// Fill in the OpDate and OpTime fields.
	for fd := range msg.Fields() {
		switch fd.Common() {
		case field.COperatorDate:
			fd.SetValue(msg, time.Now().Format("01/02/2006"))
		case field.COperatorTime:
			fd.SetValue(msg, time.Now().Format("15:04"))
		}
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
	msgID = msg.Subject().SubjectMessageID()
	if sendToOpdirect(w, r, body, msgID) {
		serveMessagePDF(w, r, msg)
	}
}

// sendToOpdirect sends the submitted message to Opdirect, and handles its
// various possible responses.  It returns true if successful (and nothing
// emitted); false if an error occurs (and an ErrPage has been emitted).
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
		ErrPage(w, "The packet software was unable to communicate with Outpost.  Are Outpost and Opdirect running?", http.StatusBadRequest)
		return false
	}
	// Check the result from Opdirect.
	if by, _ := io.ReadAll(resp.Body); true {
		response = string(by)
	}
	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.Error("post to opdirect", "url", opdirectURL, "code", resp.StatusCode, "status", resp.Status)
		ErrPage(w, "Outpost returned an error and did not accept the message.  Please report this problem to the author.", http.StatusInternalServerError)
		return false
	} else if err = opdirectReturnCode(response); err != nil {
		slog.Error("post to opdirect", "url", opdirectURL, "err", err)
		ErrPage(w, "Outpost returned an error and did not accept the message.  Please report this problem to the author.", http.StatusInternalServerError)
		return false
	} else if strings.Contains(response, "Your PacFORMS submission was successful!") {
		slog.Error("PacFORMS response from opdirect")
		ErrPage(w, "It appears you are running an obsolete version of Outpost.  Please upgrade Outpost to a current version.", http.StatusInternalServerError)
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
		ErrPage(w, "The software was unable to create a PDF file for this message.  Please report this error to the author.", http.StatusInternalServerError)
		return
	}
	// Tell the client to redirect to fetch that file.  We can't serve the
	// PDF directly because the browser will want to re-fetch it to handle
	// save, print, or reload, so it has to have its own URL.
	w.Header().Set("X-Packet-Action", "redirect:/pdf/"+filepath.Base(fname))
	w.WriteHeader(http.StatusNoContent)
	slog.Debug("redirecting to PDF", "url", "/pdf/"+filepath.Base(fname))
}

// CreateTempPDF creates a PDF rendering of the supplied message in a temporary
// file, and returns the filename.
func CreateTempPDF(msg message.Message) (fname string, err error) {
	var (
		msgid  string
		tmpdir = os.TempDir()
	)
	// Get the origin message number.
	for fd := range msg.Fields() {
		if fd.Common() == "originMessageID" {
			msgid = fd.Value(msg)
			break
		}
	}
	if msgid == "" {
		// No message number field found.  See if we can extract one
		// from the subject line.
		msgid = msg.Subject().SubjectMessageID()
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
