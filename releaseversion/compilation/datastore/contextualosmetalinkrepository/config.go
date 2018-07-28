package contextualosmetalinkrepository

import "github.com/dpb587/boshua/artifact/datastore/datastoreutil/git"

type Config struct {
	git.RepositoryConfig `yaml:",inline"`

	Release string `yaml:"release"`
	Prefix  string `yaml:"prefix"`
}
