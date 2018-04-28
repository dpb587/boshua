package releaseversion

import (
	"github.com/dpb587/boshua"
	"github.com/dpb587/boshua/checksum"
	"github.com/dpb587/boshua/util/metalinkutil"
	"github.com/dpb587/metalink"
)

type Artifact struct {
	Reference

	MetalinkFile   metalink.File
	MetalinkSource map[string]interface{}
}

var _ boshua.Artifact = &Artifact{}

func (s Artifact) ArtifactMetalink() metalink.Metalink {
	return metalink.Metalink{
		Files: []metalink.File{
			s.MetalinkFile,
		},
	}
}

func (s Artifact) ArtifactMetalinkStorage() map[string]interface{} {
	return s.MetalinkSource
}

func (s Artifact) MatchesChecksum(cs checksum.Checksum) bool {
	for _, hash := range s.MetalinkFile.Hashes {
		if metalinkutil.HashToChecksum(hash).String() == cs.String() {
			return true
		}
	}

	return false
}
