package factory

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"

	releaseartifactfilechecksumsv1 "github.com/dpb587/boshua/analysis/analyzer/releaseartifactfilechecksums.v1"
	releaseartifactfilestatv1 "github.com/dpb587/boshua/analysis/analyzer/releaseartifactfilestat.v1"
	releasemanifestsv1 "github.com/dpb587/boshua/analysis/analyzer/releasemanifests.v1"
	stemcellimagefilechecksumsv1 "github.com/dpb587/boshua/analysis/analyzer/stemcellimagefilechecksums.v1"
	stemcellimagefilestatv1 "github.com/dpb587/boshua/analysis/analyzer/stemcellimagefilestat.v1"
	stemcellmanifestv1 "github.com/dpb587/boshua/analysis/analyzer/stemcellmanifest.v1"
	stemcellpackagesv1 "github.com/dpb587/boshua/analysis/analyzer/stemcellpackages.v1"
)

type Factory struct{}

func (Factory) Create(analyzer string, path string) (analysis.Analyzer, error) {
	if analyzer == "releaseartifactfilechecksums.v1" {
		return releaseartifactfilechecksumsv1.New(path), nil
	} else if analyzer == "releaseartifactfilestat.v1" {
		return releaseartifactfilestatv1.New(path), nil
	} else if analyzer == "releasemanifests.v1" {
		return releasemanifestsv1.New(path), nil
	} else if analyzer == "stemcellimagefilechecksums.v1" {
		return stemcellimagefilechecksumsv1.New(path), nil
	} else if analyzer == "stemcellimagefilestat.v1" {
		return stemcellimagefilestatv1.New(path), nil
	} else if analyzer == "stemcellmanifest.v1" {
		return stemcellmanifestv1.New(path), nil
	} else if analyzer == "stemcellpackages.v1" {
		return stemcellpackagesv1.New(path), nil
	}

	return nil, fmt.Errorf("invalid analyzer: %s", analyzer)
}
