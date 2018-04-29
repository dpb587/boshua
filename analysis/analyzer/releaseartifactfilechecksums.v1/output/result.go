package output

import (
	"github.com/dpb587/boshua/checksum"
)

type Result struct {
	Artifact string                       `json:"artifact" yaml:"artifact"`
	Path     string                       `json:"path" yaml:"path"`
	Result   []checksum.ImmutableChecksum `json:"result" yaml:"result"`
}
