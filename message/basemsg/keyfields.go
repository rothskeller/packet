package basemsg

import "github.com/rothskeller/packet/message"

// TODO defining these functions on BaseMessage means that only messages for
// which they're appropriate can leverage BaseMessage.  Before BaseMessage can
// be extended to all message types, code using these methods will need to
// switch to using Base().F* instead.

// GetHandling returns the handling order of the message.
func (bm *BaseMessage) GetHandling() string { return *bm.FHandling }

// SetHandling sets the handling order of the message.
func (bm *BaseMessage) SetHandling(value string) { *bm.FHandling = value }

// GetOriginID returns the origin message ID of the message.
func (bm *BaseMessage) GetOriginID() string { return *bm.FOriginMsgID }

// SetOriginID sets the origin message ID of the message.
func (bm *BaseMessage) SetOriginID(value string) { *bm.FOriginMsgID = value }

// GetSubject returns the subject of the message (not including
// message ID, handling, etc.).
func (bm *BaseMessage) GetSubject() string { return *bm.FSubject }

// KeyFields returns a structure containing the values of certain key
// fields of the message that are needed for message analysis.
func (bm *BaseMessage) KeyFields() *message.KeyFields {
	return &message.KeyFields{
		PIFOVersion:   bm.PIFOVersion,
		FormVersion:   bm.Form.Version,
		ToICSPosition: *bm.FToICSPosition,
		ToLocation:    *bm.FToLocation,
		OpCall:        *bm.FOpCall,
	}
}
