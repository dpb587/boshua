package analysis

import "github.com/dpb587/boshua/artifact"

type Subject interface {
	artifact.Artifact

	SupportedAnalyzers() []string
}
