package opts

import (
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	"github.com/dpb587/boshua/stemcellversion"
	"github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/pkg/errors"
)

type Opts struct {
	AppOpts *cmdopts.Opts `no-flag:"true"`

	Stemcell   *Stemcell `long:"stemcell" description:"The stemcell name and version"`
	OS         string    `long:"stemcell-os" description:"The stemcell OS"`
	Version    string    `long:"stemcell-version" description:"The stemcell version"`
	IaaS       string    `long:"stemcell-iaas" description:"The stemcell IaaS"`
	Hypervisor string    `long:"stemcell-hypervisor" description:"The stemcell hypervisor"`
	Flavor     string    `long:"stemcell-flavor" description:"The stemcell flavor (e.g. 'light')"`

	Labels []string `long:"stemcell-label" description:"The label(s) to filter stemcells by"`
}

func (o *Opts) Artifact() (stemcellversion.Artifact, error) {
	index, err := o.AppOpts.GetStemcellIndex("default")
	if err != nil {
		return stemcellversion.Artifact{}, errors.Wrap(err, "loading stemcell index")
	}

	results, err := index.GetArtifacts(o.FilterParams())
	if err != nil {
		return stemcellversion.Artifact{}, errors.Wrap(err, "finding stemcell")
	}

	result, err := datastore.RequireSingleResult(results)
	if err != nil {
		return stemcellversion.Artifact{}, errors.Wrap(err, "finding stemcell")
	}

	return result, err
}

func (o Opts) FilterParams() datastore.FilterParams {
	f := datastore.FilterParams{
		FlavorExpected: o.Flavor != "",
		Flavor:         o.Flavor,

		LabelsExpected: len(o.Labels) > 0,
		Labels:         o.Labels,
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

	// TODO no default?
	if !f.FlavorExpected {
		f.FlavorExpected = true
		f.Flavor = "heavy"
	}

	return f
}
