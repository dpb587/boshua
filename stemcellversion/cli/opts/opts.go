package opts

import (
	"github.com/dpb587/boshua/cli/args"
	"github.com/dpb587/boshua/stemcellversion/datastore"
)

type Opts struct {
	Stemcell args.Stemcell `long:"stemcell" description:"The stemcell name and version"`
}

func (o Opts) FilterParams() *datastore.FilterParams {
	return &datastore.FilterParams{
		IaaS:       o.Stemcell.IaaS,
		Hypervisor: o.Stemcell.Hypervisor,
		OS:         o.Stemcell.OS,
		Version:    o.Stemcell.Version,
		// Light:      o.Stemcell.Light,
	}
}
