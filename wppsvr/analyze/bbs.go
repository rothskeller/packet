package analyze

func init() {
	Problems[ProbFromBBSDown.Code] = ProbFromBBSDown
	Problems[ProbToBBS.Code] = ProbToBBS
	Problems[ProbToBBSDown.Code] = ProbToBBSDown
}

// ProbFromBBSDown is raised when a message is sent from a BBS that is simulated
// down.
var ProbFromBBSDown = &Problem{
	Code:  "FromBBSDown",
	ifnot: []*Problem{ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		for _, down := range a.session.DownBBSes {
			if down == a.msg.FromBBS() {
				return true, ""
			}
		}
		return false, ""
	},
}

// ProbToBBSDown is raised when a message is sent to a BBS that is simulated
// down.
var ProbToBBSDown = &Problem{
	Code:  "ToBBSDown",
	ifnot: []*Problem{ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		for _, down := range a.session.DownBBSes {
			if down == a.toBBS {
				return true, ""
			}
		}
		return false, ""
	},
}

// ProbToBBS is raised when a message is sent to a BBS that is not a correct BBS
// for the session.
var ProbToBBS = &Problem{
	Code:  "ToBBS",
	ifnot: []*Problem{ProbToBBSDown, ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		for _, to := range a.session.ToBBSes {
			if to == a.toBBS {
				return false, ""
			}
		}
		return true, ""
	},
}
