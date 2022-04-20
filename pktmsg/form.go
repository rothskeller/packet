package pktmsg

// This file defines TxForm and RxForm.

// The format of a PackItForms message is:
//     !SCCoPIFO!
//     #T: «form html»
//     #V: «version»
//     «fields»
//     !/ADDON!
// where «fields» is any number of
//     «fieldname»: [«fieldvalue»]
//
// «fieldname» can be either a word, or an integer followed by a period.
//
// Newlines in «fieldvalue» are rendered as "\n", and backslashes are rendered
// as "\\".  Close brackets "]" are rendered as "`]".  A backtick at the end of
// the field value is rendered as "`]]" (plus the third close bracket that ends
// the «fieldvalue»).  Lines are not wrapped.

import (
	"fmt"
	"regexp"
	"strings"
)

// TxForm is the foundation for all outgoing messages containing
// PackItForms-encoded forms.
type TxForm struct {
	TxMessage
	// FormHTML is the name of the PackItForms HTML file for the form.  It
	// identifies the form type.
	FormHTML string
	// FormVersion is the version number of the form.
	FormVersion string
	// Fields is an ordered list of field-name/field-value pairs, describing
	// the fields of the form.
	fields []string
}

// SetField sets the field with the specified name to the specified value.
// Setting a field to "" removes it.
func (f *TxForm) SetField(name, value string) {
	for i := 0; i < len(f.fields); i += 2 {
		if f.fields[i] == name {
			if value != "" {
				f.fields[i+1] = value
			} else {
				f.fields = append(f.fields[:i], f.fields[i+2:]...)
			}
			return
		}
	}
	if value != "" {
		f.fields = append(f.fields, name, value)
	}
}

// Encode returns the encoded subject line and body of the message.
func (f *TxForm) Encode() (subject, body string, err error) {
	var sb strings.Builder

	if f.FormHTML == "" || f.FormName == "" || f.FormVersion == "" || len(f.fields) == 0 {
		return "", "", ErrIncomplete
	}
	sb.WriteString("!SCCoPIFO!\n")
	fmt.Fprintf(&sb, "#T: %s\n#V: 3.2-%s\n", f.FormHTML, f.FormVersion)
	for i := 0; i < len(f.fields); i += 2 {
		fmt.Fprintf(&sb, "%s: [%s]\n", f.fields[i], encodePIFOValue(f.fields[i+1]))
	}
	sb.WriteString("!/ADDON!\n")
	f.Body = sb.String()
	return f.TxMessage.Encode()
}

var encodePIFOReplacer = strings.NewReplacer(`\`, `\\`, "\n", `\n`, "]", "`]")

// encodePIFOValue encodes a value for a PackItForms field line, escaping
// newlines and close brackets.
func encodePIFOValue(s string) string {
	s = encodePIFOReplacer.Replace(s)
	if len(s) != 0 && s[len(s)-1] == '`' {
		s += "]]"
	}
	return s
}

//------------------------------------------------------------------------------

// RxForm is the foundation for all received messages containing
// PackItForms-encoded forms.
type RxForm struct {
	RxMessage
	// FormHTML is the name of the PackItForms HTML file for the form.  It
	// identifies the form type.
	FormHTML string
	// PIFOVersion is the version number of the PackItForms encoding.
	PIFOVersion string
	// FormVersion is the version number of the form.
	FormVersion string
	// Fields is a map from field name to field value, containing the values
	// of the form fields.
	Fields map[string]string
	// CorruptForm is a flag indicating that the form was not encoded
	// properly.
	CorruptForm bool
}

// Form returns a pointer to the RxForm portion of a message object.  It
// can be used to reach fields of the RxForm object that are occluded by
// types that embed RxForme.
func (f *RxForm) Form() *RxForm { return f }

// TypeCode returns the machine-readable code for the message type.
func (f *RxForm) TypeCode() string {
	if f.FormName != "" {
		return f.FormName
	}
	return "UNKNOWN"
}

// TypeName returns the human-reading name of the message type.
func (f *RxForm) TypeName() string {
	if f.FormName != "" {
		return f.FormName + " form"
	}
	return "form of unknown type"
}

// TypeArticle returns "a" or "an", whichever is appropriate for the TypeName.
func (f *RxForm) TypeArticle() string {
	if f.FormName != "" {
		// We have to guess.  If it starts with an uppercase letter,
		// we'll guess it's an acronym.
		switch f.FormName[0] {
		case 'A', 'E', 'F', 'H', 'I', 'L', 'M', 'N', 'O', 'R', 'S', 'X', 'a', 'e', 'i', 'o', 'u':
			return "an"
		}
	}
	return "a"
}

// pifoFieldRE is the regular expression for a field value line in a PackItForms
// message body.
var pifoFieldRE = regexp.MustCompile(`^([-A-Za-z0-9.]+): \[(.*)\]$`)

// parseRxForm examines an RxMessage to see if it contains a PackItForms-encoded
// form, and if so, wraps it in an RxForm and returns it.  If it is not, it
// returns nil.
func parseRxForm(m *RxMessage) *RxForm {
	var (
		f             RxForm
		seenSignature bool
		seenHTMLForm  bool
		seenVersion   bool
		seenField     bool
		seenFooter    bool
	)
	// Make sure this has the signature of a PackItForms-encoded form.
	if !strings.HasPrefix(m.Body, "!SCCoPIFO!\n") {
		return nil
	}
	f.RxMessage = *m
	f.Fields = make(map[string]string)
	// Parse each line of the form message.
	for _, line := range strings.Split(m.Body, "\n") {
		switch {
		case line == "!SCCoPIFO!": // signature line
			if seenSignature {
				f.CorruptForm = true
			}
			seenSignature = true
		case strings.HasPrefix(line, "#T: "): // form type line
			if seenHTMLForm || seenField || seenFooter {
				f.CorruptForm = true
			}
			f.FormHTML = strings.TrimSpace(line[4:])
			seenHTMLForm = true
		case strings.HasPrefix(line, "#V: "): // version line
			if seenVersion || seenField || seenFooter {
				f.CorruptForm = true
			}
			if parts := strings.Split(line[4:], "-"); len(parts) == 2 {
				f.PIFOVersion = parts[0]
				f.FormVersion = parts[1]
			} else {
				f.CorruptForm = true
			}
			seenVersion = true
		case line == "!/ADDON!": // footer line
			if seenFooter {
				f.CorruptForm = true
			}
			seenFooter = true
		case line == "":
			break
		default: // anything else should be a field setting line
			match := pifoFieldRE.FindStringSubmatch(line)
			if match == nil {
				f.CorruptForm = true // not a field setting line
				break
			}
			if seenFooter {
				f.CorruptForm = true
			}
			if _, ok := f.Fields[match[1]]; ok {
				f.CorruptForm = true // multiple values for same field
			}
			f.Fields[match[1]] = decodePIFOValue(match[2])
			seenField = true
		}
	}
	if !seenSignature || !seenHTMLForm || !seenVersion || !seenFooter {
		f.CorruptForm = true
	}
	return &f
}

var decodePIFOReplacer = strings.NewReplacer(`\\`, `\`, `\n`, "\n", "`]", "]")

// decodePIFOValue decodes a value in a PackItForms field line, by unescaping
// newlines and close brackets.
func decodePIFOValue(s string) string {
	if strings.HasSuffix(s, "`]]") {
		s = s[:len(s)-2]
	}
	return decodePIFOReplacer.Replace(s)
}
