package common

// StdFields holds the standard form header and footer fields.
type StdFields struct {
	PIFOVersion      string
	FormVersion      string
	OriginMsgID      string
	DestinationMsgID string
	MessageDate      string
	MessageTime      string
	Handling         string
	ToICSPosition    string
	ToLocation       string
	ToName           string
	ToContact        string
	FromICSPosition  string
	FromLocation     string
	FromName         string
	FromContact      string
	OpRelayRcvd      string
	OpRelaySent      string
	OpName           string
	OpCall           string
	OpDate           string
	OpTime           string
}
