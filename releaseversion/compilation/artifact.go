package compilation

import (
	"github.com/dpb587/metalink"
)

type Artifact struct {
	reference    Reference
	metalinkFile metalink.File
}

func (s Artifact) MetalinkFile() metalink.File {
	return s.metalinkFile
}

func (s Artifact) Reference() interface{} {
	return s.reference
}
