package boshioindex

import "github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"

type Config struct {
	repository.RepositoryConfig `yaml:"repository"`

	Labels []string `yaml:"labels"`
}

func (c *Config) ApplyDefaults() {
	if c.RepositoryConfig.URI == "" {
		c.RepositoryConfig.URI = "https://github.com/bosh-io/releases-index.git"
	}

	c.RepositoryConfig.ApplyDefaults()
}
