package boshreleasedir

import (
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"
)

type Config struct {
	repository.RepositoryConfig `yaml:"repository"`
	Release                     string   `yaml:"release"`
	Labels                      []string `yaml:"labels"` // TODO no way for labels to be applied to only final/non-dev releases
	DevReleases                 bool     `yaml:"dev_releases"`
	DevLabels                   []string `yaml:"dev_labels"`
}
