package analyze

func init() {
	Problems[ProbFromBBSDown.Code] = ProbFromBBSDown
	Problems[ProbToBBS.Code] = ProbToBBS
	Problems[ProbToBBSDown.Code] = ProbToBBSDown
}

// ProbFromBBSDown is raised when a message is sent from a BBS that is simulated
// down.
var ProbFromBBSDown = &Problem{
	Code: "FromBBSDown",
	detect: func(a *Analysis) bool {
		return inList(a.session.DownBBSes, a.msg.FromBBS())
	},
}

// ProbToBBSDown is raised when a message is sent to a BBS that is simulated
// down.
var ProbToBBSDown = &Problem{
	Code: "ToBBSDown",
	detect: func(a *Analysis) bool {
		return inList(a.session.DownBBSes, a.toBBS)
	},
}

// ProbToBBS is raised when a message is sent to a BBS that is not a correct BBS
// for the session (and is not down).
var ProbToBBS = &Problem{
	Code:  "ToBBS",
	ifnot: []*Problem{ProbToBBSDown},
	detect: func(a *Analysis) bool {
		return !inList(a.session.ToBBSes, a.toBBS)
	},
}
