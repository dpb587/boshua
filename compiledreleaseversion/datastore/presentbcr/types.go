package presentbcr

import (
	"github.com/dpb587/boshua/checksum"
)

type Record struct {
	Release RecordRelease `json:"release"`
	OS      RecordOS      `json:"os"`
}

type RecordRelease struct {
	Name      string                      `json:"name"`
	Version   string                      `json:"version"`
	Checksums checksum.ImmutableChecksums `json:"checksums"`
}

type RecordOS struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
