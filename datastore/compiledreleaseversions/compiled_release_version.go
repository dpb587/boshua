package compiledreleaseversions

import (
	"time"

	"github.com/dpb587/boshua/datastore/releaseversions"
)

type CompiledReleaseVersion struct {
	CompiledReleaseVersionRef

	TarballURL       string
	TarballSize      *uint64
	TarballPublished *time.Time
	TarballChecksums releaseversions.Checksums
}
