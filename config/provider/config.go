package provider

import (
	"os"

	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/config"
	pivnetfiledatastore "github.com/dpb587/boshua/pivnetfile/datastore"
	compilationdatastore "github.com/dpb587/boshua/releaseversion/compilation/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	downloaderurl "github.com/dpb587/boshua/artifact/downloader/url"
	"github.com/dpb587/boshua/task/scheduler"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	*config.Config

	logger logrus.FieldLogger

	releaseFactory releaseversiondatastore.Factory
	releaseIndices map[string]releaseversiondatastore.Index

	releaseCompilationFactory compilationdatastore.Factory
	releaseCompilationIndices map[string]compilationdatastore.Index

	stemcellFactory stemcellversiondatastore.Factory
	stemcellIndices map[string]stemcellversiondatastore.Index

	pivnetFileFactory pivnetfiledatastore.Factory
	pivnetFileIndices map[string]pivnetfiledatastore.Index


	analysisFactory analysisdatastore.Factory
	analysisIndices map[string]analysisdatastore.Index

	schedulerFactory scheduler.Factory
	scheduler        scheduler.Scheduler

	downloaderUrlFactory downloaderurl.Factory
}

func (c *Config) Marshal() ([]byte, error) {
	if c.Config.RawConfig != nil {
		return c.Config.RawConfig()
	}

	return yaml.Marshal(c.Config)
}

func (c *Config) GetLogger() logrus.FieldLogger {
	c.requireLogger()

	return c.logger
}

func (c *Config) AppendLoggerFields(fields logrus.Fields) {
	c.requireLogger()

	c.logger = c.logger.WithFields(fields)
}

func (c *Config) requireLogger() {
	if c.logger != nil {
		return
	}

	var logger = logrus.New()
	logger.Out = os.Stderr
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logrus.Level(c.Config.Global.LogLevel)

	c.logger = logger
}
