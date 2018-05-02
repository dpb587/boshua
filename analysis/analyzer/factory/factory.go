package factory

import (
	"fmt"

	"github.com/dpb587/boshua/analysis"

	releaseartifactfilesv1 "github.com/dpb587/boshua/analysis/analyzer/releaseartifactfiles.v1"
	releasemanifestsv1 "github.com/dpb587/boshua/analysis/analyzer/releasemanifests.v1"
	stemcellimagefilesv1 "github.com/dpb587/boshua/analysis/analyzer/stemcellimagefiles.v1"
	stemcellmanifestv1 "github.com/dpb587/boshua/analysis/analyzer/stemcellmanifest.v1"
	stemcellpackagesv1 "github.com/dpb587/boshua/analysis/analyzer/stemcellpackages.v1"
)

type Factory struct{}

func (Factory) Create(analyzer string, path string) (analysis.Analyzer, error) {
	if analyzer == "releaseartifactfiles.v1" {
		return releaseartifactfilesv1.New(path), nil
	} else if analyzer == "releasemanifests.v1" {
		return releasemanifestsv1.New(path), nil
	} else if analyzer == "stemcellimagefiles.v1" {
		return stemcellimagefilesv1.New(path), nil
	} else if analyzer == "stemcellmanifest.v1" {
		return stemcellmanifestv1.New(path), nil
	} else if analyzer == "stemcellpackages.v1" {
		return stemcellpackagesv1.New(path), nil
	}

	return nil, fmt.Errorf("invalid analyzer: %s", analyzer)
}
