package pktmsg

import (
	"reflect"
	"testing"
)

var parseFormTests = []struct {
	name   string
	body   string
	strict bool
	want   *Form
}{
	{
		"no form",
		"Hello, world!",
		false,
		nil,
	},
	{
		"minimal valid form",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA:\n!/ADDON!\n",
		false,
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", ""}},
		},
	},
	{
		"missing header",
		"#T: tt.html\n#V: 1-2\nA:\n!/ADDON!\n",
		false,
		nil,
	},
	{
		"missing type",
		"!SCCoPIFO!\n#V: 1-2\nA:\n!/ADDON!\n",
		false,
		nil,
	},
	{
		"invalid type",
		"!SCCoPIFO!\n#T: t\n#V: 1-2\nA:\n!/ADDON!\n",
		false,
		nil,
	},
	{
		"missing version",
		"!SCCoPIFO!\n#T: tt.html\nA:\n!/ADDON!\n",
		false,
		nil,
	},
	{
		"invalid version",
		"!SCCoPIFO!\n#T: tt.html\n#V: X\nA:\n!/ADDON!\n",
		false,
		nil,
	},
	{
		"invalid field",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA\n!/ADDON!\n",
		false,
		nil,
	},
	{
		"missing footer",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA:\n",
		false,
		nil,
	},
	{
		"extra stuff after footer",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA:\n!/ADDON!\nX\n",
		false,
		nil,
	},
	{
		"annotation - loose",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA.foo:\n!/ADDON!\n",
		false,
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A.", ""}},
		},
	},
	{
		"annotation - strict",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA.foo:\n!/ADDON!\n",
		true,
		nil,
	},
	{
		"multiple settings of same field",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA:\nA:\n!/ADDON!\n",
		false,
		nil,
	},
	{
		"strict quoting - brackets",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [nl\\nbs\\\\rb`]et`]]]\n!/ADDON!\n",
		true,
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", "nl\nbs\\rb]et`"}},
		},
	},
	{
		"strict quoting - comment after brackets",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [nl\\nbs\\\\rb`]et`]]] # foo`\n!/ADDON!\n",
		true,
		nil,
	},
	{
		"strict quoting - no brackets",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: nl\\nbs\\\\rb]et`\n!/ADDON!\n",
		true,
		nil,
	},
	{
		"loose quoting - brackets",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA:   [nl\\nbs\\\\rb`]et`]]]  \n!/ADDON!\n",
		false,
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", "nl\nbs\\rb]et`"}},
		},
	},
	{
		"loose quoting - no brackets",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA:   nl\\nbs\\\\rb]et`  \n!/ADDON!\n",
		false,
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", "nl\nbs\\rb]et`"}},
		},
	},
	{
		"loose quoting - no brackets - comment",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA:   nl\\nbs\\\\rb]et` # foo\n!/ADDON!\n",
		false,
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", "nl\nbs\\rb]et`"}},
		},
	},
	{
		"loose quoting - comment only",
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA:   # foo\n!/ADDON!\n",
		false,
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", ""}},
		},
	},
}

func TestParseForm(t *testing.T) {
	for _, tt := range parseFormTests {
		t.Run(tt.name, func(t *testing.T) {
			if gotF := ParseForm(tt.body, tt.strict); !reflect.DeepEqual(gotF, tt.want) {
				t.Errorf("Parse() = %v, want %v", gotF, tt.want)
			}
		})
	}
}

var formEncodeTests = []struct {
	name         string
	form         *Form
	annotations  map[string]string
	comments     map[string]string
	looseQuoting bool
	want         string
}{
	{
		"minimal strict",
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", ""}},
		},
		nil, nil,
		false,
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: []\n!/ADDON!\n",
	},
	{
		"minimal loose",
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", ""}},
		},
		nil, nil,
		true,
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: \n!/ADDON!\n",
	},
	{
		"loose needs quoting - whitespace",
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", " "}},
		},
		nil, nil,
		true,
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [ ]\n!/ADDON!\n",
	},
	{
		"loose needs quoting - brackets",
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", "[]"}},
		},
		nil, nil,
		true,
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [[`]]\n!/ADDON!\n",
	},
	{
		"annotations",
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A.", "x"}},
		},
		map[string]string{"A.": "foo"}, nil,
		false,
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA.foo: [x]\n!/ADDON!\n",
	},
	{
		"comments",
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", "x"}, {"B", ""}},
		},
		nil, map[string]string{"A": "comment1", "B": "comment2"},
		true,
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: x\nB: # comment2\n!/ADDON!\n",
	},
	{
		"alignment - below minimums",
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields:      []FormField{{"A", "x"}, {"BBBBBBBBBBB", "y"}},
		},
		nil, nil,
		true,
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA:           x\nBBBBBBBBBBB: y\n!/ADDON!\n",
	},
	{
		"alignment - mix",
		&Form{
			PIFOVersion: "1",
			FormType:    "tt.html",
			FormVersion: "2",
			Fields: []FormField{
				{"AAA", "a"},
				{"BBBB", "b"},
				{"CCCCC", "c"},
				{"DDDDDDDDDDDD", "d"},
				{"EEEEEEEEEEEEE", "e"},
				{"FFFFFFFFFFFF", "f"},
				{"GGGGGG", "g"},
			},
		},
		nil, nil,
		true,
		"!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nAAA:   a\nBBBB:  b\nCCCCC: c\nDDDDDDDDDDDD:  d\nEEEEEEEEEEEEE: e\nFFFFFFFFFFFF:  f\nGGGGGG:        g\n!/ADDON!\n",
	},
}

func TestFormEncode(t *testing.T) {
	for _, tt := range formEncodeTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.form.Encode(tt.annotations, tt.comments, tt.looseQuoting); got != tt.want {
				t.Errorf("Form.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
