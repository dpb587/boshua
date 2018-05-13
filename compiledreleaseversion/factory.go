package compiledreleaseversion

import (
	"github.com/dpb587/metalink"
)

func New(ref Reference, meta4File metalink.File) Artifact {
	return Artifact{
		reference:    ref,
		metalinkFile: meta4File,
	}
}
