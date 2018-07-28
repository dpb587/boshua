package params

import (
	"errors"
	"net/http"

	"github.com/dpb587/boshua/analysis"
)

func AnalyzerNameFromQuery(r *http.Request) (analysis.AnalyzerName, error) {
	q := r.URL.Query()

	v, ok := q["analyzer"]
	if ok {
		if len(v) != 1 {
			return analysis.AnalyzerName(""), errors.New("analyzer: expected single value")
		}

		return analysis.AnalyzerName(v[0]), nil
	}

	return analysis.AnalyzerName(""), errors.New("analyzer: missing")
}
