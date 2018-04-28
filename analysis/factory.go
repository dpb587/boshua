package analysis

import (
	"github.com/dpb587/boshua"
	"github.com/dpb587/metalink"
)

func New(artifact boshua.ArtifactReference, analyzer string, meta4File metalink.File, meta4Source map[string]interface{}) Artifact {
	return Artifact{
		Reference: Reference{
			Artifact: artifact,
			Analyzer: analyzer,
		},
		MetalinkFile:   meta4File,
		MetalinkSource: meta4Source,
	}
}
