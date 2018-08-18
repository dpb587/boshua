package boshreleasedir

import "github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"

type Config struct {
	repository.RepositoryConfig `yaml:"repository"`
	Labels                      []string `yaml:"labels"`
}
