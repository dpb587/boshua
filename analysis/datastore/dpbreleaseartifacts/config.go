package boshreleasedpb

import (
	"github.com/dpb587/boshua/blobstore"
	"github.com/dpb587/boshua/datastore/git"
)

type Config struct {
	git.RepositoryConfig      `yaml:",inline"`
	blobstore.BlobstoreConfig `yaml:"blobstore"`

	Release string `yaml:"release"`

	ReleasePrefix         string `yaml:"release_prefix"`
	CompiledReleasePrefix string `yaml:"compiled_release_prefix"`
}
