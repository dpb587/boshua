package stemcellversion

import (
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/metalink"
)

type Artifact struct {
	reference    Reference
	metalinkFile metalink.File
}

var _ artifact.Artifact = &Artifact{}

func (s Artifact) Reference() interface{} {
	return s.reference
}

func (s Artifact) MetalinkFile() metalink.File {
	return s.metalinkFile
}
