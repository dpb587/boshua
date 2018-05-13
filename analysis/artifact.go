package analysis

import (
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/metalink"
)

type Artifact struct {
	artifact artifact.Artifact
	analyzer AnalyzerName
	subject  Subject
}

func (a Artifact) Analyzer() AnalyzerName {
	return a.analyzer
}

func (a Artifact) MetalinkFile() metalink.File {
	return a.artifact.MetalinkFile()
}

func (a Artifact) Subject() Subject {
	return a.subject
}
