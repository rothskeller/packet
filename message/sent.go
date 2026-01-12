package message

import (
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"net/mail"
	"net/textproto"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/rothskeller/packet/message/field"
)

// A SentMessage is for a message that the local system sent to a BBS.
type SentMessage struct {
	*common
	from     string
	date     time.Time
	bulletin bool
	receipts []SentMessageReceipt
}

var _ Message = (*SentMessage)(nil)

func (m *SentMessage) init() {
	m.Tracker.OnDirty(func(reason string) {
		switch reason {
		case "message.SentMessage.Receipts", "body.FormBody.Field.DestMsgNo":
			// OK
		default:
			panic("SentMessage should not change: " + reason)
		}
	})
}

// From returns the sender of the message, i.m., the address we were using when
// we sent it.  (This includes which BBS we sent it through.)
func (m *SentMessage) From() string { return m.from }

// Date returns the timestamp when the message was sent (i.m., the Date:
// header).
func (m *SentMessage) Date() time.Time { return m.date }

// Bulletin returns whether the message was sent as a BBS bulletin (as opposed
// to a private message).
func (m *SentMessage) Bulletin() bool { return m.bulletin }

// Draft returns whether the message is a draft.
func (m *SentMessage) Draft() bool { return false }

// Received returns whether the message has been received by the local system.
// It returns false for a message that has been sent or is being prepared to be
// sent.
func (m *SentMessage) Received() bool { return false }

// Receipts returns an iterator on the receipts that we have received for the
// message.
func (m *SentMessage) Receipts() iter.Seq[SentMessageReceipt] {
	return slices.Values(m.receipts)
}

// AddReceipt adds a receipt to the list of received receipts for the sent
// message.
func (m *SentMessage) AddReceipt(r SentMessageReceipt) {
	m.receipts = append(m.receipts, r)
	m.MarkDirty("envelope.SentMessage.Receipts")
}

// SentMessageReceipt contains the details of a received receipt for a sent
// message.
type SentMessageReceipt struct {
	// ReceiverAddress is the address of the party that received the
	// message, taken from the From: line of the receipt.
	ReceiverAddress string
	// ReceiverMessageID is the message ID that the receiver assigned to
	// the message, if any.
	ReceiverMessageID string
	// ReceiptDate is the date and time when the receiver sent the receipt.
	// Note that there is no standard for the formatting of this field, so
	// it is unparsed.
	ReceiptDate string
	// HasBeenRead indicates whether the message has been read by a human
	// (i.m., we received a read receipt).  If false, this record is for a
	// delivery receipt.
	HasBeenRead bool
}

// RFC5322 returns the message encoded in RFC-5322 format for storage or email
// transmission.
func (m *SentMessage) RFC5322() string {
	hdr := make(textproto.MIMEHeader)

	if m.from != "" {
		hdr.Set("From", m.from)
	}
	hdr.Set("Date", m.date.Format(time.RFC1123Z))
	if m.bulletin {
		hdr.Set("X-Packet-Bulletin", "true")
	}
	for _, r := range m.receipts {
		if r.HasBeenRead {
			hdr.Add("X-Packet-Receipt", fmt.Sprintf("READ BY %s AT %s",
				r.ReceiverAddress, r.ReceiptDate))
		} else {
			hdr.Add("X-Packet-Receipt", fmt.Sprintf("DELIVERED TO %s AT %s RMI %s",
				r.ReceiverAddress, r.ReceiptDate, r.ReceiverMessageID))
		}
	}
	return m.rfc5322(hdr)
}

var (
	readReceiptHeaderRE     = regexp.MustCompile(`^READ BY (.*) AT (.*)$`)
	deliveryReceiptHeaderRE = regexp.MustCompile(`^DELIVERED TO (.*) AT (.*) RMI (.*)$`)
)

// readSentMessage creates a SentMessage with the supplied common message
// fields and based on the supplied headers.
func readSentMessage(filename string, hdr mail.Header, c *common) (_ Message, err error) {
	m := SentMessage{common: c}
	m.init()

	m.from = strings.Join(hdr["From"], ", ")
	if t, err := time.Parse(time.RFC1123Z, hdr.Get("Date")); err == nil {
		m.date = t
	}
	m.bulletin = hdr.Get("X-Packet-Bulletin") != ""
	for _, rlist := range hdr["X-Packet-Receipt"] {
		for rstr := range strings.SplitSeq(rlist, ",") {
			rstr = strings.TrimSpace(rstr)
			if match := readReceiptHeaderRE.FindStringSubmatch(rstr); match != nil {
				m.receipts = append(m.receipts, SentMessageReceipt{
					ReceiverAddress: match[1],
					ReceiptDate:     match[2],
					HasBeenRead:     true,
				})
			} else if match := deliveryReceiptHeaderRE.FindStringSubmatch(rstr); match != nil {
				m.receipts = append(m.receipts, SentMessageReceipt{
					ReceiverAddress:   match[1],
					ReceiptDate:       match[2],
					ReceiverMessageID: match[3],
				})
			} else {
				slog.Error("incorrect X-Packet-Receipt", "f", filename, "X-Packet-Receipt", rstr)
				return nil, errors.New("incorrect X-Packet-Receipt: header format for stored sent message")
			}
		}
	}
	return &m, nil
}

func (m *SentMessage) Fields() iter.Seq[field.Field] {
	return func(yield func(field.Field) bool) {
		for _, f := range sentFields {
			if !yield(f) {
				return
			}
		}
		for f := range m.Subject().Fields() {
			if !yield(f) {
				return
			}
		}
		for f := range m.Body().Fields() {
			if !yield(f) {
				return
			}
		}
	}
}

var sentFields = []field.Field{
	field.NewField("", "From").
		Common(field.CHeaderFrom).
		ValueFunc(func(m Message) string { return m.(*SentMessage).from }).
		MakeField(),
	field.NewField("", "To").
		Common(field.CHeaderTo).
		ValueFunc(func(m Message) string { return m.To() }).
		MakeField(),
	field.NewField("", "Sent").
		Common(field.CHeaderDate).
		ValueFunc(func(m Message) string { return m.(*SentMessage).date.Format("01/02/2006 15:04") }).
		MakeField(),
	field.NewField("", "Received").
		Common(field.CHeaderReceived).
		ValueFunc(func(m Message) string {
			var rstr []string
			for _, r := range m.(*SentMessage).receipts {
				if !r.HasBeenRead {
					rstr = append(rstr, fmt.Sprintf("by %s at %s as %s", r.ReceiverAddress, r.ReceiptDate, r.ReceiverMessageID))
				}
			}
			return strings.Join(rstr, ", ")
		}).MakeField(),
}
