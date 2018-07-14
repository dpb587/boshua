package v2

import (
	"net/http"

	"github.com/dpb587/boshua/analysis"
	"github.com/dpb587/boshua/server/httputil"
)

func ApplyAnalysisAnalyzerToQuery(r *http.Request, analyzer string) {
	q := r.URL.Query()

	q.Add("analysis.analyzer", analyzer)

	r.URL.RawQuery = q.Encode()
}

func AnalysisAnalyzerFromParam(r *http.Request) (analysis.AnalyzerName, error) {
	v, err := httputil.SimpleQueryLookup(r, "analysis.analyzer")

	return analysis.AnalyzerName(v), err
}
