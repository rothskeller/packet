package envelope

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// RenderSaved renders the supplied envelope and body, in a form suitable for
// later reading by ParseSaved.  Note that a few envelope fields are not
// preserved, as noted in their documentation.
func (env *Envelope) RenderSaved(body string) string {
	var sb strings.Builder
	if env.ReceivedBBS != "" {
		if env.ReceivedArea != "" {
			fmt.Fprintf(&sb, "Received: FROM %s.ampr.org BY pktmsg.local FOR %s;\n\t%s\n",
				env.ReceivedBBS, env.ReceivedArea, env.ReceivedDate.Format(time.RFC1123Z))
		} else {
			fmt.Fprintf(&sb, "Received: FROM %s.ampr.org BY pktmsg.local; %s\n",
				env.ReceivedBBS, env.ReceivedDate.Format(time.RFC1123Z))
		}
	}
	if !env.IsFinal() && env.ReadyToSend {
		sb.WriteString("X-Packet-Queued: true\n")
	}
	if env.From != "" {
		fmt.Fprintf(&sb, "From: %s\n", env.From)
	}
	if len(env.To) != 0 {
		fmt.Fprintf(&sb, "To: %s\n", strings.Join(env.To, ",\n\t"))
	}
	if env.SubjectLine != "" {
		fmt.Fprintf(&sb, "Subject: %s\n", env.SubjectLine)
	}
	if !env.Date.IsZero() {
		fmt.Fprintf(&sb, "Date: %s\n", env.Date.Format(time.RFC1123Z))
	}
	sb.WriteByte('\n')
	sb.WriteString(env.RenderBody(body))
	return sb.String()
}

// RenderBody renders just the body part of the message according to the
// parameters in the envelope.
func (env *Envelope) RenderBody(body string) string {
	needB64 := strings.IndexFunc(body, nonASCII) >= 0
	if !needB64 && !env.OutpostUrgent && !env.RequestDeliveryReceipt && !env.RequestReadReceipt {
		if body != "" && !strings.HasSuffix(body, "\n") {
			body += "\n"
		}
		return body
	}
	var sb strings.Builder
	if env.OutpostUrgent {
		sb.WriteString("!URG!")
	}
	if env.RequestDeliveryReceipt {
		sb.WriteString("!RDR!")
	}
	if env.RequestReadReceipt {
		sb.WriteString("!RRR!")
	}
	sb.WriteString(body)
	body = sb.String()
	if needB64 {
		return "!B64!" + base64.StdEncoding.EncodeToString([]byte(body)) + "\n"
	}
	return body
}
func nonASCII(r rune) bool {
	return r > 126 || (r < 32 && r != '\t' && r != '\n')
}
