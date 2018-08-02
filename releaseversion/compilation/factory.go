package compilation

import (
	"github.com/dpb587/metalink"
)

func New(ref Reference, meta4File metalink.File) Artifact {
	return Artifact{
		Release: ref.ReleaseVersion,
		OS:      ref.OSVersion,
		Tarball: meta4File,
	}
}
