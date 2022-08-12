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
	Label: "message from incorrect BBS (simulated outage)",
	ifnot: []*Problem{ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		for _, down := range a.session.DownBBSes {
			if down == a.msg.FromBBS() {
				return true, ""
			}
		}
		return false, ""
	},
	references: refWeeklyPractice,
}

// ProbToBBSDown is raised when a message is sent to a BBS that is simulated
// down.
var ProbToBBSDown = &Problem{
	Code:  "ToBBSDown",
	Label: "message to incorrect BBS (simulated outage)",
	ifnot: []*Problem{ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		for _, down := range a.session.DownBBSes {
			if down == a.toBBS {
				return true, ""
			}
		}
		return false, ""
	},
	references: refWeeklyPractice,
}

// ProbToBBS is raised when a message is sent to a BBS that is not a correct BBS
// for the session.
var ProbToBBS = &Problem{
	Code:  "ToBBS",
	Label: "message to incorrect BBS",
	ifnot: []*Problem{ProbToBBSDown, ProbBounceMessage, ProbDeliveryReceipt, ProbReadReceipt},
	detect: func(a *Analysis) (bool, string) {
		for _, to := range a.session.ToBBSes {
			if to == a.toBBS {
				return false, ""
			}
		}
		return true, ""
	},
	references: refWeeklyPractice,
}
