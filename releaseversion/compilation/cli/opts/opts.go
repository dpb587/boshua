package opts

import (
	"github.com/dpb587/boshua/cli/args"
	"github.com/dpb587/boshua/config/provider"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	releaseversionopts "github.com/dpb587/boshua/releaseversion/cli/opts"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	"github.com/pkg/errors"
)

type Opts struct {
	ReleaseOpts *releaseversionopts.Opts `no-flag:"true"`
	OS          args.OS                  `long:"os" description:"The OS in name/version format"`
}

func (o *Opts) Artifact(config *provider.Config) (compilation.Artifact, error) {
	index, err := config.GetCompiledReleaseIndexScheduler("default")
	if err != nil {
		return compilation.Artifact{}, errors.Wrap(err, "loading index")
	}

	f := o.FilterParams()

	result, err := datastore.GetCompilationArtifact(index, f)
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
