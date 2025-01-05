package incident

import (
	"bytes"
	_ "embed" // .
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/pdf/pdftext"
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

var receiptExtRE = regexp.MustCompile(`\.[DR]R\d*\.txt$`)

// GenerateICS309 generates an ICS-309 communications log covering all of the
// messages in the directory.  The log is generated in CSV format (ics309.csv),
// and if PDF rendering support is built into the program, it is also generated
// in PDF format (ics309.pdf).  GenerateICS309 returns an error if the directory
// could not be read or the files could not be written.  Note that the generated
// files are removed by any call to SaveMessage or SaveReceipt, since they could
// be stale.
func GenerateICS309(header *ICS309Header) (err error) {
	var (
		dir   *os.File
		files []os.FileInfo
		msgs  []*envelope.Envelope
		form  [][]string
		lmis  = make(map[*envelope.Envelope]string)
	)
	if dir, err = os.Open("."); err != nil {
		return err
	}
	defer dir.Close()
	if files, err = dir.Readdir(0); err != nil {
		return err
	}
	for _, fi := range files {
		var (
			lmi     string
			rcpt    bool
			content []byte
			env     *envelope.Envelope
		)
		if !strings.HasSuffix(fi.Name(), ".txt") || !fi.Mode().IsRegular() {
			continue
		}
		if idxs := receiptExtRE.FindStringIndex(fi.Name()); idxs != nil {
			lmi, rcpt = fi.Name()[:idxs[0]], true
		} else {
			lmi = fi.Name()[:len(fi.Name())-4]
		}
		if !MsgIDRE.MatchString(lmi) {
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
	}
	// Sort the list chronologically.
	sort.Slice(msgs, func(i, j int) bool { return envelopeLess(msgs[i], msgs[j]) })
	// Generate the form data.
	for _, m := range msgs {
		if lines, err := make309Lines(m, lmis[m]); err != nil {
			return err
		} else {
			form = append(form, lines...)
		}
	}
	// Render the form.
	RemoveICS309s()
	if err = render309CSV(header, form); err != nil {
		return err
	}
	if err = render309PDF(header, form); err != nil {
		return err
	}
	return nil
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

// make309Lines generates one or more ICS-309 form lines for a single message.
func make309Lines(m *envelope.Envelope, lmi string) (lines [][]string, err error) {
	if m.IsReceived() {
		return [][]string{make309Line(m, lmi, nil)}, nil
	} else if lmi == "" { // outgoing receipt
		return [][]string{make309Line(m, "", &DeliveryInfo{Recipient: m.To})}, nil
	} else {
		if delivs, err := Deliveries(lmi); err != nil {
			return nil, err
		} else {
			lines = make([][]string, len(delivs))
			for i, deliv := range delivs {
				lines[i] = make309Line(m, lmi, deliv)
			}
		}
		return lines, nil
	}
}

func make309Line(m *envelope.Envelope, lmi string, deliv *DeliveryInfo) []string {
	var t time.Time
	var from, oid, to, did, sub string
	if deliv == nil {
		t = m.ReceivedDate
		from = m.From
		if addrs, err := envelope.ParseAddressList(from); err == nil && len(addrs) != 0 {
			from = addrs[0].Address
		}
		from, _, _ = strings.Cut(from, "@")
		from = strings.ToUpper(from)
		oid, _, _, _, _ = message.DecodeSubject(m.SubjectLine)
		did = lmi
		if m.ReceivedArea != "" {
			to = strings.ToUpper(m.ReceivedArea)
		}
	} else {
		t = m.Date
		oid, did = lmi, deliv.RemoteMessageID
		to = deliv.Recipient
		if addr, err := envelope.ParseAddress(to); err == nil {
			to = addr.Address
		}
		to, _, _ = strings.Cut(to, "@")
		to = strings.ToUpper(to)
	}
	sub = m.SubjectLine
	if strings.HasPrefix(m.SubjectLine, "DELIVERED: ") {
		if msgid, _, _, _, _ := message.DecodeSubject(m.SubjectLine[11:]); msgid != "" {
			sub = "Delivery receipt for " + msgid
		}
	} else if strings.HasPrefix(m.SubjectLine, "READ: ") {
		if msgid, _, _, _, _ := message.DecodeSubject(m.SubjectLine[6:]); msgid != "" {
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
func render309CSV(header *ICS309Header, form [][]string) (err error) {
	const filename = "ics309.csv"
	var (
		fh *os.File
		w  *csv.Writer
	)
	if fh, err = os.Create(filename); err != nil {
		return err
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
	return nil
}

// render309PDF renders the ICS-309 in PDF format.
func render309PDF(header *ICS309Header, form [][]string) (err error) {
	const filename = "ics309.pdf"
	var (
		rdr   io.ReadSeeker
		pdf   *gofpdf.Fpdf
		imp   *gofpdi.Importer
		tpl   int
		pages = (len(form) + 30) / 31
		page  = 1
	)
	if ics309pdf == nil { // built without PDF support
		return nil
	}
	// Create the output PDF and the importer from the base PDF.
	rdr = bytes.NewReader(ics309pdf)
	pdf = gofpdf.New("P", "pt", "Letter", "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.SetMargins(0, 0, 0)
	imp = gofpdi.NewImporter()
	// Add the form pages.
	tpl = imp.ImportPageFromStream(pdf, &rdr, 1, "/MediaBox")
	for len(form) != 0 {
		pdf.AddPage()
		imp.UseImportedTemplate(pdf, tpl, 0, 0, 612, 792)
		render309PDFHeaderFooter(pdf, header, page, pages)
		for i := 0; i < len(form) && i < 31; i++ {
			render309PDFLine(pdf, i, form[i])
		}
		form = form[min(len(form), 31):]
		page++
	}
	// Add the instructions page.
	pdf.AddPage()
	imp.UseImportedTemplate(pdf, imp.ImportPageFromStream(pdf, &rdr, 2, "/MediaBox"), 0, 0, 612, 792)
	// Write the file.
	if err = pdf.OutputFileAndClose(filename); err != nil {
		os.Remove(filename)
		return err
	}
	return nil
}

func render309PDFHeaderFooter(pdf *gofpdf.Fpdf, header *ICS309Header, page, pages int) {
	render309PDFString(pdf, header.IncidentName, 140, 50, 200, 12)
	render309PDFString(pdf, header.ActivationNum, 158, 65, 182, 12)
	render309PDFString(pdf, header.OpStartDate, 380, 50, 70, 12)
	render309PDFString(pdf, header.OpStartTime, 380, 65, 70, 12)
	render309PDFString(pdf, header.OpEndDate, 473, 50, 100, 12)
	render309PDFString(pdf, header.OpEndTime, 473, 65, 100, 12)
	render309PDFString(pdf, header.TacName+" "+header.TacCall, 41, 95, 267, 16)
	render309PDFString(pdf, header.OpName+" "+header.OpCall, 315, 95, 257, 16)
	render309PDFString(pdf, header.OpName+" "+header.OpCall, 42, 665, 167, 13)
	render309PDFString(pdf, time.Now().Format("01/02/2006 15:04"), 343, 665, 117, 13)
	render309PDFString(pdf, strconv.Itoa(page), 498, 665, 23, 13)
	render309PDFString(pdf, strconv.Itoa(pages), 540, 665, 24, 13)
}

func render309PDFLine(pdf *gofpdf.Fpdf, lnum int, fields []string) {
	y := 15.69*float64(lnum) + 166.15
	render309PDFString(pdf, fields[0][11:], 39, y, 48, 12) // time only, no date
	render309PDFString(pdf, fields[1], 93, y, 53, 12)
	render309PDFString(pdf, fields[2], 151, y, 58, 12)
	render309PDFString(pdf, fields[3], 215, y, 53, 12)
	render309PDFString(pdf, fields[4], 272, y, 63, 12)
	render309PDFString(pdf, fields[5], 340, y, 233, 12)
}

func render309PDFString(pdf *gofpdf.Fpdf, s string, x, y, w, h float64) {
	/* red background of rectangle for layout validation
	pdf.SetAlpha(0.5, "")
	pdf.SetFillColor(255, 0, 0)
	pdf.Rect(x, y, w, h, "F")
	pdf.SetAlpha(1.0, "")
	*/
	pdftext.Draw(pdf, s, x, y, w, h, pdftext.Style{
		MinFontSize: 10, VAlign: "baseline", Clip: 1, Color: []byte{0, 0, 153},
	})
}

// RemoveICS309s removes generated ICS-309 communication log files.
func RemoveICS309s() {
	os.Remove("ics309.csv")
	os.Remove("ics309.pdf")
}
