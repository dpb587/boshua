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

	if !strings.HasPrefix(s.OS, "windows") {
		analyzers = append(analyzers,
			"stemcellimagefilechecksums.v1",
			"stemcellimagefilestat.v1",
			"stemcellpackages.v1",
		)
	}

	return analyzers
}
