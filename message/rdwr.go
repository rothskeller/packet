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
		by []byte
	)
	if by, err = os.ReadFile(filename); err != nil {
		slog.Error("os.ReadFile", "f", filename, "err", err)
		return nil, err
	}
	by = bytes.ReplaceAll(by, crlf, lf)
	if noHeaderOK && !hasHeadersRE.Match(by) {
		nby := make([]byte, len(by)+6)
		copy(nby, []byte("X: X\n\n"))
		copy(nby[6:], by)
		by = nby
	}
	return Parse(string(by), filename)
}

// Parse parses a message in RFC-5322 format and returns the corresponding
// Message and any associated non-fatal decoding issues.  The supplied filename
// is used only in error messages.  Parse returns nil and an error if the
// message cannot be decoded.
func Parse(str, filename string) (m Message, err error) {
	var (
		by     []byte
		msg    *mail.Message
		to     []string
		c      common
		issues error
		parser func(string, mail.Header, *common) (Message, error)
	)
	if msg, err = mail.ReadMessage(strings.NewReader(str)); err != nil {
		slog.Error("mail.ReadMessage", "msg", string(str), "err", err)
		return nil, fmt.Errorf("%s: %w", filename, err)
	}
	to = msg.Header["To"]
	to = append(to, msg.Header["Cc"]...)
	to = append(to, msg.Header["Bcc"]...)
	c.to = strings.Join(to, ", ")
	if ct := msg.Header.Get("Content-Type"); ct != "" {
		if mt, params, err := mime.ParseMediaType(ct); err != nil {
			slog.Error("mime.ParseMediaType", "f", filename, "ct", ct, "err", err)
			return nil, fmt.Errorf("%s: unsupported Content-Type %q", filename, ct)
		} else if mt != "text/plain" || (params["charset"] != "" && strings.ToLower(params["charset"]) != "utf8") {
			slog.Error("unsupported Content-Type", "f", filename, "ct", ct)
			return nil, fmt.Errorf("%s: unsupported Content-Type %q", filename, ct)
		}
	}
	if by, err = io.ReadAll(msg.Body); err != nil {
		slog.Error("io.ReadAll msg.Body", "f", filename, "err", err)
		return nil, fmt.Errorf("%s: body: %w", filename, err)
	}
	if c.payload, issues = payload.Decode(msg.Header, string(by)); c.payload == nil {
		slog.Error("payload.Decode", "f", filename, "err", issues)
		return nil, fmt.Errorf("%s: %w", filename, issues)
	}
	c.subject = subject.DecodePlainSubject(msg.Header.Get("Subject"))
	c.init()
	if msg.Header.Get("Received") != "" {
		parser = readReceivedMessage
	} else if msg.Header.Get("Date") != "" {
		parser = readSentMessage
	} else {
		parser = readDraftMessage
	}
	if m, err = parser(filename, msg.Header, &c); m == nil {
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
