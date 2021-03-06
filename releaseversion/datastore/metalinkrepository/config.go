package metalinkrepository

import "github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"

type Config struct {
	repository.RepositoryConfig `yaml:"repository"`

	Labels  []string `yaml:"labels"`
	Release string   `yaml:"release"`
	Path    string   `yaml:"path"`
}
