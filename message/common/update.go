package common

import "time"

// UpdateReceived updates the message contents to reflect the fact that it has
// just been received.
func (s *StdFields) UpdateReceived(dmi, opcall, opname string) {
	s.DestinationMsgID = dmi
	s.OpCall, s.OpName = opcall, opname
	s.OpDate, s.OpTime = time.Now().Format("01/02/2006"), time.Now().Format("15:04")
}

// UpdateSent updates the message contents to reflect the fact that it is about
// to be sent.
func (s *StdFields) UpdateSent(opcall, opname string) {
	s.OpCall, s.OpName = opcall, opname
	s.OpDate, s.OpTime = time.Now().Format("01/02/2006"), time.Now().Format("15:04")
}

// UpdateDelivered updates the message contents to reflect the fact that it has
// been delivered.
func (s *StdFields) UpdateDelivered(dmi string) {
	s.DestinationMsgID = dmi
}
