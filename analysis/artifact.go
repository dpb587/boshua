package analysis

import (
	"time"

	"github.com/dpb587/boshua/checksum"
)

type Artifact struct {
	TarballURL       string
	TarballSize      *uint64
	TarballPublished *time.Time
	TarballChecksums checksum.ImmutableChecksums
}
