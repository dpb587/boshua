package releaseversion

import (
	"github.com/dpb587/boshua/checksum"
)

type Reference struct {
	Name     string
	Version  string
	Checksum checksum.ImmutableChecksum
}
