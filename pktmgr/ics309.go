package pktmgr

import (
	"encoding/csv"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/rothskeller/packet/xscmsg/delivrcpt"
	"github.com/rothskeller/packet/xscmsg/readrcpt"
)

func (i *Incident) ics309() {
	i.ics309CSV()
	// TODO i.ics309PDF()
}

// ics309CSV (re-)generates the CSV version of the ICS-309 log for the incident.
func (i *Incident) ics309CSV() {
	var (
		fh       *os.File
		incName  string
		actNum   string
		opPeriod string
		w        *csv.Writer
		msgs     []*Message
		err      error
	)
	// It's possible that someone has added metadata to the report on disk.
	// Read it and preserve the metadata if any.
	if fh, err = os.Open("ics309.csv"); err == nil {
		r := csv.NewReader(fh)
		r.FieldsPerRecord = -1 // variable
		if rows, err := r.ReadAll(); err == nil {
			for _, row := range rows {
				if len(row) == 2 {
					switch row[0] {
					case "Incident Name":
						incName = row[1]
					case "Activation Number:":
						actNum = row[1]
					case "Operational Period:":
						opPeriod = row[1]
					}
				}
			}

		}
		fh.Close()
	}
	// Create the new file.
	if fh, err = os.Create("ics309.csv"); err != nil {
		return
	}
	defer fh.Close()
	w = csv.NewWriter(fh)
	defer w.Flush()
	w.Write([]string{"ICS 309 COMMUNICATIONS LOG"})
	w.Write([]string{"Incident Name:", incName})
	w.Write([]string{"Activation Number:", actNum})
	w.Write([]string{"Operational Period:", opPeriod})
	if i.config.TacCall != "" {
		w.Write([]string{"Tactical Station:", i.config.TacName, i.config.TacCall})
	}
	w.Write([]string{"Radio Operator:", i.config.OpName, i.config.OpCall})
	w.Write([]string{"Prepared:", time.Now().Format("01/02/2006 15:04")})
	w.Write([]string{})
	w.Write([]string{"Date/Time", "From Station", "Origin Msg ID", "To Station", "Dest Msg ID", "Subject"})
	// Now, build an ordered list of all messages that have been sent or
	// received (i.e., not draft or queued).
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
	// Now render each message in the list.
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
		w.Write([]string{t.Format("01/02/2006 15:04"), from, oid, to, did, sub})
	}
}
