package artifact

import (
	"github.com/dpb587/metalink"
)

type StaticArtifact struct {
	StaticMetalinkFile metalink.File
}

var _ Artifact = &StaticArtifact{}

func (StaticArtifact) Reference() interface{} {
	return nil
}

func (a StaticArtifact) MetalinkFile() metalink.File {
	return a.StaticMetalinkFile
}

func (StaticArtifact) GetLabels() []string {
	return nil
}

func (StaticArtifact) GetDatastoreName() string {
	return "unknown"
}
