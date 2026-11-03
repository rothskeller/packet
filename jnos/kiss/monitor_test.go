package kiss

import (
	"strings"
	"sync"
	"testing"
	"time"
)

// TestMonitorOtherTraffic checks that frames between other stations (and ones
// that can't be decoded) appear in the monitor, marked as other traffic, and
// don't disturb our link.
func TestMonitorOtherTraffic(t *testing.T) {
	var (
		mu    sync.Mutex
		lines []string
	)
	monitor := func(line string) { mu.Lock(); lines = append(lines, line); mu.Unlock() }
	tr, b := connectTest(t, testDevice, nil, monitor)
	other := address{call: "N6XYZ"}
	digi := address{call: "RELAY", ssid: 2}

	// A text I frame from another station to the BBS, by way of a
	// digipeater that has repeated it.
	raw := appendAddress(nil, testBBS, true, false)
	raw = appendAddress(raw, other, false, false)
	raw = appendAddress(raw, digi, true, true) // H bit set: repeated
	raw = append(raw, 2<<5|1<<1, pidNoL3)
	raw = append(raw, "hello\rbell\x07 here\r"...)
	b.conn.Write(kissEncode(raw))
	// A NET/ROM frame (binary data) from the BBS to another station.
	nr := appendAddress(nil, other, true, false)
	nr = appendAddress(nr, testBBS, false, true)
	nr = append(nr, ctlUI, 0xCF, 0xFF, 0x00, 0x01)
	b.conn.Write(kissEncode(nr))
	// Garbage.
	b.conn.Write(kissEncode([]byte{1, 2, 3}))
	// Our link still works.
	go b.write("prompt\n")
	if data, err := tr.ReadUntil("prompt\n"); err != nil || data != "prompt\n" {
		t.Fatalf("ReadUntil = %q, %v", data, err)
	}
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	got := strings.Join(lines, "\n")
	mu.Unlock()
	for _, want := range []string{
		" (other) N6XYZ>W6XSC-1,RELAY-2* <I C R2 S1>\nhello\nbell. here\n",
		" (other) W6XSC-1>N6XYZ <UI C>\n(3 bytes of protocol CF data)\n",
		" (other) undecodable frame, 3 bytes\n",
	} {
		if !strings.Contains(got+"\n", want) {
			t.Errorf("monitor lacks %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "R N6XYZ") || strings.Contains(got, "R W6XSC-1>N6XYZ") {
		t.Errorf("other traffic shown as ours:\n%s", got)
	}
}
