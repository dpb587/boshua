package releaseversion

import (
	"github.com/Masterminds/semver"
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/util/checksum"
	"github.com/dpb587/metalink"
)

type Artifact struct {
	Name          string        `json:"name"`
	Version       string        `json:"version"`
	SourceTarball metalink.File `json:"source_tarball"`

	Labels []string `json:"labels"`

	semver       *semver.Version
	semverParsed bool
}

var _ artifact.Artifact = &Artifact{}

func (s Artifact) MetalinkFile() metalink.File {
	return s.SourceTarball
}

func (s Artifact) Reference() interface{} {
	return Reference{
		Name:    s.Name,
		Version: s.Version,
	}
}

func (s Artifact) PreferredChecksum() checksum.ImmutableChecksum {
	// TODO should not panic; should be nillable
	return metalinkutil.HashToChecksum(metalinkutil.PreferredHash(s.SourceTarball.Hashes))
}

func (s Artifact) MatchesChecksum(cs checksum.Checksum) bool {
	for _, hash := range s.SourceTarball.Hashes {
		if metalinkutil.HashToChecksum(hash).String() == cs.String() {
			return true
		}
	}

	return false
}

func (s Artifact) Semver() *semver.Version {
	if s.semverParsed {
		return s.semver
	}

	semver, err := semver.NewVersion(s.Version)
	if err == nil {
		s.semver = semver
	}

	s.semverParsed = true

	return s.semver
}
