package message

import (
	"iter"
	"net/mail"
	"net/textproto"
	"strings"
	"time"

	"github.com/rothskeller/packet/v4/errors"
	"github.com/rothskeller/packet/v4/message/address"
	"github.com/rothskeller/packet/v4/message/field"
	"github.com/rothskeller/packet/v4/message/msgifc"
	"github.com/rothskeller/packet/v4/message/payload"
	"github.com/rothskeller/packet/v4/message/subject"
)

var ErrInvalidTo = errors.New("The \"To\" address list is not valid.")

// A DraftMessage is an envelope for a message that the local system is
// preparing to send but has not yet sent.
type DraftMessage struct {
	*common
	bulletin    bool
	readyToSend bool
}

var _ Message = (*DraftMessage)(nil)

// NewDraftMessage creates a new DraftMessage with the specified MType, Subject,
// and Payload.  It is marked as a bulletin if the bulletin flag is set.  It is
// not initially marked ReadyToSend, and has no initial To address.
func NewDraftMessage(mtype MType, subject subject.Subject, payload payload.Payload, bulletin bool) (m *DraftMessage) {
	m = &DraftMessage{
		common: &common{
			MType:   mtype,
			subject: subject,
			payload: payload,
		},
		bulletin: bulletin,
	}
	m.init()
	return m
}

func (m *DraftMessage) Bulletin() bool    { return m.bulletin }
func (m *DraftMessage) ReadyToSend() bool { return m.readyToSend }
func (m *DraftMessage) SetBulletin(bull bool) {
	if m.bulletin != bull {
		m.bulletin = bull
		m.MarkDirty("envelope.DraftMessage.Bulletin")
	}
}
func (m *DraftMessage) SetReadyToSend(ready bool) {
	if m.readyToSend != ready {
		m.readyToSend = ready
		m.MarkDirty("envelope.DraftMessage.ReadyToSend")
	}
}
func (m *common) SetTo(to string) (err error) {
	to = strings.ReplaceAll(to, "\n", " ")
	if _, err = address.ParseList(to); err != nil {
		err = ErrInvalidTo
	}
	m.to = to
	m.MarkDirty("envelope.DraftMessage.To")
	return err
}

func (m *DraftMessage) RFC5322() string {
	hdr := make(textproto.MIMEHeader)

	if m.bulletin {
		hdr.Set("X-Packet-Bulletin", "true")
	}
	if m.readyToSend {
		hdr.Set("X-Packet-Queued", "true")
	}
	return m.rfc5322(hdr)
}

// readDraftMessage creates a DraftMessage with the supplied common message
// fields and based on the supplied headers.
func readDraftMessage(filename string, hdr mail.Header, c *common) (_ Message, err error) {
	m := DraftMessage{common: c}

	m.bulletin = hdr.Get("X-Packet-Bulletin") != ""
	m.readyToSend = hdr.Get("X-Packet-Queued") != ""
	return &m, nil
}

// ToSentMessage converts a DraftMessage to a SentMessage.  It would be
// called when the message has been sent.
func (m *DraftMessage) ToSentMessage(from string, date time.Time) (s *SentMessage) {
	s = &SentMessage{
		common:   m.common,
		from:     from,
		date:     date,
		bulletin: m.bulletin,
	}
	s.init()
	return s
}

func (m *DraftMessage) Fields() iter.Seq[msgifc.Field] {
	return func(yield func(msgifc.Field) bool) {
		if !yield(draftToField) {
			return
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

var draftToField = field.NewField("", "To Address").
	Common(field.CHeaderTo).
	ValueFunc(func(m Message) string { return m.To() }).
	FromHumanFunc(func(_ msgifc.Message, s string) string {
		if addrs, err := address.ParseList(s); err == nil {
			trim := make([]string, len(addrs))
			for i := range addrs {
				trim[i] = addrs[i].Address
			}
			return strings.Join(trim, ", ")
		}
		return s
	}).
	SetValueFunc(func(m msgifc.Message, s string) { m.(*DraftMessage).SetTo(s) }).
	Required().
	ValidateFunc(func(m msgifc.Message, f msgifc.Field, vf msgifc.ValidateFlags) error {
		if _, err := address.ParseList(m.To()); err != nil {
			return errors.New(`Field "To Address" contains an invalid packet address.`)
		}
		return nil
	}).
	EditHelp("This is the comma-separated list of packet addresses to which the message will be sent.  Packet addresses usually have the form callsign@bbsname, and there is usually only one of them.").
	MakeField()

// CopyFields copies the field values from one message to another.  Message
// numbers, operator fields, and header info are not copied.  Aside from that,
// common fields are copied to like common fields, and non-common fields are
// copied only if the messages are the same type.
func CopyFields(from Message, to *DraftMessage) {
	if from.Type() != to.Type() {
		var cfs = make(map[string]string)
		for f := range from.Fields() {
			switch f.Common() {
			case "", field.COriginMessageID, field.CDestinationMessageID, field.COperatorCall, field.COperatorDate,
				field.COperatorMethod, field.COperatorMethodOther, field.COperatorName, field.COperatorTime,
				field.CReceiverSender, field.CHeaderDate, field.CHeaderFrom, field.CHeaderReceived,
				field.CSubjectFormTag, field.CSubjectMessageID:
				// ignore
			default:
				cfs[f.Common()] = f.Value(from)
			}
		}
		for f := range to.Fields() {
			if value, ok := cfs[f.Common()]; ok {
				f.SetValue(to, value)
			}
		}
	} else {
		for f := range from.Fields() {
			switch f.Common() {
			case field.COriginMessageID, field.CDestinationMessageID, field.COperatorCall, field.COperatorDate,
				field.COperatorMethod, field.COperatorMethodOther, field.COperatorName, field.COperatorTime,
				field.CReceiverSender, field.CHeaderDate, field.CHeaderFrom, field.CHeaderReceived,
				field.CSubjectFormTag, field.CSubjectMessageID:
				// ignore
			default:
				if f.Settable() && len(f.Children()) == 0 {
					f.SetValue(to, f.Value(from))
				}
			}
		}
	}
}

// MakeReply copies some field values from the first message to the second,
// appropriate for making the second message a reply to the first.
func MakeReply(from *ReceivedMessage, to *DraftMessage) {
	var cfs = make(map[string]string)
	for f := range from.Fields() {
		if c := f.Common(); c != "" {
			cfs[c] = f.Value(from)
		}
	}
	for f := range to.Fields() {
		switch f.Common() {
		case field.CDefaultBody:
			maybeSetValue(to, f, cfs[field.CDefaultBody])
		case field.CFromContact:
			maybeSetValue(to, f, cfs[field.CToContact])
		case field.CFromICSPosition:
			maybeSetValue(to, f, cfs[field.CToICSPosition])
		case field.CFromLocation:
			maybeSetValue(to, f, cfs[field.CToLocation])
		case field.CFromName:
			maybeSetValue(to, f, cfs[field.CToName])
		case field.CHandling:
			maybeSetValue(to, f, cfs[field.CHandling])
		case field.CHeaderTo:
			maybeSetValue(to, f, cfs[field.CHeaderFrom])
		case field.CMessageSummary:
			maybeSetValue(to, f, cfs[field.CMessageSummary])
		case field.CReference:
			maybeSetValue(to, f, cfs[field.COriginMessageID])
		case field.CSubjectHandling:
			maybeSetValue(to, f, cfs[field.CSubjectHandling])
		case field.CSubjectSummary:
			maybeSetValue(to, f, cfs[field.CSubjectSummary])
		case field.CToContact:
			maybeSetValue(to, f, cfs[field.CFromContact])
		case field.CToICSPosition:
			maybeSetValue(to, f, cfs[field.CFromICSPosition])
		case field.CToLocation:
			maybeSetValue(to, f, cfs[field.CFromLocation])
		case field.CToName:
			maybeSetValue(to, f, cfs[field.CFromName])
		}
	}
}

func maybeSetValue(msg *DraftMessage, f field.Field, value string) {
	if value != "" {
		f.SetValue(msg, value)
	}
}
