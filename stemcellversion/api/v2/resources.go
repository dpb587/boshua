package v2

import (
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/metalink"
)

type InfoResponse struct {
	Data InfoResponseData `json:"data"`
}

type InfoResponseData struct {
	Reference Reference     `json:"reference"`
	Artifact  metalink.File `json:"file"`
}

type GETIndexResponse struct {
	Data []Reference `json:"data"`
}

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
