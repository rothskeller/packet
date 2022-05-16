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
	r.htmlParams(&sb)
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
}

func (r *Report) htmlParams(sb *strings.Builder) {
	sb.WriteString(`<div id="description">`)
	sb.WriteString(html.EscapeString(r.Parameters))
	if r.Modified {
		sb.WriteString("<br><br>NOTE: The practice session expectations were changed after some check-in messages were received.  The earlier check-in messages may have been evaluated with different criteria.")
	}
	sb.WriteString(`</div>`)
}

func (r *Report) htmlStatistics(sb *strings.Builder) {
	sb.WriteString(`<div id="statistics">`)
	fmt.Fprintf(sb, `<div class="stat">Total messages:</div><div class="stat-count">%d</div>`, r.TotalMessages)
	fmt.Fprintf(sb, `<div class="stat">Unique addresses:</div><div class="stat-count">%d</div>`, r.UniqueAddresses)
	if r.UniqueAddresses != 0 {
		fmt.Fprintf(sb, `<div class="stat">Correct messages:</div><div class="stat-count">%d</div><div class="stat-percent">(%d%%)</div>`, r.UniqueAddressesCorrect, r.PercentCorrect)
	}
	fmt.Fprintf(sb, `<div class="stat">Unique call signs:</div><div class="stat-count">%d</div><div class="stat-note">(reported<span class="omitabbr"> to net</span>)</div>`, r.UniqueCallSigns)
	if r.UniqueCallSignsWeek != 0 {
		fmt.Fprintf(sb, `<div class="stat-indent">for the week:</div><div class="stat-count">%d</div>`, r.UniqueCallSignsWeek)
	}
	if len(r.Sources) != 0 {
		sb.WriteString(`<div class="stat-head">Messages from:</div>`)
	}
	for _, source := range r.Sources {
		fmt.Fprintf(sb, `<div class="stat-indent">%s:</div><div class="stat-count">%d</div>`,
			html.EscapeString(source.Name), source.Count)
		if source.SimulatedDown {
			sb.WriteString(`<div class="stat-note">(<span class="omitabbr">simulated </span>down)</div>`)
		}
	}
	sb.WriteString(`</div>`)
}

func (r *Report) htmlMessages(sb *strings.Builder, links string) {
	sb.WriteString(`<div id="messages">`)
	if len(r.CountedMessages) != 0 {
		sb.WriteString(`<div class="heading">The following messages were counted in this report:</div>`)
		for _, m := range r.CountedMessages {
			if links == "" || (links != "" && links == m.FromCallSign) {
				fmt.Fprintf(sb, `<div class="from"><a href="/message/%s">%s</a></div>`,
					m.ID, html.EscapeString(m.FromAddress))
			} else {
				fmt.Fprintf(sb, `<div class="from">%s</div>`, html.EscapeString(m.FromAddress))
			}
			fmt.Fprintf(sb, `<div class="subject">%s</div>`, html.EscapeString(m.Subject))
			if len(m.Problems) != 0 {
				fmt.Fprintf(sb, `<div class="error">%s</div>`,
					html.EscapeString(strings.Join(m.Problems, "\n")))
			}
		}
	}
	if len(r.InvalidMessages) != 0 {
		sb.WriteString(`<div class="heading">The following messages were <span style="color:red">not</span> counted in this report:</div>`)
		for _, m := range r.InvalidMessages {
			if links == "" || (links != "" && links == m.FromCallSign) {
				fmt.Fprintf(sb, `<div class="from"><a href="/message/%s">%s</a></div>`,
					m.ID, html.EscapeString(m.FromAddress))
			} else {
				fmt.Fprintf(sb, `<div class="from">%s</div>`, html.EscapeString(m.FromAddress))
			}
			fmt.Fprintf(sb, `<div class="subject">%s</div>`, html.EscapeString(m.Subject))
			if len(m.Problems) != 0 {
				fmt.Fprintf(sb, `<div class="error">%s</div>`,
					html.EscapeString(strings.Join(m.Problems, "\n")))
			}
		}
	}
	sb.WriteString(`</div>`)
}

func (r *Report) htmlGenInfo(sb *strings.Builder) {
	fmt.Fprintf(sb, `<div id="generation">%s</div>`, html.EscapeString(r.GenerationInfo))
}
