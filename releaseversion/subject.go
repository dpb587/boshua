package releaseversion

import (
	"github.com/dpb587/boshua/checksum"
)

type Subject struct {
	Reference

	Checksums      checksum.ImmutableChecksums
	MetalinkSource map[string]interface{}
}
