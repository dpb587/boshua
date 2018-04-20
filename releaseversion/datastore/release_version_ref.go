package releaseversions

import (
	"github.com/dpb587/boshua/checksum"
)

type ReleaseVersionRef struct {
	Name     string
	Version  string
	Checksum checksum.ImmutableChecksum
}
