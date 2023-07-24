package checkin

// New creates a new check-in message with default values.
func New() *CheckIn {
	return &CheckIn{
		Handling: "ROUTINE",
	}
}
