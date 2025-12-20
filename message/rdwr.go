package message

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/mail"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/rothskeller/packet/message/payload"
	"github.com/rothskeller/packet/message/subject"
)

var (
	lf           = []byte{'\n'}
	crlf         = []byte{'\r', '\n'}
	hasHeadersRE = regexp.MustCompile(`^[-A-Za-z0-9]+:`)
)

// Read reads a message in RFC-5322 format from the specified file and returns
// the corresponding Message and any associated non-fatal decoding issues.  It
// returns nil and an error if the file cannot be read or the message cannot be
// decoded.
func Read(filename string) (m Message, err error) {
	return read(filename, false)
}

// ReadNoHeader reads a message body, without RFC-5322 headers, from the
// specified file and returns the corresponding Message (as a DraftMessage) and
// any associated non-fatal decoding issues.  It returns nil and an error if
// the file cannot be read or the message cannot be decoded.  (This method is
// used when parsing message files provided by Outpost, which have no headers.)
func ReadNoHeader(filename string) (m Message, err error) {
	return read(filename, true)
}

func read(filename string, noHeaderOK bool) (m Message, err error) {
	var (
		by     []byte
		msg    *mail.Message
		to     []string
		c      common
		issues error
		parser func(mail.Header, *common) (Message, error)
	)
	if by, err = os.ReadFile(filename); err != nil {
		return nil, err
	}
	by = bytes.ReplaceAll(by, crlf, lf)
	if noHeaderOK && !hasHeadersRE.Match(by) {
		nby := make([]byte, len(by)+6)
		copy(nby, []byte("X: X\n\n"))
		copy(nby[6:], by)
		by = nby
	}
	if msg, err = mail.ReadMessage(bytes.NewReader(by)); err != nil {
		return nil, fmt.Errorf("%s: %w", filename, err)
	}
	to = msg.Header["To"]
	to = append(to, msg.Header["Cc"]...)
	to = append(to, msg.Header["Bcc"]...)
	c.to = strings.Join(to, ", ")
	c.subject, issues = subject.Decode(msg.Header.Get("Subject"))
	if ct := msg.Header.Get("Content-Type"); ct != "" {
		if mt, params, err := mime.ParseMediaType(ct); err != nil ||
			mt != "text/plain" || (params["charset"] != "" && strings.ToLower(params["charset"]) != "utf8") {
			return nil, fmt.Errorf("%s: unsupported Content-Type %q", filename, ct)
		}
	}
	if by, err = io.ReadAll(msg.Body); err != nil {
		return nil, fmt.Errorf("%s: body: %w", filename, err)
	}
	if c.payload, err = payload.Decode(msg.Header, string(by)); c.payload == nil {
		return nil, fmt.Errorf("%s: %w", filename, err)
	} else {
		issues = errors.Join(issues, err)
	}
	if b, ok := c.Body().(interface{ ShouldHaveFormSubject() bool }); ok && b.ShouldHaveFormSubject() {
		if s, ok := c.subject.(interface {
			ToFormSubject() (subject.Subject, error)
		}); ok {
			c.subject, err = s.ToFormSubject()
			issues = errors.Join(issues, err)
		}
	}
	c.init()
	if msg.Header.Get("Received") != "" {
		parser = readReceivedMessage
	} else if msg.Header.Get("Date") != "" {
		parser = readSentMessage
	} else {
		parser = readDraftMessage
	}
	if m, err = parser(msg.Header, &c); m == nil {
		return nil, err
	} else {
		SetType(m)
		return m, errors.Join(issues, err)
	}
}

// Write writes a message in RFC-5322 format to the specified file.
func Write(m Message, filename string) (err error) {
	encoded := []byte(m.RFC5322())
	if runtime.GOOS == "windows" {
		encoded = bytes.ReplaceAll(encoded, lf, crlf)
	}
	if err = os.WriteFile(filename, encoded, 0666); err != nil {
		slog.Error("message.Write", "fname", filename, "err", err)
		return err
	}
	return nil
}
