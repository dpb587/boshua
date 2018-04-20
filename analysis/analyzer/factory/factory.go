package factory

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"

	releaseartifactchecksumsv1 "github.com/dpb587/boshua/analysis/analyzer/releaseartifactchecksums.v1"
	releaseartifactfilestatv1 "github.com/dpb587/boshua/analysis/analyzer/releaseartifactfilestat.v1"
	releasemanifestsv1 "github.com/dpb587/boshua/analysis/analyzer/releasemanifests.v1"
)

type Factory struct{}

func (Factory) Create(analyzer string, path string) (analysis.Analyzer, error) {
	if analyzer == "releaseartifactchecksums.v1" {
		return releaseartifactchecksumsv1.New(path), nil
	} else if analyzer == "releaseartifactfilestat.v1" {
		return releaseartifactfilestatv1.New(path), nil
	} else if analyzer == "releasemanifests.v1" {
		return releasemanifestsv1.New(path), nil
	}

	return nil, fmt.Errorf("invalid analyzer: %s", analyzer)
}
