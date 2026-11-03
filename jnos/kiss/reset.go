package kiss

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// A resetPort is the connection to the radio used for resetting its TNC.
// (go.bug.st/serial's Port satisfies it.)
type resetPort interface {
	io.ReadWriteCloser
	SetReadTimeout(time.Duration) error
}

// Timing for resetting a Kenwood TNC.  They are variables so that tests can
// shorten them.
var (
	// resetSettle is how long after our last transmission to wait before
	// resetting the TNC, so that a final frame (such as our UA to the
	// BBS's DISC) has time to get on the air rather than being discarded.
	resetSettle = 3 * time.Second
	// resetExitWait is how long to wait after the KISS exit command for
	// the radio to leave KISS mode.
	resetExitWait = 300 * time.Millisecond
	// resetReplyWait is how long to wait for the radio to answer a command.
	resetReplyWait = 1500 * time.Millisecond
)

// tnRE matches the radio's answer to the TN command: the TNC mode (0 off, 1
// APRS, 2 KISS) and a second setting, which we preserve.
var tnRE = regexp.MustCompile(`^TN (\d),(\d)$`)

// kenwoodReset resets a Kenwood TNC (such as the TH-D75's) to a clean state:
// it leaves KISS mode, which discards any frames the TNC still holds, and then
// re-enters KISS mode with Kenwood's TN command, keeping the second TN setting
// that the radio reports.  It returns the command used to re-enter KISS mode.
// It sends nothing that transmits.
func kenwoodReset(port resetPort) (reentered string, err error) {
	// The KISS "return" command: leave KISS mode.
	if _, err = port.Write([]byte{fend, 0xFF, fend}); err != nil {
		return "", fmt.Errorf("send KISS exit: %w", err)
	}
	time.Sleep(resetExitWait)
	// If the radio wasn't in KISS mode, it took the KISS exit bytes as the
	// start of a command, and would answer our next one with "?".  A bare
	// carriage return ends that garbage command; its answer is drained.
	if _, err = port.Write([]byte{'\r'}); err != nil {
		return "", fmt.Errorf("send CR: %w", err)
	}
	time.Sleep(resetExitWait)
	drain(port)
	// Ask for the TNC setting.  Out of KISS mode, the radio answers "TN
	// 0,n"; if it doesn't answer, it isn't in the state we expect, so stop
	// rather than guess.  (A "?" gets one more try.)
	reply, err := catCommand(port, "TN")
	if err == nil && reply == "?" {
		drain(port)
		reply, err = catCommand(port, "TN")
	}
	if err != nil {
		return "", err
	}
	match := tnRE.FindStringSubmatch(reply)
	if match == nil {
		return "", fmt.Errorf("unexpected answer to TN: %q", reply)
	}
	reentered = "TN 2," + match[2]
	if reply, err = catCommand(port, reentered); err != nil {
		return "", err
	}
	if reply != reentered {
		return "", fmt.Errorf("unexpected answer to %s: %q", reentered, reply)
	}
	return reentered, nil
}

var errNoAnswer = errors.New("no answer")

// catCommand sends a command to the radio and returns its one-line answer.
func catCommand(port resetPort, cmd string) (reply string, err error) {
	if _, err = port.Write([]byte(cmd + "\r")); err != nil {
		return "", fmt.Errorf("send %s: %w", cmd, err)
	}
	var (
		got      []byte
		buf      = make([]byte, 64)
		deadline = time.Now().Add(resetReplyWait)
	)
	for {
		if idx := bytes.IndexByte(got, '\r'); idx >= 0 {
			return string(bytes.TrimSpace(got[:idx])), nil
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return "", fmt.Errorf("%s: %w", cmd, errNoAnswer)
		}
		if err = port.SetReadTimeout(remaining); err != nil {
			return "", err
		}
		count, err := port.Read(buf)
		got = append(got, buf[:count]...)
		if err != nil {
			return "", fmt.Errorf("%s: %w", cmd, err)
		}
	}
}

// drain discards anything the radio has already sent.
func drain(port resetPort) {
	buf := make([]byte, 256)
	for range 16 {
		if port.SetReadTimeout(50*time.Millisecond) != nil {
			return
		}
		if count, err := port.Read(buf); count == 0 || err != nil {
			return
		}
	}
}

// sessionMarker returns the name of the file whose existence marks a session
// in progress on serialPort, or "" if there's nowhere to put it.  The file is
// created when a session starts, and removed when it ends with the TNC in good
// order.  If it's there when a session starts, the last one never finished
// (the program was killed, or the computer crashed or slept, in the middle of
// it), and the TNC may still hold its frames.
func sessionMarker(serialPort string) string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	name := []byte(filepath.Base(serialPort))
	for i, b := range name {
		if !('a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || '0' <= b && b <= '9' || b == '.' || b == '-') {
			name[i] = '_'
		}
	}
	return filepath.Join(dir, "packet", "kiss-session-"+string(name))
}

// markSession creates the session marker file (see sessionMarker).  It's best
// effort: without it, we just can't tell later that the session never ended.
func markSession(marker string) {
	if err := os.MkdirAll(filepath.Dir(marker), 0o755); err == nil {
		err = os.WriteFile(marker, []byte(time.Now().Format(time.RFC3339)+"\n"), 0o644)
		if err != nil {
			slog.Warn("KISS session marker", "err", err)
		}
	}
}

// clearSession removes the session marker file (see sessionMarker).
func clearSession(marker string) {
	if marker != "" {
		os.Remove(marker)
	}
}

// resetBeforeStart resets a Kenwood TNC on the (just opened) port, because the
// last session on it never finished.  It reports to the monitor.  If the reset
// fails, the TNC may be out of KISS mode, so it returns an error saying how to
// fix that rather than letting the connection fail obscurely.
func resetBeforeStart(port resetPort, monitor func(string)) error {
	monitorNote(monitor, "The last session on this port didn't finish; resetting the TNC first.")
	cmd, err := kenwoodReset(port)
	if err != nil {
		slog.Warn("TNC reset before start", "err", err)
		monitorNote(monitor, fmt.Sprintf("TNC reset failed: %s.", err))
		return fmt.Errorf("the radio's TNC needed a reset, which failed (%s); turn the TNC off and back to KISS mode on the radio, then try again", err)
	}
	slog.Info("TNC reset before start", "reentered", cmd)
	monitorNote(monitor, "TNC reset: left KISS mode, discarding anything it still held, and re-entered it with "+cmd)
	return nil
}
