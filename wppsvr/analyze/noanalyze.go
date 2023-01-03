package analyze

func init() {
	Problems[ProbBounceMessage.Code] = ProbBounceMessage
	Problems[ProbDeliveryReceipt.Code] = ProbDeliveryReceipt
	Problems[ProbMessageCorrupt.Code] = ProbMessageCorrupt
	Problems[ProbMultipleMessagesFromAddress.Code] = ProbMultipleMessagesFromAddress
	Problems[ProbReadReceipt.Code] = ProbReadReceipt
}

// ProbMessageCorrupt is raised for messages that can't be parsed as valid email
// messages.
var ProbMessageCorrupt = &Problem{
	Code: "MessageCorrupt",
	// detection is hard-coded in Analyze
}

// ProbBounceMessage is raised for messages that appear to be auto-responses.
var ProbBounceMessage = &Problem{
	Code: "BounceMessage",
	// detection is hard-coded in Analyze
}

// ProbDeliveryReceipt is raised for any delivery receipt message.  This check
// has the side effect of determining the message type and setting a.xsc.
var ProbDeliveryReceipt = &Problem{
	Code: "DeliveryReceipt",
	// detection is hard-coded in Analyze
}

// ProbReadReceipt is raised for any READ receipt message.
var ProbReadReceipt = &Problem{
	Code: "ReadReceipt",
	// detection is hard-coded in Analyze
}

// ProbMultipleMessagesFromAddress is raised when multiple messages are received
// from the same address.  It is not raised by the analyze code; it's raised by
// the reporting code.  But it's convenient to have it defined the same way as
// the other problems.
var ProbMultipleMessagesFromAddress = &Problem{
	Code: "MultipleMessagesFromAddress",
	// detection happens in the reporting code.
}
