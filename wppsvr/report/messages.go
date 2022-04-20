package report

import (
	"fmt"
	"strings"

	"steve.rothskeller.net/packet/wppsvr/analyze"
	"steve.rothskeller.net/packet/wppsvr/store"
)

// reportMessages generates the lists of valid and invalid check-in messages
// that appear in the report.
func reportMessages(sb *strings.Builder, session *store.Session, messages []*store.Message) {
	var (
		invalid   []*store.Message
		seenValid bool
	)
	messages = removeReplaced(messages)
	for _, m := range messages {
		if !m.Valid {
			invalid = append(invalid, m)
			continue
		}
		if !seenValid {
			sb.WriteString("---- The following messages were counted in this report: ----\n")
			seenValid = true
		}
		fmt.Fprintf(sb, "%-30s %s\n", m.FromAddress, m.Subject)
		for _, p := range m.Problems {
			fmt.Fprintf(sb, "  ^ %s\n", analyze.ProblemLabel[p])
		}
	}
	if seenValid {
		sb.WriteByte('\n')
	}
	if len(invalid) != 0 {
		sb.WriteString("---- The following messages were not counted in this report: ----\n")
		for _, m := range invalid {
			if m.FromAddress == "" && m.Subject == "" {
				fmt.Fprintf(sb, "[unparseable message with hash %08x]\n", m.Hash)
			} else {
				fmt.Fprintf(sb, "%-30s %s\n", m.FromAddress, m.Subject)
			}
			for _, p := range m.Problems {
				fmt.Fprintf(sb, "  ^ %s\n", analyze.ProblemLabel[p])
			}
		}
		sb.WriteByte('\n')
	}
}

// removeReplaced removes all but the last message from each address.  If more
// than one message is found from a given address, a MultipleMessagesFromAddress
// problem code is added to the one that is kept.
func removeReplaced(messages []*store.Message) (out []*store.Message) {
	var (
		msgidx    int
		outidx    int
		addresses = make(map[string]*store.Message)
	)
	out = make([]*store.Message, len(messages))
	outidx = len(messages)
	for msgidx = len(messages) - 1; msgidx >= 0; msgidx-- {
		m := messages[msgidx]
		if m.FromAddress == "" {
			outidx--
			out[outidx] = m
			continue
		}
		if keeper := addresses[m.FromAddress]; keeper == nil {
			outidx--
			out[outidx] = m
			addresses[m.FromAddress] = m
		} else if len(keeper.Problems) == 0 ||
			keeper.Problems[len(keeper.Problems)-1] != analyze.ProblemMultipleMessagesFromAddress {
			keeper.Problems = append(keeper.Problems, analyze.ProblemMultipleMessagesFromAddress)
		}
	}
	return out[outidx:]
}
