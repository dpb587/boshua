package output

import (
	"github.com/dpb587/boshua/analysis/analyzer/filescommon.v1/output"
)

type Result struct {
	Artifact string        `json:"artifact" yaml:"artifact"`
	Result   output.Result `json:"result" yaml:"result"`
}
