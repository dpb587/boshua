package urlutil

import "net/http"

func ApplyAnalysisAnalyzerToQuery(r *http.Request, analyzer string) {
	q := r.URL.Query()

	q.Add("analysis.analyzer", analyzer)

	r.URL.RawQuery = q.Encode()
}

func AnalysisAnalyzerFromParam(r *http.Request) (string, error) {
	return simpleQueryLookup(r, "analysis.analyzer")
}
