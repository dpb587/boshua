package result

import (
	"github.com/dpb587/boshua/analysis/analyzer/filescommon.v1/result"
)

type Record struct {
	Artifact string        `json:"artifact" yaml:"artifact"`
	Result   result.Record `json:"result" yaml:"result"`
}
