package analysis

import "github.com/dpb587/boshua"

type Subject interface {
	boshua.Subject

	SupportedAnalyzers() []string
}
