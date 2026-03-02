package field

import "k8s.io/apimachinery/pkg/util/sets"

// Values for Field.Common().
const (
	CDefaultBody          = "defaultBody"
	CDestinationMessageID = "destinationMessageID"
	CFormDate             = "formDate"
	CFromContact          = "fromContact"
	CFromICSPosition      = "fromICSPosition"
	CFromLocation         = "fromLocation"
	CFromName             = "fromName"
	CHandling             = "handling"
	CHeaderDate           = "headerDate"
	CHeaderFrom           = "headerFrom"
	CHeaderReceived       = "headerReceived"
	CHeaderTo             = "headerTo"
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
	CToContact            = "toContact"
	CToICSPosition        = "toICSPosition"
	CToLocation           = "toLocation"
	CToName               = "toName"
	CUseTactical          = "useTactical"
)

// CommonTags contains the allowed values for Field.Common.
var CommonTags = sets.New(
	CDefaultBody,
	CDestinationMessageID,
	CFormDate,
	CFromContact,
	CFromICSPosition,
	CFromLocation,
	CFromName,
	CHandling,
	CHeaderDate,
	CHeaderFrom,
	CHeaderReceived,
	CHeaderTo,
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
	CToContact,
	CToICSPosition,
	CToLocation,
	CToName,
	CUseTactical,
)
