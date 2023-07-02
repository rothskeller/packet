package main

import (
	"strings"

	xscmsg "github.com/rothskeller/packet/xscmsg/forms"
	_ "github.com/rothskeller/packet/xscmsg/forms/ahtest"
)

func main() {
	ah := new(xscmsg.ICS213)
	probs := xscmsg.Validate(ah)
	println(strings.Join(probs, ", "))
}
