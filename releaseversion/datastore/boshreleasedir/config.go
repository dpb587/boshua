package boshreleasedir

import (
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"
)

type Config struct {
	repository.RepositoryConfig `yaml:"repository"`
	Release                     string   `yaml:"release"`
	Labels                      []string `yaml:"labels"`
}
