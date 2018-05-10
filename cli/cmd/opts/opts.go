package opts

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type Opts struct {
	Config string `long:"config" description:"" env:"BOSHUA_CONFIG"`

	Quiet    bool          `long:"quiet" description:"Suppress informational output"`
	LogLevel args.LogLevel `long:"log-level" description:"Show additional levels of log messages" default:"FATAL" env:"BOSHUA_LOG_LEVEL"`

	logger logrus.FieldLogger

	parsedConfig         *config.Config
	releaseIndex         releaseversiondatastore.Index
	compiledReleaseIndex compiledreleaseversiondatastore.Index
	stemcellIndex        stemcellversiondatastore.Index
	osIndex              osversiondatastore.Index
}

func (o *Opts) getParsedConfig() (config.Config, error) {
	if o.parsedConfig != nil {
		return *o.parsedConfig, nil
	}

	o.parsedConfig = &config.Config{}

	configBytes, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".config", "boshua", "config.yml"))
	if err != nil {
		return config.Config{}, fmt.Errorf("reading config: %v", err)
	}

	err = yaml.Unmarshal(configBytes, o.parsedConfig)
	if err != nil {
		return config.Config{}, fmt.Errorf("unmarshalling config: %v", err)
	}

	return *o.parsedConfig, nil
}

func (o *Opts) GetReleaseIndex(name string) (releaseversiondatastore.Index, error) {
	if name != "default" {
		panic("TODO")
	}

	config, err := o.getParsedConfig()
	if err != nil {
		return nil, fmt.Errorf("loading config: %v", err)
	}

	var all []releaseversiondatastore.Index
	factory := releaseversionfactory.New(o.GetLogger())

	for _, cfg := range config.Releases {
		idx, err := factory.Create(cfg.Type, cfg.Name, cfg.Options)
		if err != nil {
			return nil, fmt.Errorf("creating release version datastore: %v", err)
		}

		all = append(all, idx)
	}

	return releaseversionaggregate.New(all...), nil
}

// func (o *Opts) GetCompiledReleaseManager() (*manager.Manager, error) {
// 	rvi, err := o.GetReleaseIndex(name)
// 	if err != nil {
// 		return nil, fmt.Errorf("loading release index: %v", err)
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
		return nil, fmt.Errorf("loading config: %v", err)
	}

	releaseIndex, err := o.GetReleaseIndex("default")
	if err != nil {
		return nil, fmt.Errorf("loading release index: %v", err)
	}

	var all []compiledreleaseversiondatastore.Index
	factory := compiledreleaseversionfactory.New(o.GetLogger(), releaseIndex)

	for _, cfg := range config.CompiledReleases {
		idx, err := factory.Create(cfg.Type, cfg.Name, cfg.Options)
		if err != nil {
			return nil, fmt.Errorf("creating compiled release version datastore: %v", err)
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
		return nil, fmt.Errorf("loading config: %v", err)
	}

	var all []stemcellversiondatastore.Index
	factory := stemcellversionfactory.New(o.GetLogger())

	for _, cfg := range config.Stemcells {
		idx, err := factory.Create(cfg.Type, cfg.Name, cfg.Options)
		if err != nil {
			return nil, fmt.Errorf("creating stemcell version datastore: %v", err)
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
		return nil, fmt.Errorf("loading stemcell index: %v", err)
	}

	return osversionstemcellversionindex.New(stemcellVersionIndex, o.GetLogger()), nil
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
