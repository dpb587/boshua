package releaseversion

import (
	"fmt"

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

func (s Artifact) ArtifactReference() boshua.Reference {
	return s.Reference.ArtifactReference()
}

func (s Artifact) ArtifactStorageDir() string {
	ref := s.ArtifactReference()

	return fmt.Sprintf(
		"%s/%s/%s/%s",
		ref.Context,
		ref.ID[0:2],
		ref.ID[2:4],
		ref.ID[4:],
	)
}

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
