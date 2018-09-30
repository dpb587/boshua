package analyzer

import (
	"github.com/dpb587/boshua/analysis"
	analyzerpkg "github.com/dpb587/boshua/analysis/analyzer"
	"github.com/dpb587/boshua/task"
)

const AnalyzerName analysis.AnalyzerName = "tilearchivefiles.v1"

type analyzer struct{}

var _ analysis.Analyzer = &analyzer{}

func (analyzer) Name() analysis.AnalyzerName {
	return AnalyzerName
}

func (analyzer) BuildTask(subject analysis.Subject) (*task.Task, error) {
	return analyzerpkg.NewSimpleTask(subject, AnalyzerName, false)
}

var Analyzer = &analyzer{}
