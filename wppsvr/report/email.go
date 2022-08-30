package report

import (
	"fmt"
	"html"
	"io"
	"mime/quotedprintable"
	"strings"

	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/english"
)

// RenderEmail renders the receiver report in a form suitable for emailing.  The
// email is a multipart/alternative email with plain text and HTML variants.
// The parameter is the To: line of the email.
func (r *Report) RenderEmail(to string) string {
	var sb strings.Builder
	var qw = quotedprintable.NewWriter(&sb)

	fmt.Fprintf(&sb, "From: %s\r\nTo: %s\r\n", config.Get().SMTP.From, to)
	sb.WriteString("Subject: SCCo Packet Practice Report\r\nContent-Type: multipart/alternative; boundary=\"BOUNDARY\"\r\n\r\n\r\n--BOUNDARY\r\nContent-Type: text/plain\r\n\r\n")
	var textReport = r.RenderPlainText()
	sb.WriteString(strings.ReplaceAll(textReport, "\n", "\r\n"))
	sb.WriteString("\r\n--BOUNDARY\r\nContent-Type: text/html; charset=utf-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n")
	r.renderEmail(qw)
	sb.WriteString("\r\n--BOUNDARY--\r\n")
	return sb.String()
}

func (r *Report) renderEmail(w *quotedprintable.Writer) {
	io.WriteString(w, `<!DOCTYPE html><html><head><meta charset="utf-8"></head><body><div style="width:800px;font-size:16px;line-height:1.25;padding:16px"><div style="color:#444;font-size:20px;font-weight:bold;text-align:center">Santa Clara County ARES<sup>®</sup>/RACES</div><div style="color:#D00;font-size:28px;font-weight:bold;text-align:center">Weekly Packet Practice</div>`)
	r.emailTitle(w)
	r.emailExpectsResults(w)
	r.emailStatistics(w)
	r.emailMessages(w)
	r.emailGenInfo(w)
	io.WriteString(w, "</div></body></html>\n")
}

func (r *Report) emailTitle(w *quotedprintable.Writer) {
	fmt.Fprintf(w, `<div style="font-size:28px;font-weight:bold;text-align:center">%s — %s</div>`, r.SessionName, r.SessionDate)
	if r.Preliminary {
		io.WriteString(w, `<div style="font-size:24px;font-weight:bold;text-align:center;color:#d00">PRELIMINARY</div>`)
	}
	if r.UniqueCallSigns != 0 {
		fmt.Fprintf(w, `<div style="margin-top:16px;font-size:24px;font-weight:bold;text-align:center">%d Unique Call Signs</div>`, r.UniqueCallSigns)
		if r.UniqueCallSignsWeek != 0 {
			fmt.Fprintf(w, `<div style="font-size:20px;font-weight:bold;color:#888;text-align:center">%d for the week</div>`, r.UniqueCallSignsWeek)
		}
	}
}

func (r *Report) emailExpectsResults(w *quotedprintable.Writer) {
	io.WriteString(w, `<table cellspacing="0" cellpadding="0" style="margin-top:24px"><tr><td style="vertical-align:top"><div style="max-width:480px;margin-bottom:24px"><div style="font-size:20px;font-weight:bold;color:#444">Expectations`)
	if r.Modified {
		io.WriteString(w, `*`)
	}
	io.WriteString(w, `</div><table cellspacing="0" cellpadding="0"><tr><td style="white-space:nowrap;vertical-align:top;color:#666">Message type:</td><td style="padding-left:16px">`)
	io.WriteString(w, html.EscapeString(english.Conjoin(r.MessageTypes, "or")))
	io.WriteString(w, `</td></tr><tr><td style="white-space:nowrap;padding-top:2px;color:#666">Sent to:</td><td style="padding:2px 0 0 16px">`)
	io.WriteString(w, html.EscapeString(r.SentTo))
	io.WriteString(w, `</td></tr><tr><td style="white-space:nowrap;vertical-align:top;padding-top:2px;color:#666">Sent between:</td><td style="padding:2px 0 0 16px">`)
	io.WriteString(w, noBreakReplacer.Replace(r.SentAfter))
	io.WriteString(w, `&nbsp;and `)
	io.WriteString(w, noBreakReplacer.Replace(r.SentBefore))
	io.WriteString(w, `</td></tr>`)
	if r.NotSentFrom != "" {
		io.WriteString(w, `<tr><td style="white-space:nowrap;padding-top:2px;color:#666">Not sent from:</td><td style="padding:2px 0 0 16px">`)
		io.WriteString(w, r.NotSentFrom)
		io.WriteString(w, `</td></tr>`)
	}
	io.WriteString(w, `</table>`)
	if r.Modified {
		io.WriteString(w, `<div>*modified during session</div>`)
	}
	io.WriteString(w, `</div></td><td style="padding-left:32px;vertical-align:top"><div style="max-width:640px;margin-bottom:24px"><div style="font-size:20px;font-weight:bold;color:#444">Results</div><table cellspacing="0" cellpadding="0">`)
	if r.OKCount+r.WarningCount+r.ErrorCount+r.InvalidCount+r.ReplacedCount+r.DroppedCount != 0 {
		if r.OKCount != 0 {
			fmt.Fprintf(w, `<tr><td style="padding-top:2px;color:#666">OK</td><td style="padding:2px 0 0 16px;text-align:right">%d</td></tr>`, r.OKCount)
		}
		if r.WarningCount != 0 {
			fmt.Fprintf(w, `<tr><td style="padding-top:2px;color:#666">WARNING</td><td style="padding:2px 0 0 16px;text-align:right">%d</td></tr>`, r.WarningCount)
		}
		if r.ErrorCount != 0 {
			fmt.Fprintf(w, `<tr><td style="padding-top:2px;color:#666">ERROR</td><td style="padding:2px 0 0 16px;text-align:right">%d</td></tr>`, r.ErrorCount)
		}
		if r.InvalidCount != 0 {
			fmt.Fprintf(w, `<tr><td style="padding-top:2px;color:#666">NOT COUNTED</td><td style="padding:2px 0 0 16px;color:#888;text-align:right">%d</td></tr>`, r.InvalidCount)
		}
		if r.ReplacedCount != 0 {
			fmt.Fprintf(w, `<tr><td style="padding-top:2px;color:#666">Duplicate</td><td style="padding:2px 0 0 16px;color:#888;text-align:right">%d</td></tr>`, r.ReplacedCount)
		}
		if r.DroppedCount != 0 {
			fmt.Fprintf(w, `<tr><td style="padding-top:2px;color:#666">Deliv. rcpt.</td><td style="padding:2px 0 0 16px;color:#888;text-align:right">%d</td></tr>`, r.DroppedCount)
		}
	} else {
		io.WriteString(w, `<tr><td style="padding-top:2px;color:#666">Messages</td><td style="padding:2px 0 0 16px">0</td></tr>`)
	}
	io.WriteString(w, `</table></div></td></tr></table>`)
}

func (r *Report) emailStatistics(w *quotedprintable.Writer) {
	if len(r.Sources) == 0 && len(r.Jurisdictions) == 0 && len(r.MTypeCounts) == 0 {
		return
	}
	io.WriteString(w, `<table cellspacing="0" cellpadding="0"><tr>`)
	if len(r.Sources) != 0 {
		var hasDown bool
		io.WriteString(w, `<td style="vertical-align:top"><div style="max-width:640px;margin-bottom:24px"><div style="font-size:20px;font-weight:bold;color:#444">Sources</div><table cellspacing="0" cellpadding="0">`)
		for _, source := range r.Sources {
			var down string
			if source.SimulatedDown {
				down, hasDown = `*`, true
			}
			fmt.Fprintf(w, `<tr><td style="padding-top:2px;color:#666">%s%s</td><td style="padding:2px 0 0 16px;text-align:right">%d</td></tr>`, html.EscapeString(source.Name), down, source.Count)
		}
		io.WriteString(w, `</table>`)
		if hasDown {
			io.WriteString(w, `<div>*Simulated outage</div>`)
		}
		io.WriteString(w, `</div></td>`)
	}
	if len(r.Jurisdictions) != 0 {
		var cols = (len(r.Jurisdictions) + 5) / 6
		var rows = (len(r.Jurisdictions) + cols - 1) / cols
		io.WriteString(w, `<td style="padding-left:32px;vertical-align:top"><div style="max-width:640px;margin-bottom:24px"><div style="font-size:20px;font-weight:bold;color:#444">Jurisdictions</div><table cellspacing="0" cellpadding="0"><tr>`)
		for col := 0; col < len(r.Jurisdictions); col += rows {
			io.WriteString(w, `<td style="vertical-align:top"><table cellspacing="0" cellpadding="0">`)
			for i := col; i < len(r.Jurisdictions) && i < col+rows; i++ {
				fmt.Fprintf(w, `<tr><td style="padding-top:2px;color:#666">%s</td><td style="padding:2px 0 0 16px;text-align:right">%d</td></tr>`, html.EscapeString(r.Jurisdictions[i].Name), r.Jurisdictions[i].Count)
			}
			io.WriteString(w, `</table></td>`)
		}
		io.WriteString(w, `</tr></table></div></td>`)
	}
	if len(r.MTypeCounts) != 0 {
		io.WriteString(w, `<td style="padding-left:32px;vertical-align:top"><div style="max-width:640px;margin-bottom:24px"><div style="font-size:20px;font-weight:bold;color:#444">Types</div><table cellspacing="0" cellpadding="0">`)
		for _, mtype := range r.MTypeCounts {
			fmt.Fprintf(w, `<tr><td style="padding-top:2px;color:#666">%s</td><td style="padding:2px 0 0 16px;text-align:right">%d</td></tr>`, html.EscapeString(mtype.Name), mtype.Count)
		}
		io.WriteString(w, `</table></div></td>`)
	}
	io.WriteString(w, `</tr></table>`)
}

func (r *Report) emailMessages(w *quotedprintable.Writer) {
	var hasMultiple bool
	io.WriteString(w, `<div style="max-width:640px;margin-bottom:24px"><div style="font-size:20px;font-weight:bold;color:#444">Messages</div><table cellspacing="0" cellpadding="0">`)
	for _, m := range r.Messages {
		var multiple string
		fmt.Fprintf(w, `<tr><td style="padding-top:4px;text-align:right">%s</td><td style="padding-top:4px;font-weight:bold">%s</td>`, html.EscapeString(m.Prefix), html.EscapeString(m.Suffix))
		if m.Multiple {
			multiple, hasMultiple = `*`, true
		}
		fmt.Fprintf(w, `<td style="padding:4px 0 0 16px">%s%s</td><td style="padding:4px 0 0 16px">%s</td><td style="padding:4px 0 0 16px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;color:%s">%s%s</td></tr>`,
			m.Source, multiple, m.Jurisdiction, classColor[m.Class], classLabel[m.Class], m.Problem)
	}
	io.WriteString(w, `</table>`)
	if hasMultiple {
		io.WriteString(w, `<div>*multiple messages from this address; only the last one counts</div>`)
	}
	io.WriteString(w, `</div>`)
}

func (r *Report) emailGenInfo(w *quotedprintable.Writer) {
	fmt.Fprintf(w, `<div>%s</div>`, html.EscapeString(r.GenerationInfo))
}

var classColor = map[string]string{
	"ok":      "green",
	"warning": "#ed7d31",
	"error":   "red",
	"invalid": "#888",
}

var classLabel = map[string]string{
	"ok":      "OK",
	"warning": "WARNING: ",
	"error":   "ERROR: ",
	"invalid": "NOT COUNTED: ",
}
