package report

import (
	"fmt"
	"strings"

	"github.com/rothskeller/packet/wppsvr/english"
)

// RenderPlainText renders a report in plain text format.
func (r *Report) RenderPlainText() string {
	var sb strings.Builder

	r.plainTextTitle(&sb)
	r.plainTextParams(&sb)
	r.plainTextStatistics(&sb)
	r.plainTextMessages(&sb)
	r.plainTextGenInfo(&sb)
	return sb.String()
}

func (r *Report) plainTextTitle(sb *strings.Builder) {
	var (
		line1  string
		line2  string
		maxlen int
		pad1   int
		pad2   int
	)
	line1 = "SCCo ARES/RACES Packet Practice Report"
	if r.Preliminary {
		line1 = "***PRELIMINARY*** " + line1
	}
	maxlen = len(line1)
	line2 = fmt.Sprintf("for %s on %s", r.SessionName, r.SessionDate)
	if len(line2) > maxlen {
		maxlen = len(line2)
	}
	pad1 = (maxlen - len(line1)) / 2
	pad2 = (maxlen - len(line2)) / 2
	fmt.Fprintf(sb, "==== %*s%-*s ====\n==== %*s%-*s ====\n\n",
		pad1, "", maxlen-pad1, line1, pad2, "", maxlen-pad2, line2)

}

func (r *Report) plainTextParams(sb *strings.Builder) {
	var wr = english.NewWrapper(sb)

	wr.WriteString(r.Parameters)
	if r.Modified {
		wr.WriteString("\n\nNOTE: The practice session expectations were changed after some check-in messages were received.  The earlier check-in messages may have been evaluated with different criteria.")
	}
	wr.WriteString("\n\n")
	wr.Close()
}

func (r *Report) plainTextStatistics(sb *strings.Builder) {
	var countWidth int

	countWidth = numberWidth(r.TotalMessages)
	if cw := numberWidth(r.UniqueCallSignsWeek); cw > countWidth {
		countWidth = cw
	}
	fmt.Fprintf(sb, "Total messages:     %*d\n", countWidth, r.TotalMessages)
	fmt.Fprintf(sb, "Unique addresses:   %*d\n", countWidth, r.UniqueAddresses)
	if r.UniqueAddresses != 0 {
		fmt.Fprintf(sb, "Correct messages:   %*d  (%d%%)\n", countWidth, r.UniqueAddressesCorrect, r.PercentCorrect)
	}
	fmt.Fprintf(sb, "Unique call signs:  %*d  [report this count to the net]\n", countWidth, r.UniqueCallSigns)
	if r.UniqueCallSignsWeek != 0 {
		fmt.Fprintf(sb, "  for the week:     %*d\n", countWidth, r.UniqueCallSignsWeek)
	}
	if len(r.Sources) != 0 {
		sb.WriteString("Messages from:\n")
	}
	for _, source := range r.Sources {
		if source.SimulatedDown {
			fmt.Fprintf(sb, "  %s:%*s%*d  (simulated down)\n",
				source.Name, 17-len(source.Name), "", countWidth, source.Count)
		} else {
			fmt.Fprintf(sb, "  %s:%*s%*d\n", source.Name, 17-len(source.Name), "", countWidth, source.Count)
		}
	}
	sb.WriteByte('\n')
}
func numberWidth(i int) int {
	if i > 99 {
		return 3 // none of the numbers we deal with are over 999
	}
	if i > 9 {
		return 2
	}
	return 1 // none of them are negative, either
}

func (r *Report) plainTextMessages(sb *strings.Builder) {
	if len(r.CountedMessages) != 0 {
		sb.WriteString("---- The following messages were counted in this report: ----\n")
		for _, m := range r.CountedMessages {
			fmt.Fprintf(sb, "%-30s %s\n", m.FromAddress, m.Subject)
			for _, p := range m.Problems {
				fmt.Fprintf(sb, "  ^ %s\n", p)
			}
		}
		sb.WriteByte('\n')
	}
	if len(r.InvalidMessages) != 0 {
		sb.WriteString("---- The following messages were not counted in this report: ----\n")
		for _, m := range r.InvalidMessages {
			fmt.Fprintf(sb, "%-30s %s\n", m.FromAddress, m.Subject)
			for _, p := range m.Problems {
				fmt.Fprintf(sb, "  ^ %s\n", p)
			}
		}
		sb.WriteByte('\n')
	}
}

func (r *Report) plainTextGenInfo(sb *strings.Builder) {
	var wr = english.NewWrapper(sb)
	wr.WriteString(r.GenerationInfo)
	wr.WriteString("\n")
	wr.Close()
}
