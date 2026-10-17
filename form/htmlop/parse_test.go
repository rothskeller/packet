package htmlop

import (
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestParseKeepsRepeatedAttributes(t *testing.T) {
	doc, err := Parse(strings.NewReader(`<input type=button disabled VALUE="Submit &amp; Go" value={{label}} data-x='a>b'>`))
	if err != nil {
		t.Fatal(err)
	}
	want := []html.Attribute{
		{Key: "type", Val: "button"},
		{Key: "disabled"},
		{Key: "value", Val: "Submit & Go"},
		{Key: "value", Val: "{{label}}"},
		{Key: "data-x", Val: "a>b"},
	}
	if got := doc.FirstChild.Attr; !reflect.DeepEqual(got, want) {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}

func TestParseOrdinaryAttributes(t *testing.T) {
	// Tags without repeated attributes get the tokenizer's attributes
	// unchanged.
	const input = `<a href="x?a=1&amp;b=2" title=t>x</a><br/>`
	doc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	tz := html.NewTokenizer(strings.NewReader(input))
	tz.Next()
	if got, want := doc.FirstChild.Attr, tz.Token().Attr; !reflect.DeepEqual(got, want) {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}
