package provider

import (
	"os"

	analysisdatastore "github.com/dpb587/boshua/analysis/datastore"
	"github.com/dpb587/boshua/config"
	compilationdatastore "github.com/dpb587/boshua/releaseversion/compilation/datastore"
	releaseversiondatastore "github.com/dpb587/boshua/releaseversion/datastore"
	stemcellversiondatastore "github.com/dpb587/boshua/stemcellversion/datastore"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	*config.Config
	RawConfig []byte

	logger                    logrus.FieldLogger
	releaseFactory            releaseversiondatastore.Factory
	releaseCompilationFactory compilationdatastore.Factory
	stemcellFactory           stemcellversiondatastore.Factory
	analysisFactory           analysisdatastore.Factory
}

func (c *Config) Marshal() ([]byte, error) {
	if c.RawConfig != nil {
		return c.RawConfig, nil
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
	logger.Level = logrus.Level(c.Config.General.LogLevel)

	c.logger = logger
}
