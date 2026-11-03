package telnet

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rothskeller/packet/v4/jnos"
)

// prompt is how long a call that should return promptly may take.
const prompt = 2 * time.Second

// newPair starts a loopback listener, connects a Transport to it (without the
// login sequence), and returns the Transport and the server side of the
// connection.  Both are cleaned up at the end of the test.
func newPair(t *testing.T, log io.Writer) (tr *Transport, server net.Conn) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	accepted := make(chan net.Conn, 1)
	go func() {
		c, err := ln.Accept()
		if err != nil {
			close(accepted)
			return
		}
		accepted <- c
	}()
	tr = &Transport{log: log}
	if err = tr.connectTelnet(context.Background(), ln.Addr().String()); err != nil {
		t.Fatal(err)
	}
	select {
	case server = <-accepted:
		if server == nil {
			t.Fatal("accept failed")
		}
	case <-time.After(prompt):
		t.Fatal("accept timed out")
	}
	t.Cleanup(func() {
		server.Close()
		within(t, "cleanup Close", func() { tr.Close() })
	})
	return tr, server
}

// within runs f and fails the test if it does not return within prompt.
func within(t *testing.T, what string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	select {
	case <-done:
	case <-time.After(prompt):
		t.Fatalf("%s did not return within %s", what, prompt)
	}
}

func TestReadUntilAfterPeerClose(t *testing.T) {
	tr, server := newPair(t, nil)
	server.Write([]byte("hello\r\n"))
	server.Close()

	var data string
	var err error
	within(t, "first ReadUntil", func() { data, err = tr.ReadUntil("never") })
	if err != jnos.ErrDisconnected {
		t.Errorf("first ReadUntil err = %v, want ErrDisconnected", err)
	}
	if data != "hello\n" {
		t.Errorf("first ReadUntil data = %q, want %q", data, "hello\n")
	}
	within(t, "second ReadUntil", func() { data, err = tr.ReadUntil("never") })
	if err != jnos.ErrDisconnected {
		t.Errorf("second ReadUntil err = %v, want ErrDisconnected", err)
	}
	if data != "" {
		t.Errorf("second ReadUntil data = %q, want empty", data)
	}
	within(t, "third ReadUntil", func() { _, err = tr.ReadUntil("never") })
	if err != jnos.ErrDisconnected {
		t.Errorf("third ReadUntil err = %v, want ErrDisconnected", err)
	}
	within(t, "Send", func() { err = tr.Send("hi") })
	if err != jnos.ErrDisconnected {
		t.Errorf("Send err = %v, want ErrDisconnected", err)
	}
	within(t, "Close", func() { tr.Close() })
}

// checkWriter fails the test if Write is entered concurrently, or after
// closed is set.
type checkWriter struct {
	inWrite atomic.Bool
	closed  atomic.Bool
	bad     atomic.Int32
	count   atomic.Int64
}

func (w *checkWriter) Write(b []byte) (int, error) {
	if w.closed.Load() {
		w.bad.Add(1)
	}
	if !w.inWrite.CompareAndSwap(false, true) {
		w.bad.Add(1)
		return len(b), nil
	}
	w.count.Add(1)
	time.Sleep(20 * time.Microsecond) // widen the window for overlap
	w.inWrite.Store(false)
	return len(b), nil
}

func TestNoConcurrentLogWrites(t *testing.T) {
	w := &checkWriter{}
	tr, server := newPair(t, w)
	stop := make(chan struct{})
	var wg sync.WaitGroup

	// Server: drain whatever the client sends, and flood it with lines.
	wg.Add(2)
	go func() { defer wg.Done(); io.Copy(io.Discard, server) }()
	go func() {
		defer wg.Done()
		bw := bufio.NewWriter(server)
		for i := 0; ; i++ {
			select {
			case <-stop:
				return
			default:
			}
			server.SetWriteDeadline(time.Now().Add(prompt))
			if _, err := fmt.Fprintf(bw, "line %d\r\n", i); err != nil {
				return
			}
			if i%8 == 0 {
				if bw.Flush() != nil {
					return
				}
			}
		}
	}()
	// Client: read lines on one goroutine...
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if _, err := tr.ReadUntil("\n"); err != nil {
				return
			}
		}
	}()
	// ... and send on another.
	deadline := time.Now().Add(300 * time.Millisecond)
	for time.Now().Before(deadline) {
		if err := tr.Send("some outgoing text"); err != nil {
			t.Fatalf("Send: %v", err)
		}
	}
	within(t, "Close", func() { tr.Close() })
	w.closed.Store(true)
	close(stop)
	server.Close()
	within(t, "helpers exit", wg.Wait)
	time.Sleep(20 * time.Millisecond)
	if n := w.bad.Load(); n != 0 {
		t.Errorf("log writer called concurrently or after Close %d times", n)
	}
	if w.count.Load() < 10 {
		t.Errorf("log writer called only %d times; test did not exercise it", w.count.Load())
	}
}

func TestCloseWithUnreadData(t *testing.T) {
	tr, server := newPair(t, nil)
	// Flood the client, which never reads, so that the reader goroutine
	// fills readch and blocks.
	go func() {
		chunk := []byte(strings.Repeat("x", 1024))
		for {
			server.SetWriteDeadline(time.Now().Add(prompt))
			if _, err := server.Write(chunk); err != nil {
				return
			}
		}
	}()
	time.Sleep(200 * time.Millisecond)
	if len(tr.readch) != cap(tr.readch) {
		t.Errorf("readch has %d of %d entries; reader not blocked as intended", len(tr.readch), cap(tr.readch))
	}
	var err error
	within(t, "Close", func() { err = tr.Close() })
	if err != nil {
		t.Errorf("Close: %v", err)
	}
	select {
	case <-tr.readerDone:
	case <-time.After(prompt):
		t.Fatal("reader goroutine did not exit after Close")
	}
	within(t, "second Close", func() { err = tr.Close() })
	if err != nil {
		t.Errorf("second Close: %v", err)
	}
	within(t, "ReadUntil after Close", func() { _, err = tr.ReadUntil("x") })
	if err != jnos.ErrDisconnected {
		t.Errorf("ReadUntil after Close err = %v, want ErrDisconnected", err)
	}
}

func TestAbortDuringReadUntil(t *testing.T) {
	tr, _ := newPair(t, nil)
	result := make(chan error, 1)
	go func() {
		_, err := tr.ReadUntil("never")
		result <- err
	}()
	time.Sleep(100 * time.Millisecond)
	start := time.Now()
	tr.Abort()
	select {
	case err := <-result:
		if err != jnos.ErrDisconnected {
			t.Errorf("ReadUntil err = %v, want ErrDisconnected", err)
		}
		if d := time.Since(start); d > time.Second {
			t.Errorf("ReadUntil took %s to return after Abort", d)
		}
	case <-time.After(prompt):
		t.Fatal("ReadUntil did not return after Abort")
	}
	var err error
	within(t, "ReadUntil after Abort", func() { _, err = tr.ReadUntil("never") })
	if err != jnos.ErrDisconnected {
		t.Errorf("ReadUntil after Abort err = %v, want ErrDisconnected", err)
	}
	within(t, "Send after Abort", func() { err = tr.Send("hi") })
	if err != jnos.ErrDisconnected {
		t.Errorf("Send after Abort err = %v, want ErrDisconnected", err)
	}
	within(t, "second Abort", tr.Abort)
	within(t, "Close after Abort", func() { err = tr.Close() })
	if err != nil {
		t.Errorf("Close after Abort: %v", err)
	}
	within(t, "Abort after Close", tr.Abort)
}

func TestAbortDuringSend(t *testing.T) {
	tr, _ := newPair(t, nil)
	// The server never reads, so a large enough Send blocks once the
	// socket buffers fill.
	big := strings.Repeat("y", 64<<20)
	result := make(chan error, 1)
	go func() { result <- tr.Send(big) }()
	time.Sleep(200 * time.Millisecond)
	select {
	case err := <-result:
		t.Fatalf("Send returned before Abort: %v", err)
	default:
	}
	tr.Abort()
	select {
	case err := <-result:
		if !errors.Is(err, jnos.ErrDisconnected) {
			t.Errorf("Send err = %v, want ErrDisconnected", err)
		}
	case <-time.After(prompt):
		t.Fatal("Send did not return after Abort")
	}
	within(t, "Close after Abort", func() { tr.Close() })
}

func TestConcurrentAbortAndClose(t *testing.T) {
	tr, _ := newPair(t, nil)
	var wg sync.WaitGroup
	for range 4 {
		wg.Add(2)
		go func() { defer wg.Done(); tr.Abort() }()
		go func() { defer wg.Done(); tr.Close() }()
	}
	within(t, "concurrent Abort and Close", wg.Wait)
	select {
	case <-tr.readerDone:
	default:
		t.Error("reader goroutine still running after Close")
	}
}

// TestAbortWhileConnecting checks that canceling the context while waiting for
// the BBS's login prompt ends the attempt promptly.
func TestAbortWhileConnecting(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go func() {
		if c, err := ln.Accept(); err == nil {
			defer c.Close()
			io.Copy(io.Discard, c) // never sends a login prompt
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(200*time.Millisecond, cancel)
	start := time.Now()
	if _, err = ConnectContext(ctx, ln.Addr().String(), "MBOX", "pw", nil); err == nil {
		t.Fatal("connect succeeded")
	}
	if elapsed := time.Since(start); elapsed > 3*time.Second {
		t.Errorf("connect took %v to notice the abort", elapsed)
	}
}
