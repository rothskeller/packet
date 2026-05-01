package cio

import (
	"encoding/csv"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/rothskeller/packet/incident"
)

// EmitLogList emits a list of log entries.
func (cio *CIO) EmitLogList(entries []*incident.LogEntry, full, numbers, receipts bool, myID string) {
	if !cio.OutputIsTerm {
		emitLogCSV(entries)
		return
	}
	var toshow []*incident.LogEntry
	for _, e := range entries {
		switch {
		case e.Status == incident.StatusDeleted: // never show
		case e.Flags&incident.FIsReceipt != 0 && !receipts && (e.Status == incident.StatusReceived || e.Status == incident.StatusSent): // hide receipts
		default:
			toshow = append(toshow, e)
		}
	}
	if len(toshow) == 0 {
		cio.Confirm("There are no log entries to show.")
		return
	}
	if full {
		cio.emitLogFull(toshow, numbers)
	} else {
		cio.emitLogCompact(toshow, numbers, myID)
	}
}

func (cio *CIO) emitLogFull(entries []*incident.LogEntry, numbers bool) {
	var (
		nw, fcw, fmw, tcw, tmw, mw = 1, 7, 7, 7, 7, 0
		fu                         bool
	)
	for _, e := range entries {
		text, _ := fullLine(e)
		nw = max(nw, len(text[0]))
		fcw = max(fcw, len(text[2]))
		fmw = max(fmw, len(text[3]))
		tcw = max(tcw, len(text[4]))
		tmw = max(tmw, len(text[5]))
		if e.Flags&incident.FFollowup != 0 {
			fu = true
		}
	}
	mw = cio.Width - fcw - fmw - tcw - tmw - 15
	if numbers {
		mw -= nw + 2
	}
	if fu {
		mw -= 2
	}
	if numbers {
		cio.print(colorLabel, setLength("N", nw)+"  ")
	}
	cio.print(colorLabel, "TIME   ")
	cio.print(colorLabel, setLength("FROM", fcw)+"  ")
	cio.print(colorLabel, setLength("MSG ID", fmw)+"  ")
	cio.print(colorLabel, setLength("TO", tcw)+"  ")
	cio.print(colorLabel, setLength("MSG ID", tmw)+"  ")
	cio.print(colorLabel, setMaxLength("MESSAGE", mw)+"\n")
	for _, e := range entries {
		text, colors := fullLine(e)
		if numbers {
			cio.print(colors[0], setLength(text[0], nw)+"  ")
		}
		cio.print(colors[1], text[1])
		cio.print(colors[0], "  ")
		cio.print(colors[2], setLength(text[2], fcw)+"  ")
		cio.print(colors[3], setLength(text[3], fmw)+"  ")
		cio.print(colors[4], setLength(text[4], tcw)+"  ")
		cio.print(colors[5], setMaxLength(text[5], tmw))
		if len(text[5]) < tmw {
			cio.print(colors[0], spaces[:tmw-len(text[5])])
		}
		cio.print(colors[0], "  ")
		w := mw
		if e.Flags&incident.FUnread != 0 {
			cio.print(colorAlertBG, "NEW")
			cio.print(colors[6], " ")
			w -= 4
		}
		if e.Flags&incident.FVoice != 0 {
			cio.print(colorSuccessBG, "VOICE")
			cio.print(colors[6], " ")
			w -= 6
		}
		switch {
		case w > 0 && e.Flags&incident.FFollowup != 0:
			cio.print(colors[6], setLength(text[6], w)+" ")
			cio.print(colorImmediate, "★\n")
		case w > 0:
			cio.print(colors[6], setMaxLength(text[6], w)+"\n")
		default:
			cio.print(colors[0], "\n")
		}
	}

}

func fullLine(e *incident.LogEntry) (columns []string, colors []int) {
	columns = make([]string, 7)
	var color = colorNormal
	switch {
	case e.Flags&incident.FImmediate != 0:
		color = colorImmediate
	case e.Flags&incident.FPriority != 0:
		color = colorPriority
	case e.Flags&incident.FBulletin != 0:
		color = colorBulletin
	case e.Flags&incident.FIsReceipt != 0:
		color = colorReceipt
	}
	colors = slices.Repeat([]int{color}, 7)
	columns[0] = strconv.Itoa(e.Ident)
	switch e.Status {
	case incident.StatusDraft:
		columns[1] = "DRAFT"
		colors[1] = colorAlertBG
	case incident.StatusQueued:
		columns[1] = "READY"
		colors[1] = colorWarningBG
	default:
		columns[1] = e.Time.Format("15:04")
	}
	columns[2] = e.FromCall
	columns[3] = e.FromMsgID
	columns[4] = e.ToCall
	if e.ToMsgID != "" {
		columns[5] = e.ToMsgID
	} else if e.Flags&incident.FNeedsReceipt != 0 {
		columns[5] = "NO RCPT"
		colors[5] = colorWarningBG
	}
	columns[6] = e.Subject
	return columns, colors
}

func (cio *CIO) emitLogCompact(entries []*incident.LogEntry, numbers bool, myID string) {
	var (
		nw, fw, lw, tw, mw = 1, 7, 7, 7, 0
		fu                 bool
	)
	for _, e := range entries {
		text, _ := compactLine(e, myID)
		nw = max(nw, len(text[0]))
		fw = max(fw, len(text[2]))
		lw = max(lw, len(text[3]))
		tw = max(tw, len(text[4]))
		if e.Flags&incident.FFollowup != 0 {
			fu = true
		}
	}
	mw = cio.Width - fw - lw - tw - 15
	if numbers {
		mw -= nw + 2
	}
	if fu {
		mw -= 2
	}
	if numbers {
		cio.print(colorLabel, setLength("N", nw)+"  ")
	}
	cio.print(colorLabel, "TIME   ")
	cio.print(colorLabel, setLength("FROM", fw)+"   ")
	cio.print(colorLabel, setLength("MSG ID", lw)+"   ")
	cio.print(colorLabel, setLength("TO", tw)+"  ")
	cio.print(colorLabel, setMaxLength("MESSAGE", mw)+"\n")
	for _, e := range entries {
		text, colors := compactLine(e, myID)
		if numbers {
			cio.print(colors[0], setLength(text[0], nw)+"  ")
		}
		cio.print(colors[1], text[1])
		cio.print(colors[0], "  ")
		cio.print(colors[2], setMaxLength(text[2], fw))
		if len(text[2]) < fw {
			cio.print(colors[0], spaces[:fw-len(text[2])])
		}
		if e.Status == incident.StatusReceived {
			cio.print(colors[0], " → ")
		} else {
			cio.print(colors[0], "   ")
		}
		cio.print(colors[3], setLength(text[3], lw))
		switch e.Status {
		case incident.StatusDraft, incident.StatusQueued, incident.StatusSent:
			cio.print(colors[0], " → ")
		default:
			cio.print(colors[0], "   ")
		}
		cio.print(colors[4], setMaxLength(text[4], tw))
		if len(text[4]) < tw {
			cio.print(colors[0], spaces[:tw-len(text[4])])
		}
		cio.print(colors[0], "  ")
		w := mw
		if e.Flags&incident.FVoice != 0 {
			cio.print(colorSuccessBG, "VOICE")
			cio.print(colors[5], " ")
			w -= 6
		}
		switch {
		case w > 0 && e.Flags&incident.FFollowup != 0:
			cio.print(colors[5], setLength(text[5], w)+" ")
			cio.print(colorImmediate, "★\n")
		case w > 0:
			cio.print(colors[5], setMaxLength(text[5], w)+"\n")
		default:
			cio.print(colors[0], "\n")
		}
	}
}

func compactLine(e *incident.LogEntry, myID string) (columns []string, colors []int) {
	var inbound, outbound, comment bool

	switch e.Status {
	case incident.StatusDraft, incident.StatusQueued, incident.StatusSent:
		outbound = true
	case incident.StatusReceived:
		inbound = true
	case incident.StatusHandEntered:
		inbound = e.ToCall == myID || (e.ToMsgID != "" && e.ToCall == "")
		outbound = e.FromCall == myID || (e.FromMsgID != "" && e.FromCall == "")
		comment = e.FromCall == "" && e.FromMsgID == "" && e.ToCall == "" && e.ToMsgID == ""
	}
	columns = make([]string, 6)
	var color = colorNormal
	switch {
	case e.Flags&incident.FImmediate != 0:
		color = colorImmediate
	case e.Flags&incident.FPriority != 0:
		color = colorPriority
	case e.Flags&incident.FBulletin != 0:
		color = colorBulletin
	case e.Flags&incident.FIsReceipt != 0:
		color = colorReceipt
	}
	colors = slices.Repeat([]int{color}, 6)
	columns[0] = strconv.Itoa(e.Ident)
	switch e.Status {
	case incident.StatusDraft:
		columns[1] = "DRAFT"
		colors[1] = colorAlertBG
	case incident.StatusQueued:
		columns[1] = "READY"
		colors[1] = colorWarningBG
	default:
		columns[1] = e.Time.Format("15:04")
	}
	if e.Flags&incident.FNeedsReceipt != 0 {
		columns[2] = "NO RCPT"
		colors[2] = colorWarningBG
	} else if !outbound {
		if e.FromMsgID != "" {
			columns[2] = e.FromMsgID
		} else if e.FromCall != "" {
			columns[2] = e.FromCall
		}
	}
	columns[3] = e.LocalMsgID
	if inbound {
		if e.Flags&incident.FUnread != 0 {
			columns[4] = "NEW"
			colors[4] = colorAlertBG
		}
	} else if !comment {
		if e.ToMsgID != "" {
			columns[4] = e.ToMsgID
		} else if e.ToCall != "" {
			columns[4] = e.ToCall
		} else {
			columns[4] = "??????"
		}
	}
	subject := e.Subject
	if e.FromMsgID != "" && strings.HasPrefix(subject, e.FromMsgID+"_") {
		subject = strings.TrimPrefix(subject, e.FromMsgID+"_")
	} else if e.LocalMsgID != "" && strings.HasPrefix(subject, e.LocalMsgID+"_") {
		subject = strings.TrimPrefix(subject, e.LocalMsgID+"_")
	} else if e.ToMsgID != "" && strings.HasPrefix(subject, e.ToMsgID+"_") {
		subject = strings.TrimPrefix(subject, e.ToMsgID+"_")
	}
	columns[5] = subject
	return columns, colors
}

func emitLogCSV(entries []*incident.LogEntry) {
	cw := csv.NewWriter(os.Stdout)
	cw.Write([]string{"N", "FLAG", "TIME", "FROM CALL", "FROM MSG #", "TO CALL", "TO MSG #", "MESSAGE", "NF"})
	for _, e := range entries {
		if e.Status == incident.StatusDeleted {
			continue
		}
		var line = make([]string, 9)
		line[0] = strconv.Itoa(e.Ident)
		switch {
		case e.Status == incident.StatusDraft:
			line[1] = "DRAFT"
		case e.Status == incident.StatusQueued:
			line[1] = "READY"
		case e.Flags&incident.FNeedsReceipt != 0:
			line[1] = "NO RCPT"
		case e.Flags&incident.FUnread != 0:
			line[1] = "NEW"
		case e.Flags&incident.FVoice != 0:
			line[1] = "VOICE"
		}
		if !e.Time.IsZero() {
			line[2] = e.Time.Format("2006-01-02T15:04")
		}
		line[3] = e.FromCall
		line[4] = e.FromMsgID
		line[5] = e.ToCall
		line[6] = e.ToMsgID
		line[7] = e.Subject
		if e.Flags&incident.FFollowup != 0 {
			line[8] = "X"
		}
		cw.Write(line)
	}
	cw.Flush()
}
