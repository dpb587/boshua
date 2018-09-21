package cfdeployment

import (
	"time"

	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"
	"github.com/dpb587/boshua/util/marshaltypes"
)

type Config struct {
	repository.RepositoryConfig `yaml:"repository"`
}

func (c *Config) ApplyDefaults() {
	if c.RepositoryConfig.URI == "" {
		c.RepositoryConfig.URI = "https://github.com/cloudfoundry/cf-deployment.git"
	}

	if c.RepositoryConfig.PullInterval == nil {
		// this datastore is more expensive to rebuild
		d := marshaltypes.Duration(15 * time.Minute)
		c.RepositoryConfig.PullInterval = &d
	}

	c.RepositoryConfig.ApplyDefaults()
}
