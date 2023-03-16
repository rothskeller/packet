package pktmgr

import (
	"errors"
	"fmt"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/ics213"
)

// receivedHeaderRE is a regular expression that matches the contents of the
// Received: headers that we put in received messages.  It returns the BBS
// domain, the local ID if any, the bulletin area if any, and the received date.
var receivedHeaderRE = regexp.MustCompile(`^FROM (\S+) BY pktmgr.local(?: ID \S+)?(?: FOR (\S+))?; (.*)$`)

// Message is a single physical message.  It extends xscmsg.Message.
type Message struct {
	// Message is the parsed XSC message type.
	*xscmsg.Message
	// BBS is the domain name of the BBS from which the message was
	// received; it is empty for outgoing messages.
	BBS string
	// Bulletin is the category or area from which the received bulletin
	// message was retrieved; it is empty for other types of messages.
	Bulletin string
	// Received is the time at which the message was retrieved; it is zero
	// for outgoing messages.
	Received time.Time
	// From is the address of the sender of the message.  For outgoing
	// messages, it is set only when the message is ready to be sent; the
	// absence of a From address in an outgoing message indicates that it is
	// a draft message.
	From *mail.Address
	// To is the set of addresses to which the message is addressed.
	To []*mail.Address
	// SubjectLine is the raw subject line of the message.
	SubjectLine string
	// Sent is the time at which the message was sent; it is zero for
	// messages that have not yet been sent.
	Sent time.Time
	// mandr is a pointer to the MAndR that contains the message, so that
	// methods on the Message can reach MAndR and Incident details.
	mandr *MAndR
}

// readMessage reads a physical message from a file.  It returns an error if the
// message cannot be read or parsed.
func readMessage(filename string) (m *Message, err error) {
	var contents []byte

	if contents, err = os.ReadFile(filename); err != nil {
		return nil, err
	}
	return newMessage(string(contents))
}

// newMessage parses a message.  It returns an error if the message cannot be
// parsed.
func newMessage(raw string) (m *Message, err error) {
	var pkt *pktmsg.Message

	if pkt, err = pktmsg.ParseMessage(raw); err != nil {
		return nil, err
	}
	if pkt.Flags&pktmsg.AutoResponse != 0 {
		return nil, errors.New("message is an auto-response")
	}
	if pkt.Flags&pktmsg.NotPlainText != 0 && pkt.Body == "" {
		return nil, errors.New("message has no plain-text body")
	}
	m = new(Message)
	m.Message = xscmsg.Recognize(pkt)
	if match := receivedHeaderRE.FindStringSubmatch(pkt.Header.Get("Received")); match != nil {
		m.BBS = match[1]
		m.Bulletin = match[2]
		m.Received, _ = mail.ParseDate(match[3])
	}
	m.From, _ = mail.ParseAddress(pkt.Header.Get("From"))
	m.To, _ = mail.ParseAddressList(pkt.Header.Get("To"))
	m.Sent, _ = mail.ParseDate(pkt.Header.Get("Date"))
	m.SubjectLine = pkt.Header.Get("Subject")
	return m, nil
}

// IsReceived returns whether the message has been received.  Note this is true
// for both bulletins and regular messages.
func (m *Message) IsReceived() bool { return m.BBS != "" }

// IsBulletin returns whether the message is a received bulletin.
func (m *Message) IsBulletin() bool { return m.Bulletin != "" }

// IsSent returns whether the message has been sent.
func (m *Message) IsSent() bool { return m.BBS == "" && !m.Sent.IsZero() }

// IsReady returns whether the message is ready to be sent.
func (m *Message) IsReady() bool { return m.BBS == "" && m.Sent.IsZero() && m.From != nil }

// Save saves a physical message to the named file.  The local message ID is
// encoded in received messages if provided.
func (m *Message) save(filename, lmi string) (err error) {
	var fh *os.File

	if fh, err = os.Create(filename); err != nil {
		return err
	}
	defer fh.Close()
	// Write the Received: header if appropriate.
	if m.BBS != "" {
		fmt.Fprintf(fh, "Received: FROM %s BY pktmgr.local", m.BBS)
		linelen := len(m.BBS) + 64
		if lmi != "" {
			fmt.Fprintf(fh, " ID %s", lmi)
			linelen += len(lmi) + 4
		}
		if m.Bulletin != "" {
			fmt.Fprintf(fh, " FOR %s", m.Bulletin)
			linelen += len(m.Bulletin) + 5
		}
		if linelen > 78 {
			fmt.Fprint(fh, ";\n\t")
		} else {
			fmt.Fprint(fh, "; ")
		}
		fmt.Fprintln(fh, m.Received.Format(time.RFC1123Z))
	}
	// Write the From header if appropriate.
	if m.From != nil {
		fmt.Fprintf(fh, "From: %s\n", m.From)
	}
	// Write the To header if appropriate.
	if len(m.To) != 0 {
		tostr := make([]string, len(m.To))
		for i, to := range m.To {
			if to.Name != "" {
				tostr[i] = to.String()
			} else {
				tostr[i] = to.Address // without the <> that to.String produces
			}
		}
		fmt.Fprintf(fh, "To: %s\n", strings.Join(tostr, ", "))
	}
	// Write the Subject header if appropriate.
	if !m.IsReceived() && !m.IsSent() {
		m.SubjectLine = m.Subject()
	}
	if m.SubjectLine != "" {
		fmt.Fprintf(fh, "Subject: %s\n", m.SubjectLine)
	}
	// Write the Date header if appropriate.
	if !m.Sent.IsZero() {
		fmt.Fprintf(fh, "Date: %s\n", m.Sent.Format(time.RFC1123Z))
	}
	// Newline to end the headers.
	fmt.Fprintln(fh)
	// Write the message body.
	fmt.Fprint(fh, m.Body())
	return nil
}

// MarkSent marks a message as having been sent.
func (m *Message) MarkSent() {
	m.Sent = time.Now()
}

// MarkReady marks a message as queued to send.
func (m *Message) MarkReady() {
	addr := fmt.Sprintf("%s@%s.ampr.org",
		strings.ToLower(m.mandr.incident.config.callsign),
		strings.ToLower(m.mandr.incident.config.BBS))
	m.From = &mail.Address{Name: m.mandr.incident.config.name, Address: addr}
}

// MarkDraft marks a message as draft (not ready to send).
func (m *Message) MarkDraft() {
	m.From = nil
}

// SetOperatorFields sets the OpCall, OpName, OpDate, and OpTime fields of the
// message, if it has them.  For use with ICS-213 messages, it also sets the
// receiver/sender field.
func (m *Message) SetOperatorFields() {
	if f := m.KeyField(xscmsg.FOpCall); f != nil {
		f.Value = m.mandr.incident.config.OpCall
	}
	if f := m.KeyField(xscmsg.FOpName); f != nil {
		f.Value = m.mandr.incident.config.OpName
	}
	if f := m.KeyField(xscmsg.FOpDate); f != nil {
		f.Value = time.Now().Format("01/02/2006")
	}
	if f := m.KeyField(xscmsg.FOpTime); f != nil {
		f.Value = time.Now().Format("15:04")
	}
	if m.Type.Tag == ics213.Tag {
		if m.IsReceived() {
			m.Field("Rec-Sent").Value = "receiver"
		} else {
			m.Field("Rec-Sent").Value = "sender"
		}
		m.Field("Method").Value = "Other"
		m.Field("Other").Value = "Packet"
	}
}
