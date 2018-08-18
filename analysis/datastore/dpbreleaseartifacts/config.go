package dpbreleaseartifacts

import (
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/storage"
)

type Config struct {
	repository.RepositoryConfig `yaml:"repository"`
	storage.StorageConfig       `yaml:"storage"`

	Release string `yaml:"release"`

	CompiledReleasePath string `yaml:"compiled_release_path"`
	ReleasePath         string `yaml:"release_path"`
	StemcellPath        string `yaml:"stemcell_path"`
}
