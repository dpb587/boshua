package opts

import (
	"github.com/dpb587/boshua/cli/args"
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/pkg/errors"
)

type Opts struct {
	AppOpts *cmdopts.Opts `no-flag:"true"`

	Stemcell   *args.Stemcell `long:"stemcell" description:"The stemcell name and version"`
	OS         string         `long:"stemcell-os" description:"The stemcell OS"`
	Version    string         `long:"stemcell-version" description:"The stemcell version"`
	IaaS       string         `long:"stemcell-iaas" description:"The stemcell IaaS"`
	Hypervisor string         `long:"stemcell-hypervisor" description:"The stemcell hypervisor"`
	// Light      bool           `long:"stemcell-light" description:"The stemcell as a light version"` // TODO add filtering support; heavy?
}

func (o *Opts) Artifact() (stemcellversion.Artifact, error) {
	index, err := o.AppOpts.GetStemcellIndex("default")
	if err != nil {
		return stemcellversion.Artifact{}, errors.Wrap(err, "loading stemcell index")
	}

	res, err := index.Filter(o.FilterParams())
	if err != nil {
		return stemcellversion.Artifact{}, errors.Wrap(err, "finding stemcell")
	} else if err = datastore.RequireSingleResult(res); err != nil {
		return stemcellversion.Artifact{}, errors.Wrap(err, "finding stemcell")
	}

	return res[0], err
}

func (o Opts) FilterParams() *datastore.FilterParams {
	f := &datastore.FilterParams{
		// Light: o.Light,
	}

	if o.Stemcell != nil {
		f.OSExpected = true
		f.OS = o.Stemcell.OS

		f.VersionExpected = true
		f.Version = o.Stemcell.Version

		f.IaaSExpected = true
		f.IaaS = o.Stemcell.IaaS

		f.HypervisorExpected = true
		f.Hypervisor = o.Stemcell.Hypervisor
	} else {
		f.OSExpected = o.OS != ""
		f.OS = o.OS

		f.VersionExpected = o.Version != ""
		f.Version = o.Version

		f.IaaSExpected = o.IaaS != ""
		f.IaaS = o.IaaS

		f.HypervisorExpected = o.Hypervisor != ""
		f.Hypervisor = o.Hypervisor
	}

	return f
}
