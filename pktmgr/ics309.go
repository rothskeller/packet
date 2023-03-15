package pktmgr

import (
	_ "embed" // .
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/xscmsg/delivrcpt"
	"github.com/rothskeller/packet/xscmsg/readrcpt"
	"github.com/rothskeller/pdf/pdfform"
	"github.com/rothskeller/pdf/pdfstruct"
)

//go:embed ICS-309.pdf
var ics309pdf []byte

func (i *Incident) ics309() {
	var (
		msgs []*Message
		form [][]string
	)
	// Build an ordered list of all messages that have been sent or received
	// (i.e., not draft or queued).
	for _, mr := range i.list {
		if mr.M.IsReceived() || mr.M.IsSent() {
			msgs = append(msgs, mr.M)
			if mr.DR != nil && (mr.DR.IsReceived() || mr.DR.IsSent()) {
				msgs = append(msgs, mr.DR)
			}
			if mr.RR != nil {
				msgs = append(msgs, mr.RR)
			}
		}
	}
	sort.Slice(msgs, func(a, b int) bool {
		var at, bt time.Time
		if msgs[a].IsReceived() {
			at = msgs[a].Received
		} else {
			at = msgs[a].Sent
		}
		if msgs[b].IsReceived() {
			bt = msgs[b].Received
		} else {
			bt = msgs[b].Sent
		}
		return at.Before(bt)
	})
	// Generate the form data.
	for _, m := range msgs {
		var t time.Time
		var from, oid, to, did, sub string
		if m.IsReceived() {
			t = m.Received
			from, _, _ = strings.Cut(m.From.Address, "@")
			from = strings.ToUpper(from)
			if m == m.mandr.M { // i.e., not a received receipt
				oid = m.mandr.RMI
				did = m.mandr.LMI
			}
			if m.IsBulletin() {
				to = strings.ToUpper(m.Bulletin)
			}
			switch m.Type.Tag {
			case delivrcpt.Tag:
				sub = "Delivery receipt for " + m.mandr.LMI
			case readrcpt.Tag:
				sub = "Read receipt for " + m.mandr.LMI
			default:
				sub = m.RawMessage.Header.Get("Subject")
			}
		} else {
			t = m.Sent
			if m == m.mandr.M { // i.e., not a sent receipt
				oid = m.mandr.LMI
				did = m.mandr.RMI
			}
			if len(m.To) != 0 {
				to, _, _ = strings.Cut(m.To[0].Address, "@")
				to = strings.ToUpper(to)
			}
			switch m.Type.Tag {
			case delivrcpt.Tag:
				if m.mandr.RMI != "" {
					sub = "Delivery receipt for " + m.mandr.RMI
					break
				}
				fallthrough
			default:
				sub = m.Subject()
			}
		}
		if strings.HasPrefix(sub, oid+"_") {
			sub = sub[len(oid)+1:]
		}
		if strings.HasPrefix(sub, did+"_") {
			sub = sub[len(did)+1:]
		}
		form = append(form, []string{t.Format("01/02/2006 15:04"), from, oid, to, did, sub})
	}
	// Render the form in CSV and PDF formats.
	i.ics309CSV(form)
	i.ics309PDF(form)
}

// ics309CSV (re-)generates the CSV version of the ICS-309 log for the incident.
func (i *Incident) ics309CSV(form [][]string) {
	var (
		fh  *os.File
		w   *csv.Writer
		err error
	)
	if fh, err = os.Create("ics309.csv"); err != nil {
		return
	}
	defer fh.Close()
	w = csv.NewWriter(fh)
	defer w.Flush()
	w.Write([]string{"ICS 309 COMMUNICATIONS LOG"})
	w.Write([]string{"Incident Name:", i.config.IncidentName})
	w.Write([]string{"Activation Number:", i.config.ActivationNum})
	w.Write([]string{"Operational Period:", fmt.Sprintf("%s %s to %s %s",
		i.config.OpStartDate, i.config.OpStartTime, i.config.OpEndDate, i.config.OpEndTime)})
	if i.config.TacCall != "" {
		w.Write([]string{"Tactical Station:", i.config.TacName, i.config.TacCall})
	}
	w.Write([]string{"Radio Operator:", i.config.OpName, i.config.OpCall})
	w.Write([]string{"Prepared:", time.Now().Format("01/02/2006 15:04")})
	w.Write([]string{})
	w.Write([]string{"Date/Time", "From Station", "Origin Msg ID", "To Station", "Dest Msg ID", "Subject"})
	w.WriteAll(form)
}

// ics309PDF (re-)generates the PDF version of the ICS-309 log for the incident.
func (i *Incident) ics309PDF(form [][]string) {
	var (
		fh    *os.File
		pdf   *pdfstruct.PDF
		pages int
		err   error
	)
	pages = (len(form) + 30) / 31
	if pages == 0 {
		pages = 1
	}
	if pages > 1 {
		os.Remove("ics309.pdf") // just in case there's a previous single-page version around
	}
	for page := 1; page <= pages; page++ {
		filename := "ics309.pdf"
		if pages > 1 {
			filename = fmt.Sprintf("ics309-p%d.pdf", page)
		}
		if fh, err = os.Create(filename); err != nil {
			return
		}
		defer fh.Close()
		if _, err = fh.Write(ics309pdf); err != nil {
			os.Remove(filename)
			return
		}
		if pdf, err = pdfstruct.Open(fh); err != nil {
			os.Remove(filename)
			return
		}
		pdfform.SetField(pdf, "Incident Name", i.config.IncidentName, 0)
		pdfform.SetField(pdf, "Activation Number", i.config.ActivationNum, 0)
		pdfform.SetField(pdf, "OpPeriod Start Date", i.config.OpStartDate, 0)
		pdfform.SetField(pdf, "OpPeriod Start Time", i.config.OpStartTime, 0)
		pdfform.SetField(pdf, "OpPeriod End Date", i.config.OpEndDate, 0)
		pdfform.SetField(pdf, "OpPeriod End Time", i.config.OpEndTime, 0)
		pdfform.SetField(pdf, "Tactical Station", i.config.TacName+" "+i.config.TacCall, 0)
		pdfform.SetField(pdf, "Operator", i.config.OpName+" "+i.config.OpCall, 0)
		pdfform.SetField(pdf, "Prepared By", i.config.OpName+" "+i.config.OpCall, 0)
		pdfform.SetField(pdf, "Prepared Time", time.Now().Format("01/02/2006 15:04"), 0)
		pdfform.SetField(pdf, "Page Number", strconv.Itoa(page), 0)
		pdfform.SetField(pdf, "Page Count", strconv.Itoa(pages), 0)
		for i := 1; i <= 31; i++ {
			idx := (page-1)*31 + i - 1
			if idx >= len(form) {
				break
			}
			pdfform.SetField(pdf, fmt.Sprintf("L%02d Time", i), form[idx][0][11:], 0) // time only, no date
			pdfform.SetField(pdf, fmt.Sprintf("L%02d Fm Stn", i), form[idx][1], 0)
			pdfform.SetField(pdf, fmt.Sprintf("L%02d Fm OID", i), form[idx][2], 0)
			pdfform.SetField(pdf, fmt.Sprintf("L%02d To Stn", i), form[idx][3], 0)
			pdfform.SetField(pdf, fmt.Sprintf("L%02d To DID", i), form[idx][4], 0)
			pdfform.SetField(pdf, fmt.Sprintf("L%02d Msg", i), form[idx][5], 0)
		}
		if err = pdf.Write(); err != nil {
			os.Remove(filename)
			return
		}
	}
}
