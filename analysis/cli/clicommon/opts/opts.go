package opts

import (
	"github.com/dpb587/boshua/analysis"
)

type Opts struct {
	Analyzer analysis.AnalyzerName `long:"analyzer" description:"The analyzer name"`
}
