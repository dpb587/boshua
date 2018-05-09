package releaseversion

import (
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/util/checksum"
	"github.com/dpb587/metalink"
)

type Artifact struct {
	Reference

	metalinkFile   metalink.File
	metalinkSource map[string]interface{}
}

var _ artifact.Artifact = &Artifact{}

func (s Artifact) ArtifactMetalinkFile() metalink.File {
	return s.metalinkFile
}

func (s Artifact) ArtifactMetalinkStorage() map[string]interface{} {
	return s.metalinkSource
}

func (s Artifact) MatchesChecksum(cs checksum.Checksum) bool {
	for _, hash := range s.metalinkFile.Hashes {
		if metalinkutil.HashToChecksum(hash).String() == cs.String() {
			return true
		}
	}

	return false
}
