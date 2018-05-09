package presentbcr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dpb587/boshua/compiledreleaseversion"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore"
	"github.com/dpb587/boshua/compiledreleaseversion/datastore/inmemory"
	"github.com/dpb587/boshua/datastore/git"
	"github.com/dpb587/boshua/osversion"
	"github.com/dpb587/boshua/releaseversion"
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

func (i *index) Filter(ref compiledreleaseversion.Reference) ([]compiledreleaseversion.Artifact, error) {
	return i.inmemory.Filter(ref)
}

func (i *index) Find(ref compiledreleaseversion.Reference) (compiledreleaseversion.Artifact, error) {
	return i.inmemory.Find(ref)
}

func (i *index) Store(artifact compiledreleaseversion.Artifact) error {
	return datastore.UnsupportedOperationErr
}

func (i *index) loader() ([]compiledreleaseversion.Artifact, error) {
	paths, err := filepath.Glob(fmt.Sprintf("%s/compiledreleaseversion/**/**/**/reference.json", i.config.RepositoryConfig.LocalPath))
	if err != nil {
		return nil, fmt.Errorf("globbing: %v", err)
	}

	var inmemory = []compiledreleaseversion.Artifact{}

	for _, bcrPath := range paths {
		bcrBytes, err := ioutil.ReadFile(bcrPath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", bcrPath, err)
		}

		var bcrJSON Record

		err = json.Unmarshal(bcrBytes, &bcrJSON)
		if err != nil {
			// TODO warn and continue?
			return nil, fmt.Errorf("unmarshalling %s: %v", bcrPath, err)
		}

		meta4Path := path.Join(path.Dir(bcrPath), "artifact.meta4")

		meta4Bytes, err := ioutil.ReadFile(meta4Path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", meta4Path, err)
		}

		var meta4 metalink.Metalink

		err = metalink.Unmarshal(meta4Bytes, &meta4)
		if err != nil {
			// TODO warn and continue?
			return nil, fmt.Errorf("unmarshalling %s: %v", meta4Path, err)
		}

		inmemory = append(
			inmemory,
			compiledreleaseversion.New(
				releaseversion.Reference{
					Name:      bcrJSON.Release.Name,
					Version:   bcrJSON.Release.Version,
					Checksums: bcrJSON.Release.Checksums,
				},
				osversion.Reference{
					Name:    bcrJSON.OS.Name,
					Version: bcrJSON.OS.Version,
				},
				meta4.Files[0],
				map[string]interface{}{
					"uri": fmt.Sprintf(
						"%s//%s",
						i.config.RepositoryConfig.Repository,
						strings.TrimPrefix(path.Dir(strings.TrimPrefix(meta4Path, i.config.RepositoryConfig.LocalPath)), "/"),
					),
					"version": meta4.Files[0].Version,
					"options": map[string]interface{}{
						"private_key": "((index_private_key))",
					},
				},
			),
		)
	}

	i.logger.Infof("found %d entries", len(inmemory))

	return inmemory, nil
}
