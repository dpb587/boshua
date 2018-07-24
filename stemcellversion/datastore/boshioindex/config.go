package boshioindex

import "github.com/dpb587/boshua/datastore/git"

type Config struct {
	git.RepositoryConfig `yaml:",inline"`

	Prefix string `yaml:"prefix"`
}
