package compilation

import (
	"github.com/dpb587/metalink"
)

func New(ref Reference, meta4File metalink.File) Artifact {
	return Artifact{
		reference: ref,
		Tarball:   meta4File,
	}
}
