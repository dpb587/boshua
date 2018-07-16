package opts

import (
	"github.com/dpb587/boshua/stemcellversion/datastore"
)

type Opts struct {
	// Stemcell args.Stemcell `long:"stemcell" description:"The stemcell name and version"` // TODO resurrect Name parsing?
	OS         string `long:"stemcell-os" description:"The stemcell OS"`
	Version    string `long:"stemcell-version" description:"The stemcell version"`
	IaaS       string `long:"stemcell-iaas" description:"The stemcell IaaS"`
	Hypervisor string `long:"stemcell-hypervisor" description:"The stemcell hypervisor"`
	Light      bool   `long:"stemcell-light" description:"The stemcell as a light version"`
}

func (o Opts) FilterParams() *datastore.FilterParams {
	return &datastore.FilterParams{
		IaaS:       o.IaaS,
		Hypervisor: o.Hypervisor,
		OS:         o.OS,
		Version:    o.Version,
		Light:      o.Light,
	}
}
