package message

import (
	"bytes"
	_ "embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/form/htmlop"
	"github.com/rothskeller/packet/message/body"
	"github.com/rothskeller/packet/message/payload"
	"github.com/rothskeller/packet/message/subject"
	"golang.org/x/net/html"
)

var (
	ErrHasFormTag         = errors.New("The subject line of a plain text message must not contain a form name.")
	ErrNonStandardSubject = errors.New("The subject line is not in SCCo-standard format.")
	ErrNoBody             = errors.New("The message body is empty.")
)

type ErrNotPlainTextBody string

func (e ErrNotPlainTextBody) Error() string {
	return fmt.Sprintf("A plain text message must have a plain text body, not %s.", string(e))
}

// A PlainMessage is a plain text message with no special structure or
// interpretation.
type plainMessage struct{ *BaseEditableMType }

var _ EditableMType = plainMessage{}

var PlainMessage = plainMessage{NewBaseEditableMType("a plain text message", "plain", "p")}

// Recognize does nothing.  Messages are assigned the PlainMessage type by
// SetType when nothing else matches.
func (mt plainMessage) Recognize(m Message) {}

/*
// Validate validates a plain message.
func (m *PlainMessage) Validate(pifo bool) (err error) {
	// Unless the message is a bulletin, it should have an SCCo-standard
	// subject line.
	if m.envelope.Bulletin() {
		if m.envelope.Subject().EncodedSubject() == "" && !pifo {
			err = errors.Join(err, subject.ErrNoSubject)
		}
	} else if !pifo {
		switch s := m.envelope.Subject().(type) {
		case *subject.SCCoSubject:
			err = errors.Join(err, s.Validate(!m.envelope.Received()))
		case *subject.SCCoFormSubject:
			err = errors.Join(err,
				ErrHasFormTag,
				s.Validate(!m.envelope.Received()),
			)
		default:
			err = errors.Join(err, ErrNonStandardSubject)
		}
	}
	switch b := m.envelope.Body().(type) {
	case *body.PlainBody:
		if b.EncodedBody() == "" && !pifo {
			err = errors.Join(err, ErrNoBody)
		}
	default:
		err = errors.Join(err, ErrNotPlainTextBody(fmt.Sprintf("%T", b)))
	}
	return err
}
*/

func (mt plainMessage) NewDraft() (msg Message) {
	s, _ := subject.NewPlainSubject("", "", "")
	return NewDraftMessage(PlainMessage, s, payload.NewOutpostPayload(body.NewPlainBody("")), false)
}

//go:embed plain.html
var plainHTML []byte

// EditHTML returns the HTML form for editing the message.
func (mt plainMessage) EditHTML(msg Message, vars EditHTMLVars) (out []byte, err error) {
	var (
		formHTML *html.Node
		formBuf  bytes.Buffer
		fields   = make(map[string]string)
		values   = make(url.Values)
	)
	// Read and parse the HTML for the form.
	if formHTML, err = html.Parse(bytes.NewReader(plainHTML)); err != nil {
		slog.Error("html.Parse", "err", err)
		return nil, err
	}
	fields["submit-url"] = vars.SubmitURL
	fields["submit-label"] = vars.SubmitLabel
	fields["save-label"] = vars.SaveLabel
	if vars.ShowAddressFields {
		fields["show-addrs"] = "true"
	}
	// Expand the templates in the form HTML, using the supplied fields.
	htmlop.Expand(formHTML, fields)
	// Fill in the form using the fields from the message.
	values.Set("ToAddr", msg.To())
	values.Set("FromAddr", vars.FromAddress)
	values.Set("MsgNo", msg.Subject().SubjectMessageID())
	values.Set("handling", msg.Subject().SubjectHandling())
	values.Set("subject", msg.Subject().SubjectSummary())
	values.Set("body", msg.Body().EncodedBody())
	htmlop.FillForm(formHTML, values)
	// Render and minimize the result.
	htmlop.Minify(&formBuf, formHTML)
	return formBuf.Bytes(), nil
}

// EditAssets returns the file system containing the form assets.
func (mt plainMessage) EditAssets() (assets fs.FS) { return nil }

// FromPOST translates the HTML response back into a DraftMessage.
func (mt plainMessage) FromPOST(r *http.Request) (msg Message, err error) {
	var (
		bdy  *body.PlainBody
		subj *subject.PlainSubject
		payl *payload.OutpostPayload
		dm   *DraftMessage
	)
	bdy = body.NewPlainBody(r.FormValue("body"))
	payl = payload.NewOutpostPayload(bdy)
	subj, _ = subject.NewPlainSubject(r.FormValue("MsgNo"), r.FormValue("handling"), r.FormValue("subject"))
	dm = NewDraftMessage(PlainMessage, subj, payl, false)
	dm.SetTo(r.FormValue("ToAddr"))
	return dm, nil
}
