package sheltstat

import (
	"time"

	"github.com/rothskeller/packet/message/common"
)

// New creates a new OA shelter status form with default values.
func New(opcall, opname string) *SheltStat {
	return &SheltStat{
		StdFields: common.StdFields{
			FormVersion: "2.2",
			MessageDate: time.Now().Format("01/02/2006"),
			Handling:    "PRIORITY",
			OpCall:      opcall,
			OpName:      opname,
		},
	}
}
