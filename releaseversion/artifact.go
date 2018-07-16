package releaseversion

import (
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/util/checksum"
	"github.com/dpb587/metalink"
)

type Artifact struct {
	Name          string        `json:"name"`
	Version       string        `json:"version"`
	SourceTarball metalink.File `json:"source_tarball"`
}

var _ artifact.Artifact = &Artifact{}

func (s Artifact) MetalinkFile() metalink.File {
	return s.SourceTarball
}

func (s Artifact) Reference() interface{} {
	ref := Reference{
		Name:    s.Name,
		Version: s.Version,
	}

	return ref
}

func (s Artifact) MatchesChecksum(cs checksum.Checksum) bool {
	for _, hash := range s.SourceTarball.Hashes {
		if metalinkutil.HashToChecksum(hash).String() == cs.String() {
			return true
		}
	}

	return false
}
