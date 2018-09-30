package result

import (
	"github.com/dpb587/boshua/analysis/analyzer/filescommon.v1/result"
)

type Record struct {
	Parents []string      `json:"parents" yaml:"parents"`
	Result  result.Record `json:"result" yaml:"result"`
}
