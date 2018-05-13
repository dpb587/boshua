package stemcellversion

import (
	"strings"

	"github.com/dpb587/boshua/analysis"
)

var _ analysis.Subject = &Artifact{}

func (s Artifact) SupportedAnalyzers() []string {
	analyzers := []string{
		"stemcellmanifest.v1",
	}

	if !strings.HasPrefix(s.reference.OS, "windows") {
		analyzers = append(analyzers,
			"stemcellimagefiles.v1",
			"stemcellpackages.v1",
		)
	}

	return analyzers
}
