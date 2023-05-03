package xscmsg

// MessageSeverity is the SCCo-standard message severity code.
type MessageSeverity byte

// Values for MessageSeverity:
const (
	SeverityEmergency MessageSeverity = 'E'
	SeverityUrgent    MessageSeverity = 'U'
	SeverityOther     MessageSeverity = 'O'
)

// ParseSeverity parses an SCCo-standard message severity name or code.  It
// returns false if the input string is not a valid severity name or code.
func ParseSeverity(s string) (MessageSeverity, bool) {
	switch s {
	case "E", "EMERGENCY":
		return SeverityEmergency, true
	case "U", "URGENT":
		return SeverityUrgent, true
	case "O", "OTHER":
		return SeverityOther, true
	}
	return 0, false
}

// String returns the full string form of the severity.
func (ms MessageSeverity) String() string {
	switch ms {
	case SeverityEmergency:
		return "EMERGENCY"
	case SeverityUrgent:
		return "URGENT"
	case SeverityOther:
		return "OTHER"
	}
	return ""
}

// Code returns the single-letter code for the severity.
func (ms MessageSeverity) Code() string {
	if ms != 0 {
		return string(ms)
	}
	return ""
}
