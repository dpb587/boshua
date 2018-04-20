package compiledreleaseversion

import (
	"time"

	"github.com/dpb587/boshua/checksum"
)

type Subject struct {
	Reference

	TarballURL       string
	TarballSize      *uint64
	TarballPublished *time.Time
	TarballChecksums checksum.ImmutableChecksums
}
