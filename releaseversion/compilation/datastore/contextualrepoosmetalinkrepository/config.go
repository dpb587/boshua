package contextualrepoosmetalinkrepository

import (
	"github.com/dpb587/boshua/artifact/datastore/datastoreutil/git"
	"github.com/dpb587/boshua/blobstore"
)

type Config struct {
	git.RepositoryConfig      `yaml:",inline"`
	blobstore.BlobstoreConfig `yaml:"blobstore"`

	Prefix string `yaml:"prefix"`
}
