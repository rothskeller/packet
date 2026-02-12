package incident

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/receipt"
)

// GetLogEntryByIdent returns the log entry with the specified ident, if any.
func (i *Incident) GetLogEntryByIdent(ident int) *LogEntry {
	for _, le := range i.Log {
		if le.Ident == ident {
			return le
		}
	}
	return nil
}

// GetMessageByLMI returns the message with the specified local message ID, if
// any.  If the returned message is non-nil, the returned error is non-fatal.
func (i *Incident) GetMessageByLMI(lmi string) (msg message.Message, err error) {
	if e := i.logEntryForLMI(lmi); e != nil {
		return i.GetMessageFromLogEntry(e)
	}
	return nil, nil
}

// logEntryForLMI returns the first log entry for the message with the
// specified LMI, if any.
func (i *Incident) logEntryForLMI(lmi string) (entry *LogEntry) {
	for _, e := range i.Log {
		if e.LocalMsgID == lmi && e.Flags&FIsReceipt == 0 {
			return e
		}
	}
	return nil
}

// GetMessageFromLogEntry returns the message described by the provided log
// entry, if any.  If the returned  message is non-nil, any returned error is
// non-fatal.
func (i *Incident) GetMessageFromLogEntry(entry *LogEntry) (msg message.Message, err error) {
	if fname := entry.Filename(); fname == "" {
		return nil, nil
	} else {
		return message.Read(filepath.Join(i.Dir, fname))
	}
}

// saveMessage writes a message to disk, in both text and PDF formats, and
// creates symlinks to both from the RMI if appropriate.
func (i *Incident) saveMessage(msg message.Message, le *LogEntry) (err error) {
	var basename, fname, baselink, lname string

	basename = le.Filename()
	fname = filepath.Join(i.Dir, basename)
	if err = message.Write(msg, fname); err != nil {
		return err
	}
	if baselink = le.Linkname(); baselink != "" {
		lname = filepath.Join(i.Dir, baselink)
		if err = replaceLink(basename, lname); err != nil {
			return err
		}
	}
	if _, ok := msg.(*message.DraftMessage); ok {
		return nil // no PDF generation for unsent messages
	}
	switch msg.Body().(type) {
	case *receipt.DeliveryReceiptBody, *receipt.ReadReceiptBody:
		return nil // no PDF generation for receipts
	}
	basename = ToPDF(basename)
	fname = filepath.Join(i.Dir, basename)
	switch err = msg.Type().RenderPDF(msg, fname, ""); err.(type) {
	case nil:
	case message.Warning:
		slog.Warn("RenderPDF", "fname", fname, "warn", err)
	default:
		slog.Error("RenderPDF", "fname", fname, "err", err)
		return err
	}
	if baselink != "" {
		baselink = ToPDF(baselink)
		lname = filepath.Join(i.Dir, baselink)
		if err = replaceLink(basename, lname); err != nil {
			return err
		}
	}
	return nil
}

// replaceLink is like os.Symlink, but if a link already exists there, it
// overwrites it, and if something other than a link already exists there, it
// is a silent no-op.
func replaceLink(to, from string) (err error) {
	if stat, err := os.Stat(from); err != nil && !os.IsNotExist(err) {
		slog.Error("os.Stat", "f", from, "err", err)
		return err
	} else if err == nil {
		if stat.Mode().Type() != os.ModeSymlink {
			return nil // don't overwrite anything other than a symlink
		} else {
			os.Remove(from)
		}
	}
	if err = os.Symlink(to, from); err != nil {
		slog.Error("os.Symlink", "from", from, "to", to, "err", err)
		// Errors creating symoblic links are logged but do not cause
		// an abort.  On Windows, base users do not have permission to
		// create symlinks.
	}
	return nil
}
