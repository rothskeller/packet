package server

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/address"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/messageid"
)

// serveGetNewMessage handles GET /new-message requests, which are sent from
// the Create New Message dialog in the GUI.  They have dir= and tag= params,
// with tag= specifying the create tag of the message type to create.  This
// semantically should be a POST, because it changes state, but window.open()
// issues a GET, so that's what it is.  This function creates a new draft
// message of the correct type in the incident, saves it, and then redirects to
// an edit URL to edit it.  That way refreshes of the URL don't create multiple
// messages.
func (s *Server) serveGetNewMessage(w http.ResponseWriter, r *http.Request) {
	var (
		dir    string
		tag    string
		mt     message.EditableMType
		msg    *message.DraftMessage
		ident  int
		params url.Values
		err    error
	)
	dir, tag = r.FormValue("dir"), r.FormValue("tag")
	if mt = message.FindTypeTag(tag); mt == nil {
		slog.Error("no such message tag", "tag", tag)
		ErrPage(w, fmt.Sprintf("The message type tag %q is not recognized.  Please report this error to the author.", tag), http.StatusInternalServerError)
		return
	}
	msg = mt.NewDraft().(*message.DraftMessage)
	err = incident.Write(dir, func(i *incident.Incident) (err error) {
		if ident, err = i.AddDraftMessage(msg, true); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		ErrPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	params = make(url.Values)
	params.Set("dir", dir)
	params.Set("id", strconv.Itoa(ident))
	http.Redirect(w, r, "/edit-message?"+params.Encode(), http.StatusSeeOther)
}

// serveGetEditMessage handles GET /edit-message requests.  They have dir= and
// id= params, and open an editor window on the message identified by id (which
// must be a draft).
func (s *Server) serveGetEditMessage(w http.ResponseWriter, r *http.Request) {
	var (
		dir    string
		ident  int
		msg    message.Message
		tag    string
		vars   message.EditHTMLVars
		out    []byte
		err    error
		params = url.Values{}
	)
	dir = r.FormValue("dir")
	ident, _ = strconv.Atoi(r.FormValue("id"))
	params.Set("dir", dir)
	params.Set("id", strconv.Itoa(ident))
	if err = incident.Read(dir, func(i *incident.Incident) (err error) {
		vars.FromAddress = i.Config.FromAddress()
		if le := i.GetLogEntryByIdent(ident); le == nil {
			slog.Error("no such message", "dir", dir, "id", ident)
			return fmt.Errorf("no such message %d", ident)
		} else if msg, err = i.GetMessageFromLogEntry(le); msg == nil && err != nil {
			return err
		} else if msg == nil {
			slog.Error("no such message", "dir", dir, "id", ident)
			return fmt.Errorf(" The message with ID %d was not found.", ident)
		} else {
			return nil
		}
	}); err != nil {
		ErrPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, ok := msg.(*message.DraftMessage); !ok {
		slog.Error("message not editable", "dir", dir, "id", ident, "t", fmt.Sprintf("%T", msg))
		ErrPage(w, fmt.Sprintf("Message %d is not editable: it is not an unsent outgoing message.", ident), http.StatusBadRequest)
		return
	}
	tag = msg.Type().(message.EditableMType).CreateTag()
	params.Set("tag", tag)
	vars.AssetBase = "/assets/" + url.PathEscape(tag)
	vars.SaveLabel = "Save as Draft"
	vars.ShowAddressFields = true
	vars.SubmitLabel = "Send Message"
	vars.SubmitURL = "/send-message?" + params.Encode()
	if out, err = msg.Type().(message.EditableMType).EditHTML(msg.(*message.DraftMessage), vars); err != nil {
		ErrPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, private")
	w.Write(out)
}

// serveGetAsset handles GET /assets/{tag}/{path...} requests, for assets used
// by the edit HTML.
func (s *Server) serveGetAsset(w http.ResponseWriter, r *http.Request) {
	var (
		tag string
		mt  message.EditableMType
		afs fs.FS
	)
	tag = r.PathValue("tag")
	if mt = message.FindTypeTag(tag); mt == nil {
		slog.Error("no such message type", "tag", tag)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if afs = mt.EditAssets(); afs == nil {
		slog.Error("no assets for message type", "tag", tag)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	r.URL.Path = "/" + r.PathValue("asset")
	w.Header().Set("Cache-Control", "no-store, private")
	http.FileServerFS(afs).ServeHTTP(w, r)
}

// servePostSendMessage handles POST /send-message requests, to send a message
// or save it as a draft.  In addition to the form fields, the requests have
// dir, tag, id, and readyToSend parameters, which are the incident directory,
// message type tag, log entry ident, and ready-to-send flags, respectively.
// The request returns either an error status with a text/plain body containing
// an error message or a 204 No Content indicating that the message was saved.
// In the latter case, it may also send an X-Packet-Action header set to
// "manual-send", indicating that the manual send dialog should be invoked for
// the message.
func (s *Server) servePostSendMessage(w http.ResponseWriter, r *http.Request) {
	var (
		dir    string
		tag    string
		ident  int
		mt     message.EditableMType
		msg    *message.DraftMessage
		manual bool
		err    error
	)
	dir, tag = r.FormValue("dir"), r.FormValue("tag")
	ident, _ = strconv.Atoi(r.FormValue("id"))
	if mt = message.FindTypeTag(tag); mt == nil {
		slog.Error("no such message type", "tag", tag)
		http.Error(w, fmt.Sprintf("no such message type %q", tag), http.StatusBadRequest)
		return
	}
	if m, err := mt.FromPOST(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		msg = m.(*message.DraftMessage)
	}
	if r.FormValue("readyToSend") == "true" {
		// It's not ready to send unless we have a valid To: address.
		if addrs, err := address.ParseList(msg.To()); err != nil {
			http.Error(w, `The "To Address" field contains an invalid packet address.`, http.StatusBadRequest)
			return
		} else if len(addrs) == 0 {
			http.Error(w, `The "To Address" field is required.`, http.StatusBadRequest)
			return
		} else if idx := slices.IndexFunc(addrs, addressIsBulletin); idx >= 0 {
			if len(addrs) > 1 {
				http.Error(w, fmt.Sprintf(`The "To Address" field contains %q, which is a bulletin address.  A bulletin can only be addressed to a single address.`, addrs[idx].Address), http.StatusBadRequest)
				return
			}
			msg.SetBulletin(true)
		}
		msg.SetReadyToSend(true)
	}
	// Update the incident and save the message.
	if err = incident.Write(dir, func(i *incident.Incident) (err error) {
		manual = i.Config.ConnectType == incident.ConnectNone
		return i.UpdateDraftMessage(ident, msg)
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if manual && msg.ReadyToSend() {
		w.Header().Set("X-Packet-Action", fmt.Sprintf("manual-send:%d", ident))
	} else {
		w.Header().Set("X-Packet-Action", fmt.Sprintf("select:%d", ident))
	}
	w.WriteHeader(http.StatusNoContent)
}

func addressIsBulletin(addr *address.Address) bool {
	user, domain, _ := strings.Cut(strings.ToLower(addr.Address), "@")
	switch user {
	case "xscperm", "xscevent", "xsctest":
		return true
	}
	domain = strings.TrimPrefix(domain, "all")
	switch domain {
	case "local", "xsc", "bay", "nca", "ca", "wusa", "usa", "noam", "ww":
		return true
	}
	return false
}

func (s *Server) serveGetViewMessage(w http.ResponseWriter, r *http.Request) {
	var (
		dir   string
		ident int
		msg   message.Message
		pdf   string
		dpdf  string
		err   error
	)
	dir = r.FormValue("dir")
	ident, _ = strconv.Atoi(r.FormValue("id"))
	err = incident.Read(dir, func(i *incident.Incident) (err error) {
		if le := i.GetLogEntryByIdent(ident); le == nil {
			slog.Error("no such message", "dir", dir, "id", ident)
			return fmt.Errorf(" The message with ID %d was not found.", ident)
		} else if msg, err = i.GetMessageFromLogEntry(le); msg == nil && err != nil {
			return err
		} else if msg == nil {
			slog.Error("no such message", "dir", dir, "id", ident)
			return fmt.Errorf(" The message with ID %d was not found.", ident)
		} else {
			pdf = incident.ToPDF(le.Filename())
			return nil
		}
	})
	if err != nil {
		ErrPage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// It's possible that the PDF doesn't exist yet.  If so we need to
	// create it.
	dpdf = filepath.Join(dir, pdf)
	if _, err := os.Stat(dpdf); os.IsNotExist(err) {
		if err = msg.Type().RenderPDF(msg, dpdf, ""); err != nil {
			slog.Error("RenderPDF", "dir", dir, "id", ident, "err", err)
			ErrPage(w, fmt.Sprintf("Unable to create PDF: %s", err), http.StatusInternalServerError)
			return
		}

	} else if err != nil {
		ErrPage(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename="+pdf)
	http.ServeFile(w, r, dpdf)
}

// servePostPrintMessage handles POST /print-message requests, asking for a
// message to be sent to the default printer.  These should be issued only if
// serverPrintCmd is non-empty, which basically means anywhere except Windows.
// These will have dir= and id= parameters.  They respond with an error status
// and text/plain error message, or with 204 No Content on success.
func (s *Server) servePostPrintMessage(w http.ResponseWriter, r *http.Request) {
	var (
		dir   string
		ident int
		msg   message.Message
		pdf   string
		dpdf  string
		cmd   *exec.Cmd
		err   error
	)
	serverPrintOnce.Do(setServerCanPrint)
	if serverPrintCmd == "" {
		slog.Error("server print not supported")
		http.Error(w, "Server-side printing is not supported on this system.", http.StatusBadRequest)
		return
	}
	dir = r.FormValue("dir")
	ident, _ = strconv.Atoi(r.FormValue("id"))
	err = incident.Read(dir, func(i *incident.Incident) (err error) {
		if le := i.GetLogEntryByIdent(ident); le == nil {
			slog.Error("no such message", "dir", dir, "id", ident)
			return fmt.Errorf(" The message with ID %d was not found.", ident)
		} else if msg, err = i.GetMessageFromLogEntry(le); msg == nil && err != nil {
			return err
		} else if msg == nil {
			slog.Error("no such message", "dir", dir, "id", ident)
			return fmt.Errorf(" The message with ID %d was not found.", ident)
		} else {
			pdf = incident.ToPDF(le.Filename())
			return nil
		}
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// It's possible that the PDF doesn't exist yet.  If so we need to
	// create it.
	dpdf = filepath.Join(dir, pdf)
	if _, err := os.Stat(dpdf); os.IsNotExist(err) {
		if err = msg.Type().RenderPDF(msg, dpdf, ""); err != nil {
			slog.Error("RenderPDF", "dir", dir, "id", ident, "err", err)
			http.Error(w, fmt.Sprintf("Unable to create PDF: %s", err), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	cmd = exec.Command(serverPrintCmd, dpdf)
	if err = cmd.Start(); err == nil {
		go cmd.Wait()
	}
	w.WriteHeader(http.StatusNoContent)
}

// servePostNewMessageFrom handles POST /new-message-from requests, which
// create new draft messages based on the content of existing messages.  They
// have dir=, id=, and action= parameters, which give the incident directory,
// the log entry ID for the source message, and the way the draft message should
// be generated from the source ("resend", "copy", or "reply").  They respond
// with an error status and text/plain error message body, or a 204 No Content
// with an X-Packet-Action header to edit the resulting draft.
func (s *Server) servePostNewMessageFrom(w http.ResponseWriter, r *http.Request) {
	var action = r.FormValue("action")

	serveMessage(w, r, true, func(i *incident.Incident, le *incident.LogEntry, msg message.Message) error {
		var dr *message.DraftMessage

		if emt, ok := msg.Type().(message.EditableMType); !ok {
			_, typename, _ := strings.Cut(msg.Type().Name(), " ")
			return errors.NewF("Creation of new %ss is not supported.", typename)
		} else {
			dr = emt.NewDraft().(*message.DraftMessage)
		}
		switch action {
		case "resend":
			if _, ok := msg.(*message.SentMessage); !ok {
				return errors.New(`The "Resend" action can only be used with a sent message.`)
			}
			for f := range msg.Fields() {
				switch f.Common() {
				case field.CDestinationMessageID: // ignore
				case field.COriginMessageID:
					f.SetValue(dr, incrementMsgIDSuffix(i, f.Value(msg)))
				default:
					f.SetValue(dr, f.Value(msg))
				}
			}
			dr.SetTo(msg.To())
		case "copy":
			switch msg.(type) {
			case *message.ReceivedMessage:
				// OK
			case *message.SentMessage:
				dr.SetTo(msg.To())
			default:
				return errors.New(`The "Send Copy" action can only be used with a received or sent message.`)
			}
			for f := range msg.Fields() {
				switch f.Common() {
				case field.CDestinationMessageID, field.COriginMessageID: // ignore
				default:
					f.SetValue(dr, f.Value(msg))
				}
			}
		case "reply":
			var fref, ftopos, ftoloc field.Field

			if _, ok := msg.(*message.ReceivedMessage); !ok {
				return errors.New(`The "Reply" action can only be used with a received message.`)
			}
			for f := range msg.Fields() {
				switch f.Common() {
				case field.CReference:
					fref = f
				case field.CToICSPosition:
					ftopos = f
				case field.CToLocation:
					ftoloc = f
				case field.CHandling, field.CMessageSummary:
					f.SetValue(dr, f.Value(msg))
				}
			}
			for f := range msg.Fields() {
				switch f.Common() {
				case field.COriginMessageID:
					if fref != nil {
						fref.SetValue(dr, f.Value(msg))
					}
				case field.CFromICSPosition:
					if ftopos != nil {
						ftopos.SetValue(dr, f.Value(msg))
					}
				case field.CFromLocation:
					if ftoloc != nil {
						ftoloc.SetValue(dr, f.Value(msg))
					}
				}
			}
		// TODO: reply fields
		default:
			return errors.NewF("%q is not a recognized action for the /new-message-from request.", action)
		}
		if ident, err := i.AddDraftMessage(dr, false); err != nil {
			return err
		} else {
			w.Header().Set("X-Packet-Action", "edit:"+strconv.Itoa(ident))
		}
		return nil
	}, nil)
}

// copyFields copies the values of all fields from the from message to the to
// message.  The two messages are assumed to be of the same type.
func copyFields(to *message.DraftMessage, from message.Message) {
	for f := range from.Fields() {
		f.SetValue(to, f.Value(from))
	}
}

// incrementMsgIDSuffix increments the suffix of the origin message ID of the
// given message.  If R is an available suffix, it is used; failing that, S, T,
// etc. until an available suffix is found.
func incrementMsgIDSuffix(i *incident.Incident, msgid string) string {
	pfx, seq, suf, err := messageid.Decode(msgid, true, false)
	if err != nil {
		return msgid + "R" // best we can do
	}
	suf = "R"
	for {
		msgid, _ = messageid.Encode(pfx, seq, suf)
		if !slices.ContainsFunc(i.Log, func(e *incident.LogEntry) bool { return e.LocalMsgID == msgid }) {
			return msgid
		}
		suf = string(suf[0] + 1)
	}
}

func (s *Server) servePostDeleteMessage(w http.ResponseWriter, r *http.Request) {
	serveLogIdent(w, r, true, func(i *incident.Incident, le *incident.LogEntry) error {
		return i.DeleteMessage(le)
	}, nil)
}
