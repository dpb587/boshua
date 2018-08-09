package opts

import (
	"fmt"
	"os"
	"time"

	"github.com/dpb587/boshua/cli/args"
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	releaseversionopts "github.com/dpb587/boshua/releaseversion/cli/opts"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/dpb587/boshua/task"
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

	result, err := datastore.GetCompilationArtifact(index, o.FilterParams())
	if err == datastore.NoMatchErr {
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

		stemcellVersion, err := stemcellversiondatastore.GetArtifact(stemcellVersionIndex, stemcellversiondatastore.FilterParams{
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

		scheduler, err := o.AppOpts.GetScheduler()
		if err != nil {
			return compilation.Artifact{}, errors.Wrap(err, "loading scheduler")
		}

		scheduledTask, err := scheduler.ScheduleCompilation(releaseVersion, stemcellVersion)
		if err != nil {
			return compilation.Artifact{}, errors.Wrap(err, "creating compilation")
		}

		status, err := task.WaitForScheduledTask(scheduledTask, func(status task.Status) {
			if o.AppOpts.Quiet {
				return
			}

			fmt.Fprintf(os.Stderr, "%s [%s/%s %s/%s] compilation is %s\n", time.Now().Format("15:04:05"), stemcellVersion.OS, stemcellVersion.Version, releaseVersion.Name, releaseVersion.Version, status)
		})
		if err != nil {
			return compilation.Artifact{}, errors.Wrap(err, "checking task")
		} else if status != task.StatusSucceeded {
			return compilation.Artifact{}, fmt.Errorf("task did not succeed: %s", status)
		}

		result, err = datastore.GetCompilationArtifact(index, o.FilterParams())
		if err != nil {
			return compilation.Artifact{}, errors.Wrap(err, "finding finished compilation")
		}
	} else if err != nil {
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
