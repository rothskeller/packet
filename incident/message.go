package incident

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/rothskeller/packet/message"
)

// saveMessage writes a message to disk, in both text and PDF formats, and
// creates symlinks to both from the RMI if appropriate.
func (i *Incident) saveMessage(msg message.Message, le *LogEntry) (err error) {
	var fname, lname, pname string

	fname = filepath.Join(i.Dir, le.Filename)
	if err = message.Write(msg, fname); err != nil {
		return err
	}
	if le.Linkname != "" {
		lname = filepath.Join(i.Dir, le.Linkname)
		os.Remove(lname)
		if err = os.Symlink(le.Filename, lname); err != nil {
			slog.Error("os.Symlink", "from", lname, "to", le.Filename, "err", err)
			return err
		}
	}
	pname = strings.TrimSuffix(le.Filename, ".txt") + ".pdf"
	fname = filepath.Join(i.Dir, pname)
	switch err = msg.Type().RenderPDF(msg, fname, ""); err.(type) {
	case nil:
	case message.Warning:
		slog.Warn("RenderPDF", "fname", fname, "warn", err)
	default:
		slog.Error("RenderPDF", "fname", fname, "err", err)
		return err
	}
	if le.Linkname != "" {
		lname = filepath.Join(i.Dir, strings.TrimSuffix(le.Linkname, ".txt")+".pdf")
		os.Remove(lname)
		if err = os.Symlink(pname, lname); err != nil {
			slog.Error("os.Symlink", "from", lname, "to", pname, "err", err)
			return err
		}
	}
	return nil
}

func (i *Incident) makeUniqueFilename(base, ext string) string {
	var (
		seq   = 1
		fname = filepath.Join(i.Dir, base+ext)
	)
	for {
		if _, err := os.Stat(fname); err != nil {
			return filepath.Base(fname)
		}
		seq++
		fname = filepath.Join(i.Dir, fmt.Sprintf("%s-%d%s", base, seq, ext))
	}
}
