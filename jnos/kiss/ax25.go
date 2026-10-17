package kiss

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// An address is an AX.25 station address: a call sign and an SSID.
type address struct {
	call string
	ssid byte
}

// parseAddress parses an address string like "W6XSC-1" or "KC6RSC".
func parseAddress(s string) (a address, err error) {
	call, ssid, hasSSID := strings.Cut(strings.ToUpper(s), "-")
	if call == "" || len(call) > 6 || strings.TrimFunc(call, isAlnum) != "" {
		return a, fmt.Errorf("%q is not a valid AX.25 address", s)
	}
	a.call = call
	if hasSSID {
		if n, err := strconv.Atoi(ssid); err != nil || n < 0 || n > 15 {
			return a, fmt.Errorf("%q is not a valid AX.25 address", s)
		} else {
			a.ssid = byte(n)
		}
	}
	return a, nil
}

func isAlnum(r rune) bool { return (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') }

func (a address) String() string {
	if a.ssid == 0 {
		return a.call
	}
	return fmt.Sprintf("%s-%d", a.call, a.ssid)
}

// appendAddress appends the 7-byte encoding of a to b.  cbit is the
// command/response bit for the address; last indicates that this is the last
// address in the address field.
func appendAddress(b []byte, a address, cbit, last bool) []byte {
	for i := range 6 {
		c := byte(' ')
		if i < len(a.call) {
			c = a.call[i]
		}
		b = append(b, c<<1)
	}
	ssid := 0x60 | a.ssid<<1
	if cbit {
		ssid |= 0x80
	}
	if last {
		ssid |= 0x01
	}
	return append(b, ssid)
}

// decodeAddress decodes a 7-byte encoded address.
func decodeAddress(b []byte) (a address, cbit, last bool) {
	var call [6]byte

	for i := range 6 {
		call[i] = b[i] >> 1
	}
	a.call = strings.TrimRight(string(call[:]), " ")
	a.ssid = (b[6] >> 1) & 0x0F
	return a, b[6]&0x80 != 0, b[6]&0x01 != 0
}

// Control field values (modulo 8 operation).  For I and S frames, the N(S)
// and N(R) sequence numbers are added to these.
const (
	ctlRR   = 0x01
	ctlRNR  = 0x05
	ctlREJ  = 0x09
	ctlSABM = 0x2F
	ctlDISC = 0x43
	ctlDM   = 0x0F
	ctlUA   = 0x63
	ctlFRMR = 0x87
	ctlUI   = 0x03
	pfBit   = 0x10
	// pidNoL3 is the protocol ID for I frames carrying plain text.
	pidNoL3 = 0xF0
)

// A frame is a decoded AX.25 frame.
type frame struct {
	dest, src address
	// command is true for command frames and false for response frames.
	command bool
	// digis is true if the frame has a digipeater path.
	digis   bool
	control byte
	info    []byte
}

func (f *frame) isI() bool     { return f.control&0x01 == 0 }
func (f *frame) isS() bool     { return f.control&0x03 == 0x01 }
func (f *frame) isUI() bool    { return f.control&^pfBit == ctlUI }
func (f *frame) ns() byte      { return f.control >> 1 & 7 }
func (f *frame) nr() byte      { return f.control >> 5 & 7 }
func (f *frame) pf() bool      { return f.control&pfBit != 0 }
func (f *frame) sType() byte   { return f.control & 0x0F }
func (f *frame) uType() byte   { return f.control &^ pfBit }
func (f *frame) hasInfo() bool { return f.isI() || f.isUI() }

// String describes the frame in the style of a packet monitor, e.g.
// "KO6NVY>W3XSC-1 <I C P R1 S0>".
func (f *frame) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s>%s <", f.src, f.dest)
	switch {
	case f.isI():
		sb.WriteString("I")
	case f.isS():
		sb.WriteString(map[byte]string{ctlRR: "RR", ctlRNR: "RNR", ctlREJ: "REJ"}[f.sType()])
	default:
		name, ok := map[byte]string{ctlSABM: "SABM", ctlDISC: "DISC", ctlDM: "DM", ctlUA: "UA", ctlFRMR: "FRMR", ctlUI: "UI"}[f.uType()]
		if !ok {
			name = fmt.Sprintf("U%02X", f.uType())
		}
		sb.WriteString(name)
	}
	if f.command {
		sb.WriteString(" C")
		if f.pf() {
			sb.WriteString(" P")
		}
	} else {
		sb.WriteString(" R")
		if f.pf() {
			sb.WriteString(" F")
		}
	}
	if f.isI() || f.isS() {
		fmt.Fprintf(&sb, " R%d", f.nr())
	}
	if f.isI() {
		fmt.Fprintf(&sb, " S%d", f.ns())
	}
	sb.WriteString(">")
	return sb.String()
}

// encode returns the encoded frame (without any FCS, which the TNC adds).
func (f *frame) encode() []byte {
	b := make([]byte, 0, 16+len(f.info))
	b = appendAddress(b, f.dest, f.command, false)
	b = appendAddress(b, f.src, !f.command, true)
	b = append(b, f.control)
	if f.hasInfo() {
		b = append(b, pidNoL3)
		b = append(b, f.info...)
	}
	return b
}

var errBadFrame = errors.New("malformed AX.25 frame")

// decodeFrame decodes an AX.25 frame.
func decodeFrame(b []byte) (f *frame, err error) {
	var dcbit, scbit, last bool

	if len(b) < 15 {
		return nil, errBadFrame
	}
	f = new(frame)
	f.dest, dcbit, _ = decodeAddress(b[0:7])
	f.src, scbit, last = decodeAddress(b[7:14])
	b = b[14:]
	for !last {
		if len(b) < 8 {
			return nil, errBadFrame
		}
		f.digis = true
		_, _, last = decodeAddress(b[0:7])
		b = b[7:]
	}
	// AX.25 2.0 marks commands with the C bit set in the destination
	// address and clear in the source address.  Older stations set both
	// bits the same; we treat those frames as commands.
	f.command = dcbit || !scbit
	f.control = b[0]
	b = b[1:]
	if f.hasInfo() {
		if len(b) < 1 {
			return nil, errBadFrame
		}
		f.info = b[1:] // skip the PID
	}
	return f, nil
}

// KISS framing bytes.
const (
	fend  = 0xC0
	fesc  = 0xDB
	tfend = 0xDC
	tfesc = 0xDD
)

// KISS commands (for TNC port 0).
const (
	kissData     = 0x00
	kissTXDelay  = 0x01
	kissPersist  = 0x02
	kissSlotTime = 0x03
	// kissFullDuplex set to nonzero makes the TNC transmit without waiting
	// for the channel to be clear (ignoring carrier detect).
	kissFullDuplex = 0x05
)

// kissEncode wraps an AX.25 frame in a KISS data frame for TNC port 0.
func kissEncode(data []byte) []byte { return kissFrame(kissData, data) }

// kissFrame encodes a KISS frame with the specified command and contents.
func kissFrame(cmd byte, data []byte) []byte {
	b := make([]byte, 0, len(data)+8)
	b = append(b, fend, cmd)
	for _, c := range data {
		switch c {
		case fend:
			b = append(b, fesc, tfend)
		case fesc:
			b = append(b, fesc, tfesc)
		default:
			b = append(b, c)
		}
	}
	return append(b, fend)
}

// maxKISSFrame is the longest KISS frame kissDecoder accepts.  Real AX.25
// frames are far shorter (about 330 bytes at most); the limit just keeps
// garbage input without FENDs from growing the buffer without bound.
const maxKISSFrame = 4096

// A kissDecoder extracts AX.25 frames from a stream of KISS-encoded bytes.
type kissDecoder struct {
	buf []byte
	esc bool
	// overflow is set when the current frame has grown too long; the
	// rest of it, up to the next FEND, is discarded.
	overflow bool
	// command, if set, is called for each non-data KISS frame.  (TNCs
	// don't send these; it is used in testing.)
	command func(cmd byte, data []byte)
}

// feed processes received bytes, calling emit for each complete data frame.
func (d *kissDecoder) feed(data []byte, emit func([]byte)) {
	for _, c := range data {
		switch {
		case c == fend:
			// The first byte is the KISS command; the low nibble is 0
			// for data frames.  Anything else (including stray bytes
			// received before the first FEND, and frames that were too
			// long) is ignored.
			if d.overflow {
				// discard it
			} else if len(d.buf) > 1 && d.buf[0]&0x0F == 0 {
				emit(append([]byte(nil), d.buf[1:]...))
			} else if len(d.buf) != 0 && d.command != nil {
				d.command(d.buf[0], append([]byte(nil), d.buf[1:]...))
			}
			d.buf, d.esc, d.overflow = d.buf[:0], false, false
		case d.esc:
			switch c {
			case tfend:
				c = fend
			case tfesc:
				c = fesc
			}
			d.esc = false
			d.add(c)
		case c == fesc:
			d.esc = true
		default:
			d.add(c)
		}
	}
}

// add adds a byte to the frame being decoded, unless the frame is too long.
func (d *kissDecoder) add(c byte) {
	switch {
	case d.overflow:
	case len(d.buf) >= maxKISSFrame:
		d.buf, d.overflow = d.buf[:0], true
	default:
		d.buf = append(d.buf, c)
	}
}
