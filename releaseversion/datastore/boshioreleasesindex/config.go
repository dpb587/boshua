package boshioreleasesindex

import "github.com/dpb587/boshua/artifact/datastore/datastoreutil/git"

type Config struct {
	git.RepositoryConfig `yaml:",inline"`

	Labels []string `yaml:"labels"`
}
