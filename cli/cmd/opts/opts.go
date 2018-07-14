package opts

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	analysisfactory "github.com/dpb587/boshua/analysis/datastore/factory"
	"github.com/dpb587/boshua/cli/args"
	compiledreleaseversiondatastore "github.com/dpb587/boshua/compiledreleaseversion/datastore"
	compiledreleaseversionaggregate "github.com/dpb587/boshua/compiledreleaseversion/datastore/aggregate"
	compiledreleaseversionfactory "github.com/dpb587/boshua/compiledreleaseversion/datastore/factory"
	"github.com/dpb587/boshua/config"
	osversiondatastore "github.com/dpb587/boshua/osversion/datastore"
	osversionstemcellversionindex "github.com/dpb587/boshua/osversion/datastore/stemcellversionindex"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	releaseversionaggregate "github.com/dpb587/boshua/releaseversion/datastore/aggregate"
	releaseversionfactory "github.com/dpb587/boshua/releaseversion/datastore/factory"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	stemcellversionaggregate "github.com/dpb587/boshua/stemcellversion/datastore/aggregate"
	stemcellversionfactory "github.com/dpb587/boshua/stemcellversion/datastore/factory"
	"github.com/dpb587/boshua/task/scheduler"
	schedulerfactory "github.com/dpb587/boshua/task/scheduler/factory"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Opts struct {
	Config string `long:"config" description:"Path to configuration file" env:"BOSHUA_CONFIG" default:"~/.config/boshua/config.yml"`

	Quiet    bool          `long:"quiet" description:"Suppress informational output"`
	LogLevel args.LogLevel `long:"log-level" description:"Show additional levels of log messages" default:"FATAL" env:"BOSHUA_LOG_LEVEL"`

	logger logrus.FieldLogger

	parsedConfig         *config.Config
	releaseIndex         releaseversiondatastore.Index
	compiledReleaseIndex compiledreleaseversiondatastore.Index
	stemcellIndex        stemcellversiondatastore.Index
	osIndex              osversiondatastore.Index
}

func (o *Opts) getConfigPath() string {
	configPath := o.Config

	if strings.HasPrefix(configPath, "~/") {
		configPath = filepath.Join(os.Getenv("HOME"), configPath[1:])
	}

	configPath, err := filepath.Abs(configPath)
	if err != nil {
		panic(err)
	}

	return configPath
}

func (o *Opts) getParsedConfig() (config.Config, error) {
	if o.parsedConfig != nil {
		return *o.parsedConfig, nil
	}

	configBytes, err := ioutil.ReadFile(o.getConfigPath())
	if err != nil {
		return config.Config{}, errors.Wrap(err, "reading config")
	}

	cfg := config.Config{}
	err = config.UnmarshalYAML(configBytes, &cfg)
	if err != nil {
		return config.Config{}, errors.Wrap(err, "loading options")
	}

	o.parsedConfig = &cfg

	return *o.parsedConfig, nil
}

func (o *Opts) GetAnalysisIndex(name string) (analysisdatastore.Index, error) {
	config, err := o.getParsedConfig()
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	factory := analysisfactory.New(o.GetLogger())

	for _, cfg := range config.Analyses {
		if cfg.Name != name {
			continue
		}

		idx, err := factory.Create(cfg.Type, cfg.Name, cfg.Options)
		if err != nil {
			return nil, errors.Wrap(err, "creating analysis datastore")
		}

		return idx, nil
	}

	return nil, fmt.Errorf("failed to find analysis index: %s", name)
}

func (o *Opts) GetReleaseIndex(name string) (releaseversiondatastore.Index, error) {
	if name != "default" {
		panic("TODO")
	}

	config, err := o.getParsedConfig()
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	var all []releaseversiondatastore.Index
	factory := releaseversionfactory.New(o.GetLogger())

	for _, cfg := range config.Releases {
		idx, err := factory.Create(cfg.Type, cfg.Name, cfg.Options)
		if err != nil {
			return nil, errors.Wrap(err, "creating release version datastore")
		}

		all = append(all, idx)
	}

	return releaseversionaggregate.New(all...), nil
}

// func (o *Opts) GetCompiledReleaseManager() (*manager.Manager, error) {
// 	rvi, err := o.GetReleaseIndex(name)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "loading release index")
// 	}
//
// 	return manager.NewManager(rv, ov)
// }

func (o *Opts) GetCompiledReleaseIndex(name string) (compiledreleaseversiondatastore.Index, error) {
	if name != "default" {
		panic("TODO")
	}

	config, err := o.getParsedConfig()
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	releaseIndex, err := o.GetReleaseIndex("default")
	if err != nil {
		return nil, errors.Wrap(err, "loading release index")
	}

	var all []compiledreleaseversiondatastore.Index
	factory := compiledreleaseversionfactory.New(o.GetLogger(), releaseIndex)

	for _, cfg := range config.CompiledReleases {
		var analysisIndex analysisdatastore.Index

		if cfg.Analysis != nil {
			analysisIndex, err = o.GetAnalysisIndex(cfg.Analysis.Name)
			if err != nil {
				return nil, errors.Wrap(err, "loading compiled release analysis")
			}
		}

		idx, err := factory.Create(cfg.Type, cfg.Name, cfg.Options, analysisIndex)
		if err != nil {
			return nil, errors.Wrap(err, "creating compiled release version datastore")
		}

		all = append(all, idx)
	}

	return compiledreleaseversionaggregate.New(all...), nil
}

func (o *Opts) GetStemcellIndex(name string) (stemcellversiondatastore.Index, error) {
	if name != "default" {
		panic("TODO")
	}

	config, err := o.getParsedConfig()
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	var all []stemcellversiondatastore.Index
	factory := stemcellversionfactory.New(o.GetLogger())

	for _, cfg := range config.Stemcells {
		var analysisIndex analysisdatastore.Index

		if cfg.Analysis != nil {
			analysisIndex, err = o.GetAnalysisIndex(cfg.Analysis.Name)
			if err != nil {
				return nil, errors.Wrap(err, "loading compiled release analysis")
			}
		}

		idx, err := factory.Create(cfg.Type, cfg.Name, cfg.Options, analysisIndex)
		if err != nil {
			return nil, errors.Wrap(err, "creating stemcell version datastore")
		}

		all = append(all, idx)
	}

	return stemcellversionaggregate.New(all...), nil
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
	config, err := o.getParsedConfig()
	if err != nil {
		return nil, errors.Wrap(err, "loading config")
	}

	factory := schedulerfactory.New(
		o.GetLogger(),
		func(args ...string) *exec.Cmd {
			return exec.Command(
				"boshua",
				append([]string{
					"--config", o.getConfigPath(),
					// "--log-level", string(o.LogLevel), // TODO stringify
				}, args...)...,
			)
		},
	)

	return factory.Create(config.Scheduler.Type, config.Scheduler.Options)
}

func (o *Opts) GetServerConfig() (config.ServerConfig, error) {
	parsed, err := o.getParsedConfig()
	if err != nil {
		return config.ServerConfig{}, errors.Wrap(err, "loading config")
	}

	return parsed.Server, nil
}

func (o *Opts) GetLogger() logrus.FieldLogger {
	if o.logger == nil {
		panic("logger is not configured")
	}

	return o.logger
}

func (o *Opts) ConfigureLogger(command string) {
	if o.logger != nil {
		panic("logger is already configured")
	}

	var logger = logrus.New()
	logger.Out = os.Stderr
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logrus.Level(o.LogLevel)

	o.logger = logger.WithField("cli.command", command)
}
