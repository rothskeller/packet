package analyze

func init() {
	Problems[ProbMultipleMessagesFromAddress.Code] = ProbMultipleMessagesFromAddress
}

// ProbMultipleMessagesFromAddress is raised when multiple messages are received
// from the same address.  It is not raised by the analyze code; it's raised by
// the reporting code.  But it's convenient to have it defined the same way as
// the other problems.
var ProbMultipleMessagesFromAddress = &Problem{
	Code:   "MultipleMessagesFromAddress",
	detect: func(*Analysis) (bool, string) { return false, "" },
}
