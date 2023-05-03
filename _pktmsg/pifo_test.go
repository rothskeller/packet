package pktmsg

import (
	"reflect"
	"testing"
)

var parseFormTests = []struct {
	name string
	body string
	want *Form
}{
	{
		"no form",
		"Hello, world!",
		nil,
	},
	{
		"minimal valid form",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [x]\n!/ADDON!\n",
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", "x"}},
		},
	},
	{
		"form with stuff before it",
		"Hello, world!\n!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [x]\n!/ADDON!\n",
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", "x"}},
		},
	},
	{
		"form with stuff after it",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [x]\n!/ADDON!\nGoodbye, cruel world!\n",
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", "x"}},
		},
	},
	{
		"missing header",
		"#T: tt.html\n#V: 1-2\nA: [x]\n!/ADDON!\n",
		nil,
	},
	{
		"missing type",
		"!SCCoPIFO!\n#V: 1-2\nA: [x]\n!/ADDON!\n",
		nil,
	},
	{
		"invalid type",
		"!SCCoPIFO!\n#T: t\n#V: 1-2\nA: [x]\n!/ADDON!\n",
		nil,
	},
	{
		"missing version",
		"!SCCoPIFO!\n#T: tt.html\nA: [x]\n!/ADDON!\n",
		nil,
	},
	{
		"invalid version",
		"!SCCoPIFO!\n#T: tt.html\n#V: X\nA: [x]\n!/ADDON!\n",
		nil,
	},
	{
		"invalid field",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA\n!/ADDON!\n",
		nil,
	},
	{
		"missing footer",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [x]\n",
		nil,
	},
	{
		"multiple settings of same field",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [x]\nA: [x]\n!/ADDON!\n",
		nil,
	},
	{
		"bracket quoting",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [nl\\nbs\\\\rb`]et`]]]\n!/ADDON!\n",
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", "nl\nbs\\rb]et`"}},
		},
	},
	{
		"line continuation",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [this is \na test]\n!/ADDON!\n",
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", "this is a test"}},
		},
	},
}

func TestParseForm(t *testing.T) {
	for _, tt := range parseFormTests {
		t.Run(tt.name, func(t *testing.T) {
			if gotF := ParseForm(tt.body); !reflect.DeepEqual(gotF, tt.want) {
				t.Errorf("Parse() = %v, want %v", gotF, tt.want)
			}
		})
	}
}

var formEncodeTests = []struct {
	name string
	form *Form
	want string
}{
	{
		"minimal strict",
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", ""}, {"B", "b"}},
		},
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nB: [b]\n!/ADDON!\n",
	},
	{
		"line continuation",
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields: []FormField{
				{"A", "1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 "},
			},
		},
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 123\n4567890 ]\n!/ADDON!\n",
	},
}

func TestFormEncode(t *testing.T) {
	for _, tt := range formEncodeTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.form.Encode(); got != tt.want {
				t.Errorf("Form.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
