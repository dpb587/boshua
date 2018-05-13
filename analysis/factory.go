package analysis

import (
	"github.com/dpb587/boshua/artifact"
)

func New(artifact artifact.Artifact, analyzer AnalyzerName, subject Subject) Artifact {
	return Artifact{
		artifact: artifact,
		analyzer: analyzer,
		subject:  subject,
	}
}
