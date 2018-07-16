package stemcellversion

import (
	"fmt"

	"github.com/dpb587/boshua/artifact"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/util/checksum"
	"github.com/dpb587/metalink"
)

type Artifact struct {
	OS         string        `json:"os"`
	Version    string        `json:"version"`
	IaaS       string        `json:"iaas"`
	Hypervisor string        `json:"hypervisor"`
	DiskFormat string        `json:"diskFormat"`
	Light      bool          `json:"light"`
	Tarball    metalink.File `json:"tarball"`
}

var _ artifact.Artifact = &Artifact{}

func (r Artifact) FullName() string {
	// TODO rename to Name()
	// TODO breaks with light prefix; should match name from `bosh stemcells`
	var prefix string

	if r.Light {
		prefix = "light-"
	}

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
		Light:      s.Light,
	}
}

func (s Artifact) MetalinkFile() metalink.File {
	return s.Tarball
}
