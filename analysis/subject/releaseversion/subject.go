package releaseversion

import "github.com/dpb587/boshua/releaseversion/datastore"

type Subject struct {
	input releaseversions.ReleaseVersion
}

func (s Subject) SupportedAnalyzers() []string {
	return []string{
		"releaseartifactchecksums.v1",
		"releaseartifactfilestat.v1",
		"releasemanifests.v1",
	}
}

func (s Subject) Input() map[string]interface{} {
	return s.input.MetalinkSource
}
