package kiss

import (
	"time"
)

// A Device describes a KISS TNC and the AX.25 link parameters to use with it.
type Device struct {
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

// Validate returns a copy of d with its link parameters forced into their
// valid ranges.
func (d *Device) Validate() error {
	d.MaxFrame = min(max(d.MaxFrame, 1), maxMaxFrame)
	d.PacLen = min(max(d.PacLen, 1), maxPacLen)
	d.Frack = max(d.Frack, minFrack)
	d.InitialFrack = min(max(d.InitialFrack, d.Frack), max(maxFrack, d.Frack))
	d.Retry = max(d.Retry, 1)
	d.RespTime = max(d.RespTime, 0)
	d.Turnaround = max(d.Turnaround, 0)
	return nil
}
