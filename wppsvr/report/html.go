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

var noBreakReplacer = strings.NewReplacer(" ", "&nbsp;", "-", "&#8209;")

func (r *Report) htmlParams(sb *strings.Builder) {
	sb.WriteString(`<div class="block"><div class="block-title">Message Expectations`)
	if r.Modified {
		sb.WriteString(`<sup>*</sup>`)
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
		sb.WriteString(`<div><sup>*</sup>NOTE: The message expectations were changed after some messages were received.  The earlier messages were evaluated with different expectations.</div>`)
	}
	sb.WriteString(`</div>`)
}

func (r *Report) htmlStatistics(sb *strings.Builder) {
	var lines = 4
	sb.WriteString(`<div class="blocks-line"><div class="block"><div class="block-title">Message Counts</div><div class="key-value-note">`)
	fmt.Fprintf(sb, `<div>Total messages:</div><div>%d</div><div></div>`, r.TotalMessages)
	fmt.Fprintf(sb, `<div>Unique addresses:</div><div>%d</div><div></div>`, r.UniqueAddresses)
	if r.UniqueAddresses != 0 {
		fmt.Fprintf(sb, `<div>Correct messages:</div><div>%d</div><div>(%d%%)</div>`, r.UniqueAddressesCorrect, r.PercentCorrect)
		lines++
	}
	fmt.Fprintf(sb, `<div>Unique call signs:</div><div>%d</div><div>(reported)</div>`, r.UniqueCallSigns)
	if r.UniqueCallSignsWeek != 0 {
		fmt.Fprintf(sb, `<div class="indent">for the week:</div><div>%d</div><div></div>`, r.UniqueCallSignsWeek)
		lines++
	}
	sb.WriteString(`</div></div>`)
	if len(r.Sources) != 0 {
		var hasDown bool
		sb.WriteString(`<div class="block"><div class="block-title">Source</div><div class="key-value">`)
		for _, source := range r.Sources {
			var down string
			if source.SimulatedDown {
				down, hasDown = `<sup>*</sup>`, true
			}
			fmt.Fprintf(sb, `<div>%s%s</div><div>%d</div>`, html.EscapeString(source.Name), down, source.Count)
		}
		sb.WriteString(`</div>`)
		if hasDown {
			sb.WriteString(`<div><sup>*</sup>Simulated "down"</div>`)
		}
		sb.WriteString(`</div>`)
		if len(r.Sources) > lines {
			lines = len(r.Sources)
		}
	}
	if len(r.Jurisdictions) != 0 {
		var cols = (len(r.Jurisdictions) + lines - 1) / lines
		var rows = (len(r.Jurisdictions) + cols - 1) / cols
		sb.WriteString(`<div class="block"><div class="block-title">Jurisdiction</div><div class="key-value-columns">`)
		for col := 0; col < len(r.Jurisdictions); col += rows {
			sb.WriteString(`<div class="key-value">`)
			for i := col; i < len(r.Jurisdictions) && i < col+rows; i++ {
				fmt.Fprintf(sb, `<div>%s</div><div>%d</div>`, html.EscapeString(r.Jurisdictions[i].Name), r.Jurisdictions[i].Count)
			}
			sb.WriteString(`</div>`)
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
			multiple, hasMultiple = `<sup>*</sup>`, true
		}
		fmt.Fprintf(sb, `<div>%s%s</div><div>%s</div><div class="%s">%s</div>`,
			m.Source, multiple, m.Jurisdiction, m.Class, m.Problem)
	}
	sb.WriteString(`</div>`)
	if hasMultiple {
		sb.WriteString(`<div><sup>*</sup>multiple messages from this address; only the last one counts</div>`)
	}
	sb.WriteString(`</div>`)
}

func (r *Report) htmlGenInfo(sb *strings.Builder) {
	fmt.Fprintf(sb, `<div id="generation">%s</div>`, html.EscapeString(r.GenerationInfo))
}

/*
   <div class="block">
     <div class="block-title">Messages</div>
     <div id="messages">
       <div><a href="...">W6</a></div>
       <div><a href="...">BG</a></div>
       <div>W1XSC</div>
       <div>CUP</div>
       <div class="ok"></div>
       <div><a href="...">AK6</a></div>
       <div><a href="...">BY</a></div>
       <div>W1XSC</div>
       <div>SJC</div>
       <div class="ok"></div>
       <div><a href="...">KK6</a></div>
       <div><a href="...">EBL</a></div>
       <div>W1XSC</div>
       <div>LGT</div>
       <div class="warning">invalid Practice subject format</div>
       <div><a href="...">W6</a></div>
       <div><a href="...">ESL</a></div>
       <div>W1XSC*</div>
       <div>SJC</div>
       <div class="error">wrong something</div>
       <div><a href="...">AJ6</a></div>
       <div><a href="...">LG</a></div>
       <div>W1XSC</div>
       <div>SJC</div>
       <div class="invalid">message to incorrect BBS</div>
     </div>
   </div>
   <div id="generation">This report was generated on Tuesday, May 24, 2022 at 08:32 by wppsvr version (devel).</div>

*/
