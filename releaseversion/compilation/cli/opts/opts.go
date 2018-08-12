package opts

import (
	"fmt"
	"time"

	"github.com/dpb587/boshua/cli/args"
	cmdopts "github.com/dpb587/boshua/main/boshua/cmd/opts"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	releaseversionopts "github.com/dpb587/boshua/releaseversion/cli/opts"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	schedulerpkg "github.com/dpb587/boshua/task/scheduler"
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

	f := o.FilterParams()

	result, err := datastore.GetCompilationArtifact(index, f)
	if err == datastore.NoMatchErr {
		if o.NoWait {
			return compilation.Artifact{}, errors.New("none found")
		}

		scheduler, err := o.AppOpts.GetScheduler()
		if err != nil {
			return compilation.Artifact{}, errors.Wrap(err, "loading scheduler")
		}

		scheduledTask, err := scheduler.ScheduleCompilation(f)
		if err != nil {
			return compilation.Artifact{}, errors.Wrap(err, "creating compilation")
		}

		status, err := schedulerpkg.WaitForScheduledTask(scheduledTask, schedulerpkg.DefaultStatusChangeCallback)
		if err != nil {
			return compilation.Artifact{}, errors.Wrap(err, "checking task")
		} else if status != schedulerpkg.StatusSucceeded {
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
