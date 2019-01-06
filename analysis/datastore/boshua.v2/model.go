package boshuaV2

import (
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/metalink"
)

type filterResponse struct {
	Releases  []filterReleasesResponse  `json:"releases"`
	Stemcells []filterStemcellsResponse `json:"stemcells"`
}

type filterReleasesResponse struct {
	Compilations []filterReleasesCompilationsResponse `json:"compilations"`
	Analysis     *filterAnalysisResponse              `json:"analysis"`
}

type filterReleasesCompilationsResponse struct {
	Analysis *filterAnalysisResponse `json:"analysis"`
}

type filterStemcellsResponse struct {
	Analysis *filterAnalysisResponse `json:"analysis"`
}

type filterAnalysisResponse struct {
	Results []filterAnalysisResultResponse `json:"results"`
}

type filterAnalysisResultResponse struct {
	Analyzer analysis.AnalyzerName `json:"analyzer"`
	Artifact metalink.File         `json:"artifact"`
}

func (r filterResponse) GetAnalysis() []filterAnalysisResultResponse {
	// TODO this is bad; it's assuming a single result for these lookups; could technically batch results
	if len(r.Releases) > 0 {
		if r.Releases[0].Analysis != nil {
			return r.Releases[0].Analysis.Results
		} else if len(r.Releases[0].Compilations) > 0 {
			return r.Releases[0].Compilations[0].Analysis.Results
		}
	} else if len(r.Stemcells) > 0 {
		if r.Stemcells[0].Analysis != nil {
			return r.Stemcells[0].Analysis.Results
		}
	}

	panic("unexpected results") // !panic
}
