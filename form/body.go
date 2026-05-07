package form

import (
	"cmp"
	"fmt"
	"iter"
	"maps"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/form/formdef"
	"github.com/rothskeller/packet/form/pifover"
	"github.com/rothskeller/packet/message/body"
	"github.com/rothskeller/packet/message/cachetrack"
	"github.com/rothskeller/packet/message/field"
)

var (
	fieldLineRE           = regexp.MustCompile(`(?i)^([A-Z0-9][-A-Z0-9.]*): \[`)
	formRecognizeRE       = regexp.MustCompile(`(?s)^\n*!([A-Z][A-Za-z0-9_]*)!\n#T:.*\n!/ADDON!(\n|$)`)
	headerRE              = regexp.MustCompile(`^#T: ([a-z][-a-z0-9]+\.html)\n#V: (\d+(?:\.\d+)*[A-Za-z]?)-(\d+(?:\.\d+)*[A-Za-z]*)\n`)
	quoteSCCoPIFO         = strings.NewReplacer(`\`, `\\`, "\n", `\n`, "]", "`]")
	ErrInvalidAddonName   = errors.New("The form addon name must start with an uppercase letter and contain only letters, digits, and underscores.")
	ErrInvalidFormHTML    = errors.New("The form HTML name must start with a lowercase letter, contain only lowercase letters, digits, and dashes, and end with \".html\".")
	ErrInvalidPIFOVersion = errors.New("The PackItForms version number must be a dot-separated sequence of one or more non-negative integers, followed by zero or more letters.")
	ErrInvalidFormVersion = errors.New("The form version number must be a dot-separated sequence of one or more non-negative integers, followed by zero or more letters.")
	ErrInvalidFormHeader  = errors.New("The body contains a PackItForms form that cannot be decoded because it does not start with properly-formatted #T and #V lines.")
	ErrInvalidFormSyntax  = errors.New("The body contains a PackItForms form that cannot be decoded because it contains a line that does not conform to PackItForms syntax.")
)

type ErrRedundantFormField string

func (e ErrRedundantFormField) Error() string {
	return fmt.Sprintf("The body contains a PackItForms form that cannot be decoded because it has more than one entry for the %q field.", string(e))
}

// fieldOrders provides ordering for the field tags that aren't numeric.  It
// maps them to fake numeric tags that sort correctly.
var fieldOrders = map[string]string{
	"MsgNo":     "0a.",
	"DestMsgNo": "0b.",
	"ToName":    "9c.",
	"FmName":    "9d.",
	"ToTel":     "9e.",
	"FmTel":     "9f.",

	"OpRelayRcvd": "100a.",
	"OpRelaySent": "100b.",
	"Rec-Sent":    "100c.",
	"OpCall":      "100d.",
	"OpName":      "100e.",
	"Method":      "100f.",
	"Other":       "100g.",
	"OpDate":      "100h.",
	"OpTime":      "100i.",
}

// A FormBody is a message body containing a form encoded in PackItForms
// encoding.
type FormBody struct {
	cachetrack.Tracker
	body        string
	def         *formdef.FormDef
	addonName   string
	formHTML    string
	pifoVersion string
	formVersion string
	textBefore  string
	fields      map[string]string
	textAfter   string
}

var _ body.Body = (*FormBody)(nil)

// NewFormBody creates a new form message body with the specified HTML
// identifier and version numbers.  It returns an error if the arguments are
// invalid.
func NewFormBody(def *formdef.FormDef) (b *FormBody, err error) {
	b = &FormBody{
		def:         def,
		addonName:   def.AddonName,
		formHTML:    def.HTMLName,
		pifoVersion: pifover.PIFOVersion,
		formVersion: def.Version,
		fields:      make(map[string]string),
	}
	b.MarkDirty("")
	return b, nil
}

func init() { body.RegisterDecoder(decodeFormBody) }
func decodeFormBody(body string) (_ body.Body, err error) {
	var b FormBody

	if match := formRecognizeRE.FindStringSubmatchIndex(body); match == nil {
		return nil, nil
	} else {
		b.body = body
		if match[2] > 2 {
			b.textBefore = body[:match[2]-2]
		}
		b.addonName = body[match[2]:match[3]]
		body = body[match[3]+2:]
	}
	if match := headerRE.FindStringSubmatch(body); match == nil {
		return nil, ErrInvalidFormHeader
	} else {
		b.formHTML, b.pifoVersion, b.formVersion = match[1], match[2], match[3]
		body = body[len(match[0]):]
	}
	b.fields = make(map[string]string)
	for {
		var (
			match []string
			tag   string
			value string
			ok    bool
		)
		if strings.HasPrefix(body, "\n") {
			// Blank lines are allowed.  They can result from JNOS
			// inserting line breaks in the wrong place.
			body = body[1:]
			continue
		}
		if match = fieldLineRE.FindStringSubmatch(body); match == nil {
			break
		}
		tag, body = match[1], body[len(match[0]):]
		if _, ok = b.fields[tag]; ok {
			err = errors.Join(err, ErrRedundantFormField(tag))
		}
		if value, body, ok = parseBracketedValue(body); !ok {
			return nil, ErrInvalidFormSyntax
		}
		b.fields[tag] = value
	}
	if !strings.HasPrefix(body, "!/ADDON!") || (len(body) > 8 && body[8] != '\n') {
		return nil, ErrInvalidFormSyntax
	}
	if err != nil {
		return nil, err
	}
	b.textAfter = body[8:]
	if len(b.textAfter) != 0 && b.textAfter[0] == '\n' {
		b.textAfter = b.textAfter[1:]
	}
	return &b, nil
}

// parseBracketedValue parses a field value in brackets.  Within the brackets,
// \n represents a newline, \\ represents a backslash, `] represents a close
// bracket, and `]]] represents a backtick at the end of the string.  Literal
// newlines are ignored, even in the middle of the above sequences.  The final
// close bracket must be followed by a newline.
func parseBracketedValue(body string) (value, nbody string, ok bool) {
	for body != "" {
		if body[0] == ']' { // close bracket ends the value
			ok = true
			body = body[1:]
			break
		}
		if body[0] == '\n' { // newlines are ignored
			body = body[1:]
			continue
		}
		if len(body) > 1 && body[0] == '\\' && body[1] == '\\' { // escaped backslashes are a single backslash
			value += "\\"
			body = body[2:]
			continue
		}
		if len(body) > 2 && body[0] == '\\' && body[1] == '\n' && body[2] == '\\' { // escaped backslashes are a single backslash
			value += "\\"
			body = body[3:]
			continue
		}
		if len(body) > 1 && body[0] == '\\' && body[1] == 'n' { // escaped 'n's are newlines
			value += "\n"
			body = body[2:]
			continue
		}
		if len(body) > 2 && body[0] == '\\' && body[1] == '\n' && body[2] == 'n' { // escaped 'n's are newlines
			value += "\n"
			body = body[3:]
			continue
		}
		if len(body) > 3 && body[0] == '`' && body[1] == ']' && body[2] == ']' && body[3] == ']' {
			// backtick followed by three close brackets is a backtick that ends the string
			ok = true
			value += "`"
			body = body[4:]
			break
		}
		if len(body) > 4 &&
			((body[0] == '`' && body[1] == '\n' && body[2] == ']' && body[3] == ']' && body[4] == ']') ||
				(body[0] == '`' && body[1] == ']' && body[2] == '\n' && body[3] == ']' && body[4] == ']') ||
				(body[0] == '`' && body[1] == ']' && body[2] == ']' && body[3] == '\n' && body[4] == ']')) {
			// backtick followed by three close brackets is a backtick that ends the string
			ok = true
			value += "`"
			body = body[5:]
			break
		}
		if len(body) > 1 && body[0] == '`' && body[1] == ']' { // backtick, close bracket is a literal close bracket
			value += "]"
			body = body[2:]
			continue
		}
		if len(body) > 2 && body[0] == '`' && body[1] == '\n' && body[2] == ']' { // backtick, close bracket is a literal close bracket
			value += "]"
			body = body[3:]
			continue
		}
		// anything else is copied literally
		r, sz := utf8.DecodeRuneInString(body)
		value += string(r)
		body = body[sz:]
	}
	if !ok { // end of body before end of value
		return "", body, false
	}
	if body == "" || body[0] != '\n' { // no newline, or extra text after value
		return "", body, false
	}
	return value, body[1:], true
}

// AddonName returns the identifier of the Outpost addon that's supposed to
// handle this form.
func (b *FormBody) AddonName() string { return b.addonName }

// FormHTML returns the identifier of the form type (an HTML filename).
func (b *FormBody) FormHTML() string { return b.formHTML }

// PIFOVersion returns the PackItForms version number in the form.
func (b *FormBody) PIFOVersion() string { return b.pifoVersion }

// FormVersion returns the form version number in the form.
func (b *FormBody) FormVersion() string { return b.formVersion }

// TextBefore returns any text that appears before the start of the form.
func (b *FormBody) TextBefore() string { return b.textBefore }

// Fields returns an iterator of the message fields.
func (b *FormBody) Fields() iter.Seq[field.Field] {
	return func(yield func(field.Field) bool) {
		for fd := range b.def.AllFields() {
			if !yield(ff2mf{fd}) {
				return
			}
		}
	}
}

// FieldList returns a list of the field tags defined in the form.
func (b *FormBody) FieldList() []string { return slices.Collect(maps.Keys(b.fields)) }

// Field returns the value of the form field with the specified tag, or an
// empty string if the specified field is not specified in the form.
func (b *FormBody) Field(tag string) string { return b.fields[tag] }

// SetField sets the value of the form field with the specified tag.  Setting a
// field to an empty string removes it.
func (b *FormBody) SetField(tag, value string) {
	b.SetFieldR(tag, value, "form.FormBody.Field."+tag)
}

// SetFieldR sets the value of the form field with the specified tag.  Setting
// a field to an empty string removes it.  If the result is a change, the body
// is marked dirty with the specified reason.
func (b *FormBody) SetFieldR(tag, value, reason string) {
	if value == b.fields[tag] {
		return
	}
	if value == "" {
		delete(b.fields, tag)
	} else {
		b.fields[tag] = value
	}
	b.MarkDirty(reason)
}

// TextAfter returns any text that appears after the end of the form.
func (b *FormBody) TextAfter() string { return b.textAfter }

// EncodedBody returns the encoded representation of the form.
func (b *FormBody) EncodedBody() string {
	if b.Dirty() {
		var sb strings.Builder

		fmt.Fprintf(&sb, "!%s!\n#T: %s\n#V: %s-%s\n", b.addonName, b.formHTML, b.pifoVersion, b.formVersion)
		tags := slices.Collect(maps.Keys(b.fields))
		slices.SortFunc(tags, tagSort)
		for _, tag := range tags {
			encodeField(&sb, tag, b.fields[tag])
		}
		sb.WriteString("!/ADDON!\n")
		b.body = sb.String()
		b.MarkClean()
	}
	return b.body
}

// Clone returns a copy of the FormBody.
func (b *FormBody) Clone() body.Body {
	nb := *b
	nb.Tracker = cachetrack.Tracker{}
	nb.fields = maps.Clone(b.fields)
	return &nb
}

// tagSort sorts tags in the order they should appear in the PackItForms
// encoding.  It converts common non-numeric tags to fake numeric tags that
// sort correctly.  Then it extracts the initial numeric part from each tag.
// If those are different sizes, the shortest one wins.  Otherwise, the tags
// are compared with regular alphanumeric sort.
func tagSort(a, b string) int {
	if fake, ok := fieldOrders[a]; ok {
		a = fake
	}
	if fake, ok := fieldOrders[b]; ok {
		b = fake
	}
	ia := strings.IndexFunc(a, nondigit)
	ib := strings.IndexFunc(b, nondigit)
	if ia != ib {
		return cmp.Compare(ia, ib)
	}
	return cmp.Compare(a, b)
}
func nondigit(r rune) bool { return r < '0' || r > '9' }

// encodeField encodes a field tag and value in PackItForms encoding.
func encodeField(sb *strings.Builder, tag, value string) {
	if value == "" {
		return
	}
	value = quoteSCCoPIFO.Replace(value)
	if strings.HasSuffix(value, "`") {
		value += "]]"
	}
	fmt.Fprintf(sb, "%s: [%s]\n", tag, value)
}
