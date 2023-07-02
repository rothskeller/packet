package common

import "time"

// UpdateReceived updates the message contents to reflect the fact that it has
// just been received.
func (s *StdFields) UpdateReceived(opcall, opname string) {
	s.OpCall, s.OpName = opcall, opname
	s.OpDate, s.OpTime = time.Now().Format("01/02/2006"), time.Now().Format("15:04")
}
