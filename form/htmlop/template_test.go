package htmlop

import (
	"bytes"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

var expandTests = []struct {
	input    string
	expected string
}{
	{
		"<div>Hello, world</div>",
		"<div>Hello, world</div>",
	},
	{
		"<def:x><div>Hello, world</div></def:x><x></x>",
		"<div>Hello, world</div>",
	},
	{
		"<def:x><div>Hello, {{var}}</div></def:x><x var=world></x>",
		"<div>Hello, world</div>",
	},
	{
		"<def:x><div>Hello<if:var>, {{var}}</if:var></div></def:x><x var=world></x>",
		"<div>Hello, world</div>",
	},
	{
		"<def:x><div>Hello<if:var>, {{var}}</if:var><ifnot:var>.</ifnot:var></div></def:x><x var=world></x>",
		"<div>Hello, world</div>",
	},
	{
		"<def:x><div>Hello<if:var>, {{var}}</if:var><ifnot:var>.</ifnot:var></div></def:x><x></x>",
		"<div>Hello.</div>",
	},
	{
		"<def:x><div>Hello, <for:var>{{.}}</for:var></div></def:x><x var='w ; o;r;l; d'></x>",
		"<div>Hello, world</div>",
	},
	{
		"<def:x><div>Hello, world<for:novar>{{.}}</for:novar></div></def:x><x></x>",
		"<div>Hello, world</div>",
	},
	{
		"<def:x><div class={{var}}>Hello, world</div></def:x><x var=foo></x>",
		`<div class="foo">Hello, world</div>`,
	},
	{
		"<def:x><div class={{var}}>Hello, world</div></def:x><x></x>",
		`<div>Hello, world</div>`,
	},
	{
		"<def:x><div class=default class={{var}}>Hello, world</div></def:x><x var=foo></x>",
		`<div class="foo">Hello, world</div>`,
	},
	{
		"<def:x><div class=default class={{var}}>Hello, world</div></def:x><x></x>",
		`<div class="default">Hello, world</div>`,
	},
	{
		"<def:x><input type=checkbox checked={{var}}></def:x><x var></x>",
		`<input type="checkbox" checked=""/>`,
	},
	{
		"<def:x><input type=checkbox checked={{var}}></def:x><x></x>",
		`<input type="checkbox"/>`,
	},
	{
		"<def:x><div {{required}}>Hello, world</div></def:x><x required=yes></x>",
		`<div required="">Hello, world</div>`,
	},
	{
		"<def:x><div {{required}}>Hello, world</div></def:x><x required></x>",
		`<div required="">Hello, world</div>`,
	},
	{
		"<def:x><div {{required}}>Hello, world</div></def:x><x></x>",
		`<div>Hello, world</div>`,
	},
	{
		"<def:x><div>Hello, world {{var++}} {{var++}}</div></def:x><x var=a></x>",
		`<div>Hello, world a a</div>`,
	},
	{
		"<def:x><div>Hello, world {{var++}} {{var++}}</div></def:x><x var=1></x>",
		`<div>Hello, world 1 2</div>`,
	},
	{
		"<def:x>{{i++}}</def:x><def:y><x></x><x></x></def:y><y i=2></y>",
		"23",
	},
	{
		"<def:x><div>Hello, <slot></slot>!</div></def:x><x>world</x>",
		"<div>Hello, world!</div>",
	},
}

func TestExpand(t *testing.T) {
	for _, tt := range expandTests {
		input, err := html.Parse(strings.NewReader(tt.input))
		if err != nil {
			t.Fatal(err)
		}
		Expand(input, nil)
		var buf bytes.Buffer
		err = html.Render(&buf, input)
		if err != nil {
			t.Fatal(err)
		}
		actual := buf.String()
		actual = strings.TrimPrefix(actual, "<html><head></head><body>")
		actual = strings.TrimSuffix(actual, "</body></html>")
		if actual != tt.expected {
			t.Errorf("Output mismatch:\n    input:    %s\n    expected: %s\n    actual:   %s\n", tt.input, tt.expected, actual)
		}
	}
}
