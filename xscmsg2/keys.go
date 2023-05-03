package pktmsg

const (
	// FBBSRxDate is the time the message was received by the BBS.  It is
	// present only on instantly-received messages; it is not persisted in
	// local storage.
	FBBSRxDate FieldKey = "BBSRXDATE"
	// FBody is the body of the message.
	FBody FieldKey = "BODY"
	// FFormHTML is the PackItForms HTML file for the form.  It is present
	// only on form messages.
	FFormHTML FieldKey = "FORMHTML"
	// FFormVersion is the form version number.  It is present only on form
	// messages.
	FFormVersion FieldKey = "FORMVERSION"
	// FFromAddr is the origin address (From: header).  It may contain a
	// name as well as an address.
	FFromAddr FieldKey = "FROMADDR"
	// FNotPlainText indicates by its presence that the instantly-received
	// message was not in plain text.  It is not persisted in local storage.
	FNotPlainText FieldKey = "NOTPLAINTEXT"
	// FOutpostFlags contains the Outpost message flags.
	FOutpostFlags FieldKey = "OUTPOSTFLAGS"
	// FPIFOVersion is the PackItForms encoding version number.  It is
	// present only on form messages.
	FPIFOVersion FieldKey = "PIFOVERSION"
	// FReturnAddr is the return address of the message.  It is present only
	// on instantly-received messages; it is not persisted in local storage.
	FReturnAddr FieldKey = "RETURNADDR"
	// FRxArea is the BBS bulletin area from which the message was
	// retrieved.  It is present only on received bulletin messages.
	FRxArea FieldKey = "RXAREA"
	// FRxBBS is the BBS from which the message was retrieved.  It is
	// present only on received messages.
	FRxBBS FieldKey = "RXBBS"
	// FRxDate is the time the message was received locally.  It is present
	// only on received messages.
	FRxDate FieldKey = "RXDATE"
	// FSentDate is the time sent (Date: header).
	FSentDate FieldKey = "SENTDATE"
	// FSubject is the message subject (Subject: header).
	FSubject FieldKey = "SUBJECT"
	// FToAddrs is the list of destination addresses (To: header).
	FToAddrs FieldKey = "TOADDRS"
)
