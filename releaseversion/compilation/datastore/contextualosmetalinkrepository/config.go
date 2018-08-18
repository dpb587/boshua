package contextualosmetalinkrepository

import (
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/storage"
)

type Config struct {
	repository.RepositoryConfig `yaml:"repository"`
	storage.StorageConfig       `yaml:"storage"`

	Release string `yaml:"release"`
	Path    string `yaml:"path"`
}
