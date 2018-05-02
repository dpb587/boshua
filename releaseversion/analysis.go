package releaseversion

import "github.com/dpb587/boshua/analysis"

var _ analysis.Subject = &Artifact{}

func (Reference) SupportedAnalyzers() []string {
	return []string{
		"releaseartifactfiles.v1",
		"releasemanifests.v1",
	}
}
