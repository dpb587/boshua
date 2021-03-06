package stemcellversion

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/util/checksum"
	"github.com/dpb587/metalink"
)

type Artifact struct {
	Datastore string `json:"-"`

	OS         string        `json:"os"`
	Version    string        `json:"version"`
	IaaS       string        `json:"iaas"`
	Hypervisor string        `json:"hypervisor"`
	DiskFormat string        `json:"diskFormat"`
	Flavor     string        `json:"flavor"`
	Tarball    metalink.File `json:"tarball"`

	Labels []string `json:"labels"`

	semver       *semver.Version
	semverParsed bool
}

var _ artifact.Artifact = &Artifact{}

func (r Artifact) FullName() string {
	// TODO rename to Name()
	var prefix string

	return fmt.Sprintf("%sbosh-%s-%s-%s-go_agent", prefix, r.IaaS, r.Hypervisor, r.OS)
}

func (s Artifact) PreferredChecksum() checksum.ImmutableChecksum {
	// TODO should not panic; should be nillable
	return metalinkutil.HashToChecksum(metalinkutil.PreferredHash(s.Tarball.Hashes))
}

func (s Artifact) Reference() interface{} {
	return Reference{
		IaaS:       s.IaaS,
		Hypervisor: s.Hypervisor,
		OS:         s.OS,
		Version:    s.Version,
		Flavor:     s.Flavor,
	}
}

func (s Artifact) MetalinkFile() metalink.File {
	return s.Tarball
}

func (s Artifact) GetLabels() []string {
	return s.Labels
}

func (s Artifact) GetDatastoreName() string {
	return s.Datastore
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
