package compilation

import (
	"github.com/dpb587/metalink"
)

type Artifact struct {
	reference Reference
	Tarball   metalink.File `json:"tarball"`
}

func (s Artifact) MetalinkFile() metalink.File {
	return s.Tarball
}

func (s Artifact) Reference() interface{} {
	return s.reference
}
