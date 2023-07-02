package ics213

import "time"

// UpdateReceived updates the message contents to reflect the fact that it has
// just been received.
func (f *ICS213) UpdateReceived(opcall, opname string) {
	f.OpCall, f.OpName = opcall, opname
	f.OpDate, f.OpTime = time.Now().Format("01/02/2006"), time.Now().Format("15:04")
	f.ReceivedSent = "receiver"
}
