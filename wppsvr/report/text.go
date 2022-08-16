package report

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rothskeller/packet/wppsvr/english"
)

const dashes = "--------------------------------------------------------------------------------"
const spaces = "                                                                                "

// RenderPlainText renders a report in plain text format.
func (r *Report) RenderPlainText() string {
	var sb strings.Builder

	r.plainTextTitle(&sb)
	r.plainTextExpectsResults(&sb)
	r.plainTextMessages(&sb)
	r.plainTextStatistics(&sb)
	r.plainTextGenInfo(&sb)
	return sb.String()
}

func (r *Report) plainTextTitle(sb *strings.Builder) {
	fmt.Fprintf(sb, "==== SCCo ARES/RACES Packet Practice Report\n==== for %s on %s", r.SessionName, r.SessionDate)
	if r.Preliminary {
		sb.WriteString(" (PRELIMINARY)")
	}
	if r.UniqueCallSignsWeek != 0 {
		fmt.Fprintf(sb, "\n\n%d unique call signs (%d for the week)\n\n", r.UniqueCallSigns, r.UniqueCallSignsWeek)
	} else if r.UniqueCallSigns != 0 {
		fmt.Fprintf(sb, "\n\n%d unique call signs\n\n", r.UniqueCallSigns)
	} else {
		sb.WriteString("\n\n")
	}
}

func (r *Report) plainTextExpectsResults(sb *strings.Builder) {
	var wrapper = english.NewWrapper(sb)
	fmt.Fprintf(wrapper, "EXPECTATIONS:  %s sent to %s between %s and %s",
		english.Conjoin(r.MessageTypes, "or"), r.SentTo, r.SentAfter, r.SentBefore)
	if r.NotSentFrom != "" {
		fmt.Fprintf(wrapper, "; not sent from %s", r.NotSentFrom)
	}
	wrapper.WriteString(".")
	if r.Modified {
		wrapper.WriteString("  Expectations were modified during session; some early messages may have been evaluated against different expectations.")
	}
	wrapper.WriteString("\n\n")
	wrapper.Close()

	sb.WriteString("---- RESULTS\n")
	if r.OKCount+r.WarningCount+r.ErrorCount+r.InvalidCount+r.ReplacedCount+r.DroppedCount != 0 {
		var lines, col1, col2 []string
		if r.OKCount != 0 {
			col1 = append(col1, strconv.Itoa(r.OKCount))
			col2 = append(col2, "OK")
		}
		if r.WarningCount != 0 {
			col1 = append(col1, strconv.Itoa(r.WarningCount))
			col2 = append(col2, "WARNING")
		}
		if r.ErrorCount != 0 {
			col1 = append(col1, strconv.Itoa(r.ErrorCount))
			col2 = append(col2, "ERROR")
		}
		if r.InvalidCount != 0 {
			col1 = append(col1, strconv.Itoa(r.InvalidCount))
			col2 = append(col2, "NOT COUNTED")
		}
		if r.ReplacedCount != 0 {
			col1 = append(col1, strconv.Itoa(r.ReplacedCount))
			col2 = append(col2, "Duplicate")
		}
		if r.DroppedCount != 0 {
			col1 = append(col1, strconv.Itoa(r.DroppedCount))
			col2 = append(col2, "Delivery receipt")
		}
		rightAlign(col1)
		lines = sideBySide(col1, col2, 2)
		for _, line := range lines {
			sb.WriteString(line)
			sb.WriteByte('\n')
		}
	} else {
		sb.WriteString("0  Messages")
	}
	sb.WriteByte('\n')
}

var plainTextClassLabels = map[string]string{
	"ok":      "OK",
	"warning": "WARNING: ",
	"error":   "ERROR: ",
	"invalid": "NOT COUNTED: ",
}

func (r *Report) plainTextMessages(sb *strings.Builder) {
	var col1, col2, col3, col4, col5 []string
	var hasMultiple bool

	if len(r.Messages) == 0 {
		return
	}
	sb.WriteString("---- MESSAGES\n")
	for _, m := range r.Messages {
		var multiple string
		col1 = append(col1, m.Prefix)
		col2 = append(col2, m.Suffix)
		if m.Multiple {
			multiple, hasMultiple = `*`, true
		}
		col3 = append(col3, "@"+m.Source+multiple)
		col4 = append(col4, "("+m.Jurisdiction+")")
		col5 = append(col5, plainTextClassLabels[m.Class]+m.Problem)
	}
	rightAlign(col1)
	col1 = sideBySide(col1, col2, 0)
	col1 = sideBySide(col1, col3, 2)
	col1 = sideBySide(col1, col4, 2)
	col1 = sideBySide(col1, col5, 2)
	for _, line := range col1 {
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	if hasMultiple {
		sb.WriteString("* multiple messages from this address; only the last one counts\n")
	}
	sb.WriteString("\n")
}

func (r *Report) plainTextStatistics(sb *strings.Builder) {
	if len(r.Sources) != 0 {
		var lines, col1, col2 []string

		for _, source := range r.Sources {
			var down string
			if source.SimulatedDown {
				down = " (simulated outage)"
			}
			col1 = append(col1, strconv.Itoa(source.Count))
			col2 = append(col2, source.Name+down)
		}
		rightAlign(col1)
		lines = sideBySide(col1, col2, 2)
		sb.WriteString("---- SENT FROM\n")
		for _, line := range lines {
			sb.WriteString(line)
			sb.WriteByte('\n')
		}
		sb.WriteString("\n")
	}
	if len(r.Jurisdictions) != 0 {
		var lines, col1, col2 []string

		for _, juris := range r.Jurisdictions {
			col1 = append(col1, strconv.Itoa(juris.Count))
			col2 = append(col2, juris.Name)
		}
		rightAlign(col1)
		lines = sideBySide(col1, col2, 2)
		sb.WriteString("---- JURISDICTION\n")
		for _, line := range lines {
			sb.WriteString(line)
			sb.WriteByte('\n')
		}
		sb.WriteString("\n")
	}
	if len(r.MTypeCounts) != 0 {
		var lines, col1, col2 []string

		for _, mtype := range r.MTypeCounts {
			col1 = append(col1, strconv.Itoa(mtype.Count))
			col2 = append(col2, mtype.Name)
		}
		rightAlign(col1)
		lines = sideBySide(col1, col2, 2)
		sb.WriteString("---- MESSAGE TYPE\n")
		for _, line := range lines {
			sb.WriteString(line)
			sb.WriteByte('\n')
		}
		sb.WriteString("\n")
	}
}

func (r *Report) plainTextGenInfo(sb *strings.Builder) {
	var wr = english.NewWrapper(sb)
	wr.WriteString(r.GenerationInfo)
	wr.WriteString("\n")
	wr.Close()
}

func leftAlign(ss []string) {
	var maxlen = maxlength(ss)
	for i, s := range ss {
		if len(s) < maxlen {
			ss[i] = s + spaces[:maxlen-len(s)]
		}
	}
}

func rightAlign(ss []string) {
	var maxlen = maxlength(ss)
	for i, s := range ss {
		if len(s) < maxlen {
			ss[i] = spaces[:maxlen-len(s)] + s
		}
	}
}

func sideBySide(block1, block2 []string, gap int) (combined []string) {
	var (
		maxlen  int
		linenum int
	)
	maxlen = maxlength(block1)
	for linenum = 0; linenum < len(block1) && linenum < len(block2); linenum++ {
		combined = append(combined, fmt.Sprintf("%-*s%s", maxlen+gap, block1[linenum], block2[linenum]))
	}
	for ; linenum < len(block1); linenum++ {
		combined = append(combined, block1[linenum])
	}
	for ; linenum < len(block2); linenum++ {
		combined = append(combined, fmt.Sprintf("%-*s%s", maxlen+gap, "", block2[linenum]))
	}
	return combined
}

func addPlainTextHeading(ss []string, head string) []string {
	var maxlen = maxlength(ss)
	if maxlen > len(head) {
		head = head + dashes[:maxlen-len(head)]
	}
	return append([]string{head}, ss...)
}

func maxlength(ss []string) (maxlen int) {
	for _, s := range ss {
		if len(s) > maxlen {
			maxlen = len(s)
		}
	}
	return maxlen
}
