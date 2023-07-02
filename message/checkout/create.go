package checkout

// New creates a new check-out message with default values.
func New(opcall, opname string) *CheckOut {
	return &CheckOut{
		Handling:         "ROUTINE",
		OperatorCallSign: opcall,
		OperatorName:     opname,
	}
}
