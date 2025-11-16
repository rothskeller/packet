package formdef

var typeHandlers = map[string]any{
	"addressList":    nil,
	"cardinalNumber": nil,
	"checkbox":       nil,

	// checkboxGroup should display the (child) labels of the selected items
	// in the group, separated by comma and space.  If presence is required,
	// one of them must be selected.
	"checkboxGroup": nil,

	"date": nil,

	// dateTime displays and edits the date and time together.
	"dateTime": nil,

	"fccCallSign":     nil,
	"frequency":       nil,
	"frequencyOffset": nil,

	// join should display the values of the children as specified in its
	// value parameter, which has "{1}" references to the children plus
	// other content.  If any "{1}" reference is empty, other than at the
	// start of the string, characters between it and the preceding "{1}"
	// reference (or start of string) are omitted.  If a "{1}" reference at
	// the start of the string is empty, characters between it // and the
	// next "{1}" reference are omitted.  As a special case, when only one
	// child has a value, display is delegated to that child.  If the join
	// has no value parameter, the values of non-empty children are joined
	// with a single space.
	"join": nil,

	"messageID":        nil,
	"multiline":        nil,
	"password":         nil,
	"phoneNumber":      nil,
	"realNumber":       nil,
	"restricted":       nil,
	"static":           nil,
	"tacticalCallSign": nil,
	"text":             nil,
	"time":             nil,
	"zipCode":          nil,
}
