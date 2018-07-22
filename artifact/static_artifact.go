package artifact

import (
	"github.com/dpb587/metalink"
)

type StaticArtifact struct {
	StaticMetalinkFile metalink.File
}

func (StaticArtifact) Reference() interface{} {
	return nil
}

func (a StaticArtifact) MetalinkFile() metalink.File {
	return a.StaticMetalinkFile
}
