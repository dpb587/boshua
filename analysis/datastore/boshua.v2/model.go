package boshuaV2

import (
	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/metalink"
)

type filterResponse struct {
	Release  filterReleaseResponse  `json:"release"`
	Stemcell filterStemcellResponse `json:"stemcell"`
}

type filterReleaseResponse struct {
	Compilations []filterReleaseCompilationResponse `json:"compilations"`
	Analysis     *filterAnalysisResponse            `json:"analysis"`
}

type filterReleaseCompilationResponse struct {
	Analysis *filterAnalysisResponse `json:"analysis"`
}

type filterStemcellResponse struct {
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
	if r.Release.Analysis != nil {
		return r.Release.Analysis.Results
	} else if len(r.Release.Compilations) > 0 {
		return r.Release.Compilations[0].Analysis.Results
	} else if r.Stemcell.Analysis != nil {
		return r.Stemcell.Analysis.Results
	}

	panic("unexpected results") // !panic
}
