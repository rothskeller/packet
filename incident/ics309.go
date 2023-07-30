package incident

import (
	_ "embed" // .
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message/common"
	"github.com/rothskeller/pdf/pdfform"
	"github.com/rothskeller/pdf/pdfstruct"
)

// An ICS309Header structure contains all of the information needed for the
// header of an ICS-309 communications log.
type ICS309Header struct {
	IncidentName  string
	ActivationNum string
	OpStartDate   string
	OpStartTime   string
	OpEndDate     string
	OpEndTime     string
	OpCall        string
	OpName        string
	TacCall       string
	TacName       string
}

// GenerateICS309 generates an ICS-309 communications log covering all of the
// messages in the directory.  It returns the names of the generated files, or
// an error if the directory could not be read or the files could not be
// written.  Note that the generated files are removed by any call to
// SaveMessage or SaveReceipt, since they could be stale.
func GenerateICS309(header *ICS309Header) (csv string, pdfs []string, err error) {
	var (
		dir   *os.File
		files []os.FileInfo
		rmis  map[string]string
		lmis  map[*envelope.Envelope]string
		msgs  []*envelope.Envelope
		form  [][]string
	)
	if dir, err = os.Open("."); err != nil {
		return "", nil, err
	}
	defer dir.Close()
	if files, err = dir.Readdir(0); err != nil {
		return "", nil, err
	}
	for _, fi := range files {
		if !strings.HasSuffix(fi.Name(), ".txt") {
			continue
		}
		switch fi.Mode().Type() {
		case 0: // regular file
			var (
				lmi     string
				rcpt    bool
				content []byte
				env     *envelope.Envelope
			)
			if strings.HasSuffix(fi.Name(), ".DR.txt") || strings.HasSuffix(fi.Name(), ".RR.txt") {
				lmi = fi.Name()[:len(fi.Name())-7]
				rcpt = true
			} else {
				lmi = fi.Name()[:len(fi.Name())-4]
			}
			if !msgIDRE.MatchString(lmi) {
				continue
			}
			if content, err = os.ReadFile(fi.Name()); err != nil {
				continue
			}
			if env, _, err = envelope.ParseSaved(string(content)); err != nil {
				continue
			}
			if !env.IsFinal() {
				continue
			}
			msgs = append(msgs, env)
			if !rcpt {
				lmis[env] = lmi
			}
		case os.ModeSymlink:
			var (
				rmi string
				lmi string
			)
			rmi = fi.Name()[:len(fi.Name())-4]
			if !msgIDRE.MatchString(rmi) {
				continue
			}
			if lmi, err = os.Readlink(fi.Name()); err != nil {
				continue
			}
			if !strings.HasSuffix(lmi, ".txt") {
				continue
			}
			lmi = lmi[:len(lmi)-4]
			if !msgIDRE.MatchString(lmi) {
				continue
			}
			rmis[lmi] = rmi
		}
	}
	// Sort the list chronologically.
	sort.Slice(msgs, func(i, j int) bool { return envelopeLess(msgs[i], msgs[j]) })
	// Generate the form data.
	for _, m := range msgs {
		form = append(form, make309Line(m, lmis[m], rmis[lmis[m]]))
	}
	// Render the form.
	RemoveICS309s()
	if csv, err = render309CSV(header, form); err != nil {
		return "", nil, err
	}
	if pdfs, err = render309PDF(header, form); err != nil {
		return "", nil, err
	}
	return csv, pdfs, nil
}

// envelopeLess is the comparison function for sorting the message list.
func envelopeLess(a, b *envelope.Envelope) bool {
	var at, bt time.Time
	if a.IsReceived() {
		at = a.ReceivedDate
	} else {
		at = a.Date
	}
	if b.IsReceived() {
		bt = b.ReceivedDate
	} else {
		bt = b.Date
	}
	return at.Before(bt)
}

// make309Line generates the ICS-309 form data for a single message.
func make309Line(m *envelope.Envelope, lmi, rmi string) []string {
	var t time.Time
	var from, oid, to, did, sub string
	if m.IsReceived() {
		t = m.ReceivedDate
		from, _, _ = strings.Cut(m.From, "@")
		from = strings.ToUpper(from)
		oid, did = rmi, lmi
		if m.ReceivedArea != "" {
			to = strings.ToUpper(m.ReceivedArea)
		}
	} else {
		t = m.Date
		oid, did = lmi, rmi
		if ta := m.To; len(ta) != 0 {
			to = ta[0]
		}
		to, _, _ = strings.Cut(to, "@")
		to = strings.ToUpper(to)
	}
	sub = m.SubjectLine
	if strings.HasPrefix(m.SubjectLine, "DELIVERED: ") {
		if msgid, _, _, _, _ := common.DecodeSubject(m.SubjectLine[11:]); msgid != "" {
			sub = "Delivery receipt for " + msgid
		}
	} else if strings.HasPrefix(m.SubjectLine, "READ: ") {
		if msgid, _, _, _, _ := common.DecodeSubject(m.SubjectLine[6:]); msgid != "" {
			sub = "Read receipt for " + msgid
		}
	}
	if strings.HasPrefix(sub, oid+"_") {
		sub = sub[len(oid)+1:]
	}
	if strings.HasPrefix(sub, did+"_") {
		sub = sub[len(did)+1:]
	}
	return []string{t.Format("01/02/2006 15:04"), from, oid, to, did, sub}
}

// render309CSV renders the ICS-309 in CSV format.
func render309CSV(header *ICS309Header, form [][]string) (filename string, err error) {
	var (
		fh *os.File
		w  *csv.Writer
	)
	filename = "ics309.csv"
	if fh, err = os.Create(filename); err != nil {
		return "", err
	}
	defer fh.Close()
	w = csv.NewWriter(fh)
	defer w.Flush()
	w.Write([]string{"ICS 309 COMMUNICATIONS LOG"})
	w.Write([]string{"Incident Name:", header.IncidentName})
	w.Write([]string{"Activation Number:", header.ActivationNum})
	w.Write([]string{"Operational Period:", fmt.Sprintf("%s %s to %s %s",
		header.OpStartDate, header.OpStartTime, header.OpEndDate, header.OpEndTime)})
	if header.TacCall != "" {
		w.Write([]string{"Tactical Station:", header.TacName, header.TacCall})
	}
	w.Write([]string{"Radio Operator:", header.OpName, header.OpCall})
	w.Write([]string{"Prepared:", time.Now().Format("01/02/2006 15:04")})
	w.Write([]string{})
	w.Write([]string{"Date/Time", "From Station", "Origin Msg ID", "To Station", "Dest Msg ID", "Subject"})
	w.WriteAll(form)
	return filename, nil
}

// render309PDF renders the ICS-309 in PDF format.
func render309PDF(header *ICS309Header, form [][]string) (filenames []string, err error) {
	var (
		fh    *os.File
		pdf   *pdfstruct.PDF
		pages int
	)
	if ics309pdf == nil { // built without PDF support
		return nil, err
	}
	pages = (len(form) + 30) / 31
	if pages == 0 {
		pages = 1
	}
	for page := 1; page <= pages; page++ {
		filename := "ics309.pdf"
		if pages > 1 {
			filename = fmt.Sprintf("ics309-p%d.pdf", page)
		}
		if fh, err = os.Create(filename); err != nil {
			return nil, err
		}
		defer fh.Close()
		if _, err = fh.Write(ics309pdf); err != nil {
			os.Remove(filename)
			return nil, err
		}
		if pdf, err = pdfstruct.Open(fh); err != nil {
			os.Remove(filename)
			return nil, err
		}
		pdfform.SetField(pdf, "Incident Name", header.IncidentName, 0)
		pdfform.SetField(pdf, "Activation Number", header.ActivationNum, 0)
		pdfform.SetField(pdf, "OpPeriod Start Date", header.OpStartDate, 0)
		pdfform.SetField(pdf, "OpPeriod Start Time", header.OpStartTime, 0)
		pdfform.SetField(pdf, "OpPeriod End Date", header.OpEndDate, 0)
		pdfform.SetField(pdf, "OpPeriod End Time", header.OpEndTime, 0)
		pdfform.SetField(pdf, "Tactical Station", header.TacName+" "+header.TacCall, 0)
		pdfform.SetField(pdf, "Operator", header.OpName+" "+header.OpCall, 0)
		pdfform.SetField(pdf, "Prepared By", header.OpName+" "+header.OpCall, 0)
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
			return nil, err
		}
		filenames = append(filenames, filename)
	}
	return filenames, nil
}

// RemoveICS309s removes generated ICS-309 communication log files.
func RemoveICS309s() {
	os.Remove("ics309.csv")
	os.Remove("ics309.pdf")
	pages, _ := filepath.Glob("ics309-p*.pdf")
	for _, page := range pages {
		os.Remove(page)
	}
}
