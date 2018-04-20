package releaseversions

import (
	"github.com/dpb587/boshua/checksum"
)

type ReleaseVersion struct {
	ReleaseVersionRef

	Checksums      checksum.ImmutableChecksums
	MetalinkSource map[string]interface{}
}
