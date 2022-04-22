package config

import "steve.rothskeller.net/packet/pktmsg"

// ValidMessageTypes contains an empty message of each type that can be used
// for packet check-ins.
var ValidMessageTypes = []pktmsg.ParsedMessage{
	new(pktmsg.RxMessage),
	new(pktmsg.RxICS213Form),
	new(pktmsg.RxAHFacStatForm),
	new(pktmsg.RxEOC213RRForm),
	new(pktmsg.RxMuniStatForm),
	new(pktmsg.RxRACESMARForm),
	new(pktmsg.RxSheltStatForm),
}

// LookupMessageType finds the message type with the specified code, if it
// exists in ValidMessageTypes.
func LookupMessageType(code string) pktmsg.ParsedMessage {
	for _, msg := range ValidMessageTypes {
		if msg.TypeCode() == code {
			return msg
		}
	}
	return nil
}
