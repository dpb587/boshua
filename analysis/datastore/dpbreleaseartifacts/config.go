package dpbreleaseartifacts

import (
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/repository"
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/storage"
)

type Config struct {
	// RepositoryConfig defines how to access the repository.
	repository.RepositoryConfig `yaml:"repository"`

	// StorageConfig defines where results should be stored.
	storage.StorageConfig `yaml:"storage"`

	// Release defines a static release name for release-related results.
	Release string `yaml:"release"`

	// ReleasePath defines a custom prefix when storing release analyses.
	ReleasePath string `yaml:"release_path"`

	// ReleaseCompilationPath defines a custom prefix when storing release
	// compilation analyses.
	ReleaseCompilationPath string `yaml:"release_compilation_path"`

	// StemcellPath defines a custom prefix when storing stemcell analyses.
	StemcellPath string `yaml:"stemcell_path"`
}
