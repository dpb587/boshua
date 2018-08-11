package opts

import (
	"github.com/dpb587/boshua/analysis"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	analysisfactory "github.com/dpb587/boshua/analysis/datastore/defaultfactory"
	"github.com/dpb587/boshua/cli/opts"
	configprovider "github.com/dpb587/boshua/config/provider"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	osversionstemcellversionindex "github.com/dpb587/boshua/osversion/datastore/stemcellversionindex"
	compiledreleaseversiondatastore "github.com/dpb587/boshua/releaseversion/compilation/datastore"
	compilationfactory "github.com/dpb587/boshua/releaseversion/compilation/datastore/defaultfactory"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	releaseversionfactory "github.com/dpb587/boshua/releaseversion/datastore/defaultfactory"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	stemcellversionfactory "github.com/dpb587/boshua/stemcellversion/datastore/defaultfactory"
	"github.com/dpb587/boshua/task/scheduler"
	schedulerfactory "github.com/dpb587/boshua/task/scheduler/factory"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
		cfg.SetSchedulerFactory(schedulerfactory.New(cfg.Marshal, cfg.GetLogger()))

		o.parsedConfig = cfg
	}

	return o.parsedConfig, nil
}

func (o *Opts) GetAnalysisIndex(r analysis.Reference) (analysisdatastore.Index, error) {
	o.mustConfig() // TODO cleaner error

	return o.parsedConfig.GetAnalysisIndex(r)
}

func (o *Opts) GetReleaseIndex(name string) (releaseversiondatastore.Index, error) {
	o.mustConfig() // TODO cleaner error

	return o.parsedConfig.GetReleaseIndex(name)
}

func (o *Opts) GetCompiledReleaseIndex(name string) (compiledreleaseversiondatastore.Index, error) {
	o.mustConfig() // TODO cleaner error

	return o.parsedConfig.GetCompiledReleaseIndex(name)
}

func (o *Opts) GetStemcellIndex(name string) (stemcellversiondatastore.Index, error) {
	o.mustConfig() // TODO cleaner error

	return o.parsedConfig.GetStemcellIndex(name)
}

func (o *Opts) GetOSIndex(name string) (osversiondatastore.Index, error) {
	if name != "default" {
		panic("TODO")
	}

	stemcellVersionIndex, err := o.GetStemcellIndex("default")
	if err != nil {
		return nil, errors.Wrap(err, "loading stemcell index")
	}

	return osversionstemcellversionindex.New(stemcellVersionIndex, o.GetLogger()), nil
}

func (o *Opts) GetScheduler() (scheduler.Scheduler, error) {
	o.mustConfig()

	return o.parsedConfig.GetScheduler()
}

func (o *Opts) GetLogger() logrus.FieldLogger {
	o.mustConfig()

	return o.parsedConfig.GetLogger()
}

func (o *Opts) ConfigureLogger(command string) {
	o.mustConfig()

	o.parsedConfig.AppendLoggerFields(logrus.Fields{"cli.command": command})
}

func (o *Opts) mustConfig() {
	_, err := o.GetConfig()
	if err != nil {
		panic(errors.Wrap(err, "loading config"))
	}
}
