package pivnet

import (
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/storage"
)

type Config struct {
	// RepositoryConfig defines how to access the repository.
	repository.RepositoryConfig `yaml:"repository"`

	// StorageConfig defines where results should be stored.
	storage.StorageConfig `yaml:"storage"`
}
