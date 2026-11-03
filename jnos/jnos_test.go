package jnos

import (
	"errors"
	"strings"
	"sync"
	"testing"
	"time"
)

// fakeBlock is how long the fake transport's ReadUntil waits, when it has no
// data to return, before giving up with ErrTimeout.  It bounds every test even
// if the code under test is broken.
const fakeBlock = 300 * time.Millisecond

// fakeTransport is a scripted, in-memory Transport.  Each Send appends the
// scripted response for that exact string (if any) to the pending input.
// ReadUntil returns pending input through the delimiter.  When the pending
// input has no delimiter, ReadUntil returns endErr immediately if it is set,
// otherwise it blocks (until Abort or fakeBlock) and returns ErrTimeout.  Once
// the fake has returned ErrDisconnected, every later call is counted in
// afterGone and ReadUntil blocks, mimicking a telnet read after a disconnect.
type fakeTransport struct {
	mu        sync.Mutex
	input     string
	responses map[string]string
	endErr    error
	sendErr   error
	sent      []string
	reads     int
	gone      bool
	afterGone int
	closed    int
	blocked   chan struct{} // receives when ReadUntil starts blocking
	abortCh   chan struct{}
	abortOnce sync.Once
}

func newFake(initial string, responses map[string]string) *fakeTransport {
	return &fakeTransport{
		input:     initial,
		responses: responses,
		blocked:   make(chan struct{}, 10),
		abortCh:   make(chan struct{}),
	}
}

func (f *fakeTransport) ReadUntil(delim string) (string, error) {
	f.mu.Lock()
	f.reads++
	if f.gone {
		f.afterGone++
		f.mu.Unlock()
		return f.block()
	}
	select {
	case <-f.abortCh:
		f.mu.Unlock()
		return "", errors.New("fake: aborted")
	default:
	}
	if idx := strings.Index(f.input, delim); idx >= 0 {
		data := f.input[:idx+len(delim)]
		f.input = f.input[idx+len(delim):]
		f.mu.Unlock()
		return data, nil
	}
	if f.endErr != nil {
		data, err := f.input, f.endErr
		f.input = ""
		if errors.Is(err, ErrDisconnected) {
			f.gone = true
		}
		f.mu.Unlock()
		return data, err
	}
	f.mu.Unlock()
	return f.block()
}

func (f *fakeTransport) block() (string, error) {
	select {
	case f.blocked <- struct{}{}:
	default:
	}
	select {
	case <-f.abortCh:
		return "", errors.New("fake: aborted")
	case <-time.After(fakeBlock):
		return "", ErrTimeout
	}
}

func (f *fakeTransport) Send(s string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.gone {
		f.afterGone++
		return ErrDisconnected
	}
	f.sent = append(f.sent, s)
	if f.sendErr != nil {
		if errors.Is(f.sendErr, ErrDisconnected) {
			f.gone = true
		}
		return f.sendErr
	}
	f.input += f.responses[s]
	return nil
}

func (f *fakeTransport) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.closed++
	return nil
}

func (f *fakeTransport) UseVerboseReads() bool { return false }

func (f *fakeTransport) Abort() { f.abortOnce.Do(func() { close(f.abortCh) }) }

func (f *fakeTransport) set(fn func(f *fakeTransport)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	fn(f)
}

func (f *fakeTransport) sentList() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.sent...)
}

func (f *fakeTransport) counts() (reads, sends, afterGone, closed int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.reads, len(f.sent), f.afterGone, f.closed
}

// connect returns a Conn over a fake transport that has completed the
// connection handshake.  The responses map is used for later commands.
func connect(t *testing.T, ident string, responses map[string]string) (*Conn, *fakeTransport) {
	t.Helper()
	if responses == nil {
		responses = map[string]string{}
	}
	responses["XM 0\n"] = "(#1) >\n"
	if ident != "" {
		responses["# DE "+ident+"\n"] = "(#1) >\n"
	}
	f := newFake("[JNOS-2.0-B1FHIM$]\nWelcome\n(#1) >\n", responses)
	c, err := Connect(f, ident)
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}
	return c, f
}

const msg5 = "Date: Fri, 17 Oct 2026 10:00:00 PDT\nFrom: a@b\nTo: c@d\nSubject: hi\n\nbody line\n"

// readMsg5 checks that the Conn is in sync by reading message 5.
func readMsg5(t *testing.T, c *Conn, f *fakeTransport, first string) {
	t.Helper()
	f.set(func(f *fakeTransport) { f.responses["R 5\n"] = first + "\n" + msg5 + "(#1) >\n" })
	got, err := c.Read(5)
	if err != nil {
		t.Fatalf("Read(5): %v", err)
	}
	if got != msg5 {
		t.Fatalf("Read(5) = %q, want %q", got, msg5)
	}
}

func TestStickyDisconnect(t *testing.T) {
	ops := map[string]func(c *Conn) error{
		"List":         func(c *Conn) error { _, err := c.List(""); return err },
		"Read":         func(c *Conn) error { _, err := c.Read(1); return err },
		"Kill":         func(c *Conn) error { return c.Kill(1) },
		"SetArea":      func(c *Conn) error { return c.SetArea("xsc") },
		"Send":         func(c *Conn) error { return c.Send("s", "b", "x@y") },
		"SendBulletin": func(c *Conn) error { return c.SendBulletin("s", "b", "x@y") },
	}
	for _, viaSend := range []bool{false, true} {
		for name, first := range ops {
			t.Run(name, func(t *testing.T) {
				c, f := connect(t, "KC6RSC", nil)
				if viaSend {
					f.set(func(f *fakeTransport) { f.sendErr = ErrDisconnected })
				} else {
					f.set(func(f *fakeTransport) { f.endErr = ErrDisconnected })
				}
				c.nextIdent = time.Time{} // ident would be due
				if err := first(c); !errors.Is(err, ErrDisconnected) {
					t.Fatalf("first call: got %v, want ErrDisconnected", err)
				}
				for oname, op := range ops {
					if err := op(c); !errors.Is(err, ErrDisconnected) {
						t.Errorf("%s after disconnect: got %v, want ErrDisconnected", oname, err)
					}
				}
				start := time.Now()
				if err := c.Close(); err != nil {
					t.Errorf("Close: %v", err)
				}
				if d := time.Since(start); d > fakeBlock/2 {
					t.Errorf("Close took %v", d)
				}
				if _, _, after, closed := f.counts(); after != 0 || closed != 1 {
					t.Errorf("transport calls after disconnect = %d, closes = %d", after, closed)
				}
				for _, s := range f.sentList() {
					if s == "B\r" {
						t.Errorf("B sent after disconnect")
					}
				}
			})
		}
	}
}

func TestAbortDuringRead(t *testing.T) {
	c, f := connect(t, "", nil)
	done := make(chan error, 1)
	go func() { _, err := c.List(""); done <- err }()
	select {
	case <-f.blocked:
	case <-time.After(5 * time.Second):
		t.Fatal("List never blocked in ReadUntil")
	}
	start := time.Now()
	c.Abort()
	var err error
	select {
	case err = <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("List did not return after Abort")
	}
	if d := time.Since(start); d > fakeBlock/2 {
		t.Errorf("List returned %v after Abort; want prompt return", d)
	}
	if err == nil || !errors.Is(err, ErrDisconnected) {
		t.Errorf("List after Abort: got %v, want aborted error", err)
	}
	r0, s0, _, _ := f.counts()
	if err = c.Kill(1); err == nil {
		t.Error("Kill after Abort succeeded")
	}
	if err = c.Close(); err != nil {
		t.Errorf("Close after Abort: %v", err)
	}
	if r, s, _, closed := f.counts(); r != r0 || s != s0 || closed != 1 {
		t.Errorf("transport used after Abort: reads %d->%d sends %d->%d closes %d", r0, r, s0, s, closed)
	}
}

func TestListNoneDisconnect(t *testing.T) {
	for name, resp := range map[string]string{
		"line1": "None to list.\n",
		"line4": "Mail area: xsc\n0 messages  -  0 new\n\nNo messages\n",
	} {
		t.Run(name, func(t *testing.T) {
			c, f := connect(t, "", map[string]string{"LA\n": resp})
			f.set(func(f *fakeTransport) { f.endErr = ErrDisconnected })
			ml, err := c.List("")
			if ml != nil || !errors.Is(err, ErrDisconnected) {
				t.Fatalf("List = %v, %v; want nil, ErrDisconnected", ml, err)
			}
		})
	}
}

func TestListNone(t *testing.T) {
	c, f := connect(t, "", map[string]string{"LA\n": "Mail area: xsc\n0 messages  -  0 new\n\nNone to list.\n(#1) >\n"})
	ml, err := c.List("")
	if ml != nil || err != nil {
		t.Fatalf("List = %v, %v; want nil, nil", ml, err)
	}
	readMsg5(t, c, f, "Message #5 ")
}

func TestListFormatErrorDrains(t *testing.T) {
	const hdr = "St.  #  TO            FROM     DATE   SIZE SUBJECT\n"
	const msgline = "> N   5 xsc           a        Oct 17  123 hello\n"
	for name, resp := range map[string]string{
		"header":  "Bogus\nextra line\n(#1) >\n",
		"line2":   "Mail area: xsc\nbogus\n\n" + hdr + msgline + "(#1) >\n",
		"line3":   "Mail area: xsc\n1 message  -  1 new\nnot blank\n" + hdr + msgline + "(#1) >\n",
		"listhdr": "Mail area: xsc\n1 message  -  1 new\n\nbogus header\n" + msgline + "(#1) >\n",
		"msgline": "Mail area: xsc\n1 message  -  1 new\n\n" + hdr + msgline + "garbage line\nmore\n(#1) >\n",
	} {
		t.Run(name, func(t *testing.T) {
			c, f := connect(t, "", map[string]string{"LA\n": resp})
			if ml, err := c.List(""); err == nil {
				t.Fatalf("List = %v, nil; want format error", ml)
			}
			readMsg5(t, c, f, "Message #5 ")
		})
	}
}

func TestListOK(t *testing.T) {
	c, _ := connect(t, "", map[string]string{"LA\n": "Mail area: xsc_news\n1 message  -  1 new\n\n" +
		"St.  #  TO            FROM     DATE   SIZE SUBJECT\n" +
		"> N   5 xsc           a        Oct 17  123 hello\n" +
		"Area: xsc_news (#1) >\n"})
	ml, err := c.List("")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if ml.Area != "xsc_news" || ml.Count != 1 || len(ml.Messages) != 1 || ml.Messages[0].Number != 5 || ml.Messages[0].SubjectPrefix != "hello" {
		t.Fatalf("List = %+v", ml)
	}
}

func TestReadFlags(t *testing.T) {
	for _, first := range []string{"Message #5 ", "Message #5 [Held]", "Message #5 [Deleted]"} {
		t.Run(first, func(t *testing.T) {
			c, f := connect(t, "", nil)
			readMsg5(t, c, f, first)
		})
	}
}

func TestAreaPrompt(t *testing.T) {
	for _, prompt := range []string{"Area: xsc_news (#1) >", "Area: xsc-news.2 (#12) >", "(#1) >"} {
		if !jnosPromptRE.MatchString(prompt) {
			t.Errorf("prompt %q not recognized", prompt)
		}
	}
	c, f := connect(t, "", map[string]string{
		"A xsc_news\n": "Mail area: xsc_news\n1 message  -  1 new\nArea: xsc_news (#1) >\n",
	})
	if err := c.SetArea("xsc_news"); err != nil {
		t.Fatalf("SetArea: %v", err)
	}
	readMsg5(t, c, f, "Message #5 [Held]")
}

func TestNoIdentAfterFailedSend(t *testing.T) {
	c, f := connect(t, "KC6RSC", map[string]string{"SP x@y\n": "Subject:\n"})
	c.nextIdent = time.Time{} // ident is due
	if err := c.Send("subj", "body", "x@y"); !errors.Is(err, ErrTimeout) {
		t.Fatalf("Send: got %v, want ErrTimeout", err)
	}
	sent := f.sentList()
	if last := sent[len(sent)-1]; last != "subj\n" {
		t.Errorf("last sent %q; sent %q", last, sent)
	}
	for _, s := range sent[2:] { // skip connect-time ident and XM 0
		if strings.HasPrefix(s, "# DE") {
			t.Errorf("ident sent after failed Send: %q", sent)
		}
	}
}

func TestIdentAfterSuccessfulSend(t *testing.T) {
	c, f := connect(t, "KC6RSC", map[string]string{
		"SP x@y\n":       "Subject:\n",
		"subj\n":         "Enter message.  End with /EX or ^Z in first column (^A aborts):\n",
		"body\n/EX\n":    "Msg queued\n(#1) >\n",
		"# DE KC6RSC\n":  "(#1) >\n",
		"XM 0\n":         "(#1) >\n",
		"unused-command": "",
	})
	c.nextIdent = time.Time{}
	if err := c.Send("subj", "body", "x@y"); err != nil {
		t.Fatalf("Send: %v", err)
	}
	sent := f.sentList()
	if last := sent[len(sent)-1]; last != "# DE KC6RSC\n" {
		t.Errorf("no ident after successful Send; sent %q", sent)
	}
}

func TestSetAreaNoSuchDrains(t *testing.T) {
	c, f := connect(t, "", map[string]string{"A bogus\n": "No such message area: bogus\n(#1) >\n"})
	if err := c.SetArea("bogus"); err == nil || !strings.Contains(err.Error(), "no such message area") {
		t.Fatalf("SetArea: got %v", err)
	}
	readMsg5(t, c, f, "Message #5 ")
}

func TestCloseReadErrorIsDisconnect(t *testing.T) {
	for name, endErr := range map[string]error{
		"disconnected": ErrDisconnected,
		"other":        errors.New("no response from BBS"),
		"wrapped":      errors.Join(errors.New("x"), ErrDisconnected),
	} {
		t.Run(name, func(t *testing.T) {
			c, f := connect(t, "", nil)
			f.set(func(f *fakeTransport) { f.endErr = endErr })
			if err := c.Close(); err != nil {
				t.Errorf("Close: %v", err)
			}
			sent := f.sentList()
			if sent[len(sent)-1] != "B\r" {
				t.Errorf("B not sent; sent %q", sent)
			}
			if _, _, _, closed := f.counts(); closed != 1 {
				t.Errorf("closes = %d", closed)
			}
			if err := c.Close(); err != nil {
				t.Errorf("second Close: %v", err)
			}
			if _, _, _, closed := f.counts(); closed != 1 {
				t.Errorf("second Close closed transport again")
			}
		})
	}
}
