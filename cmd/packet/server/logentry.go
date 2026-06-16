package server

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/rothskeller/packet/v4/incident"
	"github.com/rothskeller/packet/v4/message"
	"github.com/rothskeller/packet/v4/message/address"
	"github.com/rothskeller/packet/v4/message/msgifc"
)

// servePostEditLogEntry handles POST /edit-log-entry requests, which edit or
// add log entries.  They have dir= and id= parameters, with the incident
// directory and either the id of a log entry to edit or "NEW" to create a new
// one.  They return an error status with a text/plain error message, or 204
// No Content, possibly with an X-Packet-Action haeder.
func (s *Server) servePostEditLogEntry(w http.ResponseWriter, r *http.Request) {
	var (
		dir   string
		idstr string
		id    int
		err   error
	)
	s.outpost = false
	// Get the message from the incident and make sure it's proper.
	dir, idstr = r.FormValue("dir"), r.FormValue("id")
	if idstr != "NEW" {
		if id, err = strconv.Atoi(idstr); err != nil || id <= 0 {
			http.Error(w, fmt.Sprintf("invalid message ident %q", r.FormValue("id")), http.StatusBadRequest)
			return
		}
		err = incident.Write(dir, func(i *incident.Incident) (err error) {
			if le := i.GetLogEntryByIdent(id); le == nil {
				return fmt.Errorf("invalid message ident %q", r.FormValue("id"))
			} else if le.Status == incident.StatusDraft || le.Status == incident.StatusQueued {
				return errors.New("The log entries for unsent messages cannot be edited.")
			} else {
				if le.Time, err = time.ParseInLocation("01/02/2006 15:04", r.FormValue("date")+" "+r.FormValue("time"), time.Local); err != nil {
					return errors.New("The date and/or time does not have the proper MM/DD/YYYY HH:MM format.")
				}
				le.FromCall = r.FormValue("fromcall")
				le.FromMsgID = r.FormValue("fromid")
				le.ToCall = r.FormValue("tocall")
				le.ToMsgID = r.FormValue("toid")
				le.Subject = r.FormValue("summ")
				if r.FormValue("fu") != "" {
					le.Flags |= incident.FFollowup
				} else {
					le.Flags &^= incident.FFollowup
				}
				if le.Status == incident.StatusSent && le.Flags&incident.FHasReceipt == 0 {
					if r.FormValue("nr") != "" {
						le.Flags |= incident.FNeedsReceipt
					} else {
						le.Flags &^= incident.FNeedsReceipt
					}
				}
				if le.Status == incident.StatusHandEntered {
					if r.FormValue("vm") != "" {
						le.Flags |= incident.FVoice
					} else {
						le.Flags &^= incident.FVoice
					}
				}
				i.UpdateLogEntry(le)
				return nil
			}
		})
	} else {
		var le incident.LogEntry

		if le.Time, err = time.ParseInLocation("01/02/2006 15:04", r.FormValue("date")+" "+r.FormValue("time"), time.Local); err != nil {
			slog.Error("invalid date/time")
			http.Error(w, "The date and/or time does not have the proper MM/DD/YYYY HH:MM format.", http.StatusBadRequest)
			return
		}
		le.FromCall = r.FormValue("fromcall")
		le.FromMsgID = r.FormValue("fromid")
		le.ToCall = r.FormValue("tocall")
		le.ToMsgID = r.FormValue("toid")
		le.Subject = r.FormValue("summ")
		if r.FormValue("fu") != "" {
			le.Flags |= incident.FFollowup
		}
		if r.FormValue("vm") != "" {
			le.Flags |= incident.FVoice
		}
		err = incident.Write(dir, func(i *incident.Incident) error {
			i.AddLogEntry(&le)
			w.Header().Set("X-Packet-Action", fmt.Sprintf("select:%d", le.Ident))
			return nil
		})
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// servePostResetLogEntry handles POST /reset-log-entry requests, which reset
// the log data for a message back to values derived from the message content.
// They take dir= and id= parameters, which are the incident directory and log
// entry ident.  They respond with an error status and text/plain error message,
// or with a 204 No Content.
func (s *Server) servePostResetLogEntry(w http.ResponseWriter, r *http.Request) {
	var (
		dir string
		id  int
		err error
	)
	s.outpost = false
	// Get the message from the incident and make sure it's proper.
	dir = r.FormValue("dir")
	if id, err = strconv.Atoi(r.FormValue("id")); err != nil || id <= 0 {
		http.Error(w, fmt.Sprintf("invalid message ident %q", r.FormValue("id")), http.StatusBadRequest)
		return
	}
	err = incident.Write(dir, func(i *incident.Incident) (err error) {
		if le := i.GetLogEntryByIdent(id); le == nil {
			return fmt.Errorf("invalid message ident %q", r.FormValue("id"))
		} else {
			return i.ResetLogEntry(le)
		}
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// servePostDeleteLogEntry handles POST /delete-log-entry requests, which
// delete human-entered log entries.  They take dir= and id= parameters, which
// are the incident directory and log entry ident.  They respond with an error
// status and text/plain error message, or with a 204 No Content.
func (s *Server) servePostDeleteLogEntry(w http.ResponseWriter, r *http.Request) {
	var (
		dir string
		id  int
		err error
	)
	s.outpost = false
	// Get the message from the incident and make sure it's proper.
	dir = r.FormValue("dir")
	if id, err = strconv.Atoi(r.FormValue("id")); err != nil || id <= 0 {
		http.Error(w, fmt.Sprintf("invalid message ident %q", r.FormValue("id")), http.StatusBadRequest)
		return
	}
	err = incident.Write(dir, func(i *incident.Incident) (err error) {
		if le := i.GetLogEntryByIdent(id); le == nil {
			return fmt.Errorf("invalid message ident %q", r.FormValue("id"))
		} else {
			return i.DeleteLogEntry(le)
		}
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// servePostToggleFlag handles POST /toggle-flag requests, which toggle a flag
// in a log entry.  They take dir=, id=, and flag= parameters, which are the
// incident directory, log entry ident, and flag character to be toggled.  They
// respond with an error status and text/plain error message, or with a 204 No
// Content.
func (s *Server) servePostToggleFlag(w http.ResponseWriter, r *http.Request) {
	var (
		dir  string
		id   int
		flag string
		warn error
		err  error
	)
	s.outpost = false
	// Get the message from the incident and make sure it's proper.
	dir, flag = r.FormValue("dir"), r.FormValue("flag")
	if flag != "F" && flag != "-" && flag != "Q" {
		http.Error(w, fmt.Sprintf("Cannot toggle the %q flag: that flag is not recognized.", flag), http.StatusBadRequest)
		return
	}
	if id, err = strconv.Atoi(r.FormValue("id")); err != nil || id <= 0 {
		http.Error(w, fmt.Sprintf("invalid message ident %q", r.FormValue("id")), http.StatusBadRequest)
		return
	}
	err = incident.Write(dir, func(i *incident.Incident) (err error) {
		if le := i.GetLogEntryByIdent(id); le == nil {
			return fmt.Errorf("invalid message ident %q", r.FormValue("id"))
		} else {
			switch flag {
			case "F":
				le.Flags ^= incident.FFollowup
				i.UpdateLogEntry(le)
			case "Q":
				warn, err = toggleQueued(i, le, r.FormValue("force") != "")
			case "-":
				if le.Status != incident.StatusSent {
					return errors.New("The Needs Receipt flag can only be set on a sent message.")
				} else if le.Flags&incident.FIsReceipt != 0 && le.Flags&incident.FNeedsReceipt == 0 {
					return errors.New("The Needs Receipt flag cannot be set on a receipt message.")
				} else if le.Flags&incident.FHasReceipt != 0 && le.Flags&incident.FNeedsReceipt == 0 {
					return errors.New("The Needs Receipt flag cannot be set on this message, because a receipt has been received for it.")
				} else {
					le.Flags ^= incident.FNeedsReceipt
				}
				i.UpdateLogEntry(le)
			}
			return err
		}
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if warn != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		io.WriteString(w, warn.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// toggleQueued toggles the ready-to-send flag for an unsent outgoing message.
// If the flag is being turned on, the message has validation problems, and
// force is false, the validation problems are returned in "warn" and no change
// is made.
func toggleQueued(i *incident.Incident, le *incident.LogEntry, force bool) (warn, err error) {
	var msg message.Message

	if msg, err = i.GetMessageFromLogEntry(le); err != nil {
		return nil, err
	}
	if dm, ok := msg.(*message.DraftMessage); !ok {
		return nil, errors.New("The Ready-to-Send flag can only be set on an unsent outgoing message.")
	} else if !dm.ReadyToSend() {
		if addrs, err := address.ParseList(dm.To()); err != nil || len(addrs) == 0 {
			return nil, errors.New("The message cannot be marked ready to send because it does not have a valid To: address.")
		}
		if !force {
			if warn = message.ValidateMessage(dm, msgifc.VPacket); warn != nil {
				return warn, nil
			}
		}
		dm.SetReadyToSend(true)
		return nil, i.UpdateDraftMessage(le.Ident, dm)
	} else {
		dm.SetReadyToSend(false)
		return nil, i.UpdateDraftMessage(le.Ident, dm)
	}
}
