package kiss

import (
	"maps"
	"slices"
	"time"
)

// A Device describes a KISS TNC and the AX.25 link parameters to use with it.
type Device struct {
	// Name is the name of the device.
	Name string
	// PacLen is the maximum number of data bytes in each packet.
	PacLen int
	// MaxFrame is the maximum number of unacknowledged packets that can be
	// outstanding at once.
	MaxFrame int
	// Frack is how long to wait for an acknowledgment from the BBS before
	// polling it, not counting the time it takes to transmit the
	// unacknowledged packets.
	Frack time.Duration
	// InitialFrack is the Frack to use until the BBS's answer time has been
	// measured.  It is longer than Frack, so that a slow BBS doesn't get a
	// retransmission while its answer to the first copy is still on the
	// way; as the link measures the BBS's actual answer time, the wait
	// adapts to it (but never drops below Frack).  Zero means Frack.
	InitialFrack time.Duration
	// Retry is the number of times to poll the BBS without a response
	// before deciding the link has failed.
	Retry int
	// RespTime is how long the channel must be quiet after a packet from
	// the BBS before we send a plain acknowledgment, so that the BBS can
	// finish sending a burst of packets first.  Acknowledgments go sooner
	// when the BBS polls for one or when we have data to send, which
	// carries it.
	RespTime time.Duration
	// KenwoodReset says that the TNC is a Kenwood one, whose KISS mode can
	// be left and re-entered from the computer (with the KISS return
	// command and Kenwood's TN command).  After a session that didn't end
	// cleanly, we do that, which discards any frames the TNC still holds,
	// so that it doesn't transmit them after we've given up and starts the
	// next session clean.
	KenwoodReset bool
	// Turnaround is the minimum time after a packet from the BBS before we
	// transmit anything, giving the BBS's radio time to switch from
	// transmitting back to receiving.  (The TNC's carrier detect can't tell
	// when the BBS is ready to hear us: it sees only that the channel is
	// clear.)  Packets sent sooner than this are lost: on W3XSC, packets
	// sent within a second of its own were never heard.  It also keeps us
	// from transmitting into a burst: W3XSC sometimes sends two or three
	// frames up to a second apart (an RR, then the data it carries).  With
	// FullDuplex, it is all that keeps us off the BBS's transmissions.
	// This is the equivalent of a TNC's DWAIT or RESPTIME setting.
	Turnaround time.Duration
	// TXDelay, Persist and SlotTime are sent to the TNC at connect time
	// (as KISS commands), to control how it decides when to transmit.  The
	// TNC keys up only when its carrier detect shows the channel is clear,
	// then (each slot time) with probability (Persist+1)/256, and then
	// sends TXDelay worth of preamble before the data, which also gives a
	// BBS that just stopped transmitting time to start receiving.  Zero
	// values are not sent, leaving the TNC's own settings in effect.
	TXDelay  time.Duration
	Persist  int
	SlotTime time.Duration
	// FullDuplex makes the TNC transmit as soon as it has a frame, without
	// waiting for its carrier detect to show the channel clear (KISS
	// command 05).  We still wait out Turnaround after the BBS's frames.
	// It is restored to half duplex when the session ends.
	FullDuplex bool
}

var devices = map[string]*Device{
	// The TH-D75's built-in TNC supports only KISS mode for packet
	// connections.  It is reached over USB or Bluetooth, whichever is
	// selected in the radio's KISS interface menu setting.
	//
	// Its TNC must be in full-duplex mode (ignoring its carrier detect).
	// In half duplex, it often held our frames for seconds, or until our
	// next write, though the channel was quiet (its BUSY light off): the
	// BBS got our first copy together with our retransmission, and
	// rejected the duplicate, or nothing at all (the radio stuck showing
	// STA).  Idle, it transmitted within a second either way.  In full
	// duplex, a whole session with W3XSC went without a single retry.  We
	// still wait out Turnaround after the BBS's frames, so we don't
	// transmit over it.
	//
	// The other settings suit a slow, lossy path to W3XSC: one frame at a
	// time (so a lost frame doesn't take the next one with it), and
	// patient: T1 starts at 15 seconds, then adapts to the measured answer
	// time (but no lower than 8 seconds), backing off for each retry, up to
	// 30 seconds; 6 retries cover about 3 minutes of silence from the BBS.
	"Kenwood TH-D75": {
		Name:         "Kenwood TH-D75",
		KenwoodReset: true,
		PacLen:       128,
		MaxFrame:     1,
		Frack:        8 * time.Second,
		InitialFrack: 15 * time.Second,
		Retry:        6,
		RespTime:     500 * time.Millisecond,
		Turnaround:   2 * time.Second,
		TXDelay:      400 * time.Millisecond,
		Persist:      63,
		SlotTime:     100 * time.Millisecond,
		FullDuplex:   true,
	},
	"Generic KISS TNC": {
		Name:     "Generic KISS TNC",
		PacLen:   128,
		MaxFrame: 2,
		Frack:    6 * time.Second,
		Retry:    8,
		// Without measurements, wait twice Frack before resending.
		InitialFrack: 12 * time.Second,
		RespTime:     500 * time.Millisecond,
		Turnaround:   2 * time.Second,
		TXDelay:      400 * time.Millisecond,
		Persist:      63,
		SlotTime:     100 * time.Millisecond,
	},
}

// AllDevices returns an alphabetical list of the names of all supported KISS
// devices.
func AllDevices() []string {
	return slices.Sorted(maps.Keys(devices))
}

// GetDevice returns the definition of the KISS device with the specified name,
// or nil if there is none.
func GetDevice(name string) *Device {
	return devices[name]
}

// Limits on device parameters, enforced by validated.
const (
	// maxMaxFrame is the largest usable MaxFrame: with modulo 8 sequence
	// numbers, a window of 8 or more would make N(R) ambiguous.
	maxMaxFrame = 7
	// maxPacLen is the largest information field AX.25 allows.
	maxPacLen = 256
	// minFrack is the shortest Frack we allow.  It is far too short for
	// real use at 1200 baud; it just keeps a missing or nonsensical value
	// from making the link spin.
	minFrack = 100 * time.Millisecond
)

// validated returns a copy of d with its link parameters forced into their
// valid ranges.
func (d *Device) validated() *Device {
	v := *d
	v.MaxFrame = min(max(v.MaxFrame, 1), maxMaxFrame)
	v.PacLen = min(max(v.PacLen, 1), maxPacLen)
	v.Frack = max(v.Frack, minFrack)
	v.InitialFrack = min(max(v.InitialFrack, v.Frack), max(maxFrack, v.Frack))
	v.Retry = max(v.Retry, 1)
	v.RespTime = max(v.RespTime, 0)
	v.Turnaround = max(v.Turnaround, 0)
	return &v
}
