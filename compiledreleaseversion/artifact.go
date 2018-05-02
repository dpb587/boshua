package compiledreleaseversion

import (
	"github.com/dpb587/metalink"
)

type Artifact struct {
	Reference

	metalinkFile   metalink.File
	metalinkSource map[string]interface{}
}

func (s Artifact) ArtifactMetalink() metalink.Metalink {
	return metalink.Metalink{
		Files: []metalink.File{
			s.metalinkFile,
		},
	}
}

func (s Artifact) ArtifactMetalinkStorage() map[string]interface{} {
	return s.metalinkSource
}
