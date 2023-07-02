package xscmsg

import (
	"fmt"

	"github.com/rothskeller/packet/xscmsg/forms/pifo"
)

// StdHeader holds the standard form header fields.
type StdHeader struct {
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
}

// PullTags pulls the standard form header fields from the tags map.
func (sh *StdHeader) PullTags(tags map[string]string) {
	sh.OriginMsgID = PullTag(tags, "MsgNo")
	sh.DestinationMsgID = PullTag(tags, "DestMsgNo")
	sh.MessageDate = PullTag(tags, "1a.")
	sh.MessageTime = PullTag(tags, "1b.")
	sh.Handling = PullTag(tags, "5.")
	sh.ToICSPosition = PullTag(tags, "7a.")
	sh.ToLocation = PullTag(tags, "7b.")
	sh.ToName = PullTag(tags, "7c.")
	sh.ToContact = PullTag(tags, "7d.")
	sh.FromICSPosition = PullTag(tags, "8a.")
	sh.FromLocation = PullTag(tags, "8b.")
	sh.FromName = PullTag(tags, "8c.")
	sh.FromContact = PullTag(tags, "8d.")
}

// FromTags populates the fields from a tag-value map.  It removes all matched
// tags from the map.
func (sh *StdHeader) FromTags(tags map[string]string) {
	for tag, value := range tags {
		switch tag {
		case "MsgNo":
			sh.OriginMsgID = value
		case "DestMsgNo":
			sh.DestinationMsgID = value
		case "1a.":
			sh.MessageDate = value
		case "1b.":
			sh.MessageTime = value
		case "5.":
			sh.Handling = value
		case "7a.":
			sh.ToICSPosition = value
		case "7b.":
			sh.ToLocation = value
		case "7c.":
			sh.ToName = value
		case "7d.":
			sh.ToContact = value
		case "8a.":
			sh.FromICSPosition = value
		case "8b.":
			sh.FromLocation = value
		case "8c.":
			sh.FromName = value
		case "8d.":
			sh.FromContact = value
		default:
			continue
		}
		delete(tags, tag)
	}
}

// EncodeBody encodes the standard form header fields to the supplied PIFO
// encoder.
func (sh *StdHeader) EncodeBody(enc *pifo.Encoder) {
	enc.Write("MsgNo", sh.OriginMsgID)
	enc.Write("DestMsgNo", sh.DestinationMsgID)
	enc.Write("1a.", sh.MessageDate)
	enc.Write("1b.", sh.MessageTime)
	enc.Write("5.", sh.Handling)
	enc.Write("7a.", sh.ToICSPosition)
	enc.Write("8a.", sh.FromICSPosition)
	enc.Write("7b.", sh.ToLocation)
	enc.Write("8b.", sh.FromLocation)
	enc.Write("7c.", sh.ToName)
	enc.Write("8c.", sh.FromName)
	enc.Write("7d.", sh.ToContact)
	enc.Write("8d.", sh.FromContact)
}

// StdFooter holds the standard form footer fields.
type StdFooter struct {
	OpRelayRcvd   string
	OpRelaySent   string
	OpName        string
	OpCall        string
	OpDate        string
	OpTime        string
	UnknownFields map[string]string
}

// PullTags pulls the standard footer fields from the tags map.
func (sf *StdFooter) PullTags(tags map[string]string) {
	sf.OpRelayRcvd = PullTag(tags, "OpRelayRcvd")
	sf.OpRelaySent = PullTag(tags, "OpRelaySent")
	sf.OpName = PullTag(tags, "OpName")
	sf.OpCall = PullTag(tags, "OpCall")
	sf.OpDate = PullTag(tags, "OpDate")
	sf.OpTime = PullTag(tags, "OpTime")
}

// FromTags populates the fields from a tag-value map.  The map is consumed.
func (sf *StdFooter) FromTags(tags map[string]string) {
	for tag, value := range tags {
		switch tag {
		case "OpRelayRcvd":
			sf.OpRelayRcvd = value
		case "OpRelaySent":
			sf.OpRelaySent = value
		case "OpName":
			sf.OpName = value
		case "OpCall":
			sf.OpCall = value
		case "OpDate":
			sf.OpDate = value
		case "OpTime":
			sf.OpTime = value
		default:
			continue
		}
		delete(tags, tag)
	}
	sf.UnknownFields = tags
}

// EncodeBody encodes the standard form footer fields to the supplied PIFO
// encoder.
func (sf *StdFooter) EncodeBody(enc *pifo.Encoder) {
	enc.Write("OpRelayRcvd", sf.OpRelayRcvd)
	enc.Write("OpRelaySent", sf.OpRelaySent)
	enc.Write("OpName", sf.OpName)
	enc.Write("OpCall", sf.OpCall)
	enc.Write("OpDate", sf.OpDate)
	enc.Write("OpTime", sf.OpTime)
	for tag, value := range sf.UnknownFields {
		enc.Write(tag, value)
	}
}

// PullTag returns the value of the specified tag in the supplied tag map, and
// removes the tag from the map.  If the tag does not exist in the map, PullTag
// returns an empty string.
func PullTag(tags map[string]string, tag string) (value string) {
	value = tags[tag]
	delete(tags, tag)
	return value
}

// LeftoverTagProblems translates each tag remaining in the tags map into a
// problem string reporting the unexpected tag.
func LeftoverTagProblems(formTag, version string, tags map[string]string) (problems []string) {
	for tag := range tags {
		problems = append(problems, fmt.Sprintf(
			"This form contains a field %q, which is not defined for %s version %s.", tag, formTag, version))
	}
	return problems
}
