package racesmar

import (
	"fmt"

	"github.com/rothskeller/packet/message/basemsg"
	"github.com/rothskeller/packet/message/common"
)

// A Resource is the description of a single resource in a RACES mutual aid
// request form.
type Resource struct {
	Qty           string
	Role          string // Added in v2.3
	Position      string // Added in v2.3
	RolePos       string
	PreferredType string
	MinimumType   string
}

func (r *Resource) Fields21(index int) []*basemsg.Field {
	var presence func() (basemsg.Presence, string)
	if index == 1 {
		presence = basemsg.Required
	}
	return []*basemsg.Field{
		{
			Label:     fmt.Sprintf("Resource %d Quantity", index),
			Value:     &r.Qty,
			Presence:  presence,
			PIFOTag:   fmt.Sprintf("18.%da.", index),
			PIFOValid: basemsg.ValidCardinal,
			Compare:   common.CompareCardinal,
			PDFMap:    basemsg.PDFName(fmt.Sprintf("Qty%d", index)),
			EditWidth: 2,
			EditHelp:  `This is the number of people needed for the role and position requested on this row.`,
			EditApply: basemsg.ApplyCardinal,
			EditSkip:  r.skip,
		},
		{
			Label:     fmt.Sprintf("Resource %d Role/Position", index),
			Value:     &r.RolePos,
			Presence:  presence,
			PIFOTag:   fmt.Sprintf("18.%db.", index),
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName(fmt.Sprintf("Position%d", index)),
			EditWidth: 31,
			EditHelp:  `This is the role and position to be filled by the people requested on this row.`,
			EditSkip:  r.skip,
		},
		{
			Label:     fmt.Sprintf("Resource %d Preferred Type", index),
			Value:     &r.PreferredType,
			Presence:  presence,
			PIFOTag:   fmt.Sprintf("18.%dc.", index),
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName(fmt.Sprintf("Pref%d", index)),
			EditWidth: 7,
			EditHelp:  `This is the preferred resource type (credential) for the people requested on this row.`,
			EditSkip:  r.skip,
		},
		{
			Label:     fmt.Sprintf("Resource %d Minimum Type", index),
			Value:     &r.MinimumType,
			Presence:  presence,
			PIFOTag:   fmt.Sprintf("18.%dd.", index),
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName(fmt.Sprintf("Min%d", index)),
			EditWidth: 7,
			EditHelp:  `This is the minimum resource type (credential) for the people requested on this row.`,
			EditSkip:  r.skip,
		},
	}
}

var typeMap = map[string]basemsg.ChoiceMapper{
	"Field Communicator":   basemsg.Choices{"F1", "F2", "F3", "Type IV", "Type V"},
	"Net Control Operator": basemsg.Choices{"N1", "N2", "N3", "Type IV", "Type V"},
	"Packet Operator":      basemsg.Choices{"P1", "P2", "P3", "Type IV", "Type V"},
	"Shadow Communicator":  basemsg.Choices{"S1", "S2", "S3", "Type IV", "Type V"},
}

func (r *Resource) Fields23(index int) []*basemsg.Field {
	var qtyPresence, rolePresence, posPresence, typePresence func() (basemsg.Presence, string)
	if index == 1 {
		qtyPresence = basemsg.Required
		rolePresence = basemsg.Required
		typePresence = basemsg.Required
	} else {
		rolePresence = r.requiredIfQtyElseNotAllowed
		posPresence = r.notAllowedWithoutQty
		typePresence = r.requiredIfQtyElseNotAllowed
	}
	return []*basemsg.Field{
		{
			Label:     fmt.Sprintf("Resource %d Quantity", index),
			Value:     &r.Qty,
			Presence:  qtyPresence,
			PIFOTag:   fmt.Sprintf("18.%da.", index),
			PIFOValid: basemsg.ValidCardinal,
			Compare:   common.CompareCardinal,
			PDFMap:    basemsg.PDFName(fmt.Sprintf("Qty%d", index)),
			EditWidth: 2,
			EditHelp:  `This is the number of people needed for the role and position requested on this row.`,
			EditApply: basemsg.ApplyCardinal,
		},
		{
			Label:     fmt.Sprintf("Resource %d Role", index),
			Value:     &r.Role,
			Presence:  rolePresence,
			PIFOTag:   fmt.Sprintf("18.%de.", index),
			PIFOValid: basemsg.ValidRestricted,
			Choices:   basemsg.Choices{"Field Communicator", "Net Control Operator", "Packet Operator", "Shadow Communicator"},
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName(fmt.Sprintf("Role%d", index)),
			EditWidth: 19,
			EditHelp:  `This is the role of the people requested on this row.  It is required when there is a quantity on the row.`,
		},
		{
			Label:     fmt.Sprintf("Resource %d Position", index),
			Value:     &r.Position,
			Presence:  posPresence,
			PIFOTag:   fmt.Sprintf("18.%df.", index),
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName(fmt.Sprintf("Position%d", index)),
			EditWidth: 31,
			EditHelp:  `This is the position to be held by the people requested on this row.`,
		},
		{
			Label:      fmt.Sprintf("Resource %d Role/Position", index),
			Value:      &r.RolePos,
			PIFOTag:    fmt.Sprintf("18.%db.", index),
			TableValue: basemsg.OmitFromTable,
		},
		{
			Label:     fmt.Sprintf("Resource %d Preferred Type", index),
			Value:     &r.PreferredType,
			Presence:  typePresence,
			Choices:   r,
			PIFOTag:   fmt.Sprintf("18.%dc.", index),
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName(fmt.Sprintf("Pref%d", index)),
			EditWidth: 7,
			EditHelp:  `This is the preferred resource type (credential) for the people requested on this row.  It is required when there is a quantity on the row.`,
		},
		{
			Label:     fmt.Sprintf("Resource %d Minimum Type", index),
			Value:     &r.MinimumType,
			Presence:  typePresence,
			Choices:   r,
			PIFOTag:   fmt.Sprintf("18.%dd.", index),
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName(fmt.Sprintf("Min%d", index)),
			EditWidth: 7,
			EditHelp:  `This is the minimum resource type (credential) for the people requested on this row.  It is required when there is a quantity on the row.`,
		},
	}
}
func (r *Resource) requiredIfQtyElseNotAllowed() (basemsg.Presence, string) {
	if r.Qty != "" {
		return basemsg.PresenceRequired, "there is a quantity for the resource"
	} else {
		return basemsg.PresenceNotAllowed, "there is no quantity for the resource"
	}
}
func (r *Resource) notAllowedWithoutQty() (basemsg.Presence, string) {
	if r.Qty == "" {
		return basemsg.PresenceNotAllowed, "there is no quantity for the resource"
	}
	return basemsg.PresenceOptional, ""
}
func (r *Resource) skip() bool {
	return r.Qty == ""
}

// Implement ChoiceMapper for Resource, providing the choices for the Preferred
// Type and Minimum Type fields based on the value of the Role field.

func (r *Resource) IsHuman(s string) bool {
	if cm := typeMap[r.Role]; cm != nil {
		return cm.IsHuman(s)
	}
	return false
}
func (r *Resource) IsPIFO(s string) bool {
	if cm := typeMap[r.Role]; cm != nil {
		return cm.IsPIFO(s)
	}
	return false
}
func (r *Resource) ToHuman(s string) string {
	if cm := typeMap[r.Role]; cm != nil {
		return cm.ToHuman(s)
	}
	return s
}
func (r *Resource) ToPIFO(s string) string {
	if cm := typeMap[r.Role]; cm != nil {
		return cm.ToPIFO(s)
	}
	return s
}
func (r *Resource) ListHuman() []string {
	if cm := typeMap[r.Role]; cm != nil {
		return cm.ListHuman()
	}
	return nil
}
