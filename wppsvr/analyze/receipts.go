package analyze

func init() {
	Problems[ProbDeliveryReceipt.Code] = ProbDeliveryReceipt
	Problems[ProbReadReceipt.Code] = ProbReadReceipt
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
