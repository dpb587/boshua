package analysisutil

import "github.com/dpb587/boshua/api/v2/models/analysis"

type AnalysisLoader func() (*analysis.GETAnalysisResponse, error)
