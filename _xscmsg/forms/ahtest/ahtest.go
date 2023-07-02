package ahtest

import xscmsg "github.com/rothskeller/packet/xscmsg/forms"

func init() {
	xscmsg.Register(AHFacStatValidator{})
}

type AHFacStatValidator struct {
	*xscmsg.AHFacStat
}

func (AHFacStatValidator) Validate() []string {
	println("here!")
	return []string{"oops"}
}
