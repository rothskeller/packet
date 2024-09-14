package racesmar

import (
	"fmt"
	"strings"

	"github.com/rothskeller/packet/message"
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

func (r *Resource) Fields21(index int) []*message.Field {
	var presence func() (message.Presence, string)
	if index == 1 {
		presence = message.Required
	}
	return []*message.Field{
		message.NewCardinalNumberField(&message.Field{
			Label:       fmt.Sprintf("Resource %d Quantity", index),
			Value:       &r.Qty,
			Presence:    presence,
			PIFOTag:     fmt.Sprintf("18.%da.", index),
			PDFRenderer: &message.PDFTextRenderer{X: 139, Y: 332.6 + 18.6*float64(index), R: 150, H: 11, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			EditWidth:   2,
			EditHelp:    `This is the number of people needed for the role and position requested on this row.`,
			EditSkip:    r.skip,
		}),
		message.NewTextField(&message.Field{
			Label:       fmt.Sprintf("Resource %d Role/Position", index),
			Value:       &r.RolePos,
			Presence:    presence,
			PIFOTag:     fmt.Sprintf("18.%db.", index),
			PDFRenderer: &message.PDFTextRenderer{X: 170, Y: 332.6 + 18.6*float64(index), R: 263, H: 11, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   31,
			EditHelp:    `This is the role and position to be filled by the people requested on this row.`,
			EditSkip:    r.skip,
		}),
		message.NewTextField(&message.Field{
			Label:       fmt.Sprintf("Resource %d Preferred Type", index),
			Value:       &r.PreferredType,
			Presence:    presence,
			PIFOTag:     fmt.Sprintf("18.%dc.", index),
			PDFRenderer: &message.PDFTextRenderer{X: 279, Y: 332.6 + 18.6*float64(index), R: 418, H: 11, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   7,
			EditHelp:    `This is the preferred resource type (credential) for the people requested on this row.`,
			EditSkip:    r.skip,
		}),
		message.NewTextField(&message.Field{
			Label:       fmt.Sprintf("Resource %d Minimum Type", index),
			Value:       &r.MinimumType,
			Presence:    presence,
			PIFOTag:     fmt.Sprintf("18.%dd.", index),
			PDFRenderer: &message.PDFTextRenderer{X: 509, Y: 332.6 + 18.6*float64(index), R: 567, H: 11, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   7,
			EditHelp:    `This is the minimum resource type (credential) for the people requested on this row.`,
			EditSkip:    r.skip,
		}),
	}
}

var typeMap = map[string]message.ChoiceMapper{
	"Field Communicator":   message.Choices{"F1", "F2", "F3", "Type IV", "Type V"},
	"Net Control Operator": message.Choices{"N1", "N2", "N3", "Type IV", "Type V"},
	"Packet Operator":      message.Choices{"P1", "P2", "P3", "Type IV", "Type V"},
	"Shadow Communicator":  message.Choices{"S1", "S2", "S3", "Type IV", "Type V"},
}

func (r *Resource) Fields23(index int) []*message.Field {
	var qtyPresence, rolePresence, posPresence, typePresence func() (message.Presence, string)
	if index == 1 {
		qtyPresence = message.Required
		rolePresence = message.Required
		typePresence = message.Required
	} else {
		rolePresence = r.requiredIfQtyElseNotAllowed
		posPresence = r.notAllowedWithoutQty
		typePresence = r.requiredIfQtyElseNotAllowed
	}
	return []*message.Field{
		message.NewCardinalNumberField(&message.Field{
			Label:       fmt.Sprintf("Resource %d Quantity", index),
			Value:       &r.Qty,
			Presence:    qtyPresence,
			PIFOTag:     fmt.Sprintf("18.%da.", index),
			PDFRenderer: &message.PDFTextRenderer{X: 139, Y: 332.6 + 18.6*float64(index), R: 150, H: 11, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			EditWidth:   2,
			EditHelp:    `This is the number of people needed for the role and position requested on this row.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       fmt.Sprintf("Resource %d Role", index),
			Value:       &r.Role,
			Presence:    rolePresence,
			PIFOTag:     fmt.Sprintf("18.%de.", index),
			Choices:     message.Choices{"Field Communicator", "Net Control Operator", "Packet Operator", "Shadow Communicator"},
			PDFRenderer: &message.PDFTextRenderer{X: 170, Y: 332.6 + 18.6*float64(index), R: 263, H: 11, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the role of the people requested on this row.  It is required when there is a quantity on the row.`,
			EditApply: func(f *message.Field, s string) {
				*f.Value = f.Choices.ToPIFO(strings.TrimSpace(s))
				if r.Position != "" {
					r.RolePos = message.SmartJoin(r.Role, "/ "+r.Position, " ")
				}
			},
		}),
		message.NewTextField(&message.Field{
			Label:       fmt.Sprintf("Resource %d Position", index),
			Value:       &r.Position,
			Presence:    posPresence,
			PIFOTag:     fmt.Sprintf("18.%df.", index),
			PDFRenderer: &message.PDFTextRenderer{X: 279, Y: 332.6 + 18.6*float64(index), R: 418, H: 11, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   31,
			EditHelp:    `This is the position to be held by the people requested on this row.`,
			EditApply: func(_ *message.Field, s string) {
				r.Position = strings.TrimSpace(s)
				if r.Position != "" {
					r.RolePos = message.SmartJoin(r.Role, "/ "+r.Position, " ")
				}
			},
		}),
		message.NewCalculatedField(&message.Field{
			Label:   fmt.Sprintf("Resource %d Role/Position", index),
			Value:   &r.RolePos,
			PIFOTag: fmt.Sprintf("18.%db.", index),
		}),
		message.NewRestrictedField(&message.Field{
			Label:       fmt.Sprintf("Resource %d Preferred Type", index),
			Value:       &r.PreferredType,
			Presence:    typePresence,
			Choices:     r,
			PIFOTag:     fmt.Sprintf("18.%dc.", index),
			PDFRenderer: &message.PDFTextRenderer{X: 431, Y: 332.6 + 18.6*float64(index), R: 489, H: 11, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the preferred resource type (credential) for the people requested on this row.  It is required when there is a quantity on the row.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       fmt.Sprintf("Resource %d Minimum Type", index),
			Value:       &r.MinimumType,
			Presence:    typePresence,
			Choices:     r,
			PIFOTag:     fmt.Sprintf("18.%dd.", index),
			PDFRenderer: &message.PDFTextRenderer{X: 509, Y: 332.6 + 18.6*float64(index), R: 567, H: 11, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the minimum resource type (credential) for the people requested on this row.  It is required when there is a quantity on the row.`,
		}),
	}
}
func (r *Resource) requiredIfQtyElseNotAllowed() (message.Presence, string) {
	if r.Qty != "" {
		return message.PresenceRequired, "there is a quantity for the resource"
	} else {
		return message.PresenceNotAllowed, "there is no quantity for the resource"
	}
}
func (r *Resource) notAllowedWithoutQty() (message.Presence, string) {
	if r.Qty == "" {
		return message.PresenceNotAllowed, "there is no quantity for the resource"
	}
	return message.PresenceOptional, ""
}
func (r *Resource) skip(*message.Field) bool {
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
