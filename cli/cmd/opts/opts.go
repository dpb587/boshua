package opts

import (
	"time"

	"github.com/dpb587/boshua/analysis"
	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	analysisfactory "github.com/dpb587/boshua/analysis/datastore/factory"
	"github.com/dpb587/boshua/cli/args"
	"github.com/dpb587/boshua/config"
	configloader "github.com/dpb587/boshua/config/loader"
	configprovider "github.com/dpb587/boshua/config/provider"
	configtypes "github.com/dpb587/boshua/config/types"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	osversionstemcellversionindex "github.com/dpb587/boshua/osversion/datastore/stemcellversionindex"
	compiledreleaseversiondatastore "github.com/dpb587/boshua/releaseversion/compilation/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	releaseversionfactory "github.com/dpb587/boshua/releaseversion/datastore/factory"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	stemcellversionfactory "github.com/dpb587/boshua/stemcellversion/datastore/factory"
	"github.com/dpb587/boshua/task/scheduler"
	schedulerfactory "github.com/dpb587/boshua/task/scheduler/factory"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Opts struct {
	Config string `long:"config" description:"Path to configuration file" env:"BOSHUA_CONFIG" default:"~/.config/boshua/config.yml"`

	DefaultServer string        `long:"default-server" description:"Default boshua API server" env:"BOSHUA_SERVER"`
	DefaultWait   args.Duration `long:"default-wait" description:"Maximum time to wait for scheduled tasks; 0 to disable scheduling" env:"BOSHUA_WAIT" default:"30m"` // TODO better name

	Quiet    bool          `long:"quiet" description:"Suppress informational output" env:"BOSHUA_QUIET"`
	LogLevel args.LogLevel `long:"log-level" description:"Show additional levels of log messages" default:"FATAL" env:"BOSHUA_LOG_LEVEL"`

	parsedConfig *configprovider.Config
}

func (o *Opts) GetConfig() (*configprovider.Config, error) {
	if o.parsedConfig != nil {
		return o.parsedConfig, nil
	}

	cfg, err := configloader.LoadFromFile(
		o.Config,
		&config.Config{
			General: config.GeneralConfig{
				DefaultServer: o.DefaultServer,
				DefaultWait:   time.Duration(o.DefaultWait),
				LogLevel:      configtypes.LogLevel(o.LogLevel),
			},
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "loading file")
	}

	cfg.SetAnalysisFactory(analysisfactory.New(cfg.GetLogger()))
	cfg.SetReleaseFactory(releaseversionfactory.New(cfg.GetLogger()))
	cfg.SetStemcellFactory(stemcellversionfactory.New(cfg.GetLogger()))

	o.parsedConfig = cfg

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
	config, err := o.GetConfig()
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	factory := schedulerfactory.New(config.Marshal, o.GetLogger())

	return factory.Create(config.Scheduler.Type, config.Scheduler.Options)
}

func (o *Opts) GetServerConfig() (config.ServerConfig, error) {
	parsed, err := o.GetConfig()
	if err != nil {
		return config.ServerConfig{}, errors.Wrap(err, "loading config")
	}

	return parsed.Server, nil
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
