package analysis

import (
	"fmt"

	"github.com/dpb587/metalink"
)

type Artifact struct {
	Reference

	MetalinkFile   metalink.File
	MetalinkSource map[string]interface{}
}

func (s Artifact) ArtifactMetalink() metalink.Metalink {
	return metalink.Metalink{
		Files: []metalink.File{
			s.MetalinkFile,
		},
	}
}

func (s Artifact) ArtifactMetalinkStorage() map[string]interface{} {
	return map[string]interface{}{
		"uri": fmt.Sprintf("git@github.com:dpb587/bosh-compiled-releases-index.git//%s", s.ArtifactStorageDir()),
		"options": map[string]string{
			"private_key": "((index_private_key))",
		},
	}
}
