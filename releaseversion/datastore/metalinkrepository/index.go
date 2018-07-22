package metalinkrepository

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"

	"github.com/dpb587/boshua/datastore/git"
	"github.com/dpb587/boshua/releaseversion"
	"github.com/dpb587/boshua/releaseversion/datastore"
	"github.com/dpb587/metalink"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Index struct {
	logger     logrus.FieldLogger
	config     Config
	repository *git.Repository
}

var _ datastore.Index = &Index{}

func New(config Config, logger logrus.FieldLogger) *Index {
	return &Index{
		logger:     logger.WithField("build.package", reflect.TypeOf(Index{}).PkgPath()),
		config:     config,
		repository: git.NewRepository(logger, config.RepositoryConfig),
	}
}

func (i *Index) Filter(f *datastore.FilterParams) ([]releaseversion.Artifact, error) {
	if !f.NameSatisfied(i.config.Release) {
		return nil, nil
	}

	_, err := i.repository.Reload()
	if err != nil {
		return nil, errors.Wrap(err, "reloading repository")
	}

	paths, err := filepath.Glob(filepath.Join(i.config.RepositoryConfig.LocalPath, i.config.Prefix, "*.meta4"))
	if err != nil {
		return nil, errors.Wrap(err, "globbing")
	}

	var results = []releaseversion.Artifact{}

	for _, releasePath := range paths {
		releaseBytes, err := ioutil.ReadFile(releasePath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", releasePath, err)
		}

		var releaseMeta4 metalink.Metalink

		err = metalink.Unmarshal(releaseBytes, &releaseMeta4)
		if err != nil {
			// TODO warn and continue?
			return nil, fmt.Errorf("unmarshalling %s: %v", releasePath, err)
		}

		err = metalink.Unmarshal(releaseBytes, &releaseMeta4)
		if err != nil {
			// TODO warn and continue?
			return nil, fmt.Errorf("unmarshalling %s: %v", releasePath, err)
		}

		if !f.VersionSatisfied(releaseMeta4.Files[0].Version) {
			continue
		} else if !f.ChecksumSatisfied(releaseMeta4.Files[0].Hashes) {
			continue
		}

		// TODO sanity checks? version match? files = 1?
		results = append(results, releaseversion.Artifact{
			Name:          i.config.Release,
			Version:       releaseMeta4.Files[0].Version,
			SourceTarball: releaseMeta4.Files[0],
		})
	}

	return results, nil
}
