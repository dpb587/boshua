package checksum

import (
	"github.com/dpb587/bosh-compiled-releases/checksum/algorithm"
)

func New(a algorithm.Algorithm) *WritableChecksum {
	return &WritableChecksum{
		algorithm: a,
		hasher:    a.NewHash(),
	}
}

func NewExisting(a algorithm.Algorithm, d []byte) ImmutableChecksum {
	return ImmutableChecksum{
		algorithm: a,
		data:      d,
	}
}
