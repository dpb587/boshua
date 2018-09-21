package releaseversion

import (
	"github.com/Masterminds/semver"
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/util/checksum"
	"github.com/dpb587/metalink"
)

type Artifact struct {
	Datastore string `json:"-"`

	Name          string        `json:"name"`
	Version       string        `json:"version"`
	SourceTarball metalink.File `json:"tarball"` // TODO rename to Tarball

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
		Name:      s.Name,
		Version:   s.Version,
		Checksums: metalinkutil.HashesToChecksums(s.SourceTarball.Hashes),
	}
}

func (s Artifact) GetLabels() []string {
	return s.Labels
}

func (s Artifact) GetDatastoreName() string {
	return s.Datastore
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
