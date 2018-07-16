package stemcellversion

import (
	"github.com/dpb587/boshua/artifact"
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
