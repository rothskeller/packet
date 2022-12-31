package xscform

// Most of the functions in this file correspond to the validation checks called
// out in the "required" and "class" attributes of the various input controls in
// PackItForms HTML files.

/*
// ValidateComputedChoice handles a common pattern where the value of the target
// field is computed based on the value of another field.  Specifically, if the
// value of the other field is one of the allowed values for this field, it is
// kept, and otherwise, the last allowed value of this field is used.  This
// "validation" is used for fields that are implemented in the PackItForms HTML
// with a combo box and a "compatible_values" function call.
func ValidateComputedChoice(f *XSCForm, fd *FieldDefinition, value string, _ bool) (newval, _ string) {
	other := f.Get(fd.ComputedFromField)
	other, _ = ValidateSelect(f, fd, other, false)
	for _, allowed := range fd.Values {
		if other == allowed {
			return other, ""
		}
	}
	return fd.Values[len(fd.Values)-1], ""
}

// ValidateFrequency verifies that the value is a frequency.
func ValidateFrequency(_ *XSCForm, fd *FieldDefinition, value string, _ bool) (newval, problem string) {
	if value != "" && !frequencyRE.MatchString(value) {
		return value, fmt.Sprintf("%q is not a valid frequency value for field %q", value, fd.Tag)
	}
	return value, ""
}

// ValidateFrequencyOffset verifies that the value is a frequency offset.
func ValidateFrequencyOffset(_ *XSCForm, fd *FieldDefinition, value string, _ bool) (newval, problem string) {
	if value != "" && !frequencyOffsetRE.MatchString(value) {
		return value, fmt.Sprintf("%q is not a valid frequency offset value for field %q", value, fd.Tag)
	}
	return value, ""
}
*/
