package main

import (
	"fmt"
	"io"

	"github.com/dpb587/boshua/analysis"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/config"
	"github.com/dpb587/boshua/config/presets/defaults"
	"github.com/dpb587/boshua/metalink/analysisprocessor"
	releasemanifestsV1 "github.com/dpb587/boshua/releaseversion/analyzers/releasemanifests.v1/result"
	releasedatastore "github.com/dpb587/boshua/releaseversion/datastore"
)

func main() {
	cfg, err := defaults.NewConfig()
	if err != nil {
		panic(err)
	}

	releases, err := cfg.GetReleaseIndex(config.DefaultName)
	if err != nil {
		panic(err)
	}

	release, err := releasedatastore.GetArtifact(releases, releasedatastore.FilterParamsFromSlug("openvpn/5.1.0"))
	if err != nil {
		panic(err)
	}

	analyses, err := cfg.GetReleaseAnalysisIndex(release.GetDatastoreName())
	if err != nil {
		panic(err)
	}

	analysis, err := analysisdatastore.GetAnalysisArtifact(analyses, analysis.Reference{
		Subject:  release,
		Analyzer: "releasemanifests.v1",
	})
	if err != nil {
		panic(err)
	}

	err = analysisprocessor.Process(analysis, func(reader io.Reader) error {
		return releasemanifestsV1.NewProcessor(reader, func(r releasemanifestsV1.Record) error {
			if r.Path == "release.MF" {
				fmt.Printf("%s\n", r.Raw)
			}

			return nil
		})
	})
	if err != nil {
		panic(err)
	}
}
