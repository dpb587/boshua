package releaseartifactchecksums

import (
	"github.com/dpb587/bosh-compiled-releases/checksum"
)

type Record struct {
	Artifact string                       `json:"artifact"`
	Path     string                       `json:"path"`
	Result   []checksum.ImmutableChecksum `json:"result"`
}
