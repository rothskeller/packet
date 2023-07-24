package checkout

// New creates a new check-out message with default values.
func New() *CheckOut {
	return &CheckOut{
		Handling: "ROUTINE",
	}
}
