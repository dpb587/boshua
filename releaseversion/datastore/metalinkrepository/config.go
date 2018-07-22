package metalinkrepository

import "github.com/dpb587/boshua/datastore/git"

type Config struct {
	git.RepositoryConfig `yaml:",inline"`

	Release string `yaml:"release"`
	Prefix  string `yaml:"prefix"`
}
