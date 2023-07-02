package checkin

// New creates a new check-in message with default values.
func New(opcall, opname string) *CheckIn {
	return &CheckIn{
		Handling:         "ROUTINE",
		OperatorCallSign: opcall,
		OperatorName:     opname,
	}
}
