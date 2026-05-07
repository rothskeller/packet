package htmlb_test

import (
	"strings"
	"testing"

	"github.com/rothskeller/packet/wppsvr/htmlb"
)

const expected = `<!DOCTYPE html><html><div a1=v1 a2 a3=v&lt;3&gt; a4=/foo/4 a5="v5 x" a6=6 class="foo bar">bar<li>baz<li a=b>&lt;br&gt;<div>q</div></div>`

func Test(t *testing.T) {
	buf := new(strings.Builder)
	h := htmlb.HTML(buf)
	h.E("div class=foo a1=v1 a2 a3=%s a4=/foo/%d", "v<3>", 4, "class=bar a5='v5 x' a6=%d", 6, false, "a7=v7").R("bar").
		E(`li>baz`).P().E(`li a=b>%s`, "<br>").E("div").R("q")
	h.Close()
	if buf.String() != expected {
		t.Errorf("wanted\n%s\ngot\n%s", expected, buf.String())
	}
}
