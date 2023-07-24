package shell

import (
	_ "embed" // .
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message/common"
	"github.com/rothskeller/pdf/pdfform"
	"github.com/rothskeller/pdf/pdfstruct"
)

//go:embed ICS-309.pdf
var ics309pdf []byte

// cmdICS309 implements the ics309 command.
func cmdICS309(args []string) bool {
	var (
		rmis map[string]string
		lmis map[*envelope.Envelope]string
		msgs []*envelope.Envelope
		form [][]string
	)
	if len(args) != 0 {
		io.WriteString(os.Stderr, "usage: packet ics309\n")
		return false
	}
	// If the generated ICS-309 exists, just show it.
	if _, err := os.Stat("ics309.csv"); !errors.Is(err, os.ErrNotExist) {
		return showICS309()
	}
	// Make sure we have the incident settings.
	if !requestConfig("incident", "activation", "period", "tactical", "operator") {
		return false
	}
	// Build a map of RMIs.
	rmis = make(map[string]string)
	forEachMessageSymlink(func(lmi, rmi string) {
		rmis[lmi] = rmi
	})
	// Build a list of message envelopes.
	lmis = make(map[*envelope.Envelope]string)
	forEachMessageFile(func(lmi string) {
		env, _, err := readMessage(lmi)
		if err != nil || !env.IsFinal() {
			return
		}
		msgs = append(msgs, env)
		lmis[env] = lmi
		if env, _, err = readMessage(lmi + ".DR"); err == nil {
			msgs = append(msgs, env)
		}
		if env, _, err = readMessage(lmi + ".RR"); err == nil {
			msgs = append(msgs, env)
		}
	})
	// Sort the list chronologically.
	sort.Slice(msgs, func(i, j int) bool { return envelopeLess(msgs[i], msgs[j]) })
	// Generate the form data.
	for _, m := range msgs {
		form = append(form, make309Line(m, lmis[m], rmis[lmis[m]]))
	}
	// Render the form.
	render309CSV(form)
	render309PDF(form)
	// Show the PDF form.
	return showICS309()
}

func helpICS309() {
	io.WriteString(os.Stdout, `The "ics309" command generates and displays an ICS-309 log.
    usage: packet ics309
The "ics309" command generates an ICS-309 communications log, in both CSV and
PDF formats.  It lists all sent and received messages in the current incident
(i.e., current working directory), including receipts.  The generated log is
stored in "ics309.csv" and "ics309.pdf".  (If multiple pages are needed, each
page is stored in "ics309-p##.pdf".)  After generating the log, the "ics309"
command opens the formatted PDF version in the system PDF viewer.
    NOTE:  Packet commands automatically remove the saved ICS-309 files after
any change to any message, to avoid reliance on a stale communications log.
Simply run "ics309" again to generate a new one.
`)
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
func render309CSV(form [][]string) {
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
	w.Write([]string{"Incident Name:", config.IncidentName})
	w.Write([]string{"Activation Number:", config.ActivationNum})
	w.Write([]string{"Operational Period:", fmt.Sprintf("%s %s to %s %s",
		config.OpStartDate, config.OpStartTime, config.OpEndDate, config.OpEndTime)})
	if config.TacCall != "" {
		w.Write([]string{"Tactical Station:", config.TacName, config.TacCall})
	}
	w.Write([]string{"Radio Operator:", config.OpName, config.OpCall})
	w.Write([]string{"Prepared:", time.Now().Format("01/02/2006 15:04")})
	w.Write([]string{})
	w.Write([]string{"Date/Time", "From Station", "Origin Msg ID", "To Station", "Dest Msg ID", "Subject"})
	w.WriteAll(form)
}

// render309PDF renders the ICS-309 in PDF format.
func render309PDF(form [][]string) {
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
		pdfform.SetField(pdf, "Incident Name", config.IncidentName, 0)
		pdfform.SetField(pdf, "Activation Number", config.ActivationNum, 0)
		pdfform.SetField(pdf, "OpPeriod Start Date", config.OpStartDate, 0)
		pdfform.SetField(pdf, "OpPeriod Start Time", config.OpStartTime, 0)
		pdfform.SetField(pdf, "OpPeriod End Date", config.OpEndDate, 0)
		pdfform.SetField(pdf, "OpPeriod End Time", config.OpEndTime, 0)
		pdfform.SetField(pdf, "Tactical Station", config.TacName+" "+config.TacCall, 0)
		pdfform.SetField(pdf, "Operator", config.OpName+" "+config.OpCall, 0)
		pdfform.SetField(pdf, "Prepared By", config.OpName+" "+config.OpCall, 0)
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

// showICS309 opens the system PDF viewer to show the generated ICS-309 log.
func showICS309() bool {
	var pages []string

	if _, err := os.Stat("ics309.pdf"); err == nil {
		pages = []string{"ics309.pdf"}
	} else {
		pages, _ = filepath.Glob("ics309-p*.pdf")
	}
	if len(pages) == 0 {
		io.WriteString(os.Stderr, "ERROR: generated ICS-309 PDF files are missing\n")
		return false
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd.exe", append([]string{"/C"}, pages...)...)
	case "darwin":
		cmd = exec.Command("open", pages...)
	default:
		cmd = exec.Command("xdg-open", pages...)
	}
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: PDF viewer could not be started: %s\n", err)
		return false
	}
	go func() { cmd.Wait() }()
	return true
}

// removeICS309s removes generated ICS-309 communication log files.  It is
// called whenever there is any change that would invalidate the log.
func removeICS309s() {
	os.Remove("ics309.csv")
	os.Remove("ics309.pdf")
	pages, _ := filepath.Glob("ics309-p*.pdf")
	for _, page := range pages {
		os.Remove(page)
	}
}
