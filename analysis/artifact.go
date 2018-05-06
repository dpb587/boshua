package analysis

import (
	"github.com/dpb587/metalink"
)

type Artifact struct {
	Reference

	metalinkFile   metalink.File
	metalinkSource map[string]interface{}
}

func (s Artifact) ArtifactMetalinkFile() metalink.File {
	return s.metalinkFile
}

func (s Artifact) ArtifactMetalinkStorage() map[string]interface{} {
	return s.metalinkSource
}
