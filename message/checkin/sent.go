package checkin

// UpdateSent updates the message contents to reflect the fact that it is about
// to be sent.
func (f *CheckIn) UpdateSent(opcall, opname string) {
	f.OperatorCallSign, f.OperatorName = opcall, opname
}
