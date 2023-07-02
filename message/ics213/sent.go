package ics213

import "time"

// UpdateSent updates the message contents to reflect the fact that it is about
// to be sent.
func (f *ICS213) UpdateSent(opcall, opname string) {
	f.OpCall, f.OpName = opcall, opname
	f.OpDate, f.OpTime = time.Now().Format("01/02/2006"), time.Now().Format("15:04")
	f.ReceivedSent = "sender"
}
