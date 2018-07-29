package analysis

import (
	"github.com/dpb587/boshua/task"
)

type Analyzer interface {
	Name() AnalyzerName
	BuildTask(subject Subject) (*task.Task, error)
}

type AnalyzerName string

type AnalysisGenerator interface {
	Analyze(Writer) error
}
