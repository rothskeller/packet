package receipt

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/body"
	"github.com/rothskeller/packet/message/cachetrack"
	"github.com/rothskeller/packet/message/payload"
	"github.com/rothskeller/packet/message/subject"
)

var (
	readReceiptMatchRE = regexp.MustCompile(`^\n*!RR!`)
	readReceiptRE      = regexp.MustCompile(`^\n*!RR!(.+)\nYour Message\n\nTo: (.+)\nSubject: (.+)\n\nwas read on .*\n`)
	ErrNoReadTime      = errors.New("The read time must be set.")
	ErrRRSyntax        = errors.New("The body cannot be decoded as a read receipt because it does not conform to read receipt syntax.")
)

// NewReadReceipt creates a new read receipt message (in draft state) with the
// specified parameters.  It returns an error if the parameters are invalid.
// rmTo is the destination address for the message being receipted.  rmSubject
// is the encoded subject line of the message being receipted.  readTime is the
// time that the message being receipted was displayed to our user.  extraText
// is free-form text to be added to the end of the receipt message body.  (This
// is rarely used since Outpost users will never see it.)
func NewReadReceipt(rmTo, rmSubject string, readTime time.Time, extraText string) (m message.Message, err error) {
	var (
		subject *ReadReceiptSubject
		body    *ReadReceiptBody
	)
	if subject, err = newReadReceiptSubject(rmSubject); err != nil {
		return nil, err
	}
	if body, err = newReadReceiptBody(rmTo, rmSubject, readTime, extraText); err != nil {
		return nil, err
	}
	return message.NewDraftMessage(ReadReceipt, subject, payload.NewOutpostPayload(body), false), nil
}

// --- MTYPE -------------------------------------------------------------------

type readReceipt struct{ *message.BaseMType }

var ReadReceipt readReceipt
var readReceiptOnce sync.Once

func init() {
	ReadReceipt.BaseMType = message.NewBaseMType("a read receipt", "")
	message.RegisterType(ReadReceipt)
}

func (mt readReceipt) Recognize(m message.Message) {
	if _, ok := m.Body().(*ReadReceiptBody); ok {
		m.SetType(ReadReceipt)
		readReceiptOnce.Do(readReceiptInit)
	}
}

func readReceiptInit() {
	ReadReceipt.AddField() // TODO
}

//--- BODY --------------------------------------------------------------------

// A ReadReceiptBody is a body of a read receipt message.
type ReadReceiptBody struct {
	cachetrack.NoTracker
	body           string
	messageTo      string
	messageSubject string
	readTime       string
	extraText      string
}

var _ body.Body = (*ReadReceiptBody)(nil)

// newReadReceiptBody creates a new body for a read receipt message with the
// supplied contents.
func newReadReceiptBody(rmTo, rmSubject string, readTime time.Time, extraText string) (b *ReadReceiptBody, err error) {
	if rmTo == "" {
		err = errors.Join(err, ErrNoRMTo)
	} else if strings.ContainsAny(rmTo, "\r\n") {
		err = errors.Join(err, ErrNewlineInRMTo)
	}
	if rmSubject == "" {
		err = errors.Join(err, ErrNoRMSubject)
	} else if strings.ContainsAny(rmSubject, "\r\n") {
		err = errors.Join(err, ErrNewlineInRMSubject)
	}
	if readTime.IsZero() {
		err = errors.Join(err, ErrNoReadTime)
	}
	if err != nil {
		return nil, err
	}
	ftime := readTime.Format("01/02/2006 15:04")
	return &ReadReceiptBody{
		body: fmt.Sprintf("!RR!%s\nYour Message\n\nTo: %s\nSubject: %s\n\nwas read on %[1]s\n%s",
			ftime, rmTo, rmSubject, extraText),
		messageTo:      rmTo,
		messageSubject: rmSubject,
		readTime:       ftime,
		extraText:      extraText,
	}, nil
}

func init() { body.RegisterDecoder(decodeReadReceiptBody) }
func decodeReadReceiptBody(body string) (body.Body, error) {
	if !readReceiptMatchRE.MatchString(body) {
		return nil, nil
	}
	if match := readReceiptRE.FindStringSubmatch(body); match == nil {
		return nil, ErrRRSyntax
	} else {
		return &ReadReceiptBody{
			body:           body,
			messageTo:      match[2],
			messageSubject: match[3],
			readTime:       match[1],
			extraText:      body[len(match[0]):],
		}, nil
	}
}

// EncodedBody returns the encoded read receipt message body.
func (b *ReadReceiptBody) EncodedBody() string { return b.body }

// MessageTo returns the address to which the message being receipted was sent.
func (b *ReadReceiptBody) MessageTo() string { return b.messageTo }

// MessageSubject returns the subject of the message being receipted.
func (b *ReadReceiptBody) MessageSubject() string { return b.messageSubject }

// ReadTime returns the date and time when the message being receipted was
// read by its receiver.  The date and time format are indeterminate.
func (b *ReadReceiptBody) ReadTime() string { return b.readTime }

// ExtraText returns any extra text in the receipt message following the
// receipt itself.
func (b *ReadReceiptBody) ExtraText() string { return b.extraText }

//--- SUBJECT -----------------------------------------------------------------

// A ReadReceiptSubject is a subject for a read receipt.
type ReadReceiptSubject struct {
	cachetrack.NoTracker
	rmSubject string
}

var _ subject.Subject = (*ReadReceiptSubject)(nil)

// newReadReceiptSubject creates a new subject line for a read receipt.
// rmSubject is the encoded subject line of the message being receipted.  The
// function returns an error if the rmSubject is invalid.
func newReadReceiptSubject(rmSubject string) (s *ReadReceiptSubject, err error) {
	if rmSubject == "" {
		return nil, ErrNoRMSubject
	} else if strings.ContainsAny(rmSubject, "\r\n") {
		return nil, ErrNewlineInRMSubject
	}
	return &ReadReceiptSubject{rmSubject: rmSubject}, nil
}

// EncodedSubject returns the encoded subject line.
func (s *ReadReceiptSubject) EncodedSubject() string { return "READ: " + s.rmSubject }

// ReceiptedMessageSubject returns the receipted message subject encoded in the
// receipt message subject line.
func (s *ReadReceiptSubject) ReceiptedMessageSubject() string { return s.rmSubject }

func init() { subject.RegisterDecoder(decodeReadReceiptSubject) }
func decodeReadReceiptSubject(subject string) (subject.Subject, error) {
	if strings.HasPrefix(subject, "READ: ") {
		return &ReadReceiptSubject{rmSubject: subject[6:]}, nil
	}
	return nil, nil
}
