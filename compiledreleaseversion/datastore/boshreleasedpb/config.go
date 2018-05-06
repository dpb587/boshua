package boshreleasedpb

import "github.com/dpb587/boshua/datastore/git"

type Config struct {
	git.RepositoryConfig `yaml:",inline"`

	Release string `yaml:"release"`
	Channel string `yaml:"channel"`
}
