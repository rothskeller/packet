package analyze

import (
	"fmt"
	"strings"

	"github.com/rothskeller/packet/wppsvr/config"
)

const multipleProblemsSubject = "Issues with packet practice message"

func (a *Analysis) reportProblem(problemCode string, references reference, report string, args ...any) {
	if _, ok := a.problems[problemCode]; ok {
		panic(problemCode + " reported twice")
	}
	a.problems[problemCode] = struct{}{}
	actions := config.Get().ProblemActionFlags[problemCode]
	if actions&config.ActionDontCount != 0 {
		a.invalid = true
	}
	if actions&config.ActionRespond == 0 {
		return // no need to generate response
	}
	if a.reportSubject == "" {
		if a.reportSubject = ProblemLabels[problemCode]; a.reportSubject == "" {
			panic("missing problem label for " + problemCode)
		}
		a.reportSubject = strings.ToUpper(a.reportSubject[:1]) + a.reportSubject[1:]
	} else {
		a.reportSubject = multipleProblemsSubject
	}
	a.references |= references
	if len(args) == 0 {
		a.reportText.WriteString(report)
	} else {
		fmt.Fprintf(&a.reportText, report, args...)
	}
	a.reportText.WriteString("\n\n")
}

// ProblemLabels is a map from problem code to problem label (a short string
// describing the problem).
var ProblemLabels = map[string]string{}

func inList[T comparable](list []T, item T) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}
