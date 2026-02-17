package incident

import (
	"bytes"
	_ "embed"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rothskeller/packet/form/formdef"
	"github.com/rothskeller/pdf/v2"
)

//go:embed ics-309.pdf
var ics309Template []byte

//go:embed ics-309.form
var ics309FormFile []byte
var ics309FormDef *formdef.FormDef
var ics309FormInit sync.Once

// ICS309FormDef returns the form definition for the ICS-309 form.  This is
// used by other code to determine edit widths for the ICS-309 fields.
func ICS309FormDef() *formdef.FormDef {
	ics309FormInit.Do(func() {
		var err error
		if ics309FormDef, err = formdef.ReadFH(nil, bytes.NewReader(ics309FormFile), ""); err != nil {
			panic(err.Error())
		}
	})
	return ics309FormDef
}

// GenerateICS309 creates ics309.pdf in the incident directory.
func (i *Incident) GenerateICS309(signature string) (err error) {
	var (
		fname string
		out   *os.File
		op    *pdf.PDF
		ip    *pdf.PDF
		imp   *pdf.Importer
	)
	ics309FormInit.Do(func() {
		if ics309FormDef, err = formdef.ReadFH(nil, bytes.NewReader(ics309FormFile), ""); err != nil {
			panic(err.Error())
		}
	})
	fname = filepath.Join(i.Dir, "ics309.pdf")
	if out, err = os.Create(fname); err != nil {
		slog.Error("os.Create", "f", fname, "err", err)
		return err
	}
	op = pdf.New(out)
	op.Info["Title"] = "ICS-309 Communications Log"
	op.Info["Producer"] = "https://github.com/rothskeller/packet"
	if ip, err = pdf.Open(bytes.NewReader(ics309Template)); err != nil {
		slog.Error("pdf.Open", "err", err)
		return err
	}
	if imp, err = op.NewImporter(ip); err != nil {
		slog.Error("pdf.NewImporter", "err", err)
		return err
	}
	if err = i.render309(op, fname, imp, signature, false); err != nil {
		return err
	}
	if err = i.render309(op, fname, imp, signature, true); err != nil {
		return err
	}
	if err = op.Write(); err != nil {
		slog.Error("pdf.Write", "f", fname, "err", err)
		return err
	}
	if err = out.Close(); err != nil {
		slog.Error("os.Close", "f", fname, "err", err)
		return err
	}
	return nil
}

// render309 renders the ICS-309 pages, if any, for one message type.  It
// renders for voice messages if voice is true and for packet messages
// otherwise.
func (i *Incident) render309(out *pdf.PDF, fname string, in *pdf.Importer, signature string, voice bool) (err error) {
	var (
		log        []*LogEntry
		pageOffset int
		pageCount  int
	)
	// Gather the log entries for the specified type.
	for _, e := range i.Log {
		if (e.Flags&FVoice != 0) == voice && (e.Status == StatusHandEntered || e.Status == StatusReceived || e.Status == StatusSent) {
			log = append(log, e)
		}
	}
	if len(log) == 0 && (voice || len(i.Log) != 0) {
		return
	}
	if pageCount = (len(log) + 30) / 31; pageCount == 0 {
		pageCount = 1
	}
	if pageOffset, err = out.NumPages(); err != nil {
		slog.Error("pdf.NumPages", "f", fname, "err", err)
		return err
	}
	for line, e := range log {
		page := line/31 + 1
		if line%31 == 0 {
			if err = i.start309Page(out, fname, in, signature, voice, pageOffset, page, pageCount); err != nil {
				return err
			}
		}
		// Find the form "field" for the line.
		tag := "Line" + strconv.Itoa(line%31+1)
		idx := slices.IndexFunc(ics309FormDef.Fields, func(fd *formdef.FieldDef) bool { return fd.Tag == tag })
		fd := ics309FormDef.Fields[idx]
		// Emit the line.
		for i, r := range fd.PDF {
			var value string

			switch i {
			case 0:
				value = e.Time.Format("15:04")
			case 1:
				value = e.FromCall
			case 2:
				value = e.FromMsgID
			case 3:
				value = e.ToCall
			case 4:
				value = e.ToMsgID
			case 5:
				value = e.Subject
				if e.FromMsgID != "" {
					value = strings.TrimPrefix(value, e.FromMsgID+"_")
				}
				if e.ToMsgID != "" {
					value = strings.TrimPrefix(value, e.ToMsgID+"_")
				}
			}
			if value == "" {
				continue
			}
			if tr, ok := r.Renderer.(formdef.TextRenderer); ok {
				tr.Page = pageOffset + page
				switch err := tr.Draw(out, value); err.(type) {
				case nil, pdf.ErrTextRendering: // OK
				default:
					slog.Error("pdf.Text.Draw", "f", fname, "err", err)
					return err
				}
			}
		}
	}
	return nil
}

// start309Page starts a new page of an ICS-309 form.
func (i *Incident) start309Page(out *pdf.PDF, fname string, in *pdf.Importer, signature string, voice bool, pageOffset, page, pageCount int) (err error) {
	if err = out.AddPage(pdf.USLetterPortrait); err != nil {
		slog.Error("pdf.AddPage", "f", fname, "err", err)
		return err
	}
	if err = in.ImportPage(1, page+pageOffset); err != nil {
		slog.Error("pdf.ImportPage", "f", fname, "err", err)
		return err
	}
	for _, fd := range ics309FormDef.Fields {
		var value string

		switch fd.Tag {
		case "IncidentName":
			value = i.Config.IncidentName
		case "ActivationNum":
			value = i.Config.ActivationNum
		case "OpStartDate":
			if !i.Config.OpStart.IsZero() {
				value = i.Config.OpStart.Format("01/02/2006")
			}
		case "OpStartTime":
			if !i.Config.OpStart.IsZero() {
				value = i.Config.OpStart.Format("15:04")
			}
		case "OpEndDate":
			if !i.Config.OpEnd.IsZero() {
				value = i.Config.OpEnd.Format("01/02/2006")
			}
		case "OpEndTime":
			if !i.Config.OpEnd.IsZero() {
				value = i.Config.OpEnd.Format("15:04")
			}
		case "TacNameCall":
			var vm string
			if voice {
				vm = "(VOICE)"
			}
			value = strings.Join(
				slices.DeleteFunc(
					[]string{i.Config.TacName, i.Config.TacCall, vm},
					func(s string) bool { return s == "" }),
				" ")
		case "OpNameCall":
			value = strings.Join(
				slices.DeleteFunc(
					[]string{i.Config.OpName, i.Config.OpCall},
					func(s string) bool { return s == "" }),
				" ")
		case "Signature":
			if signature != "" {
				value = "/s/ " + signature
			}
		case "DateTime":
			value = time.Now().Format("01/02/2006 15:04")
		case "PageNum":
			value = strconv.Itoa(page)
		case "PageCount":
			value = strconv.Itoa(pageCount)
		}
		if value == "" {
			continue
		}
		for _, r := range fd.PDF {
			if tr, ok := r.Renderer.(formdef.TextRenderer); ok {
				tr.Page = pageOffset + page
				switch err := tr.Draw(out, value); err.(type) {
				case nil, pdf.ErrTextRendering: // OK
				default:
					slog.Error("pdf.Text.Draw", "f", fname, "err", err)
					return err
				}
			}
		}
	}
	return nil
}
