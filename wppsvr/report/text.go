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
	r.plainTextStatistics(&sb)
	r.plainTextMessages(&sb)
	r.plainTextGenInfo(&sb)
	return sb.String()
}

func (r *Report) plainTextTitle(sb *strings.Builder) {
	var (
		line string
	)
	sb.WriteString("--------------- SCCo ARES/RACES Packet Practice Report ----------------\n")
	line = fmt.Sprintf("for %s on %s", r.SessionName, r.SessionDate)
	fmt.Fprintf(sb, "%*s%s\n", 36-len(line)/2, "", line)
	if r.Preliminary {
		sb.WriteString("                         *** PRELIMINARY ***\n")
	}
	if r.UniqueCallSignsWeek != 0 {
		line = fmt.Sprintf("%d unique call signs (%d for the week)", r.UniqueCallSigns, r.UniqueCallSignsWeek)
	} else if r.UniqueCallSigns != 0 {
		line = fmt.Sprintf("%d unique call signs", r.UniqueCallSigns)
	} else {
		line = ""
	}
	if line != "" {
		fmt.Fprintf(sb, "\n%*s%s\n", 36-len(line)/2, "", line)
	}
	sb.WriteString("\n\n")
}

func (r *Report) plainTextExpectsResults(sb *strings.Builder) {
	var (
		lines []string
		col1  []string
	)
	switch len(r.MessageTypes) {
	case 0:
		break
	case 1:
		lines = append(lines, fmt.Sprintf("Message type:  %s", r.MessageTypes[0]))
	case 2:
		lines = append(lines, fmt.Sprintf("Message type:  %s or", r.MessageTypes[0]))
		lines = append(lines, fmt.Sprintf("               %s", r.MessageTypes[1]))
	default:
		lines = append(lines, fmt.Sprintf("Message type:  %s,", r.MessageTypes[0]))
		for i := 1; i < len(r.MessageTypes)-2; i++ {
			lines = append(lines, fmt.Sprintf("               %s,", r.MessageTypes[i]))
		}
		lines = append(lines, fmt.Sprintf("               %s, or", r.MessageTypes[len(r.MessageTypes)-2]))
		lines = append(lines, fmt.Sprintf("               %s", r.MessageTypes[len(r.MessageTypes)-1]))
	}
	lines = append(lines, fmt.Sprintf("Sent to:       %s", r.SentTo))
	lines = append(lines, fmt.Sprintf("Sent between:  %s and", r.SentAfter))
	lines = append(lines, fmt.Sprintf("               %s", r.SentBefore))
	if r.NotSentFrom != "" {
		lines = append(lines, fmt.Sprintf("Not sent from: %s", r.NotSentFrom))
	}
	if r.Modified {
		lines = append(lines, "(*) modified during session")
	}
	if r.Modified {
		lines = addPlainTextHeading(lines, "EXPECTATIONS(*)")
	} else {
		lines = addPlainTextHeading(lines, "EXPECTATIONS")
	}
	if r.OKCount+r.WarningCount+r.ErrorCount+r.InvalidCount+r.ReplacedCount+r.DroppedCount != 0 {
		var col2 []string
		if r.OKCount != 0 {
			col1 = append(col1, "OK")
			col2 = append(col2, strconv.Itoa(r.OKCount))
		}
		if r.WarningCount != 0 {
			col1 = append(col1, "WARNING")
			col2 = append(col2, strconv.Itoa(r.WarningCount))
		}
		if r.ErrorCount != 0 {
			col1 = append(col1, "ERROR")
			col2 = append(col2, strconv.Itoa(r.ErrorCount))
		}
		if r.InvalidCount != 0 {
			col1 = append(col1, "NOT COUNTED")
			col2 = append(col2, strconv.Itoa(r.InvalidCount))
		}
		if r.ReplacedCount != 0 {
			col1 = append(col1, "Duplicate")
			col2 = append(col2, strconv.Itoa(r.ReplacedCount))
		}
		if r.DroppedCount != 0 {
			col1 = append(col1, "Deliv. rcpt.")
			col2 = append(col2, strconv.Itoa(r.DroppedCount))
		}
		rightAlign(col2)
		col1 = sideBySide(col1, col2, 2)
	} else {
		col1 = append(col1, "Messages  0")
	}
	col1 = addPlainTextHeading(col1, "RESULTS")
	lines = sideBySide(lines, col1, 6)
	for _, line := range lines {
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	sb.WriteByte('\n')
}

func (r *Report) plainTextStatistics(sb *strings.Builder) {
	var lines []string

	if len(r.Sources) == 0 && len(r.Jurisdictions) == 0 && len(r.MTypeCounts) == 0 {
		return
	}
	if len(r.Sources) != 0 {
		var col1, col2 []string
		var hasDown bool

		for _, source := range r.Sources {
			var down string
			if source.SimulatedDown {
				down, hasDown = `*`, true
			}
			col1 = append(col1, source.Name+down)
			col2 = append(col2, strconv.Itoa(source.Count))
		}
		rightAlign(col2)
		lines = sideBySide(col1, col2, 2)
		if hasDown {
			lines = append(lines, `* simulated outage`)
		}
		lines = addPlainTextHeading(lines, "SOURCES")
	}
	if len(r.Jurisdictions) != 0 {
		var jlines []string
		var cols = (len(r.Jurisdictions) + 5) / 6
		var rows = (len(r.Jurisdictions) + cols - 1) / cols
		for col := 0; col < len(r.Jurisdictions); col += rows {
			var col1, col2 []string
			for i := col; i < len(r.Jurisdictions) && i < col+rows; i++ {
				col1 = append(col1, r.Jurisdictions[i].Name)
				col2 = append(col2, strconv.Itoa(r.Jurisdictions[i].Count))
			}
			rightAlign(col2)
			col1 = sideBySide(col1, col2, 2)
			if col == 0 {
				jlines = col1
			} else {
				jlines = sideBySide(jlines, col1, 4)
			}
		}
		jlines = addPlainTextHeading(jlines, "JURISDICTIONS")
		if lines == nil {
			lines = jlines
		} else {
			lines = sideBySide(lines, jlines, 6)
		}
	}
	if len(r.MTypeCounts) != 0 {
		var col1, col2 []string
		for _, mtype := range r.MTypeCounts {
			col1 = append(col1, mtype.Name)
			col2 = append(col2, strconv.Itoa(mtype.Count))
		}
		rightAlign(col2)
		col1 = sideBySide(col1, col2, 2)
		col1 = addPlainTextHeading(col1, "TYPES")
		if lines == nil {
			lines = col1
		} else {
			lines = sideBySide(lines, col1, 6)
		}
	}
	for _, line := range lines {
		sb.WriteString(line)
		sb.WriteByte('\n')
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

	sb.WriteString("MESSAGES----------------------------------------------------------------\n")
	for _, m := range r.Messages {
		var multiple string
		col1 = append(col1, m.Prefix)
		col2 = append(col2, m.Suffix)
		if m.Multiple {
			multiple, hasMultiple = `*`, true
		}
		col3 = append(col3, m.Source+multiple)
		col4 = append(col4, m.Jurisdiction)
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
	sb.WriteString("\n\n")
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
