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
	// Retry is the number of times to poll the BBS without a response
	// before deciding the link has failed.
	Retry int
	// RespTime is how long to wait after receiving a packet before
	// acknowledging it, so that the BBS can finish sending a burst of
	// packets without the acknowledgment colliding with them.
	RespTime time.Duration
}

var devices = map[string]*Device{
	// The TH-D75's built-in TNC supports only KISS mode for packet
	// connections.  It is reached over USB or Bluetooth, whichever is
	// selected in the radio's KISS interface menu setting.
	"Kenwood TH-D75": {
		Name:     "Kenwood TH-D75",
		PacLen:   128,
		MaxFrame: 2,
		Frack:    6 * time.Second,
		Retry:    8,
		RespTime: 1500 * time.Millisecond,
	},
	"Generic KISS TNC": {
		Name:     "Generic KISS TNC",
		PacLen:   128,
		MaxFrame: 2,
		Frack:    6 * time.Second,
		Retry:    8,
		RespTime: 1500 * time.Millisecond,
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
