package meta4

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dpb587/boshua/datastore/git"
	"github.com/dpb587/boshua/metalink/metalinkutil"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/boshua/releaseversion/datastore/inmemory"
	"github.com/dpb587/metalink"
	"github.com/sirupsen/logrus"
)

type index struct {
	logger     logrus.FieldLogger
	config     Config
	repository *git.Repository
	inmemory   datastore.Index
}

var _ datastore.Index = &index{}

func New(config Config, logger logrus.FieldLogger) datastore.Index {
	idx := &index{
		logger:     logger.WithField("build.package", reflect.TypeOf(index{}).PkgPath()),
		config:     config,
		repository: git.NewRepository(logger, config.RepositoryConfig),
	}

	idx.inmemory = inmemory.New(idx.loader, idx.repository.Reload)

	return idx
}

func (i *index) Filter(ref releaseversion.Reference) ([]releaseversion.Artifact, error) {
	return i.inmemory.Filter(ref)
}

func (i *index) Find(ref releaseversion.Reference) (releaseversion.Artifact, error) {
	return i.inmemory.Find(ref)
}

func (i *index) loader() ([]releaseversion.Artifact, error) {
	paths, err := filepath.Glob(filepath.Join(i.config.RepositoryConfig.LocalPath, "*.meta4"))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

	var inmemory = []releaseversion.Artifact{}

	for _, meta4Path := range paths {
		meta4Bytes, err := ioutil.ReadFile(meta4Path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", meta4Path, err)
		}

		var meta4 metalink.Metalink

		err = metalink.Unmarshal(meta4Bytes, &meta4)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling %s: %v", meta4Path, err)
		}

		meta4File := meta4.Files[0]

		ref := releaseversion.Reference{
			Name:      i.config.Release,
			Version:   meta4File.Version,
			Checksums: metalinkutil.HashesToChecksums(meta4File.Hashes),
		}

		inmemory = append(
			inmemory,
			releaseversion.New(
				ref,
				meta4File,
				map[string]interface{}{
					"uri": fmt.Sprintf(
						"%s//%s",
						i.config.RepositoryConfig.Repository,
						strings.TrimPrefix(path.Dir(strings.TrimPrefix(meta4Path, i.config.RepositoryConfig.LocalPath)), "/"),
					),
					"version": ref.Version,
				},
			),
		)
	}

	i.logger.Infof("found %d entries", len(inmemory))

	return inmemory, nil
}
