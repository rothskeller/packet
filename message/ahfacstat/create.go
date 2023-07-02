package ahfacstat

import (
	"time"

	"github.com/rothskeller/packet/message/common"
)

// New creates a new allied health facility status form with default values.
func New(opcall, opname string) *AHFacStat {
	return &AHFacStat{
		StdFields: common.StdFields{
			FormVersion: "2.3",
			MessageDate: time.Now().Format("01/02/2006"),
			Handling:    "ROUTINE",
			OpCall:      opcall,
			OpName:      opname,
		},
		Date: time.Now().Format("01/02/2006"),
	}
}
