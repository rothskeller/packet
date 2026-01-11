package field

import "k8s.io/apimachinery/pkg/util/sets"

// Values for Field.Common().
const (
	CDefaultBody          = "defaultBody"
	CDestinationMessageID = "destinationMessageID"
	CFormDate             = "formDate"
	CFromICSPosition      = "fromICSPosition"
	CFromLocation         = "fromLocation"
	CHandling             = "handling"
	CMessageDate          = "messageDate"
	CMessageSummary       = "messageSummary"
	CMessageTime          = "messageTime"
	COperatorCall         = "operatorCall"
	COperatorDate         = "operatorDate"
	COperatorMethod       = "operatorMethod"
	COperatorMethodOther  = "operatorMethodOther"
	COperatorName         = "operatorName"
	COperatorTime         = "operatorTime"
	COriginMessageID      = "originMessageID"
	CReceiverSender       = "receiverSender"
	CReference            = "reference"
	CSubjectFormTag       = "subjectFormTag"
	CSubjectHandling      = "subjectHandling"
	CSubjectMessageID     = "subjectMessageID"
	CSubjectSummary       = "subjectSummary"
	CTacticalCall         = "tacticalCall"
	CTacticalName         = "tacticalName"
	CToICSPosition        = "toICSPosition"
	CToLocation           = "toLocation"
	CUseTactical          = "useTactical"
)

// CommonTags contains the allowed values for Field.Common.
var CommonTags = sets.New(
	CDefaultBody,
	CDestinationMessageID,
	CFormDate,
	CFromICSPosition,
	CFromLocation,
	CHandling,
	CMessageDate,
	CMessageSummary,
	CMessageTime,
	COperatorCall,
	COperatorDate,
	COperatorMethod,
	COperatorMethodOther,
	COperatorName,
	COperatorTime,
	COriginMessageID,
	CReceiverSender,
	CReference,
	CSubjectFormTag,
	CSubjectHandling,
	CSubjectMessageID,
	CSubjectSummary,
	CTacticalCall,
	CTacticalName,
	CToICSPosition,
	CToLocation,
	CUseTactical,
)
