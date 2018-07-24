package boshioindex

import "github.com/dpb587/boshua/datastore/git"

type Config struct {
	git.RepositoryConfig `yaml:",inline"`

	Labels []string `yaml:"labels"`
	Prefix string   `yaml:"prefix"`
}
