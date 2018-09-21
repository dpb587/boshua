package cfdeployment

import (
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"
)

type Config struct {
	repository.RepositoryConfig `yaml:"repository"`
}

func (c *Config) ApplyDefaults() {
	if c.RepositoryConfig.URI == "" {
		c.RepositoryConfig.URI = "https://github.com/cloudfoundry/cf-deployment.git"
	}

	c.RepositoryConfig.ApplyDefaults()
}
