package releaseversion

import "github.com/dpb587/boshua/analysis"

var _ analysis.Subject = &Artifact{}

func (Artifact) SupportedAnalyzers() []string {
	return []string{
		"releaseartifactchecksums.v1",
		"releaseartifactfilestat.v1",
		"releasemanifests.v1",
	}
}
