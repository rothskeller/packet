package report

import (
	_ "embed" // -
	"fmt"
	"html"
	"strings"
)

var contentMarker = "@@CONTENT@@"

//go:embed "html.html"
var reportHTML string

// RenderHTML renders a report in HTML format.  If links is set to a call sign,
// only messages from that call sign have embedded links.  If links is an empty
// string, all messages have embedded links.
func (r *Report) RenderHTML(links string) string {
	var sb strings.Builder

	content := strings.Index(reportHTML, contentMarker)
	sb.WriteString(reportHTML[:content])
	r.htmlTitle(&sb)
	r.htmlExpectsResults(&sb)
	r.htmlStatistics(&sb)
	r.htmlMessages(&sb, links)
	r.htmlGenInfo(&sb)
	sb.WriteString(reportHTML[content+len(contentMarker):])
	return sb.String()
}

func (r *Report) htmlTitle(sb *strings.Builder) {
	fmt.Fprintf(sb, `<div id="date">%s — %s</div>`, r.SessionName, r.SessionDate)
	if r.Preliminary {
		fmt.Fprintf(sb, `<div id="preliminary">PRELIMINARY</div>`)
	}
	if r.UniqueCallSigns != 0 {
		fmt.Fprintf(sb, `<div id="unique">%d Unique Call Signs</div>`, r.UniqueCallSigns)
		if r.UniqueCallSignsWeek != 0 {
			fmt.Fprintf(sb, `<div id="unique-week">%d for the week</div>`, r.UniqueCallSignsWeek)
		}
	}
}

var noBreakReplacer = strings.NewReplacer(" ", "&nbsp;", "-", "&#8209;")

func (r *Report) htmlExpectsResults(sb *strings.Builder) {
	sb.WriteString(`<div class="blocks-line"><div class="block"><div class="block-title">Expectations`)
	if r.Modified {
		sb.WriteString(`*`)
	}
	sb.WriteString(`</div><div class="key-text"><div>Message type:</div><div>`)
	sb.WriteString(r.MessageTypes)
	sb.WriteString(`</div><div>Sent to:</div><div>`)
	sb.WriteString(r.SentTo)
	sb.WriteString(`</div><div>Sent between:</div><div style="white-space:normal">`)
	sb.WriteString(noBreakReplacer.Replace(r.SentAfter))
	sb.WriteString(`&nbsp;and `)
	sb.WriteString(noBreakReplacer.Replace(r.SentBefore))
	sb.WriteString(`</div>`)
	if r.NotSentFrom != "" {
		sb.WriteString(`<div>Not sent from:</div><div>`)
		sb.WriteString(r.NotSentFrom)
		sb.WriteString(`</div>`)
	}
	sb.WriteString(`</div>`)
	if r.Modified {
		sb.WriteString(`<div>*modified during session</div>`)
	}
	sb.WriteString(`</div>`)
	sb.WriteString(`<div class="block"><div class="block-title">Results</div><div class="key-value">`)
	if r.OKCount+r.WarningCount+r.ErrorCount+r.InvalidCount+r.ReplacedCount+r.DroppedCount != 0 {
		if r.OKCount != 0 {
			fmt.Fprintf(sb, `<div>OK:</div><div>%d</div>`, r.OKCount)
		}
		if r.WarningCount != 0 {
			fmt.Fprintf(sb, `<div>WARNING:</div><div>%d</div>`, r.WarningCount)
		}
		if r.ErrorCount != 0 {
			fmt.Fprintf(sb, `<div>ERROR:</div><div>%d</div>`, r.ErrorCount)
		}
		if r.InvalidCount != 0 {
			fmt.Fprintf(sb, `<div class="gray">NOT COUNTED:</div><div class="gray">%d</div>`, r.InvalidCount)
		}
		if r.ReplacedCount != 0 {
			fmt.Fprintf(sb, `<div class="gray">Duplicate:</div><div class="gray">%d</div>`, r.ReplacedCount)
		}
		if r.DroppedCount != 0 {
			fmt.Fprintf(sb, `<div class="gray">Delivery rcpt:</div><div class="gray">%d</div>`, r.DroppedCount)
		}
	} else {
		sb.WriteString(`<div>Messages:</div><div>0</div>`)
	}
	sb.WriteString(`</div></div></div>`)
}

func (r *Report) htmlStatistics(sb *strings.Builder) {
	if len(r.Sources) == 0 && len(r.Jurisdictions) == 0 && len(r.MTypeCounts) == 0 {
		return
	}
	sb.WriteString(`<div class="blocks-line">`)
	if len(r.Sources) != 0 {
		var hasDown bool
		sb.WriteString(`<div class="block"><div class="block-title">Sources</div><div class="key-value">`)
		for _, source := range r.Sources {
			var down string
			if source.SimulatedDown {
				down, hasDown = `*`, true
			}
			fmt.Fprintf(sb, `<div>%s%s</div><div>%d</div>`, html.EscapeString(source.Name), down, source.Count)
		}
		sb.WriteString(`</div>`)
		if hasDown {
			sb.WriteString(`<div>*Simulated "down"</div>`)
		}
		sb.WriteString(`</div>`)
	}
	if len(r.Jurisdictions) != 0 {
		var cols = (len(r.Jurisdictions) + 5) / 6
		var rows = (len(r.Jurisdictions) + cols - 1) / cols
		sb.WriteString(`<div class="block"><div class="block-title">Jurisdictions</div><div class="key-value-columns">`)
		for col := 0; col < len(r.Jurisdictions); col += rows {
			sb.WriteString(`<div class="key-value">`)
			for i := col; i < len(r.Jurisdictions) && i < col+rows; i++ {
				fmt.Fprintf(sb, `<div>%s</div><div>%d</div>`, html.EscapeString(r.Jurisdictions[i].Name), r.Jurisdictions[i].Count)
			}
			sb.WriteString(`</div>`)
		}
		sb.WriteString(`</div></div>`)
	}
	if len(r.MTypeCounts) != 0 {
		sb.WriteString(`<div class="block"><div class="block-title">Types</div><div class="key-value">`)
		for _, mtype := range r.MTypeCounts {
			fmt.Fprintf(sb, `<div>%s</div><div>%d</div>`, html.EscapeString(mtype.Name), mtype.Count)
		}
		sb.WriteString(`</div></div>`)
	}
	sb.WriteString(`</div>`)
}

func (r *Report) htmlMessages(sb *strings.Builder, links string) {
	var hasMultiple bool
	sb.WriteString(`<div class="block"><div class="block-title">Messages</div><div id="messages">`)
	for _, m := range r.Messages {
		var multiple string
		if links == "" || (links != "" && links == m.FromCallSign) {
			fmt.Fprintf(sb, `<div><a href="/message?id=%s">%s</a></div><div><a href="/message?id=%s">%s</a></div>`,
				m.ID, html.EscapeString(m.Prefix), m.ID, html.EscapeString(m.Suffix))
		} else {
			fmt.Fprintf(sb, `<div>%s</div><div>%s</div>`, html.EscapeString(m.Prefix), html.EscapeString(m.Suffix))
		}
		if m.Multiple {
			multiple, hasMultiple = `*`, true
		}
		fmt.Fprintf(sb, `<div>%s%s</div><div>%s</div><div class="%s">%s</div>`,
			m.Source, multiple, m.Jurisdiction, m.Class, m.Problem)
	}
	sb.WriteString(`</div>`)
	if hasMultiple {
		sb.WriteString(`<div>*multiple messages from this address; only the last one counts</div>`)
	}
	sb.WriteString(`</div>`)
}

func (r *Report) htmlGenInfo(sb *strings.Builder) {
	fmt.Fprintf(sb, `<div id="generation">%s</div>`, html.EscapeString(r.GenerationInfo))
}
