package compiledreleaseversions

import (
	"time"

	"github.com/dpb587/boshua/checksum"
)

type CompiledReleaseVersion struct {
	CompiledReleaseVersionRef

	TarballURL       string
	TarballSize      *uint64
	TarballPublished *time.Time
	TarballChecksums checksum.ImmutableChecksums
}
