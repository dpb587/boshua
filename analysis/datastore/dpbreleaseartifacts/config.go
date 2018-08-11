package dpbreleaseartifacts

import (
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/git"
	"github.com/dpb587/boshua/blobstore"
)

type Config struct {
	git.RepositoryConfig      `yaml:",inline"`
	blobstore.BlobstoreConfig `yaml:"blobstore"`

	Release string `yaml:"release"`

	CompiledReleasePrefix string `yaml:"compiled_release_prefix"`
	ReleasePrefix         string `yaml:"release_prefix"`
	StemcellPrefix        string `yaml:"stemcell_prefix"`
}
