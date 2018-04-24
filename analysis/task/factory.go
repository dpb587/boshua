package task

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"
)

func New(subject analysis.Subject, analyzer string) *Task {
	var found bool

	for _, expectedAnalyzer := range subject.SupportedAnalyzers() {
		if expectedAnalyzer == analyzer {
			found = true

			break
		}
	}

	if !found {
		panic(fmt.Errorf("invalid analyzer: %s", analyzer))
	}

	return &Task{
		subject:  subject,
		analyzer: analyzer,
	}
}
