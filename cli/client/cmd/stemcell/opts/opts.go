package opts

import (
	"github.com/dpb587/boshua/cli/client/args"
	"github.com/dpb587/boshua/stemcellversion"
)

type Opts struct {
	Stemcell args.Stemcell `long:"stemcell" description:"The stemcell name and version"`
}

func (o Opts) Reference() stemcellversion.Reference {
	return stemcellversion.Reference{
		IaaS:       o.Stemcell.IaaS,
		Hypervisor: o.Stemcell.Hypervisor,
		OS:         o.Stemcell.OS,
		Version:    o.Stemcell.Version,
		// Light:      o.Stemcell.Light,
	}
}
