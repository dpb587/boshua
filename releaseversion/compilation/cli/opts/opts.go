package opts

import (
	"github.com/dpb587/boshua/config"
	"github.com/dpb587/boshua/config/provider"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	releaseversionopts "github.com/dpb587/boshua/releaseversion/cli/opts"
	"github.com/dpb587/boshua/releaseversion/compilation"
	"github.com/dpb587/boshua/releaseversion/compilation/datastore"
	stemcellversionopts "github.com/dpb587/boshua/stemcellversion/cli/opts"
	"github.com/pkg/errors"
)

type Opts struct {
	ReleaseOpts  *releaseversionopts.Opts `no-flag:"true"`
	StemcellOpts *stemcellversionopts.Opts
}

func (o *Opts) Artifact(cfg *provider.Config) (compilation.Artifact, error) {
	index, err := cfg.GetReleaseCompilationIndex(config.DefaultName)
	if err != nil {
		return compilation.Artifact{}, errors.Wrap(err, "loading release compilation index")
	}

	f := o.FilterParams()

	result, err := datastore.GetCompilationArtifact(index, f)
	if err != nil {
		return compilation.Artifact{}, errors.Wrap(err, "getting compilation artifact")
	}

	return result, nil
}

func (o Opts) FilterParams() datastore.FilterParams {
	releaseFilterParams, _ := o.ReleaseOpts.ArtifactParams()
	stemcellFilterParams, _ := o.StemcellOpts.ArtifactParams()

	return datastore.FilterParams{
		Release: releaseFilterParams,
		OS: osversiondatastore.FilterParams{
			NameExpected:      stemcellFilterParams.OSExpected,
			Name:              stemcellFilterParams.OS,
			VersionExpected:   stemcellFilterParams.VersionExpected,
			VersionConstraint: stemcellFilterParams.VersionConstraint,
			Version:           stemcellFilterParams.Version,
		},
	}
}
