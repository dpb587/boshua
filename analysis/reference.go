package analysis

import (
	"github.com/dpb587/boshua/artifact"
)

type Reference struct {
	Subject  artifact.Artifact
	Analyzer AnalyzerName
}
