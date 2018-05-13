package releaseversion

import (
	"github.com/dpb587/boshua/util/checksum"
)

type Reference struct {
	Name      string                      `json:"name"`
	Version   string                      `json:"version"`
	Checksums checksum.ImmutableChecksums `json:"checksums"`
	URLs      []string                    `json:"urls"`
}
