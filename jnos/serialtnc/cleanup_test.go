package serialtnc

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestCloseAfterAbortCleansUp checks that Close after an Abort reopens the port
// and gets the TNC back to command mode, disconnected, with its settings
// restored, rather than leaving it in converse mode with our settings.
func TestCloseAfterAbortCleansUp(t *testing.T) {
	setRFTimeout(t, 30*time.Second)
	fp := newFakePort(tncSim("\r\n*** CONNECTED to BBS-1\r\n"))
	fp2 := newFakePort(tncSim(""))
	var reopened int
	reopen := func() (port, error) { reopened++; return fp2, nil }
	var (
		tr  *Transport
		err error
	)
	within(t, 10*time.Second, func() { fp.Close() }, func() {
		tr, err = start(testTNC(), fp, "BBS-1", "MBOX", nil, withReopen(reopen))
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	tr.Abort()
	elapsed := within(t, 10*time.Second, func() { fp2.Close() }, func() { tr.Close() })
	if elapsed > 3*time.Second {
		t.Errorf("Close after Abort took %s", elapsed)
	}
	if reopened != 1 {
		t.Fatalf("reopened the port %d times, want 1", reopened)
	}
	lines := fp2.writtenLines()
	if !contains(lines, "^C") || !contains(lines, "D") || !contains(lines, "CR ON") {
		t.Errorf("cleanup after abort sent %q, want ^C, D, and the post-disconnect commands", lines)
	}
	if !fp2.isClosed() {
		t.Error("reopened port not closed")
	}
}

// TestCloseAfterAbortReopenFails checks that a failure to reopen the port
// doesn't hold up Close.
func TestCloseAfterAbortReopenFails(t *testing.T) {
	fp := newFakePort(nil)
	tr := newTransport(testTNC(), fp, nil)
	tr.connected = true
	tr.reopen = func() (port, error) { return nil, errPortClosed }
	tr.Abort()
	within(t, 2*time.Second, func() {}, func() { tr.Close() })
	if !fp.isClosed() {
		t.Error("port not closed")
	}
}

// TestCloseCarriesOnAfterFailedDisconnect checks that when the TNC never
// reports the disconnection, Close still restores the TNC settings.
func TestCloseCarriesOnAfterFailedDisconnect(t *testing.T) {
	setRFTimeout(t, 300*time.Millisecond)
	fp := newFakePort(func(line string) string {
		switch line {
		case "^C":
			return "\r\ncmd:"
		case "D":
			return "\r\ncmd:" // but never *** DISCONNECTED
		default:
			return "cmd:"
		}
	})
	tr := newTransport(testTNC(), fp, nil)
	tr.connected = true
	var err error
	within(t, 10*time.Second, func() { fp.Close() }, func() { err = tr.Close() })
	if err == nil {
		t.Error("Close didn't report the failed disconnection")
	}
	if lines := fp.writtenLines(); !contains(lines, "CR ON") {
		t.Errorf("post-disconnect commands not sent after failed disconnection: %q", lines)
	}
}

// TestCloseRetriesCommandMode checks that Close tries Ctrl-C again when the
// first gets no prompt, and restores the settings even if neither does.
func TestCloseRetriesCommandMode(t *testing.T) {
	setRFTimeout(t, 300*time.Millisecond)
	fp := newFakePort(func(line string) string { return "" }) // a TNC that never prompts
	tr := newTransport(testTNC(), fp, nil)
	tr.connected = true
	within(t, 10*time.Second, func() { fp.Close() }, func() { tr.Close() })
	lines := fp.writtenLines()
	var ctrlC int
	for _, l := range lines {
		if l == "^C" {
			ctrlC++
		}
	}
	if ctrlC != 2 {
		t.Errorf("sent Ctrl-C %d times, want 2: %q", ctrlC, lines)
	}
	if !contains(lines, "CR ON") {
		t.Errorf("post-disconnect commands not sent: %q", lines)
	}
}

// TestSetupRecoversFromConverseMode checks that a TNC left in converse mode
// (connected) by a session that didn't end cleanly is brought back to command
// mode and disconnected before we connect.
func TestSetupRecoversFromConverseMode(t *testing.T) {
	setRFTimeout(t, 5*time.Second)
	old := staleDisconnectWait
	staleDisconnectWait = 2 * time.Second
	t.Cleanup(func() { staleDisconnectWait = old })
	var (
		mu       sync.Mutex
		converse = true
	)
	fp := newFakePort(nil)
	fp.respond = func(line string) string {
		mu.Lock()
		defer mu.Unlock()
		switch {
		case line == "^C":
			converse = false
			return "\r\ncmd:"
		case converse:
			return "" // sent to the BBS
		case line == "D":
			// The prompt comes right away; the disconnection
			// report after the disconnection is done.
			go func() {
				time.Sleep(100 * time.Millisecond)
				fp.feed("\r\n*** DISCONNECTED\r\n")
			}()
			return "\r\ncmd:"
		case strings.HasPrefix(line, "C "):
			converse = true
			return "\r\n*** CONNECTED to BBS-1\r\n"
		default:
			return "cmd:"
		}
	}
	var (
		tr  *Transport
		err error
	)
	elapsed := within(t, 10*time.Second, func() { fp.Close() }, func() {
		tr, err = start(testTNC(), fp, "BBS-1", "MBOX", nil)
	})
	if err != nil {
		t.Fatalf("start: %v (sent %q)", err, fp.writtenLines())
	}
	if elapsed > 2*time.Second {
		t.Errorf("start took %s", elapsed)
	}
	lines := fp.writtenLines()
	idxD, idxC := -1, -1
	for i, l := range lines {
		if l == "D" && idxD < 0 {
			idxD = i
		}
		if strings.HasPrefix(l, "C ") {
			idxC = i
		}
	}
	if idxD < 0 || idxC < idxD {
		t.Errorf("didn't disconnect the stale connection before connecting: %q", lines)
	}
	tr.Abort()
	tr.Close()
}

// TestAbortWhileConnecting checks that canceling the context while waiting for
// the BBS to answer the connect command ends the attempt promptly, and cleans
// up the TNC (which may still be trying to connect).
func TestAbortWhileConnecting(t *testing.T) {
	setRFTimeout(t, 30*time.Second)
	fp := newFakePort(tncSim("")) // the connect is never answered
	fp2 := newFakePort(tncSim(""))
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(200*time.Millisecond, cancel)
	var err error
	elapsed := within(t, 10*time.Second, func() { fp.Close(); fp2.Close() }, func() {
		_, err = start(testTNC(), fp, "BBS-1", "MBOX", nil,
			withReopen(func() (port, error) { return fp2, nil }), withContext(ctx))
	})
	if err == nil {
		t.Fatal("connect succeeded")
	}
	if elapsed > 3*time.Second {
		t.Errorf("connect took %s to notice the abort", elapsed)
	}
	if lines := fp2.writtenLines(); !contains(lines, "^C") || !contains(lines, "D") || !contains(lines, "CR ON") {
		t.Errorf("cleanup after aborted connect sent %q", lines)
	}
}
