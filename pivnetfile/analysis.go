package pivnetfile

import "github.com/dpb587/boshua/analysis"

var _ analysis.Subject = &Artifact{}

func (Artifact) SupportedAnalyzers() []analysis.AnalyzerName {
	return []analysis.AnalyzerName{
		"tilereleasemanifests.v1",
	}
}
