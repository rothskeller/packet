package server

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/address"
	"github.com/rothskeller/packet/message/msgifc"
)

// servePostManualReceive handles POST /manual-receive requests, which contain
// a received message that was manually provided.
func (s *Server) servePostManualReceive(w http.ResponseWriter, r *http.Request) {
	var (
		makedr bool
		mtext  string
		msg    *message.JustReceivedMessage
		err    error
	)
	// Check parameters.
	makedr = r.FormValue("mrdr") != ""
	mtext = strings.ReplaceAll(r.FormValue("mrmsg"), "\r\n", "\n")
	if msg, err = message.NewJustReceivedMessage(mtext, r.FormValue("mrbbs"), ""); err != nil {
		slog.Error("NewJustReceivedMessage", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = incident.Write(r.FormValue("dir"), func(i *incident.Incident) error {
		var dr *message.DraftMessage

		if dr, err = i.ReceiveMessage(msg); err != nil {
			slog.Error("ReceiveMessage", "err", err)
			return err
		}
		if dr != nil && makedr {
			if drid, err := i.AddDraftMessage(dr, false); err != nil {
				return fmt.Errorf("queueing delivery receipt: %s", err)
			} else {
				w.Header().Set("X-Packet-Action", fmt.Sprintf("manual-send:%d", drid))
			}
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// serveGetManualSendCommand handles GET /manual-send-command requests.  These
// have dir= and id= parameters, with the incident directory and log entry
// ident number for a draft message in the incident.  The handler responds with
// a text/plain body with the JNOS command to send the message.
func (s *Server) serveGetManualSendCommand(w http.ResponseWriter, r *http.Request) {
	var (
		dir   string
		id    int
		msg   message.Message
		ident string
		to    []string
		cmd   strings.Builder
		eb    string
		err   error
	)
	// Get the message from the incident and make sure it's proper.
	dir = r.FormValue("dir")
	if id, err = strconv.Atoi(r.FormValue("id")); err != nil || id <= 0 {
		http.Error(w, fmt.Sprintf("invalid message ident %q", r.FormValue("id")), http.StatusBadRequest)
		return
	}
	err = incident.Read(dir, func(i *incident.Incident) (err error) {
		if i.Config.TacCall != "" {
			ident = i.Config.OpCall
		}
		if le := i.GetLogEntryByIdent(id); le != nil {
			msg, err = i.GetMessageFromLogEntry(le)
		} else {
			err = fmt.Errorf("invalid message ident %q", r.FormValue("id"))
		}
		return err
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if dm, ok := msg.(*message.DraftMessage); !ok {
		http.Error(w, "not an unsent outgoing message", http.StatusBadRequest)
		return
	} else {
		if addrs, err := address.ParseList(dm.To()); err != nil || len(addrs) == 0 {
			http.Error(w, "invalid or empty To: address list", http.StatusBadRequest)
			return
		} else {
			for _, addr := range addrs {
				to = append(to, addr.Address)
			}
		}
	}
	if r.FormValue("force") == "" {
		if err = message.ValidateMessage(msg, msgifc.VPacket); err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
	}
	// Generate the actual send command.
	if msg.Bulletin() {
		cmd.WriteString("SB ")
	} else if len(to) > 1 {
		cmd.WriteString("SC ")
	} else {
		cmd.WriteString("SP ")
	}
	cmd.WriteString(to[0])
	cmd.WriteByte('\n')
	if len(to) > 1 {
		cmd.WriteString(strings.Join(to[1:], ","))
		cmd.WriteByte('\n')
	}
	cmd.WriteString(msg.Subject().EncodedSubject())
	cmd.WriteByte('\n')
	eb = msg.Payload().Encode()
	cmd.WriteString(eb)
	if !strings.HasSuffix(eb, "\n") {
		cmd.WriteByte('\n')
	}
	cmd.WriteString("/EX\n")
	if ident != "" {
		fmt.Fprintf(&cmd, "# DE %s\n", ident)
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, cmd.String())
}

// servePostMarkSent handles POST /mark-sent requests, which mark a message as
// having been sent (manually).  They have dir= and id= parameters with the
// incident directory and log entry ident number.  They return 204 No Content
// unless an error occurs.
func (s *Server) servePostMarkSent(w http.ResponseWriter, r *http.Request) {
	var (
		dir string
		id  int
		msg message.Message
		err error
	)
	// Get the message from the incident and make sure it's proper.
	dir = r.FormValue("dir")
	if id, err = strconv.Atoi(r.FormValue("id")); err != nil || id <= 0 {
		http.Error(w, fmt.Sprintf("invalid message ident %q", r.FormValue("id")), http.StatusBadRequest)
		return
	}
	err = incident.Write(dir, func(i *incident.Incident) (err error) {
		if le := i.GetLogEntryByIdent(id); le == nil {
			return fmt.Errorf("invalid message ident %q", r.FormValue("id"))
		} else if msg, err = i.GetMessageFromLogEntry(le); err != nil {
			return err
		} else if dm, ok := msg.(*message.DraftMessage); !ok {
			return errors.New("not an unsent outgoing message")
		} else {
			return i.MarkMessageSent(dm, le)
		}
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// servePostMakeReceipt handles POST /make-receipt requests, which create a
// delivery receipt for a received message.  It has dir= and id= parameters.
// It responds with an error status and text/plain error message, or 204 No
// Content if successful.  In the latter case, the X-Packet-Action header is
// set to "manual-send:$IDENT" where $IDENT is the ident of the new receipt.
func (s *Server) servePostMakeReceipt(w http.ResponseWriter, r *http.Request) {
	var (
		dir  string
		id   int
		drid int
		err  error
	)
	// Get the message from the incident and make sure it's proper.
	dir = r.FormValue("dir")
	if id, err = strconv.Atoi(r.FormValue("id")); err != nil || id <= 0 {
		http.Error(w, fmt.Sprintf("invalid message ident %q", r.FormValue("id")), http.StatusBadRequest)
		return
	}
	err = incident.Write(dir, func(i *incident.Incident) (err error) {
		if le := i.GetLogEntryByIdent(id); le == nil {
			return fmt.Errorf("invalid message ident %q", r.FormValue("id"))
		} else if msg, err := i.GetMessageFromLogEntry(le); err != nil {
			return err
		} else if dr, err := i.MakeDeliveryReceipt(msg, le); err != nil {
			return err
		} else if drid, err = i.AddDraftMessage(dr, false); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("X-Packet-Action", fmt.Sprintf("manual-send:%d", drid))
	w.WriteHeader(http.StatusNoContent)
}
