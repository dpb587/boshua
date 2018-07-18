package opts

import (
	"time"

	"github.com/dpb587/boshua/cli/args"
	cmdopts "github.com/dpb587/boshua/cli/cmd/opts"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	releaseversionopts "github.com/dpb587/boshua/releaseversion/cli/opts"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	"github.com/pkg/errors"
)

type Opts struct {
	AppOpts *cmdopts.Opts `no-flag:"true"`

	ReleaseOpts *releaseversionopts.Opts `no-flag:"true"`
	OS          args.OS                  `long:"os" description:"The OS and version"`

	NoWait      bool          `long:"no-wait" description:"Do not request and wait for compilation if not already available"`
	WaitTimeout time.Duration `long:"wait-timeout" description:"Timeout duration when waiting for compilations" default:"30m"`
}

func (o *Opts) Artifact() (compilation.Artifact, error) {
	index, err := o.AppOpts.GetCompiledReleaseIndex("default")
	if err != nil {
		return compilation.Artifact{}, errors.Wrap(err, "loading compiled release index")
	}

	results, err := index.Filter(o.FilterParams())
	if err != nil {
		return compilation.Artifact{}, errors.Wrap(err, "finding compiled release")
	}

	result, err := datastore.RequireSingleResult(results)
	if err != nil {
		return compilation.Artifact{}, errors.Wrap(err, "finding compiled release")
	}

	return result, err
}

func (o Opts) FilterParams() *datastore.FilterParams {
	return &datastore.FilterParams{
		Release: o.ReleaseOpts.FilterParams(),
		OS: &osversiondatastore.FilterParams{
			NameExpected:    o.OS.Name != "",
			Name:            o.OS.Name,
			VersionExpected: o.OS.Version != "",
			Version:         o.OS.Version,
		},
	}
}
