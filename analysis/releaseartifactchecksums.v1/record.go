package releaseartifactchecksums

import (
	"github.com/dpb587/bosh-compiled-releases/checksum"
)

type Record struct {
	Artifact string                       `json:"artifact" yaml:"artifact"`
	Path     string                       `json:"path" yaml:"path"`
	Result   []checksum.ImmutableChecksum `json:"result" yaml:"result"`
}
