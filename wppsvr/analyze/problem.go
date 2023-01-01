package analyze

import (
	"github.com/rothskeller/packet/wppsvr/english"
)

// A Problem object represents a particular problem that can be found in
// analysis of a practice message.
type Problem struct {
	// Code is the short string that identifies this problem.  It is a
	// single word in PascalCase.
	Code string
	// Variables is a map from variable names to functions that return the
	// values of those variables for a particular message.  These variables
	// can be interpolated into response messages.
	Variables variableMap
	// ifnot is a set of zero or more other problems that are exclusive to
	// this one.  If any of those other problems are found with a message,
	// we won't even attempt to detect this problem.
	ifnot []*Problem
	// after is a set of zero or more other problems that must be checked
	// before checking this one.
	after []*Problem
	// detect is a function that examines a practice message and returns
	// whether this problem exists with the message, and if so, which
	// response message should be used.
	detect func(*Analysis) (bool, string)
}

type variableMap map[string]func(*Analysis) string

// Problems is a map from problem code to Problem object.
var Problems = map[string]*Problem{}

// cacheOrderedProblems caches the result of the orderedProblems function.
var cacheOrderedProblems []*Problem

// orderedProblems returns an ordered list of all Problem objects.  The order is
// random except that it honors the 'ifnot' and 'after' fields of each problem.
func orderedProblems() (list []*Problem) {
	if cacheOrderedProblems != nil {
		return cacheOrderedProblems
	}
	var mp = make(map[*Problem]struct{}, len(Problems))
	for _, p := range Problems {
		mp[p] = struct{}{}
	}
	list = make([]*Problem, 0, len(Problems))
	var iter int
	for len(mp) != 0 {
		iter++
		if iter > 100 {
			panic("too many iterations — dependency loop in problems")
		}
		for p := range mp {
			var after = make(map[*Problem]struct{})
			for _, ap := range p.after {
				after[ap] = struct{}{}
			}
			for _, ip := range p.ifnot {
				after[ip] = struct{}{}
			}
			for _, op := range list {
				delete(after, op)
			}
			if len(after) != 0 {
				continue
			}
			list = append(list, p)
			delete(mp, p)
		}
	}
	cacheOrderedProblems = list
	return list
}

// Variables is the set of variables that can be interpolated into any response
// message for any problem.
var Variables = variableMap{
	"AMSGTYPE": func(a *Analysis) string {
		return a.xsc.Type.Article + " " + a.xsc.Type.Name
	},
	"FROMBBS": func(a *Analysis) string {
		return a.msg.FromBBS()
	},
	"FROMCALLSIGN": func(a *Analysis) string {
		return a.fromCallSign
	},
	"MSGDATE": func(a *Analysis) string {
		return a.msg.Date().Format("2006-01-02 at 15:04")
	},
	"MSGTYPE": func(a *Analysis) string {
		return a.xsc.Type.Name
	},
	"SESSIONBBSES": func(a *Analysis) string {
		return english.Conjoin(a.session.ToBBSes, "or")
	},
	"SESSIONDATE": func(a *Analysis) string {
		return a.session.End.Format("January 2")
	},
	"SESSIONNAME": func(a *Analysis) string {
		return a.session.Name
	},
	"TOBBS": func(a *Analysis) string {
		return a.toBBS
	},
	"TOCALLSIGN": func(a *Analysis) string {
		return a.session.CallSign
	},
}

func inList[T comparable](list []T, item T) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

var knownProblems map[string]map[string]struct{}
var knownVariables map[string]struct{}

// KnownProblems returns a map giving the names of known problems and the
// interpolated variables they support, and a map giving the names of global
// interpolated variables.  This is retrieved by main and passed into
// config.Validate, to avoid a circular import dependency between config and
// analyze.
func KnownProblems() (map[string]map[string]struct{}, map[string]struct{}) {
	if knownProblems == nil {
		knownProblems = make(map[string]map[string]struct{}, len(Problems))
		for _, prob := range Problems {
			if len(prob.Variables) != 0 {
				vars := make(map[string]struct{}, len(prob.Variables))
				for varname, fn := range prob.Variables {
					if fn != nil {
						vars[varname] = struct{}{}
					}
				}
				knownProblems[prob.Code] = vars
			} else {
				knownProblems[prob.Code] = nil
			}
		}
		knownVariables = make(map[string]struct{}, len(Variables))
		for varname, fn := range Variables {
			if fn != nil {
				knownVariables[varname] = struct{}{}
			}
		}
	}
	return knownProblems, knownVariables
}
