package opts

import (
	analysisfactory "github.com/dpb587/boshua/analysis/datastore/defaultfactory"
	"github.com/dpb587/boshua/cli/opts"
	configprovider "github.com/dpb587/boshua/config/provider"
	compilationfactory "github.com/dpb587/boshua/releaseversion/compilation/datastore/defaultfactory"
	releaseversionfactory "github.com/dpb587/boshua/releaseversion/datastore/defaultfactory"
	stemcellversionfactory "github.com/dpb587/boshua/stemcellversion/datastore/defaultfactory"
	schedulerfactory "github.com/dpb587/boshua/task/scheduler/factory"
)

type Opts struct {
	*opts.Opts

	parsedConfig *configprovider.Config
}

func (o *Opts) GetConfig() (*configprovider.Config, error) {
	if o.parsedConfig == nil {
		cfg, err := o.Opts.GetConfig()
		if err != nil {
			return nil, err
		}

		cfg.SetAnalysisFactory(analysisfactory.New(cfg.GetLogger()))
		cfg.SetReleaseFactory(releaseversionfactory.New(cfg.GetLogger()))
		cfg.SetReleaseCompilationFactory(compilationfactory.New(cfg.GetLogger()))
		cfg.SetStemcellFactory(stemcellversionfactory.New(cfg.GetLogger()))
		cfg.SetSchedulerFactory(schedulerfactory.New(cfg, cfg.GetLogger()))

		o.parsedConfig = cfg
	}

	return o.parsedConfig, nil
}
