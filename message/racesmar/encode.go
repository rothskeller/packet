package racesmar

import (
	"fmt"
	"strings"

	"github.com/rothskeller/packet/message/common"
)

// EncodeSubject encodes the message subject.
func (f *RACESMAR) EncodeSubject() string {
	return common.EncodeSubject(f.OriginMsgID, f.Handling, Type.Tag, f.AgencyName)
}

// EncodeBody encodes the message body.
func (f *RACESMAR) EncodeBody() string {
	var (
		sb  strings.Builder
		enc *common.PIFOEncoder
	)
	if f.FormVersion == "" {
		f.FormVersion = "2.3"
	}
	if f.FormVersion == "1.6" {
		enc = common.NewPIFOEncoder(&sb, "form-oa-mutual-aid-request.html", f.FormVersion)
	} else {
		enc = common.NewPIFOEncoder(&sb, "form-oa-mutual-aid-request-v2.html", f.FormVersion)
	}
	f.StdFields.EncodeHeader(enc)
	enc.Write("15.", f.AgencyName)
	enc.Write("16a.", f.EventName)
	enc.Write("16b.", f.EventNumber)
	enc.Write("17.", f.Assignment)
	if f.FormVersion == "1.6" {
		enc.Write("18a.", f.Resources[0].Qty)
		enc.Write("18b.", f.Resources[0].RolePos)
		enc.Write("18c.", f.Resources[0].PreferredType)
		enc.Write("18d.", f.Resources[0].MinimumType)
	} else {
		for i, r := range f.Resources {
			enc.Write(fmt.Sprintf("18.%da.", i+1), r.Qty)
			enc.Write(fmt.Sprintf("18.%db.", i+1), r.RolePos)
			enc.Write(fmt.Sprintf("18.%dc.", i+1), r.PreferredType)
			enc.Write(fmt.Sprintf("18.%dd.", i+1), r.MinimumType)
			if f.FormVersion == "2.3" {
				enc.Write(fmt.Sprintf("18.%de.", i+1), r.Role)
				enc.Write(fmt.Sprintf("18.%df.", i+1), r.Position)
			}
		}
	}
	enc.Write("19a.", f.RequestedArrivalDates)
	enc.Write("19b.", f.RequestedArrivalTimes)
	enc.Write("20a.", f.NeededUntilDates)
	enc.Write("20b.", f.NeededUntilTimes)
	enc.Write("21.", f.ReportingLocation)
	enc.Write("22.", f.ContactOnArrival)
	enc.Write("23.", f.TravelInfo)
	enc.Write("24a.", f.RequestedByName)
	enc.Write("24b.", f.RequestedByTitle)
	enc.Write("24c.", f.RequestedByContact)
	enc.Write("25a.", f.ApprovedByName)
	enc.Write("25b.", f.ApprovedByTitle)
	enc.Write("25c.", f.ApprovedByContact)
	enc.Write("26a.", f.ApprovedByDate)
	enc.Write("26b.", f.ApprovedByTime)
	f.StdFields.EncodeFooter(enc)
	if err := enc.Close(); err != nil {
		panic(err)
	}
	return sb.String()
}
