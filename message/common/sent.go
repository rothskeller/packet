package common

import "time"

// UpdateSent updates the message contents to reflect the fact that it is about
// to be sent.
func (s *StdFields) UpdateSent(opcall, opname string) {
	s.OpCall, s.OpName = opcall, opname
	s.OpDate, s.OpTime = time.Now().Format("01/02/2006"), time.Now().Format("15:04")
}
