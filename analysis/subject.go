package analysis

import "github.com/dpb587/boshua"

type Subject interface {
	boshua.Artifact

	SupportedAnalyzers() []string
}
