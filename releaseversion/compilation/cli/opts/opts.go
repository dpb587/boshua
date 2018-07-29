package opts

import (
	"time"

	"github.com/dpb587/boshua/cli/args"
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	releaseversionopts "github.com/dpb587/boshua/releaseversion/cli/opts"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/boshua/task/scheduler/schedulerutil"
	"github.com/pkg/errors"
)

type Opts struct {
	AppOpts *cmdopts.Opts `no-flag:"true"`

	ReleaseOpts *releaseversionopts.Opts `no-flag:"true"`
	OS          args.OS                  `long:"os" description:"The OS in name/version format"`

	NoWait      bool          `long:"no-wait" description:"Do not request and wait for compilation if not already available"`
	WaitTimeout time.Duration `long:"wait-timeout" description:"Timeout duration when waiting for compilations" default:"30m"`
}

func (o *Opts) Artifact() (compilation.Artifact, error) {
	index, err := o.AppOpts.GetCompiledReleaseIndex("default")
	if err != nil {
		return compilation.Artifact{}, errors.Wrap(err, "loading index")
	}

	results, err := index.GetCompilationArtifacts(o.FilterParams())
	if err != nil {
		return compilation.Artifact{}, errors.Wrap(err, "filtering")
	}

	if len(results) == 0 {
		if o.NoWait {
			return compilation.Artifact{}, errors.New("none found")
		}

		releaseVersion, err := o.ReleaseOpts.Artifact()
		if err != nil {
			return compilation.Artifact{}, errors.New("finding release")
		}

		stemcellVersionIndex, err := o.AppOpts.GetStemcellIndex("default")
		if err != nil {
			return compilation.Artifact{}, errors.Wrap(err, "loading stemcell index")
		}

		stemcellVersions, err := stemcellVersionIndex.GetArtifacts(stemcellversiondatastore.FilterParams{
			OSExpected:      true,
			OS:              o.OS.Name,
			VersionExpected: true,
			Version:         o.OS.Version,
			// TODO dynamic
			IaaSExpected:   true,
			IaaS:           "aws",
			FlavorExpected: true,
			Flavor:         "light",
		})
		if err != nil {
			return compilation.Artifact{}, errors.Wrap(err, "filtering stemcell")
		}

		stemcellVersion, err := stemcellversiondatastore.RequireSingleResult(stemcellVersions)
		if err != nil {
			return compilation.Artifact{}, errors.Wrap(err, "filtering stemcell")
		}

		scheduler, err := o.AppOpts.GetScheduler()
		if err != nil {
			return compilation.Artifact{}, errors.Wrap(err, "loading scheduler")
		}

		err = schedulerutil.CreateCompilation(scheduler, releaseVersion, stemcellVersion)
		if err != nil {
			return compilation.Artifact{}, errors.Wrap(err, "creating compilation")
		}

		results, err = index.GetCompilationArtifacts(o.FilterParams())
		if err != nil {
			return compilation.Artifact{}, errors.Wrap(err, "finding finished compilation")
		}
	}

	result, err := datastore.RequireSingleResult(results)
	if err != nil {
		return compilation.Artifact{}, errors.Wrap(err, "filtering")
	}

	return result, nil
}

func (o Opts) FilterParams() datastore.FilterParams {
	return datastore.FilterParams{
		Release: o.ReleaseOpts.FilterParams(),
		OS: osversiondatastore.FilterParams{
			NameExpected:    o.OS.Name != "",
			Name:            o.OS.Name,
			VersionExpected: o.OS.Version != "",
			Version:         o.OS.Version,
		},
	}
}
