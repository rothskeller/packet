package xscmsg

// HandlingOrder is the SCCo-standard message handling order code.
type HandlingOrder byte

// Values for HandlingOrder:
const (
	HandlingImmediate HandlingOrder = 'I'
	HandlingPriority  HandlingOrder = 'P'
	HandlingRoutine   HandlingOrder = 'R'
)

// ParseHandlingOrder parses an SCCo-standard message handling order name or
// code.  It returns false if the input string is not a valid handling order
// name or code.
func ParseHandlingOrder(s string) (HandlingOrder, bool) {
	switch s {
	case "I", "IMMEDIATE":
		return HandlingImmediate, true
	case "P", "PRIORITY":
		return HandlingPriority, true
	case "R", "ROUTINE":
		return HandlingRoutine, true
	}
	return 0, false
}

// String returns the full string form of the handling order.
func (ho HandlingOrder) String() string {
	switch ho {
	case HandlingImmediate:
		return "IMMEDIATE"
	case HandlingPriority:
		return "PRIORITY"
	case HandlingRoutine:
		return "ROUTINE"
	}
	return ""
}

// Code returns the single-letter code for the handling order.
func (ho HandlingOrder) Code() string {
	if ho != 0 {
		return string(ho)
	}
	return ""
}
