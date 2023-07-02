package sheltstat

import (
	"strings"

	"github.com/rothskeller/packet/message/common"
)

// EncodeSubject encodes the message subject.
func (f *SheltStat) EncodeSubject() string {
	return common.EncodeSubject(f.OriginMsgID, f.Handling, Type.Tag, f.ShelterName)
}

// EncodeBody encodes the message body.
func (f *SheltStat) EncodeBody() string {
	var (
		sb  strings.Builder
		enc *common.PIFOEncoder
	)
	if f.FormVersion == "" {
		f.FormVersion = "2.2"
	}
	enc = common.NewPIFOEncoder(&sb, "form-oa-shelter-status.html", f.FormVersion)
	f.StdFields.EncodeHeader(enc)
	enc.Write("19.", f.ReportType)
	enc.Write("32.", f.ShelterName)
	enc.Write("30.", f.ShelterType)
	enc.Write("31.", f.ShelterStatus)
	enc.Write("33a.", f.ShelterAddress)
	if f.FormVersion >= "2.2" {
		enc.Write("33b.", f.ShelterCityCode)
		enc.Write("34b.", f.ShelterCity)
	} else {
		enc.Write("33b.", f.ShelterCity)
	}
	enc.Write("33c.", f.ShelterState)
	enc.Write("33d.", f.ShelterZip)
	enc.Write("37a.", f.Latitude)
	enc.Write("37b.", f.Longitude)
	enc.Write("40a.", f.Capacity)
	enc.Write("40b.", f.Occupancy)
	enc.Write("41.", f.MealsServed)
	enc.Write("42.", f.NSSNumber)
	enc.Write("43a.", f.PetFriendly)
	enc.Write("43b.", f.BasicSafetyInspection)
	enc.Write("43c.", f.ATC20Inspection)
	enc.Write("44.", f.AvailableServices)
	enc.Write("45.", f.MOU)
	enc.Write("46.", f.FloorPlan)
	if f.FormVersion >= "2.2" {
		enc.Write("50a.", f.ManagedByCode)
		enc.Write("49a.", f.ManagedBy)
	} else {
		enc.Write("50a.", f.ManagedBy)
	}
	enc.Write("50b.", f.ManagedByDetail)
	enc.Write("51a.", f.PrimaryContact)
	enc.Write("51b.", f.PrimaryPhone)
	enc.Write("52a.", f.SecondaryContact)
	enc.Write("52b.", f.SecondaryPhone)
	enc.Write("60.", f.TacticalCallSign)
	enc.Write("61.", f.RepeaterCallSign)
	enc.Write("62a.", f.RepeaterInput)
	enc.Write("62b.", f.RepeaterInputTone)
	enc.Write("63a.", f.RepeaterOutput)
	enc.Write("63b.", f.RepeaterOutputTone)
	enc.Write("62c.", f.RepeaterOffset)
	enc.Write("70.", f.Comments)
	enc.Write("71.", f.RemoveFromList)
	f.StdFields.EncodeFooter(enc)
	if err := enc.Close(); err != nil {
		panic(err)
	}
	return sb.String()
}
