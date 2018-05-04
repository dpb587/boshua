package stemcellversion

import (
	"github.com/dpb587/boshua/stemcellversion"
)

type Reference struct {
	IaaS       string `json:"iaas"`
	Hypervisor string `json:"hypervisor"`
	OS         string `json:"os"`
	Version    string `json:"version"`
	Light      bool   `json:"light"`
}

func FromReference(ref stemcellversion.Reference) Reference {
	return Reference{
		IaaS:       ref.IaaS,
		Hypervisor: ref.Hypervisor,
		OS:         ref.OS,
		Version:    ref.Version,
		Light:      ref.Light,
	}
}
