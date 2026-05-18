// Package pifover defines the PIFOVersion constant.
package pifover

// PIFOVersion is the PackItForms engine version number.  It must be modified
// whenever there is a change to how forms are defined or generated.  (Usually
// this would be because of a change in the *.form file format, or the variable
// definitions provided to form-*.html files.)
//
// For historical reasons, the PIFOVersion is encoded at the beginning of the
// #V: line in every message, and must therefore be in <number><dot><number>
// format.  However, in this context it is not used by anything.  In
// particular, there's no ordering semantic to the value.
//
// Its more significant use is in the pathname of the forms directory, and in
// the URL for fetching forms updates.  Because of those uses, each different
// PIFOVersion has its own independent set of forms.
const (
	PIFOVersion      = "4.0"
	PIFOVersionMajor = 4
	PIFOVersionMinor = 0
	FormsDirVersion  = "4.0.9"
)
