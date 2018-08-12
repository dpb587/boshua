package defaults

import (
	"os"

	analysisfactory "github.com/dpb587/boshua/analysis/datastore/defaultfactory"
	"github.com/dpb587/boshua/config/loader"
	"github.com/dpb587/boshua/config/provider"
	compilationfactory "github.com/dpb587/boshua/releaseversion/compilation/datastore/defaultfactory"
	releaseversionfactory "github.com/dpb587/boshua/releaseversion/datastore/defaultfactory"
	stemcellversionfactory "github.com/dpb587/boshua/stemcellversion/datastore/defaultfactory"
	schedulerfactory "github.com/dpb587/boshua/task/scheduler/factory"
)

func NewConfig() (*provider.Config, error) {
	cfg, err := loader.LoadFromFile(os.Getenv("BOSHUA_CONFIG"), nil)
	if err != nil {
		return nil, err
	}

	cfg.SetAnalysisFactory(analysisfactory.New(cfg.GetLogger()))
	cfg.SetReleaseFactory(releaseversionfactory.New(cfg.GetLogger()))
	cfg.SetReleaseCompilationFactory(compilationfactory.New(cfg.GetLogger()))
	cfg.SetStemcellFactory(stemcellversionfactory.New(cfg.GetLogger()))
	cfg.SetSchedulerFactory(schedulerfactory.New(cfg, cfg.GetLogger()))

	return cfg, nil
}
