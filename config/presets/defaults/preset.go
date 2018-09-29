package defaults

import (
	"os"

	analysisfactory "github.com/dpb587/boshua/analysis/datastore/defaultfactory"
	"github.com/dpb587/boshua/config/loader"
	"github.com/dpb587/boshua/config/provider"
	compilationfactory "github.com/dpb587/boshua/releaseversion/compilation/datastore/defaultfactory"
	releaseversionfactory "github.com/dpb587/boshua/releaseversion/datastore/defaultfactory"
	pivnetfilefactory "github.com/dpb587/boshua/pivnetfile/datastore/defaultfactory"
	stemcellversionfactory "github.com/dpb587/boshua/stemcellversion/datastore/defaultfactory"
	schedulerfactory "github.com/dpb587/boshua/task/scheduler/factory"
	downloaderurlfactory "github.com/dpb587/boshua/artifact/downloader/url/defaultfactory"
)

func NewConfig() (*provider.Config, error) {
	rawcfg, err := loader.LoadFromFile(os.Getenv("BOSHUA_CONFIG"), nil)
	if err != nil {
		return nil, err
	}

	cfg := &provider.Config{Config: rawcfg}

	UseDefaultFactories(cfg)

	return cfg, nil
}

func UseDefaultFactories(cfg *provider.Config) {
	cfg.SetAnalysisFactory(analysisfactory.New(cfg.GetLogger()))
	cfg.SetReleaseFactory(releaseversionfactory.New(cfg.GetLogger()))
	cfg.SetReleaseCompilationFactory(compilationfactory.New(cfg.GetReleaseIndex, cfg.GetLogger()))
	cfg.SetStemcellFactory(stemcellversionfactory.New(cfg.GetLogger()))
	cfg.SetPivnetFileFactory(pivnetfilefactory.New(cfg.GetLogger()))
	cfg.SetSchedulerFactory(schedulerfactory.New(cfg, cfg.GetLogger()))
	cfg.SetDownloaderURLFactory(downloaderurlfactory.New(cfg.GetLogger()))
}
