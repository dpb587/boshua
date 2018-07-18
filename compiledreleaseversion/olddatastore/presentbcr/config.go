package presentbcr

import "github.com/dpb587/boshua/datastore/git"

type Config struct {
	git.RepositoryConfig `yaml:",inline"`
}
